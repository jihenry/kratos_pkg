package kafka

import (
	"context"

	"gitlab.yeahka.com/gaas/pkg/log"

	"github.com/Shopify/sarama"
)

type (
	KafkaConf struct {
		Name  string   `yaml:"Name"`
		Addr  []string `yaml:"Addr"`
		Topic string   `yaml:"Topic"`
		Group string   `yaml:"Group"`
	}
	ProduceCli interface {
		Sync(context.Context, *sarama.ProducerMessage) error
		Syncs(context.Context, []*sarama.ProducerMessage) error
		Async(msg *sarama.ProducerMessage)
	}
)

var (
	_ ProduceCli = (*produceClient)(nil)
)

type produceClient struct {
	asyncProduce sarama.AsyncProducer
	syncProduce  sarama.SyncProducer
}

func InitProduce(kc KafkaConf, funHandle func(anyErr *sarama.ProducerError)) ProduceCli {
	produceCli := new(produceClient)
	pConf := sarama.NewConfig()
	pConf.Producer.RequiredAcks = sarama.WaitForAll
	pConf.Producer.Return.Successes = true
	pConf.Producer.Partitioner = sarama.NewRandomPartitioner
	var (
		oCli  sarama.Client
		async sarama.AsyncProducer
		sync  sarama.SyncProducer
		err   error
	)
	if oCli, err = sarama.NewClient(kc.Addr, pConf); err != nil {
		panic(err)
	}
	//异步模式
	if async, err = sarama.NewAsyncProducerFromClient(oCli); err != nil {
		panic(err)
	}
	//同步模式
	if sync, err = sarama.NewSyncProducerFromClient(oCli); err != nil {
		panic(err)
	}
	go receiveErr(async, funHandle)
	go receiveSuccess(async, funHandle)
	produceCli.asyncProduce = async
	produceCli.syncProduce = sync
	return produceCli
}

// Sync 同步发送消息
func (p *produceClient) Sync(ctx context.Context, msg *sarama.ProducerMessage) error {
	partition, offset, err := p.syncProduce.SendMessage(msg)
	log.FromContext(ctx).Infof("[KafkaSendMessage] partition:%d offset:%d err:%v", partition, offset, err)
	return err
}

// Syncs 批量同步发送消息
func (p *produceClient) Syncs(ctx context.Context, msgs []*sarama.ProducerMessage) error {
	return p.syncProduce.SendMessages(msgs)
}

// Async 异步发送消息
func (p *produceClient) Async(msg *sarama.ProducerMessage) {
	p.asyncProduce.Input() <- msg
}

func receiveErr(asyncP sarama.AsyncProducer, funHandle func(err *sarama.ProducerError)) {
	for err := range asyncP.Errors() {
		funHandle(err)
	}
}

func receiveSuccess(asyncProducer sarama.AsyncProducer, fnHandle func(err *sarama.ProducerError)) {
	for msg := range asyncProducer.Successes() {
		err := new(sarama.ProducerError)
		err.Msg = msg
		err.Err = nil
		fnHandle(err)
	}
}
