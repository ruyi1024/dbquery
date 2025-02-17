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

package users

import (
	"dbmcloud/log"
	"dbmcloud/setting"
	"dbmcloud/src/database"
	"dbmcloud/src/model"
	"dbmcloud/src/utils"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FindUser struct {
	Id          int64     `json:"id"`
	Username    string    `json:"username"`
	ChineseName string    `json:"chineseName"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Admin       bool      `json:"admin"`
	Remark      string    `json:"remark"`
}

func GetUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))   // 当前页数
	offset, _ := strconv.Atoi(c.Query("offset")) // 分页
	sorterField := c.Query("sorterField")        // 排序字段
	sorterOrder := c.Query("sorterOrder")        // 排序方式
	searchValue := c.Query("keyword")            // 搜索

	log.Info("debug get user -->", zap.Int("limit", limit), zap.Int("offset", offset), zap.Any("sorterField:", sorterField))
	log.Info("debug get user searchValue -->", zap.Any("searchValue:", searchValue))

	order := "ASC"
	if sorterOrder == "descend" {
		order = "DESC"
	}
	if sorterField == "" {
		sorterField = "username"
	}

	// get db data
	var users []FindUser
	result := database.DB.Model(&model.Users{})
	if searchValue != "" {
		result.Where("username LIKE ?", "%"+searchValue+"%")
	}

	result.Order(fmt.Sprintf("%s %s", SnakeString(sorterField), order)).Limit(limit).Offset(offset).Find(&users)
	//result.Order(fmt.Sprintf("%s %s", SnakeString(sorterField), order)).Find(&users)
	if result.Error != nil {
		c.JSON(200, gin.H{"success": false, "msg": "query db users error " + result.Error.Error()})
		return
	}

	var total int64
	database.DB.Model(&model.Users{}).Count(&total)
	c.JSON(200, gin.H{"success": true, "data": users, "total": total})
	//c.JSON(200, gin.H{"userinfo": count})
	//c.String(200, "user list <br> session:"+fmt.Sprintf("%#v", session))
}

func PostUser(c *gin.Context) {
	var user model.Users
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "msg": "query users error " + err.Error()})
		return
	}
	user.Password = utils.Md5plus(user.Password, setting.Setting.Md5Iteration)
	result := database.DB.Create(&user)
	if result.Error != nil {
		c.JSON(200, gin.H{"success": false, "msg": "insert db users error " + result.Error.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
	return
}

func PutUser(c *gin.Context) {
	var user model.Users
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "msg": "query users error " + err.Error()})
		return
	}
	log.Info("user info ", zap.Any("%#v", user))

	//fmt.Print(user.Admin)
	//result := database.DB.Model(&model.Users{}).Omit("id", "password").Where("id = ?", user.ID).Updates(user)
	result := database.DB.Model(&model.Users{}).Omit("id", "password").Where("id = ?", user.Id).Updates(map[string]interface{}{"username": user.Username, "chinese_name": user.ChineseName, "remark": user.Remark, "admin": user.Admin})
	if result.Error != nil {
		c.JSON(200, gin.H{"success": false, "msg": "update db users error " + result.Error.Error()})
		return
	}

	if user.Password != "" {
		user.Password = utils.Md5plus(user.Password, setting.Setting.Md5Iteration)
		database.DB.Model(&model.Users{}).Select("password").Where("id = ?", user.Id).Updates(map[string]interface{}{"password": user.Password})
	}

	c.JSON(200, gin.H{"success": true})
	return
}

func DeleteUser(c *gin.Context) {
	var user model.Users
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(200, gin.H{"success": false, "msg": "query users error " + err.Error()})
		return
	}
	result := database.DB.Model(&model.Users{}).Where("username = ?", user.Username).Delete(user)
	if result.Error != nil {
		c.JSON(200, gin.H{"success": false, "msg": "Delete Error:" + result.Error.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
	return
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
