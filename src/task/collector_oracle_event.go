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
	"dbmcloud/src/libary/oracle"
	"dbmcloud/src/libary/tool"
	"dbmcloud/src/model"
	"dbmcloud/src/mq"
	"dbmcloud/src/utils"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func init() {
	go collectorOracleEventTask()
}

var (
	queryInstanceSQL      = "select * from v$instance"
	queryDatabaseSQL      = "select * from v$database"
	queryParameterSQL     = "select name,type,value from v$parameter"
	queryProcessesSQL     = "select name,type,value from v$parameter where name='processes' "
	querySysStatSQL       = "select name,value from v$sysstat "
	querySessionTotalSQL  = "Select a.SID,a.SERIAL#,a.STATUS,a.USERNAME,a.MACHINE,a.MODULE,a.EVENT,b.SQL_ID,b.SQL_TEXT from v$session a, v$sqlarea b where a.sql_hash_value = b.HASH_VALUE and  a.username not in('SYS','SYSTEM') and a.username is not null"
	querySessionActiveSQL = "Select a.SID,a.SERIAL#,a.STATUS,a.USERNAME,a.MACHINE,a.MODULE,a.EVENT,b.SQL_ID,b.SQL_TEXT from v$session a, v$sqlarea b where a.sql_hash_value = b.HASH_VALUE and  a.username not in('SYS','SYSTEM') and a.username is not null and a.status='ACTIVE'"
	querySessionWaitSQL   = "Select a.SID,a.SERIAL#,a.STATUS,a.USERNAME,a.MACHINE,a.MODULE,a.EVENT,b.SQL_ID,b.SQL_TEXT from v$session a, v$sqlarea b where a.sql_hash_value = b.HASH_VALUE  and a.username is not null and  ( a.event like 'library%' or a.event like 'cursor%' or a.event like 'latch%'  or a.event like 'enq%' or a.event like 'log file%')"
	queryDataGuardSQL     = "SELECT substr((SUBSTR(VALUE,5)),0,2)*3600 + substr((SUBSTR(VALUE,5)),4,2)*60 + substr((SUBSTR(VALUE,5)),7,2) AS seconds,VALUE FROM v$dataguard_stats a WHERE NAME ='apply lag'"
)

func collectorOracleEventTask() {
	time.Sleep(time.Second * time.Duration(30))
	var db = database.DB
	var record model.TaskOption
	db.Select("crontab").Where("task_key=?", "collector_oracle_event").Take(&record)
	c := cron.New()
	c.AddFunc(record.Crontab, func() {
		db.Select("enable").Where("task_key=?", "collector_oracle_event").Take(&record)
		if record.Enable == 1 {
			db.Model(model.TaskHeartbeat{}).Where("heartbeat_key='collector_oracle_event'").Updates(map[string]interface{}{"heartbeat_time": time.Now().Format("2006-01-02 15:04:05.999")})
			doCollectorOracleEventTask()
			db.Model(model.TaskHeartbeat{}).Where("heartbeat_key='collector_oracle_event'").Updates(map[string]interface{}{"heartbeat_end_time": time.Now().Format("2006-01-02 15:04:05.999")})
		}
	})
	c.Start()
}

func doCollectorOracleEventTask() {
	var db = database.DB
	var dataList []model.Datasource
	result := db.Where("enable=1").Where("type = ? ", "Oracle").Order("type asc").Find(&dataList)
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
		sid := datasource.Dbid

		var origPass string
		if pass != "" {
			var err error
			origPass, err = utils.AesPassDecode(pass, setting.Setting.DbPassKey)
			if err != nil {
				fmt.Println("Encrypt Password Error.")
				return
			}
		}

		startCollectOracleEvent(datasourceName, datasourceType, env, host, port, user, origPass, sid)

	}

}

