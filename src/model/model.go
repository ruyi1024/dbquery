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

package model

import (
	"time"
)

type Users struct {
	Id          int64     `gorm:"primaryKey" json:"id"`
	Username    string    `gorm:"size:30;uniqueIndex" json:"username"`
	Password    string    `gorm:"size:200" json:"password"`
	ChineseName string    `gorm:"size:50" json:"chineseName"`
	Admin       bool      `gorm:"default:false" json:"admin"`
	Remark      string    `gorm:"size:200" json:"remark"`
	CreatedAt   time.Time `json:"createAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Token struct {
	TokenKey  string    `gorm:"primaryKey;size:180"`
	Value     []byte    `gorm:"type:bytes;size:1000"`
	CreatedAt time.Time `json:"createAt"`
	Expired   time.Time `json:"expired"`
}

type DatasourceType struct {
	Id          int       `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"size:30;uniqueIndex" json:"name"`
	Description string    `gorm:"size:100" json:"description"`
	Sort        int8      `gorm:"default:1" json:"sort"`
	Enable      int8      `gorm:"default:1" json:"enable"`
	CreatedAt   time.Time `gorm:"column:gmt_created" json:"gmt_created"`
	UpdatedAt   time.Time `gorm:"column:gmt_updated" json:"gmt_updated"`
}

func (DatasourceType) TableName() string {
	return "datasource_type"
}

type Datasource struct {
	Id              int       `gorm:"primarykey" json:"id"`
	Name            string    `gorm:"size:50;uniqueIndex" json:"name"`
	GroupName       string    `gorm:"size:50" json:"group_name"`
	Idc             string    `gorm:"size:30" json:"idc"`
	Env             string    `gorm:"size:30" json:"env"`
	Type            string    `gorm:"size:30" json:"type"`
	Host            string    `gorm:"size:100;index:uniq_host_port_dbid,unique" json:"host"`
	Port            string    `gorm:"size:30;index:uniq_host_port_dbid,unique" json:"port"`
	User            string    `gorm:"size:30" json:"user"`
	Pass            string    `gorm:"size:100" json:"pass"`
	Dbid            string    `gorm:"size:50;index:uniq_host_port_dbid,unique" json:"dbid"`
	Role            int32     `gorm:"default:1" json:"role"`
	Enable          int32     `gorm:"default:1" json:"enable"`
	Status          int32     `gorm:"default:1" json:"status"`
	StatusText      string    `gorm:"size:500" json:"status_text"`
	DbmetaEnable    int32     `gorm:"default:0" json:"dbmeta_enable"`
	SensitiveEnable int32     `gorm:"default:0" json:"sensitive_enable"`
	ExecuteEnable   int32     `gorm:"default:0" json:"execute_enable"`
	MonitorEnable   int32     `gorm:"default:0" json:"monitor_enable"`
	AlarmEnable     int32     `gorm:"default:0" json:"alarm_enable"`
	DmlBackupEnable int32     `gorm:"default:0" json:"dml_backup_enable"`
	DmlBackupDir    string    `gorm:"size:100" json:"dml_backup_dir"`
	CreatedAt       time.Time `gorm:"column:gmt_created" json:"gmt_created"`
	UpdatedAt       time.Time `gorm:"column:gmt_updated" json:"gmt_updated"`
}

func (Datasource) TableName() string {
	return "datasource"
}

type Privilege struct {
	Id             int64     `gorm:"primarykey" json:"id"`
	Username       string    `gorm:"size:50;index" json:"username"`
	DatasourceType string    `gorm:"size:50" json:"datasource_type"`
	Datasource     string    `gorm:"size:100" json:"datasource"`
	GrantType      string    `gorm:"size:50" json:"grant_type"`
	DatabaseName   string    `gorm:"size:50;index" json:"database_name"`
	TableName      string    `gorm:"size:50;index" json:"table_name"`
	DoSelect       int8      `gorm:"default:0" json:"do_select"`
	DoInsert       int8      `gorm:"default:0" json:"do_insert"`
	DoUpdate       int8      `gorm:"default:0" json:"do_update"`
	DoDelete       int8      `gorm:"default:0" json:"do_delete"`
	DoCreate       int8      `gorm:"default:0" json:"do_create"`
	DoAlter        int8      `gorm:"default:0" json:"do_alter"`
	MaxSelect      int       `gorm:"default:0" json:"max_select"`
	MaxUpdate      int       `gorm:"default:0" json:"max_update"`
	MaxDelete      int       `gorm:"default:0" json:"max_delete"`
	ExpireDate     time.Time `gorm:"type:date" json:"expire_date"`
	Reason         string    `gorm:"size:500" json:"reason"`
	Enable         int       `gorm:"default:1" json:"enable"`
	UserCreated    string    `gorm:"size:50;index" json:"user_created"`
	UserUpdated    string    `gorm:"size:50" json:"user_updated"`
	CreatedAt      time.Time `gorm:"column:gmt_created" json:"gmt_created"`
	UpdatedAt      time.Time `gorm:"column:gmt_updated" json:"gmt_updated"`
}

type QueryLog struct {
	Id             int64     `gorm:"primarykey"`
	Username       string    `gorm:"size:50;index"`
	DatasourceType string    `gorm:"size:50"`
	Datasource     string    `gorm:"size:100" json:"datasource"`
	Database       string    `gorm:"size:50;index" json:"database"`
	QueryType      string    `gorm:"size:50" json:"query_type"`
	SqlType        string    `gorm:"size:50" json:"sql_type"`
	Status         string    `gorm:"size:50" json:"status"`
	Times          int64     `gorm:"size:10" json:"times"`
	Content        string    `gorm:"size:1000" json:"content"`
	Result         string    `gorm:"size:1000" json:"result"`
	CreatedAt      time.Time `gorm:"column:gmt_created;index" json:"gmt_created"`
	UpdatedAt      time.Time `gorm:"column:gmt_updated" json:"gmt_updated"`
}

func (QueryLog) TableName() string {
	return "query_log"
}

type Favorite struct {
	ID             int       `gorm:"primarykey" json:"id"`
	Username       string    `gorm:"size:50;index:idx_user_datasource" json:"username"`
	DatasourceType string    `gorm:"size:50;index:idx_user_datasource" json:"datasource_type"`
	Datasource     string    `gorm:"size:100;index:idx_user_datasource" json:"datasource"`
	DatabaseName   string    `gorm:"size:50" json:"database_name"`
	Content        string    `gorm:"size:1000" json:"content"`
	CreatedAt      time.Time `gorm:"column:gmt_created" json:"gmt_created"`
	UpdatedAt      time.Time `gorm:"column:gmt_updated" json:"gmt_updated"`
}

func (Favorite) TableName() string {
	return "favorite"
}

type MetaDatabase struct {
	Id             int64     `gorm:"primarykey" json:"id"`
	DatasourceType string    `gorm:"size:50;index" json:"datasource_type"`
	Host           string    `gorm:"size:100;index:idx_host_port" json:"host"`
	Port           string    `gorm:"size:10;index:idx_host_port" json:"port"`
	DatabaseName   string    `gorm:"size:50;index" json:"database_name"`
	SchemaName     string    `gorm:"size:50" json:"schema_name"`
	Characters     string    `gorm:"size:50" json:"characters"`
	IsDeleted      int       `gorm:"default:0" json:"is_deleted"`
	CreatedAt      time.Time `gorm:"column:gmt_created;index" json:"gmt_created"`
	UpdatedAt      time.Time `gorm:"column:gmt_updated" json:"gmt_updated"`
}

func (MetaDatabase) TableName() string {
	return "meta_database"
}

type MetaTable struct {
	Id             int       `gorm:"primarykey" json:"id"`
	DatasourceType string    `gorm:"size:50;index" json:"datasource_type"`
	Host           string    `gorm:"size:100;index:idx_host_port" json:"host"`
	Port           string    `gorm:"size:10;index:idx_host_port" json:"port"`
	DatabaseName   string    `gorm:"size:50;index" json:"database_name"`
	TableType      string    `gorm:"size:50" json:"table_type"`
	TableNameX     string    `gorm:"column:table_name;size:50;index" json:"table_name"`
	TableComment   string    `gorm:"size:50" json:"table_comment"`
	Characters     string    `gorm:"size:50" json:"characters"`
	IsDeleted      int8      `gorm:"default:0" json:"is_deleted"`
	CreatedAt      time.Time `gorm:"column:gmt_created;index" json:"gmt_created"`
	UpdatedAt      time.Time `gorm:"column:gmt_updated" json:"gmt_updated"`
}

func (MetaTable) TableName() string {
	return "meta_table"
}

type MetaColumn struct {
	Id               int       `gorm:"primarykey" json:"id"`
	DatasourceType   string    `gorm:"size:50;index" json:"datasource_type"`
	Host             string    `gorm:"size:100;index:idx_host_port" json:"host"`
	Port             string    `gorm:"size:10;index:idx_host_port" json:"port"`
	DatabaseName     string    `gorm:"size:50;index" json:"database_name"`
	TableNameX       string    `gorm:"column:table_name;size:50;index" json:"table_name"`
	ColumnName       string    `gorm:"size:50;index" json:"column_name"`
	ColumnComment    string    `gorm:"size:50" json:"column_comment"`
	DataType         string    `gorm:"size:50" json:"data_type"`
	IsNullable       string    `gorm:"size:50" json:"is_nullable"`
	DefaultValue     string    `gorm:"size:50" json:"default_value"`
	Ordinal_Position int       `gorm:"default:0" json:"ordinal_position"`
	Characters       string    `gorm:"size:100" json:"characters"`
	IsDeleted        int8      `gorm:"default:0" json:"is_deleted"`
	CreatedAt        time.Time `gorm:"column:gmt_created;index" json:"gmt_created"`
	UpdatedAt        time.Time `gorm:"column:gmt_updated" json:"gmt_updated"`
}

func (MetaColumn) TableName() string {
	return "meta_column"
}

type Idc struct {
	Id          int64     `gorm:"primarykey" json:"id"`
	IdcKey      string    `gorm:"size:30;index" json:"idc_key"`
	IdcName     string    `gorm:"size:30" json:"idc_name"`
	City        string    `gorm:"size:30" json:"city"`
	Description string    `gorm:"size:300" json:"description"`
	CreatedAt   time.Time `gorm:"column:gmt_created" json:"gmt_created"`
	UpdatedAt   time.Time `gorm:"column:gmt_updated" json:"gmt_updated"`
}

func (Idc) TableName() string {
	return "idc"
}

type Env struct {
	Id          int64     `gorm:"primarykey" json:"id"`
	EnvKey      string    `gorm:"size:30;index" json:"env_key"`
	EnvName     string    `gorm:"size:30" json:"env_name"`
	Description string    `gorm:"size:300" json:"description"`
	CreatedAt   time.Time `gorm:"column:gmt_created" json:"gmt_created"`
	UpdatedAt   time.Time `gorm:"column:gmt_updated" json:"gmt_updated"`
}

func (Env) TableName() string {
	return "env"
}

type TaskOption struct {
	TaskKey         string    `gorm:"size:50;primarykey" json:"task_key"`
	TaskName        string    `gorm:"size:50," json:"task_name"`
	TaskDescription string    `gorm:"size:500," json:"task_description"`
	Crontab         string    `gorm:"size:100," json:"crontab"`
	Enable          int8      `gorm:"default:1" json:"enable"`
	CreatedAt       time.Time `gorm:"column:gmt_created" json:"gmt_created"`
	UpdatedAt       time.Time `gorm:"column:gmt_updated" json:"gmt_updated"`
}

func (TaskOption) TableName() string {
	return "task_option"
}

type TaskHeartbeat struct {
	HeartbeatKey     string    `gorm:"size:50;primarykey" json:"heartbeat_key"`
	HeartbeatTime    time.Time `gorm:"column:heartbeat_time" json:"heartbeat_time"`
	HeartbeatEndTime time.Time `gorm:"column:heartbeat_end_time" json:"heartbeat_end_time"`
	CreatedAt        time.Time `gorm:"column:gmt_created" json:"gmt_created"`
	UpdatedAt        time.Time `gorm:"column:gmt_updated" json:"gmt_updated"`
	//TaskOption       TaskOption `gorm:"references:task"`
}

func (TaskHeartbeat) TableName() string {
	return "task_heartbeat"
}

type SensitiveRule struct {
	Id          int64     `gorm:"primarykey" json:"id"`
	RuleType    string    `gorm:"size:50," json:"rule_type"`
	RuleKey     string    `gorm:"size:50;index" json:"rule_key"`
	RuleName    string    `gorm:"size:50" json:"rule_name"`
	RuleExpress string    `gorm:"size:500" json:"rule_express"`
	RulePct     int       `gorm:"default:0" json:"rule_pct"`
	Level       int8      `gorm:"default:0;comment:'0:低敏,1:高敏'" json:"level"`
	Status      int8      `gorm:"default:0;comment:'0:非敏感,-1:疑似敏感,1:确认敏感'" json:"status"`
	Enable      int8      `gorm:"default:1" json:"enable"`
	CreatedAt   time.Time `gorm:"column:gmt_created" json:"gmt_created"`
	UpdatedAt   time.Time `gorm:"column:gmt_updated" json:"gmt_updated"`
}

func (SensitiveRule) TableName() string {
	return "sensitive_rule"
}

type SensitiveMeta struct {
	Id             int64     `gorm:"primarykey" json:"id"`
	DatasourceType string    `gorm:"size:50;index" json:"datasource_type"`
	Host           string    `gorm:"size:100;index:idx_host_port" json:"host"`
	Port           string    `gorm:"size:50;index:idx_host_port" json:"port"`
	DatabaseName   string    `gorm:"size:50;index" json:"database_name"`
	TableNameX     string    `gorm:"column:table_name;size:50;index" json:"table_name"`
	TableComment   string    `gorm:"size:50" json:"table_comment"`
	ColumnName     string    `gorm:"size:50;index" json:"column_name"`
	ColumnComment  string    `gorm:"size:50" json:"column_comment"`
	RuleType       string    `gorm:"size:50" json:"rule_type"`
	RuleKey        string    `gorm:"size:50" json:"rule_key"`
	RuleName       string    `gorm:"size:50" json:"rule_name"`
	SensitiveCount int       `gorm:"default:0" json:"sensitive_count"`
	SimpleCount    int       `gorm:"default:0" json:"simple_count"`
	Level          int8      `gorm:"default:1;comment:'1:低敏,2:高敏'" json:"level"`
	Status         int8      `gorm:"default:0;comment:'0:非敏感,-1:疑似敏感,1:确认敏感'" json:"status"`
	CreatedAt      time.Time `gorm:"column:gmt_created;index" json:"gmt_created"`
	UpdatedAt      time.Time `gorm:"column:gmt_updated" json:"gmt_updated"`
}

func (SensitiveMeta) TableName() string {
	return "sensitive_meta"
}

type EventGlobal struct {
	Id             int64     `gorm:"primarykey" json:"id"`
	DatasourceType string    `gorm:"size:50;index" json:"datasource_type"`
	DatasourceName string    `gorm:"size:50;index" json:"datasource_name"`
	Host           string    `gorm:"size:100;index:idx_host_port" json:"host"`
	Port           string    `gorm:"size:50;index:idx_host_port" json:"port"`
	Version        string    `gorm:"size:50" json:"version"`
	Role           string    `gorm:"size:50" json:"role"`
	Uptime         int64     `gorm:"default:-1" json:"uptime"`
	Connect        int       `gorm:"default:-1" json:"connect"`
	Session        int       `gorm:"default:-1" json:"session"`
	Active         int       `gorm:"default:-1" json:"active"`
	Wait           int       `gorm:"default:-1" json:"wait"`
	Qps            int       `gorm:"default:-1" json:"qps"`
	Tps            int       `gorm:"default:-1" json:"tps"`
	Repl           int       `gorm:"default:-1" json:"repl"`
	Delay          int       `gorm:"default:-1" json:"delay"`
	Remark         string    `gorm:"size:1000" json:"remark"`
	CreatedAt      time.Time `gorm:"column:gmt_created;index" json:"gmt_created"`
	UpdatedAt      time.Time `gorm:"column:gmt_updated" json:"gmt_updated"`
}

func (EventGlobal) TableName() string {
	return "event_global"
}

type Event struct {
	EventUuid   string    `gorm:"primary_key" json:"event_uuid"`
	EventTime   time.Time `gorm:"column:event_time" json:"event_time"`
	EventType   string    `gorm:"" json:"event_type"`
	EventGroup  string    `gorm:"" json:"event_group"`
	EventEntity string    `gorm:"" json:"event_entity"`
	EventKey    string    `gorm:"" json:"event_key"`
	EventValue  float32   `gorm:"type:decimal(20,2)" json:"event_value"`
	EventTag    string    `gorm:"" json:"event_tag"`
	EventUnit   string    `gorm:"" json:"event_unit"`
	EventDetail string    `gorm:"" json:"event_detail"`
	//CreatedAt   time.Time `gorm:"column:gmt_created;autoCreateTime" json:"gmt_created"`
}

type EventsDescription struct {
	ID          int64  `gorm:"primarykey" json:"id"`
	EventType   string `gorm:"column:event_type" json:"eventType"`
	EventKey    string `gorm:"column:event_key" json:"eventKey"`
	Description string `gorm:"size:300" json:"description"`
}

type EventDescription struct {
	ID          int64  `gorm:"primarykey" json:"id"`
	EventType   string `gorm:"column:event_type" json:"eventType"`
	EventKey    string `gorm:"column:event_key" json:"eventKey"`
	Description string `gorm:"size:300" json:"description"`
}

func (EventDescription) TableName() string {
	return "event_description"
}

type AlarmChannel struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	Name          string    `gorm:"size:30" json:"name"`
	Description   string    `gorm:"size:1000" json:"description"`
	Enable        int       `gorm:"default:0" json:"enable"`
	MailEnable    int       `gorm:"default:0" json:"mail_enable"`
	SmsEnable     int       `gorm:"default:0" json:"sms_enable"`
	PhoneEnable   int       `gorm:"default:0" json:"phone_enable"`
	WechatEnable  int       `gorm:"default:0" json:"wechat_enable"`
	WebhookEnable int       `gorm:"default:0" json:"webhook_enable"`
	MailList      string    `gorm:"size:500" json:"mail_list"`
	SmsList       string    `gorm:"size:500" json:"sms_list"`
	PhoneList     string    `gorm:"size:500" json:"phone_list"`
	WechatList    string    `gorm:"size:500" json:"wechat_list"`
	WebhookUrl    string    `gorm:"size:500" json:"webhook_url"`
	CreatedAt     time.Time `gorm:"column:gmt_created" json:"gmt_created"`
	UpdatedAt     time.Time `gorm:"column:gmt_updated" json:"gmt_updated"`
}

func (AlarmChannel) TableName() string {
	return "alarm_channel"
}

type AlarmLevel struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	LevelName   string    `gorm:"size:30" json:"level_name"`
	Description string    `gorm:"size:1000" json:"description"`
	Enable      int       `gorm:"default:0" json:"enable"`
	CreatedAt   time.Time `gorm:"column:gmt_created" json:"gmt_created"`
	UpdatedAt   time.Time `gorm:"column:gmt_updated" json:"gmt_updated"`
}

func (AlarmLevel) TableName() string {
	return "alarm_level"
}

type AlarmRule struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:30" json:"title"`
	EventType   string    `gorm:"size:50" json:"event_type"`
	EventGroup  string    `gorm:"size:50" json:"event_group"`
	EventKey    string    `gorm:"size:50" json:"event_key"`
	EventEntity string    `gorm:"size:50" json:"event_entity"`
	AlarmRule   string    `gorm:"size:50" json:"alarm_rule"`
	AlarmValue  string    `gorm:"size:50" json:"alarm_value"`
	AlarmSleep  int       `gorm:"default:3600" json:"alarm_sleep"`
	AlarmTimes  int       `gorm:"default:3" json:"alarm_times"`
	ChannelId   int       `gorm:"default:1" json:"channel_id"`
	LevelId     int       `gorm:"default:1" json:"level_id"`
	Enable      int       `gorm:"default:1" json:"enable"`
	CreatedAt   time.Time `gorm:"column:gmt_created" json:"gmt_created"`
	UpdatedAt   time.Time `gorm:"column:gmt_updated" json:"gmt_updated"`
}

