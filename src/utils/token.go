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

package utils

import (
	"dbmcloud/log"
	"dbmcloud/setting"
	"dbmcloud/src/database"
	"dbmcloud/src/model"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"time"
)

type token struct {
	Username     string `json:"username"`
	ChineseName  string `json:"chineseName"`
	RemoteIp     string `json:"remoteIp"`
	UserId       int64  `json:"userId"`
	CreateOn     int64  `json:"createOn"`
	TokenExpired int64  `json:"tokenExpired"`
}

// TokenCreate 创建token
// 返回 key 和 value 以及错误
func TokenCreate(username string, userId int64, expired int64) (k string, v []byte, err error) {
	var aeKey = []byte(setting.Setting.TokenKey)
	var tk token
	// tk.UserId = userId
	tk.Username = username
	// tk.CreateOn = time.Now().Unix()
	tk.TokenExpired = expired
	tb, err := json.Marshal(tk)
	v, err = AesEncrypt(tb, aeKey)
	if err != nil {
		log.Error("generate token error.", zap.Error(err))
		return
	}
	log.Info("token " + base64.StdEncoding.EncodeToString(v))
	k = base64.StdEncoding.EncodeToString(v)
	return
}

// TokenVerify token 验证, 并返回用户信息
func TokenVerify(key string) (user model.Users, err error) {
	var tk model.Token
	result := database.DB.Where("token_key = ?", key).First(&tk)
	if result.Error != nil {
		return user, result.Error
	}
	if result.RowsAffected == 1 {
		if tk.Expired.Unix() < time.Now().Unix() {
			// token 过期
			RemoveToken(key)
			result.Error = fmt.Errorf("%s", "token expired!")
			return user, result.Error
		}
		log.Info("debug TokenVerify token value -> " + tk.TokenKey)
		t, err := tokenDecrypt(tk.Value)
		if err != nil {
			return user, err
		}
		log.Info("debug TokenVerify decrypt token value -> "+fmt.Sprintf("%v", t), zap.String("username", t.Username))
		database.DB.Where("username = ?", t.Username).First(&user)
		log.Info("debug TokenVerify query user is -> " + fmt.Sprintf("%v", user))
		return user, nil
	}
	return
}

// RemoveToken 移除指定token
func RemoveToken(key string) {
	database.DB.Where("token_key = ?", key).Delete(&model.Token{})
}

// 解密
func tokenDecrypt(v []byte) (tk token, err error) {
	var aesKey = []byte(setting.Setting.TokenKey)
	val, err := AesDecrypt(v, aesKey)
	if err != nil {
		log.Error("decrypt token error.", zap.Error(err))
		return
	}
	_ = json.Unmarshal(val, &tk)
	return
}
