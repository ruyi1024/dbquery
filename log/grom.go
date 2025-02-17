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
	"context"
	"database/sql/driver"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	glogger "gorm.io/gorm/logger"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"
)

var (
	// GormLoggerName gorm logger 名称
	GormLoggerName = "gorm"
	// GormLoggerCallerSkip caller skip
	GormLoggerCallerSkip = 3
)

type GormLogger struct {
	// 日志级别
	logLevel zapcore.Level
	// 指定慢查询时间
	slowThreshold time.Duration
	// Trace 方法打印日志是使用的日志 level
	traceWithLevel zapcore.Level
}

var gormLogLevelMap = map[glogger.LogLevel]zapcore.Level{
	glogger.Info:  zap.InfoLevel,
	glogger.Warn:  zap.WarnLevel,
	glogger.Error: zap.ErrorLevel,
}

func (g GormLogger) LogMode(gormLogLevel glogger.LogLevel) glogger.Interface {
	zl, exists := gormLogLevelMap[gormLogLevel]
	if !exists {
		zl = zap.DebugLevel
	}
	newlogger := g
	newlogger.logLevel = zl
	return &newlogger
}

// Info 实现 gorm logger 接口方法
func (g GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel <= zap.InfoLevel {
		//g.CtxLogger(ctx).Sugar().Infof(msg, data...)
	}
}

// Warn 实现 gorm logger 接口方法
func (g GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel <= zap.WarnLevel {
		//g.CtxLogger(ctx).Sugar().Warnf(msg, data...)
	}
}

// Error 实现 gorm logger 接口方法
func (g GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if g.logLevel <= zap.ErrorLevel {
		//g.CtxLogger(ctx).Sugar().Errorf(msg, data...)
	}
}

// Trace 实现 gorm logger 接口方法
func (g GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	now := time.Now()
	latency := now.Sub(begin).Seconds()
	sql, rows := fc()
	sql = removeDuplicateWhitespace(sql, true)
	//logger := g.CtxLogger(ctx)
	switch {
	case err != nil:
		Error("sql: "+sql, zap.Float64("latency", latency), zap.Int64("rows", rows), zap.String("error", err.Error()))
	case g.slowThreshold != 0 && latency > g.slowThreshold.Seconds():
		Warn("sql: "+sql, zap.Float64("latency", latency), zap.Int64("rows", rows), zap.Float64("threshold", g.slowThreshold.Seconds()))
	default:
		log := Debug
		if g.traceWithLevel == zap.InfoLevel {
			log = Info
		} else if g.traceWithLevel == zap.WarnLevel {
			log = Warn
		} else if g.traceWithLevel == zap.ErrorLevel {
			log = Error
		}
		log("sql: "+sql, zap.Float64("latency", latency), zap.Int64("rows", rows))
	}
}

func NewGormLogger(logLevel zapcore.Level, traceWithLevel zapcore.Level, slowThreshold time.Duration) GormLogger {
	return GormLogger{
		logLevel:       logLevel,
		slowThreshold:  slowThreshold,
		traceWithLevel: traceWithLevel,
	}
}

func removeDuplicateWhitespace(s string, trim bool) string {
	ws, err := regexp.Compile(`\s+`)
	if err != nil {
		return s
	}
	s = ws.ReplaceAllString(s, " ")
	if trim {
		s = strings.TrimSpace(s)
	}
	return s
}

func GormPrint(values ...interface{}) {
	if len(values) < 2 {
		Info(fmt.Sprint(values...))
	}

	if len(values) == 2 {
		Info(fmt.Sprintf("%v", values[1]))
	}

	level := values[0]
	if level == "log" {
		Debug("log", zap.Any("gorm", fmt.Sprint(values[2:]...)))
	}

	if level == "sql" {
		Debug(
			"sql query",
			zap.String("module", "gorm"),
			zap.String("type", "sql"),
			zap.Any("src", values[1]),
			zap.Any("duration", values[2]),
			zap.Any("sql", formatSQL(values[3].(string), values[4].([]interface{}))),
			zap.Any("rows_returned", values[5]),
		)
	}

}

func formatSQL(sql string, values []interface{}) string {
	size := len(values)

	replacements := make([]string, size*2)

	var indexFunc func(int) string
	if strings.Contains(sql, "$1") {
		indexFunc = formatNumbered
	} else {
		indexFunc = formatQuestioned
	}

	for i := size - 1; i >= 0; i-- {
		replacements[(size-i-1)*2] = indexFunc(i)
		replacements[(size-i-1)*2+1] = formatValue(values[i])
	}

	r := strings.NewReplacer(replacements...)
	return r.Replace(sql)
}

func formatNumbered(index int) string {
	return fmt.Sprintf("$%d", index+1)
}

func formatQuestioned(index int) string {
	return "?"
}

func formatValue(value interface{}) string {
	indirectValue := reflect.Indirect(reflect.ValueOf(value))
	if !indirectValue.IsValid() {
		return "NULL"
	}

	value = indirectValue.Interface()

	switch v := value.(type) {
	case time.Time:
		return fmt.Sprintf("'%v'", v.Format("2006-01-02 15:04:05"))
	case []byte:
		s := string(v)
		if isPrintable(s) {
			return redactLong(fmt.Sprintf("'%s'", s))
		}
		return "'<binary>'"
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case driver.Valuer:
		if dv, err := v.Value(); err == nil && dv != nil {
			return formatValue(dv)
		}
		return "NULL"
	default:
		return redactLong(fmt.Sprintf("'%v'", value))
	}
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

func redactLong(s string) string {
	if len(s) > maxLen {
		return "'<redacted>'"
	}
	return s
}

const maxLen = 255
