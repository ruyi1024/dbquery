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
	"dbmcloud/src/model"
	"dbmcloud/src/utils"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

func init() {
	go sensitiveTask()
}

func sensitiveTask() {
	time.Sleep(time.Second * time.Duration(30))
	var db = database.DB
	var record model.TaskOption
	db.Select("crontab").Where("task_key=?", "gather_sensitive").Take(&record)
	c := cron.New()
	c.AddFunc(record.Crontab, func() {
		db.Select("enable").Where("task_key=?", "gather_sensitive").Take(&record)
		if record.Enable == 1 {
			db.Model(model.TaskHeartbeat{}).Where("heartbeat_key='gather_sensitive'").Updates(map[string]interface{}{"heartbeat_time": time.Now().Format("2006-01-02 15:04:05.999")})
			doSensitiveTask()
			db.Model(model.TaskHeartbeat{}).Where("heartbeat_key='gather_sensitive'").Updates(map[string]interface{}{"heartbeat_end_time": time.Now().Format("2006-01-02 15:04:05.999")})
		}
	})
	c.Start()
}

func doSensitiveTask() {

	var db = database.DB
	var dataList []model.Datasource
	result := db.Where("enable=1").Where("sensitive_enable=1").Where("type in ? ", strings.Split("MySQL,MariaDB,GreatSQL,TiDB,OceanBase,Doris", ",")).Order("type asc").Find(&dataList)
	if result.Error != nil {
		log.Logger.Error(result.Error.Error())
		return

	}
	for _, datasource := range dataList {
		datasourceType := datasource.Type
		host := datasource.Host
		port := datasource.Port
		user := datasource.User
		pass := datasource.Pass
		dbid := datasource.Dbid
		var origPass string
		if pass != "" {
			var err error
			origPass, err = utils.AesPassDecode(pass, setting.Setting.DbPassKey)
			if err != nil {
				fmt.Println("Encrypt Password Error.")
				return
			}
		}
		//fmt.Println(datasourceType, host, port, user, pass, dbid)
		go startSensitiveCollectorTask(datasourceType, host, port, user, origPass, dbid)
	}

}

func startSensitiveCollectorTask(datasourceType, host, port, user, origPass, dbid string) {
	//获取敏感规则列表
	var db = database.DB
	var ruleList []model.SensitiveRule
	res := db.Where("enable=1").Find(&ruleList)
	if res.Error != nil {
		log.Logger.Error(res.Error.Error())
	}
	log.Debug(fmt.Sprintln("Get sensitive rule list:", ruleList))
	//fmt.Println("Get sensitive rule list:", ruleList)
	if datasourceType == "MySQL" || datasourceType == "TiDB" || datasourceType == "Doris" || datasourceType == "MariaDB" || datasourceType == "GreatSQL" || datasourceType == "OceanBase" {
		queryDatabaseSql := "select schema_name schema_name from information_schema.schemata where schema_name not in ('information_schema','performance_schema','sys','mysql','metrics_schema','__internal_schema','sys_audit','lbacsys','oceanbase','ocs','oraauditor')"
		dbCon, err = database.Connect(database.WithDriver("mysql"), database.WithHost(host), database.WithPort(port), database.WithUsername(user), database.WithPassword(origPass), database.WithDatabase("information_schema"))
		if err != nil {
			log.Logger.Error(fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err))
			return
		}
		defer dbCon.Close()
		schemaList, err := database.QueryRemote(dbCon, queryDatabaseSql)
		if err != nil {
			log.Logger.Error(fmt.Sprintf("Can't query server on %s:%s, %s", host, port, err))
			return
		}
		for _, item := range schemaList {
			schemaName := item["schema_name"]
			queryTableSql := fmt.Sprintf("select table_name table_name,table_comment table_comment from information_schema.tables where table_schema='%s' ", schemaName)
			dbCon, err = database.Connect(database.WithDriver("mysql"), database.WithHost(host), database.WithPort(port), database.WithUsername(user), database.WithPassword(origPass), database.WithDatabase("information_schema"))
			if err != nil {
				log.Logger.Error(fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err))
				return
			}
			defer dbCon.Close()
			tableList, err := database.QueryRemote(dbCon, queryTableSql)
			if err != nil {
				log.Logger.Error(fmt.Sprintf("Can't query server on %s:%s, %s", host, port, err))
				return
			}
			//fmt.Println(tableList)
			for _, item := range tableList {
				tableName := item["table_name"]
				tablecomment := item["table_comment"]
				dbCon, err = database.Connect(database.WithDriver("mysql"), database.WithHost(host), database.WithPort(port), database.WithUsername(user), database.WithPassword(origPass), database.WithDatabase("information_schema"))
				if err != nil {
					log.Logger.Error(fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err))
					return
				}
				defer dbCon.Close()
				data, err := database.QueryRemote(dbCon, fmt.Sprintf("select * from (select * from %s.%s limit 1000) t order by rand() limit 100", schemaName, tableName))
				if err != nil {
					log.Logger.Error(fmt.Sprintf("Can't query server on %s:%s, %s", host, port, err))
					return
				}
				if len(data) > 0 {
					//获取字段，取第一行数据使用range拿到字段
					for column := range data[0] {
						simpleData := make([]string, 0, len(data))
						for _, item := range data {
							if item[column] != nil {
								simpleData = append(simpleData, item[column].(string))
							} else {
								simpleData = append(simpleData, "")
							}
						}
						startDataRuleScan(datasourceType, host, port, schemaName.(string), tableName.(string), tablecomment.(string), column, simpleData, ruleList)
					}
				}
			}
		}
	}

}

