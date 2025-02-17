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

package main

import (
	"dbmcloud/log"
	"dbmcloud/router"
	"dbmcloud/setting"
	"dbmcloud/src/database"
	_ "dbmcloud/src/task" // run task
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

const Version = "6.0"

func help() {
	h := `
Usage: [OPTION] ...
Used to perform some operation commands on the lepus, if there is no command, start directly.

Mandatory arguments to long options are mandatory for short options too.
  -h        dispaly help info.
  -v        display this version and exit
  -c        specify the configuration file path, the default is './setting.yml'
  -l        display local machine id and license info.
`
	fmt.Println(h)
}

func init() {
	path := "./setting.yml"
	args := os.Args
	if len(args) >= 2 {
		switch args[1] {
		case "-h":
			help()
			os.Exit(0)
		case "-v":
			fmt.Println(Version)
			os.Exit(0)
		case "-c":
			path = args[2]
		default:
			help()
			os.Exit(0)
		}
	}

	// init setting
	err := setting.InitSetting(path)
	if err != nil {
		fmt.Println(err)
	}

	// init log
	log.InitLogs()

	// init database
	database.DB = database.InitDb()
	database.SQL = database.InitConnect()
}

//go:embed index.html
var indexHtml embed.FS

//go:embed static/*
var staticAsset embed.FS

func main() {
	r := router.Router()
	r.Use(log.HandleLogger(log.Logger), log.HandleRecovery(log.Logger, true))

	r.SetHTMLTemplate(template.Must(template.New("").ParseFS(indexHtml, "index.html")))
	r.StaticFS("/public/", http.FS(staticAsset))

	//解决前后台打包后找不到logo和avarar路径的问题
	//r.StaticFile("/logo.png", "./static/logo.png")
	//r.StaticFile("/avatar.jpg", "./static/avatar.jpg")
	r.GET("/logo.png", func(c *gin.Context) {
		c.Request.URL.Path = "/public/static/logo.png"
		r.HandleContext(c)
	})
	r.GET("/avatar.jpg", func(c *gin.Context) {
		c.Request.URL.Path = "/public/static/avatar.jpg"
		r.HandleContext(c)
	})

	r.GET("/", func(c *gin.Context) {
		//c.String(200, "Lepus Home Test!")
		c.HTML(200, "index.html", "")
	})
	_ = r.Run(":8088")
}
