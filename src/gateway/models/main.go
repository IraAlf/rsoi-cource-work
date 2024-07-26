package models

import (
	"gateway/utils"
	"log"
	"net/http"
	"os"

	"github.com/Shopify/sarama"
)

type KafkaSettings struct {
	Topic    string
	Producer sarama.SyncProducer
}

type Models struct {
	Client     *http.Client
	Flights    *FlightsM
	Privileges *PrivilegesM
	Tickets    *TicketsM
	Statistics *StatisticsM

	Kafka *KafkaSettings
}

func InitKafka() *KafkaSettings {
	kafkaBrokers := utils.Config.Kafka.Endpoints
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	config := sarama.NewConfig()
	config.Net.TLS.Enable = false
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(kafkaBrokers, config)
	if err != nil {
		log.Printf("Error creating Kafka producer: %v", err)
	}

	return &KafkaSettings{
		Topic:    utils.Config.Kafka.Topics[0],
		Producer: producer,
	}
}

func InitModels() *Models {
	models := new(Models)
	client := &http.Client{}

	models.Client = client
	models.Flights = NewFlightsM(client)
	models.Privileges = NewPrivilegesM(client)
	models.Tickets = NewTicketsM(client, models.Flights)
	models.Statistics = NewStatisticsM(client)
	models.Kafka = InitKafka()

	return models
}
