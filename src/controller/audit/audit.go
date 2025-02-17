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

package audit

import (
	"dbmcloud/src/database"
	"dbmcloud/src/model"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type FindQuerLog struct {
	ID             int64     `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time `gorm:"column:gmt_created" json:"gmt_created"`
	Username       string    `gorm:"size:100" json:"username"`
	DatasourceType string    `gorm:"size:200" json:"datasource_type"`
	Datasource     string    `gorm:"size:50" json:"datasource"`
	Database       string    `gorm:"size:50" json:"database"`
	QueryType      string    `gorm:"size:50" json:"query_type"`
	SqlType        string    `gorm:"size:50" json:"sql_type"`
	Status         string    `gorm:"size:30" json:"status"`
	Times          int64     `gorm:"size:10" json:"times"`
	Content        string    `gorm:"size:1000" json:"content"`
	Result         string    `gorm:"size:500" json:"result"`
}

func GetQueryLog(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))   // 当前页数
	offset, _ := strconv.Atoi(c.Query("offset")) // 分页
	sorterField := c.Query("sorterField")        // 排序字段
	sorterOrder := c.Query("sorterOrder")        // 排序方式
	searchValue := c.Query("keyword")            // 搜索

	order := "DESC"
	if sorterOrder == "ascend" {
		order = "ASC"
	}
	if sorterField == "" {
		sorterField = "id"
	}

	// get db data
	var datalist []FindQuerLog
	result := database.DB.Model(&model.QueryLog{})
	if searchValue != "" {
		result.Where("username LIKE ?", "%"+searchValue+"%")
	}
	result.Order(fmt.Sprintf("%s %s", SnakeString(sorterField), order)).Limit(limit).Offset(offset).Find(&datalist)
	if result.Error != nil {
		c.JSON(200, gin.H{"success": false, "msg": "query db users error " + result.Error.Error()})
		return
	}
	var total int64
	database.DB.Model(&model.QueryLog{}).Count(&total)
	c.JSON(200, gin.H{"success": true, "data": datalist, "total": total})
	//c.JSON(200, gin.H{"userinfo": count})
	//c.String(200, "user list <br> session:"+fmt.Sprintf("%#v", session))
}

// snake string, XxYy to xx_yy , XxYY to xx_yy
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}
