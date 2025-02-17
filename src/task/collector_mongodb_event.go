/*
Copyright 2014-2022 The Lepus Team Group, website: https://www.lepus.cc
Licensed under the GNU General Public License, Version 3.0 (the "GPLv3 License");
You may not use this file except in compliance with the License.
You may obtain a copy of the License at
    https://www.gnu.org/licenses/gpl-3.0.html
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Special note:
Please do not use this source code for any commercial purpose,
or use it for commercial purposes after secondary development, otherwise you may bear legal risks.
*/

package task

import (
	"context"
	"dbmcloud/log"
	"dbmcloud/setting"
	"dbmcloud/src/database"
	"dbmcloud/src/libary/mongodb"
	"dbmcloud/src/libary/tool"
	"dbmcloud/src/model"
	"dbmcloud/src/mq"
	"dbmcloud/src/utils"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
)

func init() {
	go collectorMongodbEventTask()
}

func collectorMongodbEventTask() {
	time.Sleep(time.Second * time.Duration(30))
	var db = database.DB
	var record model.TaskOption
	db.Select("crontab").Where("task_key=?", "collector_mongodb_event").Take(&record)
	c := cron.New()
	c.AddFunc(record.Crontab, func() {
		db.Select("enable").Where("task_key=?", "collector_mongodb_event").Take(&record)
		if record.Enable == 1 {
			db.Model(model.TaskHeartbeat{}).Where("heartbeat_key='collector_mongodb_event'").Updates(map[string]interface{}{"heartbeat_time": time.Now().Format("2006-01-02 15:04:05.999")})
			doCollectorMongodbEventTask()
			db.Model(model.TaskHeartbeat{}).Where("heartbeat_key='collector_mongodb_event'").Updates(map[string]interface{}{"heartbeat_end_time": time.Now().Format("2006-01-02 15:04:05.999")})
		}
	})
	c.Start()
}

func doCollectorMongodbEventTask() {
	var db = database.DB
	var dataList []model.Datasource
	result := db.Where("enable=1").Where("type = ? ", "Mongodb").Order("type asc").Find(&dataList)
	if result.Error != nil {
		log.Logger.Error(result.Error.Error())
		return

	}

	for _, datasource := range dataList {
		datasourceName := datasource.Name
		datasourceType := datasource.Type
		env := datasource.Env
		host := datasource.Host
		port := datasource.Port
		user := datasource.User
		pass := datasource.Pass

		var origPass string
		if pass != "" {
			var err error
			origPass, err = utils.AesPassDecode(pass, setting.Setting.DbPassKey)
			if err != nil {
				fmt.Println("Encrypt Password Error.")
				return
			}
		}
		startCollectMongodbEvent(datasourceName, datasourceType, env, host, port, user, origPass)

	}

}

