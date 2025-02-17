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

package conf

import (
	"flag"
	"os"

	"github.com/larspensjo/config"
)

var (
	configFile = flag.String("config", "./etc/config.ini", "General configuration file")
	Option     = make(map[string]string)
)

func init() {
	flag.Parse()
	cfg, err := config.ReadDefault(*configFile)
	if err != nil {
		panic(err)
		os.Exit(0)
	}
	if cfg.HasSection("main") {
		section, err := cfg.SectionOptions("main")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("main", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}
	if cfg.HasSection("mysql") {
		section, err := cfg.SectionOptions("mysql")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("mysql", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}
	if cfg.HasSection("mongodb") {
		section, err := cfg.SectionOptions("mongodb")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("mongodb", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}
	if cfg.HasSection("redis") {
		section, err := cfg.SectionOptions("redis")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("redis", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}
	if cfg.HasSection("kafka") {
		section, err := cfg.SectionOptions("kafka")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("kafka", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}
	if cfg.HasSection("mail") {
		section, err := cfg.SectionOptions("mail")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("mail", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}
	if cfg.HasSection("task") {
		section, err := cfg.SectionOptions("task")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("task", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}

	if cfg.HasSection("event") {
		section, err := cfg.SectionOptions("event")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("event", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}

	if cfg.HasSection("nsq") {
		section, err := cfg.SectionOptions("nsq")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("nsq", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}

	if cfg.HasSection("influxdb") {
		section, err := cfg.SectionOptions("influxdb")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("influxdb", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}

	if cfg.HasSection("clickhouse") {
		section, err := cfg.SectionOptions("clickhouse")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("clickhouse", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}

	if cfg.HasSection("doris") {
		section, err := cfg.SectionOptions("doris")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("doris", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}

	if cfg.HasSection("aliyun") {
		section, err := cfg.SectionOptions("aliyun")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("aliyun", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}

	if cfg.HasSection("wechat") {
		section, err := cfg.SectionOptions("wechat")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("wechat", v)
				if err == nil {
					Option[v] = options
				}
			}
		}
	}

}
