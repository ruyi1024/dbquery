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

package utils

import (
	"crypto/md5"
	"fmt"
	"io"
)

const SALT = "___LePus_2021_02_26 && :) ___Happy_ZhongQiu___"

// md5加盐
func Md5plus(text string, cost int) string {
	for i := 0; i < cost; i++ {
		if i%10 == 0 {
			text += SALT
		}
		hash := md5.New()
		_, _ = io.WriteString(hash, text)
		text = fmt.Sprintf("%x", hash.Sum(nil))
	}
	return text
}
