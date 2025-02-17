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
	"database/sql"
	"dbmcloud/log"
	"dbmcloud/setting"
	"dbmcloud/src/database"
	"dbmcloud/src/model"
	"dbmcloud/src/utils"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

var dbCon *sql.DB
var err error

func init() {
	go dbMetaCrontabTask()
}

func dbMetaCrontabTask() {
	time.Sleep(time.Second * time.Duration(30))
	var db = database.DB
	var record model.TaskOption
	db.Select("crontab").Where("task_key=?", "gather_dbmeta").Take(&record)
	c := cron.New()
	c.AddFunc(record.Crontab, func() {
		db.Select("enable").Where("task_key=?", "gather_dbmeta").Take(&record)
		if record.Enable == 1 {
			db.Model(model.TaskHeartbeat{}).Where("heartbeat_key='gather_dbmeta'").Updates(map[string]interface{}{"heartbeat_time": time.Now().Format("2006-01-02 15:04:05.999")})
			doDbMetaTask()
			db.Model(model.TaskHeartbeat{}).Where("heartbeat_key='gather_dbmeta'").Updates(map[string]interface{}{"heartbeat_end_time": time.Now().Format("2006-01-02 15:04:05.999")})
		}
	})
	c.Start()
}

func formatInterface(inter interface{}) string {
	if inter != nil {
		return inter.(string)
	} else {
		return ""
	}
}

func doDbMetaTask() {

	var db = database.DB
	var dataList []model.Datasource
	result := db.Where("enable=1").Where("dbmeta_enable=1").Order("type asc").Find(&dataList)
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
		doDbMetaCollectorTask(datasourceType, host, port, user, origPass, dbid)
	}
	// clear expire meta data
	database.DB.Model(model.MetaDatabase{}).Where("gmt_updated <= ?", time.Now().Add(-time.Minute*10).Format("2006-01-02 15:04:05")).Updates(map[string]interface{}{"is_deleted": 1})
	database.DB.Model(model.MetaTable{}).Where("gmt_updated <= ?", time.Now().Add(-time.Minute*10).Format("2006-01-02 15:04:05")).Updates(map[string]interface{}{"is_deleted": 1})
	database.DB.Model(model.MetaColumn{}).Where("gmt_updated <= ?", time.Now().Add(-time.Minute*10).Format("2006-01-02 15:04:05")).Updates(map[string]interface{}{"is_deleted": 1})

}

func getDbCon(datasourceType, host, port, user, origPass, dbid string) *sql.DB {
	if datasourceType == "MySQL" || datasourceType == "TiDB" || datasourceType == "Doris" || datasourceType == "MariaDB" || datasourceType == "GreatSQL" || datasourceType == "OceanBase" {
		dbCon, err = database.Connect(database.WithDriver("mysql"), database.WithHost(host), database.WithPort(port), database.WithUsername(user), database.WithPassword(origPass), database.WithDatabase("information_schema"))
		if err != nil {
			log.Logger.Error(fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err))
			return nil
		}
	} else if datasourceType == "ClickHouse" {
		dbCon, err = database.Connect(database.WithDriver("clickhouse"), database.WithHost(host), database.WithPort(port), database.WithUsername(user), database.WithPassword(origPass), database.WithDatabase("system"))
		if err != nil {
			log.Logger.Error(fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err))
			return nil
		}
	}
	return dbCon
}

