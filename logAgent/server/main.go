package main

import (
	"gopkg.in/ini.v1"
	"log"
	myEtcd "logAgent/etcd"
	myKafka "logAgent/kafaka"
	myTailf "logAgent/tailf"
	"sync"
)

type LogEntry struct {
	Topic    string
	FilePath string
}

type KafkaIni struct {
	Addrs []string `ini:"addrs"`
}
type EtcdIni struct {
	Addrs []string `ini:"addrs"`
}
type Conf struct {
	Kafka     KafkaIni `ini:"kafka"`
	Etcd      EtcdIni  `ini:"etcd"`
	Logs      []LogEntry
	DebugMode bool `ini:"debugMode"`
}

var kafkaProducerObj = new(myKafka.KafkaProducer)
var cfg = new(Conf)

func initSeverConf() error {
	iniFile := "./conf.ini"
	return ini.MapTo(cfg, iniFile)
}

func main() {
	log.Println("test server")
	var err error
	err = initSeverConf()
	if err != nil {
		DebugLog("initSeverConf failed ,err:", err)
	}
	DebugLog(*cfg)
	kafkaProducerObj, err = myKafka.InitKafkaProducer(cfg.Kafka.Addrs)
	if err != nil {
		DebugLog("InitKafkaProducer failed ,err:", err)
	}
	myEtcd.InitEtcdClient(cfg.Etcd.Addrs)
	t, err := myTailf.NewTailTask("./log", "giegie")
	if err != nil {
		DebugLog("NewTailTask failed ,err:", err)
	}
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go run(t, wg)

	go SendToKafka()
	wg.Wait()
}