func (AlarmRule) TableName() string {
	return "alarm_rule"
}

type AlarmEvent struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	AlarmTitle  string    `gorm:"size:50" json:"alarm_title"`
	AlarmLevel  string    `gorm:"size:50" json:"alarm_level"`
	AlarmRule   string    `gorm:"size:50" json:"alarm_rule"`
	AlarmValue  string    `gorm:"size:50" json:"alarm_value"`
	EventTime   time.Time `gorm:"column:event_time" json:"event_time"`
	EventUuid   string    `gorm:"size:200" json:"event_uuid"`
	EventType   string    `gorm:"size:50" json:"event_type"`
	EventGroup  string    `gorm:"size:50" json:"event_group"`
	EventKey    string    `gorm:"size:50" json:"event_key"`
	EventValue  float64   `gorm:"type:decimal(20,2)" json:"event_value"`
	EventUnit   string    `gorm:"size:50" json:"event_unit"`
	EventEntity string    `gorm:"size:50" json:"event_entity"`
	EventTag    string    `gorm:"size:50" json:"event_tag"`
	RuleId      int64     `gorm:"default:0" json:"rule_id"`
	LevelId     int       `gorm:"default:0" json:"level_id"`
	ChannelId   int       `gorm:"default:0" json:"channel_id"`
	SendMail    int       `gorm:"default:0" json:"send_mail"`
	SendSms     int       `gorm:"default:0" json:"send_sms"`
	SendPhone   int       `gorm:"default:0" json:"send_phone"`
	SendWechat  int       `gorm:"default:0" json:"send_wechat"`
	SendWebhook int       `gorm:"default:0" json:"send_webhook"`
	Status      int       `gorm:"default:0" json:"status"`
	CreatedAt   time.Time `gorm:"column:gmt_created" json:"gmt_created"`
}

