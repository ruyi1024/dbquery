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

package logger

import (
	"dbmcloud/src/libary/conf"

	go_logger "github.com/phachon/go-logger"
)

var logger *go_logger.Logger

func InitLog() *go_logger.Logger {
	var Logger = go_logger.NewLogger()
	Logger.Detach("console")
	consoleConfig := &go_logger.ConsoleConfig{
		Color:      true,
		JsonFormat: false,
		Format:     "",
	}
	Logger.Attach("console", go_logger.LOGGER_LEVEL_ERROR, consoleConfig)

	fileConfig := &go_logger.FileConfig{
		Filename:   conf.Option["log"],
		MaxSize:    1024 * 1024, //KB
		MaxLine:    10000,
		DateSlice:  "d",
		JsonFormat: false,
		Format:     "",
	}
	if conf.Option["debug"] == "1" {
		Logger.Attach("file", go_logger.LOGGER_LEVEL_DEBUG, fileConfig)
	} else {
		Logger.Attach("file", go_logger.LOGGER_LEVEL_INFO, fileConfig)
	}
	return Logger
}

func NewLog(file string, debug int) *go_logger.Logger {
	var Logger = go_logger.NewLogger()
	Logger.Detach("console")
	consoleConfig := &go_logger.ConsoleConfig{
		Color:      true,
		JsonFormat: false,
		Format:     "",
	}
	Logger.Attach("console", go_logger.LOGGER_LEVEL_ERROR, consoleConfig)

	fileConfig := &go_logger.FileConfig{
		Filename:   file,
		MaxSize:    1024 * 1024, //KB
		MaxLine:    10000,
		DateSlice:  "d",
		JsonFormat: false,
		Format:     "",
	}
	if debug == 1 {
		Logger.Attach("file", go_logger.LOGGER_LEVEL_DEBUG, fileConfig)
	} else {
		Logger.Attach("file", go_logger.LOGGER_LEVEL_INFO, fileConfig)
	}
	return Logger
}
