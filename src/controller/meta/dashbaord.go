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

package meta

import (
	"dbmcloud/src/database"
	"dbmcloud/src/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// var db = database.InitConnect()

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	//允许跨域访问
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func EventInfo(c *gin.Context) {
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Failed to set websocket upgrade: %+v", err)
		return
	}
	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		//fmt.Printf("websocket read message:  %s\n" ,msg)

		var d map[string]interface{} /*创建集合 */
		d = make(map[string]interface{})

		eventList, _ := database.QueryAll("select event_time,event_type,event_group,event_entity,event_key,event_value from events order by id desc limit 8")
		alarmList, _ := database.QueryAll("select alarm_title,alarm_level,event_type,event_entity,send_mail,send_phone,event_time,gmt_created from alarm_events order by id desc limit 8")
		eventCount, _ := database.QueryAll("select count(*) as count from events where event_time>date_sub(now(),interval 1 minute) limit 1")
		eventMinuteData, _ := database.QueryAll("select DATE_FORMAT(event_time,'%Y-%m-%d %H:%i') x ,count(*) as y from events where gmt_created>date_sub(now(),interval 15 minute) group by DATE_FORMAT(event_time,'%Y-%m-%d %H:%i')")

		alarmCount, _ := database.QueryAll("select count(*) as count from alarm_events where event_time>date_sub(now(),interval 10 second) limit 1")
		alarmMinuteData, _ := database.QueryAll("select DATE_FORMAT(event_time,'%Y-%m-%d %H:%i') x ,count(*) as y from alarm_events where gmt_created>date_sub(now(),interval 60 minute) group by DATE_FORMAT(event_time,'%Y-%m-%d %H:%i')")

		lastEventTime, _ := database.QueryAll("select event_time from events order by id desc limit 1")
		lastAlarmTime, _ := database.QueryAll("select event_time from alarm_events order by id desc limit 1")

		taskCount, _ := database.QueryAll("select count(*) as count from task_run where gmt_create>date_sub(now(),interval 1 minute) limit 1")
		taskHourCount, _ := database.QueryAll("select count(*) as count from task_run where gmt_create>date_sub(now(),interval 1 hour) limit 1")
		taskMinuteData, _ := database.QueryAll("select DATE_FORMAT(gmt_create,'%Y-%m-%d %H:%i') x ,count(*) as y from task_run where gmt_create>date_sub(now(),interval 60 minute) group by DATE_FORMAT(gmt_create,'%Y-%m-%d %H:%i')")

		disEvents, _ := database.QueryAll("select count(distinct event_entity,event_key) as count from events where gmt_created>date_sub(now(),INTERVAL 5 MINUTE) limit 1")
		curAlarms, _ := database.QueryAll("select count(distinct event_entity,event_key) as count from alarm_events where gmt_created>= DATE_FORMAT(NOW(), '%Y-%m-%d 00:00:00') limit 1")
		healthPct := 0.00
		if utils.StrToInt(disEvents[0]["count"].(string)) != 0 {
			healthPct = 1 - utils.StrToFloat64(curAlarms[0]["count"].(string))/utils.StrToFloat64(disEvents[0]["count"].(string))
		}

		alarmPieData, _ := database.QueryAll("select CONCAT(event_type,':',event_key) type,count(*) value from alarm_events where gmt_created>= DATE_FORMAT(NOW(), '%Y-%m-%d 00:00:00') group by type order by value desc limit 50")

		pieDataList := make([]map[string]interface{}, 0)
		for _, item := range alarmPieData {
			pieData := make(map[string]interface{})
			//pieData[item["type"].(string)] = utils.StrToInt(item["value"].(string))
			pieData["type"] = item["type"].(string)
			pieData["value"] = utils.StrToInt(item["value"].(string))
			pieDataList = append(pieDataList, pieData)
		}

		nodeCount, _ := database.QueryAll("select count(*) as count from meta_nodes limit 1;")
		taskNextCount, _ := database.QueryAll("select count(*) as count from task where  next_time > date_sub(now(), interval 3 minute) limit 1;")
		taskFailCount, _ := database.QueryAll("select count(*) as count from task_run where gmt_create > date_sub(now(), interval 3 minute) and run_status='failed' limit 1;")
		sqlModeUnSupportCount, _ := database.QueryAll("select count(*) as count from information_schema.global_variables where variable_name='sql_mode' and (variable_value like '%ONLY_FULL_GROUP_BY%' or  variable_value like '%only_full_group_by%')limit 1;")

		/* map插入key - value对,各个国家对应的首都 */
		d["receiveMessage"] = msg
		d["eventList"] = eventList
		d["alarmList"] = alarmList
		d["alarmCount"] = alarmCount[0]["count"]
		d["eventCount"] = eventCount[0]["count"]
		d["taskCount"] = taskCount[0]["count"]
		d["taskHourCount"] = taskHourCount[0]["count"]
		d["alarmMinuteData"] = alarmMinuteData
		d["eventMinuteData"] = eventMinuteData
		d["taskMinuteData"] = taskMinuteData
		d["nodeCount"] = nodeCount[0]["count"]
		d["taskNextCount"] = taskNextCount[0]["count"]
		d["disEvents"] = disEvents[0]["count"]
		d["taskFailCount"] = taskFailCount[0]["count"]
		d["sqlModeUnSupportCount"] = sqlModeUnSupportCount[0]["count"]

		if len(lastEventTime) > 0 {
			d["lastEventTime"] = lastEventTime[0]["event_time"]
		} else {
			d["lastEventTime"] = ""
		}
		if len(lastAlarmTime) > 0 {
			d["lastAlarmTime"] = lastAlarmTime[0]["event_time"]
		} else {
			d["lastAlarmTime"] = ""
		}

		d["healthPct"] = healthPct
		d["alarmPieData"] = pieDataList
		//fmt.Printf("%f",healthPct)
		j, _ := json.Marshal(d)
		conn.WriteMessage(t, []byte(j))
	}

}

