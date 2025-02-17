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

package datasource

import (
	"dbmcloud/log"
	"dbmcloud/setting"
	"dbmcloud/src/database"
	"dbmcloud/src/libary/clickhouse"
	"dbmcloud/src/libary/mongodb"
	"dbmcloud/src/libary/mssql"
	"dbmcloud/src/libary/mysql"
	"dbmcloud/src/libary/oracle"
	"dbmcloud/src/libary/postgres"
	"dbmcloud/src/libary/redis"
	"dbmcloud/src/model"
	"dbmcloud/src/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func List(c *gin.Context) {
	var db = database.DB
	method := c.Request.Method
	if method == "GET" {
		if c.Query("type") != "" {
			db = db.Where("type=?", c.Query("type"))
		}
		if c.Query("enable") != "" {
			db = db.Where("enable=?", c.Query("enable"))
		}
		if c.Query("name") != "" {
			db = db.Where("name like ? ", "%"+c.Query("name")+"%")
		}
		if c.Query("host") != "" {
			db = db.Where("host like ? ", "%"+c.Query("host")+"%")
		}
		sorterMap := make(map[string]string)
		sorterData := c.Query("sorter")
		json.Unmarshal([]byte(sorterData), &sorterMap)
		for sortField, sortOrder := range sorterMap {
			if sortField != "" && sortOrder != "" {
				db = db.Order(fmt.Sprintf("%s %s", sortField, strings.Replace(sortOrder, "end", "", 1)))
			}
		}

		var dataList []model.Datasource
		result := db.Find(&dataList)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "Query Error:" + result.Error.Error()})
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
	if method == "POST" {
		var record model.Datasource
		err := c.BindJSON(&record)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "Bind record error: " + err.Error()})
			return
		}

		if len(record.Pass) > 0 {
			//log.Info("debug orig pass -->", zap.Any("pass", record.Pass))
			encryptPass, err := utils.AesPassEncode(record.Pass, setting.Setting.DbPassKey)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"success": false, "msg": "Encrypt Password Error."})
				return
			}
			record.Pass = encryptPass
		} else {
			record.Pass = ""
		}

		result := database.DB.Create(&record)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "Insert Error: " + result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
		return

	}

	if method == "PUT" {
		var record model.Datasource
		err := c.BindJSON(&record)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "Bind record error: " + err.Error()})
			return
		}

		if len(record.Pass) > 0 {
			log.Info("debug orig pass -->", zap.Any("pass", record.Pass))
			encryptPass, err := utils.AesPassEncode(record.Pass, setting.Setting.DbPassKey)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"success": false, "msg": "Encrypt Password Error."})
				return
			}
			record.Pass = encryptPass
		} else {
			record.Pass = ""
		}

		result := database.DB.Model(&record).Omit("id").Where("id = ?", record.Id).Updates(record)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "Update Error: " + result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}

	if method == "DELETE" {
		var record model.Datasource
		c.BindJSON(&record)
		result := database.DB.Model(&model.Datasource{}).Where("id = ?", record.Id).Delete(record)
		if result.Error != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "Delete Error:" + result.Error.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}

}

/*
添加数据源时检查连接
*/
func Check(c *gin.Context) {
	method := c.Request.Method
	if method == "POST" {
		var record model.Datasource
		err := c.BindJSON(&record)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"success": false, "msg": "Bind Record Error:" + err.Error()})
			return
		}
		datasourceType := record.Type
		host := record.Host
		port := record.Port
		user := record.User
		pass := record.Pass
		dbid := record.Dbid

		//更新场景，密码为空时从数据库读取密码，检查数据源是否连通
		if pass == "" {
			userPass, _ := database.QueryAll(fmt.Sprintf("select pass from datasource where host='%s' and port='%s' limit 1 ", host, port))
			passInDb := userPass[0]["pass"].(string)
			if passInDb != "" {
				var err error
				pass, err = utils.AesPassDecode(passInDb, setting.Setting.DbPassKey)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{"success": false, "msg": "Encrypt Password Error."})
					return
				}
			}
		}

		if datasourceType == "MySQL" || datasourceType == "TiDB" || datasourceType == "Doris" || datasourceType == "MariaDB" || datasourceType == "GreatSQL" || datasourceType == "OceanBase" {
			db, err := mysql.Connect(host, port, user, pass, "")
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err)})
				return
			}
			defer db.Close()
		}
		if datasourceType == "SQLServer" {
			//fmt.Println(host, port, user, pass)
			db, err := mssql.Connect(host, port, user, pass, "")
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err)})
				return
			}
			defer db.Close()
		}
		if datasourceType == "Oracle" {
			db, err := oracle.Connect(host, port, user, pass, dbid)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err)})
				return
			}
			defer db.Close()
		}
		if datasourceType == "PostgreSQL" {
			db, err := postgres.Connect(host, port, user, pass, "postgres")
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err)})
				return
			}
			defer db.Close()
		}
		if datasourceType == "ClickHouse" {
			db, err := clickhouse.Connect(host, port, user, pass, "system")
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err)})
				return
			}
			defer db.Close()
		}
		if datasourceType == "Redis" {
			db, err := redis.Connect(host, port, pass)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("Can't connect server on %s:%s, %s", host, port, err)})
				return
			}
			defer db.Close()
		}
		if datasourceType == "MongoDB" {
			_, err := mongodb.Connect(host, port, user, pass, "local")
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"success": false, "msg": fmt.Sprintf("Can't connect mongo server on %s:%s, %s", host, port, err)})
				return
			}
			//defer db.Disconnect()
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}
}
