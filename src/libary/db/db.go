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

/*
封装了所有数据库的连接
*/

package db

import (
	"database/sql"
	"fmt"

	//_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
	_ "github.com/lib/pq"
)

var err error

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
		url = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?timeout=3s&readTimeout=5s", opt.username, opt.password, opt.host, opt.port, opt.database)
	}
	if opt.driver == "godror" {
		url = fmt.Sprintf(`user="%s" password="%s" connectString="%s:%s/%s"`, opt.username, opt.password, opt.host, opt.port, opt.sid)
	}
	if opt.driver == "mssql" {
		url = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;encrypt=disable;connection timeout=6;", opt.host, opt.username, opt.password, opt.port, opt.database)
	}
	if opt.driver == "postgres" {
		url = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", opt.host, opt.port, opt.username, opt.password, opt.database)
	}
	if opt.driver == "clickhouse" {
		url = fmt.Sprintf("tcp://%s:%s/%s?username=%s&password=%s&read_timeout=30s", opt.host, opt.port, opt.database, opt.username, opt.password)
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

func QueryOne(db *sql.DB, sql string) (data string, err error) {
	row := db.QueryRow(sql)
	if err := row.Scan(); err != nil {
		return "", nil
	}
	return
}

func QueryAll(db *sql.DB, sql string) ([]map[string]interface{}, error) {
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
QueryAllNew方法会返回columns，columns顺序是稳定的
*/
func QueryAllNew(db *sql.DB, sql string) ([]string, []map[string]interface{}, error) {
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
