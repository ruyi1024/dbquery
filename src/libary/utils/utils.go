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
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GetCurrentTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func StrToInt(str string) int {
	nonFractionalPart := strings.Split(str, ".")
	result, _ := strconv.Atoi(nonFractionalPart[0])
	return result
}

func StrToInt64(str string) int64 {
	result, _ := strconv.ParseInt(str, 10, 64)
	return result
}

func StrToFloat(str string) float32 {
	result, _ := strconv.ParseFloat(str, 32)
	return float32(result)
}

func StrToFloat64(str string) float64 {
	result, _ := strconv.ParseFloat(str, 64)
	return result
}

func FloatToStr(f float64) string {
	result := strconv.FormatFloat(f, 'f', 2, 64)
	return result
}

func FormatFloat64(f float64) float64 {
	result, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", f), 64) //保留2位小数，四舍五入
	return result
}

func TimeStrToTimestamp(timeStr string, flag int) int64 {
	var t int64
	loc, _ := time.LoadLocation("Local")
	if flag == 1 {
		t1, _ := time.ParseInLocation("2006.01.02 15:04:05", timeStr, loc)
		t = t1.Unix()
	} else if flag == 2 {
		t1, _ := time.ParseInLocation("2006-01-02 15:04", timeStr, loc)
		t = t1.Unix()
	} else if flag == 3 {
		t1, _ := time.ParseInLocation("2006-01-02", timeStr, loc)
		t = t1.Unix()
	} else if flag == 4 {
		t1, _ := time.ParseInLocation("2006.01.02", timeStr, loc)
		t = t1.Unix()
	} else {
		t1, _ := time.ParseInLocation("2006-01-02 15:04:05", timeStr, loc)
		t = t1.Unix()
	}
	return t
}

func FormatDateTime(timeStr string) (newTimeStr string) {
	// 给定的时间字符串
	timeStr = "2024-06-06T14:46:02+08:00"
	// 解析时间字符串，这里使用了ISO 8601格式
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return "------"
	}
	// 使用Format方法输出不带时区的格式，例如 "2006-01-02 15:04:05"
	// 注意这里的格式字符串与RFC3339不同，去掉了时区部分
	newTimeStr = t.Format("2006-01-02 15:04:05")
	return newTimeStr
}

//分割数组，根据传入的数组和分割大小，将数组分割为大小等于指定大小的多个数组，如果不够分，则最后一个数组元素小于其他数组
//数组：[1, 2, 3, 4, 5, 6, 7, 8, 9]，正整数：2
//期望结果: [[1, 2], [3, 4], [5, 6], [7, 8], [9]]

func SplitArray(arr []int, num int64) [][]int {
	max := int64(len(arr))
	//判断数组大小是否小于等于指定分割大小的值，是则把原数组放入二维数组返回
	if max <= num {
		return [][]int{arr}
	}
	//获取应该数组分割为多少份
	var quantity int64
	if max%num == 0 {
		quantity = max / num
	} else {
		quantity = (max / num) + 1
	}
	//声明分割好的二维数组
	var segments = make([][]int, 0)
	//声明分割数组的截止下标
	var start, end, i int64
	for i = 1; i <= quantity; i++ {
		end = i * num
		if i != quantity {
			segments = append(segments, arr[start:end])
		} else {
			segments = append(segments, arr[start:])
		}
		start = i * num
	}
	return segments
}

func MapToStr(data []map[string]interface{}) string {
	// 序列化为 JSON 字符串
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err) // 处理错误
	}
	// 将 JSON 字符串转换为可打印的字符串
	str := string(jsonData)
	return str
}

func SplitArrayMap(arr []map[string]interface{}, num int64) [][]map[string]interface{} {
	max := int64(len(arr))
	//判断数组大小是否小于等于指定分割大小的值，是则把原数组放入二维数组返回
	if max <= num {
		return [][]map[string]interface{}{arr}
	}
	//获取应该数组分割为多少份
	var quantity int64
	if max%num == 0 {
		quantity = max / num
	} else {
		quantity = (max / num) + 1
	}
	//声明分割好的二维数组
	var segments = make([][]map[string]interface{}, 0)
	//声明分割数组的截止下标
	var start, end, i int64
	for i = 1; i <= quantity; i++ {
		end = i * num
		if i != quantity {
			segments = append(segments, arr[start:end])
		} else {
			segments = append(segments, arr[start:])
		}
		start = i * num
	}
	return segments
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	} else {
		n = n + len(start) // 增加了else，不加的会把start带上
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}
