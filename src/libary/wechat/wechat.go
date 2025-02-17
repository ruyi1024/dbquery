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

package wechat

import (
	"dbmcloud/setting"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

var (
	APPID          = setting.Setting.WechatAppId
	APPSECRET      = setting.Setting.WechatAppSecret
	SentTemplateID = setting.Setting.WechatSendTemplateId //每监控告警通知 模板ID
)

type token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type sentence struct {
	Content     string `json:"content"`
	Note        string `json:"note"`
	Translation string `json:"translation"`
}

// 获取微信accessToken
func getAccessToken() (string, error) {
	APPID = setting.Setting.WechatAppId
	APPSECRET = setting.Setting.WechatAppSecret
	fmt.Println(APPID)
	fmt.Println(setting.Setting.WechatAppId)
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%v&secret=%v", APPID, APPSECRET)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("获取微信token失败", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("微信token读取失败", err)
		return "", err
	}

	token := token{}
	err = json.Unmarshal(body, &token)
	if err != nil {
		fmt.Println("微信token解析json失败", err)
		return "", err
	}
	fmt.Println(token)
	return token.AccessToken, nil
}

//获取关注者列表

func getFollowerList(accessToken string) []gjson.Result {
	url := "https://api.weixin.qq.com/cgi-bin/user/get?access_token=" + accessToken + "&next_openid="
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("获取关注列表失败", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取内容失败", err)
		return nil
	}
	followerList := gjson.Get(string(body), "data.openid").Array()
	return followerList
}

// 发送模板消息
func templatePost(accessToken string, reqData string, url string, templateId string, openId string) (err error) {
	sendUrl := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + accessToken

	reqBody := "{\"touser\":\"" + openId + "\", \"template_id\":\"" + templateId + "\", \"url\":\"" + url + "\", \"data\": " + reqData + "}"

	resp, err := http.Post(sendUrl,
		"application/x-www-form-urlencoded",
		strings.NewReader(string(reqBody)))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	/*
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println(body)
	*/
	return nil
}

func Send(userStrList, templateData string) error {
	accessToken, err := getAccessToken()
	SentTemplateID = setting.Setting.WechatSendTemplateId //每监控告警通知 模板ID
	if err != nil {
		return err
	}
	if accessToken == "" {
		return errors.New("access token empty")
	}
	userList := strings.Split(userStrList, ",")
	if len(userList) > 0 {
		for _, userId := range userList {
			openId := userId
			err := templatePost(accessToken, templateData, "", SentTemplateID, openId)
			if err != nil {
				return err
			}
		}
	}
	return nil
	/*
		followerList := getFollowerList(accessToken)
		if followerList == nil {
			return
		}
	*/

	/*
		var city string
		for _, v := range followerList {
			go sendAlarm(accessToken, city, v.Str)
		}
	*/

}

/*
func main() {
	userStrList := "o0OjWwQTikvoazf8-OKHaxDMAV6c,o0OjWwT3mAUEJwWMm3ZwI_qhRsks"
	templateData := "{\"first\":{\"value\":\"[MySQL]数据库连接数异常\", \"color\":\"#0000CD\"},\"keyword1\":{\"value\":\"2022-03-25 18:55:48\", \"color\":\"#0000CD\"},\"keyword2\":{\"value\":\"192.168.10.100:3306\", \"color\":\"#0000CD\"},\"keyword3\":{\"value\":\"警告\", \"color\":\"#CC6633\"},\"keyword4\":{\"value\":\"ThreadConnected(381)>100\", \"color\":\"#0000CD\"},\"remark\":{\"value\":\"Lepus通知您尽快关注和处理。\", \"color\":\"#0000CD\"}}"
	//fmt.Println(reqdata)
	Send(userStrList, templateData)
}
*/
