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
	"dbmcloud/setting"
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func formatEncodeTime(t time.Time, en zapcore.PrimitiveArrayEncoder) {
	en.AppendString(fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d.%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond()))
}

func level(l string) zapcore.Level {
	switch strings.ToLower(l) {
	case "info":
		return zapcore.InfoLevel
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	default:
	}
	return zapcore.InfoLevel
}

func InitLogs() {
	hook := lumberjack.Logger{
		Filename:   setting.Setting.Path,
		MaxSize:    10,
		MaxBackups: 7,
		MaxAge:     10,
		Compress:   false,
	}

	ec := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeTime:     formatEncodeTime,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// set log level
	l := zap.NewAtomicLevel()
	l.SetLevel(level(setting.Setting.Level))

	c := zapcore.NewCore(zapcore.NewConsoleEncoder(ec), zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), l)
	if setting.Setting.Debug {
		Logger = zap.New(c, zap.AddCaller(), zap.Development())
	} else {
		Logger = zap.New(c)
	}
	HandleLogger(Logger)
	Logger.Info("logger load finish")
}
