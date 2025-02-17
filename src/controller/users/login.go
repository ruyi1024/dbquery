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
	"bytes"
	"dbmcloud/log"
	"dbmcloud/setting"
	"dbmcloud/src/database"
	"dbmcloud/src/model"
	"dbmcloud/src/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type longin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	var data bytes.Buffer
	_, err := io.Copy(&data, c.Request.Body)
	if err != nil {
		return
	}
	res := make([]byte, data.Len())
	_, err = data.Read(res)
	if err != nil {
		c.JSON(200, gin.H{"successLogin": false, "msg": "login error. " + err.Error()})
		return
	}
	h, _ := hex.DecodeString(string(res))
	d, err := utils.AesDecrypt(h, []byte(setting.Setting.SignKey))
	if err != nil {
		c.JSON(200, gin.H{"successLogin": false, "msg": "login error. " + err.Error()})
		return
	}
	var l longin
	err = json.Unmarshal(d, &l)
	if err != nil {
		c.JSON(200, gin.H{"successLogin": false, "msg": "login error. " + err.Error()})
		return
	}
	log.Info("login", zap.String("username", l.Username))
	// 验证密码
	var user model.Users
	user.Username = l.Username
	result := database.DB.Where("username = ?", l.Username).First(&user)
	if result.Error == nil && result.RowsAffected == 1 {
		if user.Password == utils.Md5plus(l.Password, setting.Setting.Md5Iteration) {
			// set token

			k, v, err := utils.TokenCreate(user.Username, user.Id, setting.Setting.TokenExpired)
			if err != nil {
				c.JSON(200, gin.H{"successLogin": false, "msg": "create session fail. " + err.Error()})
				return
			}
			/*s := sessions.Default(c)
			s.Options()*/
			c.SetCookie(setting.Setting.TokenName, k, int(setting.Setting.TokenExpired), "/", "", false, true)
			database.DB.Create(model.Token{
				TokenKey:  k,
				Value:     v,
				CreatedAt: time.Now(),
				Expired:   time.Unix(time.Now().Unix()+setting.Setting.TokenExpired, 0),
			})
			c.JSON(200, gin.H{"successLogin": true, "msg": "登录成功!"})
			return
		} else {
			c.JSON(200, gin.H{"successLogin": false, "msg": "密码错误！"})
		}
	} else {
		c.JSON(200, gin.H{"successLogin": false, "msg": "账号不存在！"})
	}
	return
}

// Logout 登出用户
func Logout(c *gin.Context) {
	tk, _ := c.Get("tokenKey")
	if tk != nil {
		utils.RemoveToken(fmt.Sprintf("%s", tk))
	}
}

// CurrentUser 获取当前已登录用户信息
func CurrentUser(c *gin.Context) {
	// 获取有已登录用户信息
	user, _ := c.Get("loginUser")
	c.JSON(200, gin.H{"success": true, "errorMsg": "", "errorCode": 0, "data": user})
	return
}
