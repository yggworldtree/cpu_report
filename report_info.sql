/*
 Navicat Premium Data Transfer

 Source Server         : N550JK-vpn
 Source Server Type    : MySQL
 Source Server Version : 50732
 Source Host           : linux.ydtjxw.com:13306
 Source Schema         : cpudb

 Target Server Type    : MySQL
 Target Server Version : 50732
 File Encoding         : 65001

 Date: 21/06/2021 14:40:47
*/
-- ----------------------------
-- Table structure for report_info
-- ----------------------------
DROP TABLE IF EXISTS `report_info`;
CREATE TABLE `report_info` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `create_time` datetime(0) NULL DEFAULT NULL,
  `update_time` datetime(0) NULL DEFAULT NULL,
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `start_time` datetime(0) NULL DEFAULT NULL,
  `end_time` datetime(0) NULL DEFAULT NULL,
  `infos` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `warn_len` int(11) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
);
-- ----------------------------
-- Table structure for report_log
-- ----------------------------
DROP TABLE IF EXISTS `report_log`;
CREATE TABLE `report_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT 'primary key',
  `create_time` datetime(0) NULL DEFAULT NULL COMMENT 'create time',
  `update_time` datetime(0) NULL DEFAULT NULL COMMENT 'update time',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT 'dev name',
  `cpu_avg` double NULL DEFAULT NULL,
  `proc_len` int(11) NULL DEFAULT NULL,
  `mem_total` bigint(20) NULL DEFAULT NULL,
  `swap_total` bigint(20) NULL DEFAULT NULL,
  `mem_used` bigint(20) NULL DEFAULT NULL,
  `swap_used` bigint(20) NULL DEFAULT NULL,
  `mem_per` double NULL DEFAULT NULL,
  `swap_per` double NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
);
-- ----------------------------
-- Table structure for report_param
-- ----------------------------
DROP TABLE IF EXISTS `report_param`;
CREATE TABLE `report_param` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `data` longblob NULL,
  `crate_time` datetime(0) NULL DEFAULT NULL,
  `update_time` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
);
-- ----------------------------
-- Table structure for report_warn
-- ----------------------------
DROP TABLE IF EXISTS `report_warn`;
CREATE TABLE `report_warn` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `gid` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `create_time` datetime(0) NULL DEFAULT NULL,
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `val` double NULL DEFAULT NULL,
  `wval` double NULL DEFAULT NULL,
  `lev` int(11) NULL DEFAULT NULL,
  `warns` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
);