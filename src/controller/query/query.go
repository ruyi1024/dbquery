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

package query

import (
	"context"
	"database/sql"
	"dbmcloud/setting"
	"encoding/json"
	"fmt"
	"net/http"
	_ "reflect"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"dbmcloud/src/database"
	"dbmcloud/src/libary/db"
	"dbmcloud/src/libary/mongodb"
	"dbmcloud/src/libary/redis"
	"dbmcloud/src/utils"
)

var dbCon *sql.DB
var err error

var (
	success   = "success"
	failed    = "failed"
	intercept = "intercept"
)

/*
判断字符是否在数组里面的方法
*/
func in(target string, str_array []string) bool {
	for _, element := range str_array {
		if target == element {
			return true
		}
	}
	return false
}

func DoQuery(c *gin.Context) {
	params := make(map[string]string)
	c.BindJSON(&params)
	if len(params) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "msg": "params error."})
		return
	}

	datasourceType := params["datasource_type"]
	datasource := params["datasource"]
	databaseName := params["database"]
	table := params["table"]
	sql := params["sql"]
	queryType := params["query_type"]

	username, _ := c.Get("username")

	var (
		queryTableList []string
		sqlType        string
		backupName     string
	)

	//执行SQL规则检查
	if queryType == "execute" && (datasourceType == "MySQL" || datasourceType == "TiDB" || datasourceType == "Doris" || datasourceType == "MariaDB" || datasourceType == "GreatSQL" || datasourceType == "PostgreSQL" || datasourceType == "Oracle" || datasourceType == "ClickHouse") {
		var (
			queryTable     string
			findQueryTable [][]string
		)
		//语句合法性检查
		r, err := regexp.MatchString("^(?i)select |^(?i)insert |^(?i)update |^(?i)delete |^(?i)create |^(?i)alter |^(?i)drop |^(?i)truncate |^(?i)rename ", sql)
		if err != nil || !r {
			WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, intercept, 0, sql, "不是合法的SQL命令.")
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "不是合法的SQL命令."})
			return
		}

		//提取解析SQL类型
		findSqlType := regexp.MustCompile(`^\s*(?s:(.*?)) `).FindAllStringSubmatch(sql, -1)
		sqlType = strings.ToLower(findSqlType[0][1])

		//高危命令拦截
		r1, _ := regexp.MatchString(`.*(?i)drop\s+|.*(?i)truncate\s+|.*(?i)rename\s+|.*(?i)shutdown\s+`, sql)
		if r1 {
			WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, intercept, 0, sql, "存在高风险命令.")
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "存在高风险命令."})
			return
		}

		//判断查询需要limit限制
		matchLimit, _ := regexp.MatchString(`\s+(?i)limit\s+`, sql)
		if sqlType == "select" && !matchLimit && (datasourceType == "MySQL" || datasourceType == "TiDB" || datasourceType == "PostgreSQL" || datasourceType == "Doris" || datasourceType == "ClickHouse") {
			WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, intercept, 0, sql, "数据查询请使用limit限制行数.")
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据查询请使用limit限制行数."})
			return
		}
		if (sqlType == "update" || sqlType == "delete") && !matchLimit && (datasourceType == "MySQL" || datasourceType == "TiDB") {
			WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, intercept, 0, sql, "数据变更请使用limit限制行数.")
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据变更请使用limit限制行数."})
			return
		}
		matchRownum, _ := regexp.MatchString(`\s+(?i)rownum\s*<`, sql)
		if (sqlType == "select" || sqlType == "update" || sqlType == "delete") && !matchRownum && (datasourceType == "Oracle") {
			WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, intercept, 0, sql, "数据操作请使用rownum限制行数.")
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据操作请使用rownum限制行数."})
			return
		}

		//update/delete需要where条件
		matchWhere, _ := regexp.MatchString(`\s+(?i)where\s+`, sql)
		if (sqlType == "update" || sqlType == "delete") && !matchWhere {
			WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, intercept, 0, sql, "数据修改删除必须使用where条件.")
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "数据修改删除必须使用where条件."})
			return
		}

		//提取limit行数
		var queryNumber int
		findQueryNumber := regexp.MustCompile(`(?i)limit\s+(?s:(.\d*))\s*|(?i)rownum\s*<\s*(?s:(.\d*))\s*`).FindAllStringSubmatch(sql, -1)
		if len(findQueryNumber) > 0 {
			queryNumber = utils.StrToInt(findQueryNumber[0][1])
			if queryNumber > 100000 {
				WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, intercept, 0, sql, "limit/rownum最大上限为10万.")
				c.JSON(http.StatusOK, gin.H{"success": false, "msg": "limit/rownum最大上限为10万."})
				return
			}
		}

		//解析表名
		if sqlType == "select" {
			findQueryTable = regexp.MustCompile(`(?i)select.+from\s+(?s:(.*?)) |(?i)join\s+(?s:(.*?)) `).FindAllStringSubmatch(sql, -1)
		}
		if sqlType == "insert" {
			findQueryTable = regexp.MustCompile(`(?i)insert\s+into\s+(?s:(.*?)) `).FindAllStringSubmatch(sql, -1)
		}
		if sqlType == "update" {
			findQueryTable = regexp.MustCompile(`(?i)update\s+(?s:(.*?)) `).FindAllStringSubmatch(sql, -1)
		}
		if sqlType == "delete" {
			findQueryTable = regexp.MustCompile(`(?i)delete\s+from\s+(?s:(.*?)) `).FindAllStringSubmatch(sql, -1)
		}

		var i = 0
		for _, item := range findQueryTable {
			//select join 查询的时候join表通过item[2]获取
			if i == 0 {
				queryTable = item[1]
			} else {
				queryTable = item[2]
			}
			queryTableList = append(queryTableList, queryTable)
			i++
		}

	}

	//查询数据源
	dbHostPort := strings.Split(datasource, ":")
	host := dbHostPort[0]
	port := dbHostPort[1]
	userPass, _ := database.QueryAll(fmt.Sprintf("select user,pass,dbid,dml_backup_enable,dml_backup_dir from datasource where host='%s' and port='%s' limit 1 ", host, port))
	user := userPass[0]["user"].(string)
	pass := userPass[0]["pass"].(string)
	dmlBackupEnable := userPass[0]["dml_backup_enable"].(string)
	dmlBackupDir := userPass[0]["dml_backup_dir"].(string)

	var origPass string
	if pass != "" {
		var err error
		origPass, err = utils.AesPassDecode(pass, setting.Setting.DbPassKey)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "Encrypt Password Error."})
			return
		}
	}

	if datasourceType == "MySQL" || datasourceType == "TiDB" || datasourceType == "Doris" || datasourceType == "MariaDB" || datasourceType == "GreatSQL" || datasourceType == "OceanBase" {

		if queryType == "doExplain" {
			if sql == "" || (!strings.HasPrefix(sql, "select") && !strings.HasPrefix(sql, "SELECT")) {
				c.JSON(http.StatusOK, gin.H{"success": false, "msg": "Explain的SQL语句必须以select开头"})
				return
			}
			sql = "explain " + sql
		}
		if queryType == "showColumn" {
			sql = "show columns from " + table
		}
		if queryType == "showIndex" {
			sql = "show index from " + table
		}
		if queryType == "showCreate" {
			sql = "show create table " + table
		}
		if queryType == "showTableSize" {
			sql = fmt.Sprintf("select table_schema,table_name,table_rows,data_length,data_free/1024/1024 data_free, index_length/1024/1024 index_length from information_schema.tables where table_schema='%s' and table_name='%s'", databaseName, table)
		}

		dbCon, err = db.Connect(db.WithDriver("mysql"), db.WithHost(host), db.WithPort(port), db.WithUsername(user), db.WithPassword(origPass), db.WithDatabase(databaseName))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err)})
			return
		}
		defer dbCon.Close()

	}
	if datasourceType == "Oracle" {
		sid := userPass[0]["dbid"].(string)
		dbCon, err = db.Connect(db.WithDriver("godror"), db.WithHost(host), db.WithPort(port), db.WithUsername(user), db.WithPassword(origPass), db.WithSid(sid))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("Can't connect postgres server on %s:%s, %s", host, port, err)})
			return
		}
		defer dbCon.Close()
		if queryType == "doExplain" {
			sql = "explain plan for " + sql
			db.QueryAllNew(dbCon, sql)
			sql = "select plan_table_output  \"执行计划\" from table(dbms_xplan.display) "
		}
		if queryType == "showColumn" {
			sql = fmt.Sprintf("select column_id \"字段ID\", column_name \"字段名\", data_type \"数据类型\", nullable \"允许为空\", default_length \"默认长度\", data_default \"默认值\", table_name \"表名称\" from dba_tab_columns where owner='%s' and table_name='%s' order by column_id asc", databaseName, table)
		}
		if queryType == "showIndex" {
			sql = fmt.Sprintf("select * from pg_indexes where tablename='%s'", strings.Split(table, ".")[1])
		}
		if queryType == "showCreate" {
			sql = fmt.Sprintf("%s %s", "\\d", table)
		}
		if queryType == "showTableSize" {
			sql = fmt.Sprintf("select pg_size_pretty(pg_relation_size('%s')) as size", table)
		}

	}
	if datasourceType == "PostgreSQL" {

		if queryType == "doExplain" {
			sql = "explain " + sql
		}
		if queryType == "showColumn" {
			sql = fmt.Sprintf("SELECT distinct a.attnum as num,a.attname as name,format_type(a.atttypid,a.atttypmod) as type,a.attlen as length,a.attnotnull as notnull,com.description as comment,coalesce(i.indisprimary,false) as primary_key,def.adsrc as default FROM pg_attribute a JOIN pg_class pgc ON pgc.oid = a.attrelid LEFT JOIN pg_index i ON (pgc.oid = i.indrelid AND i.indkey[0] = a.attnum) LEFT JOIN pg_description com on (pgc.oid = com.objoid AND a.attnum = com.objsubid) LEFT JOIN pg_attrdef def ON (a.attrelid = def.adrelid AND a.attnum = def.adnum) WHERE a.attnum > 0 AND pgc.oid = a.attrelid AND NOT a.attisdropped AND pgc.relname = '%s' ORDER BY a.attnum;", strings.Split(table, ".")[1])
		}
		if queryType == "showIndex" {
			sql = fmt.Sprintf("select * from pg_indexes where tablename='%s'", strings.Split(table, ".")[1])
		}
		if queryType == "showCreate" {
			sql = fmt.Sprintf("%s %s", "\\d", table)
		}
		if queryType == "showTableSize" {
			sql = fmt.Sprintf("select pg_size_pretty(pg_relation_size('%s')) as size", table)
		}

		dbCon, err = db.Connect(db.WithDriver("postgres"), db.WithHost(host), db.WithPort(port), db.WithUsername(user), db.WithPassword(origPass), db.WithDatabase(databaseName))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("Can't connect postgres server on %s:%s, %s", host, port, err)})
			return
		}
		defer dbCon.Close()
	}

	if datasourceType == "ClickHouse" {
		if queryType == "showColumn" {
			sql = fmt.Sprintf("select name as `字段名`, type  as `类型`, comment as `备注`,table as `所属表`,database as `所属库`,is_in_primary_key as `是否主键`,is_in_partition_key as `是否分区字段`,is_in_sorting_key as `是否排序字段`,is_in_sampling_key as `是否抽样字段` from system.columns where database='%s' and table='%s'", databaseName, table)
		}
		if queryType == "showCreate" {
			sql = "show create table " + table
		}
		if queryType == "showTableSize" {
			sql = fmt.Sprintf("select column as `字段名`,any(type) as `类型`, formatReadableSize(sum(column_data_uncompressed_bytes)) as `原始大小`,formatReadableSize(sum(column_data_compressed_bytes)) as `压缩大小`,sum(rows) as `行数` from system.parts_columns where database='%s' and table='%s' group by column ", databaseName, table)
		}
		dbCon, err = db.Connect(db.WithDriver("clickhouse"), db.WithHost(host), db.WithPort(port), db.WithUsername(user), db.WithPassword(origPass), db.WithDatabase(databaseName))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("Can't connect clickhouse server on %s:%s, %s", host, port, err)})
			return
		}
		defer dbCon.Close()
	}

	if datasourceType == "SQLServer" {
		if queryType == "showColumn" {
			sql = fmt.Sprintf("select name as `字段名`, type  as `类型`, comment as `备注`,table as `所属表`,database as `所属库`,is_in_primary_key as `是否主键`,is_in_partition_key as `是否分区字段`,is_in_sorting_key as `是否排序字段`,is_in_sampling_key as `是否抽样字段` from system.columns where database='%s' and table='%s'", databaseName, table)
		}
		if queryType == "showCreate" {
			sql = "show create table " + table
		}
		if queryType == "showTableSize" {
			sql = fmt.Sprintf("select column as `字段名`,any(type) as `类型`, formatReadableSize(sum(column_data_uncompressed_bytes)) as `原始大小`,formatReadableSize(sum(column_data_compressed_bytes)) as `压缩大小`,sum(rows) as `行数` from system.parts_columns where database='%s' and table='%s' group by column ", databaseName, table)
		}
		dbCon, err = db.Connect(db.WithDriver("mssql"), db.WithHost(host), db.WithPort(port), db.WithUsername(user), db.WithPassword(origPass), db.WithDatabase(databaseName))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("Can't connect clickhouse server on %s:%s, %s", host, port, err)})
			return
		}
		defer dbCon.Close()
	}

	if datasourceType != "Redis" && datasourceType != "MongoDB" {
		//执行dml语句，使用execute方法，可以获得执行行数
		if queryType == "execute" && (sqlType == "update" || sqlType == "insert" || sqlType == "delete" || sqlType == "create" || sqlType == "alter") {
			startTime := time.Now().UnixNano() / 1e6
			rowsAffected, err := db.Execute(dbCon, sql)
			endTime := time.Now().UnixNano() / 1e6
			times := endTime - startTime
			if err != nil {
				WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, failed, times, sql, fmt.Sprintf("%s", err))
				c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("%s", err)})
				return
			} else {
				if dmlBackupEnable == "1" && (sqlType == "update" || sqlType == "delete") {
					WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, success, times, sql, fmt.Sprintf("执行完成，影响%d行数据,备份文件:%s/%s", rowsAffected, dmlBackupDir, backupName))
				} else {
					WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, success, times, sql, fmt.Sprintf("执行完成，影响%d行数据", rowsAffected))
				}
				c.JSON(http.StatusOK, gin.H{"success": true, "times": times, "msg": fmt.Sprintf("执行完成，影响%d行数据", rowsAffected)})
				return
			}
		}

		//数据查询，使用queryAll方法
		startTime := time.Now().UnixNano() / 1e6
		columnList, dataList, err := db.QueryAllNew(dbCon, sql)
		endTime := time.Now().UnixNano() / 1e6
		times := endTime - startTime
		if err != nil {
			WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, failed, times, sql, fmt.Sprintf("%s", err))
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("%s", err)})
			return
		}

		//如果使用queryAll方法需要解析column，但是因为datalist的map排序会不稳定，需要使用sort排序
		/*
			var columns = make([]map[string]string, 0)
			if len(dataList) > 0 {
				//map排序不稳定，转换成数组排序
				var keys = make([]string, 0)
				for key, _ := range dataList[0] {
					keys = append(keys, key)
				}
				sort.Strings(keys)
				for _, key := range keys {
					columns = append(columns, map[string]string{"title": key, "dataIndex": key})
				}
			}
		*/

		var columns = make([]map[string]string, 0)
		if len(columnList) > 0 {
			for _, col := range columnList {
				columns = append(columns, map[string]string{"title": col, "dataIndex": col})
			}
		}

		WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, success, times, sql, fmt.Sprintf("查询完成，查询到%d行数据", len(dataList)))
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"msg":     "OK",
			"times":   times,
			"columns": columns,
			"data":    dataList,
			"total":   len(dataList),
		})
		return
	}

	if datasourceType == "Redis" {
		db, err := redis.Connect(host, port, origPass)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err)})
			return
		}
		defer db.Close()

		cmdList := strings.Split(sql, " ")
		if len(cmdList) < 2 && strings.ToLower(cmdList[0]) != "randomkey" {
			WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, intercept, 0, sql, "Redis命令不正确")
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "Redis命令不正确"})
			return
		}
		cmd := strings.ToLower(cmdList[0])
		allowCmd := strings.Split("randomkey,type,ttl,exists,get,hlen,hget,hkeys,hgetall,llen,lindex,lrange,scard,smembers,sismember,zcard,zcount,zrange", ",")
		if !in(cmd, allowCmd) {
			WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, intercept, 0, sql, "不允许的Redis命令")
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "不允许的Redis命令"})
			return
		}

		var result string
		var dataList = make([]map[string]string, 0)

		startTime := time.Now().UnixNano() / 1e6

		if cmd == "randomkey" && len(cmdList) == 1 {
			var result string
			result, err = db.RandomKey().Result()
			dataList = append(dataList, map[string]string{"执行结果": result})
		} else if cmd == "exists" && len(cmdList) == 2 {
			var result int64
			result, err = db.Exists(cmdList[1]).Result()
			dataList = append(dataList, map[string]string{"执行结果": fmt.Sprintf("%d", result)})
		} else if cmd == "type" && len(cmdList) == 2 {
			var result string
			result, err = db.Type(cmdList[1]).Result()
			dataList = append(dataList, map[string]string{"执行结果": result})
		} else if cmd == "ttl" && len(cmdList) == 2 {
			var result time.Duration
			result, err = db.TTL(cmdList[1]).Result()
			dataList = append(dataList, map[string]string{"执行结果": result.String()})
		} else if cmd == "get" && len(cmdList) == 2 {
			var result string
			result, err = db.Get(cmdList[1]).Result()
			dataList = append(dataList, map[string]string{"执行结果": result})
		} else if cmd == "hlen" && len(cmdList) == 2 {
			var result int64
			result, err = db.HLen(cmdList[1]).Result()
			dataList = append(dataList, map[string]string{"执行结果": fmt.Sprintf("%d", result)})
		} else if cmd == "hkeys" && len(cmdList) == 2 {
			var result []string
			result, err = db.HKeys(cmdList[1]).Result()
			for _, key := range result {
				dataList = append(dataList, map[string]string{"执行结果": key})
			}
		} else if cmd == "hget" && len(cmdList) == 3 {
			var result string
			result, err = db.HGet(cmdList[1], cmdList[2]).Result()
			dataList = append(dataList, map[string]string{"执行结果": result})
		} else if cmd == "hgetall" && len(cmdList) == 2 {
			var result map[string]string
			result, err = db.HGetAll(cmdList[1]).Result()
			resultJson, _ := json.Marshal(result)
			dataList = append(dataList, map[string]string{"执行结果": string(resultJson)})
		} else if cmd == "llen" && len(cmdList) == 2 {
			var result int64
			result, err = db.LLen(cmdList[1]).Result()
			dataList = append(dataList, map[string]string{"执行结果": fmt.Sprintf("%d", result)})
		} else if cmd == "lindex" && len(cmdList) == 3 {
			var result string
			result, err = db.LIndex(cmdList[1], utils.StrToInt64(cmdList[2])).Result()
			dataList = append(dataList, map[string]string{"执行结果": result})
		} else if cmd == "lrange" && len(cmdList) == 4 {
			var result []string
			result, err = db.LRange(cmdList[1], utils.StrToInt64(cmdList[2]), utils.StrToInt64(cmdList[3])).Result()
			for _, key := range result {
				dataList = append(dataList, map[string]string{"执行结果": key})
			}
		} else if cmd == "scard" && len(cmdList) == 2 {
			var result int64
			result, err = db.SCard(cmdList[1]).Result()
			dataList = append(dataList, map[string]string{"执行结果": fmt.Sprintf("%d", result)})
		} else if cmd == "smembers" && len(cmdList) == 2 {
			var result []string
			result, err = db.SMembers(cmdList[1]).Result()
			for _, key := range result {
				dataList = append(dataList, map[string]string{"执行结果": key})
			}
		} else if cmd == "sismember" && len(cmdList) == 3 {
			var result bool
			result, err = db.SIsMember(cmdList[1], cmdList[2]).Result()
			dataList = append(dataList, map[string]string{"执行结果": fmt.Sprintf("%v", result)})
		} else if cmd == "zcard" && len(cmdList) == 2 {
			var result int64
			result, err = db.ZCard(cmdList[1]).Result()
			dataList = append(dataList, map[string]string{"执行结果": fmt.Sprintf("%d", result)})
		} else if cmd == "zcount" && len(cmdList) == 4 {
			var result int64
			result, err = db.ZCount(cmdList[1], cmdList[2], cmdList[3]).Result()
			dataList = append(dataList, map[string]string{"执行结果": fmt.Sprintf("%d", result)})
		} else if cmd == "zrange" && len(cmdList) == 4 {
			var result []string
			result, err = db.ZRange(cmdList[1], utils.StrToInt64(cmdList[2]), utils.StrToInt64(cmdList[3])).Result()
			for _, key := range result {
				dataList = append(dataList, map[string]string{"执行结果": key})
			}
		} else {
			WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, intercept, 0, sql, "不正确的Redis命令")
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "不正确的Redis命令"})
			return
		}

		endTime := time.Now().UnixNano() / 1e6
		times := endTime - startTime
		if fmt.Sprintf("%s", err) == "redis: nil" {
			result = "empty"
			//dataList = append(dataList, map[string]string{"执行结果": result})
		}
		if err != nil && result != "empty" {
			WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, intercept, 0, sql, fmt.Sprintf("%s", err))
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("%s", err)})
			return
		}

		var columns = make([]map[string]string, 0)
		columns = append(columns, map[string]string{"title": "执行结果", "dataIndex": "执行结果"})
		WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, success, times, sql, "执行完成")
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"msg":     "OK",
			"times":   times,
			"columns": columns,
			"data":    dataList,
			"total":   len(dataList),
		})
		return

	}

	if datasourceType == "MongoDB" {
		sqlType = "select"
		cmdList := strings.Split(sql, ".")
		if len(cmdList) < 4 && strings.ToLower(cmdList[0]) != "select" {
			WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, intercept, 0, sql, "SQL命令不正确")
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "SQL命令不正确"})
			return
		}
		client, err := mongodb.Connect(host, port, user, origPass, "")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err)})
			return
		}
		table := strings.Replace(utils.GetBetweenStr(sql, ".from(", ")"), "'", "", -1)
		where := utils.GetBetweenStr(sql, ".where(", ")")
		limit := utils.StrToInt64(utils.GetBetweenStr(sql, ".limit(", ")"))
		whereSplit := strings.Split(where, ",")
		whereKey := strings.Replace(whereSplit[0], "'", "", -1)
		whereType := strings.Replace(whereSplit[1], "'", "", -1)
		whereVal := strings.Replace(whereSplit[2], "'", "", -1)
		fmt.Println(databaseName)
		fmt.Println(table)
		fmt.Println(where)
		fmt.Println(whereKey)
		fmt.Println(whereType)
		fmt.Println(whereVal)
		fmt.Println(limit)
		if whereType == "=" {
			whereType = "$eq"
		}
		if whereType == "!=" {
			whereType = "$ne"
		}
		if whereType == ">" {
			whereType = "$gt"
		}
		if whereType == "<" {
			whereType = "$lt"
		}

		filter := bson.D{}
		if utils.IsNumber(whereVal) {
			filter = bson.D{{"age", bson.D{{whereType, utils.StrToInt(whereVal)}}}}
		} else {
			filter = bson.D{{"age", bson.D{{whereType, whereVal}}}}
		}
		startTime := time.Now().UnixNano() / 1e6
		coll := client.Database(databaseName).Collection(table)
		// Creates a query filter to match documents in which the "cuisine"
		// is "Italian"
		//filter := bson.D{{whereKey, bson.D{{whereType, whereVal}}}}
		opts := options.Find().SetLimit(limit)
		// Retrieves documents that match the query filer
		cursor, err := coll.Find(context.TODO(), filter, opts)
		if err != nil {
			WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, failed, 0, sql, err.Error())
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": err.Error()})
			return
		}
		// end find
		var results []map[string]interface{}
		if err = cursor.All(context.TODO(), &results); err != nil {
			WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, failed, 0, sql, err.Error())
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": err.Error()})
			return
		}

		var columns = make([]map[string]string, 0)
		if len(results) > 0 {
			columnList := []string{}
			for key := range results[0] {
				columnList = append(columnList, key)
			}
			for _, col := range columnList {
				columns = append(columns, map[string]string{"title": col, "dataIndex": col})
			}
		}
		endTime := time.Now().UnixNano() / 1e6
		times := endTime - startTime
		// Prints the results of the find operation as structs
		// for _, result := range results {
		// 	cursor.Decode(&result)
		// 	output, err := json.MarshalIndent(result, "", "    ")
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	fmt.Printf("%s\n", output)
		// }
		WriteLog(username.(string), datasourceType, datasource, queryType, sqlType, databaseName, success, times, sql, "执行完成")
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"msg":     "OK",
			"times":   times,
			"columns": columns,
			"data":    results,
			"total":   len(results),
		})
		return
	}

}
