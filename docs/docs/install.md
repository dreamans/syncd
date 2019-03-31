# 安装

安装 Syncd 只需几分钟时间，若您在安装过程中遇到问题或无法找到解决方式，请[提交Issue](https://github.com/dreamans/syncd/issues)，我们会尽力解决您的问题。

## 环境需求

**操作系统**

Linux / macOS + Bash. 需要注意的是Syncd不支持Win系统。

**Go 编译环境**

Syncd依赖 `Go1.11+` 编译环境，可前往[官方网站](https://golang.org/dl/) 或 [国内镜像](https://golang.google.cn/dl/) 下载安装。

**MySQL**

MySQL 5.6+

**Git**

升级操作系统Git到最新版本。

## 编译Syncd

执行以下命令
```shell
$ curl https://syncd.cc/install.sh | bash
```

若编译程序未报错并且在当前目录中生成 `syncd-deploy` 子目录，则表示安装成功。

syncd-deploy目录结构

```shell
syncd-deploy // syncd-deploy可修改为任意其他目录名称
├── bin // bin目录存放Syncd的可执行文件
│   └── syncd
├── etc // bin/syncd 程序运行时若不指定配置文件，则会在etc目录中查找syncd.ini作为默认配置
│   └── syncd.ini
├── log
├── public // 静态资源目录
    ├── css
    ├── favicon.ico
    ├── fonts
    ├── img
    ├── index.html
    └── js
└── resource // 资源目录
    └── sql
```

> 生成的 syncd-deploy 目录可拷贝或移动到你想要的地方，但不要试图将此目录拷贝到其他服务器上运行，会造成不可预料的结果。

## 导入数据库

Syncd依赖的MySQL数据表结构存在于 `syncd-deploy/resource/sql/syncd_{版本号}.sql` 中

```
$ cd ./syncd-deploy
$ mysql --default-character-set=utf8mb4 < ./resource/sql/syncd_{版本号}.sql # 导入MySQL表结构
```

## 配置文件

配置文件位置: `syncd-deploy/etc/syncd.ini`

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
$ cd syncd-deploy
$ ./bin/syncd

   _____   __  __   ____     _____   ____/ /
  / ___/  / / / /  / __ \   / ___/  / __  /
 (__  )  / /_/ /  / / / /  / /__   / /_/ /
/____/   \__, /  /_/ /_/   \___/   \__,_/
        /____/

Service:              syncd
Version:              v2.0.0
Config Loaded:        ./etc/syncd.ini
Log:                  stdout
Mail Enable:          0
HTTP Service:         :8878
Start Running...
```

打开浏览器，访问 `http://localhost:8878`

初始账号：
```
用户名: syncd
密码: 111111
```
**!!!登录后尽快修改默认密码**

[filename](include/footer.md ':include')