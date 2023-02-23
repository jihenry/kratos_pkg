package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
)

var (
	pvProducer KafkaProducer
	pv         sync.Mutex
)

func SetGlobalProducer(producer KafkaProducer) {
	pv.Lock()
	pvProducer = producer
	pv.Unlock()
}

func GetGlobalProdcuer() KafkaProducer {
	return pvProducer
}

func SendSyncMsg(ctx context.Context, msgs ...*sarama.ProducerMessage) error {
	if pvProducer == nil {
		return fmt.Errorf("global kafka producer is nil")
	}
	return pvProducer.SendSyncMsg(ctx, msgs...)
}

func SendAsyncMsg(ctx context.Context, msgs ...*sarama.ProducerMessage) error {
	if pvProducer == nil {
		return fmt.Errorf("global kafka producer is nil")
	}
	return pvProducer.SendAsyncMsg(ctx, msgs...)
}
