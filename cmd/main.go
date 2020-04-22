package main

import (
	"github.com/Shopify/sarama"
	"github.com/asavt92/kafka-stub/internal/configs"
	"github.com/asavt92/kafka-stub/internal/json_service"
	"github.com/asavt92/kafka-stub/internal/kafka"
	"github.com/asavt92/kafka-stub/internal/processors"
)

var (
	kafkaConfig         *configs.KafkaConfig
	jsonConfig          *configs.JsonConfig
	jsonService         *json_service.JsonService
	consumer            kafka.Consumer
	producer            kafka.Producer
	fromConsumerChannel = make(chan sarama.ConsumerMessage, 100)
	toProducerChannel   = make(chan sarama.ProducerMessage, 100)
	mainProcService     *processors.MainService
)

func initApp() {

	configs.InitConfig()
	kafkaConfig = configs.InitKafkaConfig()
	jsonConfig = configs.InitJsonConfig()
	jsonService = json_service.NewJsonService(jsonConfig)

	mainProcService = processors.NewMainService(jsonService, fromConsumerChannel, toProducerChannel)

	producer = kafka.NewProducer(kafkaConfig, toProducerChannel)
	consumer = kafka.NewConsumer(kafkaConfig, fromConsumerChannel)

}

func runApp() {

	go mainProcService.Start()
	go producer.Start()
	consumer.Start()

}

func main() {
	initApp()
	runApp()
}
