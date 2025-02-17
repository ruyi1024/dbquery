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
	"dbmcloud/log"
	"dbmcloud/setting"
	"dbmcloud/src/database"

	//_ "dbmcloud/src/libary/redis"
	"dbmcloud/src/libary/tool"
	"dbmcloud/src/model"
	"dbmcloud/src/mq"
	"dbmcloud/src/utils"
	"fmt"
	_ "reflect"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/robfig/cron/v3"
)

func init() {
	go collectorRedisEventTask()
}

func collectorRedisEventTask() {

	time.Sleep(time.Second * time.Duration(35))
	var db = database.DB
	var record model.TaskOption
	db.Select("crontab").Where("task_key=?", "collector_redis_event").Take(&record)
	c := cron.New()
	c.AddFunc(record.Crontab, func() {
		db.Select("enable").Where("task_key=?", "collector_redis_event").Take(&record)
		if record.Enable == 1 {
			db.Model(model.TaskHeartbeat{}).Where("heartbeat_key='collector_redis_event'").Updates(map[string]interface{}{"heartbeat_time": time.Now().Format("2006-01-02 15:04:05.999")})
			doCollectorRedisEventTask()
			db.Model(model.TaskHeartbeat{}).Where("heartbeat_key='collector_redis_event'").Updates(map[string]interface{}{"heartbeat_end_time": time.Now().Format("2006-01-02 15:04:05.999")})
		}
	})
	c.Start()
}

func doCollectorRedisEventTask() {
	var db = database.DB
	var dataList []model.Datasource
	result := db.Where("enable=1").Where("type = ? ", "Redis").Order("type asc").Find(&dataList)
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

		go startCollectRedisEvent(datasourceName, datasourceType, env, host, port, origPass)

	}

}

