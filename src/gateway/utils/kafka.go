package utils

import (
	"encoding/json"
	"gateway/objects"
	"log"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
	"github.com/urfave/negroni"
)

func sendRequestStatToKafka(stat *objects.RequestStat, topic string, producer sarama.SyncProducer) {
	statBytes, err := json.Marshal(stat)
	if err != nil {
		log.Printf("Error encoding request stats")
		return
	}

	// Создаем сообщение Kafka
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(statBytes),
	}

	// Отправляем сообщение в Kafka
	_, _, err = producer.SendMessage(msg)
	if err != nil {
		log.Printf("Error sending request stat to Kafka: %v", err)
		return
	}

	log.Printf("Request stat sent to Kafka: %s", string(statBytes))
}

// Обертка для обработчиков HTTP, чтобы сохранять статистику запросов
func RequestStatMiddleware(next http.Handler, topic string, producer sarama.SyncProducer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := negroni.NewResponseWriter(w)

		stat := &objects.RequestStat{}
		stat.Path = r.URL.Path
		stat.Method = r.Method
		stat.UserName = r.Header.Get("X-User-Name")
		stat.StartedAt = time.Now()

		next.ServeHTTP(lrw, r)

		stat.FinishedAt = time.Now()
		stat.Duration = stat.FinishedAt.Sub(stat.StartedAt)
		stat.ResponceCode = lrw.Status()

		go sendRequestStatToKafka(stat, topic, producer)
	})
}