func DashboardInfo(c *gin.Context) {
	datasourceCount, _ := database.QueryAll("select count(*) as count from datasource limit 1")
	datasourceTypeCount, _ := database.QueryAll("select count(*) as count from (select distinct type from datasource) t limit 1")
	datasourceIdcCount, _ := database.QueryAll("select count(*) as count from (select distinct idc from datasource where idc is not null  and idc !='') t limit 1")
	datasourceEnvCount, _ := database.QueryAll("select count(*) as count from (select distinct env from datasource where env is not null  and env !='') t limit 1")
	databaseCount, _ := database.QueryAll("select count(*) as count from meta_database limit 1")
	tableCount, _ := database.QueryAll("select count(*) as count from meta_table limit 1")
	columnCount, _ := database.QueryAll("select count(*) as count from meta_column limit 1")

	datasourcePieData, _ := database.QueryAll("select type,count(*) value from datasource  group by type order by value desc limit 30")
	datasourcePieDataList := make([]map[string]interface{}, 0)
	for _, item := range datasourcePieData {
		pieData := make(map[string]interface{})
		pieData["type"] = item["type"].(string)
		pieData["value"] = utils.StrToInt(item["value"].(string))
		datasourcePieDataList = append(datasourcePieDataList, pieData)
	}

	databasePieData, _ := database.QueryAll("select datasource_type as type,count(*) value from meta_database  group by type order by value desc limit 30")
	databasePieDataList := make([]map[string]interface{}, 0)
	for _, item := range databasePieData {
		pieData := make(map[string]interface{})
		pieData["type"] = item["type"].(string)
		pieData["value"] = utils.StrToInt(item["value"].(string))
		databasePieDataList = append(databasePieDataList, pieData)
	}

	tablePieData, _ := database.QueryAll("select datasource_type as type,count(*) value from meta_table  group by type order by value desc limit 30")
	tablePieDataList := make([]map[string]interface{}, 0)
	for _, item := range tablePieData {
		pieData := make(map[string]interface{})
		pieData["type"] = item["type"].(string)
		pieData["value"] = utils.StrToInt(item["value"].(string))
		tablePieDataList = append(tablePieDataList, pieData)
	}

	columnPieData, _ := database.QueryAll("select datasource_type as type,count(*) value from meta_column  group by type order by value desc limit 30")
	columnPieDataList := make([]map[string]interface{}, 0)
	for _, item := range columnPieData {
		pieData := make(map[string]interface{})
		pieData["type"] = item["type"].(string)
		pieData["value"] = utils.StrToInt(item["value"].(string))
		columnPieDataList = append(columnPieDataList, pieData)
	}

	var data map[string]interface{}
	data = make(map[string]interface{})
	data["datasourceCount"] = datasourceCount[0]["count"]
	data["datasourceTypeCount"] = datasourceTypeCount[0]["count"]
	data["datasourceIdcCount"] = datasourceIdcCount[0]["count"]
	data["datasourceEnvCount"] = datasourceEnvCount[0]["count"]
	data["databaseCount"] = databaseCount[0]["count"]
	data["tableCount"] = tableCount[0]["count"]
	data["columnCount"] = columnCount[0]["count"]
	data["datasourcePieDataList"] = datasourcePieDataList
	data["databasePieDataList"] = databasePieDataList
	data["tablePieDataList"] = tablePieDataList
	data["columnPieDataList"] = columnPieDataList
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}
