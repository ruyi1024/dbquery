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

package setting

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type setting struct {
	Log        `yam:"log"`
	DataSource `yaml:"dataSource"`
	Notice     `yaml:"notice"`
	Decrypt    `yaml:"decrypt"`
	Token      `yaml:"token"`
}

type Log struct {
	Path  string `yaml:"path"`
	Level string `yaml:"level"`
	Debug bool   `yaml:"debug"`
}

type DataSource struct {
	Host               string `yaml:"host"`
	Port               string `yaml:"port"`
	User               string `yaml:"user"`
	Password           string `yaml:"password"`
	Database           string `yaml:"database"`
	RedisHost          string `yaml:"redisHost"`
	RedisPort          string `yaml:"redisPort"`
	RedisPassword      string `yaml:"redisPassword"`
	ClickhouseHost     string `yaml:"clickhouseHost"`
	ClickhousePort     string `yaml:"clickhousePort"`
	ClickhouseUser     string `yaml:"clickhouseUser"`
	ClickhousePassword string `yaml:"clickhousePassword"`
	ClickhouseDatabase string `yaml:"clickhouseDatabase"`
	NsqServer          string `yaml:"nsqServer"`
}

type Notice struct {
	MailHost             string `yaml:"mailHost"`
	MailPort             string `yaml:"mailPort"`
	MailUser             string `yaml:"mailUser"`
	MailPass             string `yaml:"mailPass"`
	MailFrom             string `yaml:"mailFrom"`
	AccessKeyId          string `yaml:"accessKeyId"`
	AccessKeySecret      string `yaml:"accessKeySecret"`
	SmsSignName          string `yaml:"smsSignName"`
	SmsTemplateCode      string `yaml:"smsTemplateCode"`
	PhoneTemplateCode    string `yaml:"phoneTemplateCode"`
	PhonePlayTimes       string `yaml:"phonePlayTimes"`
	WechatAppId          string `yaml:"wechatAppId"`
	WechatAppSecret      string `yaml:"wechatAppSecret"`
	WechatSendTemplateId string `yaml:"wechatSendTemplateId"`
}

type Decrypt struct {
	SignKey      string `yaml:"signKey"`
	DbPassKey    string `yaml:"dbPassKey"`
	Md5Iteration int
}

type Token struct {
	TokenKey     string `yaml:"key"`
	TokenName    string `yaml:"name"`
	Expired      string `yaml:"expired"`
	TokenExpired int64
}

var Setting = new(setting)

func InitSetting(path string) (err error) {
	if f, err := os.Open(path); err == nil {
		defer func() {
			_ = f.Close()
		}()
		c, err := ioutil.ReadAll(f)
		if err != nil {
			return errors.Wrap(err, "init setting read file")
		}
		err = yaml.Unmarshal(c, &Setting)
		if err != nil {
			return errors.Wrap(err, "init setting unmarshal data")
		}
		//fmt.Println(fmt.Sprintf("config:%#v", Setting))
	}
	Setting.Md5Iteration = 1500
	// token expired
	if strings.HasSuffix(strings.ToLower(Setting.Expired), "h") {
		s := strings.Replace(Setting.Expired, "h", "", -1)
		h, _ := strconv.ParseInt(s, 10, 64)
		Setting.TokenExpired = h * 60 * 60
	}

	if strings.HasSuffix(strings.ToLower(Setting.Expired), "d") {
		s := strings.Replace(Setting.Expired, "d", "", -1)
		h, _ := strconv.ParseInt(s, 10, 64)
		Setting.TokenExpired = h * 60 * 60 * 24
	}
	return nil
}

func DataSourceInfo() DataSource {
	return Setting.DataSource
}
