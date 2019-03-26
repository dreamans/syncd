#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin
export PATH

#--------------------------------------------
# syncdz自动安装脚本 v1.0
# 作者：周末 i@mrjooz.com
# 功能：syncd在centos7下自动安装
# 修改日期：20190212
#--------------------------------------------

#-----------------Variable-------------------
# 需要修改的配置，其它基本上默认就行
web_port=80 #端口
web_host=172.21.100.10  #域名或服务器IP
run_user=www            #运行用户，没特殊要求可以不改
workspace=/opt/checkout #代码检出目录
sql_file=/root/syncd.sql
cipher_key=pG1L62EM0cPIIOwusQsbcV8Cs6j/M1RxLoXIZylWUC4= #加密私钥，base64编码

# 邮件
#0 - 关闭
#1 - 开启
enable=0
smtp_host=smtp.exmail.qq.com
smtp_port=465
smtp_user=
smtp_pass= 

# 软件版本
go_url=https://studygolang.com/dl/golang/go1.11.linux-amd64.tar.gz
git_url=https://github.com/git/git/archive/v2.20.1.tar.gz
git_ver=git-2.20.1
nginx_url=http://nginx.org/download/nginx-1.15.8.tar.gz
nginx_ver=nginx-1.15.8
syncd_url=https://github.com/dreamans/syncd/releases/download/1.1.2/Syncd-v1.1.2-linux-amd64.tar.gz
syncd_ver=Syncd-v1.1.2-linux-amd64
mysql_url=https://repo.mysql.com/mysql57-community-release-el7-9.noarch.rpm

# syncd配置
app_port=8868     #HTTP服务监听的端口
app_host=127.0.0.1
read_timeout=300  #读超时时间设置, 单位秒
write_timeout=300 #写超时时间设置, 单位秒
idle_timeout=300  #空闲连接超时设置, 单位秒

# log
#stdout - 打印到标准输出
#/path/file - 输出到文件
log=stdout

# mysql
host=127.0.0.1
port=3306
user=syncd
password=syncd
dbname=syncd

#-----------------Function-------------------
Color_Text()
{
  echo -e " \e[0;$2m$1\e[0m"
}

Echo_Blue()
{
  echo $(Color_Text "$1" "34")
}

Echo_Green()
{
  echo $(Color_Text "$1" "32")
}

Echo_Yellow()
{
  echo $(Color_Text "$1" "33")
}

Set_Timezone()
{
  Echo_Blue "时区设置..."
  rm -rf /etc/localtime
  ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
}

Install_Go()
{
  Echo_Blue "[+] 安装go..."
  if [ -d "/usr/local/go/" ];then
	Echo_Yellow "go已经安装！"
  else
    cd /usr/local/src
    wget -cq $go_url -O go.tar.gz
    tar zxf go.tar.gz -C /usr/local
    echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
    source /etc/profile
  fi
}

Install_Git()
{
  Echo_Blue "[+] 安装git..."
  if [ -d "/usr/local/git/" ];then
	Echo_Yellow "git已经安装！"
  else
    cd /usr/local/src
    yum remove git -y
    wget -cq $git_url -O ${git_ver}.tar.gz
    tar zxf ${git_ver}.tar.gz
    cd ${git_ver}
    make configure
    ./configure --prefix=/usr/local/git
    make
    make install
    echo 'export PATH=$PATH:/usr/local/git/bin' >> /etc/profile
    source /etc/profile
  fi
}

Install_Mysql()
{
  Echo_Blue "[+] 安装mysql..."
  result=`yum list installed | grep mysql | wc -l`
  if [ $result -ne 0 ];then
    Echo_Yellow "mysql已经安装！"
  else
    yum localinstall $mysql_url -y
    yum install mysql-community-server -y
    systemctl start mysqld.service
    systemctl enable mysqld
  
    #获取初始密码
    rootpwd=`grep 'temporary password' /var/log/mysqld.log | cut -d ' ' -f 11`
    #修改初始密码
    mysql -uroot -p$rootpwd --connect-expired-password << EOF
SET GLOBAL validate_password_length = 0;
SET GLOBAL validate_password_number_count = 0;
SET GLOBAL validate_password_policy = LOW;
ALTER USER 'root'@'localhost' IDENTIFIED BY 'root';
flush privileges; 
EOF
    #创建数据库
    mysql -uroot -proot --connect-expired-password << EOF
create database $dbname default charset utf8;
grant all privileges on $dbname.* to $user@'%localhost' identified by "$password";
flush privileges; 
EOF

    mysql -u$user -p$password $dbname < $sql_file  --connect-expired-password
  fi
}

