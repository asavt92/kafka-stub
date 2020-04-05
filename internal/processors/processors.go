package processors

import (
	"github.com/Shopify/sarama"
	"github.com/asavt92/kafka-stub/internal/json"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type MainService struct {
	JsonService   json.JsonService
	InputChannel  chan sarama.ConsumerMessage
	OutputChannel chan sarama.ProducerMessage
}

func NewMainService(JsonService *json.JsonService, InputChannel chan sarama.ConsumerMessage, OutputChannel chan sarama.ProducerMessage) (m *MainService) {
	return &MainService{
		JsonService:   *JsonService,
		InputChannel:  InputChannel,
		OutputChannel: OutputChannel,
	}
}

func (m *MainService) Start() {

	go func() {
		for {
			msg := <-m.InputChannel
			msgJsonString := string(msg.Value)

			ok := m.JsonService.IsJsonStringValid(msgJsonString)
			if ok {

				headers := convertHeaders(msg.Headers)

				outMsg := sarama.ProducerMessage{
					Key:     sarama.StringEncoder(strconv.Itoa(int(time.Now().Unix()))),
					Value:   sarama.StringEncoder(m.JsonService.GetRandomJsonResponseString()),
					Headers: headers,
				}
				m.OutputChannel <- outMsg
			} else {
				log.Errorf("VALIDATION ERROR : %s", msg)
			}

		}

	}()

}

func convertHeaders(headers []*sarama.RecordHeader) []sarama.RecordHeader {
	var out []sarama.RecordHeader
	for i := 0; i < len(headers); i++ {
		if string((*headers[i]).Key) == "kafka_correlationId" {
			out = append(out, *headers[i])
		}
	}
	return out
}
