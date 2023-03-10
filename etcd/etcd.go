package myEtcd

import (
	"context"
	"errors"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type EtcdClient struct {
	cli *clientv3.Client
}

func InitEtcdClient(addrs []string) *EtcdClient {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   addrs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil
	}
	return &EtcdClient{cli: client}
}

func (c *EtcdClient) GetValueByKey(key string) (value string, err error) {
	resp, err := c.cli.Get(context.Background(), key)
	if err != nil {
		return "", err
	}
	for _, ev := range resp.Kvs {
		return string(ev.Value), nil
	}
	return "", errors.New("no value by key")
}

func (c *EtcdClient) PutKeyValue(key string, value string) error {
	_, err := c.cli.Put(context.Background(), key, value)
	return err
}

func (c *EtcdClient) DeleteKey(key string) error {
	_, err := c.cli.Delete(context.Background(), key)
	return err
}

func (c *EtcdClient) WatchKeyInChanStr(key string, chanStr chan<- string) {
	for resp := range c.cli.Watch(context.Background(), key) {
		for _, ev := range resp.Events {
			chanStr <- string(ev.Kv.Value)
		}
	}
}
