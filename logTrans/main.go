package main

import (
	"gopkg.in/ini.v1"
	"log"
	"logAgent/es8"
	myEtcd "logAgent/etcd"
	myKafka "logAgent/kafaka"
)

type Conf struct {
	Kafka     KafkaIni `ini:"kafka"`
	Etcd      EtcdIni  `ini:"etcd"`
	Es8       Es8Ini   `ini:"es8"`
	DebugMode bool     `ini:"debugMode"`
}

type Es8Ini struct {
	Addrs []string `ini:"addrs"`
}
type KafkaIni struct {
	Addrs []string `ini:"addrs"`
}
type EtcdIni struct {
	Addrs       []string `ini:"addrs"`
	LogEntryKey string   `ini:"logEntryKey"`
}

var cfg = new(Conf)

func initTransConf() error {
	iniFile := "../conf.ini"
	return ini.MapTo(cfg, iniFile)
}

var msgCh = make(chan myKafka.Msg, 100)

func main() {
	var err error
	err = initTransConf()
	if DebugErr(err) {
		DebugLog("initTransConf failed ,err:", err)
		return
	}
	log.Println(cfg)

	tMgrInitLogEntry(myEtcd.InitEtcdClient(cfg.Etcd.Addrs))
	tMgr.RunLogEntriesAndWatch()

	es, err := es8.InitEs8Client(cfg.Es8.Addrs)
	if DebugErr(err) {
		DebugLog("InitEs8Client failed ,err:", err)
		return
	}
	for {
		msg := <-msgCh
		es.SendToEs(msg.Topic, map[string]any{
			"data": msg.Text,
		})
	}
}
