package main

import (
	myTailf "logAgent/tailf"
	"sync"
)

type LineContent struct {
	Topic string
	Text  string
}

// 日志发往kafka的缓冲区
var sendChan = make(chan LineContent, 100)

func run(task *myTailf.TailTask, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case line := <-task.ReadChan():
			sendChan <- LineContent{
				Topic: task.GetTopic(),
				Text:  line.Text,
			}
			DebugLog("topic:%s,line:=%s\n", task.GetTopic(), line.Text)
		}
	}
}

func SendToKafka() {
	for {
		select {
		case line := <-sendChan:
			err := kafkaProducerObj.SendMsgKafka(line.Topic, line.Text)
			if err != nil {
				DebugLog("SendMsgKafka failed ,err:", err)
			}
		}
	}
}
