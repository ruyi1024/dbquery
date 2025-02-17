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

package database

import (
	"database/sql"
	"dbmcloud/log"
	"dbmcloud/setting"
	"dbmcloud/src/aes"
	"dbmcloud/src/model"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	//_ "github.com/ClickHouse/clickhouse-go"
	//_ "github.com/go-sql-driver/mysql"
	//_ "github.com/lib/pq"
	"github.com/go-redis/redis"
)

var DB *gorm.DB
var CK *gorm.DB
var SQL *sql.DB
var RDS *redis.Client

func InitDb() *gorm.DB {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
			os.Exit(0)

		}
	}()

	ds := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", setting.Setting.User, setting.Setting.Password, setting.Setting.Host, setting.Setting.Port, setting.Setting.Database)
	log.Info("debug mysql: " + fmt.Sprintf("%s", ds))
	sqlDB, err := sql.Open("mysql", ds)
	if err != nil {
		log.Error("open database error", zap.Error(err))
		panic(fmt.Sprintln("open database error.", zap.Error(err)))
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{
		Logger: log.NewGormLogger(zapcore.InfoLevel, zapcore.InfoLevel, time.Millisecond*200),
	})
	if err != nil {
		log.Error("grom open database error", zap.Error(err))
		panic(fmt.Sprintln("grom open database error.", zap.Error(err)))
	}

	if !db.Migrator().HasTable(&model.Users{}) {
		if err = db.AutoMigrate(&model.Users{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
		result := db.Create(&model.Users{Id: 1, Username: "admin", ChineseName: "管理员", Password: "a8a0d32f1abefd3fa996321d5e72c6d6", Admin: true})
		if result.Error != nil {
			panic(result.Error)
		}
	}

	if !db.Migrator().HasTable(&model.Token{}) {
		if err = db.AutoMigrate(&model.Token{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
	}

	if !db.Migrator().HasTable(&model.DatasourceType{}) {
		if err = db.AutoMigrate(&model.DatasourceType{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
		db.Create(&model.DatasourceType{Id: 1, Name: "MySQL", Sort: 1, Enable: 1})
		db.Create(&model.DatasourceType{Id: 2, Name: "MariaDB", Sort: 2, Enable: 1})
		db.Create(&model.DatasourceType{Id: 3, Name: "GreatSQL", Sort: 3, Enable: 1})
		db.Create(&model.DatasourceType{Id: 4, Name: "TiDB", Sort: 4, Enable: 1})
		db.Create(&model.DatasourceType{Id: 5, Name: "Doris", Sort: 5, Enable: 1})
		db.Create(&model.DatasourceType{Id: 6, Name: "OceanBase", Sort: 6, Enable: 1})
		db.Create(&model.DatasourceType{Id: 7, Name: "ClickHouse", Sort: 7, Enable: 1})
		db.Create(&model.DatasourceType{Id: 8, Name: "Oracle", Sort: 8, Enable: 1})
		db.Create(&model.DatasourceType{Id: 9, Name: "PostgreSQL", Sort: 9, Enable: 1})
		db.Create(&model.DatasourceType{Id: 10, Name: "SQLServer", Sort: 10, Enable: 1})
		db.Create(&model.DatasourceType{Id: 11, Name: "MongoDB", Sort: 11, Enable: 1})
		db.Create(&model.DatasourceType{Id: 12, Name: "Redis", Sort: 12, Enable: 1})

	}

	if !db.Migrator().HasTable(&model.Idc{}) {
		if err = db.AutoMigrate(&model.Idc{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
		db.Create(&model.Idc{Id: 1, IdcKey: "default", IdcName: "默认机房", Description: "默认未分类机房"})
	}

	if !db.Migrator().HasTable(&model.Env{}) {
		if err = db.AutoMigrate(&model.Env{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
		db.Create(&model.Env{Id: 1, EnvKey: "dev", EnvName: "开发环境", Description: "业务功能开发环境"})
		db.Create(&model.Env{Id: 2, EnvKey: "test", EnvName: "测试环境", Description: "业务功能测试环境"})
		db.Create(&model.Env{Id: 3, EnvKey: "pre", EnvName: "预发环境", Description: "准生产验证环境"})
		db.Create(&model.Env{Id: 4, EnvKey: "prod", EnvName: "生产环境", Description: "线上业务运行环境"})
	}

	if !db.Migrator().HasTable(&model.Datasource{}) {
		if err = db.AutoMigrate(&model.Datasource{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
		aesPassword, _ := aes.AesPassEncode(setting.Setting.Password, setting.Setting.DbPassKey)
		aesCkPassword, _ := aes.AesPassEncode(setting.Setting.ClickhousePassword, setting.Setting.DbPassKey)
		aesRdsPassword, _ := aes.AesPassEncode(setting.Setting.RedisPassword, setting.Setting.DbPassKey)
		db.Create(&model.Datasource{Id: 1, Name: "LEPUS-MySQL", GroupName: "Lepus", Idc: "default", Env: "prod", Type: "MySQL", Host: setting.Setting.Host, Port: setting.Setting.Port, User: setting.Setting.User, Pass: aesPassword, Enable: 1, DbmetaEnable: 1, ExecuteEnable: 1, MonitorEnable: 1, AlarmEnable: 1})
		db.Create(&model.Datasource{Id: 2, Name: "LEPUS-ClickHouse", GroupName: "Lepus", Idc: "default", Env: "prod", Type: "ClickHouse", Host: setting.Setting.ClickhouseHost, Port: setting.Setting.ClickhousePort, User: setting.Setting.ClickhouseUser, Pass: aesCkPassword, Enable: 1, DbmetaEnable: 1, ExecuteEnable: 1, MonitorEnable: 1, AlarmEnable: 1})
		db.Create(&model.Datasource{Id: 3, Name: "LEPUS-Redis", GroupName: "Lepus", Idc: "default", Env: "prod", Type: "Redis", Host: setting.Setting.RedisHost, Port: setting.Setting.RedisPort, User: "", Pass: aesRdsPassword, Enable: 1, DbmetaEnable: 0, ExecuteEnable: 1, MonitorEnable: 1, AlarmEnable: 1})

	}

	if !db.Migrator().HasTable(&model.MetaDatabase{}) {
		if err = db.AutoMigrate(&model.MetaDatabase{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
	}

	if !db.Migrator().HasTable(&model.MetaTable{}) {
		if err = db.AutoMigrate(&model.MetaTable{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
	}

	if !db.Migrator().HasTable(&model.MetaColumn{}) {
		if err = db.AutoMigrate(&model.MetaColumn{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
	}

	if !db.Migrator().HasTable(&model.TaskOption{}) {
		if err = db.AutoMigrate(&model.TaskOption{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
		db.Create(&model.TaskOption{TaskKey: "recycle_token", TaskName: "回收用户令牌", TaskDescription: "回收用户过期的ToKen", Crontab: "* * * * *"})
		db.Create(&model.TaskOption{TaskKey: "revoke_privileage", TaskName: "回收用户权限", TaskDescription: "检查用户查询数据库权限是否过期，并回收权限", Crontab: "1 * * * *"})
		db.Create(&model.TaskOption{TaskKey: "check_datasource", TaskName: "监测数据源状态", TaskDescription: "监测数据源连接状态是否正常", Crontab: "@every 30s"})
		db.Create(&model.TaskOption{TaskKey: "gather_dbmeta", TaskName: "采集元数据信息", TaskDescription: "采集数据库、数据表、数据列等元数据信息", Crontab: "*/3 * * * *"})
		db.Create(&model.TaskOption{TaskKey: "gather_sensitive", TaskName: "敏感数据探测分析", TaskDescription: "分析数据库数据，监测敏感信息", Crontab: "*/5 * * * *"})
		db.Create(&model.TaskOption{TaskKey: "collector_mysql_event", TaskName: "采集MySQL服务状态", TaskDescription: "采集MySQL等数据库的服务运行状态数据", Crontab: "*/1 * * * *"})
		db.Create(&model.TaskOption{TaskKey: "collector_redis_event", TaskName: "采集Redis服务状态", TaskDescription: "采集Redis数据库的服务运行状态数据", Crontab: "*/1 * * * *"})
		db.Create(&model.TaskOption{TaskKey: "collector_oracle_event", TaskName: "采集Oracle服务状态", TaskDescription: "采集Oracle数据库的服务运行状态数据", Crontab: "*/1 * * * *"})
		db.Create(&model.TaskOption{TaskKey: "collector_postgresql_event", TaskName: "采集PostgreSQL服务状态", TaskDescription: "采集PostgreSQL数据库的服务运行状态数据", Crontab: "*/1 * * * *"})
		db.Create(&model.TaskOption{TaskKey: "collector_mongodb_event", TaskName: "采集MongoDB服务状态", TaskDescription: "采集MongoDB数据库的服务运行状态数据", Crontab: "*/1 * * * *"})
		db.Create(&model.TaskOption{TaskKey: "collector_sqlserver_event", TaskName: "采集SQLServer服务状态", TaskDescription: "采集SQLServer数据库的服务运行状态数据", Crontab: "*/1 * * * *"})
	}

	if !db.Migrator().HasTable(&model.TaskHeartbeat{}) {
		if err = db.AutoMigrate(&model.TaskHeartbeat{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
		t, _ := time.Parse("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"))
		db.Create(&model.TaskHeartbeat{HeartbeatKey: "recycle_token", HeartbeatTime: t, HeartbeatEndTime: t})
		db.Create(&model.TaskHeartbeat{HeartbeatKey: "revoke_privileage", HeartbeatTime: t, HeartbeatEndTime: t})
		db.Create(&model.TaskHeartbeat{HeartbeatKey: "check_datasource", HeartbeatTime: t, HeartbeatEndTime: t})
		db.Create(&model.TaskHeartbeat{HeartbeatKey: "gather_dbmeta", HeartbeatTime: t, HeartbeatEndTime: t})
		db.Create(&model.TaskHeartbeat{HeartbeatKey: "gather_sensitive", HeartbeatTime: t, HeartbeatEndTime: t})
		db.Create(&model.TaskHeartbeat{HeartbeatKey: "collector_mysql_event", HeartbeatTime: t, HeartbeatEndTime: t})
		db.Create(&model.TaskHeartbeat{HeartbeatKey: "collector_redis_event", HeartbeatTime: t, HeartbeatEndTime: t})
		db.Create(&model.TaskHeartbeat{HeartbeatKey: "collector_oracle_event", HeartbeatTime: t, HeartbeatEndTime: t})
		db.Create(&model.TaskHeartbeat{HeartbeatKey: "collector_postgresql_event", HeartbeatTime: t, HeartbeatEndTime: t})
		db.Create(&model.TaskHeartbeat{HeartbeatKey: "collector_mongodb_event", HeartbeatTime: t, HeartbeatEndTime: t})
		db.Create(&model.TaskHeartbeat{HeartbeatKey: "collector_sqlserver_event", HeartbeatTime: t, HeartbeatEndTime: t})
	}

	if !db.Migrator().HasTable(&model.Favorite{}) {
		if err = db.AutoMigrate(&model.Favorite{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
	}

	if !db.Migrator().HasTable(&model.Privilege{}) {
		if err = db.AutoMigrate(&model.Privilege{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
	}

	if !db.Migrator().HasTable(&model.QueryLog{}) {
		if err = db.AutoMigrate(&model.QueryLog{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
	}

	if !db.Migrator().HasTable(&model.SensitiveRule{}) {
		if err = db.AutoMigrate(&model.SensitiveRule{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
		db.Create(&model.SensitiveRule{RuleKey: "mobile", RuleName: "手机号码", RuleType: "data", RuleExpress: "^1[356789]\\d{9}$|^\\+861\\d{10}$", Level: 1, Status: 1})
		db.Create(&model.SensitiveRule{RuleKey: "id_number", RuleName: "身份证号", RuleType: "data", RuleExpress: "^([1-9]\\d{5}[12]\\d{3}(0[1-9]|1[012])(0[1-9]|[12][0-9]|3[01])\\d{3}[0-9xX])$", Level: 1, Status: 1})
		db.Create(&model.SensitiveRule{RuleKey: "email", RuleName: "电子邮箱", RuleType: "data", RuleExpress: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}$", Level: 1, Status: 1})
		db.Create(&model.SensitiveRule{RuleKey: "bank_card", RuleName: "银行卡号", RuleType: "data", RuleExpress: "^[6]\\d{18}", Level: 1, Status: 1})
		db.Create(&model.SensitiveRule{RuleKey: "car_number", RuleName: "车牌号", RuleType: "data", RuleExpress: "^[\\x{4e00}-\\x{9fa2}][A-Z][0-9A-Z]{5}", Level: 0, Status: 1})
		db.Create(&model.SensitiveRule{RuleKey: "address", RuleName: "住址地址", RuleType: "data", RuleExpress: "[\\x{4e00}-\\x{9fa5}]{2,5}[市][\\x{4e00}-\\x{9fa5}]{2,5}[区](.+)[0-9]{1,4}[号]", Level: 0, Status: 1})
		db.Create(&model.SensitiveRule{RuleKey: "ip", RuleName: "IP地址", RuleType: "data", RuleExpress: "^(\\d{1,3}\\.){3}\\d{1,3}$", Level: 0, Status: 1})
		db.Create(&model.SensitiveRule{RuleKey: "ipport", RuleName: "IP端口服务", RuleType: "data", RuleExpress: "^(\\d{1,3}\\.){3}\\d{1,3}\\:\\d{2,6}$", Level: 0, Status: 1})
		db.Create(&model.SensitiveRule{RuleKey: "realname", RuleName: "姓名", RuleType: "data", RuleExpress: "^[\\x{4e00}-\\x{9fa2}]{2,3}$", Level: 1, Status: -1})
		db.Create(&model.SensitiveRule{RuleKey: "username", RuleName: "用户名", RuleType: "column", RuleExpress: "user|username|user_name", Level: 0, Status: -1})
		db.Create(&model.SensitiveRule{RuleKey: "password", RuleName: "密码", RuleType: "column", RuleExpress: "pass|password|pass_word", Level: 1, Status: -1})
	}

	if !db.Migrator().HasTable(&model.SensitiveMeta{}) {
		if err = db.AutoMigrate(&model.SensitiveMeta{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
	}

	if !db.Migrator().HasTable(&model.EventDescription{}) {
		if err = db.AutoMigrate(&model.EventDescription{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
	}

	if !db.Migrator().HasTable(&model.EventGlobal{}) {
		if err = db.AutoMigrate(&model.EventGlobal{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
	}

	if !db.Migrator().HasTable(&model.AlarmChannel{}) {
		if err = db.AutoMigrate(&model.AlarmChannel{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
		db.Create(&model.AlarmChannel{ID: 1, Name: "默认渠道", Description: "默认通知事件发送渠道", Enable: 1})
	}

	if !db.Migrator().HasTable(&model.AlarmLevel{}) {
		if err = db.AutoMigrate(&model.AlarmLevel{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
		db.Create(&model.AlarmLevel{ID: 1, LevelName: "停服", Description: "服务不可用", Enable: 1})
		db.Create(&model.AlarmLevel{ID: 2, LevelName: "严重", Description: "紧急的严重问题", Enable: 1})
		db.Create(&model.AlarmLevel{ID: 3, LevelName: "警告", Description: "不紧急的重要信息", Enable: 1})
		db.Create(&model.AlarmLevel{ID: 4, LevelName: "提醒", Description: "不紧急不严重需要关注的信息", Enable: 1})

	}

	if !db.Migrator().HasTable(&model.AlarmRule{}) {
		if err = db.AutoMigrate(&model.AlarmRule{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
		db.Create(&model.AlarmRule{Title: "MySQL数据源监测失败", EventType: "MySQL", EventKey: "datasourceCheck", AlarmRule: "!=", AlarmValue: "1", LevelId: 1, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MariaDB数据源监测失败", EventType: "MariaDB", EventKey: "datasourceCheck", AlarmRule: "!=", AlarmValue: "1", LevelId: 1, Enable: 1})
		db.Create(&model.AlarmRule{Title: "GreatSQL数据源监测失败", EventType: "GreatSQL", EventKey: "datasourceCheck", AlarmRule: "!=", AlarmValue: "1", LevelId: 1, Enable: 1})
		db.Create(&model.AlarmRule{Title: "TiDB数据源监测失败", EventType: "TiDB", EventKey: "datasourceCheck", AlarmRule: "!=", AlarmValue: "1", LevelId: 1, Enable: 1})
		db.Create(&model.AlarmRule{Title: "Doris数据源监测失败", EventType: "Doris", EventKey: "datasourceCheck", AlarmRule: "!=", AlarmValue: "1", LevelId: 1, Enable: 1})
		db.Create(&model.AlarmRule{Title: "OceanBase数据源监测失败", EventType: "OceanBase", EventKey: "datasourceCheck", AlarmRule: "!=", AlarmValue: "1", LevelId: 1, Enable: 1})
		db.Create(&model.AlarmRule{Title: "ClickHouse数据源监测失败", EventType: "ClickHouse", EventKey: "datasourceCheck", AlarmRule: "!=", AlarmValue: "1", LevelId: 1, Enable: 1})
		db.Create(&model.AlarmRule{Title: "Oracle数据源监测失败", EventType: "Oracle", EventKey: "datasourceCheck", AlarmRule: "!=", AlarmValue: "1", LevelId: 1, Enable: 1})
		db.Create(&model.AlarmRule{Title: "PostgreSQL数据源监测失败", EventType: "PostgreSQL", EventKey: "datasourceCheck", AlarmRule: "!=", AlarmValue: "1", LevelId: 1, Enable: 1})
		db.Create(&model.AlarmRule{Title: "SQLServer数据源监测失败", EventType: "SQLServer", EventKey: "datasourceCheck", AlarmRule: "!=", AlarmValue: "1", LevelId: 1, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MongoDB数据源监测失败", EventType: "MongoDB", EventKey: "datasourceCheck", AlarmRule: "!=", AlarmValue: "1", LevelId: 1, Enable: 1})
		db.Create(&model.AlarmRule{Title: "Redis数据源监测失败", EventType: "Redis", EventKey: "datasourceCheck", AlarmRule: "!=", AlarmValue: "1", LevelId: 1, Enable: 1})

		db.Create(&model.AlarmRule{Title: "MySQL数据库无法连接", EventType: "MySQL", EventKey: "connect", AlarmRule: "!=", AlarmValue: "1", LevelId: 1, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MySQL等待事件过高", EventType: "MySQL", EventKey: "threadsWait", AlarmRule: ">", AlarmValue: "5", LevelId: 2, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MySQL QPS过高", EventType: "MySQL", EventKey: "queries", AlarmRule: ">", AlarmValue: "1000", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MySQL连接数过高", EventType: "MySQL", EventKey: "threadsConnected", AlarmRule: ">", AlarmValue: "1000", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MySQL活动会话过高", EventType: "MySQL", EventKey: "threadsRunning", AlarmRule: ">", AlarmValue: "20", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MySQL活动事务量过高", EventType: "MySQL", EventKey: "activeTrx", AlarmRule: ">", AlarmValue: "10", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MySQL存在长事务", EventType: "MySQL", EventKey: "longTrx", AlarmRule: ">", AlarmValue: "0", LevelId: 4, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MySQL存在长时间运行的SQL", EventType: "MySQL", EventKey: "longQuery", AlarmRule: ">", AlarmValue: "0", LevelId: 4, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MySQL写入流量过高", EventType: "MySQL", EventKey: "bytesReceived", AlarmRule: ">", AlarmValue: "10000", LevelId: 4, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MySQL读取流量过高", EventType: "MySQL", EventKey: "bytesSent", AlarmRule: ">", AlarmValue: "10000", LevelId: 4, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MySQL慢查询过多", EventType: "MySQL", EventKey: "slowQueries", AlarmRule: ">", AlarmValue: "100", LevelId: 4, Enable: 1})

		db.Create(&model.AlarmRule{Title: "GreatSQL数据库无法连接", EventType: "GreatSQL", EventKey: "connect", AlarmRule: "!=", AlarmValue: "1", LevelId: 1, Enable: 1})
		db.Create(&model.AlarmRule{Title: "GreatSQL等待事件过高", EventType: "GreatSQL", EventKey: "threadsWait", AlarmRule: ">", AlarmValue: "5", LevelId: 2, Enable: 1})
		db.Create(&model.AlarmRule{Title: "GreatSQL QPS过高", EventType: "GreatSQL", EventKey: "queries", AlarmRule: ">", AlarmValue: "1000", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "GreatSQL连接数过高", EventType: "GreatSQL", EventKey: "threadsConnected", AlarmRule: ">", AlarmValue: "1000", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "GreatSQL活动会话过高", EventType: "GreatSQL", EventKey: "threadsRunning", AlarmRule: ">", AlarmValue: "20", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "GreatSQL活动事务量过高", EventType: "GreatSQL", EventKey: "activeTrx", AlarmRule: ">", AlarmValue: "10", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "GreatSQL存在长事务", EventType: "GreatSQL", EventKey: "longTrx", AlarmRule: ">", AlarmValue: "0", LevelId: 4, Enable: 1})
		db.Create(&model.AlarmRule{Title: "GreatSQL存在长时间运行的SQL", EventType: "GreatSQL", EventKey: "longQuery", AlarmRule: ">", AlarmValue: "0", LevelId: 4, Enable: 1})
		db.Create(&model.AlarmRule{Title: "GreatSQL写入流量过高", EventType: "GreatSQL", EventKey: "bytesReceived", AlarmRule: ">", AlarmValue: "10000", LevelId: 4, Enable: 1})
		db.Create(&model.AlarmRule{Title: "GreatSQL读取流量过高", EventType: "GreatSQL", EventKey: "bytesSent", AlarmRule: ">", AlarmValue: "10000", LevelId: 4, Enable: 1})
		db.Create(&model.AlarmRule{Title: "GreatSQL慢查询过多", EventType: "GreatSQL", EventKey: "slowQueries", AlarmRule: ">", AlarmValue: "100", LevelId: 4, Enable: 1})

		db.Create(&model.AlarmRule{Title: "MariaDB数据库无法连接", EventType: "MariaDB", EventKey: "connect", AlarmRule: "!=", AlarmValue: "1", LevelId: 1, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MariaDB等待事件过高", EventType: "MariaDB", EventKey: "threadsWait", AlarmRule: ">", AlarmValue: "5", LevelId: 2, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MariaDB QPS过高", EventType: "MariaDB", EventKey: "queries", AlarmRule: ">", AlarmValue: "1000", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MariaDB连接数过高", EventType: "MariaDB", EventKey: "threadsConnected", AlarmRule: ">", AlarmValue: "1000", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MariaDB活动会话过高", EventType: "MariaDB", EventKey: "threadsRunning", AlarmRule: ">", AlarmValue: "20", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MariaDB活动事务量过高", EventType: "MariaDB", EventKey: "activeTrx", AlarmRule: ">", AlarmValue: "10", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MariaDB存在长事务", EventType: "MariaDB", EventKey: "longTrx", AlarmRule: ">", AlarmValue: "0", LevelId: 4, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MariaDB存在长时间运行的SQL", EventType: "MariaDB", EventKey: "longQuery", AlarmRule: ">", AlarmValue: "0", LevelId: 4, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MariaDB写入流量过高", EventType: "MariaDB", EventKey: "bytesReceived", AlarmRule: ">", AlarmValue: "10000", LevelId: 4, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MariaDB读取流量过高", EventType: "MariaDB", EventKey: "bytesSent", AlarmRule: ">", AlarmValue: "10000", LevelId: 4, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MariaDB慢查询过多", EventType: "MariaDB", EventKey: "slowQueries", AlarmRule: ">", AlarmValue: "100", LevelId: 4, Enable: 1})

		db.Create(&model.AlarmRule{Title: "PostgreSQL数据库无法连接", EventType: "PostgreSQL", EventKey: "connect", AlarmRule: "!=", AlarmValue: "1", LevelId: 1, Enable: 1})
		db.Create(&model.AlarmRule{Title: "PostgreSQL连接数过高", EventType: "PostgreSQL", EventKey: "connections", AlarmRule: ">", AlarmValue: "1000", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "PostgreSQL活动SQL数过高", EventType: "PostgreSQL", EventKey: "activeSQL", AlarmRule: ">", AlarmValue: "20", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "PostgreSQL等待事件过多", EventType: "PostgreSQL", EventKey: "waitEvent", AlarmRule: ">", AlarmValue: "5", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "PostgreSQL发现锁等待", EventType: "PostgreSQL", EventKey: "locks", AlarmRule: ">", AlarmValue: "0", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "PostgreSQL检测到死锁", EventType: "PostgreSQL", EventKey: "deadlocks", AlarmRule: ">", AlarmValue: "0", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "PostgreSQL存在长时间运行的SQL", EventType: "PostgreSQL", EventKey: "longQuery", AlarmRule: ">", AlarmValue: "0", LevelId: 4, Enable: 1})
		db.Create(&model.AlarmRule{Title: "PostgreSQL存在长事务", EventType: "PostgreSQL", EventKey: "longTransaction", AlarmRule: ">", AlarmValue: "0", LevelId: 4, Enable: 1})

		db.Create(&model.AlarmRule{Title: "Oracle数据库无法连接", EventType: "Oracle", EventKey: "connect", AlarmRule: "!=", AlarmValue: "1", LevelId: 1, Enable: 1})
		db.Create(&model.AlarmRule{Title: "Oracle等待会话数过高", EventType: "Oracle", EventKey: "sessionWait", AlarmRule: ">", AlarmValue: "10", LevelId: 2, Enable: 1})
		db.Create(&model.AlarmRule{Title: "Oracle活动会话数过高", EventType: "Oracle", EventKey: "sessionActive", AlarmRule: ">", AlarmValue: "20", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "Oracle连接会话数过高", EventType: "Oracle", EventKey: "sessionTotal", AlarmRule: ">", AlarmValue: "1000", LevelId: 3, Enable: 1})

		db.Create(&model.AlarmRule{Title: "SQLServer数据库无法连接", EventType: "SQLServer", EventKey: "connect", AlarmRule: "!=", AlarmValue: "1", LevelId: 1, Enable: 1})
		db.Create(&model.AlarmRule{Title: "SQLServer等待进程数过高", EventType: "SQLServer", EventKey: "processWait", AlarmRule: ">", AlarmValue: "10", LevelId: 2, Enable: 1})
		db.Create(&model.AlarmRule{Title: "SQLServer活动进程数过高", EventType: "SQLServer", EventKey: "processRunning", AlarmRule: ">", AlarmValue: "20", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "SQLServer进程数过高", EventType: "SQLServer", EventKey: "process", AlarmRule: ">", AlarmValue: "1000", LevelId: 3, Enable: 1})

		db.Create(&model.AlarmRule{Title: "MongoDB数据库无法连接", EventType: "MongoDB", EventKey: "connect", AlarmRule: "!=", AlarmValue: "1", LevelId: 1, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MongoDB当前连接数过高", EventType: "MongoDB", EventKey: "connectionsCurrent", AlarmRule: ">", AlarmValue: "1000", LevelId: 2, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MongoDB可用连接不足", EventType: "MongoDB", EventKey: "connectionsAvailable", AlarmRule: "<", AlarmValue: "500", LevelId: 2, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MongoDB网络请求数过高", EventType: "MongoDB", EventKey: "networkNumRequests", AlarmRule: ">", AlarmValue: "2000", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MongoDB写入流量过高", EventType: "MongoDB", EventKey: "networkBytesIn", AlarmRule: ">", AlarmValue: "10000000", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MongoDB读取流量过高", EventType: "MongoDB", EventKey: "networkBytesOut", AlarmRule: ">", AlarmValue: "10000000", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MongoDB操作量过高", EventType: "MongoDB", EventKey: "opcounters", AlarmRule: ">", AlarmValue: "2000", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MongoDB查询操作量过高", EventType: "MongoDB", EventKey: "opcountersQuery", AlarmRule: ">", AlarmValue: "1000", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MongoDB写入操作量过高", EventType: "MongoDB", EventKey: "opcountersInsert", AlarmRule: ">", AlarmValue: "1000", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MongoDB更新操作量过高", EventType: "MongoDB", EventKey: "opcountersUpdate", AlarmRule: ">", AlarmValue: "1000", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "MongoDB删除操作量过高", EventType: "MongoDB", EventKey: "opcountersDelete", AlarmRule: ">", AlarmValue: "1000", LevelId: 3, Enable: 1})

		db.Create(&model.AlarmRule{Title: "Redis数据库无法连接", EventType: "Redis", EventKey: "connect", AlarmRule: "!=", AlarmValue: "1", LevelId: 1, Enable: 1})
		db.Create(&model.AlarmRule{Title: "Redis客户端连接数过高", EventType: "Redis", EventKey: "connectedClients", AlarmRule: ">", AlarmValue: "1000", LevelId: 2, Enable: 1})
		db.Create(&model.AlarmRule{Title: "Redis连接阻塞数过高", EventType: "Redis", EventKey: "blockedClients", AlarmRule: ">", AlarmValue: "10", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "Redis每秒请求数过高", EventType: "Redis", EventKey: "instantaneousOpsPerSec", AlarmRule: ">", AlarmValue: "1000", LevelId: 3, Enable: 1})
		db.Create(&model.AlarmRule{Title: "Redis内存使用率过高", EventType: "Redis", EventKey: "usedMemoryPct", AlarmRule: ">", AlarmValue: "70", LevelId: 3, Enable: 1})

	}

	if !db.Migrator().HasTable(&model.AlarmEvent{}) {
		if err = db.AutoMigrate(&model.AlarmEvent{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
	}

	if !db.Migrator().HasTable(&model.AlarmSendLog{}) {
		if err = db.AutoMigrate(&model.AlarmSendLog{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
	}

	if !db.Migrator().HasTable(&model.AlarmSuggest{}) {
		if err = db.AutoMigrate(&model.AlarmSuggest{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
	}

	if !db.Migrator().HasTable(&model.AlarmTrack{}) {
		if err = db.AutoMigrate(&model.AlarmTrack{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
	}

	return db
}

func InitConnect() *sql.DB {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
			os.Exit(0)

		}
	}()
	ds := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", setting.Setting.User, setting.Setting.Password, setting.Setting.Host, setting.Setting.Port, setting.Setting.Database)
	db, err := sql.Open("mysql", ds)
	if err != nil {
		log.Error(fmt.Sprintln("Init mysql connect err,", err))
		panic(fmt.Sprintln("Init mysql connect err,", err))
	}
	if err := db.Ping(); err != nil {
		log.Error(fmt.Sprintln("Init mysql ping err,", err))
		panic(fmt.Sprintln("Init mysql ping err,", err))
	}

	return db
}

func InitCk() *gorm.DB {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
			os.Exit(0)

		}
	}()
	dsn := fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s?dial_timeout=10s&read_timeout=20s", setting.Setting.ClickhouseUser, setting.Setting.ClickhousePassword, setting.Setting.ClickhouseHost, setting.Setting.ClickhousePort, setting.Setting.ClickhouseDatabase)
	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintln("open database clickhouse error.", zap.Error(err)))
	}

	//自动迁移 (这是GORM自动创建表的一种方式--译者注)
	//clickhouse: set allow_suspicious_fixed_string_types = 1
	if !db.Migrator().HasTable(&model.Event{}) {
		if err = db.AutoMigrate(&model.Event{}); err != nil {
			panic(fmt.Sprintln("auto migrate clickhouse table error.", zap.Error(err)))
		}
	}
	return db

}

func QueryAll(sql string) ([]map[string]interface{}, error) {
	rows, err := SQL.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var list []map[string]interface{}
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			continue
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			v := values[i]
			b, ok := v.([]byte)
			if ok {
				entry[col] = string(b)
			} else {
				entry[col] = v
			}
		}
		list = append(list, entry)
	}
	return list, nil
}

// 使用option结构体创建可选参数
type Option struct {
	f func(*options)
}

type options struct {
	driver   string
	host     string
	port     string
	username string
	password string
	database string
	sid      string
	timeout  int
}

func WithDriver(driver string) Option {
	return Option{func(op *options) {
		op.driver = driver
	}}
}
func WithHost(host string) Option {
	return Option{func(op *options) {
		op.host = host
	}}
}
func WithPort(port string) Option {
	return Option{func(op *options) {
		op.port = port
	}}
}
func WithUsername(username string) Option {
	return Option{func(op *options) {
		op.username = username
	}}
}
func WithPassword(password string) Option {
	return Option{func(op *options) {
		op.password = password
	}}
}
func WithDatabase(database string) Option {
	return Option{func(op *options) {
		op.database = database
	}}
}
func WithSid(sid string) Option {
	return Option{func(op *options) {
		op.sid = sid
	}}
}
func WithTimeout(timeout int) Option {
	return Option{func(op *options) {
		op.timeout = timeout
	}}
}

// 使用结构体动态传入参数
func Connect(ops ...Option) (*sql.DB, error) {
	//set option
	opt := &options{}
	for _, do := range ops {
		do.f(opt)
	}
	//不同数据库构造不同的url
	var url string
	if opt.driver == "mysql" {
		url = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=5s&readTimeout=10s", opt.username, opt.password, opt.host, opt.port, opt.database)
	}
	if opt.driver == "postgres" {
		url = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", opt.host, opt.port, opt.username, opt.password, opt.database)
	}
	if opt.driver == "clickhouse" {
		url = fmt.Sprintf("tcp://%s:%s/%s?username=%s&password=%s&read_timeout=30s", opt.host, opt.port, opt.database, opt.username, opt.password)
	}
	if opt.driver == "oracle" {
		url = fmt.Sprintf(`user="%s" password="%s" connectString="%s:%s/%s"`, opt.username, opt.password, opt.host, opt.port, opt.sid)
	}
	if opt.driver == "mssql" {
		url = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;encrypt=disable;connection timeout=6;", opt.host, opt.username, opt.password, opt.port, opt.database)
	}
	//连接数据库
	db, err := sql.Open(opt.driver, url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func Execute(db *sql.DB, sql string) (rowsAffected int64, err error) {
	res, err := db.Exec(sql)
	if err != nil {
		return 0, err
	}
	rowsAffected, _ = res.RowsAffected()
	return rowsAffected, nil
}

func QueryRemote(db *sql.DB, sql string) ([]map[string]interface{}, error) {
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var list []map[string]interface{}
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			continue
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			v := values[i]
			b, ok := v.([]byte)
			if ok {
				entry[col] = string(b)
				//entry[col] = b
			} else {
				entry[col] = v
			}
		}
		list = append(list, entry)
	}
	return list, nil
}

/*
QueryRemoteNew方法会返回columns，columns顺序是稳定的
*/
func QueryRemoteNew(db *sql.DB, sql string) ([]string, []map[string]interface{}, error) {
	rows, err := db.Query(sql)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}
	count := len(columns)
	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var list []map[string]interface{}
	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			continue
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			fmt.Println(col)
			v := values[i]
			b, ok := v.([]byte)
			if ok {
				entry[col] = string(b)
				//entry[col] = b
			} else {
				entry[col] = v
			}
		}
		list = append(list, entry)
	}
	return columns, list, nil
}

func InitRedis() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", setting.Setting.RedisHost, setting.Setting.RedisPort),
		Password:     setting.Setting.RedisPassword, // no password set
		DB:           0,                             // use default DB
		PoolSize:     128,
		ReadTimeout:  time.Millisecond * time.Duration(500),
		WriteTimeout: time.Millisecond * time.Duration(500),
		IdleTimeout:  time.Second * time.Duration(86400),
	})
	_, err := redisClient.Ping().Result()
	if err != nil {
		log.Error("open redis client error", zap.Error(err))
		panic(fmt.Sprintln("open redis client error.", zap.Error(err)))
	}
	return redisClient
}

func InitRedisCluster(host, port, password string) (*redis.ClusterClient, error) {
	redisClusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        []string{fmt.Sprintf("%s:%s", host, port)},
		Password:     password, // no password set
		PoolSize:     1000,
		ReadTimeout:  time.Millisecond * time.Duration(200),
		WriteTimeout: time.Millisecond * time.Duration(200),
		IdleTimeout:  time.Second * time.Duration(600),
	})
	_, err := redisClusterClient.Ping().Result()
	if err != nil {
		return nil, err
	}
	return redisClusterClient, nil
}
