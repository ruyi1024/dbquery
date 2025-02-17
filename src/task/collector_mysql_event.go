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
	"dbmcloud/src/libary/mysql"
	"dbmcloud/src/libary/tool"
	"dbmcloud/src/model"
	"dbmcloud/src/mq"
	"dbmcloud/src/utils"
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

func init() {
	go collectorMysqlEventTask()
}

var (
	queryVersionSQL   = "select version() as version limit 1"
	queryStatusSQL    = "show global status"
	queryVariablesSQL = "show global variables"
	queryReplSQL      = "show slave status"
)

func collectorMysqlEventTask() {

	time.Sleep(time.Second * time.Duration(30))
	var db = database.DB
	var record model.TaskOption
	db.Select("crontab").Where("task_key=?", "collector_mysql_event").Take(&record)
	c := cron.New()
	c.AddFunc(record.Crontab, func() {
		db.Select("enable").Where("task_key=?", "collector_mysql_event").Take(&record)
		if record.Enable == 1 {
			db.Model(model.TaskHeartbeat{}).Where("heartbeat_key='collector_mysql_event'").Updates(map[string]interface{}{"heartbeat_time": time.Now().Format("2006-01-02 15:04:05.999")})
			doCollectorMysqlEventTask()
			db.Model(model.TaskHeartbeat{}).Where("heartbeat_key='collector_mysql_event'").Updates(map[string]interface{}{"heartbeat_end_time": time.Now().Format("2006-01-02 15:04:05.999")})
		}
	})
	c.Start()
}

func doCollectorMysqlEventTask() {
	var db = database.DB
	var dataList []model.Datasource
	result := db.Where("enable=1").Where("type in ? ", strings.Split("MySQL,MariaDB,GreatSQL", ",")).Order("type asc").Find(&dataList)
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

		startCollectMysqlEvent(datasourceName, datasourceType, env, host, port, user, origPass)

	}

}