func (AlarmEvent) TableName() string {
	return "alarm_event"
}

type AlarmSuggest struct {
	ID        int64     `gorm:"primarykey" json:"id"`
	EventType string    `gorm:"size:50" json:"event_type"`
	EventKey  string    `gorm:"size:50" json:"event_key"`
	Content   string    `gorm:"size:3000" json:"content"`
	CreatedAt time.Time `gorm:"column:gmt_created" json:"gmt_created"`
}

func (AlarmSuggest) TableName() string {
	return "alarm_suggest"
}

type AlarmSendLog struct {
	ID        int64     `gorm:"primarykey" json:"id"`
	SendType  string    `gorm:"size:50" json:"send_type"`
	Receiver  string    `gorm:"size:300" json:"receiver"`
	Content   string    `gorm:"size:5000" json:"content"`
	Status    int       `gorm:"default:0" json:"status"`
	ErrorInfo string    `gorm:"size:500" json:"error_info"`
	CreatedAt time.Time `gorm:"column:gmt_created" json:"gmt_created"`
}

func (AlarmSendLog) TableName() string {
	return "alarm_send_log"
}

type AlarmTrack struct {
	ID        int64     `gorm:"primarykey" json:"id"`
	AlarmId   int64     `gorm:"alarm_id"`
	UserId    int64     `gorm:"user_id"`
	Content   string    `gorm:"size:1000"`
	CreatedAt time.Time `gorm:"column:gmt_created" json:"gmt_created"`
}

