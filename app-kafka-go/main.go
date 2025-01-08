package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/IBM/sarama"
)

func main() {
	brokers := []string{"127.0.0.1:9092"}
	topic := "example-topic"

	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		log.Fatalf("Error creating cluster admin: %v", err)
	}
	defer admin.Close()

	topics, err := admin.ListTopics()
	if err != nil {
		log.Fatalf("Error listing topics: %v", err)
	}

	topicName := "example-topic"
	if _, exists := topics[topicName]; exists {
		fmt.Printf("Topic %s exists\n", topicName)
	} else {
		fmt.Printf("Topic %s does not exist\n", topicName)
	}

	// Start producer
	go func() {
		producer, err := sarama.NewSyncProducer(brokers, nil)
		if err != nil {
			log.Fatalln("Failed to start Sarama producer:", err)
		}
		defer producer.Close()

		for {
			msg := &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.StringEncoder("Hello Kafka!"),
			}
			_, _, err := producer.SendMessage(msg)
			if err != nil {
				log.Println("Failed to send message:", err)
			} else {
				log.Println("Message sent")
			}
			time.Sleep(2 * time.Second)
		}
	}()

	// Start consumer
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		log.Fatalln("Failed to start Sarama consumer:", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalln("Failed to start partition consumer:", err)
	}
	defer partitionConsumer.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				fmt.Printf("Received message: %s\n", string(msg.Value))
			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()

	<-doneCh
}