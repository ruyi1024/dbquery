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

package datasource

import (
	"dbmcloud/src/database"
	"dbmcloud/src/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func EnvList(c *gin.Context) {
	var db = database.DB
	method := c.Request.Method
	if method == "GET" {
		var dataList []model.Env
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
		var record model.Env
		c.BindJSON(&record)
		result := database.DB.Create(&record)
		if result.Error != nil {
			c.JSON(200, gin.H{"success": false, "msg": "Insert Error: " + result.Error.Error()})
			return
		}
		c.JSON(200, gin.H{"success": true})
		return

	}

	if method == "PUT" {
		var record model.Env
		c.BindJSON(&record)
		result := database.DB.Model(&record).Omit("id").Where("id = ?", record.Id).Updates(record)
		if result.Error != nil {
			c.JSON(200, gin.H{"success": false, "msg": "Update Error: " + result.Error.Error()})
			return
		}
		c.JSON(200, gin.H{"success": true})
		return
	}

	if method == "DELETE" {
		var record model.Env
		c.BindJSON(&record)
		result := database.DB.Model(&model.Env{}).Where("id = ?", record.Id).Delete(record)
		if result.Error != nil {
			c.JSON(200, gin.H{"success": false, "msg": "Delete Error:" + result.Error.Error()})
			return
		}
		c.JSON(200, gin.H{"success": true})
		return
	}
}
