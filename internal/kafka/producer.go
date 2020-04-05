package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/asavt92/kafka-stub/internal/configs"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"time"
)

type Producer struct {
	KafkaBrokers    []string
	ConsumerGroupID string
	KafkaTopics     []string
	InChannel       chan sarama.ProducerMessage
}

func NewProducer(config *configs.KafkaConfig, inChannel chan sarama.ProducerMessage) (c Producer) {
	return Producer{
		KafkaTopics:     config.OutTopic,
		ConsumerGroupID: config.Group,
		KafkaBrokers:    config.BootstrapServers,
		InChannel:       inChannel,
	}
}


func (c *Producer) Start() {
	// Setup configuration
	config := sarama.NewConfig()
	// Return specifies what channels will be populated.
	// If they are set to true, you must read from
	// config.Producer.Return.Successes = true
	// The total number of times to retry sending a message (default 3).
	config.Producer.Retry.Max = 5
	// The level of acknowledgement reliability needed from the broker.
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Flush.Frequency= time.Millisecond*400
	config.Version = sarama.V2_1_0_0

	producer, err := sarama.NewAsyncProducer(c.KafkaBrokers, config)
	if err != nil {
		// Should not reach here
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			// Should not reach here
			panic(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var enqueued, errors int
	doneCh := make(chan struct{})
	go func() {
		for {
			msg := <- c.InChannel
			msg.Topic=c.KafkaTopics[0]

			select {
			case producer.Input() <- &msg:
				enqueued++
				log.Info("Produce message: ", msg.Key, msg.Value)
			case err := <-producer.Errors():
				errors++
				log.Errorf("Failed to produce message: %v", err)
			case <-signals:
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	log.Printf("Enqueued: %d; errors: %d\n", enqueued, errors)


}