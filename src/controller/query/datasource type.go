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

package query

import (
	"dbmcloud/src/database"
	"fmt"
	"net/http"
	_ "reflect"

	"github.com/gin-gonic/gin"
)

func DataSourceTypeList(c *gin.Context) {
	admin, _ := c.Get("admin")
	username, _ := c.Get("username")
	method := c.Request.Method
	if method == "GET" {
		var searchCondition string
		if admin != true {
			searchCondition = fmt.Sprintf(" and name in (select datasource_type from privileges where username='%s')", username)
		}
		sql := fmt.Sprintf("select id,name from datasource_type where enable=1 %s order by sort asc", searchCondition)
		dataList, _ := database.QueryAll(sql)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"msg":     "OK",
			"data":    dataList,
			"total":   len(dataList),
		})
		return
	}
}
