# 安装

安装 Syncd 只需几分钟时间，若您在安装过程中遇到问题或无法找到解决方式，请[提交Issue](https://github.com/dreamans/syncd/issues)，我们会尽力解决您的问题。

## 环境需求

**操作系统**

Linux / macOS + Bash. 需要注意的是Syncd不支持Win系统。

**Go 编译环境**

Syncd依赖 `Go1.11` 编译环境，可前往[官方网站](https://golang.org/dl/) 或 [国内镜像](https://golang.google.cn/dl/) 下载安装。

**MySQL**

MySQL 5.6+

**Git**

升级操作系统Git到最新版本。

## 编译Syncd

执行以下命令
```shell
$ cd ~
$ git clone git@github.com:dreamans/syncd.git
$ cd ./syncd
$ make
```

若编译程序未报错并且在 `./synce` 中生成 `output` 子目录，则表示安装成功。

output目录结构

```shell
output // output可修改为任意其他目录名称
├── bin // bin目录存放Syncd的可执行文件
│   └── syncd
├── etc // bin/syncd 程序运行时若不指定配置文件，则会在etc目录中查找syncd.ini作为默认配置
│   └── syncd.ini
├── log
└── public // 静态资源目录
    ├── css
    ├── favicon.ico
    ├── fonts
    ├── img
    ├── index.html
    └── js
```

> 生成的 output 目录可拷贝或移动到你想要的地方，但不要试图将此目录拷贝到其他服务器上运行，会造成不可预料的结果。

## 导入数据库

使用编辑器打开 `~/syncd/script/install_sql.sh` 文件

修改此处MySQL数据库连接配置
```
#-------------此处根据实际情况进行修改-------------
MYSQL_HOST=127.0.0.1 # 数据库Host
MYSQL_USER=          # 数据库用户名
MYSQL_PASS=          # 数据库密码
#--------------------------------------------
```

运行脚本（Warning错误可忽略）

```
$ /bin/bash ~/syncd/script/install_sql.sh
```

## 配置文件

配置文件位置: `~/syncd/output/etc/syncd.ini`

修改数据库连接信息（查看[配置项](setting.md)详细文档）

```ini
[database]
host = 127.0.0.1
port = 3306
user = syncd
password = syncd
dbname = syncd
```

## 运行程序

```shell
$ cd ~/syncd/output
$ ./bin/syncd

   _____   __  __   ____     _____   ____/ /
  / ___/  / / / /  / __ \   / ___/  / __  /
 (__  )  / /_/ /  / / / /  / /__   / /_/ /
/____/   \__, /  /_/ /_/   \___/   \__,_/
        /____/

Service:              syncd
Version:              v2.0.0
Config Loaded:        /Users/work/syncd/output/etc/syncd.ini
Start Running...
```

打开浏览器，访问 `http://localhost:8878`

初始账号：
```
用户名: syncd
密码: 111111
```
**!!!登录后尽快修改默认密码**