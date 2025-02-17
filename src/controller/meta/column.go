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

package meta

import (
	"dbmcloud/src/database"
	"dbmcloud/src/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ColumnList(c *gin.Context) {
	var db = database.DB
	method := c.Request.Method
	if method == "GET" {
		if c.Query("datasource_type") != "" {
			db = db.Where("datasource_type=?", c.Query("datasource_type"))
		}
		if c.Query("host") != "" {
			db = db.Where("host=?", c.Query("host"))
		}
		if c.Query("port") != "" {
			db = db.Where("port=?", c.Query("port"))
		}
		if c.Query("database_name") != "" {
			db = db.Where("database_name like ? ", c.Query("database_name")+"%")
		}
		if c.Query("table_name") != "" {
			db = db.Where("table_name like ? ", c.Query("table_name")+"%")
		}
		if c.Query("column_name") != "" {
			db = db.Where("column_name like ? ", c.Query("column_name")+"%")
		}
		sorterMap := make(map[string]string)
		sorterData := c.Query("sorter")
		json.Unmarshal([]byte(sorterData), &sorterMap)
		for sortField, sortOrder := range sorterMap {
			if sortField != "" && sortOrder != "" {
				db = db.Order(fmt.Sprintf("%s %s", sortField, strings.Replace(sortOrder, "end", "", 1)))
			}
		}

		var dataList []model.MetaColumn
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
}
