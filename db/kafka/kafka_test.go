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

	time.Sleep(time.Second * 5)
}

func TestKafkaConsumer(t *testing.T) {
	InitKafka()

	consumer := GetConsumer("user_consumer")
	go func() {
		consumer.Consume(userTopic, func(msg *sarama.ConsumerMessage) error {
			fmt.Println(string(msg.Key), "=", string(msg.Value))
			return nil
		}, func(err error) {

		})
	}()

	producer := GetProducer("user_producer")
	for i := 0; i < 10; i++ {
		producer.Produce(&sarama.ProducerMessage{
			Topic: userTopic,
			Key:   sarama.ByteEncoder(fmt.Sprint(i)),
			Value: sarama.ByteEncoder(fmt.Sprint(time.Now().Unix())),
		})
	}

	time.Sleep(time.Second * 5)
}
