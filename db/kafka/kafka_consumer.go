package kafka

import (
	"context"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/IBM/sarama"
	"github.com/liumingmin/goutils/conf"
)

const (
	ACK_BEFORE_AUTO    = 0
	ACK_AFTER_NOERROR  = 1
	ACK_AFTER_NOMATTER = 2
)

type saramaConsumerGroup struct {
	sarama.ConsumerGroup
	option  *conf.KafkaConsumer
	version int32
}

func (g *saramaConsumerGroup) Close() error {
	return g.ConsumerGroup.Close()
}

// blocking to consume the messages
func (g *saramaConsumerGroup) Consume(topic string, mh ConsumerMessageHandler, eh ConsumerErrorHandler) (err error) {
	return g.ConsumeM([]string{topic}, mh, eh)
}

// blocking to consume the messages
func (g *saramaConsumerGroup) ConsumeM(topics []string, mh ConsumerMessageHandler, eh ConsumerErrorHandler) (err error) {
	ver := atomic.AddInt32(&g.version, 1)
	go func() {
		tk := time.NewTicker(time.Second)
		defer tk.Stop()
		for ver == g.version {
			select {
			case e := <-g.Errors():
				eh(e)
			case <-tk.C:
			}
		}
	}()
	for {
		err = g.ConsumerGroup.Consume(context.Background(), topics, newSaramaConsumerGroupHandler(mh, g.option))
		if err != nil {
			return
		}
	}
}

func newSaramaConsumerGroup(opt *conf.KafkaConsumer) (ret *saramaConsumerGroup, err error) {
	grp, err := sarama.NewConsumerGroup(opt.Address, opt.Group, consumerConfig(opt))
	if err != nil {
		return
	}
	ret = &saramaConsumerGroup{
		ConsumerGroup: grp,
		option:        opt,
	}
	return
}

func consumerConfig(opt *conf.KafkaConsumer) (config *sarama.Config) {
	config = sarama.NewConfig()
	config.Version = sarama.V0_10_2_0 // consumer groups require Version to be >= V0_10_2_0
	if opt.User != "" {               // only plain
		config.Net.SASL.Enable = true
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		config.Net.SASL.User = opt.User
		config.Net.SASL.Password = opt.Password
	}
	if opt.KeepAlive > 0 {
		config.Net.KeepAlive = opt.KeepAlive
	}
	if opt.DialTimeout > 0 {
		config.Net.DialTimeout = opt.DialTimeout
	}
	if opt.ReadTimeout > 0 {
		config.Net.ReadTimeout = opt.ReadTimeout
	}
	if opt.WriteTimeout > 0 {
		config.Net.WriteTimeout = opt.WriteTimeout
	}
	if opt.Offset < 0 { // only -1(OffsetNewest), -2(OffsetOldest)
		if opt.Offset == sarama.OffsetNewest {
			config.Consumer.Offsets.Initial = sarama.OffsetNewest
		} else if opt.Offset == sarama.OffsetOldest {
			config.Consumer.Offsets.Initial = sarama.OffsetOldest
		} else {
			panic("invalid initial offset")
		}
	}

	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type saramaConsumerGroupHandler struct {
	option *conf.KafkaConsumer
	hander ConsumerMessageHandler
}

func newSaramaConsumerGroupHandler(mhandler ConsumerMessageHandler, option *conf.KafkaConsumer) *saramaConsumerGroupHandler {
	return &saramaConsumerGroupHandler{
		hander: mhandler,
		option: option,
	}
}

func (h *saramaConsumerGroupHandler) Setup(s sarama.ConsumerGroupSession) error {
	// FIXBUG: sarama 不支持latest, ResetOffset(-1)会导致server端出现"Invalid negative offset"
	if h.option.Offset == 0 || h.option.Offset == sarama.OffsetNewest {
		return nil
	}

	// 如果是OffsetOldest则不变, 否则自动前移, 因为位置从0开始.
	var realOffset int64 = h.option.Offset
	if realOffset > 0 {
		realOffset--
	}
	for t, ps := range s.Claims() {
		for _, p := range ps {
			s.ResetOffset(t, p, realOffset, "")
		}
	}
	return nil
}
func (h *saramaConsumerGroupHandler) Cleanup(s sarama.ConsumerGroupSession) error { return nil }
func (h *saramaConsumerGroupHandler) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) (err error) {
	for msg := range c.Messages() {
		switch h.option.Ack {
		case ACK_BEFORE_AUTO:
			s.MarkMessage(msg, "")
			err = h.hander(msg)
		case ACK_AFTER_NOERROR:
			if err = h.hander(msg); err == nil {
				s.MarkMessage(msg, "")
			}
		case ACK_AFTER_NOMATTER:
			err = h.hander(msg)
			s.MarkMessage(msg, "")
		default:
			panic("invalid ack type: " + strconv.Itoa(h.option.Ack))
		}
	}
	return
}
