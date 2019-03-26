#!/bin/bash

#-------------此处根据实际情况进行修改-------------
MYSQL_HOST=127.0.0.1 # 数据库Host
MYSQL_USER=          # 数据库用户名
MYSQL_PASS=          # 数据库密码
#--------------------------------------------

mysql --default-character-set=utf8mb4 -h${MYSQL_HOST} -u${MYSQL_USER} -p${MYSQL_PASS} < ./sql/syncd_v2.sql