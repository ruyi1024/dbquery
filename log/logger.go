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
	"fmt"
	"go.uber.org/zap"
	"time"
)

func init() {
	InitLogs()
}

var Logger *zap.Logger

type Log interface {
	Info(msg string, field ...zap.Field)
	Warn(msg string, field ...zap.Field)
	Error(msg string, field ...zap.Field)
	Debug(msg string, field ...zap.Field)
	DPanic(msg string, field ...zap.Field)
	Panic(msg string, field ...zap.Field)
	Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error)
}

func Info(msg string, field ...zap.Field) {
	Logger.Info(msg, field...)
}

func Warn(msg string, field ...zap.Field) {
	Logger.Warn(msg, field...)
}

func Error(msg string, field ...zap.Field) {
	Logger.Error(msg, field...)
}

func Debug(msg string, field ...zap.Field) {
	Logger.Debug(msg, field...)
}

func DPanic(msg string, field ...zap.Field) {
	Logger.DPanic(msg, field...)
}

func Panic(msg string, field ...zap.Field) {
	Logger.Panic(msg, field...)
}

func Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	Logger.Debug("trace", zap.Any("ctx", fmt.Sprintf("%#v", ctx)), zap.Time("begin", begin), zap.Error(err))
}
