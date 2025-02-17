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
	"dbmcloud/src/libary/mssql"
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
	go collectorSQLServerEventTask()
}

var (
	queryVersionInfoSQL          = "SELECT @@VERSION AS [SQL Server and OS Version Info];"
	queryVariableSQL             = "select @@VERSION as version,@@MAX_CONNECTIONS as max_connections,@@LOCK_TIMEOUT as lock_timeout,@@TRANCOUNT as trancount,@@CONNECTIONS as connections,@@PACK_RECEIVED as pack_received,@@PACK_SENT as pack_sent,@@PACKET_ERRORS as packet_errors,@@ROWCOUNT as row_count,@@CPU_BUSY as cpu_busy,@@IO_BUSY as io_busy,@@CURSOR_ROWS as cursor_rows,@@TOTAL_WRITE as total_write,@@TOTAL_READ as total_read,@@TOTAL_ERRORS as total_errors"
	queryUptimeSQL               = "SELECT crdate startup_time,GETDATE() AS time_now,DATEDIFF(mi,crdate,GETDATE())*60 AS uptime FROM master..sysdatabases WHERE name = 'tempdb';"
	queryOsSysInfoSQL            = "SELECT cpu_count AS [Logical CPU Count],cpu_count/hyperthread_ratio AS [Physical CPU Count],physical_memory_kb/1024 AS [Physical Memory (MB)], sqlserver_start_time FROM master.sys.dm_os_sys_info WITH (NOLOCK) OPTION (RECOMPILE);"
	queryProcessSQL              = "SELECT COUNT(*) as count FROM [Master].[dbo].[SYSPROCESSES] WHERE [DBID] IN ( SELECT  [dbid] FROM [Master].[dbo].[SYSDATABASES]);"
	queryProcessRunningDetailSQL = "SELECT * FROM [Master].[dbo].[SYSPROCESSES] WHERE [DBID] IN ( SELECT  [dbid] FROM [Master].[dbo].[SYSDATABASES])  AND  status !='SLEEPING' AND status !='BACKGROUND'; "
	queryProcessWaitDetailSQL    = "SELECT * FROM [Master].[dbo].[SYSPROCESSES] WHERE [DBID] IN ( SELECT  [dbid] FROM [Master].[dbo].[SYSDATABASES])  AND  status ='SUSPENDED' AND waittime >1;"
)

func collectorSQLServerEventTask() {
	time.Sleep(time.Second * time.Duration(50))
	var db = database.DB
	var record model.TaskOption
	db.Select("crontab").Where("task_key=?", "collector_sqlserver_event").Take(&record)
	c := cron.New()
	c.AddFunc(record.Crontab, func() {
		db.Select("enable").Where("task_key=?", "collector_sqlserver_event").Take(&record)
		if record.Enable == 1 {
			db.Model(model.TaskHeartbeat{}).Where("heartbeat_key='collector_sqlserver_event'").Updates(map[string]interface{}{"heartbeat_time": time.Now().Format("2006-01-02 15:04:05.999")})
			doCollectorSQLServerEventTask()
			db.Model(model.TaskHeartbeat{}).Where("heartbeat_key='collector_sqlserver_event'").Updates(map[string]interface{}{"heartbeat_end_time": time.Now().Format("2006-01-02 15:04:05.999")})
		}
	})
	c.Start()
}

func doCollectorSQLServerEventTask() {
	var db = database.DB
	var dataList []model.Datasource
	result := db.Where("enable=1").Where("type = ? ", "SQLServer").Order("type asc").Find(&dataList)
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

		startCollectSQLServerEvent(datasourceName, datasourceType, env, host, port, user, origPass)

	}

}

