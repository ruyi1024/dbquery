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
		result := db.Create(&model.Users{Id: 1, Username: "admin", ChineseName: "Administrator", Password: "a8a0d32f1abefd3fa996321d5e72c6d6", Admin: true})
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
		db.Create(&model.Idc{Id: 1, IdcKey: "default", IdcName: "Default", Description: "Default IDC"})
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
		db.Create(&model.Datasource{Id: 1, Name: "DBQuery-MySQL", GroupName: "default", Idc: "default", Env: "prod", Type: "MySQL", Host: setting.Setting.Host, Port: setting.Setting.Port, User: setting.Setting.User, Pass: aesPassword, Enable: 1, DbmetaEnable: 1, ExecuteEnable: 1})

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
	}

	if !db.Migrator().HasTable(&model.Favorite{}) {
		if err = db.AutoMigrate(&model.Favorite{}); err != nil {
			log.Error("db sync error.", zap.Error(err))
		}
	}

	if !db.Migrator().HasTable(&model.QueryLog{}) {
		if err = db.AutoMigrate(&model.QueryLog{}); err != nil {
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