func startCollectMongodbEvent(datasourceName, datasourceType, env, host, port, user, origPass string) {
	eventEntity := fmt.Sprintf("%s:%s", host, port)
	eventType := datasourceType
	eventGroup := env
	var connect int = 1

	client, err := mongodb.Connect(host, port, user, origPass, "admin")

	if err != nil {
		connect = 0
		detail := make([]map[string]interface{}, 0)
		detail = append(detail, map[string]interface{}{"Error": fmt.Sprint(err)})
		events := make([]map[string]interface{}, 0)
		event := map[string]interface{}{
			"event_uuid":   tool.GetUUID(),
			"event_time":   tool.GetNowTime(),
			"event_type":   eventType,
			"event_group":  eventGroup,
			"event_entity": eventEntity,
			"event_key":    "connect",
			"event_value":  utils.IntToDecimal(connect),
			"event_tag":    "",
			"event_unit":   "",
			"event_detail": utils.MapToStr(detail),
		}
		events = append(events, event)

		// write events to ck
		result := database.CK.Model(&model.Event{}).Create(events)
		if result.Error != nil {
			fmt.Println("Insert Event To Clickhouse Error: " + result.Error.Error())
			log.Logger.Error(fmt.Sprintf("Can't add events data to clickhouse: %s", result.Error.Error()))
			return
		}
		//send event to nsq
		for _, event := range events {
			mq.Send(event)
		}
		return
	}

	db := client.Database("admin")
	command := bson.D{{"serverStatus", 1}}
	status := bson.M{}
	err = db.RunCommand(context.TODO(), command).Decode(&status)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(status)
	time.Sleep(time.Duration(1) * time.Second)

	statusNew := bson.M{}
	err = db.RunCommand(context.TODO(), command).Decode(&statusNew)
	if err != nil {
		fmt.Println(err)
		return
	}

	// ok := status["ok"]
	// version := status["version"]
	// fmt.Println(ok)
	// fmt.Println(version)
	uptime := status["uptime"].(float64)
	//fmt.Print(uptime.(float64))
	connections := status["connections"].(bson.M)
	connectionsCurrent := connections["current"].(int32)
	connectionsAvailable := connections["available"].(int32)
	mem := status["mem"].(bson.M)
	memBits := mem["bits"].(int32)
	memResident := mem["resident"].(int32)
	memVirtual := mem["virtual"].(int32)
	//memSupported := mem["supported"]
	//memMapped := mem["mapped"].(int32)
	//memMappedWithJournal := mem["mappedWithJournal"].(int32)
	// var memSupportedInt = 0
	// if memSupported.(bool) {
	// 	memSupportedInt = 1
	// }
	network := status["network"].(bson.M)
	opcounters := status["opcounters"].(bson.M)
	networkNew := statusNew["network"].(bson.M)
	opcountersNew := statusNew["opcounters"].(bson.M)

	networkBytesIn := networkNew["bytesIn"].(int64) - network["bytesIn"].(int64)
	networkBytesOut := networkNew["bytesOut"].(int64) - network["bytesOut"].(int64)
	networkNumRequests := networkNew["numRequests"].(int64) - network["numRequests"].(int64)

	var (
		opcountersInsert  int64
		opcountersQuery   int64
		opcountersUpdate  int64
		opcountersDelete  int64
		opcountersCommand int64
		opcountersTotal   int64
	)
	//opcountersType := fmt.Sprintf("%T", opcounters["insert"])
	//if opcountersType == "int64" {
	opcountersInsert = opcountersNew["insert"].(int64) - opcounters["insert"].(int64)
	opcountersQuery = opcountersNew["query"].(int64) - opcounters["query"].(int64)
	opcountersUpdate = opcountersNew["update"].(int64) - opcounters["update"].(int64)
	opcountersDelete = opcountersNew["delete"].(int64) - opcounters["delete"].(int64)
	opcountersCommand = opcountersNew["command"].(int64) - opcounters["command"].(int64)
	opcountersTotal = (opcountersNew["insert"].(int64) - opcounters["insert"].(int64)) + (opcountersNew["query"].(int64) - opcounters["query"].(int64)) + (opcountersNew["update"].(int64) - opcounters["update"].(int64)) + (opcountersNew["delete"].(int64) - opcounters["delete"].(int64)) + (opcountersNew["command"].(int64) - opcounters["command"].(int64))
	//}

	events := make([]map[string]interface{}, 0)

	event := map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "connect",
		"event_value":  utils.IntToDecimal(connect),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "uptime",
		"event_value":  decimal.NewFromFloat(uptime),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "connectionsCurrent",
		"event_value":  utils.Int32ToDecimal(connectionsCurrent),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "connectionsAvailable",
		"event_value":  utils.Int32ToDecimal(connectionsAvailable),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "memBits",
		"event_value":  utils.Int32ToDecimal(memBits),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "memResident",
		"event_value":  utils.Int32ToDecimal(memResident),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "memVirtual",
		"event_value":  utils.Int32ToDecimal(memVirtual),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": "",
	}
	events = append(events, event)
	/*
		event = map[string]interface{}{
			"event_uuid":   tool.GetUUID(),
			"event_time":   tool.GetNowTime(),
			"event_type":   eventType,
			"event_group":  eventGroup,
			"event_entity": eventEntity,
			"event_key":    "memMapped",
			"event_value":  utils.Int32ToDecimal(memMapped),
			"event_tag":    "",
			"event_unit":   "",
			"event_detail": "",
		}
		events = append(events, event)

		event = map[string]interface{}{
			"event_uuid":   tool.GetUUID(),
			"event_time":   tool.GetNowTime(),
			"event_type":   eventType,
			"event_group":  eventGroup,
			"event_entity": eventEntity,
			"event_key":    "memMappedWithJournal",
			"event_value":  utils.Int32ToDecimal(memMappedWithJournal),
			"event_tag":    "",
			"event_unit":   "",
			"event_detail": "",
		}
		events = append(events, event)
	*/
	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "networkBytesIn",
		"event_value":  utils.Int64ToDecimal(networkBytesIn),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "networkBytesOut",
		"event_value":  utils.Int64ToDecimal(networkBytesOut),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "networkNumRequests",
		"event_value":  utils.Int64ToDecimal(networkNumRequests),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "opcountersInsert",
		"event_value":  utils.Int64ToDecimal(opcountersInsert),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "opcountersQuery",
		"event_value":  utils.Int64ToDecimal(opcountersQuery),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "opcountersUpdate",
		"event_value":  utils.Int64ToDecimal(opcountersUpdate),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "opcountersDelete",
		"event_value":  utils.Int64ToDecimal(opcountersDelete),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "opcountersCommand",
		"event_value":  utils.Int64ToDecimal(opcountersCommand),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "opcounters",
		"event_value":  utils.Int64ToDecimal(opcountersTotal),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": "",
	}
	events = append(events, event)

	// write events to ck
	result := database.CK.Model(&model.Event{}).Create(events)
	if result.Error != nil {
		fmt.Println("Insert Event To Clickhouse Error: " + result.Error.Error())
		log.Logger.Error(fmt.Sprintf("Can't add events data to clickhouse: %s", result.Error.Error()))
		return
	}

	//send event to nsq
	for _, event := range events {
		mq.Send(event)
	}

}
