package myTailf

import (
	"github.com/hpcloud/tail"
	"log"
)

type TailTask struct {
	filePath string
	topic    string
	tailObj  *tail.Tail
}

func NewTailTask(filePath, topic string) (task *TailTask, err error) {
	task = &TailTask{
		filePath: filePath,
		topic:    topic,
	}
	err = task.init(filePath)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (t *TailTask) init(fileName string) (err error) {
	myCfg := tail.Config{
		Location:    &tail.SeekInfo{Offset: 0, Whence: 2},
		ReOpen:      true,
		MustExist:   false,
		Poll:        true,
		Pipe:        false,
		RateLimiter: nil,
		Follow:      true,
		MaxLineSize: 0,
		Logger:      nil,
	}
	t.tailObj, err = tail.TailFile(fileName, myCfg)
	if err != nil {
		log.Println("tail file failed ,err:", err)
		return err
	}
	return nil
}

func (t *TailTask) ReadChan() <-chan *tail.Line {
	return t.tailObj.Lines
}

func (t *TailTask) GetTopic() string {
	return t.topic
}
func (t *TailTask) GetFilePath() string {
	return t.filePath
}
