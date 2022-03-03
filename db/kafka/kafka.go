package kafka

import (
	"fmt"
	"os"

	"github.com/liumingmin/goutils/conf"

	"github.com/Shopify/sarama"
)

var producers = make(map[string]Producer)
var consumers = make(map[string]Consumer)

func GetProducer(key string) Producer {
	if rt, ok := producers[key]; ok {
		return rt
	}
	return nil
}

func GetConsumer(key string) Consumer {
	if rt, ok := consumers[key]; ok {
		return rt
	}
	return nil
}

func InitKafka() {
	initProducers()
	initConsumers()
}

type ProducerMessageHandler func(msg *sarama.ProducerMessage)
type ProducerErrorHandler func(err *sarama.ProducerError)

type Producer interface {
	Close() error
	Produce(msgs ...*sarama.ProducerMessage) error
	AsyncHandle(mh ProducerMessageHandler, eh ProducerErrorHandler) // 必须设置 asyncReturnSuccess 或 asyncReturnError
}

type ConsumerMessageHandler func(msg *sarama.ConsumerMessage) error
type ConsumerErrorHandler func(err error)

type Consumer interface {
	Close() error
	Consume(topics string, mh ConsumerMessageHandler, eh ConsumerErrorHandler) error
	ConsumeM(topics []string, mh ConsumerMessageHandler, eh ConsumerErrorHandler) error
}

func initProducers() {
	producerConfs := conf.Conf.KafkaProducers
	if producerConfs == nil {
		fmt.Fprintf(os.Stderr, "No producers configuration")
		return
	}

	for _, producerConf := range producerConfs {
		var p Producer
		var err error
		if producerConf.Async {
			p, err = newSaramaAsyncProducer(producerConf)
		} else {
			p, err = newSaramaSyncProducer(producerConf)
		}
		if err != nil {
			continue
		}

		producers[producerConf.Key] = p
	}
}

func initConsumers() {
	consumerConfs := conf.Conf.KafkaConsumers
	if consumerConfs == nil {
		fmt.Fprintf(os.Stderr, "No consumer configuration")
		return
	}

	for _, consumerConf := range consumerConfs {
		c, err := newSaramaConsumerGroup(consumerConf)
		if err != nil {
			continue
		}

		consumers[consumerConf.Key] = c
	}
}
