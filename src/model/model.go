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
	Id            int       `gorm:"primarykey" json:"id"`
	Name          string    `gorm:"size:50;uniqueIndex" json:"name"`
	GroupName     string    `gorm:"size:50" json:"group_name"`
	Idc           string    `gorm:"size:30" json:"idc"`
	Env           string    `gorm:"size:30" json:"env"`
	Type          string    `gorm:"size:30" json:"type"`
	Host          string    `gorm:"size:100;index:uniq_host_port_dbid,unique" json:"host"`
	Port          string    `gorm:"size:30;index:uniq_host_port_dbid,unique" json:"port"`
	User          string    `gorm:"size:30" json:"user"`
	Pass          string    `gorm:"size:100" json:"pass"`
	Dbid          string    `gorm:"size:50;index:uniq_host_port_dbid,unique" json:"dbid"`
	Role          int32     `gorm:"default:1" json:"role"`
	Enable        int32     `gorm:"default:1" json:"enable"`
	Status        int32     `gorm:"default:1" json:"status"`
	StatusText    string    `gorm:"size:500" json:"status_text"`
	DbmetaEnable  int32     `gorm:"default:0" json:"dbmeta_enable"`
	ExecuteEnable int32     `gorm:"default:0" json:"execute_enable"`
	CreatedAt     time.Time `gorm:"column:gmt_created" json:"gmt_created"`
	UpdatedAt     time.Time `gorm:"column:gmt_updated" json:"gmt_updated"`
}

func (Datasource) TableName() string {
	return "datasource"
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
