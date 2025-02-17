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
	"dbmcloud/src/database"
	"dbmcloud/src/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Result struct {
	ID           int64  `gorm:"primarykey" json:"id"`
	TaskName     string `gorm:"size:50" json:"task_name"`
	ScheduleType string `gorm:"size:30" json:"schedule_type"`
	NextTime     string `gorm:"size:30" json:"next_time"`
	Status       string `gorm:"size:30" json:"status"`
	Enable       int    `gorm:"default:1" json:"enable"`
	CategoryName string `gorm:"size:300" json:"category_name"`
}

func TaskList(c *gin.Context) {
	var db = database.DB
	method := c.Request.Method
	if method == "GET" {

		sql := "SELECT a.id  ,a.task_name,a.schedule_type ,a.status ,a.next_time,a.enable, b.name as category_name FROM task a LEFT JOIN task_type b ON a.type_id = b.id where 1=1 "
		if c.Query("status") != "" {
			sql = fmt.Sprintf("%s and status='%s' ", sql, c.Query("status"))
		}
		if c.Query("task_name") != "" {
			sql = fmt.Sprintf("%s and task_name like '%s%s%s' ", sql, "%", c.Query("task_name"), "%")
		}
		sorterMap := make(map[string]string)
		sorterData := c.Query("sorter")
		json.Unmarshal([]byte(sorterData), &sorterMap)
		for sortField, sortOrder := range sorterMap {
			if sortField != "" && sortOrder != "" {
				sql = fmt.Sprintf("%s order by %s %s ", sql, sortField, strings.Replace(sortOrder, "end", "", 1))
			}
		}

		var dataList []Result

		result := db.Raw(sql).Scan(&dataList)
		if result.Error != nil {
			c.JSON(200, gin.H{"success": false, "msg": "Query Error:" + result.Error.Error()})
			return
		}

		taskCount, _ := database.QueryAll("SELECT COUNT(*) as count FROM task_run WHERE gmt_create>= DATE_FORMAT(NOW(), '%Y-%m-%d 00:00:00') LIMIT 1")
		taskSuccessCount, _ := database.QueryAll("SELECT COUNT(*) as count FROM task_run WHERE gmt_create>= DATE_FORMAT(NOW(), '%Y-%m-%d 00:00:00') AND run_status='success' LIMIT 1")
		taskRunningCount, _ := database.QueryAll("SELECT COUNT(*) as count FROM task_run WHERE gmt_create>= DATE_FORMAT(NOW(), '%Y-%m-%d 00:00:00') AND run_status='running' LIMIT 1")
		taskFailedCount, _ := database.QueryAll("SELECT COUNT(*) as count FROM task_run WHERE gmt_create>= DATE_FORMAT(NOW(), '%Y-%m-%d 00:00:00') AND run_status='failed' LIMIT 1")
		taskWaitCount, _ := database.QueryAll("SELECT COUNT(*) as count FROM task_run WHERE gmt_create>= DATE_FORMAT(NOW(), '%Y-%m-%d 00:00:00') AND run_status='waiting' LIMIT 1")

		TaskSuccessPct := 1.00
		if utils.StrToInt(taskCount[0]["count"].(string)) != 0 {
			TaskSuccessPct = 1 - utils.StrToFloat64(taskFailedCount[0]["count"].(string))/utils.StrToFloat64(taskCount[0]["count"].(string))
		}

		c.JSON(http.StatusOK, gin.H{
			"success":          true,
			"msg":              "OK",
			"data":             dataList,
			"total":            len(dataList),
			"taskCount":        taskCount[0]["count"],
			"taskSuccessCount": taskSuccessCount[0]["count"],
			"taskRunningCount": taskRunningCount[0]["count"],
			"taskFailedCount":  taskFailedCount[0]["count"],
			"taskWaitCount":    taskWaitCount[0]["count"],
			"taskSuccessPct":   TaskSuccessPct * 100,
		})
		return

	}
}
