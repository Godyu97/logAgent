package es8

import (
	"bytes"
	"encoding/json"
	elastic8 "github.com/elastic/go-elasticsearch/v8"
	"log"
)

type Es8Client struct {
	*elastic8.Client
}

func InitEs8Client(addrs []string) (*Es8Client, error) {
	es8, err := elastic8.NewClient(elastic8.Config{
		Addresses: addrs,
	})
	if err != nil {
		return nil, err
	}
	log.Println(es8.Info())
	return &Es8Client{es8}, nil
}

func (e *Es8Client) SendToEs(index string, r map[string]any) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(r)
	e.Index(index, buf)
}
