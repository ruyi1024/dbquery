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

package task

import (
	"dbmcloud/src/database"
	"dbmcloud/src/model"
	"time"

	"github.com/robfig/cron/v3"
)

func init() {
	go revokePrivileage()
}

func revokePrivileage() {
	/*
		time.Sleep(time.Second * time.Duration(rand.Intn(60)))
		timer := time.NewTicker(180 * time.Second)
		defer timer.Stop()
		for {
			<-timer.C
			database.DB.Model(model.TaskHeartbeat{}).Where("heartbeat_key='revoke_privileage'").Updates(map[string]interface{}{"heartbeat_time": time.Now().Format("2006-01-02 15:04:05.999")})
			database.DB.Delete(model.Privilege{}, "expire_date <= ?", time.Now().Format("2006-01-02"))
			database.DB.Model(model.TaskHeartbeat{}).Where("heartbeat_key='revoke_privileage'").Updates(map[string]interface{}{"heartbeat_end_time": time.Now().Format("2006-01-02 15:04:05.999")})
		}
	*/
	time.Sleep(time.Second * time.Duration(30))
	var db = database.DB
	var record model.TaskOption
	db.Select("crontab").Where("task_key=?", "revoke_privileage").Take(&record)
	c := cron.New()
	c.AddFunc(record.Crontab, func() {
		db.Select("enable").Where("task_key=?", "revoke_privileage").Take(&record)
		if record.Enable == 1 {
			db.Model(model.TaskHeartbeat{}).Where("heartbeat_key='revoke_privileage'").Updates(map[string]interface{}{"heartbeat_time": time.Now().Format("2006-01-02 15:04:05.999")})
			database.DB.Delete(model.Token{}, "expired <= ?", time.Now().Format("2006-01-02 15:04:05"))
			db.Model(model.TaskHeartbeat{}).Where("heartbeat_key='revoke_privileage'").Updates(map[string]interface{}{"heartbeat_end_time": time.Now().Format("2006-01-02 15:04:05.999")})
		}
	})
	c.Start()
}
