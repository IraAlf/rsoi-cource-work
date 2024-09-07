package utils

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"statistics/objects"

	"github.com/Shopify/sarama"
)

var messages = make(chan *objects.RequestStat)

func GetMessage() *objects.RequestStat {
	return <-messages
}

type TaskHandler struct{}

func (*TaskHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (*TaskHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (*TaskHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
			session.MarkMessage(message, "")

			dto := new(objects.RequestStat)
			json.Unmarshal(message.Value, dto)
			messages <- dto

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/Shopify/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}

type KafkaSettings struct {
	Consumer sarama.ConsumerGroup
}

func InitKafka() *KafkaSettings {
	kafkaBrokers := Config.Kafka.Endpoints
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	config := sarama.NewConfig()
	config.Net.TLS.Enable = false
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	consumer, err := sarama.NewConsumerGroup(kafkaBrokers, "1", config)
	if err != nil {
		log.Printf("Error creating Kafka consumer: %v", err)
	}

	return &KafkaSettings{
		Consumer: consumer,
	}
}

func (kafka *KafkaSettings) ConsumeLoop() {
	ctx := context.Background()
	handler := &TaskHandler{}
	for {
		err := kafka.Consumer.Consume(ctx, Config.Kafka.Topics, handler)
		if err != nil {
			log.Println(err.Error())
			panic(err)
		}

		if err = ctx.Err(); err != nil {
			log.Panic(err)
		}
	}
}

func (kafka *KafkaSettings) Close() {
	kafka.Consumer.Close()
}
