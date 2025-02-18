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
	"dbmcloud/src/libary/clickhouse"
	"dbmcloud/src/libary/mongodb"
	"dbmcloud/src/libary/mssql"
	"dbmcloud/src/libary/mysql"
	"dbmcloud/src/libary/oracle"
	"dbmcloud/src/libary/postgres"
	"dbmcloud/src/libary/redis"
	"dbmcloud/src/model"
	"dbmcloud/src/utils"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

func init() {
	go checker()
}

func checker() {
	time.Sleep(time.Second * time.Duration(30))
	var db = database.DB
	var record model.TaskOption
	db.Select("crontab").Where("task_key=?", "check_datasource").Take(&record)
	c := cron.New()
	c.AddFunc(record.Crontab, func() {
		db.Select("enable").Where("task_key=?", "check_datasource").Take(&record)
		if record.Enable == 1 {
			db.Model(model.TaskHeartbeat{}).Where("heartbeat_key='check_datasource'").Updates(map[string]interface{}{"heartbeat_time": time.Now().Format("2006-01-02 15:04:05.999")})
			doDatasourceCheck()
			db.Model(model.TaskHeartbeat{}).Where("heartbeat_key='check_datasource'").Updates(map[string]interface{}{"heartbeat_end_time": time.Now().Format("2006-01-02 15:04:05.999")})
		}
	})
	c.Start()
}

func doDatasourceCheck() {

	var db = database.DB
	var dataList []model.Datasource
	result := db.Where("enable=1").Order("type asc").Find(&dataList)
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
		env := datasource.Env
		var origPass string
		if pass != "" {
			var err error
			origPass, err = utils.AesPassDecode(pass, setting.Setting.DbPassKey)
			if err != nil {
				fmt.Println("Encrypt Password Error.")
				return
			}
		}
		checkConnectionTask(datasourceType, env, host, port, user, origPass, dbid)
	}

}

func checkConnectionTask(datasourceType, env, host, port, user, pass, dbid string) {

	var status int32 = 1
	var statusText = "数据源服务连接正常."
	if datasourceType == "MySQL" || datasourceType == "TiDB" || datasourceType == "Doris" || datasourceType == "MariaDB" || datasourceType == "GreatSQL" || datasourceType == "OceanBase" {
		db, err := mysql.Connect(host, port, user, pass, "")
		if err != nil {
			status = 0
			statusText = fmt.Sprintf("数据源通信失败: Can't connect server on %s:%s, %s", host, port, err)
			log.Logger.Error(fmt.Sprintf("Datasource check: Can't connect server on %s:%s, %s", host, port, err))
		} else {
			defer db.Close()
		}

	} else if datasourceType == "ClickHouse" {
		db, err := clickhouse.Connect(host, port, user, pass, "")
		if err != nil {
			status = 0
			statusText = fmt.Sprintf("数据源通信失败: Can't connect server on %s:%s, %s", host, port, err)
			log.Logger.Error(fmt.Sprintf("Datasource check: Can't connect server on %s:%s, %s", host, port, err))
		} else {
			defer db.Close()
		}
	} else if datasourceType == "PostgreSQL" {
		db, err := postgres.Connect(host, port, user, pass, "postgres")
		if err != nil {
			status = 0
			statusText = fmt.Sprintf("数据源通信失败: Can't connect server on %s:%s, %s", host, port, err)
			log.Logger.Error(fmt.Sprintf("Datasource check: Can't connect server on %s:%s, %s", host, port, err))
		} else {
			defer db.Close()
		}
	} else if datasourceType == "Oracle" {
		db, err := oracle.Connect(host, port, user, pass, dbid)
		if err != nil {
			status = 0
			statusText = fmt.Sprintf("数据源通信失败: Can't connect server on %s:%s, %s", host, port, err)
			log.Logger.Error(fmt.Sprintf("Datasource check: Can't connect server on %s:%s, %s", host, port, err))
		} else {
			defer db.Close()
		}
	} else if datasourceType == "SQLServer" {
		db, err := mssql.Connect(host, port, user, pass, "")
		if err != nil {
			status = 0
			statusText = fmt.Sprintf("数据源通信失败: Can't connect server on %s:%s, %s", host, port, err)
			log.Logger.Error(fmt.Sprintf("Datasource check: Can't connect server on %s:%s, %s", host, port, err))
		} else {
			defer db.Close()
		}
	} else if datasourceType == "Redis" {
		db, err := redis.Connect(host, port, pass)
		if err != nil {
			status = 0
			statusText = fmt.Sprintf("数据源通信失败: Can't connect server on %s:%s, %s", host, port, err)
			log.Logger.Error(fmt.Sprintf("Datasource check: Can't connect server on %s:%s, %s", host, port, err))
		} else {
			defer db.Close()
		}
	} else if datasourceType == "MongoDB" {
		_, err := mongodb.Connect(host, port, user, pass, "local")
		if err != nil {
			status = 0
			statusText = fmt.Sprintf("数据源通信失败: Can't connect server on %s:%s, %s", host, port, err)
			log.Logger.Error(fmt.Sprintf("Datasource check: Can't connect server on %s:%s, %s", host, port, err))
		} else {
			//defer db.Close()
		}
	} else {
		return
	}

	var db = database.DB
	var record model.Datasource
	record.Status = status
	record.StatusText = statusText
	db.Model(&record).Select("status", "status_text").Omit("id").Where("host=?", host).Where("port=?", port).Updates(&record)

}
