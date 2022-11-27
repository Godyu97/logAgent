package myKafka

import (
	"errors"
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

//func InitKafkaProducer(addrs []string) (producer *KafkaProducer, err error) {
//
//	config := sarama.NewConfig()
//	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
//	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
//	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
//
//	// 连接kafka
//  producer = new()
//	producer.syncProducer, err = sarama.NewSyncProducer(addrs, config)
//	if err != nil {
//		fmt.Println("producer closed, err:", err)
//		return nil, err
//	}
//	//defer client.Close()
//	return producer, nil
//}

//func (obj *KafkaProducer) SendMsgKafka(topic string, text string) (err error) {
//	if obj == nil {
//		err = errors.New("*KafkaProducer nil")
//		return err
//	}
//	// 构造一个消息
//	msg := &sarama.ProducerMessage{
//		Topic: topic,
//		Value: sarama.StringEncoder(text),
//	}
//	// 发送消息
//	pid, offset, err := obj.syncProducer.SendMessage(msg)
//	if err != nil {
//		fmt.Println("send msg failed, err:", err)
//		return
//	}
//	fmt.Printf("pid:%v offset:%v\n", pid, offset)
//	return nil
//}
