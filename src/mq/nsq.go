package mq

import (
	"dbmcloud/log"
	"dbmcloud/setting"
	"encoding/json"
	"fmt"

	"github.com/nsqio/go-nsq"
)

var NSQ *nsq.Producer

func InitNsq() *nsq.Producer {
	producer, err := nsq.NewProducer(setting.Setting.NsqServer, nsq.NewConfig())
	if err != nil {
		panic(err)
	}
	/*
		defer func() {
			if producer != nil {
				producer.Stop()
			}
		}()
	*/
	return producer
}

func Send(m map[string]interface{}) {
	/*
		convert event map to json string
	*/
	d, _ := json.Marshal(m)
	eventStr := string(d)
	/*
		send event json str to kafka
	*/
	log.Logger.Debug(fmt.Sprintf("Event map data:%s", m))
	log.Logger.Debug(fmt.Sprintf("Event json data:%s", eventStr))
	if err := NSQ.Publish("lepus_events", []byte(eventStr)); err != nil { // 发布消息
		panic(err)
	}

}
