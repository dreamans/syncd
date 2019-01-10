# ************************************************************
# Copyright 2018 syncd Author. All Rights Reserved.
# Use of this source code is governed by a MIT-style
# license that can be found in the LICENSE file.
# ************************************************************

CREATE TABLE `syd_deploy_apply` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `project_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '项目ID',
  `space_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '空间ID',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '名称',
  `description` varchar(2000) NOT NULL DEFAULT '' COMMENT '描述信息',
  `repo_data` varchar(10000) NOT NULL DEFAULT '' COMMENT '部署参数',
  `status` int(11) unsigned NOT NULL DEFAULT '1' COMMENT '状态，1 - 待审核，2 - 审核不通过，3 - 审核通过，4 - 上线中，5 - 上线成功，6 - 上线失败，7 - 废弃',
  `error_log` mediumtext NOT NULL COMMENT '上线单错误信息',
  `user_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `ctime` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_space_proj` (`space_id`,`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `syd_deploy_task` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `apply_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '上线单ID',
  `level` int(11) unsigned NOT NULL DEFAULT '1' COMMENT '执行优先级, 1-更新版本库, 2-检出指定版本, 3-打包代码, 4-部署代码(多)',
  `cmd` text NOT NULL COMMENT '执行的命令',
  `status` int(11) unsigned NOT NULL DEFAULT '1' COMMENT '任务运行状态，1-未开始，2-进行中，3-已结束，4-错误并退出, 5-终止',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '任务名称',
  `output` mediumtext NOT NULL COMMENT '命令输出信息',
  PRIMARY KEY (`id`),
  KEY `idx_apply_id` (`apply_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `syd_operate_log` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `data_id` int(11) unsigned NOT NULL DEFAULT '0',
  `op_type` varchar(10) NOT NULL DEFAULT '',
  `op_name` varchar(100) NOT NULL DEFAULT '',
  `op_content` varchar(1000) NOT NULL DEFAULT '',
  `user_id` int(11) unsigned NOT NULL DEFAULT '0',
  `user_name` varchar(100) NOT NULL DEFAULT '',
  `ctime` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_data_id` (`data_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `syd_project` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '项目名称',
  `description` varchar(500) NOT NULL DEFAULT '' COMMENT '项目描述',
  `space_id` int(11) unsigned NOT NULL DEFAULT '0',
  `repo_url` varchar(200) NOT NULL DEFAULT '' COMMENT '代码仓库地址',
  `repo_mode` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '代码上线模式',
  `repo_branch` varchar(20) NOT NULL DEFAULT '',
  `deploy_server` varchar(2000) NOT NULL DEFAULT '' COMMENT '服务器组列表, 半角'',''相隔',
  `deploy_user` varchar(20) NOT NULL DEFAULT '' COMMENT '部署用户',
  `deploy_path` varchar(100) NOT NULL DEFAULT '' COMMENT '部署目录',
  `deploy_timeout` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '部署脚本超时时间，单位秒',
  `audit_notice_email` varchar(1000) NOT NULL DEFAULT '' COMMENT '审核邮件通知',
  `deploy_notice_email` varchar(1000) NOT NULL DEFAULT '' COMMENT '上线邮件通知',
  `pre_deploy_cmd` varchar(2000) NOT NULL DEFAULT '' COMMENT '部署前执行脚本',
  `post_deploy_cmd` varchar(2000) NOT NULL DEFAULT '' COMMENT '部署后执行脚本',
  `need_audit` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '是否需要审核, 0 - 不需要，1 - 需要',
  `exclude_files` varchar(1000) NOT NULL DEFAULT '' COMMENT '上线需要排除的文件',
  `status` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '是否启用，0 - 不启用，1 - 启用',
  `ctime` int(11) NOT NULL DEFAULT '0' COMMENT '最后更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_space_id` (`space_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `syd_project_space` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `description` varchar(500) NOT NULL DEFAULT '',
  `ctime` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `syd_project_user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `space_id` int(11) unsigned NOT NULL DEFAULT '0',
  `user_id` int(11) unsigned NOT NULL DEFAULT '0',
  `ctime` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_space_uid` (`space_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `syd_server` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `group_id` int(11) unsigned NOT NULL DEFAULT '0',
  `name` varchar(100) NOT NULL DEFAULT '',
  `ip` varchar(50) NOT NULL DEFAULT '',
  `ssh_port` int(11) unsigned NOT NULL DEFAULT '0',
  `ctime` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_group_id` (`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `syd_server_group` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(100) NOT NULL DEFAULT '',
  `ctime` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


CREATE TABLE `syd_user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `group_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '用户组ID',
  `name` varchar(100) NOT NULL DEFAULT '' COMMENT '用户名',
  `password` char(32) NOT NULL DEFAULT '' COMMENT '密码密文',
  `email` varchar(100) NOT NULL DEFAULT '' COMMENT '邮箱',
  `true_name` varchar(20) NOT NULL DEFAULT '' COMMENT '真实姓名',
  `mobile` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号',
  `salt` char(10) NOT NULL DEFAULT '' COMMENT '盐',
  `lock_status` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '是否锁定，1 - 允许登录，0 - 禁止登录',
  `last_login_ip` char(15) NOT NULL DEFAULT '' COMMENT '上次登录IP',
  `last_login_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '上次登录时间',
  `ctime` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_unq_name` (`name`),
  UNIQUE KEY `idx_unq_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


INSERT INTO `syd_user` (`id`, `group_id`, `name`, `password`, `email`, `true_name`, `mobile`, `salt`, `lock_status`, `last_login_ip`, `last_login_time`, `ctime`)
VALUES
	(1,1,'syncd','7cd0c949b9998a4e8717ffef14054e4c','syncd@syncd.cc','Syncd','','0s6rrecC50',1,'::1',1547126728,0);

CREATE TABLE `syd_user_group` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `priv` varchar(10000) NOT NULL DEFAULT '',
  `ctime` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


INSERT INTO `syd_user_group` (`id`, `name`, `priv`, `ctime`)
VALUES
	(1,'超级管理员','1001,1002,1010,1004,1006,1008,1003,1005,1007,1009,2001,2002,2003,2004,2100,2101,2102,2201,2202,2203,2204,2205,2206,2207,3001,3002,3003,3004,3101,3102,3103,3104,4001,4002,4003,4004,4101,4102,4103,4104',0);


CREATE TABLE `syd_user_token` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '用户ID',
  `token` char(40) NOT NULL DEFAULT '' COMMENT '登录状态token',
  `expire_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '过期时间',
  `ctime` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