func startDataRuleScan(datasourceType, host, port, schemaName, tableName, tablecomment, column string, simpleData []string, ruleList []model.SensitiveRule) {
	simpleCount := len(simpleData)
	//fmt.Println(simpleCount)
	var db = database.DB
	for _, rule := range ruleList {
		ruleType := rule.RuleType
		ruleKey := rule.RuleKey
		ruleName := rule.RuleName
		ruleExpress := rule.RuleExpress
		rulePct := rule.RulePct
		level := rule.Level
		status := rule.Status
		sensitiveColumn := false
		sensitiveCount := 0
		sensitivePct := 0
		if ruleType == "column" {
			pattern := fmt.Sprintf(`%s`, ruleExpress)
			matched, _ := regexp.MatchString(pattern, column)
			if matched {
				sensitiveColumn = true
			}
		}
		if ruleType == "data" {
			for _, value := range simpleData {
				pattern := fmt.Sprintf(`%s`, ruleExpress)
				matched, _ := regexp.MatchString(pattern, value)
				if matched {
					sensitiveCount += 1
					//fmt.Println(simpleData)
					//fmt.Println(sensitiveCount)
					//fmt.Println(sensitiveCount * 100 / simpleCount)
				}
			}
		}
		sensitivePct = (sensitiveCount * 100 / simpleCount)

		if sensitiveColumn || (sensitiveCount > 0 && sensitivePct > rulePct) {
			var dataList []model.SensitiveMeta
			db.Where("datasource_type=?", datasourceType).Where("host=?", host).Where("port=?", port).Where("database_name=?", schemaName).Where("table_name=?", tableName).Where("column_name=?", column).Where("rule_type=?", ruleType).Where("rule_key=?", ruleKey).Find(&dataList)
			if len(dataList) == 0 {
				var record model.SensitiveMeta
				record.DatasourceType = datasourceType
				record.Host = host
				record.Port = port
				record.DatabaseName = schemaName
				record.TableNameX = tableName
				record.TableComment = tablecomment
				record.ColumnName = column
				record.ColumnComment = ""
				record.RuleType = ruleType
				record.RuleKey = ruleKey
				record.RuleName = ruleName
				record.SensitiveCount = sensitiveCount
				record.SimpleCount = simpleCount
				record.Level = level
				record.Status = status
				result := database.DB.Create(&record)
				if result.Error != nil {
					log.Logger.Error(fmt.Sprintf("Can't add sensitive meta on %s:%s,%s", host, port, result.Error.Error()))
				}
			} else {
				var record model.SensitiveMeta
				record.TableComment = tablecomment
				record.ColumnComment = ""
				record.RuleName = ruleName
				record.SensitiveCount = sensitiveCount
				record.SimpleCount = simpleCount
				record.Level = level
				db.Model(&record).Select("table_comment", "column_comment", "rule_name", "sensitive_count", "simple_count", "level").Omit("id").Where("host=?", host).Where("port=?", port).Where("database_name=?", schemaName).Where("table_name=?", tableName).Where("column_name=?", column).Where("rule_type=?", ruleType).Where("rule_key=?", ruleKey).Updates(&record)
			}

		}

	}
}
