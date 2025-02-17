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

package middleware

import (
	"dbmcloud/log"
	"dbmcloud/setting"
	"dbmcloud/src/utils"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	TokenExpired error = errors.New("未登录，或登录已过期")
	TokenInvalid error = errors.New("登录错误，登录用户不存在，请尝试重新登录。")
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("debug 11-->" + fmt.Sprintf("%#v", c.Request))
		if find := strings.Contains(c.Request.RequestURI, "static/"); find {
			//fmt.Println("find the character.")
			c.Next()
			return
		}
		if c.Request.RequestURI != "/api/v1/login/account" && c.Request.RequestURI != "/" && c.Request.RequestURI != "/logo.png" && c.Request.RequestURI != "/user/login" {
			v, err := c.Cookie(setting.Setting.TokenName)
			log.Info("debug --->" + fmt.Sprintf("%#v", v))
			if err != nil || v == "" {
				//未登录跳转到登录页面时，currentUser返回false会弹出一个提示框，新增以下逻辑返回true消除提示框
				if c.Request.RequestURI == "/api/v1/currentUser" {
					c.JSON(200, gin.H{"success": true, "errorMsg": "", "errorCode": 0, "data": ""})
					c.Abort()
				}
				log.Error("get token error. ", zap.Error(err))
				c.JSON(401, gin.H{"success": false, "code": 302, "errorMsg": TokenExpired.Error()})
				c.Abort()
			} else {
				user, err := utils.TokenVerify(v)
				log.Info("user info -> " + fmt.Sprintf("%v", user))
				if err != nil || user.Username == "" {
					log.Error("verify token error: ", zap.Error(err))
					//未登录跳转到登录页面时，currentUser返回false会弹出一个提示框，新增以下逻辑返回true消除提示框
					if c.Request.RequestURI == "/api/v1/currentUser" {
						c.JSON(200, gin.H{"success": true, "errorMsg": "", "errorCode": 0, "data": ""})
						c.Abort()
					}
					c.JSON(401, gin.H{"success": false, "code": 302, "errorMsg": TokenInvalid.Error()})
					c.Abort()
				}
				c.Set("loginUser", user)
				c.Set("tokenKey", v)
				c.Set("userId", user.Id)
				c.Set("username", user.Username)
				c.Set("admin", user.Admin)
			}
		}
		// before request
		c.Next()
		// after request
	}
}
