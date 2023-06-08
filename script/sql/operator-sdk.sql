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

 Date: 08/06/2023 10:35:47
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
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_network_token` (`network`,`token_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='NodeDAO Operator’s validator exit scan.';

-- ----------------------------
-- Table structure for neth_withdrawal_request
-- ----------------------------
DROP TABLE IF EXISTS `neth_withdrawal_request`;
CREATE TABLE `neth_withdrawal_request` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `network` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'network eg:mainnet,goerli',
  `operator_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT 'Nodedao operator ID',
  `request_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT 'withdrawal request id\n',
  `withdraw_height` bigint unsigned NOT NULL DEFAULT '0' COMMENT 'withdraw block height\n\n',
  `withdraw_neth_amount` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT 'neth amount\n',
  `withdraw_exchange` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT 'exchange: 1 neth = ? eth (wei)\n',
  `claim_eth_amount` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '0' COMMENT 'eth amount',
  `owner` char(42) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT 'Owner address',
  `is_exit` tinyint unsigned NOT NULL DEFAULT '0' COMMENT 'Whether to exit. 0 not exit (default); 1 exited.',
  `ctime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
  `mtime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'modify time',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_network_request` (`network`,`request_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='nETH Withdrawal Request record.';

SET FOREIGN_KEY_CHECKS = 1;
