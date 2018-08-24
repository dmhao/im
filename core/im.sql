/*
Navicat MySQL Data Transfer

Source Server         : local_playyx
Source Server Version : 50617
Source Host           : localhost:3306
Source Database       : im

Target Server Type    : MYSQL
Target Server Version : 50617
File Encoding         : 65001

Date: 2018-08-24 14:09:10
*/
CREATE database im default character set utf8mb4 collate utf8mb4_unicode_ci;
USE im;
SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `apps`
-- ----------------------------
DROP TABLE IF EXISTS `apps`;
CREATE TABLE `apps` (
  `app_id` int(11) NOT NULL AUTO_INCREMENT,
  `secret_id` char(100) CHARACTER SET utf8 NOT NULL DEFAULT '',
  `secret_key` char(100) CHARACTER SET utf8 NOT NULL DEFAULT '',
  `app_name` varchar(200) CHARACTER SET utf8 NOT NULL DEFAULT '',
  `create_time` bigint(20) NOT NULL DEFAULT '0',
  `status` tinyint(4) NOT NULL DEFAULT '1',
  PRIMARY KEY (`app_id`),
  KEY `secret_id` (`secret_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of apps
-- ----------------------------
INSERT INTO `apps` VALUES ('1', 'asddsa', 'asddsa', 'test', '0', '1');

-- ----------------------------
-- Table structure for `groups`
-- ----------------------------
DROP TABLE IF EXISTS `groups`;
CREATE TABLE `groups` (
  `group_id` int(11) NOT NULL AUTO_INCREMENT,
  `app_id` int(11) NOT NULL DEFAULT '0',
  `group_name` char(50) NOT NULL DEFAULT '',
  `group_des` varchar(500) NOT NULL DEFAULT '',
  `group_icon` varchar(300) CHARACTER SET utf8 NOT NULL DEFAULT '',
  `master_user_id` bigint(20) NOT NULL DEFAULT '0',
  `user_count` smallint(6) NOT NULL DEFAULT '1',
  `max_user_count` smallint(6) NOT NULL,
  `join_need_examine` tinyint(4) NOT NULL DEFAULT '1',
  `create_time` bigint(11) NOT NULL DEFAULT '0',
  `update_time` bigint(11) NOT NULL DEFAULT '0',
  `status` tinyint(4) NOT NULL DEFAULT '1',
  PRIMARY KEY (`group_id`),
  KEY `app_id` (`app_id`),
  KEY `master_user_id` (`master_user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of groups
-- ----------------------------
INSERT INTO `groups` VALUES ('1', '1', '测试群组1', '测试群组', 'xxxxxxxxxx', '1', '0', '200', '1', '1532333645', '1532335878', '1');

-- ----------------------------
-- Table structure for `group_examine_users`
-- ----------------------------
DROP TABLE IF EXISTS `group_examine_users`;
CREATE TABLE `group_examine_users` (
  `examine_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `app_id` int(11) NOT NULL DEFAULT '0',
  `group_id` int(11) NOT NULL,
  `group_name` char(50) CHARACTER SET utf8 NOT NULL DEFAULT '',
  `user_id` bigint(20) NOT NULL DEFAULT '0',
  `user_name` char(50) CHARACTER SET utf8 NOT NULL DEFAULT '',
  `create_time` bigint(20) NOT NULL DEFAULT '0',
  `examine_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0申请入群   1邀请入群',
  `status` tinyint(4) NOT NULL DEFAULT '0' COMMENT '-1拒绝  1允许',
  `op_user_id` bigint(20) NOT NULL DEFAULT '0',
  `examine_time` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`examine_id`),
  KEY `app_id` (`app_id`),
  KEY `group_id` (`group_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of group_examine_users
-- ----------------------------

-- ----------------------------
-- Table structure for `group_users`
-- ----------------------------
DROP TABLE IF EXISTS `group_users`;
CREATE TABLE `group_users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `app_id` int(11) NOT NULL DEFAULT '0',
  `group_id` int(11) NOT NULL DEFAULT '0',
  `user_id` bigint(20) NOT NULL DEFAULT '0',
  `user_role` tinyint(4) NOT NULL DEFAULT '0' COMMENT '0 普通用户  1群组  2管理员',
  `join_time` bigint(20) NOT NULL DEFAULT '0',
  `update_time` bigint(20) NOT NULL DEFAULT '0',
  `status` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1有效  -1删除  0群组解散',
  PRIMARY KEY (`id`),
  KEY `app_id` (`app_id`),
  KEY `group_id` (`group_id`),
  KEY `user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of group_users
-- ----------------------------
INSERT INTO `group_users` VALUES ('1', '1', '1', '1', '1', '1532333645', '0', '1');
INSERT INTO `group_users` VALUES ('5', '1', '1', '2', '2', '1532336654', '0', '1');
INSERT INTO `group_users` VALUES ('6', '1', '1', '3', '0', '1532685867', '0', '1');
INSERT INTO `group_users` VALUES ('7', '1', '1', '4', '0', '1532685867', '0', '1');
INSERT INTO `group_users` VALUES ('8', '1', '1', '5', '0', '1532687217', '0', '1');
INSERT INTO `group_users` VALUES ('9', '1', '1', '6', '0', '1532687217', '0', '1');
INSERT INTO `group_users` VALUES ('10', '1', '1', '13', '0', '1533608496', '0', '1');

-- ----------------------------
-- Table structure for `messages`
-- ----------------------------
DROP TABLE IF EXISTS `messages`;
CREATE TABLE `messages` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `sender_id` bigint(20) NOT NULL DEFAULT '0',
  `receiver_id` bigint(20) NOT NULL DEFAULT '0',
  `chart_type` tinyint(4) NOT NULL DEFAULT '0',
  `msg_type` tinyint(4) NOT NULL DEFAULT '0',
  `msg_id` char(50) CHARACTER SET utf8 NOT NULL DEFAULT '',
  `talk_id` char(50) CHARACTER SET utf8 NOT NULL DEFAULT '',
  `trace_id` char(50) CHARACTER SET utf8 NOT NULL DEFAULT '',
  `timestamp` bigint(20) NOT NULL DEFAULT '0',
  `content` varchar(5000) NOT NULL DEFAULT '',
  `app_id` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of messages
-- ----------------------------

-- ----------------------------
-- Table structure for `route_services`
-- ----------------------------
DROP TABLE IF EXISTS `route_services`;
CREATE TABLE `route_services` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `app_id` int(11) NOT NULL DEFAULT '0',
  `addr` char(20) CHARACTER SET utf8 NOT NULL DEFAULT '',
  `status` tinyint(4) NOT NULL DEFAULT '1',
  `create_time` bigint(20) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `app_id` (`app_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of route_services
-- ----------------------------
INSERT INTO `route_services` VALUES ('6', '0', '127.0.0.1:3333', '1', '0');
INSERT INTO `route_services` VALUES ('7', '1', '127.0.0.1:3333', '1', '0');
