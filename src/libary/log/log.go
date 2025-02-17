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

package log

import (
	"github.com/phachon/go-logger"
)

var Logger = go_logger.NewLogger()

func Initialize(file string) (data *go_logger.Logger) {

	Logger.Detach("console")
	consoleConfig := &go_logger.ConsoleConfig{
		Color:      true,
		JsonFormat: false,
		Format:     "",
	}
	Logger.Attach("console", go_logger.LOGGER_LEVEL_DEBUG, consoleConfig)

	fileConfig := &go_logger.FileConfig{
		Filename:   file,
		MaxSize:    1024 * 1024, //KB
		MaxLine:    0,
		DateSlice:  "d",
		JsonFormat: false,
		Format:     "",
	}
	Logger.Attach("file", go_logger.LOGGER_LEVEL_DEBUG, fileConfig)

	return Logger

}