func doDbMetaCollectorTask(datasourceType, host, port, user, origPass, dbid string) {

	var db = database.DB
	var (
		queryDatabaseSql string
		queryTableSql    string
		queryColumnSql   string
	)
	if datasourceType == "MySQL" || datasourceType == "TiDB" || datasourceType == "Doris" || datasourceType == "MariaDB" || datasourceType == "GreatSQL" || datasourceType == "OceanBase" {
		queryDatabaseSql = "select lower(schema_name) as database_name,lower(schema_name) as schema_name,default_character_set_name as characters from information_schema.schemata where lower(schema_name) not in ('information_schema','performance_schema','sys','mysql','metrics_schema','__internal_schema','sys_audit','lbacsys','oceanbase','ocs','oraauditor') order by database_name asc"
		queryTableSql = "select table_type as table_type,lower(table_schema) as database_name,lower(table_name) as table_name,table_comment as table_comment,table_collation as characters from information_schema.tables where lower(table_schema) not in ('information_schema','performance_schema','sys','mysql','metrics_schema','__internal_schema','sys_audit','lbacsys','oceanbase','ocs','oraauditor') order by database_name asc,table_name asc"
		queryColumnSql = "select lower(table_schema)  as database_name,lower(table_name) as table_name,lower(column_name) as column_name,  lower(column_comment) as column_comment, lower(data_type) as data_type,lower(is_nullable) as is_nullable,lower(column_default) as default_value,lower(ordinal_position) as ordinal_position,lower(collation_name) as characters from information_schema.COLUMNS where lower(table_schema) not in ('information_schema','performance_schema','sys','mysql','metrics_schema','__internal_schema','sys_audit','lbacsys','oceanbase','ocs','oraauditor')  order by table_name asc,ordinal_position asc"

	} else if datasourceType == "ClickHouse" {
		queryDatabaseSql = "select lower(name) as database_name,lower(name) as schema_name,'' as characters from databases where lower(name) not in ('information_schema','INFORMATION_SCHEMA','system') order by name asc"
		//queryTableSql = "select engine as table_type,lower(`database`) as database_name,lower(name) as table_name,comment as table_comment,'' as characters from tables where database_name not in ('information_schema','INFORMATION_SCHEMA','system') order by database_name asc,table_name asc limit 100"
		//queryColumnSql = "select lower(`database`) as database_name,lower(`table`) as table_name, lower(name) as column_name,comment as column_comment,type as data_type,'' as is_nullable, '' as default_value, toString(position) as ordinal_position,'' as characters from columns where database_name not in ('information_schema','INFORMATION_SCHEMA','system') order by database_name asc,table_name asc,ordinal_position asc"
		queryTableSql = "select engine as table_type,lower(`database`) as database_name,name as table_name,comment as table_comment,'' as characters from tables where database_name not in ('information_schema','INFORMATION_SCHEMA','system')  order by database_name asc,table_name asc limit 100"
		queryColumnSql = "select lower(`database`) as database_name,lower(`table`) as table_name, lower(name) as column_name,comment as column_comment,type as data_type,'' as is_nullable, '' as default_value, toString(position) as ordinal_position,'' as characters from columns where database_name not in ('information_schema','INFORMATION_SCHEMA','system') order by database_name asc,table_name asc,ordinal_position asc"

		// } else if datasourceType == "PostgreSQL" {
		// 	queryDatabaseSql = "select pg_database.datname as database_name,pg_database.datname as schema_name,pg_encoding_to_char(encoding) as characters from pg_database where datname not in ('postgres','template0','template1') order by database_name asc"
		// 	queryTableSql = ""
		// 	dbCon, err = database.Connect(database.WithDriver("postgres"), database.WithHost(host), database.WithPort(port), database.WithUsername(user), database.WithPassword(origPass), database.WithDatabase("postgres"))
		// 	if err != nil {
		// 		log.Logger.Error(fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err))
		// 		return
		// 	}
		// 	//defer dbCon.Close()
		// } else if datasourceType == "Oracle" {
		// 	queryDatabaseSql = "select username as database_name,username as schema_name,'' as characters from dba_users where username not in ('SYSTEM','SYS') order by username asc;"
		// 	queryTableSql = ""
		// 	dbCon, err = database.Connect(database.WithDriver("godror"), database.WithHost(host), database.WithPort(port), database.WithUsername(user), database.WithPassword(origPass), database.WithDatabase(dbid))
		// 	if err != nil {
		// 		log.Logger.Error(fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err))
		// 		return
		// 	}
		// 	//defer dbCon.Close()
		// } else if datasourceType == "SQLServer" {
		// 	queryDatabaseSql = "SELECT name as database_name,name as schema_name,collation_name as characters FROM sys.databases where name not in ('master','tempdb','msdb','model') order by name asc"
		// 	queryTableSql = "SELECT o.type_desc AS table_type, DB_NAME() AS database_name,  o.name AS table_name, CAST(ep.value AS NVARCHAR(MAX)) AS table_comment,'' as characters FROM  sys.objects o LEFT JOIN  sys.extended_properties ep ON o.object_id = ep.major_id AND ep.name = 'MS_Description'   WHERE  o.type IN ('U')  AND o.is_ms_shipped = 0"
		// 	queryColumnSql = "SELECT DB_NAME() AS database_name,  t.name AS table_Name, c.name AS column_name, '' as column_comment,ty.name AS data_type, '' as is_nullable, OBJECT_DEFINITION(c.default_object_id) AS default_value,'' as ordinal_position,'' as characters FROM sys.tables t INNER JOIN sys.columns c ON t.object_id = c.object_id LEFT JOIN sys.types ty ON c.system_type_id = ty.system_type_id AND c.user_type_id = ty.user_type_id WHERE t.is_ms_shipped = 0"
		// 	dbCon, err = database.Connect(database.WithDriver("mssql"), database.WithHost(host), database.WithPort(port), database.WithUsername(user), database.WithPassword(origPass), database.WithDatabase("master"))
		// 	fmt.Println(queryTableSql)
		// 	fmt.Println(dbCon)
		// 	if err != nil {
		// 		log.Logger.Error(fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err))
		// 		return
		// 	}
		// 	//defer dbCon.Close()
	} else {
		return
	}

	//采集数据库列表
	dbCon = getDbCon(datasourceType, host, port, user, origPass, dbid)
	if dbCon == nil {
		return
	}
	defer dbCon.Close()
	databaseList, err := database.QueryRemote(dbCon, queryDatabaseSql)
	if err != nil {
		fmt.Println(err)
		log.Logger.Error(fmt.Sprintf("Can't query database meta on %s:%s, %s", host, port, err))
	}
	for _, item := range databaseList {
		var dataList []model.MetaDatabase
		if item["database_name"] == nil || item["schema_name"] == nil {
			return
		}
		db.Where("host=?", host).Where("port=?", port).Where("database_name=?", item["database_name"].(string)).Where("schema_name=?", item["schema_name"].(string)).Find(&dataList)
		if (len(dataList)) == 0 {
			var record model.MetaDatabase
			record.DatasourceType = datasourceType
			record.Host = host
			record.Port = port
			record.DatabaseName = item["database_name"].(string)
			record.SchemaName = item["schema_name"].(string)
			if item["characters"] != nil {
				record.Characters = item["characters"].(string)
			} else {
				record.Characters = ""
			}
			result := database.DB.Create(&record)
			if result.Error != nil {
				fmt.Println(result.Error.Error())
				log.Logger.Error(fmt.Sprintf("Can't collector database on %s:%s, %s", host, port, result.Error.Error()))
			}
		} else {
			var record model.MetaDatabase
			record.Characters = formatInterface(item["characters"])
			record.IsDeleted = 0
			result := db.Model(&record).Select("characters", "is_deleted").Omit("id").Where("host=?", host).Where("port=?", port).Where("database_name=?", item["database_name"].(string)).Where("schema_name=?", item["schema_name"].(string)).Updates(&record)
			if result.Error != nil {
				log.Logger.Error(fmt.Sprintf("Can't collector database on %s:%s, %s", host, port, result.Error.Error()))
			}
		}
		time.Sleep(1 * time.Millisecond)
	}

	//采集数据表列表
	dbCon = getDbCon(datasourceType, host, port, user, origPass, dbid)
	if dbCon == nil {
		return
	}
	defer dbCon.Close()
	tableList, err := database.QueryRemote(dbCon, queryTableSql)
	if err != nil {
		fmt.Println(err)
		log.Logger.Error(fmt.Sprintf("Can't query table meta on %s:%s, %s", host, port, err))
	}

	for _, item := range tableList {
		var dataList []model.MetaTable
		db.Where("host=?", host).Where("port=?", port).Where("database_name=?", item["database_name"].(string)).Where("table_name=?", item["table_name"].(string)).Find(&dataList)
		if (len(dataList)) == 0 {
			//fmt.Println(dataList)
			var record model.MetaTable
			record.DatasourceType = datasourceType
			record.Host = host
			record.Port = port
			record.DatabaseName = item["database_name"].(string)
			record.TableType = formatInterface(item["table_type"])
			record.TableNameX = item["table_name"].(string)
			record.TableComment = formatInterface(item["table_comment"])
			record.Characters = formatInterface(item["characters"])
			result := database.DB.Create(&record)
			if result.Error != nil {
				log.Logger.Error(fmt.Sprintf("Can't collector table on %s:%s, %s", host, port, result.Error.Error()))
			}
		} else {
			var record model.MetaTable
			record.TableType = formatInterface(item["table_type"])
			record.TableComment = formatInterface(item["table_comment"])
			record.Characters = formatInterface(item["characters"])
			result := db.Model(&record).Select("table_comment", "table_type", "characters").Omit("id").Where("host=?", host).Where("port=?", port).Where("database_name=?", item["database_name"].(string)).Where("table_name=?", item["table_name"].(string)).Updates(&record)
			if result.Error != nil {
				log.Logger.Error(fmt.Sprintf("Can't collector table on %s:%s, %s", host, port, result.Error.Error()))
			}
		}
		time.Sleep(1 * time.Millisecond)
	}

	//采集字段列表
	dbCon = getDbCon(datasourceType, host, port, user, origPass, dbid)
	if dbCon == nil {
		return
	}
	defer dbCon.Close()
	columnList, err := database.QueryRemote(dbCon, queryColumnSql)
	if err != nil {
		fmt.Println(err)
		log.Logger.Error(fmt.Sprintf("Can't query column meta on %s:%s, %s", host, port, err))
	}
	for _, item := range columnList {
		var dataList []model.MetaColumn
		db.Where("host=?", host).Where("port=?", port).Where("database_name=?", item["database_name"].(string)).Where("table_name=?", item["table_name"].(string)).Where("column_name=?", item["column_name"].(string)).Find(&dataList)
		if (len(dataList)) == 0 {
			var record model.MetaColumn
			record.DatasourceType = datasourceType
			record.Host = host
			record.Port = port
			record.DatabaseName = item["database_name"].(string)
			record.TableNameX = item["table_name"].(string)
			record.ColumnName = item["column_name"].(string)
			record.ColumnComment = formatInterface(item["column_comment"])
			record.DataType = formatInterface(item["data_type"])
			record.IsNullable = formatInterface(item["is_nullable"])
			record.DefaultValue = formatInterface(item["default_value"])
			record.Ordinal_Position = utils.StrToInt(item["ordinal_position"].(string))
			record.Characters = formatInterface(item["characters"])
			result := database.DB.Create(&record)
			if result.Error != nil {
				log.Logger.Error(fmt.Sprintf("Can't collector column on %s:%s, %s", host, port, result.Error.Error()))
			}
		} else {
			var record model.MetaColumn
			record.ColumnComment = formatInterface(item["column_comment"])
			record.DataType = formatInterface(item["data_type"])
			record.IsNullable = formatInterface(item["is_nullable"])
			record.DefaultValue = formatInterface(item["default_value"])
			record.Ordinal_Position = utils.StrToInt(item["ordinal_position"].(string))
			record.Characters = formatInterface(item["characters"])
			result := db.Model(&record).Select("column_comment", "data_type", "is_nullable", "default_value", "ordinal_position", "characters").Omit("id").Where("host=?", host).Where("port=?", port).Where("database_name=?", item["database_name"].(string)).Where("table_name=?", item["table_name"].(string)).Where("column_name=?", item["column_name"].(string)).Updates(&record)
			if result.Error != nil {
				log.Logger.Error(fmt.Sprintf("Can't collector column on %s:%s, %s", host, port, result.Error.Error()))
			}
		}
		time.Sleep(1 * time.Millisecond)
	}

}
