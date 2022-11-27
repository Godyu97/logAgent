package myEtcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func InitEtcdClient(addrs []string) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   addrs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return
	}
	defer client.Close()
}
