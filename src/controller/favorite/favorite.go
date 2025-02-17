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

package favorite

import (
	"dbmcloud/src/database"
	"dbmcloud/src/model"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
获取当前用户
*/
func getLoginUsername(c *gin.Context) string {
	//var c *gin.Context
	userinfo, _ := c.Get("loginUser")  //获取用户cookie
	data, _ := json.Marshal(&userinfo) //userinfo返回结果是struct model.Users,需要转换成map
	userMap := make(map[string]interface{})
	json.Unmarshal(data, &userMap)
	loginUsername := userMap["username"].(string)
	return loginUsername
}

func List(c *gin.Context) {
	var db = database.DB
	method := c.Request.Method
	if method == "GET" {
		var dataList []model.Favorite
		db.Where("username=?", getLoginUsername(c))
		if c.Query("datasource_type") != "" {
			db = db.Where("datasource_type=?", c.Query("datasource_type"))
		}
		if c.Query("datasource") != "" {
			db = db.Where("datasource=?", c.Query("datasource"))
		}
		if c.Query("database_name") != "" {
			db = db.Where("database_name=?", c.Query("database_name"))
		}
		db.Order("id desc")
		result := db.Find(&dataList)
		if result.Error != nil {
			c.JSON(200, gin.H{"success": false, "msg": "Query Error:" + result.Error.Error()})
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

		//判断是否收藏过
		//没有收藏则添加收藏
		var record model.Favorite
		var dataList []model.Favorite
		record.Username = getLoginUsername(c)
		c.BindJSON(&record)
		database.DB.Where("username=?", getLoginUsername(c)).Where("datasource=?", record.Datasource).Where("database_name=?", record.DatabaseName).Where("content=?", record.Content).Find(&dataList)
		fmt.Println(dataList)
		if (len(dataList)) > 0 {
			c.JSON(200, gin.H{"success": true})
			return
		}
		result := database.DB.Create(&record)
		if result.Error != nil {
			c.JSON(200, gin.H{"success": false, "msg": "Insert Error: " + result.Error.Error()})
			return
		}
		c.JSON(200, gin.H{"success": true})
		return

	}

	if method == "PUT" {
		var record model.Favorite
		c.BindJSON(&record)
		result := database.DB.Model(&record).Omit("id").Where("id = ?", record.ID).Updates(record)
		if result.Error != nil {
			c.JSON(200, gin.H{"success": false, "msg": "Update Error: " + result.Error.Error()})
			return
		}
		c.JSON(200, gin.H{"success": true})
		return
	}

	if method == "DELETE" {
		var record model.Favorite
		c.BindJSON(&record)
		result := database.DB.Model(&model.Favorite{}).Where("id = ?", record.ID).Delete(record)
		if result.Error != nil {
			c.JSON(200, gin.H{"success": false, "msg": "Delete Error:" + result.Error.Error()})
			return
		}
		c.JSON(200, gin.H{"success": true})
		return
	}
}
