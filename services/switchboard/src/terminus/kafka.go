package terminus

import (
	"../header"
	"github.com/Shopify/sarama"
	"log"
)

type KafkaConf struct {
	Addresses []string
	Topic     string
}

type Kafka struct {
	conf     KafkaConf
	errors   chan<- error
	producer sarama.AsyncProducer
}

func NewKafka(conf KafkaConf, errors chan<- error) *Kafka {
	producer, err := sarama.NewAsyncProducer(conf.Addresses, sarama.NewConfig())
	if err != nil {
		log.Panicf("Unable to open connection to kafaka %s: %v\n", conf.Addresses, err)
	}
	go func() {
		for err := range producer.Errors() {
			errors <- err
		}
	}()
	return &Kafka{conf, errors, producer}
}

func (terminus *Kafka) Terminate(h header.Header) {
	terminus.producer.Input() <- &sarama.ProducerMessage{
		Topic: terminus.conf.Topic,
		Key:   nil,
		Value: sarama.StringEncoder(h.String()),
	}
}
