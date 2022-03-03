package kafka

import (
	"sync/atomic"
	"time"

	"github.com/liumingmin/goutils/conf"

	"github.com/Shopify/sarama"
)

type saramaAsyncProducer struct {
	sarama.AsyncProducer
	option  *conf.KafkaProducer
	version int32
}

func newSaramaAsyncProducer(opt *conf.KafkaProducer) (ret *saramaAsyncProducer, err error) {
	p, err := sarama.NewAsyncProducer(opt.Address, producerConfig(opt))
	if err != nil {
		return
	}
	ret = &saramaAsyncProducer{
		option:        opt,
		AsyncProducer: p,
	}
	// 默认异步读取,避免阻塞
	if opt.ReturnSuccess {
		go func() {
			tk := time.Tick(time.Second)
			for ret.version == 0 {
				select {
				case <-ret.Successes():
				case <-tk:
				}
			}
		}()
	}
	// 默认异步读取,避免阻塞
	if opt.ReturnError {
		go func() {
			tk := time.Tick(time.Second)
			for ret.version == 0 {
				select {
				case <-ret.Errors():
				case <-tk:
				}
			}
		}()
	}
	return
}

func (p *saramaAsyncProducer) Close() error {
	return p.AsyncProducer.Close()
}

func (p *saramaAsyncProducer) Produce(msgs ...*sarama.ProducerMessage) error {
	for _, m := range msgs {
		p.AsyncProducer.Input() <- m
	}
	return nil
}

func (p *saramaAsyncProducer) AsyncHandle(mh ProducerMessageHandler, eh ProducerErrorHandler) {
	ver := atomic.AddInt32(&p.version, 1)
	if p.option.ReturnSuccess && mh != nil {
		go func() {
			tk := time.Tick(time.Second)
			for p.version == ver {
				mh(<-p.Successes())
				select {
				case m := <-p.Successes():
					mh(m)
				case <-tk:
				}
			}
		}()
	}
	if p.option.ReturnError && eh != nil {
		go func() {
			tk := time.Tick(time.Second)
			for p.version == ver {
				select {
				case e := <-p.Errors():
					eh(e)
				case <-tk:
				}
			}
		}()
	}
}

func producerConfig(opt *conf.KafkaProducer) (config *sarama.Config) {
	config = sarama.NewConfig()
	config.Version = sarama.V0_10_2_0 // consumer groups require Version to be >= V0_10_2_0

	if opt.User != "" { // only plain
		config.Net.SASL.Enable = true
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		config.Net.SASL.User = opt.User
		config.Net.SASL.Password = opt.Password
	}

	config.Producer.Return.Successes = opt.ReturnSuccess
	config.Producer.Return.Errors = opt.ReturnError
	return
}

type saramaSyncProducer struct {
	sarama.SyncProducer
	option *conf.KafkaProducer
}

func (p *saramaSyncProducer) Close() error {
	return p.SyncProducer.Close()
}

func (p *saramaSyncProducer) Produce(msgs ...*sarama.ProducerMessage) error {
	return p.SyncProducer.SendMessages(msgs)
}

func (p *saramaSyncProducer) AsyncHandle(mh ProducerMessageHandler, eh ProducerErrorHandler) {

}

func newSaramaSyncProducer(opt *conf.KafkaProducer) (ret *saramaSyncProducer, err error) {
	// must be set true
	opt.ReturnSuccess = true
	opt.ReturnError = true
	p, err := sarama.NewSyncProducer(opt.Address, producerConfig(opt))
	if err != nil {
		return
	}
	ret = &saramaSyncProducer{
		option:       opt,
		SyncProducer: p,
	}
	return
}