Install_Syncd()
{
  Echo_Blue "[+] 安装syncd..."
  if [ -d "/opt/syncd/" ];then
	Echo_Yellow "syncd已经安装！"
  else
    cd /opt
    wget $syncd_url
    tar zxf ${syncd_ver}.tar.gz
    mv $syncd_ver syncd
    cd syncd
	
# 自动生成配置
(
cat <<EOF
[syncd]
workspace = $workspace
cipher_key = $cipher_key

[serve]
addr = ${app_host}:${app_port}
read_timeout = $read_timeout
write_timeout = $write_timeout
idle_timeout = $idle_timeout

[database]
host = $host
port = $port
user = $user
password = $password
dbname = $dbname

[log]
path = $log

[mail]
enable = $enable
smtp_host = $smtp_host
smtp_port = $smtp_port
smtp_user = $smtp_user
smtp_pass = $smtp_pass
EOF
) > /opt/syncd/etc/syncd.ini

    nohup ./bin/syncd -c ./etc/syncd.ini > syncd.log 2>&1 &
  fi
}

Install_Nginx()
{
  Echo_Blue "[+] 安装nginx..."
  if [ -d "/usr/local/nginx/" ];then
	Echo_Yellow "nginx已经安装！"
  else
    cd /usr/local/src/
    wget -cq $nginx_url
    tar zxf ${nginx_ver}.tar.gz
    cd $nginx_ver
    ./configure --user=$run_user --group=$run_user --prefix=/usr/local/nginx --with-http_stub_status_module --with-http_ssl_module --with-http_v2_module --with-http_gzip_static_module --with-http_sub_module
    make
    make install
 
# 自动生成配置
(
cat <<EOF
[Unit]
Description=nginx 
After=network.target 
   
[Service] 
Type=forking 
ExecStart=/usr/local/nginx/sbin/nginx
ExecReload=/usr/local/nginx/sbin/nginx -s reload
ExecStop=/usr/local/nginx/sbin/nginx -s quit
PrivateTmp=true 
   
[Install] 
WantedBy=multi-user.target
EOF
) > /lib/systemd/system/nginx.service

(
cat <<EOF
#主配置
user  $run_user;
worker_processes  auto;

#错误日志
error_log  logs/error.log;
error_log  logs/error.log  notice;
error_log  logs/error.log  info;

pid        logs/nginx.pid;

events {
    use epoll;
    worker_connections  1024;
    multi_accept on;
}


http {
    include       mime.types;
    default_type  application/octet-stream;

    sendfile        on;
    tcp_nopush     on;

    keepalive_timeout  65;

    #gzip
    gzip on;
    gzip_min_length  1k;
    gzip_buffers     4 16k;
    gzip_http_version 1.1;
    gzip_comp_level 2;
    gzip_types     text/plain application/javascript application/x-javascript text/javascript text/css application/xml application/xml+rss;
    gzip_vary on;
    gzip_proxied   expired no-cache no-store private auth;
    gzip_disable   "MSIE [1-6]\.";

    #隐藏版本号
    server_tokens off;
 
    #关闭默认日志
    access_log off;

    #引入虚拟主机配置
    include vhost/*.conf;

}
EOF
) > /usr/local/nginx/conf/nginx.conf

    mkdir /usr/local/nginx/conf/vhost
(
cat <<EOF
upstream syncdServer {
    server $app_host:$app_port weight=1;
}
server {
    listen       $web_port;
    server_name  $web_host;
    access_log   logs/syncd.log;

    location / {
        try_files \$uri \$uri/ /index.html;
        root /opt/syncd/public;
        index index.html index.htm;
    }

    location ^~ /api/ {
        proxy_pass          http://syncdServer;
        proxy_set_header    X-Forwarded-Host \$host:\$server_port;
        proxy_set_header    X-Real-IP     \$remote_addr;
        proxy_set_header    Origin        \$host:\$server_port;
        proxy_set_header    Referer       \$host:\$server_port;
    }
}
EOF
) > /usr/local/nginx/conf/vhost/syncd.conf

    systemctl enable nginx.service
    systemctl start nginx.service
  fi
}
#------------------Debug---------------------

#-------------------Main--------------------
# 检测是否是root用户
if [ $(id -u) != "0" ]; then
    echo "安装失败，请在root用户下执行脚本"
    exit 1
fi

# 安装依赖
yum install gcc gcc-c++ libstdc++-devel libcurl-devel expat-devel openssl-devel install autoconf automake libtool -y

# 设置时区
Set_Timezone

# 新增用户和组
groupadd $run_user
useradd --shell /sbin/nologin -g $run_user $run_user

clear
echo "+------------------------------------------------------------------------+"
echo "|                   syncdz自动安装脚本正在执行                           |"
echo "+------------------------------------------------------------------------+"

# 安装go
Install_Go

# 安装git
Install_Git

# 安装mysql
Install_Mysql

# 安装syncd
Install_Syncd

# 安装nginx
Install_Nginx

Echo_Green "+------------------------------------------------------------------------+"
Echo_Green "安装完成！"
Echo_Green "访问地址：http://$web_host:$web_port："
Echo_Green "登录帐号：syncd"
Echo_Green "登录密码：syncd.cc"
Echo_Green "nginx操作：systemctl {start|stop|reload|enable} nginx.service"
Echo_Green "+------------------------------------------------------------------------+"

