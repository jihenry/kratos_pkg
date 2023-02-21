package kafka

import (
	"context"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/go-kratos/kratos/v2/log"
)

type KafkaConsumerOption func(*consumerOption)
type consumerOption struct {
	addrs  []string
	mode   cluster.ConsumerMode
	topics []string
	group  string
	name   string
}

type KafkaConsumer interface {
	Receive(ctx context.Context, funHandle func(pm *sarama.ConsumerMessage) error)
	Stop() error
}

type kafkaConsumerImpl struct {
	consumer *cluster.Consumer
	logger   *log.Helper
	options  consumerOption
}

func WithMode(mode cluster.ConsumerMode) KafkaConsumerOption {
	return func(o *consumerOption) {
		o.mode = mode
	}
}

func WithName(name string) KafkaConsumerOption {
	return func(o *consumerOption) {
		o.name = name
	}
}

func NewKafkaConsumer(ctx context.Context, addrs []string, topics []string, group string, opts ...KafkaConsumerOption) (KafkaConsumer, error) {
	conf := cluster.NewConfig()
	conf.Consumer.Return.Errors = true
	conf.Group.Return.Notifications = true
	conf.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	conf.Consumer.Offsets.CommitInterval = 1 * time.Second
	conf.Consumer.Offsets.Initial = sarama.OffsetNewest
	options := consumerOption{
		mode:   cluster.ConsumerModeMultiplex,
		addrs:  addrs,
		topics: topics,
		group:  group,
	}
	options.name = strings.Join(topics, ",") + ":" + group
	for _, opt := range opts {
		opt(&options)
	}
	consumer, err := cluster.NewConsumer(options.addrs, options.group, options.topics, conf)
	if err != nil {
		return nil, err
	}
	return &kafkaConsumerImpl{
		consumer: consumer,
		options:  options,
		logger:   log.NewHelper(log.GetLogger(), log.WithMessageKey(options.name)),
	}, nil
}

func (c *kafkaConsumerImpl) Stop() error {
	if c.consumer != nil {
		return c.consumer.Close()
	}
	return nil
}

func (c *kafkaConsumerImpl) Receive(ctx context.Context, funHandle func(pm *sarama.ConsumerMessage) error) {
	for {
		select {
		case part, ok := <-c.consumer.Partitions():
			if !ok {
				return
			}
			for msg := range part.Messages() {
				if err := funHandle(msg); err != nil {
					continue
				} else {
					c.consumer.MarkOffset(msg, "")
				}
			}
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
				c.logger.Infof("consumer error: %v\n", err)
			}
		case ntf, ok := <-c.consumer.Notifications():
			if ok {
				c.logger.Infof("consumer notification: %v\n", ntf)
			}
		case <-ctx.Done():
			c.logger.Info("consumer receive stopped")
			return
		}
	}
}
