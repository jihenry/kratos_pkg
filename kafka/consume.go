package kafka

import (
	"context"
	"time"

	"gitlab.yeahka.com/gaas/pkg/log"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

type KafkaOption func(*options)
type options struct {
	addrs  []string
	mode   cluster.ConsumerMode
	topics []string
	group  string
}

type KafkaConsumer interface {
	Receive(ctx context.Context, funHandle func(pm *sarama.ConsumerMessage) error)
}

type kafkaConsumerImpl struct {
	consumer *cluster.Consumer
}

func WithMode(mode cluster.ConsumerMode) KafkaOption {
	return func(o *options) {
		o.mode = mode
	}
}

func NewKafkaConsumer(ctx context.Context, addrs []string, topics []string, group string, opts ...KafkaOption) (KafkaConsumer, error) {
	conf := cluster.NewConfig()
	conf.Consumer.Return.Errors = true
	conf.Group.Return.Notifications = true
	conf.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	conf.Consumer.Offsets.CommitInterval = 1 * time.Second
	conf.Consumer.Offsets.Initial = sarama.OffsetNewest
	options := options{
		mode:   cluster.ConsumerModeMultiplex,
		addrs:  addrs,
		topics: topics,
		group:  group,
	}
	for _, opt := range opts {
		opt(&options)
	}
	consumer, err := cluster.NewConsumer(options.addrs, options.group, options.topics, conf)
	if err != nil {
		return nil, err
	}
	return &kafkaConsumerImpl{
		consumer: consumer,
	}, nil
}

func (c *kafkaConsumerImpl) Receive(ctx context.Context, funHandle func(pm *sarama.ConsumerMessage) error) {
	for {
		select {
		case part, ok := <-c.consumer.Partitions():
			if !ok {
				return
			}
			go func(pc cluster.PartitionConsumer) {
				for msg := range pc.Messages() {
					if err := funHandle(msg); err != nil {
						continue
					} else {
						c.consumer.MarkOffset(msg, "")
					}
				}
			}(part)
		case msg, ok := <-c.consumer.Messages():
			if ok {
				if err := funHandle(msg); err != nil {
					continue
				} else {
					c.consumer.MarkOffset(msg, "")
				}
			} else {
				time.Sleep(1 * time.Second)
			}
		case err, ok := <-c.consumer.Errors():
			if ok {
				log.Loggers().Infof("consumer error: %v\n", err)
			}
		case ntf, ok := <-c.consumer.Notifications():
			if ok {
				log.Loggers().Infof("consumer notification: %v\n", ntf)
			}
		case <-ctx.Done():
			log.Loggers().Info("receive end signal")
			return
		}
	}
}
