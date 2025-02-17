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

package proxyCache

import (
	"dbmcloud/src/libary/cache"
	"dbmcloud/src/libary/conf"
	"dbmcloud/src/libary/http"
	"math/rand"
	"strings"
	"time"
)

func Send(Data interface{}, sendType string) (string, error) {
	proxyCluster, err := cache.Get("proxyCache")
	if err != nil || proxyCluster == "" {
		proxyCluster = conf.Option["proxy"]
		if err := cache.Set("proxyCache", proxyCluster, 300); err != nil {
			return "", err
		}
	}
	proxyList := strings.Split(proxyCluster, ";")
	rand.Seed(time.Now().Unix())
	proxyCount := len(proxyList)
	randNumber := rand.Intn(proxyCount)
	proxy := proxyList[randNumber]

	var proxyUrl string
	if sendType == "event" {
		proxyUrl = proxy + "/proxy/event"
	}
	if sendType == "sql" {
		proxyUrl = proxy + "/proxy/sql"
	}
	_, err = http.Post(proxyUrl, Data)

	if err != nil {
		if proxyCount == 1 {
			return proxyUrl, err
		}
		switch randNumber {
		case 0:
			proxy = proxyList[randNumber+1]
		case proxyCount - 1:
			proxy = proxyList[randNumber-1]
		default:
			proxy = proxyList[randNumber+1]
		}

		if sendType == "event" {
			proxyUrl = proxy + "/proxy/event"
		}
		if sendType == "sql" {
			proxyUrl = proxy + "/proxy/sql"
		}
		_, err = http.Post(proxyUrl, Data)
		if err != nil {
			return proxyUrl, err
		}
		return proxyUrl, nil

	}

	return proxyUrl, nil

}
