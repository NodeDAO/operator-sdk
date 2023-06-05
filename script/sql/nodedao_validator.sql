/*
 Navicat Premium Data Transfer

 Source Server         : local
 Source Server Type    : MySQL
 Source Server Version : 80030
 Source Host           : localhost:3306
 Source Schema         : operator-sdk

 Target Server Type    : MySQL
 Target Server Version : 80030
 File Encoding         : 65001

 Date: 05/06/2023 19:05:24
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for nodedao_validator
-- ----------------------------
DROP TABLE IF EXISTS `nodedao_validator`;
CREATE TABLE `nodedao_validator` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `network` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'network eg:mainnet,goerli',
  `pubkey` varchar(650) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'Validator pubkey',
  `operator_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT 'Nodedao operator ID',
  `token_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT 'vNFT tokenId corresponding to Validator',
  `type` tinyint unsigned NOT NULL DEFAULT '0' COMMENT '0 user(default); 1 liquidStaking',
  `is_exit` tinyint unsigned NOT NULL DEFAULT '0' COMMENT 'Whether to exit. 0 not exit (default); 1 exited.',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'modify time',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='NodeDAO Operator’s validator exit scan.';

SET FOREIGN_KEY_CHECKS = 1;
