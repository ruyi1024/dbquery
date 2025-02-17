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

package cache

import (
	"fmt"

	"github.com/coocood/freecache"
)

var cache = freecache.NewCache(25 * 1024 * 1024) //25M Cache

func Set(key, val string, expire int) error {
	keyByte := []byte(key)
	valByte := []byte(val)
	expire = 60 // expire in 60 seconds
	err := cache.Set(keyByte, valByte, expire)
	if err != nil {
		return err
	}
	return nil
}

func Get(key string) (string, error) {
	keyByte := []byte(key)
	got, err := cache.Get(keyByte)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", got), nil

}

/*
func main() {
	err := Set("proxy", "127.0.0.1;192.168.10.1", 60)
	if err != nil {
		fmt.Println(err)
	}
	data, err := Get("proxy")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(data)
	}
}
*/
