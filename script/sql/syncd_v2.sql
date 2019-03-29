# ************************************************************
# Copyright 2019 syncd Author. All Rights Reserved.
# Use of this source code is governed by a MIT-style
# license that can be found in the LICENSE file.
# ************************************************************

create database `syncd` default charset utf8mb4;

use `syncd`;

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
  PRIMARY KEY (`id`),
  KEY `idx_space_project` (`space_id`,`project_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `syd_deploy_build`;

CREATE TABLE `syd_deploy_build` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `apply_id` int(11) unsigned NOT NULL DEFAULT '0',
  `start_time` int(11) unsigned NOT NULL DEFAULT '0',
  `finish_time` int(11) unsigned NOT NULL DEFAULT '0',
  `status` int(11) unsigned NOT NULL DEFAULT '1',
  `tar` varchar(2000) NOT NULL DEFAULT '',
  `output` mediumtext NOT NULL,
  `errmsg` text NOT NULL,
  `ctime` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_apply_id` (`apply_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

DROP TABLE IF EXISTS `syd_deploy_task`;

CREATE TABLE `syd_deploy_task` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `apply_id` int(11) unsigned NOT NULL DEFAULT '0',
  `group_id` int(11) unsigned NOT NULL DEFAULT '0',
  `status` int(11) unsigned NOT NULL DEFAULT '0',
  `content` text NOT NULL,
  `ctime` int(11) unsigned NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_apply_id` (`apply_id`)
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
  `build_hook_script` text NOT NULL,
  `deploy_hook_script` text NOT NULL,
  `pre_deploy_cmd` text NOT NULL,
  `after_deploy_cmd` text NOT NULL,
  `audit_notice` varchar(2000) DEFAULT NULL,
  `deploy_notice` varchar(2000) DEFAULT NULL,
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
  PRIMARY KEY (`id`),
  KEY `idx_username` (`username`),
  KEY `idx_email` (`email`(20))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `syd_user` (`id`, `role_id`, `username`, `password`, `salt`, `truename`, `mobile`, `email`, `status`, `last_login_time`, `last_login_ip`, `ctime`)
VALUES
	(1,1,'syncd','c2a8d572366f7cf7bfc8b485f41bfd1b','u0EMxuE6qh','Syncd','','',1,0,'',0);

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
	(1,'管理员','2001,2002,2003,2004,2100,2101,2102,2201,2202,2203,2204,2205,2206,2207,3001,3002,3004,3003,3101,3102,3103,3104,4001,4002,4003,4004,4101,4102,4103,4104,1001,1002,1006,1003,1004,1005',0);

DROP TABLE IF EXISTS `syd_user_token`;

CREATE TABLE `syd_user_token` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) unsigned NOT NULL DEFAULT '0',
  `token` varchar(100) NOT NULL DEFAULT '',
  `expire` int(11) unsigned NOT NULL DEFAULT '0',
  `ctime` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;