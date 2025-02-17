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
	"dbmcloud/src/libary/conv"
	"dbmcloud/src/libary/postgres"
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
	go collectorPostgresEventTask()
}

var (
	queryPgVersionSQL             = "select version() as version limit 1"
	queryStartTimeSQL             = "select pg_postmaster_start_time();"
	queryInRecoverySQL            = "select pg_is_in_recovery()"
	queryMaxConnectionsSQL        = "show max_connections"
	queryConnectionsSQL           = "select count(*) from pg_stat_activity"
	queryConnectionsDetailSQL     = "select client_addr,datname,usename,count(*) count from pg_stat_activity group by client_addr,datname,usename order by count desc"
	queryActiveDetailSQL          = "select datname,usename,application_name,client_addr,client_hostname,backend_start,query_start,state,query from pg_stat_activity where state='active' and pid<> pg_backend_pid() order by query_start asc"
	queryLongQueryDetailSQL       = "select datname,usename,application_name,client_addr,client_hostname,backend_start,query_start,state,query from pg_stat_activity where state='active' and pid<> pg_backend_pid() and  now()-query_start > interval '10 second' order by query_start asc"
	queryLongTransactionDetailSQL = "select datname,usename,application_name,client_addr,client_hostname,backend_start,xact_start,state,query from pg_stat_activity where  pid<> pg_backend_pid() and  now()-xact_start > interval '10 second' order by xact_start asc"
	queryLockDetailSQL            = "SELECT bl.pid AS blocked_pid, a.usename AS blocked_user, kl.pid AS blocking_pid, ka.usename AS blocking_user, a.query AS blocked_statement FROM pg_locks bl JOIN pg_stat_activity a ON a.pid = bl.pid  JOIN pg_locks kl ON kl.transactionid = bl.transactionid AND kl.pid != bl.pid  JOIN pg_stat_activity ka ON ka.pid = kl.pid WHERE NOT bl.granted;"
	queryWaitEventDetailSQL       = "select datname,usename,application_name,client_addr,client_hostname,state,wait_event_type,query from pg_stat_activity where state!='idle' and pid<> pg_backend_pid()  and wait_event_type is not null"
	queryPreparedXactDetailSQL    = "select * from pg_prepared_xacts"
	queryCheckpointSQL            = "SELECT (100 * checkpoints_req) / (checkpoints_timed + checkpoints_req) AS checkpoints_req_pct,pg_size_pretty(buffers_checkpoint * block_size / (checkpoints_timed + checkpoints_req)) AS avg_checkpoint_write,pg_size_pretty(block_size * (buffers_checkpoint + buffers_clean + buffers_backend)) AS total_written,\n100 * buffers_checkpoint / (buffers_checkpoint + buffers_clean + buffers_backend) AS checkpoint_write_pct, 100 * buffers_backend / (buffers_checkpoint + buffers_clean + buffers_backend) AS backend_write_pct FROM pg_stat_bgwriter,(SELECT cast(current_setting('block_size') AS integer) AS block_size) AS bs;"
	queryStatDatabaseSQL          = "select sum(xact_commit) xact_commit,sum(xact_rollback) xact_rollback,sum(tup_returned) tup_returned,sum(tup_fetched) tup_fetched,sum(tup_inserted) tup_inserted,sum(tup_updated) tup_updated,sum(tup_deleted) tup_deleted,sum(conflicts) conflicts,sum(deadlocks) deadlocks from pg_stat_database"
)

func collectorPostgresEventTask() {

	time.Sleep(time.Second * time.Duration(40))
	var db = database.DB
	var record model.TaskOption
	db.Select("crontab").Where("task_key=?", "collector_postgresql_event").Take(&record)
	c := cron.New()
	c.AddFunc(record.Crontab, func() {
		db.Select("enable").Where("task_key=?", "collector_postgresql_event").Take(&record)
		if record.Enable == 1 {
			db.Model(model.TaskHeartbeat{}).Where("heartbeat_key='collector_postgresql_event'").Updates(map[string]interface{}{"heartbeat_time": time.Now().Format("2006-01-02 15:04:05.999")})
			doCollectorPostgresEventTask()
			db.Model(model.TaskHeartbeat{}).Where("heartbeat_key='collector_postgresql_event'").Updates(map[string]interface{}{"heartbeat_end_time": time.Now().Format("2006-01-02 15:04:05.999")})
		}
	})
	c.Start()
}

func doCollectorPostgresEventTask() {
	var db = database.DB
	var dataList []model.Datasource
	result := db.Where("enable=1").Where("type = ? ", "PostgreSQL").Order("type asc").Find(&dataList)
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

		startCollectPostgresEvent(datasourceName, datasourceType, env, host, port, user, origPass)

	}

}

