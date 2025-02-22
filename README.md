# [![pEl1emT.png](https://s21.ax1x.com/2025/02/22/pEl1emT.png)](https://imgse.com/i/pEl1emT)

## DBQuery：Professional and secure database query platform
The database query platform is developed by the author of the LEPUS open-source database monitoring system (lepus. cc), dedicated to creating a simple, intelligent, powerful, and secure one-stop query management platform for open-source databases.

[[**English**](README.md)] | [[**简体中文**](README_zh.md)]

###  ✨ Features：
- Support MySQL Oracle、MariaDB、GreatSQL、PostgreSQL、Redis、MongoDB、SQLServer、TiDB、Doris、OceanBase、ClickHouse Waiting for SQL execution and data queries of various databases.
- Support data queries such as execution plans, indexes, table structures, SQL tables, and table capacity for various databases.
- Support automatic collection and querying of metadata information for various databases.
- Support automatic detection of high-risk SQL execution and interception.
- Support SQL execution audit function.
- Dashboard report.
- Built in automated task scheduling system.
- Built in complete data source management functionality.
- Support internationalization of Mandarin, full screen mode, and switching between dark color styles.
-Supports multi platform deployment of Windows and Linux.

### 🧩 Technical

- Developed based on Golang, nodejs/Antd, with front-end and back-end separation.
- Powerful data source support, unifying various database drivers.
- The core sensitive data is encrypted with AES and salt for high security.
- Integrated front-end and back-end packaging technology, supporting one click deployment and startup

### 💬 <span style="color: #568DF4;">Dear friends, if you are interested in this project, please give me a <i style="color: #EA2626;">Star</i>first. Thank you!</span>💕
- If you have any installation or usage issues, please feel free to join the WeChat communication group (add Ruyi-1024 remark DBQuery to join the group)
- In the rapid iteration development of software, please prioritize testing and using the latest released version.
- Welcome everyone to provide valuable suggestions, raise issues, PR.💕

[![pEQL7eH.png](https://s21.ax1x.com/2025/02/21/pEQL7eH.png)](https://imgse.com/i/pEQL7eH)

## Install
### 📦  Quick start
Use our compiled binary package for quick installation, suitable for users who do not require secondary development and have no programming experience. The steps are as follows:
```bash

$ cd dbquery
$ sh install.sh  //Install software
$ vim /etc/dbquery/setting.yml  //Modify configuration files database connection addresses
$ sh start.sh  //Start service
$ sh status.sh  //View service status
$ sh stop.sh  //Stop service
```
> After the service is running, use a browser to access: http://127.0.0.1:8086 Login to the system, default account password：admin/dbqueryadmin


### 🦄  Build deployment
#### Requirements

- Node.js 20.18
- Npm 10.8
- Golang 1.19

#### Build step

```bash
# clond dbquery
$ git clone https://github.com/ruyi1024/dbquery.git

# build backend
$ go mod tidy
$ go mod vendor
$ cp setting.example.yml setting.yml  //Create configuration files from configuration templates and modify database connection addresses
$ go run main.go

# build frontend
$ cd web
$ npm install
$ npm start
```
>  Frontend：http://127.0.0.1:8000 
>  Backend：http://127.0.0.1:8086
>  Default Admin User: admin/dbqueryadmin

## License
License is [GPL 3.0](https://www.gnu.org/software/shishi/manual/html_node_db/a7966.html) 

## Contribute
If you are interested in this project, you can open an issue or submit a PR, please ensure to follow the project's code specifications.

## Thanks
- [ant.design](https://ant.design/index-cn)
- [china-alert](https://github.com/china-alert/ueh)
- [lepus](https://github.com/ruyi1024/lepus)
