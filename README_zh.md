# [![pEl1emT.png](https://s21.ax1x.com/2025/02/22/pEl1emT.png)](https://imgse.com/i/pEl1emT)

## DBQuery：简洁、专业、安全的数据库查询管理平台
DBQuery数据库查询平台，由LEPUS开源数据库监控系统(lepus.cc)作者开发，致力于打造简洁、智能、强大、安全的开源数据库一站式查询管理平台。

[[**English**](README.md)] | [[**简体中文**](README_zh.md)]

###  ✨ 特性：
- 支持MySQL、Oracle、MariaDB、GreatSQL、PostgreSQL、Redis、MongoDB、SQLServer、TiDB、Doris、OceanBase、ClickHouse等各类数据库的SQL执行和数据查询。
- 支持各类数据库的执行计划、索引、表结构、建表SQL、表容量等数据查询。
- 支持各类数据库的的元数据信息自动采集和查询。
- 支持自动发现高风险执行SQL并拦截。
- 支持SQL执行审计功能。
- Dashboard报告。
- 内置自动化任务调度系统。
- 内置完整的数据源管理功能。
- 支持国语言国际化、全屏模式、暗色风格切换。
- 支持Windows和Linux多平台部署。

### 🧩 技术
- 基于golang、nodejs/Antd开发，前后端分离。
- 强大的数据源支持，统一各类数据库驱动。
- 核心敏感数据采用AES并加盐加密，高安全性。
- 前后端一体化打包技术，支持一键部署启动。

### 💬 <span style="color: #568DF4;">各位亲爱的朋友们，如果您对这个项目感兴趣，请先给我一个<i style="color: #EA2626;">star</i>吧，谢谢！</span>💕
- 如有安装和使用问题，欢迎加入微信交流群（加Ruyi-1024备注DBQuery入群）
- 软件快速迭代开发中,请大家优先测试和使用最新发布的版本。
- 欢迎大家提出宝贵的建议，欢迎提issue,PR。当然更欢迎[赞助我](#都划到这了如果我的项目对您有帮助请赞助我吧)。💕

[![pEQLbTA.png](https://s21.ax1x.com/2025/02/21/pEQLbTA.png)](https://imgse.com/i/pEQLbTA)

## 立即体验
### 📦 快速部署
使用我们编译后的二进制包快速安装，适合不需要二次开发和无编程经验的用户，步骤如下：
```bash
$ cd dbquery
$ sh install.sh  //安装软件服务
$ vim /etc/dbquery/setting.yml  //修改配置文件数据库连接地址
$ sh start.sh  //启动服务
$ sh status.sh  //查看服务状态
$ sh stop.sh  //停止服务
```
> #服务运行后使用浏览器访问：http://127.0.0.1:8086 登录系统，默认账号密码:admin/dbqueryadmin

### 🦄 源码部署
#### 前置要求
- Node.js 20.18
- Npm 10.8
- Golang 1.19

#### 安装步骤
```bash
# 克隆仓库
$ git clone https://github.com/ruyi1024/dbquery.git

# 部署后端
$ go mod tidy
$ go mod vendor
$ cp setting.example.yml setting.yml  //从配置模板创建配置文件，从修改数据库连接地址
$ go run main.go

# 部署前端
$ cd web
$ npm install
$ npm start
```
>  前端地址：http://127.0.0.1:8000 
>  后端地址：http://127.0.0.1:8086
>  默认账号密码:admin/dbqueryadmin


## 许可证
本项目采用 [GPL 3.0](https://www.gnu.org/software/shishi/manual/html_node_db/a7966.html) 授权。

## 参与贡献
如果你对本项目感兴趣，可以打开一个Issue或者提交PR，请确保遵循项目的代码规范。


## 💕 感谢以下项目
- [ant.design](https://ant.design/index-cn)
- [china-alert](https://github.com/china-alert/ueh)
- [lepus](https://github.com/ruyi1024/lepus)

## 😘 赞助我
### 都划到这了，如果我的项目对您有帮助，请赞助我吧！😊😊😊
| 微信 | 支付宝  |
|-----|-----| 
| [![pEl1Kk4.jpg](https://s21.ax1x.com/2025/02/22/pEl1Kk4.jpg)](https://imgse.com/i/pEl1Kk4)  | [![pEl1Qh9.jpg](https://s21.ax1x.com/2025/02/22/pEl1Qh9.jpg)](https://imgse.com/i/pEl1Qh9) |

## 🐳 关于技术支持申明
- 本软件基于开源技术构建，使用Golang、Nodejs、Antd、MySQL等开源项目。 技术上如有问题，可以先向对应的开源社区请求帮助。
- 开源不易，本人精力和时间有限，如需一对一技术支持，请先赞助后联系微信Ruyi-1024(备注DBQuery技术支持)：

| 技术支持方式                          | 赞助(元) | 
|:--------------------------------|:-----:|
| 加微信：Ruyi-1024| 100/次 |
| 长期技术支持（不限次数，新功能优先体验等）            | 1000  | 