func startCollectPostgresEvent(datasourceName, datasourceType, env, host, port, user, origPass string) {
	eventEntity := fmt.Sprintf("%s:%s", host, port)
	eventType := datasourceType
	eventGroup := env
	var connect int = 1

	pgdb, err := postgres.Connect(host, port, user, origPass, "postgres")

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
	defer pgdb.Close()

	row := pgdb.QueryRow(queryVersionSQL)
	var version string
	if err := row.Scan(&version); err != nil {
		log.Error(fmt.Sprintf("Can't scan version on %s:%s, %s ", host, port, err))
		return
	}
	version = strings.Fields(version)[0] + "-" + strings.Fields(version)[1]

	queryStartTime, _ := postgres.QueryAll(pgdb, queryStartTimeSQL)
	startTime := queryStartTime[0]["pg_postmaster_start_time"]
	uptime := time.Now().Unix() - startTime.(time.Time).Unix()
	queryMaxConnections, _ := postgres.QueryAll(pgdb, queryMaxConnectionsSQL)
	maxConnections := utils.StrToInt(queryMaxConnections[0]["max_connections"].(string))

	queryConnections, _ := postgres.QueryAll(pgdb, queryConnectionsSQL)
	//fmt.Println(queryConnections)
	connections := queryConnections[0]["count"].(int64)
	queryConnectionsDetail, _ := postgres.QueryAll(pgdb, queryConnectionsDetailSQL)

	queryActiveDetail, _ := postgres.QueryAll(pgdb, queryActiveDetailSQL)
	queryWaitEventDetail, _ := postgres.QueryAll(pgdb, queryWaitEventDetailSQL)
	queryLongQueryDetail, _ := postgres.QueryAll(pgdb, queryLongQueryDetailSQL)
	queryLongTransactionDetail, _ := postgres.QueryAll(pgdb, queryLongTransactionDetailSQL)
	queryLockDetail, _ := postgres.QueryAll(pgdb, queryLockDetailSQL)
	queryPreparedXactDetail, _ := postgres.QueryAll(pgdb, queryPreparedXactDetailSQL)

	/*
		queryCheckkpoint, _ := postgres.QueryAll(pgdb, queryCheckpointSQL)
		fmt.Println(queryCheckkpoint)
		checkpointsReqPct := queryCheckkpoint[0]["checkpoints_req_pct"].(int64)
		avgCheckpointWrite := utils.StrToInt(strings.Replace(queryCheckkpoint[0]["avg_checkpoint_write"].(string), "bytes", "", -1)) / 1024
		totalWritten := utils.StrToInt(strings.Replace(queryCheckkpoint[0]["total_written"].(string), "kB", "", -1))
		checkpointWritePct := queryCheckkpoint[0]["checkpoint_write_pct"].(int64)
		backendWritePct := queryCheckkpoint[0]["backend_write_pct"].(int64)
	*/

	queryInRecovery, _ := postgres.QueryAll(pgdb, queryInRecoverySQL)
	inRecovery := queryInRecovery[0]["pg_is_in_recovery"].(bool)
	var role string
	if !inRecovery {
		role = "Primary"
	} else {
		role = "Standby"
	}

	statDatabasePrev, err := postgres.QueryAll(pgdb, queryStatDatabaseSQL)
	if err != nil {
		log.Error(fmt.Sprintf("Can't query pg_stat_database on %s:%s, %s", host, port, err))
		return
	}
	time.Sleep(time.Duration(1) * time.Second)
	statDatabase, err := postgres.QueryAll(pgdb, queryStatDatabaseSQL)
	if err != nil {
		log.Error(fmt.Sprintf("Can't query pg_stat_database on %s:%s, %s", host, port, err))
		return
	}
	tupFetched := conv.StrToInt(statDatabase[0]["tup_fetched"].(string)) - conv.StrToInt(statDatabasePrev[0]["tup_fetched"].(string))
	tupReturned := conv.StrToInt(statDatabase[0]["tup_returned"].(string)) - conv.StrToInt(statDatabasePrev[0]["tup_returned"].(string))
	tupInserted := conv.StrToInt(statDatabase[0]["tup_inserted"].(string)) - conv.StrToInt(statDatabasePrev[0]["tup_inserted"].(string))
	tupDeleted := conv.StrToInt(statDatabase[0]["tup_deleted"].(string)) - conv.StrToInt(statDatabasePrev[0]["tup_deleted"].(string))
	tupUpdated := conv.StrToInt(statDatabase[0]["tup_updated"].(string)) - conv.StrToInt(statDatabasePrev[0]["tup_updated"].(string))
	xactCommit := conv.StrToInt(statDatabase[0]["xact_commit"].(string)) - conv.StrToInt(statDatabasePrev[0]["xact_commit"].(string))
	xactRollback := conv.StrToInt(statDatabase[0]["xact_rollback"].(string)) - conv.StrToInt(statDatabasePrev[0]["xact_rollback"].(string))
	conflicts := conv.StrToInt(statDatabase[0]["conflicts"].(string)) - conv.StrToInt(statDatabasePrev[0]["conflicts"].(string))
	deadlocks := conv.StrToInt(statDatabase[0]["deadlocks"].(string)) - conv.StrToInt(statDatabasePrev[0]["deadlocks"].(string))

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
		"event_key":    "connections",
		"event_value":  utils.Int64ToDecimal(connections),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": utils.MapToStr(queryConnectionsDetail),
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "activeSQL",
		"event_value":  utils.IntToDecimal(len(queryActiveDetail)),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": utils.MapToStr(queryActiveDetail),
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "longQuery",
		"event_value":  utils.IntToDecimal(len(queryLongQueryDetail)),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": utils.MapToStr(queryLongQueryDetail),
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "longTransaction",
		"event_value":  utils.IntToDecimal(len(queryLongTransactionDetail)),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": utils.MapToStr(queryLongTransactionDetail),
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "waitEvent",
		"event_value":  utils.IntToDecimal(len(queryWaitEventDetail)),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": utils.MapToStr(queryWaitEventDetail),
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "locks",
		"event_value":  utils.IntToDecimal(len(queryLockDetail)),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": utils.MapToStr(queryLockDetail),
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "prepared_xacts",
		"event_value":  utils.IntToDecimal(len(queryPreparedXactDetail)),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": utils.MapToStr(queryPreparedXactDetail),
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "tup_fetched",
		"event_value":  utils.IntToDecimal(tupFetched),
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
		"event_key":    "tup_returned",
		"event_value":  utils.IntToDecimal(tupReturned),
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
		"event_key":    "tup_inserted",
		"event_value":  utils.IntToDecimal(tupInserted),
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
		"event_key":    "tup_deleted",
		"event_value":  utils.IntToDecimal(tupDeleted),
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
		"event_key":    "tup_updated",
		"event_value":  utils.IntToDecimal(tupUpdated),
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
		"event_key":    "xact_commit",
		"event_value":  utils.IntToDecimal(xactCommit),
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
		"event_key":    "xact_rollback",
		"event_value":  utils.IntToDecimal(xactRollback),
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
		"event_key":    "conflicts",
		"event_value":  utils.IntToDecimal(conflicts),
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
		"event_key":    "deadlocks",
		"event_value":  utils.IntToDecimal(deadlocks),
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
			"event_key":    "checkpointsReqPct",
			"event_value":  utils.Int64ToDecimal(checkpointsReqPct),
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
			"event_key":    "avgCheckpointWrite",
			"event_value":  utils.IntToDecimal(avgCheckpointWrite),
			"event_tag":    "",
			"event_unit":   "KB",
			"event_detail": "",
		}
		events = append(events, event)

		event = map[string]interface{}{
			"event_uuid":   tool.GetUUID(),
			"event_time":   tool.GetNowTime(),
			"event_type":   eventType,
			"event_group":  eventGroup,
			"event_entity": eventEntity,
			"event_key":    "totalWritten",
			"event_value":  utils.IntToDecimal(totalWritten),
			"event_tag":    "",
			"event_unit":   "KB",
			"event_detail": "",
		}
		events = append(events, event)

		event = map[string]interface{}{
			"event_uuid":   tool.GetUUID(),
			"event_time":   tool.GetNowTime(),
			"event_type":   eventType,
			"event_group":  eventGroup,
			"event_entity": eventEntity,
			"event_key":    "checkpointWritePct",
			"event_value":  utils.Int64ToDecimal(checkpointWritePct),
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
			"event_key":    "backendWritePct",
			"event_value":  utils.Int64ToDecimal(backendWritePct),
			"event_tag":    "",
			"event_unit":   "",
			"event_detail": "",
		}
		events = append(events, event)
	*/

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
		record.Session = int(connections)
		record.Active = len(queryActiveDetail)
		record.Wait = len(queryWaitEventDetail)
		record.Qps = tupFetched + tupDeleted + tupInserted + tupUpdated
		record.Tps = xactCommit + xactRollback
		record.Repl = -1
		record.Delay = -1

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
		record.Role = role
		record.Session = int(connections)
		record.Active = len(queryActiveDetail)
		record.Wait = len(queryWaitEventDetail)
		record.Qps = tupFetched + tupDeleted + tupInserted + tupUpdated
		record.Tps = xactCommit + xactRollback
		record.Repl = -1
		record.Delay = -1
		//gin里面如果更新为0则字段不会更新，可以使用select指定更新字段解决
		result := database.DB.Model(&record).Select("version", "uptime", "role", "session", "active", "wait", "qps", "tps", "repl", "delay").Omit("id").Where("host=?", host).Where("port=?", port).Updates(&record)
		if result.Error != nil {
			log.Logger.Error("Update Error:" + result.Error.Error())
		}
	}

}
