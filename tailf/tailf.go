package myTailf

import (
	"context"
	"github.com/hpcloud/tail"
	"log"
)

type TailTask struct {
	filePath  string
	topic     string
	tailObj   *tail.Tail
	textChan  chan string
	ctx       context.Context
	ctxCancel context.CancelFunc
}

// 打开文件
func NewTailTask(filePath, topic string) (task *TailTask, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	task = &TailTask{
		filePath:  filePath,
		topic:     topic,
		textChan:  make(chan string, 100),
		ctx:       ctx,
		ctxCancel: cancel,
	}
	err = task.init(filePath)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// 此处只开文件
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

func (t *TailTask) ReadChan() <-chan string {
	go func() {
		for {
			select {
			//防止goroutine泄露
			case <-t.ctx.Done():
				return
			case line := <-t.tailObj.Lines:
				log.Println("文件中读到一条数据：", line)
				if line != nil {
					t.textChan <- line.Text
				}
			}
		}
	}()
	return t.textChan
}
func (t *TailTask) CtxDone() <-chan struct{} {
	return t.ctx.Done()
}
func (t *TailTask) CtxCancelFn() {
	t.ctxCancel()
}
func (t *TailTask) GetTopic() string {
	return t.topic
}
func (t *TailTask) GetFilePath() string {
	return t.filePath
}
