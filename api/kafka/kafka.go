package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

var Producer sarama.SyncProducer
var Consumer sarama.Consumer
var broker = "kafka:9093"

// InitializeKafka initializes the Kafka producer
func InitializeKafka() error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	var err error
	Producer, err = sarama.NewSyncProducer([]string{broker}, config)
	if err != nil {
		log.Fatal("Failed to start Kafka producer: ", err)
		return err
	}

	Consumer, err = sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		log.Fatal("Failed to start Kafka consumer: ", err)
		return err
	}
	return nil
}

func CloseKafka() {
	if err := Producer.Close(); err != nil {
		log.Fatal("Failed to shut down Kafka producer: ", err)
	}
}
