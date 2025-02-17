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

package mail

import (
	"dbmcloud/setting"
	"strconv"

	"gopkg.in/gomail.v2"
)

func Send(mailTo []string, subject, body string) error {
	mailConn := map[string]string{
		"user": setting.Setting.MailUser,
		"pass": setting.Setting.MailPass,
		"host": setting.Setting.MailHost,
		"port": setting.Setting.MailPort,
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(mailConn["user"], setting.Setting.MailFrom)) //这种方式可以添加别名，即“XX官方”
	m.SetHeader("To", mailTo...)                                                     //发送给多个用户
	m.SetHeader("Subject", subject)                                                  //设置邮件主题
	m.SetBody("text/html", body)                                                     //设置邮件正文
	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])
	err := d.DialAndSend(m)
	return err
}
