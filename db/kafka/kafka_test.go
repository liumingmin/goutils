package kafka

import (
	"fmt"
	"testing"
	"time"

	"github.com/Shopify/sarama"
)

const userTopic = "user-topic"

func TestKafkaProducer(t *testing.T) {
	InitKafka()
	producer := GetProducer("user_producer")
	producer.Produce(&sarama.ProducerMessage{
		Topic: userTopic,
		Key:   sarama.ByteEncoder(fmt.Sprint(time.Now().Unix())),
		Value: sarama.ByteEncoder(fmt.Sprint(time.Now().Unix())),
	})

	time.Sleep(time.Second * 10)
}

func TestKafkaConsumer(t *testing.T) {
	InitKafka()

	consumer := GetConsumer("user_consumer")
	go func() {
		consumer.Consume(userTopic, func(msg *sarama.ConsumerMessage) error {
			fmt.Println(msg.Key, "=", msg.Value)
			return nil
		}, func(err error) {

		})
	}()

	time.Sleep(time.Second * 10)
}
