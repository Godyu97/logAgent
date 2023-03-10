package main

import (
	myTailf "logAgent/tailf"
)

type LineContent struct {
	Topic string
	Text  string
}

// 日志发往kafka的缓冲区
var sendChan = make(chan LineContent, 100)

func run(task *myTailf.TailTask) {
	for {
		select {
		case <-task.CtxDone():
			DebugLog("取消了：", task.GetTopic()+task.GetFilePath())
			return
		case text := <-task.ReadChan():
			sendChan <- LineContent{
				Topic: task.GetTopic(),
				Text:  text,
			}
			DebugLog("向缓冲区投递 topic:%s,line:=%s\n", task.GetTopic(), text)
		}
	}
}
func NewRunTask(ev LogEntry) *myTailf.TailTask {
	t, err := myTailf.NewTailTask(ev.FilePath, ev.Topic)
	if err != nil {
		DebugLog("NewTailTask failed ,err:", err)
	}
	go run(t)
	return t
}

func SendToKafka() {
	for {
		select {
		case line := <-sendChan:
			DebugLog("从缓冲区取log，SendMsgKafka：", line.Topic, line.Text)
			err := kafkaProducerObj.SendMsgKafka(line.Topic, line.Text)
			if err != nil {
				DebugLog("SendMsgKafka failed ,err:", err)
			}
		}
	}
}
