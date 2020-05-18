package processors

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/asavt92/kafka-stub/internal/json_service"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type MainService struct {
	JsonService   json_service.JsonService
	InputChannel  chan sarama.ConsumerMessage
	OutputChannel chan sarama.ProducerMessage
}

func NewMainService(JsonService *json_service.JsonService, InputChannel chan sarama.ConsumerMessage, OutputChannel chan sarama.ProducerMessage) (m *MainService) {
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

				var raw map[string]interface{}
				if err := json.Unmarshal([]byte(msg.Value), &raw); err != nil {
					panic(err)
				}

				//value := raw[reflect.ValueOf(m.JsonService.JsonConfig.IdFieldMapping).MapKeys()[0].String()]
				value := raw[m.JsonService.JsonConfig.IdFieldMapping["from"]]

				res, _ := m.JsonService.GetJsonResponseStringByConfiguredFieldMapping(fmt.Sprintf("%v", value))

				if res == nil {
					log.Infof("Not found examples by field mapping %v", m.JsonService.JsonConfig.IdFieldMapping)
					res = m.JsonService.GetRandomJsonResponseString()
				} else {
					log.Infof("Found example by field mapping %v", m.JsonService.JsonConfig.IdFieldMapping)
				}

				if res == nil {
					log.Info("Not found response examples, produce empty msg")
					res = make(map[string]interface{})
				}

				out, _ := json.Marshal(res)

				outMsg := sarama.ProducerMessage{
					Key:     sarama.StringEncoder(strconv.Itoa(int(time.Now().Unix()))),
					Value:   sarama.StringEncoder(string(out)),
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