func (AlarmTrack) TableName() string {
	return "alarm_track"
}

// type StatusMysql struct {
// 	ID                            int64     `gorm:"primarykey" json:"id"`
// 	DatasourceName                string    `gorm:"column:datasource_name;size:50;not null;default:''"`
// 	Host                          string    `gorm:"column:host;size:50;not null"`
// 	Port                          string    `gorm:"column:port;size:10;not null"`
// 	Status                        int       `gorm:"column:status;not null;default:0"`
// 	StatusText                    string    `gorm:"column:status_text;size:1000;not null;default:''"`
// 	Hostname                      string    `gorm:"column:hostname;size:100;not null;default:'-1'"`
// 	Uptime                        int       `gorm:"column:uptime;not null;default:0"`
// 	Version                       string    `gorm:"column:version;size:50;not_null;default:'-1'"`
// 	Timezone                      string    `gorm:"column:timezone;size:50;not null;default:'-1'"`
// 	Role                          string    `gorm:"column:role;size:30;not null;default:'-1'"`
// 	Readonly                      string    `gorm:"column:readonly;size:10;not null;default:'-1'"`
// 	GtidMode                      string    `gorm:"column:gtid_mode;size:30;not null;default:'-1'"`
// 	AutoPosition                  *int      `gorm:"column:auto_position;not null;default:-1"`
// 	MasterHost                    string    `gorm:"column:master_host;size:50;not null;default:'-1'"`
// 	MasterPort                    string    `gorm:"column:master_port;size:10;not null;default:'-1'"`
// 	MasterUser                    string    `gorm:"column:master_user;size:30;not null;default:'-1'"`
// 	ReplStatus                    *int      `gorm:"column:repl_status;not null;default:-1"`
// 	ReplDelay                     *int      `gorm:"column:repl_delay;not null;default:-1"`
// 	MaxConnections                *int      `gorm:"column:max_connections;not null;default:-1"`
// 	OpenFilesLimit                *int      `gorm:"column:open_files_limit;not null;default:-1"`
// 	OpenFiles                     *int      `gorm:"column:open_files;not null;default:-1"`
// 	TableOpenCache                *int      `gorm:"column:table_open_cache;not null;default:-1"`
// 	OpenTables                    *int      `gorm:"column:open_tables;not null;default:-1"`
// 	ThreadsConnected              *int      `gorm:"column:threads_connected;not null;default:-1"`
// 	ThreadsRunning                *int      `gorm:"column:threads_running;not null;default:-1"`
// 	ThreadsWait                   *int      `gorm:"column:threads_wait;not null;default:-1"`
// 	ThreadsCreated                *int      `gorm:"column:threads_created;not null;default:-1"`
// 	ThreadsCached                 *int      `gorm:"column:threads_cached;not null;default:-1"`
// 	Connections                   *int      `gorm:"column:connections;not null;default:-1"`
// 	AbortedClients                *int      `gorm:"column:aborted_clients;not null;default:-1"`
// 	AbortedConnects               *int      `gorm:"column:aborted_connects;not null;default:-1"`
// 	BytesReceived                 *int      `gorm:"column:bytes_received;not null;default:-1"`
// 	BytesSent                     *int      `gorm:"column:bytes_sent;not null;default:-1"`
// 	ComSelect                     *int      `gorm:"column:com_select;not null;default:-1"`
// 	ComInsert                     *int      `gorm:"column:com_insert;not null;default:-1"`
// 	ComUpdate                     *int      `gorm:"column:com_update;not null;default:-1"`
// 	ComDelete                     *int      `gorm:"column:com_delete;not null;default:-1"`
// 	ComCommit                     *int      `gorm:"column:com_commit;not null;default:-1"`
// 	ComRollback                   *int      `gorm:"column:com_rollback;not null;default:-1"`
// 	Questions                     *int      `gorm:"column:questions;not null;default:-1"`
// 	Queries                       *int      `gorm:"column:queries;not null;default:-1"`
// 	SlowQueries                   *int      `gorm:"column:slow_queries;not null;default:-1"`
// 	InnodbPagesCreated            *int      `gorm:"column:innodb_pages_created;not null;default:-1"`
// 	InnodbPagesRead               *int      `gorm:"column:innodb_pages_read;not null;default:-1"`
// 	InnodbPagesWritten            *int      `gorm:"column:innodb_pages_written;not null;default:-1"`
// 	InnodbRowLockCurrentWaits     *int      `gorm:"column:innodb_row_lock_current_waits;not null;default:-1"`
// 	InnodbBufferPoolReadRequests  *int      `gorm:"column:innodb_buffer_pool_read_requests;not null;default:-1"`
// 	InnodbBufferPoolWriteRequests *int      `gorm:"column:innodb_buffer_pool_write_requests;not null;default:-1"`
// 	InnodbRowsRead                *int      `gorm:"column:innodb_rows_read;not null;default:-1"`
// 	InnodbRowsInserted            *int      `gorm:"column:innodb_rows_inserted;not null;default:-1"`
// 	InnodbRowsUpdated             *int      `gorm:"column:innodb_rows_updated;not null;default:-1"`
// 	InnodbRowsDeleted             *int      `gorm:"column:innodb_rows_deleted;not null;default:-1"`
// 	GmtCreate                     time.Time `gorm:"column:gmt_create;not null;autoCreateTime"`
// }

// func (StatusMysql) TableName() string {
// 	return "status_mysql"
// }
