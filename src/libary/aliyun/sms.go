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

package aliyun

import (
	"dbmcloud/setting"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */

func _createClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

func SendSms(phoneList, TemplateParam string) (_err error) {
	client, _err := _createClient(tea.String(setting.Setting.AccessKeyId), tea.String(setting.Setting.AccessKeySecret))
	if _err != nil {
		return _err
	}
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:     tea.String(setting.Setting.SmsSignName),
		TemplateCode: tea.String(setting.Setting.SmsTemplateCode),
		PhoneNumbers: tea.String(phoneList),
		//TemplateParam: tea.String("{\"entity\":\"MySQL-10.129.100.101:3306\",\"title\":\"[告警][QPS过高]\",\"rule\":\"qps(101)>100\",\"time\":\"2022-03-22 12:00:11\"}"),
		TemplateParam: tea.String(TemplateParam),
	}
	_, _err = client.SendSms(sendSmsRequest)
	if _err != nil {
		return _err
	}
	return _err
}

/*
func main() {
	phoneList := "15216607660"
	TemplateParam := "{\"entity\":\"MySQL-10.129.100.101:3306\",\"title\":\"[告警][QPS过高]\",\"rule\":\"qps(101)>100\",\"time\":\"2022-03-22 12:00:11\"}"
	err := SendSms(phoneList, TemplateParam)
	fmt.Println(err)
}
*/
