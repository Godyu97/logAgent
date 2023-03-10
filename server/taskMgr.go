package main

import (
	"encoding/json"
	myEtcd "logAgent/etcd"
	myTailf "logAgent/tailf"
)

type LogEntry struct {
	Topic    string
	FilePath string
}

type TaskMgr struct {
	LogEntries []LogEntry
	TaskMap    map[string]*myTailf.TailTask
	Cli        *myEtcd.EtcdClient
}

var tMgr TaskMgr

// 初始化tMgr LogEntries，还没有开始任务协程
func tMgrInitLogEntry(cli *myEtcd.EtcdClient) {
	//tMgr分配内存
	tMgr.LogEntries = make([]LogEntry, 0, 5)
	tMgr.TaskMap = make(map[string]*myTailf.TailTask, 5)
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
		t := NewRunTask(ev)
		m.TaskMap[string(k)] = t
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
				t := NewRunTask(ev)
				m.TaskMap[string(k)] = t
			}
		}
		//删除无效的runtask
		for _, t := range m.TaskMap {
			delFlag := true
			for _, ev := range newEvs {
				if t.GetTopic() == ev.Topic &&
					t.GetFilePath() == ev.FilePath {
					delFlag = false
				}
			}
			if delFlag {
				t.CtxCancelFn()
			}
		}
	}
}
