package kafka

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/IBM/sarama"
)

const userTopic = "user-topic"

func TestKafkaProducer(t *testing.T) {
	InitKafka()
	producer := GetProducer("user_producer")
	err := producer.Produce(&sarama.ProducerMessage{
		Topic: userTopic,
		Key:   sarama.ByteEncoder(fmt.Sprint(time.Now().Unix())),
		Value: sarama.ByteEncoder(fmt.Sprint(time.Now().Unix())),
	})

	if err != nil {
		t.Error(err)
	}

	time.Sleep(time.Millisecond * 100)
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
		err := producer.Produce(&sarama.ProducerMessage{
			Topic: userTopic,
			Key:   sarama.ByteEncoder(fmt.Sprint(i)),
			Value: sarama.ByteEncoder(fmt.Sprint(time.Now().Unix())),
		})
		if err != nil {
			t.Error(err)
		}
	}

	time.Sleep(time.Millisecond * 100)
}

func TestMain(m *testing.M) {
	conn, err := net.DialTimeout("tcp", "127.0.0.1:9092", time.Second*2)
	if err != nil {
		fmt.Println("Please install kafka on local and start at port: 9092, then run test.")
		return
	}
	conn.Close()

	m.Run()
}
