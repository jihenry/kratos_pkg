package kafka

import (
	"context"
	"fmt"
	"time"

	"gitlab.yeahka.com/gaas/pkg/zaplog"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

func InitConsumer(ctx context.Context, kc KafkaConf, funHandle func(msg *sarama.ConsumerMessage) error, mode bool) *cluster.Consumer {
	conf := cluster.NewConfig()
	//开启错误
	conf.Consumer.Return.Errors = true
	//分组通知
	conf.Group.Return.Notifications = true
	conf.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	conf.Consumer.Offsets.CommitInterval = 1 * time.Second
	conf.Consumer.Offsets.Initial = sarama.OffsetNewest
	if mode {
		conf.Group.Mode = cluster.ConsumerModePartitions
	}
	var (
		kafkaConsumer *cluster.Consumer
		err           error
	)
	if kafkaConsumer, err = cluster.NewConsumer(kc.Addr, kc.Group, []string{kc.Topic}, conf); err != nil || kafkaConsumer == nil {
		if err != nil {
			panic(err.Error())
		} else {
			panic(fmt.Sprintf("consumer is nil kafka info: config:%+v", kc))
		}
	}
	go receiveMessage(ctx, kafkaConsumer, funHandle)

	return kafkaConsumer
}

func receiveMessage(ctx context.Context, kCon *cluster.Consumer, funHandle func(pm *sarama.ConsumerMessage) error) {
	for {
		select {
		case part, ok := <-kCon.Partitions():
			if !ok {
				return
			}
			go func(pc cluster.PartitionConsumer) {
				for msg := range pc.Messages() {
					if err := funHandle(msg); err != nil {
						continue
					} else {
						kCon.MarkOffset(msg, "")
					}
				}
			}(part)
		case msg, ok := <-kCon.Messages():
			if ok {
				if err := funHandle(msg); err != nil {
					continue
				} else {
					kCon.MarkOffset(msg, "")
				}
			} else {
				time.Sleep(1 * time.Second)
			}
		case err, ok := <-kCon.Errors():
			if ok {
				zaplog.Loggers().Infof("consumer error: %v\n", err)
			}
		case ntf, ok := <-kCon.Notifications():
			if ok {
				zaplog.Loggers().Infof("consumer notification: %v\n", ntf)
			}
		case <-ctx.Done():
			zaplog.Loggers().Info("receive end signal")
			return
		}
	}
}
