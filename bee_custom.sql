/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 50728
 Source Host           : localhost:3306
 Source Schema         : bee_custom

 Target Server Type    : MySQL
 Target Server Version : 50728
 File Encoding         : 65001

 Date: 26/10/2019 18:11:17
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for bee_custom_resource
-- ----------------------------
DROP TABLE IF EXISTS `bee_custom_resource`;
CREATE TABLE `bee_custom_resource`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(0) NULL DEFAULT NULL,
  `updated_at` datetime(0) NULL DEFAULT NULL,
  `rtype` int(11) NOT NULL DEFAULT 0,
  `name` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `icon` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `url_for` varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `parent_id` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 16 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of bee_custom_resource
-- ----------------------------
INSERT INTO `bee_custom_resource` VALUES (1, '2019-10-26 17:31:08', '2019-10-26 17:31:10', 1, '权限管理', '', 'ResourceController.Index', NULL);
INSERT INTO `bee_custom_resource` VALUES (2, '2019-10-26 17:32:08', '2019-10-26 17:32:06', 1, '新建', '', 'ResourceController.Create', 1);
INSERT INTO `bee_custom_resource` VALUES (3, '2019-10-26 17:32:15', '2019-10-26 17:32:17', 1, '编辑', '', 'ResourceController.Edit', 1);
INSERT INTO `bee_custom_resource` VALUES (4, '2019-10-26 17:32:15', '2019-10-26 17:32:15', 1, '删除', '', 'ResourceController.Delete', 1);
INSERT INTO `bee_custom_resource` VALUES (5, '2019-10-26 17:33:43', '2019-10-26 17:33:44', 1, '角色管理', '', 'RoleController.Index', NULL);
INSERT INTO `bee_custom_resource` VALUES (7, '2019-10-26 17:32:08', '2019-10-26 17:32:06', 1, '新建', '', 'RoleController.Create', 5);
INSERT INTO `bee_custom_resource` VALUES (8, '2019-10-26 17:32:15', '2019-10-26 17:32:17', 1, '编辑', '', 'RoleController.Edit', 5);
INSERT INTO `bee_custom_resource` VALUES (9, '2019-10-26 17:32:15', '2019-10-26 17:32:15', 1, '删除', '', 'RoleController.Delete', 5);
INSERT INTO `bee_custom_resource` VALUES (10, '2019-10-26 17:33:43', '2019-10-26 17:33:44', 1, '用户管理', '', 'BackendUserController.Index', NULL);
INSERT INTO `bee_custom_resource` VALUES (11, '2019-10-26 17:32:08', '2019-10-26 17:32:06', 1, '新建', '', 'BackendUserController.Create', 10);
INSERT INTO `bee_custom_resource` VALUES (12, '2019-10-26 17:32:15', '2019-10-26 17:32:17', 1, '编辑', '', 'BackendUserController.Edit', 10);
INSERT INTO `bee_custom_resource` VALUES (13, '2019-10-26 17:32:15', '2019-10-26 17:32:15', 1, '删除', '', 'BackendUserController.Delete', 10);
INSERT INTO `bee_custom_resource` VALUES (14, '2019-10-26 17:32:15', '2019-10-26 17:32:15', 1, '启用/禁用', '', 'BackendUserController.Freeze', 10);
INSERT INTO `bee_custom_resource` VALUES (15, '2019-10-26 17:32:15', '2019-10-26 17:32:15', 0, '我的信息', '', 'BackendUserController.Profile', NULL);

-- ----------------------------
-- Table structure for bee_custom_roles
-- ----------------------------
DROP TABLE IF EXISTS `bee_custom_roles`;
CREATE TABLE `bee_custom_roles`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(0) NULL DEFAULT NULL,
  `updated_at` datetime(0) NULL DEFAULT NULL,
  `name` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of bee_custom_roles
-- ----------------------------
INSERT INTO `bee_custom_roles` VALUES (1, '2019-10-26 17:30:57', '2019-10-26 17:30:59', '超级管理员');

-- ----------------------------
-- Table structure for bee_custom_roles_bee_custom_resources
-- ----------------------------
DROP TABLE IF EXISTS `bee_custom_roles_bee_custom_resources`;
CREATE TABLE `bee_custom_roles_bee_custom_resources`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `bee_custom_roles_id` bigint(20) NOT NULL,
  `bee_custom_resource_id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 19 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of bee_custom_roles_bee_custom_resources
-- ----------------------------
INSERT INTO `bee_custom_roles_bee_custom_resources` VALUES (5, 1, 1);
INSERT INTO `bee_custom_roles_bee_custom_resources` VALUES (6, 1, 2);
INSERT INTO `bee_custom_roles_bee_custom_resources` VALUES (7, 1, 3);
INSERT INTO `bee_custom_roles_bee_custom_resources` VALUES (8, 1, 4);
INSERT INTO `bee_custom_roles_bee_custom_resources` VALUES (9, 1, 5);
INSERT INTO `bee_custom_roles_bee_custom_resources` VALUES (10, 1, 7);
INSERT INTO `bee_custom_roles_bee_custom_resources` VALUES (11, 1, 8);
INSERT INTO `bee_custom_roles_bee_custom_resources` VALUES (12, 1, 9);
INSERT INTO `bee_custom_roles_bee_custom_resources` VALUES (13, 1, 10);
INSERT INTO `bee_custom_roles_bee_custom_resources` VALUES (14, 1, 11);
INSERT INTO `bee_custom_roles_bee_custom_resources` VALUES (15, 1, 12);
INSERT INTO `bee_custom_roles_bee_custom_resources` VALUES (16, 1, 13);
INSERT INTO `bee_custom_roles_bee_custom_resources` VALUES (17, 1, 14);
INSERT INTO `bee_custom_roles_bee_custom_resources` VALUES (18, 1, 15);

-- ----------------------------
-- Table structure for bee_custom_users
-- ----------------------------
DROP TABLE IF EXISTS `bee_custom_users`;
CREATE TABLE `bee_custom_users`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(0) NULL DEFAULT NULL,
  `updated_at` datetime(0) NULL DEFAULT NULL,
  `real_name` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `user_name` varchar(24) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `user_pwd` varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `mobile` varchar(16) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `email` varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `avatar` varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `i_c_code` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `chapter` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL,
  `role_id` bigint(20) NOT NULL,
  `is_super` tinyint(1) NOT NULL DEFAULT 0,
  `status` tinyint(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of bee_custom_users
-- ----------------------------
INSERT INTO `bee_custom_users` VALUES (1, NULL, NULL, 'lihaitao', 'admin', 'e10adc3949ba59abbe56e057f20f883e', '18612348765', 'lhtzbj18@126.com', '/static/upload/1.jpg', '18612348765', '', 1, 0, 0);

-- ----------------------------
-- Table structure for session
-- ----------------------------
DROP TABLE IF EXISTS `session`;
CREATE TABLE `session`  (
  `session_key` char(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `session_data` blob NULL,
  `session_expiry` int(11) UNSIGNED NOT NULL,
  PRIMARY KEY (`session_key`) USING BTREE
) ENGINE = MyISAM CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
