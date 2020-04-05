package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/asavt92/kafka-stub/internal/configs"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"
)

type Consumer struct {
	KafkaBrokers    []string
	ConsumerGroupID string
	KafkaTopics     []string
	OutChannel      chan sarama.ConsumerMessage
}

func NewConsumer(config *configs.KafkaConfig, outChannel chan sarama.ConsumerMessage) (c Consumer) {
	return Consumer{
		KafkaTopics:     config.InTopic,
		ConsumerGroupID: config.Group,
		KafkaBrokers:    config.BootstrapServers,
		OutChannel:      outChannel,
	}
}

func (c *Consumer) Start() {
	// Init config, specify appropriate version
	config := sarama.NewConfig()
	logger := log.New()
	logger.Out = os.Stderr
	sarama.Logger = logger
	config.Version = sarama.V2_1_0_0
	config.Consumer.Return.Errors = true

	master, err := sarama.NewConsumer(c.KafkaBrokers, config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := master.Close(); err != nil {
			panic(err)
		}
	}()

	consumer, errors := consume(c.KafkaTopics, master)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Count how many message processed
	msgCount := 0

	// Get signnal for finish
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-consumer:
				msgCount++
				log.Info("Received messages", string(msg.Key), string(msg.Value))

				//todo refactor push to processors
				c.OutChannel <- *msg

			case consumerError := <-errors:
				msgCount++
				log.Error("Received consumerError ", string(consumerError.Topic), string(consumerError.Partition), consumerError.Err)
				doneCh <- struct{}{}
			case <-signals:
				log.Fatal("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
	log.Info("Processed", msgCount, "messages")
}

func consume(topics []string, master sarama.Consumer) (chan *sarama.ConsumerMessage, chan *sarama.ConsumerError) {
	consumers := make(chan *sarama.ConsumerMessage)
	errors := make(chan *sarama.ConsumerError)
	for _, topic := range topics {
		if strings.Contains(topic, "__consumer_offsets") {
			continue
		}
		partitions, _ := master.Partitions(topic)
		// this only consumes partition no 1, you would probably want to consume all partitions
		for _, partition := range partitions {

			consumer, err := master.ConsumePartition(topic, partition, sarama.OffsetNewest)
			if nil != err {
				log.Infof("Topic %v Partitions: %v", topic, partitions)
				panic(err)
			}
			log.Info(" Start consuming topic ", topic)
			go func(topic string, consumer sarama.PartitionConsumer) {
				for {
					select {
					case consumerError := <-consumer.Errors():
						errors <- consumerError
						log.Error("consumerError: ", consumerError.Err)

					case msg := <-consumer.Messages():
						consumers <- msg
						//fmt.Println("Got message on topic ", topic, msg.Value)
					}
				}
			}(topic, consumer)

		}

	}

	return consumers, errors
}
