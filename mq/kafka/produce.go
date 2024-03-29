package kafka

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"

	"github.com/Shopify/sarama"
)

type KafkaProduceOption func(*produceOption)

type produceOption struct {
	asyncMsgResultFunc func(err *sarama.ProducerError)
	name               string
}

func WithAsyncMsgResult(resultFunc func(err *sarama.ProducerError)) KafkaProduceOption {
	return func(po *produceOption) {
		po.asyncMsgResultFunc = resultFunc
	}
}

func WithProducerName(name string) KafkaProduceOption {
	return func(po *produceOption) {
		po.name = name
	}
}

var _ KafkaProducer = (*kafkaProducerImpl)(nil)

type kafkaProducerImpl struct {
	asyncProducer  sarama.AsyncProducer
	syncProducer   sarama.SyncProducer
	stopCanlelFunc func()
	logger         *log.Helper
	options        produceOption
}

type KafkaProducer interface {
	SendSyncMsg(ctx context.Context, msg ...*sarama.ProducerMessage) error  //同步发送消息
	SendAsyncMsg(ctx context.Context, msg ...*sarama.ProducerMessage) error //异步发送消息
	Stop() error
}

func NewKafkaProducer(addrs []string, opts ...KafkaProduceOption) (KafkaProducer, error) {
	pConf := sarama.NewConfig()
	pConf.Producer.RequiredAcks = sarama.WaitForAll
	pConf.Producer.Return.Successes = true
	pConf.Producer.Partitioner = sarama.NewRandomPartitioner
	sclient, err := sarama.NewClient(addrs, pConf)
	if err != nil {
		return nil, err
	}
	asyncProducer, err := sarama.NewAsyncProducerFromClient(sclient)
	if err != nil {
		return nil, err
	}
	syncProducer, err := sarama.NewSyncProducerFromClient(sclient)
	if err != nil {
		return nil, err
	}
	options := produceOption{
		name: "kafka",
	}
	for _, opt := range opts {
		opt(&options)
	}
	impl := &kafkaProducerImpl{
		asyncProducer: asyncProducer,
		syncProducer:  syncProducer,
		logger:        log.NewHelper(log.GetLogger(), log.WithMessageKey(options.name)),
	}
	stopCtx, cancelFunc := context.WithCancel(context.Background())
	impl.stopCanlelFunc = cancelFunc
	impl.options = options
	go impl.receiveAsyncMsg(stopCtx)
	return impl, nil
}

func (p *kafkaProducerImpl) SendSyncMsg(ctx context.Context, msgs ...*sarama.ProducerMessage) error {
	if len(msgs) == 1 {
		partition, offset, err := p.syncProducer.SendMessage(msgs[0])
		p.logger.WithContext(ctx).Infof("SendSyncMsg one partition:%d offset:%d err:%v", partition, offset, err)
		return err
	} else {
		err := p.syncProducer.SendMessages(msgs)
		p.logger.WithContext(ctx).Infof("SendSyncMsg list err:%v", err)
		return err
	}
}

func (p *kafkaProducerImpl) SendAsyncMsg(ctx context.Context, msgs ...*sarama.ProducerMessage) error {
	for _, msg := range msgs {
		p.asyncProducer.Input() <- msg
	}
	return nil
}

func (p *kafkaProducerImpl) Stop() error {
	if p.stopCanlelFunc != nil {
		p.stopCanlelFunc()
	}
	if p.asyncProducer != nil {
		if err := p.asyncProducer.Close(); err != nil {
			return err
		}
	}
	if p.syncProducer != nil {
		if err := p.syncProducer.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (p *kafkaProducerImpl) receiveAsyncMsg(stopCtx context.Context) {
	p.logger.Infof("kafka producer receiveAsyncMsg started")
	for {
		select {
		case succMsg, ok := <-p.asyncProducer.Successes():
			if !ok {
				time.Sleep(1 * time.Second)
			} else if p.options.asyncMsgResultFunc != nil {
				p.options.asyncMsgResultFunc(&sarama.ProducerError{
					Msg: succMsg,
					Err: nil,
				})
			} else {
				p.logger.Debugf("kafka producer receiveAsyncMsg success topic:%s key:%s", succMsg.Topic, succMsg.Key)
			}
		case errMsg, ok := <-p.asyncProducer.Errors():
			p.logger.Debugf("kafka producer receiveAsyncMsg err msg:%s", errMsg.Err)
			if !ok {
				time.Sleep(1 * time.Second)
			} else if p.options.asyncMsgResultFunc != nil {
				p.options.asyncMsgResultFunc(errMsg)
			}
		case <-stopCtx.Done():
			p.logger.Infof("kafka producer receiveAsyncMsg stopped")
			return
		}
	}
}
