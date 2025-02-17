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
	"fmt"
	"strconv"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dyvmsapi20170525 "github.com/alibabacloud-go/dyvmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */

func _createClient2(accessKeyId *string, accessKeySecret *string) (_result *dyvmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dyvmsapi.aliyuncs.com")
	_result = &dyvmsapi20170525.Client{}
	_result, _err = dyvmsapi20170525.NewClient(config)
	return _result, _err
}

func CallPhone(phoneStrList, TemplateParam string) (_err error) {
	client, _err := _createClient2(tea.String(setting.Setting.AccessKeyId), tea.String(setting.Setting.AccessKeySecret))
	if _err != nil {
		fmt.Println(_err)
		return _err
	}

	// 每次只支持一个电话号码
	phoneList := strings.Split(phoneStrList, ",")
	if len(phoneList) > 0 {
		phonePlayTimes, _ := strconv.ParseInt(setting.Setting.PhonePlayTimes, 10, 32)
		phoneTemplateCode := setting.Setting.PhoneTemplateCode
		for _, phone := range phoneList {
			singleCallByTtsRequest := &dyvmsapi20170525.SingleCallByTtsRequest{
				TtsCode:      tea.String(phoneTemplateCode),
				PlayTimes:    tea.Int32(int32(phonePlayTimes)),
				CalledNumber: tea.String(phone),
				TtsParam:     tea.String(TemplateParam),
			}
			rsp, _err := client.SingleCallByTts(singleCallByTtsRequest)
			if _err != nil {
				fmt.Println(_err)
				return _err
			} else {
				fmt.Println(rsp)
				return nil
			}
		}
	}

	return _err
}

// func main() {
// 	phoneList := "15216607660"
// 	TemplateParam := "{\"title\":\"数据库服务异常\"}"
// 	err := CallPhone(phoneList, TemplateParam)
// 	fmt.Println(err)
// }
