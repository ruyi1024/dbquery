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

package router

import (
	"dbmcloud/src/controller/audit"
	"dbmcloud/src/controller/datasource"
	"dbmcloud/src/controller/favorite"
	"dbmcloud/src/controller/meta"
	"dbmcloud/src/controller/query"
	"dbmcloud/src/controller/task"
	"dbmcloud/src/controller/users"
	"dbmcloud/src/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.New()
	// session
	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("lepus-v2", store))
	r.Use(middleware.Auth())

	v1 := r.Group("api/v1")
	{
		v1.GET("/currentUser", users.CurrentUser)
		v1.POST("/login/account", users.Login)
		v1.GET("/login/outLogin", users.Logout)
		v1.GET("/users/manager/lists", users.GetUsers)
		v1.POST("/users/manager/lists", users.PostUser)
		v1.PUT("/users/manager/lists", users.PutUser)
		v1.DELETE("/users/manager/lists", users.DeleteUser)

		v1.GET("/datasource/list", datasource.List)
		v1.POST("/datasource/list", datasource.List)
		v1.PUT("/datasource/list", datasource.List)
		v1.DELETE("/datasource/list", datasource.List)
		v1.POST("/datasource/check", datasource.Check)

		v1.GET("/datasource_type/list", datasource.TypeList)
		v1.POST("/datasource_type/list", datasource.TypeList)
		v1.PUT("/datasource_type/list", datasource.TypeList)
		v1.DELETE("/datasource_type/list", datasource.TypeList)

		v1.GET("/datasource_idc/list", datasource.IdcList)
		v1.POST("/datasource_idc/list", datasource.IdcList)
		v1.PUT("/datasource_idc/list", datasource.IdcList)
		v1.DELETE("/datasource_idc/list", datasource.IdcList)

		v1.GET("/datasource_env/list", datasource.EnvList)
		v1.POST("/datasource_env/list", datasource.EnvList)
		v1.PUT("/datasource_env/list", datasource.EnvList)
		v1.DELETE("/datasource_env/list", datasource.EnvList)

		v1.GET("/task/option", task.OptionList)
		v1.POST("/task/option", task.OptionList)
		v1.PUT("/task/option", task.OptionList)
		v1.DELETE("/task/option", task.OptionList)
		v1.GET("/task/heartbeat", task.HeartbeatList)

		v1.GET("/query/datasource_type", query.DataSourceTypeList)
		v1.GET("/query/datasource", query.DataSourceList)
		v1.GET("/query/database", query.DatabaseList)
		v1.GET("/query/table", query.TableList)
		v1.POST("/query/doQuery", query.DoQuery)
		v1.POST("/query/writeLog", query.DoWriteLog)

		v1.GET("/favorite/list", favorite.List)
		v1.POST("/favorite/list", favorite.List)
		v1.PUT("/favorite/list", favorite.List)
		v1.DELETE("/favorite/list", favorite.List)

		v1.GET("/meta/instance/list", meta.InstanceList)
		v1.GET("/meta/database/list", meta.DatabaseList)
		v1.GET("/meta/table/list", meta.TableList)
		v1.GET("/meta/column/list", meta.ColumnList)
		v1.GET("/meta/dashboard/info", meta.DashboardInfo)

		v1.GET("/meta/env/list", meta.MetaEnvList)
		v1.POST("/meta/env/list", meta.MetaEnvList)
		v1.PUT("/meta/env/list", meta.MetaEnvList)
		v1.DELETE("/meta/env/list", meta.MetaEnvList)
		v1.GET("/audit/query_log", audit.GetQueryLog)
		v1.GET("/task/list", task.TaskList)

	}

	return r
}
