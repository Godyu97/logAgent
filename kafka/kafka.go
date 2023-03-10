package myKafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"strings"
)

type KafkaProducer struct {
	syncProducer *kafka.Producer
}

func InitKafkaProducer(addrs []string) (producer *KafkaProducer, err error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(addrs, ","),
	}
	p, err := kafka.NewProducer(config)
	if err != nil {
		log.Println("kafka.NewProducer failed,err:", err)
		return nil, err
	}
	producer = new(KafkaProducer)
	producer.syncProducer = p
	// Delivery report handler for produced messages
	go func() {
		for e := range producer.syncProducer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					log.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()
	return producer, nil
}

func (obj *KafkaProducer) SendMsgKafka(topic string, text string) (err error) {
	if obj == nil || obj.syncProducer == nil {
		err = errors.New("*KafkaProducer nil,need init ")
		return err
	}

	// Produce messages to topic (asynchronously)
	err = obj.syncProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(text),
	}, nil)
	if err != nil {
		log.Printf("SendMsgKafka failed, err: %v", err)
		return err
	}
	return nil
}

type KafkaConsumer struct {
	*kafka.Consumer
	context.Context
	context.CancelFunc
}

func InitKafkaConsumer(addrs []string) (consumer *KafkaConsumer, err error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(addrs, ","),
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	}
	c, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}
	consumer = new(KafkaConsumer)
	ctx, cfn := context.WithCancel(context.Background())
	consumer.Context = ctx
	consumer.CancelFunc = cfn
	consumer.Consumer = c

	return consumer, nil
}

type Msg struct {
	Topic string
	Text  string
}

func (obj *KafkaConsumer) FromMsgKafka(topic string, textCh chan<- Msg) {
	obj.Subscribe(topic, nil)
	go func() {
		defer obj.Close()
		for {
			select {
			case <-obj.Done():
				log.Println("FromMsgKafka 结束：", obj)
				return
			default:
				msg, err := obj.ReadMessage(-1)
				if err == nil {
					textCh <- Msg{
						Topic: topic,
						Text:  string(msg.Value),
					}
					fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
				} else {
					// The client will automatically try to recover from all errors.
					fmt.Printf("Consumer error: %v (%v)\n", err, msg)
				}
			}
		}
	}()
}
