package main

import (
	"gopkg.in/ini.v1"
	"log"
	myEtcd "logAgent/etcd"
	myKafka "logAgent/kafaka"
	"os"
	"os/signal"
)

type KafkaIni struct {
	Addrs []string `ini:"addrs"`
}
type EtcdIni struct {
	Addrs       []string `ini:"addrs"`
	LogEntryKey string   `ini:"logEntryKey"`
}
type Conf struct {
	Kafka     KafkaIni `ini:"kafka"`
	Etcd      EtcdIni  `ini:"etcd"`
	DebugMode bool     `ini:"debugMode"`
}

var kafkaProducerObj = new(myKafka.KafkaProducer)
var cfg = new(Conf)

func initSeverConf() error {
	iniFile := "../conf.ini"
	return ini.MapTo(cfg, iniFile)
}

func main() {
	var err error
	err = initSeverConf()
	if err != nil {
		DebugLog("initSeverConf failed ,err:", err)
		return
	}
	DebugLog("iniCFG：", *cfg)
	kafkaProducerObj, err = myKafka.InitKafkaProducer(cfg.Kafka.Addrs)
	if err != nil {
		DebugLog("InitKafkaProducer failed ,err:", err)
	}
	go SendToKafka()
	tMgrInitLogEntry(myEtcd.InitEtcdClient(cfg.Etcd.Addrs))
	tMgr.RunLogEntriesAndWatch()

	//阻塞进程
	c := make(chan os.Signal)
	signal.Notify(c)
	for ctrlC := range c {
		if ctrlC == os.Interrupt {
			log.Println(ctrlC, " Bye")
			return
		}
	}
}
