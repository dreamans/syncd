# ************************************************************
# Copyright 2019 syncd Author. All Rights Reserved.
# Use of this source code is governed by a MIT-style
# license that can be found in the LICENSE file.
# ************************************************************

DROP TABLE IF EXISTS `syd_deploy_apply`;

CREATE TABLE `syd_deploy_apply` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `space_id` int(11) unsigned NOT NULL DEFAULT '0',
  `project_id` int(11) unsigned NOT NULL DEFAULT '0',
  `name` varchar(100) NOT NULL DEFAULT '',
  `description` varchar(500) NOT NULL DEFAULT '',
  `branch_name` varchar(100) NOT NULL DEFAULT '',
  `commit_version` varchar(100) NOT NULL DEFAULT '',
  `audit_status` int(11) unsigned NOT NULL DEFAULT '0',
  `audit_refusal_reasion` varchar(500) NOT NULL DEFAULT '',
  `status` int(11) unsigned NOT NULL DEFAULT '0',
  `user_id` int(11) unsigned NOT NULL DEFAULT '0',
  `rollback_id` int(10) unsigned NOT NULL DEFAULT '0',
  `rollback_apply_id` int(10) unsigned NOT NULL DEFAULT '0',
  `is_rollback_apply` int(11) unsigned NOT NULL DEFAULT '0',
  `ctime` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


DROP TABLE IF EXISTS `syd_deploy_build`;

CREATE TABLE `syd_deploy_build` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `apply_id` int(11) unsigned NOT NULL DEFAULT '0',
  `start_time` int(11) unsigned NOT NULL DEFAULT '0',
  `finish_time` int(11) unsigned NOT NULL DEFAULT '0',
  `status` int(11) unsigned NOT NULL DEFAULT '1',
  `tar` varchar(500) NOT NULL DEFAULT '',
  `output` text NOT NULL,
  `errmsg` varchar(1000) NOT NULL DEFAULT '',
  `ctime` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `syd_deploy_task`;

CREATE TABLE `syd_deploy_task` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `apply_id` int(11) unsigned NOT NULL DEFAULT '0',
  `group_id` int(11) unsigned NOT NULL DEFAULT '0',
  `server_id` int(11) unsigned NOT NULL DEFAULT '0',
  `start_time` int(11) NOT NULL DEFAULT '0',
  `finish_time` int(11) unsigned NOT NULL DEFAULT '0',
  `status` int(11) unsigned NOT NULL DEFAULT '0',
  `output` text NOT NULL,
  `errmsg` varchar(1000) NOT NULL DEFAULT '',
  `ctime` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `syd_project`;

CREATE TABLE `syd_project` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `space_id` int(11) unsigned NOT NULL DEFAULT '0',
  `name` varchar(100) NOT NULL DEFAULT '',
  `description` varchar(500) NOT NULL DEFAULT '',
  `need_audit` int(11) unsigned NOT NULL DEFAULT '0',
  `status` int(11) unsigned NOT NULL DEFAULT '0',
  `repo_url` varchar(500) NOT NULL DEFAULT '',
  `deploy_mode` int(11) unsigned NOT NULL DEFAULT '0',
  `repo_branch` varchar(100) NOT NULL DEFAULT '',
  `pre_release_cluster` int(11) unsigned NOT NULL DEFAULT '0',
  `online_cluster` varchar(2000) NOT NULL DEFAULT '',
  `deploy_user` varchar(50) NOT NULL DEFAULT '',
  `deploy_path` varchar(500) NOT NULL DEFAULT '',
  `build_script` text NOT NULL,
  `pre_deploy_cmd` text NOT NULL,
  `after_deploy_cmd` text NOT NULL,
  `deploy_timeout` int(11) unsigned NOT NULL DEFAULT '0',
  `ctime` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `syd_project_member`;

CREATE TABLE `syd_project_member` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `space_id` int(11) unsigned NOT NULL DEFAULT '0',
  `user_id` int(11) unsigned NOT NULL DEFAULT '0',
  `ctime` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `syd_project_space`;

CREATE TABLE `syd_project_space` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `description` varchar(2000) NOT NULL DEFAULT '',
  `ctime` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `syd_server`;

CREATE TABLE `syd_server` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `group_id` int(11) unsigned NOT NULL DEFAULT '0',
  `name` varchar(100) NOT NULL DEFAULT '',
  `ip` varchar(100) NOT NULL DEFAULT '',
  `ssh_port` int(11) unsigned NOT NULL DEFAULT '0',
  `ctime` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `syd_server_group`;

CREATE TABLE `syd_server_group` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `ctime` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `syd_user`;

CREATE TABLE `syd_user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `role_id` int(11) unsigned NOT NULL DEFAULT '0',
  `username` varchar(20) NOT NULL DEFAULT '',
  `password` char(32) NOT NULL DEFAULT '',
  `salt` char(10) NOT NULL DEFAULT '',
  `truename` varchar(10) NOT NULL DEFAULT '',
  `mobile` varchar(20) NOT NULL DEFAULT '',
  `email` varchar(500) NOT NULL DEFAULT '',
  `status` int(11) unsigned NOT NULL DEFAULT '0',
  `last_login_time` int(11) unsigned NOT NULL DEFAULT '0',
  `last_login_ip` varchar(50) NOT NULL DEFAULT '',
  `ctime` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `syd_user` (`id`, `role_id`, `username`, `password`, `salt`, `truename`, `mobile`, `email`, `status`, `last_login_time`, `last_login_ip`, `ctime`)
VALUES
	(28,5,'syncd','279353caa449a2920e272e1c70bee1ed','lkglzLskNC','','','syncd@syncd.cc',1,0,'',1548856722);

DROP TABLE IF EXISTS `syd_user_role`;

CREATE TABLE `syd_user_role` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL DEFAULT '',
  `privilege` varchar(2000) NOT NULL DEFAULT '',
  `ctime` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `syd_user_role` (`id`, `name`, `privilege`, `ctime`)
VALUES
	(5,'管理员','2001,2002,2003,2004,2100,2101,2102,2201,2202,2203,2204,2205,2206,3001,3002,3004,3003,3101,3102,3103,3104,4001,4002,4003,4004,4101,4102,4103,4104,1001,1002,1006,1003,1004,1005',1548315492),
	(6,'运维','2001,2002,2003,2004,2100,2101,2102,2201,2202,2203,2204,2205,2206,4001,4002,4003,4004,4101,4102,4103,4104,1002,1006,1003,1004,1005,1001',1551526114),
	(7,'研发','2001,2100,2201,1001,1002,1006,1003,1004,1005,4001,4101',1551526207),
	(8,'测试','1001,1003,1004,1002,1006,1005',1551526226);



DROP TABLE IF EXISTS `syd_user_token`;

CREATE TABLE `syd_user_token` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned NOT NULL DEFAULT '0',
  `token` varchar(100) NOT NULL DEFAULT '',
  `expire` int(11) unsigned NOT NULL DEFAULT '0',
  `ctime` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `syd_user_token` (`id`, `user_id`, `token`, `expire`, `ctime`)
VALUES
	(11,28,'7GZKp693Z5lDFudmwWLR7znXxvsQkoOHEeNz5x31',1554126713,1551534713);
