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
	"dbmcloud/src/database"
	"dbmcloud/src/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DoWriteLog(c *gin.Context) {
	params := make(map[string]string)
	c.BindJSON(&params)
	if len(params) == 0 {
		c.JSON(http.StatusOK, gin.H{"success": false, "msg": "params error."})
		return
	}

	datasourceType := params["datasource_type"]
	datasource := params["datasource"]
	databaseName := params["database"]
	sql := params["sql"]
	queryType := params["query_type"]
	username, _ := c.Get("username")
	WriteLog(username.(string), datasourceType, datasource, queryType, "export", databaseName, success, 0, sql, "导出完成")
	c.JSON(http.StatusOK, gin.H{"success": true, "msg": "导出完成"})
	return
}

func WriteLog(username string, datasourceType string, datasource string, queryType string, sqlType string, databaseName string, status string, times int64, content string, doResult string) {
	var record model.QueryLog
	record.Username = username
	record.DatasourceType = datasourceType
	record.Datasource = datasource
	record.QueryType = queryType
	record.SqlType = sqlType
	record.Database = databaseName
	record.Status = status
	record.Times = times
	record.Content = content
	record.Result = doResult
	result := database.DB.Create(&record)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return
	}
	return
}
