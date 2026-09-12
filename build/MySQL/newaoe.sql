-- MySQL dump 10.13  Distrib 8.0.44, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: newaoe
-- ------------------------------------------------------
-- Server version	8.0.41-0ubuntu0.22.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

CREATE USER IF NOT EXISTS 'root'@'%' IDENTIFIED BY 'newaoe';
GRANT ALL PRIVILEGES ON newaoe.* TO 'root'@'%';
FLUSH PRIVILEGES;

--
-- Table structure for table `CodeRun`
--
CREATE DATABASE IF NOT EXISTS newaoe;
USE newaoe;
DROP TABLE IF EXISTS `CodeRun`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CodeRun` (
  `indices` int NOT NULL AUTO_INCREMENT COMMENT '记录索引（自增主键）',
  `id` varchar(20) NOT NULL COMMENT '用户学号（关联用户表）',
  `submittime` datetime NOT NULL COMMENT '提交时间',
  `header` varchar(255) DEFAULT NULL COMMENT '头文件存储路径（如 private/2023001/header.h）',
  `source` varchar(255) DEFAULT NULL COMMENT '源文件存储路径（如 private/2023001/main.cpp）',
  `class` tinyint DEFAULT '0' COMMENT '代码类型（1=C++ 2=Python 3=Go 等）',
  `status` text COMMENT '运行状态',
  `description` text NOT NULL,
  `version` bigint NOT NULL DEFAULT '0' COMMENT '乐观锁版本号，防止更新乱序覆盖',
  PRIMARY KEY (`indices`),
  KEY `idx_id` (`id`),
  KEY `idx_submittime` (`submittime`)
) ENGINE=InnoDB AUTO_INCREMENT=176 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `CodeSubmitAssessment`
--

DROP TABLE IF EXISTS `CodeSubmitAssessment`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CodeSubmitAssessment` (
  `indices` int NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `id` varchar(20) NOT NULL COMMENT '关联用户ID',
  `uploadtime` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
  `header` text COMMENT '文件头',
  `source` text COMMENT '源代码',
  `headersize` int DEFAULT '0' COMMENT 'header字节数',
  `sourcesize` int DEFAULT '0' COMMENT 'source字节数',
  `teacher` varchar(20) NOT NULL COMMENT '提交的学生对应的教师姓名',
  PRIMARY KEY (`indices`),
  KEY `id` (`id`),
  KEY `teacher` (`teacher`),
  CONSTRAINT `CodeSubmitAssessment_ibfk_1` FOREIGN KEY (`id`) REFERENCES `Student` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `CodeSubmitAssessment_ibfk_2` FOREIGN KEY (`teacher`) REFERENCES `Teacher` (`name`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='代码提交评测表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `CodeSubmitCommon`
--

DROP TABLE IF EXISTS `CodeSubmitCommon`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `CodeSubmitCommon` (
  `indices` int NOT NULL AUTO_INCREMENT COMMENT '记录索引（自增主键）',
  `id` varchar(20) NOT NULL COMMENT '用户学号（外键，关联 StudentAllow 表）',
  `uploadtime` datetime NOT NULL COMMENT '上传时间',
  `header` varchar(255) DEFAULT NULL COMMENT '头文件存储路径（如 private/2023001/header.h）',
  `source` varchar(255) DEFAULT NULL COMMENT '源文件存储路径（如 private/2023001/main.cpp）',
  `description` text COMMENT '代码描述（支持长文本）',
  `headersize` bigint DEFAULT '0' COMMENT '头文件大小（字节）',
  `sourcesize` bigint DEFAULT '0' COMMENT '源文件大小（字节）',
  PRIMARY KEY (`indices`),
  KEY `id` (`id`),
  CONSTRAINT `FK_Code1_Student` FOREIGN KEY (`id`) REFERENCES `Student` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `FK_Code_Student` FOREIGN KEY (`id`) REFERENCES `Student` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=42 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='用户代码文件信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Student`
--

DROP TABLE IF EXISTS `Student`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Student` (
  `id` varchar(20) NOT NULL,
  `password` text,
  `email` varchar(50) DEFAULT NULL,
  `avatar` varchar(255) NOT NULL DEFAULT 'avatar/default.png',
  `registDate` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Teacher`
--

DROP TABLE IF EXISTS `Teacher`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `Teacher` (
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-04-19 13:58:25

