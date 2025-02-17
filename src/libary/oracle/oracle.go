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

package oracle

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/godror/godror"
)

var err error

func Connect(host, port, username, password, sid string) (*sql.DB, error) {
	db, err := sql.Open("godror", fmt.Sprintf(`user="%s" password="%s" connectString="%s:%s/%s"`, username, password, host, port, sid))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func Execute(db *sql.DB, sql string) (err error) {
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}
	return
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
			col = strings.ToLower(col)
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
