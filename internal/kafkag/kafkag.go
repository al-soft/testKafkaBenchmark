package kafkag

import (
	"bcm-analyzer/internal/analyzer"
	"context"
	"encoding/json"
	"strings"
	"sync"

	kafka "github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
	yaml "github.com/tsuru/config"
)

type TKafka struct {
	Server  string
	Topic   string
	GroupID string
}

type tHealthchecksData struct {
	Data      []map[string]interface{} `json:"data"`
	BchainUID string                   `json:"bchain_uid"`
}

var (
	Kafka  TKafka
	reader *kafka.Reader
	writer kafka.Writer
	wg     sync.WaitGroup
)

// New инициализация параметров сервера Kafka
func New(configFile *string) (*kafka.Reader, *kafka.Writer) {
	yaml.ReadConfigFile(*configFile)

	k := TKafka{}
	k.Server, _ = yaml.GetString("kafka:server")
	k.Topic, _ = yaml.GetString("kafka:topic:bcmseeker")
	k.GroupID, _ = yaml.GetString("kafka:group_id")

	brokers := strings.Split(k.Server, ",")
	reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  k.GroupID,
		Topic:    k.Topic,
		MinBytes: 10,   // 10B
		MaxBytes: 10e6, // 10MB
	})

	k.Topic, _ = yaml.GetString("kafka:topic:bcmrecorder")

	writer = *kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   k.Topic,
	})

	log.Warnf("kafka reader broker: %s, topic: %s, groip-id: %s", reader.Config().Brokers, reader.Config().Topic, reader.Config().GroupID)
	log.Warnf("kafka writer server: %s, topic: %s", writer.Addr, writer.Topic)

	return reader, &writer
}

func ReadKafkaTopic() error {
	channel := make(chan kafka.Message)
	ReadMessage(channel)

	go func() {
		for {
			topicData := <-channel

			var healthchecks tHealthchecksData
			if err := json.Unmarshal(topicData.Value, &healthchecks); err != nil {
				log.Warnf("kafka ошибка десериализации прочитанного сообщения из топика %s error: %s", reader.Config().Topic, err.Error())
			}
			log.Warnf("bchain_uid: %s kafka прочитал сообщение(healthcheck) из топика %s", healthchecks.BchainUID, reader.Config().Topic)

			bchain, err := analyzer.BusinessChainMonitoring(healthchecks.Data, healthchecks.BchainUID)
			if err != nil {
				log.Warnf("bchain_uid: %s ошибка получения статуса %s", healthchecks.BchainUID, err.Error())
			}

			if err := putMessage2Topic(bchain, healthchecks.BchainUID); err != nil {
				log.Warnf("bchain_uid: %s ошибка записи to kafka topic %s", healthchecks.BchainUID, err.Error())
			}
		}
	}()

	return nil
}

func ReadMessage(channel chan kafka.Message) error {
	go func() {
		for {
			topicData, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Warnf("kafka ошибка чтения из топика %s error: %s", reader.Config().Topic, err.Error())
			}
			channel <- topicData
		}
	}()

	return nil
}

// PutMessage2Topic записывает message в топик Kafka
func putMessage2Topic(message []byte, bchainID string) error {
	wg.Add(1)

	go func() {
		msg := kafka.Message{Value: message}
		if err := writer.WriteMessages(context.Background(), msg); err != nil {
			log.Warnf("bchain_id: %s kafka ошибка записи сообщения в топик %s error: %s", bchainID, writer.Topic, err.Error())
		} else {
			log.Warnf("bchain_id: %s kafka данные записаны в топик %s для обработки bcm-recorder", bchainID, writer.Topic)
		}
		wg.Done()
	}()

	wg.Wait()

	return nil
}