func startCollectOracleEvent(datasourceName, datasourceType, env, host, port, user, origPass, sid string) {
	eventEntity := fmt.Sprintf("%s:%s", host, port)
	eventType := datasourceType
	eventGroup := env
	var connect int = 1

	oradb, err := oracle.Connect(host, port, user, origPass, sid)

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
	defer oradb.Close()

	instance, err := oracle.QueryAll(oradb, queryInstanceSQL)
	if err != nil {
		log.Error(fmt.Sprintf("Can't query oracle instance on %s, %s", eventEntity, err))
		return
	}
	databaseInfo, err := oracle.QueryAll(oradb, queryDatabaseSQL)
	if err != nil {
		log.Error(fmt.Sprintf("Can't query oracle database on %s, %s", eventEntity, err))
		return
	}
	/*
		processInfo, err := oracle.QueryAll(oradb, queryProcessesSQL)
		if err != nil {
			log.Error(fmt.Sprintf("Can't query oracle processes on %s, %s", eventEntity, err))
			return
		}
	*/
	//fmt.Println(instance)
	//fmt.Println(databaseInfo)
	// instanceName := instance[0]["INSTANCE_NAME"]
	// instanceRole := instance[0]["INSTANCE_ROLE"]
	// instanceStatus := instance[0]["STATUS"]
	startupTime := instance[0]["startup_time"]
	version := instance[0]["version"].(string)
	// databaseStatus := instance[0]["DATABASE_STATUS"]
	// hostname := instance[0]["HOST_NAME"]
	// archiver := instance[0]["ARCHIVER"]
	databaseRole := databaseInfo[0]["database_role"]

	// openMode := database[0]["OPEN_MODE"]
	// protectedMode := database[0]["PROTECTION_MODE"]
	// processes := processInfo[0]["VALUE"]
	uptime := time.Now().Unix() - startupTime.(time.Time).Unix()

	rows, err := oradb.Query(querySysStatSQL)
	if err != nil {
		log.Error(fmt.Sprintf("Can't query oracle sysstat on%s, %s", eventEntity, err))
		return
	}
	defer rows.Close()
	var key, value string
	sysStatsPrev := make(map[string]string)
	for rows.Next() {
		err := rows.Scan(&key, &value)
		if err != nil {
			log.Error(fmt.Sprintf("Can't scan oracle sysstat on%s, %s", eventEntity, err))
			return
		}
		sysStatsPrev[key] = value
	}

	time.Sleep(time.Duration(1) * time.Second)

	rows, err = oradb.Query(querySysStatSQL)
	if err != nil {
		log.Error(fmt.Sprintf("Can't query oracle sysstat on%s, %s", eventEntity, err))
		return
	}
	defer rows.Close()
	sysStats := make(map[string]string)
	for rows.Next() {
		err := rows.Scan(&key, &value)
		if err != nil {
			log.Error(fmt.Sprintf("Can't scan oracle sysstat on%s, %s", eventEntity, err))
			return
		}
		sysStats[key] = value
	}

	sessionLogicalReadsPersecond := utils.StrToInt(sysStats["session_logical_reads"]) - utils.StrToInt(sysStatsPrev["session_logical_reads"])
	physicalReadsPersecond := utils.StrToInt(sysStats["physical read"]) - utils.StrToInt(sysStatsPrev["physical read"])
	physicalWritePersecond := utils.StrToInt(sysStats["physical write"]) - utils.StrToInt(sysStatsPrev["physical write"])
	physicalWriteIoRequestsPersecond := utils.StrToInt(sysStats["physical write total IO requests"]) - utils.StrToInt(sysStatsPrev["physical write total IO requests"])
	physicalReadIoRequestsPersecond := utils.StrToInt(sysStats["physical read total IO requests"]) - utils.StrToInt(sysStatsPrev["physical read total IO requests"])
	dbBlockChangesPersecond := utils.StrToInt(sysStats["db block changes"]) - utils.StrToInt(sysStatsPrev["db block changes"])
	osCpuWaitTime := utils.StrToInt(sysStats["OS CPU Qt wait time"]) - utils.StrToInt(sysStatsPrev["OS CPU Qt wait time"])
	logonsCumulative := utils.StrToInt(sysStats["logons cumulative"]) - utils.StrToInt(sysStatsPrev["logons cumulative"])
	logonsCurrent := utils.StrToInt(sysStats["logons current"])
	openedCursorsPersecond := utils.StrToInt(sysStats["opened cursors cumulative"]) - utils.StrToInt(sysStatsPrev["opened cursors cumulative"])
	openedCursorsCurrent := utils.StrToInt(sysStats["opened cursors current"])
	userCommitsPersecond := utils.StrToInt(sysStats["user commits"]) - utils.StrToInt(sysStatsPrev["user commits"])
	userRollbacksPersecond := utils.StrToInt(sysStats["user rollbacks"]) - utils.StrToInt(sysStatsPrev["user rollbacks"])
	userCallsPersecond := utils.StrToInt(sysStats["user calls"]) - utils.StrToInt(sysStatsPrev["user calls"])
	dbBlockGetsPersecond := utils.StrToInt(sysStats["db block gets"]) - utils.StrToInt(sysStatsPrev["db block gets"])

	var (
		dgStats = -1
		dgDelay = -1
	)

	if databaseRole == "STANDBY" {
		dataGuardDetail, _ := oracle.QueryAll(oradb, queryDataGuardSQL)
		if len(dataGuardDetail) > 0 {
			dgStats = 1
			dgDelay = dataGuardDetail[0]["SECONDS"].(int)
		} else {
			dgStats = 0
			dgDelay = -1
		}
	}

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

	// event = map[string]interface{}{
	// 	"event_uuid":   tool.GetUUID(),
	// 	"event_time":   tool.GetNowTime(),
	// 	"event_type":   eventType,
	// 	"event_group":  eventGroup,
	// 	"event_entity": eventEntity,
	// 	"event_key":    "uptime",
	// 	"event_value":  utils.Int64ToDecimal(uptime),
	// 	"event_tag":    "",
	// 	"event_unit":   "",
	// 	"event_detail": "",
	// }
	// events = append(events, event)

	sessionTotalDetail, _ := oracle.QueryAll(oradb, querySessionTotalSQL)
	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "sessionTotal",
		"event_value":  utils.IntToDecimal(len(sessionTotalDetail)),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": utils.MapToStr(sessionTotalDetail),
	}
	events = append(events, event)

	sessionActiveDetail, _ := oracle.QueryAll(oradb, querySessionActiveSQL)
	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "sessionActive",
		"event_value":  utils.IntToDecimal(len(sessionActiveDetail)),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": utils.MapToStr(sessionActiveDetail),
	}
	events = append(events, event)

	sessionWaitDetail, _ := oracle.QueryAll(oradb, querySessionWaitSQL)
	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "sessionWait",
		"event_value":  utils.IntToDecimal(len(sessionWaitDetail)),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": utils.MapToStr(sessionWaitDetail),
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "sessionLogicalReadsPersecond",
		"event_value":  utils.IntToDecimal(sessionLogicalReadsPersecond),
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
		"event_key":    "physicalReadsPersecond",
		"event_value":  utils.IntToDecimal(physicalReadsPersecond),
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
		"event_key":    "physicalWritePersecond",
		"event_value":  utils.IntToDecimal(physicalWritePersecond),
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
		"event_key":    "physicalWriteIoRequestsPersecond",
		"event_value":  utils.IntToDecimal(physicalWriteIoRequestsPersecond),
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
		"event_key":    "physicalReadIoRequestsPersecond",
		"event_value":  utils.IntToDecimal(physicalReadIoRequestsPersecond),
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
		"event_key":    "dbBlockChangesPersecond",
		"event_value":  utils.IntToDecimal(dbBlockChangesPersecond),
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
		"event_key":    "osCpuWaitTime",
		"event_value":  utils.IntToDecimal(osCpuWaitTime),
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
		"event_key":    "logonsCumulative",
		"event_value":  utils.IntToDecimal(logonsCumulative),
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
		"event_key":    "logonsCurrent",
		"event_value":  utils.IntToDecimal(logonsCurrent),
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
		"event_key":    "openedCursorsPersecond",
		"event_value":  utils.IntToDecimal(openedCursorsPersecond),
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
		"event_key":    "openedCursorsCurrent",
		"event_value":  utils.IntToDecimal(openedCursorsCurrent),
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
		"event_key":    "userCommitsPersecond",
		"event_value":  utils.IntToDecimal(userCommitsPersecond),
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
		"event_key":    "userRollbacksPersecond",
		"event_value":  utils.IntToDecimal(userRollbacksPersecond),
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
		"event_key":    "userCallsPersecond",
		"event_value":  utils.IntToDecimal(userCallsPersecond),
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
		"event_key":    "dbBlockGetsPersecond",
		"event_value":  utils.IntToDecimal(dbBlockGetsPersecond),
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
		"event_key":    "dbBlockGetsPersecond",
		"event_value":  utils.IntToDecimal(dbBlockGetsPersecond),
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
		"event_key":    "dgStats",
		"event_value":  utils.IntToDecimal(dgStats),
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
		"event_key":    "dgDelay",
		"event_value":  utils.IntToDecimal(dgDelay),
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
		record.Version = version
		record.Uptime = uptime
		record.Connect = connect
		record.Role = databaseRole.(string)
		record.Session = len(sessionTotalDetail)
		record.Active = len(sessionActiveDetail)
		record.Wait = len(sessionWaitDetail)
		record.Qps = -1
		record.Tps = -1
		record.Repl = dgStats
		record.Delay = dgDelay

		result := database.DB.Create(&record)
		if result.Error != nil {
			log.Logger.Error("Insert Error:" + result.Error.Error())
		}

	} else {
		var record model.EventGlobal
		record.Version = version
		record.Uptime = uptime
		record.Version = version
		record.Uptime = uptime
		record.Connect = connect
		record.Role = databaseRole.(string)
		record.Session = len(sessionTotalDetail)
		record.Active = len(sessionActiveDetail)
		record.Wait = len(sessionWaitDetail)
		record.Qps = -1
		record.Tps = -1
		record.Repl = dgStats
		record.Delay = dgDelay
		//gin里面如果更新为0则字段不会更新，可以使用select指定更新字段解决
		result := database.DB.Model(&record).Select("version", "uptime", "role", "session", "active", "wait", "qps", "tps", "repl", "delay").Omit("id").Where("host=?", host).Where("port=?", port).Updates(&record)
		if result.Error != nil {
			log.Logger.Error("Update Error:" + result.Error.Error())
		}
	}

}
