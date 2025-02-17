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
	"dbmcloud/setting"
	"dbmcloud/src/database"
	"dbmcloud/src/libary/clickhouse"
	"dbmcloud/src/libary/mongodb"
	"dbmcloud/src/libary/mssql"
	"dbmcloud/src/libary/mysql"
	"dbmcloud/src/libary/oracle"
	"dbmcloud/src/libary/postgres"
	"dbmcloud/src/utils"
	"fmt"
	"net/http"
	_ "reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

func TableList(c *gin.Context) {
	datasourceType := c.Query("type")
	datasource := c.Query("datasource")
	databaseName := c.Query("database")
	HostPort := strings.Split(datasource, ":")
	host := HostPort[0]
	port := HostPort[1]
	userPass, _ := database.QueryAll(fmt.Sprintf("select user,pass,dbid from datasource where host='%s' and port='%s' limit 1 ", host, port))
	user := userPass[0]["user"].(string)
	pass := userPass[0]["pass"].(string)

	var origPass string
	if user != "" && pass != "" {
		var err error
		origPass, err = utils.AesPassDecode(pass, setting.Setting.DbPassKey)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "Encrypt Password Error."})
			return
		}
	}

	if datasourceType == "MySQL" || datasourceType == "TiDB" || datasourceType == "Doris" || datasourceType == "MariaDB" || datasourceType == "GreatSQL" || datasourceType == "OceanBase" {
		db, err := mysql.Connect(host, port, user, origPass, "information_schema")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err)})
			return
		}
		defer db.Close()
		dataList, err := mysql.QueryAll(db, fmt.Sprintf("select table_name as table_name from tables where table_schema='%s' order by table_name asc", databaseName))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("%s", err)})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"msg":     "OK",
			"data":    dataList,
			"total":   len(dataList),
		})
		return
	}

	if datasourceType == "Oracle" {
		sid := userPass[0]["dbid"].(string)
		db, err := oracle.Connect(host, port, user, origPass, sid)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err)})
			return
		}
		defer db.Close()
		dataList, err := oracle.QueryAll(db, fmt.Sprintf("select table_name  from dba_tables where owner='%s' order by table_name asc", databaseName))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("%s", err)})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"msg":     "OK",
			"data":    dataList,
			"total":   len(dataList),
		})
		return
	}

	if datasourceType == "PostgreSQL" {
		db, err := postgres.Connect(host, port, user, origPass, databaseName)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err)})
			return
		}
		defer db.Close()
		dataList, err := postgres.QueryAll(db, "select concat(schemaname,'.',tablename) as table_name from pg_tables where schemaname not in ('pg_catalog','information_schema') order by table_name asc")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("%s", err)})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"msg":     "OK",
			"data":    dataList,
			"total":   len(dataList),
		})
		return
	}

	if datasourceType == "ClickHouse" {
		db, err := clickhouse.Connect(host, port, user, origPass, "")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err)})
			return
		}
		defer db.Close()
		dataList, err := clickhouse.QueryAll(db, fmt.Sprintf("select name as table_name from system.tables where database='%s' order by name asc", databaseName))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("%s", err)})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"msg":     "OK",
			"data":    dataList,
			"total":   len(dataList),
		})
		return
	}

	if datasourceType == "SQLServer" {
		db, err := mssql.Connect(host, port, user, origPass, databaseName)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err)})
			return
		}
		defer db.Close()
		dataList, err := mssql.QueryAll(db, "SELECT name as table_name FROM sys.tables")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("%s", err)})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"msg":     "OK",
			"data":    dataList,
			"total":   len(dataList),
		})
		return
	}
	if datasourceType == "MongoDB" {
		client, err := mongodb.Connect(host, port, user, origPass, "")
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err)})
			return
		}
		result, err := mongodb.ListCollection(client, databaseName)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("%s", err)})
			return
		}
		var dataList []map[string]string
		for _, collection := range result {
			dataMap := make(map[string]string)
			dataMap["table_name"] = collection
			dataList = append(dataList, dataMap)
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"msg":     "OK",
			"data":    dataList,
			"total":   len(dataList),
		})
		return
	}

}
