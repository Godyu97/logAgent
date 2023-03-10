package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	log.Println("test mod!")
	//key := "logAgentServer_nick"
	//c := myEtcd.InitEtcdClient([]string{"185.183.84.249:2379"})
	//l1 := []LogEntry{
	//	{
	//		Topic:    "mysql",
	//		FilePath: "./mysql.log",
	//	},
	//	{
	//		Topic:    "redis",
	//		FilePath: "./redis.log",
	//	},
	//	{
	//		Topic:    "kafka",
	//		FilePath: "./kafka.log",
	//	},
	//	{
	//		Topic:    "pgsql",
	//		FilePath: "./pgsql.log",
	//	},
	//	{
	//		Topic:    "etcd",
	//		FilePath: "./etcd.log",
	//	},
	//}

	//l2 := []LogEntry{
	//	{
	//		Topic:    "redis",
	//		FilePath: "./redis.log",
	//	},
	//	{
	//		Topic:    "kafka",
	//		FilePath: "./kafka.log",
	//	},
	//	{
	//		Topic:    "pgsql",
	//		FilePath: "./pgsql.log",
	//	},
	//	{
	//		Topic:    "etcd",
	//		FilePath: "./etcd.log",
	//	},
	//}

	//l3 := [][]LogEntry{}

	//s, _ := json.Marshal(l1)
	//c.PutKeyValue("logAgentServer_nick", string(s))
	////c.DeleteKey(key)
	//log.Println(c.GetValueByKey(key))

	//写入文件
	f, err := os.OpenFile("../server/pgsql.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	for n := 0; ; n++ {
		_, err = f.WriteString(fmt.Sprintf("line:%d\n", n))
		if err != nil {
			panic(err)
		}
		log.Printf("line:%d\n", n)
		time.Sleep(time.Second * 3)
	}
}

type LogEntry struct {
	Topic    string
	FilePath string
}

//TODO
//1. 搭建kafka
//2. 据需看接下来的项目课程