func startCollectMysqlEvent(datasourceName, datasourceType, env, host, port, user, origPass string) {
	eventEntity := fmt.Sprintf("%s:%s", host, port)
	eventType := datasourceType
	eventGroup := env
	var connect int = 1

	db, err := mysql.Connect(host, port, user, origPass, "")
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
	defer db.Close()
	row := db.QueryRow(queryVersionSQL)
	var version string
	if err := row.Scan(&version); err != nil {
		log.Error(fmt.Sprintf("Can't scan mysql version on %s:%s, %s", host, port, err))
		return
	}

	rows, err := db.Query(queryStatusSQL)
	if err != nil {
		log.Error(fmt.Sprintf("Can't query mysql status on %s:%s, %s", host, port, err))
		return
	}

	defer rows.Close()
	var key, value string
	globalStatusPrev := make(map[string]string)
	for rows.Next() {
		err := rows.Scan(&key, &value)
		if err != nil {
			log.Error(fmt.Sprintf("Can't scan mysql status on %s:%s, %s", host, port, err))
			return
		}
		globalStatusPrev[key] = value
	}

	time.Sleep(time.Duration(1) * time.Second)

	rows, err = db.Query(queryStatusSQL)
	if err != nil {
		log.Error(fmt.Sprintf("Can't query mysql status on %s:%s, %s", host, port, err))
		return
	}
	defer rows.Close()
	globalStatus := make(map[string]string)
	for rows.Next() {
		err := rows.Scan(&key, &value)
		if err != nil {
			log.Error(fmt.Sprintf("Can't scan mysql status on %s:%s, %s", host, port, err))
			return
		}
		globalStatus[key] = value
	}

	//fmt.Println(globalStatus)
	rows, err = db.Query(queryVariablesSQL)
	if err != nil {
		log.Error(fmt.Sprintf("Can't query mysql variables on %s:%s, %s", host, port, err))
		return
	}
	defer rows.Close()
	globalVariables := make(map[string]string)
	for rows.Next() {
		err := rows.Scan(&key, &value)
		if err != nil {
			log.Error(fmt.Sprintf("Can't scan mysql variables on %s:%s, %s", host, port, err))
			return
		}
		globalVariables[key] = value
	}

	//variables
	maxConnections := utils.StrToInt(globalVariables["max_connections"])
	openFilesLimit := utils.StrToInt(globalVariables["open_files_limit"])
	tableOpenCache := utils.StrToInt(globalVariables["table_open_cache"])
	//status
	uptime := utils.StrToInt64(globalStatus["Uptime"])
	openFiles := utils.StrToInt(globalStatus["open_files"])
	openTables := utils.StrToInt(globalStatus["Open_tables"])
	threadsConnected := utils.StrToInt(globalStatus["Threads_connected"])
	//threadsRunning := utils.StrToInt(globalStatus["Threads_running"])
	threadsCreated := utils.StrToInt(globalStatus["Threads_created"])
	threadsCached := utils.StrToInt(globalStatus["Threads_cached"])
	connections := utils.StrToInt(globalStatus["Connections"])
	abortedClients := utils.StrToInt(globalStatus["Aborted_clients"])
	abortedConnects := utils.StrToInt(globalStatus["Aborted_connects"])

	bytesReceived := utils.StrToInt(globalStatus["Bytes_received"]) - utils.StrToInt(globalStatusPrev["Bytes_received"])
	bytesSent := utils.StrToInt(globalStatus["Bytes_sent"]) - utils.StrToInt(globalStatusPrev["Bytes_sent"])
	comSelect := utils.StrToInt(globalStatus["Com_select"]) - utils.StrToInt(globalStatusPrev["Com_select"])
	comInsert := utils.StrToInt(globalStatus["Com_insert"]) - utils.StrToInt(globalStatusPrev["Com_insert"])
	comUpdate := utils.StrToInt(globalStatus["Com_update"]) - utils.StrToInt(globalStatusPrev["Com_update"])
	comDelete := utils.StrToInt(globalStatus["Com_delete"]) - utils.StrToInt(globalStatusPrev["Com_delete"])
	comCommit := utils.StrToInt(globalStatus["Com_commit"]) - utils.StrToInt(globalStatusPrev["Com_commit"])
	comRollback := utils.StrToInt(globalStatus["Com_rollback"]) - utils.StrToInt(globalStatusPrev["Com_rollback"])
	questions := utils.StrToInt(globalStatus["Questions"]) - utils.StrToInt(globalStatusPrev["Questions"])
	queries := utils.StrToInt(globalStatus["Queries"]) - utils.StrToInt(globalStatusPrev["Queries"])
	slowQueries := utils.StrToInt(globalStatus["Slow_queries"])

	//innodb status
	innodbPagesCreated := utils.StrToInt(globalStatus["Innodb_pages_created"])
	innodbPagesRead := utils.StrToInt(globalStatus["Innodb_pages_read"])
	innodbPagesWritten := utils.StrToInt(globalStatus["Innodb_pages_written"])
	innodbRowLockCurrentWaits := utils.StrToInt(globalStatus["Innodb_row_lock_current_waits"])
	innodbBufferPoolReadRequests := utils.StrToInt(globalStatus["Innodb_buffer_pool_read_requests"]) - utils.StrToInt(globalStatusPrev["Innodb_buffer_pool_read_requests"])
	innodbBufferPoolWriteRequests := utils.StrToInt(globalStatus["Innodb_buffer_pool_write_requests"]) - utils.StrToInt(globalStatusPrev["Innodb_buffer_pool_write_requests"])
	innodbRowsDeleted := utils.StrToInt(globalStatus["Innodb_rows_deleted"]) - utils.StrToInt(globalStatusPrev["Innodb_rows_deleted"])
	innodbRowsInserted := utils.StrToInt(globalStatus["Innodb_rows_inserted"]) - utils.StrToInt(globalStatusPrev["Innodb_rows_inserted"])
	innodbRowsRead := utils.StrToInt(globalStatus["Innodb_rows_read"]) - utils.StrToInt(globalStatusPrev["Innodb_rows_read"])
	innodbRowsUpdated := utils.StrToInt(globalStatus["Innodb_rows_updated"]) - utils.StrToInt(globalStatusPrev["Innodb_rows_updated"])

	threadsRunningDetail, _ := mysql.QueryAll(db, "select * from information_schema.processlist where db is not null and db != 'information_schema' and  command !='Sleep' order by time desc;")
	threadsWaitDetail, _ := mysql.QueryAll(db, "select * from information_schema.processlist where db is not null and db != 'information_schema' and  command !='Sleep' and time >0 order by time desc;")
	threadsConnectedDetail, _ := mysql.QueryAll(db, "select substring_index(host,':',1) 'IP',User,DB,count(*) Count from information_schema.processlist group by ip,user,db order by count desc;")
	longQueryDetail, _ := mysql.QueryAll(db, "select host,db,user,command,info from information_schema.processlist where db is not null and db != 'information_schema' and  command !='Sleep' and time > 10 order by time desc;")
	activeTrxDetail, _ := mysql.QueryAll(db, "select trx_id,trx_mysql_thread_id,trx_started,trx_isolation_level,trx_state,trx_rows_locked,trx_lock_structs,trx_tables_locked,trx_unique_checks,trx_is_read_only,trx_query from information_schema.INNODB_TRX;")
	longTrxDetail, _ := mysql.QueryAll(db, "select trx_id,trx_mysql_thread_id,trx_started,trx_isolation_level,trx_state,trx_rows_locked,trx_lock_structs,trx_tables_locked,trx_unique_checks,trx_is_read_only,trx_query from information_schema.INNODB_TRX where timestampdiff(second, trx_started,now())>10;")

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
		"event_key":    "maxConnections",
		"event_value":  utils.IntToDecimal(maxConnections),
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
		"event_key":    "openFilesLimit",
		"event_value":  utils.IntToDecimal(openFilesLimit),
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
		"event_key":    "tableOpenCache",
		"event_value":  utils.IntToDecimal(tableOpenCache),
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
		"event_value":  utils.Int64ToDecimal(uptime),
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
		"event_key":    "openFiles",
		"event_value":  utils.IntToDecimal(openFiles),
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
		"event_key":    "openTables",
		"event_value":  utils.IntToDecimal(openTables),
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
		"event_key":    "connections",
		"event_value":  utils.IntToDecimal(connections),
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
		"event_key":    "abortedClients",
		"event_value":  utils.IntToDecimal(abortedClients),
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
		"event_key":    "threadsCreated",
		"event_value":  utils.IntToDecimal(threadsCreated),
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
		"event_key":    "abortedConnects",
		"event_value":  utils.IntToDecimal(abortedConnects),
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
		"event_key":    "threadsCached",
		"event_value":  utils.IntToDecimal(threadsCached),
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
		"event_key":    "threadsConnected",
		"event_value":  utils.IntToDecimal(threadsConnected),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": utils.MapToStr(threadsConnectedDetail),
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "threadsRunning",
		"event_value":  utils.IntToDecimal(len(threadsRunningDetail)),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": utils.MapToStr(threadsRunningDetail),
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "threadsWait",
		"event_value":  utils.IntToDecimal(len(threadsWaitDetail)),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": utils.MapToStr(threadsWaitDetail),
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "longQuery",
		"event_value":  utils.IntToDecimal(len(longQueryDetail)),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": utils.MapToStr(longQueryDetail),
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "activeTrx",
		"event_value":  utils.IntToDecimal(len(activeTrxDetail)),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": utils.MapToStr(activeTrxDetail),
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "longTrx",
		"event_value":  utils.IntToDecimal(len(longTrxDetail)),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": utils.MapToStr(longTrxDetail),
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "slowQueries",
		"event_value":  utils.IntToDecimal(slowQueries),
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
		"event_key":    "queries",
		"event_value":  utils.IntToDecimal(queries),
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
		"event_key":    "questions",
		"event_value":  utils.IntToDecimal(questions),
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
		"event_key":    "bytesReceived",
		"event_value":  utils.IntToDecimal(bytesReceived / 1024),
		"event_tag":    "",
		"event_unit":   "Kb",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "bytesSent",
		"event_value":  utils.IntToDecimal(bytesSent / 1024),
		"event_tag":    "",
		"event_unit":   "Kb",
		"event_detail": "",
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "comSelect",
		"event_value":  utils.IntToDecimal(comSelect),
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
		"event_key":    "comInsert",
		"event_value":  utils.IntToDecimal(comInsert),
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
		"event_key":    "comUpdate",
		"event_value":  utils.IntToDecimal(comUpdate),
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
		"event_key":    "comDelete",
		"event_value":  utils.IntToDecimal(comDelete),
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
		"event_key":    "comCommit",
		"event_value":  utils.IntToDecimal(comCommit),
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
		"event_key":    "comRollback",
		"event_value":  utils.IntToDecimal(comRollback),
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
		"event_key":    "innodbPagesCreated",
		"event_value":  utils.IntToDecimal(innodbPagesCreated),
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
		"event_key":    "innodbPagesRead",
		"event_value":  utils.IntToDecimal(innodbPagesRead),
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
		"event_key":    "innodbPagesWritten",
		"event_value":  utils.IntToDecimal(innodbPagesWritten),
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
		"event_key":    "innodbRowLockCurrentWaits",
		"event_value":  utils.IntToDecimal(innodbRowLockCurrentWaits),
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
		"event_key":    "innodbBufferPoolReadRequests",
		"event_value":  utils.IntToDecimal(innodbBufferPoolReadRequests),
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
		"event_key":    "innodbBufferPoolWriteRequests",
		"event_value":  utils.IntToDecimal(innodbBufferPoolWriteRequests),
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
		"event_key":    "innodbRowsDeleted",
		"event_value":  utils.IntToDecimal(innodbRowsDeleted),
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
		"event_key":    "innodbRowsInserted",
		"event_value":  utils.IntToDecimal(innodbRowsInserted),
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
		"event_key":    "innodbRowsRead",
		"event_value":  utils.IntToDecimal(innodbRowsRead),
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
		"event_key":    "innodbRowsUpdated",
		"event_value":  utils.IntToDecimal(innodbRowsUpdated),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": "",
	}
	events = append(events, event)

	var (
		role       string = "Master" //1:master,2:slave
		replStatus int    = -1
		replDelay  int    = -1
	)
	replInfoList, _ := mysql.QueryAll(db, queryReplSQL)
	if len(replInfoList) > 0 {
		role = "Slave"
		for _, replInfo := range replInfoList {
			// autoPosition = replInfo["Auto_Position"]
			// masterHost = replInfo["Master_Host"]
			// masterPort = replInfo["Master_Port"]
			// masterUser = replInfo["Master_User"]
			secondBehine := replInfo["Seconds_Behine_Master"]
			IORunning := replInfo["Slave_IO_Running"]
			SQLRunning := replInfo["Slave_SQL_Running"]
			if IORunning == "Yes" && SQLRunning == "Yes" {
				replStatus = 1
			} else {
				replStatus = 0
			}
			if secondBehine != "NULL" && secondBehine == nil {
				replDelay = secondBehine.(int)
			}

			event = map[string]interface{}{
				"event_uuid":   tool.GetUUID(),
				"event_time":   tool.GetNowTime(),
				"event_type":   eventType,
				"event_group":  eventGroup,
				"event_entity": eventEntity,
				"event_key":    "replStatus",
				"event_value":  utils.IntToDecimal(replStatus),
				"event_tag":    "",
				"event_unit":   "",
				"event_detail": utils.MapToStr(replInfoList),
			}
			events = append(events, event)

			event = map[string]interface{}{
				"event_uuid":   tool.GetUUID(),
				"event_time":   tool.GetNowTime(),
				"event_type":   eventType,
				"event_group":  eventGroup,
				"event_entity": eventEntity,
				"event_key":    "replDelay",
				"event_value":  utils.IntToDecimal(replDelay),
				"event_tag":    "",
				"event_unit":   "秒",
				"event_detail": utils.MapToStr(replInfoList),
			}
			events = append(events, event)

		}
	}

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
		record.Role = role
		record.Session = threadsConnected
		record.Active = len(threadsRunningDetail)
		record.Wait = len(threadsWaitDetail)
		record.Qps = queries
		record.Tps = comInsert + comUpdate + comDelete + comCommit
		record.Repl = replStatus
		record.Delay = replDelay

		result := database.DB.Create(&record)
		if result.Error != nil {
			log.Logger.Error("Insert Error:" + result.Error.Error())
		}

	} else {
		var record model.EventGlobal
		record.Version = version
		record.Uptime = uptime
		record.Connect = connect
		record.Role = role
		record.Session = threadsConnected
		record.Active = len(threadsRunningDetail)
		record.Wait = len(threadsWaitDetail)
		record.Qps = queries
		record.Tps = comInsert + comUpdate + comDelete + comCommit
		record.Repl = replStatus
		record.Delay = replDelay
		//gin里面如果更新为0则字段不会更新，可以使用select指定更新字段解决
		result := database.DB.Model(&record).Select("version", "uptime", "role", "session", "active", "wait", "qps", "tps", "repl", "delay").Omit("id").Where("host=?", host).Where("port=?", port).Updates(&record)
		if result.Error != nil {
			log.Logger.Error("Update Error:" + result.Error.Error())
		}
	}

}
