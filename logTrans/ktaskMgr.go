package main

import (
	"encoding/json"
	myEtcd "logAgent/etcd"
	myKafka "logAgent/kafaka"
)

type LogEntry struct {
	Topic    string
	FilePath string
}

type TaskMgr struct {
	LogEntries []LogEntry
	TaskMap    map[string]*myKafka.KafkaConsumer
	Cli        *myEtcd.EtcdClient
}

var tMgr TaskMgr

// 初始化tMgr LogEntries，还没有开始任务协程
func tMgrInitLogEntry(cli *myEtcd.EtcdClient) {
	//tMgr分配内存
	tMgr.LogEntries = make([]LogEntry, 0, 5)
	tMgr.TaskMap = make(map[string]*myKafka.KafkaConsumer, 5)
	tMgr.Cli = cli
	//获取etcd已有的配置
	bStr, err := cli.GetValueByKey(cfg.Etcd.LogEntryKey)
	if err != nil || len(bStr) == 0 {
		DebugLog("tMgrInitLogEntry cli.GetValueByKey failed or len(bStr) == 0,err:", err)
		return
	}
	err = json.Unmarshal([]byte(bStr), &tMgr.LogEntries)
	if err != nil {
		DebugLog("json 失败")
	}
}

func (m TaskMgr) RunLogEntriesAndWatch() {
	//进程启动时Init etcd已有的配置，runtask
	for _, ev := range m.LogEntries {
		k, err := json.Marshal(ev)
		if err != nil {
			DebugLog("json failed,err:", err)
			continue
		}
		c := NewRunTask(ev.Topic)
		m.TaskMap[string(k)] = c
	}
	//Watch etcd热更新
	chanStr := make(chan string, 5)
	go m.Cli.WatchKeyInChanStr(cfg.Etcd.LogEntryKey, chanStr)
	go m.diffEvsFromChan(chanStr)
}

func (m TaskMgr) diffEvsFromChan(chanStr <-chan string) {
	for evs := range chanStr {
		newEvs := make([]LogEntry, 0)
		if len(evs) != 0 {
			err := json.Unmarshal([]byte(evs), &newEvs)
			if err != nil {
				DebugLog("json.Unmarshal([]byte(evs), &newEvs) err:", err)
				continue
			}
		}
		m.LogEntries = newEvs
		//新开runtask
		for _, ev := range newEvs {
			k, err := json.Marshal(ev)
			if err != nil {
				DebugLog("json failed,err:", err)
				continue
			}
			_, ok := m.TaskMap[string(k)]
			if !ok {
				t := NewRunTask(ev.Topic)
				m.TaskMap[string(k)] = t
			}
		}
		//删除无效的runtask
		for k, t := range m.TaskMap {
			delFlag := true
			for _, ev := range newEvs {
				var entry LogEntry
				err := json.Unmarshal([]byte(k), &entry)
				if err != nil {
					DebugLog("json.Unmarshal([]byte(evs), &newEvs) err:", err)
					continue
				}
				if entry.Topic == ev.Topic &&
					entry.FilePath == ev.FilePath {
					delFlag = false
				}
			}
			if delFlag {
				t.CancelFunc()
			}
		}
	}
}

func NewRunTask(topic string) *myKafka.KafkaConsumer {
	c, err := myKafka.InitKafkaConsumer(cfg.Kafka.Addrs)
	if DebugErr(err) {
		DebugLog("InitKafkaConsumer fiall,err:", err, topic)
		return nil
	}
	c.FromMsgKafka(topic, msgCh)
	return c
}