func startCollectSQLServerEvent(datasourceName, datasourceType, env, host, port, user, origPass string) {
	eventEntity := fmt.Sprintf("%s:%s", host, port)
	eventType := datasourceType
	eventGroup := env
	var connect int = 1

	msdb, err := mssql.Connect(host, port, user, origPass, "master")

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
	defer msdb.Close()

	queryUptime, _ := mssql.QueryAll(msdb, queryUptimeSQL)
	//startupTime := queryUptime[0]["startup_time"].(time.Time)
	uptime := queryUptime[0]["uptime"].(int64)

	queryProcess, _ := mssql.QueryAll(msdb, queryProcessSQL)
	queryProcessRunningDetail, _ := mssql.QueryAll(msdb, queryProcessRunningDetailSQL)
	queryProcessWaitDetail, _ := mssql.QueryAll(msdb, queryProcessWaitDetailSQL)
	process := queryProcess[0]["count"].(int64)
	processRunning := len(queryProcessRunningDetail)
	processWait := len(queryProcessWaitDetail)

	queryVariablesPrev, _ := mssql.QueryAll(msdb, queryVariableSQL)
	time.Sleep(time.Duration(1) * time.Second)
	queryVariables, _ := mssql.QueryAll(msdb, queryVariableSQL)
	fmt.Println(queryVariables)
	// versionInfo := queryVariables[0]["SQL Server and OS Version Info"]
	// var version string
	// if versionInfo != nil {
	// 	version = versionInfo.(string)
	// }
	version := strings.Replace(strings.Split(queryVariables[0]["version"].(string), "-")[0], "Microsoft SQL Server ", "", -1)
	lockTimeOut := queryVariables[0]["lock_timeout"].(int64)
	tranCount := queryVariables[0]["trancount"].(int64)
	maxConnections := queryVariables[0]["max_connections"].(int64)
	packReceived := queryVariables[0]["pack_received"].(int64)
	packSent := queryVariables[0]["pack_sent"].(int64)
	packetErrors := queryVariables[0]["packet_errors"].(int64)
	rowCount := queryVariables[0]["row_count"].(int64)
	cpuBusy := queryVariables[0]["cpu_busy"].(int64)
	ioBusy := queryVariables[0]["io_busy"].(int64)
	cursorRows := queryVariables[0]["cursor_rows"].(int64)
	currentWrite := queryVariables[0]["total_write"].(int64) - queryVariablesPrev[0]["total_write"].(int64)
	currentRead := queryVariables[0]["total_read"].(int64) - queryVariablesPrev[0]["total_read"].(int64)
	currentError := queryVariables[0]["total_errors"].(int64) - queryVariablesPrev[0]["total_errors"].(int64)
	totalErrors := queryVariables[0]["total_errors"].(int64)

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
		"event_value":  utils.Int64ToDecimal(maxConnections),
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
		"event_key":    "process",
		"event_value":  utils.Int64ToDecimal(process),
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
		"event_key":    "processRunning",
		"event_value":  utils.IntToDecimal(processRunning),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": utils.MapToStr(queryProcessRunningDetail),
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "processWait",
		"event_value":  utils.IntToDecimal(processWait),
		"event_tag":    "",
		"event_unit":   "",
		"event_detail": utils.MapToStr(queryProcessWaitDetail),
	}
	events = append(events, event)

	event = map[string]interface{}{
		"event_uuid":   tool.GetUUID(),
		"event_time":   tool.GetNowTime(),
		"event_type":   eventType,
		"event_group":  eventGroup,
		"event_entity": eventEntity,
		"event_key":    "lockTimeOut",
		"event_value":  utils.Int64ToDecimal(lockTimeOut),
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
		"event_key":    "tranCount",
		"event_value":  utils.Int64ToDecimal(tranCount),
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
		"event_key":    "packReceived",
		"event_value":  utils.Int64ToDecimal(packReceived),
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
		"event_key":    "packSent",
		"event_value":  utils.Int64ToDecimal(packSent),
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
		"event_key":    "packetErrors",
		"event_value":  utils.Int64ToDecimal(packetErrors),
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
		"event_key":    "rowCount",
		"event_value":  utils.Int64ToDecimal(rowCount),
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
		"event_key":    "cpuBusy",
		"event_value":  utils.Int64ToDecimal(cpuBusy),
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
		"event_key":    "ioBusy",
		"event_value":  utils.Int64ToDecimal(ioBusy),
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
		"event_key":    "cursorRows",
		"event_value":  utils.Int64ToDecimal(cursorRows),
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
		"event_key":    "currentWrite",
		"event_value":  utils.Int64ToDecimal(currentWrite),
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
		"event_key":    "currentRead",
		"event_value":  utils.Int64ToDecimal(currentRead),
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
		"event_key":    "currentError",
		"event_value":  utils.Int64ToDecimal(currentError),
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
		"event_key":    "totalErrors",
		"event_value":  utils.Int64ToDecimal(totalErrors),
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
		record.Role = "-1"
		record.Session = int(process)
		record.Active = processRunning
		record.Wait = processWait
		record.Qps = int(currentRead)
		record.Tps = int(tranCount)
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
		record.Connect = connect
		record.Role = "-1"
		record.Session = int(process)
		record.Active = processRunning
		record.Wait = processWait
		record.Qps = int(currentRead)
		record.Tps = int(tranCount)
		record.Repl = -1
		record.Delay = -1
		//gin里面如果更新为0则字段不会更新，可以使用select指定更新字段解决
		result := database.DB.Model(&record).Select("version", "uptime", "role", "session", "active", "wait", "qps", "tps", "repl", "delay").Omit("id").Where("host=?", host).Where("port=?", port).Updates(&record)
		if result.Error != nil {
			log.Logger.Error("Update Error:" + result.Error.Error())
		}
	}

}