func startCollectRedisEvent(datasourceName, datasourceType, env, host, port, origPass string) {
	eventEntity := fmt.Sprintf("%s:%s", host, port)
	eventType := datasourceType
	eventGroup := env
	var connect int = 1
	rdb, err := redis.Dial("tcp", host+":"+port)
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
	defer rdb.Close()

	if origPass != "" {
		if _, err := rdb.Do("AUTH", origPass); err != nil {
			rdb.Close()
			log.Error(fmt.Sprintln("Redis Auth error, ", err))
			return
		}
	}

	info, err := redis.String(rdb.Do("INFO"))
	if err != nil {
		log.Error(fmt.Sprintln("Can't do redis info query, ", err))
		return
	}

	infoMap := make(map[string]string)
	infoArray := strings.Split(info, "\n")
	for _, item := range infoArray {
		if strings.Contains(item, ":") {
			v := strings.Split(item, ":")
			infoMap[v[0]] = v[1]
		}

	}

	role := strings.Replace(infoMap["role"], "\r", "", -1)
	redisVersion := strings.Replace(infoMap["redis_version"], "\r", "", -1)
	// redisMode := infoMap["redis_mode"]
	// os := infoMap["os"]
	// archBits := infoMap["arch_bits"]
	// gccVersion := infoMap["gcc_version"]
	// processId := infoMap["process_id"]
	// runId := infoMap["run_id"]
	// tcpPort := infoMap["tcp_port"]
	uptimeInSeconds := utils.StrToInt(strings.Replace(infoMap["uptime_in_seconds"], "\r", "", -1))
	uptimeInDays := utils.StrToInt(strings.Replace(infoMap["uptime_in_days"], "\r", "", -1))
	connectedClients := utils.StrToInt(strings.Replace(infoMap["connected_clients"], "\r", "", -1))
	blockedClients := utils.StrToInt(strings.Replace(infoMap["blocked_clients"], "\r", "", -1))
	usedMemory := utils.StrToInt(strings.Replace(infoMap["used_memory"], "\r", "", -1))
	//usedMemoryHuman := infoMap["used_memory_human"]
	usedMemoryRss := utils.StrToInt(strings.Replace(infoMap["used_memory_rss"], "\r", "", -1)) / 1024 / 1024
	//usedMemoryRssHuman := infoMap["used_memory_rss_human"]
	usedMemoryPeak := utils.StrToInt(strings.Replace(infoMap["used_memory_peak"], "\r", "", -1)) / 1024 / 1024
	//usedMemoryPeakHuman := infoMap["used_memory_peak_human"]
	usedMemoryLua := utils.StrToInt(strings.Replace(infoMap["used_memory_lua"], "\r", "", -1)) / 1024 / 1024
	//usedMemoryLuaHuman := infoMap["used_memory_lua_human"]
	memFragmentationRatio := utils.StrToInt(strings.Replace(infoMap["mem_fragmentation_ratio"], "\r", "", -1))
	//memAllocator := infoMap["mem_allocator"]
	// rdbBgsaveInProgress := infoMap["rdb_bgsave_in_progress"]
	// rdbLastSaveTime := infoMap["rdb_last_save_time"]
	// rdbLastBgsaveStatus := infoMap["rdb_last_bgsave_status"]
	// rdbLastBgsaveTimeSec := infoMap["rdb_last_bgsave_time_sec"]
	// aofEnabled := infoMap["aof_enabled"]
	// aofRewriteInProgress := infoMap["aof_rewrite_in_progress"]
	// aofRewriteScheduled := infoMap["aof_rewrite_scheduled"]
	// aofLastRewriteTimeSec := infoMap["aof_last_rewrite_time_sec"]
	// aofLastBgrewriteStatus := infoMap["aof_last_bgrewrite_status"]
	totalConnectionsReceived := utils.StrToInt(strings.Replace(infoMap["total_connections_received"], "\r", "", -1))
	totalCommandsProcessed := utils.StrToInt(strings.Replace(infoMap["total_commands_processed"], "\r", "", -1))
	instantaneousOpsPerSec := utils.StrToInt(strings.Replace(infoMap["instantaneous_ops_per_sec"], "\r", "", -1))
	rejectedConnections := utils.StrToInt(strings.Replace(infoMap["rejected_connections"], "\r", "", -1))
	expiredKeys := utils.StrToInt(strings.Replace(infoMap["expired_keys"], "\r", "", -1))
	evictedKeys := utils.StrToInt(strings.Replace(infoMap["evicted_keys"], "\r", "", -1))
	keyspaceHits := utils.StrToInt(strings.Replace(infoMap["keyspace_hits"], "\r", "", -1))
	keyspaceMisses := utils.StrToInt(strings.Replace(infoMap["keyspace_misses"], "\r", "", -1))
	usedCpuSys := utils.StrToInt(strings.Replace(infoMap["used_cpu_sys"], "\r", "", -1))
	usedCpuUser := utils.StrToInt(strings.Replace(infoMap["used_cpu_user"], "\r", "", -1))
	usedCpuSysChildren := utils.StrToInt(strings.Replace(infoMap["used_cpu_sys_children"], "\r", "", -1))
	usedCpuUserChildren := utils.StrToInt(strings.Replace(infoMap["used_cpu_user_children"], "\r", "", -1))

	maxClientsConfig, err := redis.Strings(rdb.Do("config", "GET", "maxclients"))
	if err != nil {
		log.Error(fmt.Sprintln("Can't do redis maxclients query, ", err))
		return
	}
	maxClients := utils.StrToInt(maxClientsConfig[1])

	maxMemoryConfig, err := redis.Strings(rdb.Do("config", "GET", "maxmemory"))
	if err != nil {
		log.Error(fmt.Sprintln("Can't do redis maxMemory query, ", err))
		return
	}

	maxMemory := utils.StrToInt(maxMemoryConfig[1]) / 1024 / 1024
	var usedMemoryPct int
	if maxMemory == 0 {
		usedMemoryPct = 0
	} else {
		usedMemoryPct = (usedMemory / maxMemory) * 100
	}

	events := make([]map[string]interface{}, 0)
	//emptyDetail := make([]map[string]interface{},0)

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
		"event_key":    "uptimeInSeconds",
		"event_value":  utils.IntToDecimal(uptimeInSeconds),
		"event_tag":    "",
		"event_unit":   "秒",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "uptimeInDays",
		"event_value":  utils.IntToDecimal(uptimeInDays),
		"event_tag":    "",
		"event_unit":   "天",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "connectedClients",
		"event_value":  utils.IntToDecimal(connectedClients),
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
		"event_key":    "blockedClients",
		"event_value":  utils.IntToDecimal(blockedClients),
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
		"event_key":    "maxClients",
		"event_value":  utils.IntToDecimal(maxClients),
		"event_tag":    "",
		"event_unit":   "%",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "usedMemoryPct",
		"event_value":  utils.IntToDecimal(usedMemoryPct),
		"event_tag":    "",
		"event_unit":   "%",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "usedMemory",
		"event_value":  utils.IntToDecimal(usedMemory),
		"event_tag":    "",
		"event_unit":   "MB",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "usedMemoryRss",
		"event_value":  utils.IntToDecimal(usedMemoryRss),
		"event_tag":    "",
		"event_unit":   "MB",
		"event_detail": "",
	}

	events = append(events, event)
	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "usedMemoryPeak",
		"event_value":  utils.IntToDecimal(usedMemoryPeak),
		"event_tag":    "",
		"event_unit":   "MB",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "usedMemoryLua",
		"event_value":  utils.IntToDecimal(usedMemoryLua),
		"event_tag":    "",
		"event_unit":   "MB",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "memFragmentationRatio",
		"event_value":  utils.IntToDecimal(memFragmentationRatio),
		"event_unit":   "%",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "totalConnectionsReceived",
		"event_value":  utils.IntToDecimal(totalConnectionsReceived),
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
		"event_key":    "totalCommandsProcessed",
		"event_value":  utils.IntToDecimal(totalCommandsProcessed),
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
		"event_key":    "instantaneousOpsPerSec",
		"event_value":  utils.IntToDecimal(instantaneousOpsPerSec),
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
		"event_key":    "rejectedConnections",
		"event_value":  utils.IntToDecimal(rejectedConnections),
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
		"event_key":    "expiredKeys",
		"event_value":  utils.IntToDecimal(expiredKeys),
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
		"event_key":    "evictedKeys",
		"event_value":  utils.IntToDecimal(evictedKeys),
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
		"event_key":    "keyspaceHits",
		"event_value":  utils.IntToDecimal(keyspaceHits),
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
		"event_key":    "keyspaceMisses",
		"event_value":  utils.IntToDecimal(keyspaceMisses),
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
		"event_key":    "usedCpuSys",
		"event_value":  utils.IntToDecimal(usedCpuSys),
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
		"event_key":    "usedCpuUser",
		"event_value":  utils.IntToDecimal(usedCpuUser),
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
		"event_key":    "usedCpuSysChildren",
		"event_value":  utils.IntToDecimal(usedCpuSysChildren),
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
		"event_key":    "usedCpuUserChildren",
		"event_value":  utils.IntToDecimal(usedCpuUserChildren),
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

	//write mysql
	var dataList []model.EventGlobal
	database.DB.Where("host=?", host).Where("port=?", port).Find(&dataList)
	if (len(dataList)) == 0 {
		var record model.EventGlobal
		record.DatasourceType = datasourceType
		record.DatasourceName = datasourceName
		record.Host = host
		record.Port = port
		record.Version = redisVersion
		record.Uptime = int64(uptimeInSeconds)
		record.Connect = connect
		record.Role = role
		record.Session = connectedClients
		record.Active = -1
		record.Wait = blockedClients
		record.Qps = instantaneousOpsPerSec
		record.Tps = -1
		record.Repl = -1
		record.Delay = -1

		result := database.DB.Create(&record)
		if result.Error != nil {
			log.Logger.Error("Insert Error:" + result.Error.Error())
		}

	} else {
		var record model.EventGlobal
		record.Version = redisVersion
		record.Uptime = int64(uptimeInSeconds)
		record.Connect = connect
		record.Role = role
		record.Session = connectedClients
		record.Active = -1
		record.Wait = blockedClients
		record.Qps = instantaneousOpsPerSec
		record.Tps = -1
		record.Repl = -1
		record.Delay = -1
		//gin里面如果更新为0则字段不会更新，可以使用select指定更新字段解决
		result := database.DB.Model(&record).Select("version", "uptime", "role", "session", "active", "wait", "qps", "tps", "repl", "delay").Omit("id").Where("host=?", host).Where("port=?", port).Updates(&record)
		if result.Error != nil {
			log.Logger.Error("Update Error:" + result.Error.Error())
		}
	}

}
