/*
 Navicat Premium Data Transfer

 Source Server         : 127.0.0.1
 Source Server Type    : MySQL
 Source Server Version : 80018
 Source Host           : 127.0.0.1:3306
 Source Schema         : bee_custom

 Target Server Type    : MySQL
 Target Server Version : 80018
 File Encoding         : 65001

 Date: 19/11/2019 17:08:54
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
) ENGINE = InnoDB AUTO_INCREMENT = 47 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

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
INSERT INTO `bee_custom_resource` VALUES (16, '2019-11-18 01:05:08', '2019-11-18 02:01:22', 0, '客户管理', '', 'CompanyController.Index', NULL);
INSERT INTO `bee_custom_resource` VALUES (17, '2019-11-18 01:05:27', '2019-11-18 01:57:03', 0, '编辑', '', 'CompanyController.Edit', 16);
INSERT INTO `bee_custom_resource` VALUES (18, '2019-11-18 01:05:40', '2019-11-18 01:56:55', 0, '新建', '', 'CompanyController.Create', 16);
INSERT INTO `bee_custom_resource` VALUES (19, '2019-11-18 01:05:55', '2019-11-18 01:56:48', 0, '删除', '', 'CompanyController.Delete', 16);
INSERT INTO `bee_custom_resource` VALUES (20, '2019-11-18 02:02:10', '2019-11-18 02:02:10', 0, '通关参数', '', 'ClearanceController.Index', NULL);
INSERT INTO `bee_custom_resource` VALUES (21, '2019-11-18 02:15:15', '2019-11-18 02:15:15', 0, '新建', '', 'ClearanceController.Create', 20);
INSERT INTO `bee_custom_resource` VALUES (22, '2019-11-18 02:15:25', '2019-11-18 02:15:25', 0, '编辑', '', 'ClearanceController.Edit', 20);
INSERT INTO `bee_custom_resource` VALUES (23, '2019-11-18 02:15:36', '2019-11-18 02:16:04', 0, '删除', '', 'ClearanceController.Delete', 20);
INSERT INTO `bee_custom_resource` VALUES (24, '2019-11-18 02:15:50', '2019-11-18 02:15:50', 0, '导入', '', 'ClearanceController.Import', 20);
INSERT INTO `bee_custom_resource` VALUES (25, '2019-11-18 02:16:59', '2019-11-18 02:16:59', 0, '商检编码管理', '', 'CiqController.Index', NULL);
INSERT INTO `bee_custom_resource` VALUES (26, '2019-11-18 02:17:07', '2019-11-18 02:17:50', 0, '商品编码管理', '', 'HsCodeController.Index', NULL);
INSERT INTO `bee_custom_resource` VALUES (27, '2019-11-18 02:17:38', '2019-11-18 02:17:38', 0, '金关二期手账册管理', '', 'HandBookController.Index', NULL);
INSERT INTO `bee_custom_resource` VALUES (28, '2019-11-18 02:18:53', '2019-11-18 02:18:53', 0, '导入', '', 'CiqController.Import', 25);
INSERT INTO `bee_custom_resource` VALUES (29, '2019-11-18 02:19:01', '2019-11-18 02:19:01', 0, '导入', '', 'HsCodeController.Import', 26);
INSERT INTO `bee_custom_resource` VALUES (30, '2019-11-18 02:19:15', '2019-11-18 02:19:15', 0, '导入', '', 'HandBookController.Import', 27);
INSERT INTO `bee_custom_resource` VALUES (31, '2019-11-18 02:21:05', '2019-11-18 02:31:58', 0, '进口清单管理', '', 'AnnotationController.IIndex', NULL);
INSERT INTO `bee_custom_resource` VALUES (32, '2019-11-18 02:32:21', '2019-11-18 02:38:32', 0, '代客下单', '', 'AnnotationController.ICreate', 31);
INSERT INTO `bee_custom_resource` VALUES (33, '2019-11-18 02:32:42', '2019-11-18 02:34:28', 0, '开始审核', '', 'AnnotationController.IEdit', 31);
INSERT INTO `bee_custom_resource` VALUES (34, '2019-11-18 02:33:00', '2019-11-18 02:34:34', 0, '开始制单', '', 'AnnotationController.IMake', 31);
INSERT INTO `bee_custom_resource` VALUES (35, '2019-11-18 02:34:51', '2019-11-18 02:34:51', 0, '取消', '', 'AnnotationController.ICancel', 31);
INSERT INTO `bee_custom_resource` VALUES (36, '2019-11-18 02:35:06', '2019-11-18 02:35:06', 0, '审核通过', '', 'AnnotationController.IAudit', 31);
INSERT INTO `bee_custom_resource` VALUES (37, '2019-11-18 02:35:20', '2019-11-18 02:35:20', 0, '派单', '', 'AnnotationController.IDistribute', 31);
INSERT INTO `bee_custom_resource` VALUES (38, '2019-11-18 02:35:35', '2019-11-18 02:35:35', 0, '删除', '', 'AnnotationController.IDelete', 31);
INSERT INTO `bee_custom_resource` VALUES (39, '2019-11-18 02:21:05', '2019-11-18 02:31:58', 0, '出口清单管理', '', 'AnnotationController.EIndex', NULL);
INSERT INTO `bee_custom_resource` VALUES (40, '2019-11-18 02:32:21', '2019-11-18 02:38:32', 0, '代客下单', '', 'AnnotationController.ECreate', 39);
INSERT INTO `bee_custom_resource` VALUES (41, '2019-11-18 02:32:42', '2019-11-18 02:34:28', 0, '开始审核', '', 'AnnotationController.EEdit', 39);
INSERT INTO `bee_custom_resource` VALUES (42, '2019-11-18 02:33:00', '2019-11-18 02:34:34', 0, '开始制单', '', 'AnnotationController.EMake', 39);
INSERT INTO `bee_custom_resource` VALUES (43, '2019-11-18 02:34:51', '2019-11-18 02:34:51', 0, '取消', '', 'AnnotationController.ECancel', 39);
INSERT INTO `bee_custom_resource` VALUES (44, '2019-11-18 02:35:06', '2019-11-18 02:35:06', 0, '审核通过', '', 'AnnotationController.EAudit', 39);
INSERT INTO `bee_custom_resource` VALUES (45, '2019-11-18 02:35:20', '2019-11-18 02:35:20', 0, '派单', '', 'AnnotationController.EDistribute', 39);
INSERT INTO `bee_custom_resource` VALUES (46, '2019-11-18 02:35:35', '2019-11-18 02:35:35', 0, '删除', '', 'AnnotationController.EDelete', 39);

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
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of bee_custom_roles
-- ----------------------------
INSERT INTO `bee_custom_roles` VALUES (1, '2019-10-26 17:30:57', '2019-11-19 04:44:28', '超级管理员');
INSERT INTO `bee_custom_roles` VALUES (2, '2019-10-29 02:33:54', '2019-11-19 04:46:32', 'fdgdfgd');

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
  `is_super` tinyint(1) NOT NULL DEFAULT 0,
  `status` tinyint(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 3 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of bee_custom_users
-- ----------------------------
INSERT INTO `bee_custom_users` VALUES (1, '2019-11-19 11:56:09', '2019-11-19 03:30:21', 'lihaitao', 'admin', 'e10adc3949ba59abbe56e057f20f883e', '18612348765', 'lhtzbj18@126.com', '/static/upload/1.jpg', '18612348765', '', 1, 1);
INSERT INTO `bee_custom_users` VALUES (2, '2019-10-28 10:13:40', '2019-11-19 03:30:14', 'fsdfds', 'sdfdsf', 'a75d4930656ba67e7761ba235d69df25', '13577777773', 'sdfsdf@dfsd.com', '', 'sdfsdfsd', '', 1, 1);

-- ----------------------------
-- Table structure for casbin_rule
-- ----------------------------
DROP TABLE IF EXISTS `casbin_rule`;
CREATE TABLE `casbin_rule`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `p_type` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `v0` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `v1` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `v2` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `v3` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `v4` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `v5` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 191 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of casbin_rule
-- ----------------------------
INSERT INTO `casbin_rule` VALUES (1, 'g', 'fsdfds', '超级管理员', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (2, 'g', '1', '1,2', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (3, 'g', '1', '1,2', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (4, 'g', '1', '1,2', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (5, 'g', '2', '1,2', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (6, 'g', '2', '1', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (7, 'g', '2', '2', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (8, 'g', '1', '1', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (9, 'g', '1', '2', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (101, 'p', '1', 'ResourceController.Index', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (102, 'p', '1', 'ResourceController.Create', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (103, 'p', '1', 'ResourceController.Edit', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (104, 'p', '1', 'ResourceController.Delete', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (105, 'p', '1', 'RoleController.Index', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (106, 'p', '1', 'RoleController.Create', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (107, 'p', '1', 'RoleController.Edit', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (108, 'p', '1', 'RoleController.Delete', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (109, 'p', '1', 'BackendUserController.Index', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (110, 'p', '1', 'BackendUserController.Create', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (111, 'p', '1', 'BackendUserController.Edit', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (112, 'p', '1', 'BackendUserController.Delete', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (113, 'p', '1', 'BackendUserController.Freeze', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (114, 'p', '1', 'BackendUserController.Profile', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (115, 'p', '1', 'CompanyController.Index', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (116, 'p', '1', 'CompanyController.Edit', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (117, 'p', '1', 'CompanyController.Create', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (118, 'p', '1', 'CompanyController.Delete', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (119, 'p', '1', 'ClearanceController.Index', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (120, 'p', '1', 'ClearanceController.Create', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (121, 'p', '1', 'ClearanceController.Edit', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (122, 'p', '1', 'ClearanceController.Delete', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (123, 'p', '1', 'ClearanceController.Import', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (124, 'p', '1', 'CiqController.Index', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (125, 'p', '1', 'CiqController.Import', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (126, 'p', '1', 'HsCodeController.Index', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (127, 'p', '1', 'HsCodeController.Import', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (128, 'p', '1', 'HandBookController.Index', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (129, 'p', '1', 'HandBookController.Import', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (130, 'p', '1', 'AnnotationController.IIndex', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (131, 'p', '1', 'AnnotationController.ICreate', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (132, 'p', '1', 'AnnotationController.IEdit', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (133, 'p', '1', 'AnnotationController.IMake', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (134, 'p', '1', 'AnnotationController.ICancel', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (135, 'p', '1', 'AnnotationController.IAudit', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (136, 'p', '1', 'AnnotationController.IDistribute', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (137, 'p', '1', 'AnnotationController.IDelete', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (138, 'p', '1', 'AnnotationController.EIndex', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (139, 'p', '1', 'AnnotationController.ECreate', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (140, 'p', '1', 'AnnotationController.EEdit', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (141, 'p', '1', 'AnnotationController.EMake', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (142, 'p', '1', 'AnnotationController.ECancel', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (143, 'p', '1', 'AnnotationController.EAudit', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (144, 'p', '1', 'AnnotationController.EDistribute', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (145, 'p', '1', 'AnnotationController.EDelete', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (146, 'p', '2', 'ResourceController.Index', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (147, 'p', '2', 'ResourceController.Create', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (148, 'p', '2', 'ResourceController.Edit', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (149, 'p', '2', 'ResourceController.Delete', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (150, 'p', '2', 'RoleController.Index', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (151, 'p', '2', 'RoleController.Create', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (152, 'p', '2', 'RoleController.Edit', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (153, 'p', '2', 'RoleController.Delete', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (154, 'p', '2', 'BackendUserController.Index', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (155, 'p', '2', 'BackendUserController.Create', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (156, 'p', '2', 'BackendUserController.Edit', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (157, 'p', '2', 'BackendUserController.Delete', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (158, 'p', '2', 'BackendUserController.Freeze', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (159, 'p', '2', 'BackendUserController.Profile', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (160, 'p', '2', 'CompanyController.Index', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (161, 'p', '2', 'CompanyController.Edit', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (162, 'p', '2', 'CompanyController.Create', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (163, 'p', '2', 'CompanyController.Delete', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (164, 'p', '2', 'ClearanceController.Index', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (165, 'p', '2', 'ClearanceController.Create', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (166, 'p', '2', 'ClearanceController.Edit', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (167, 'p', '2', 'ClearanceController.Delete', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (168, 'p', '2', 'ClearanceController.Import', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (169, 'p', '2', 'CiqController.Index', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (170, 'p', '2', 'CiqController.Import', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (171, 'p', '2', 'HsCodeController.Index', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (172, 'p', '2', 'HsCodeController.Import', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (173, 'p', '2', 'HandBookController.Index', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (174, 'p', '2', 'HandBookController.Import', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (175, 'p', '2', 'AnnotationController.IIndex', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (176, 'p', '2', 'AnnotationController.ICreate', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (177, 'p', '2', 'AnnotationController.IEdit', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (178, 'p', '2', 'AnnotationController.IMake', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (179, 'p', '2', 'AnnotationController.ICancel', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (180, 'p', '2', 'AnnotationController.IAudit', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (181, 'p', '2', 'AnnotationController.IDistribute', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (182, 'p', '2', 'AnnotationController.IDelete', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (183, 'p', '2', 'AnnotationController.EIndex', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (184, 'p', '2', 'AnnotationController.ECreate', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (185, 'p', '2', 'AnnotationController.EEdit', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (186, 'p', '2', 'AnnotationController.EMake', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (187, 'p', '2', 'AnnotationController.ECancel', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (188, 'p', '2', 'AnnotationController.EAudit', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (189, 'p', '2', 'AnnotationController.EDistribute', '', '', '', '');
INSERT INTO `casbin_rule` VALUES (190, 'p', '2', 'AnnotationController.EDelete', '', '', '', '');

-- ----------------------------
-- Table structure for session
-- ----------------------------
DROP TABLE IF EXISTS `session`;
CREATE TABLE `session`  (
  `session_key` char(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL,
  `session_data` blob NULL,
  `session_expiry` int(11) UNSIGNED NOT NULL,
  PRIMARY KEY (`session_key`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
