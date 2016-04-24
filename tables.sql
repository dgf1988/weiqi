-- --------------------------------------------------------
-- 主机:                           112.124.30.222
-- 服务器版本:                        10.1.13-MariaDB - MariaDB Server
-- 服务器操作系统:                      Linux
-- HeidiSQL 版本:                  9.1.0.4867
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;

-- 导出  表 weiqi_2.player 结构
CREATE TABLE IF NOT EXISTS `player` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` char(50) NOT NULL,
  `sex` int(11) NOT NULL DEFAULT '0',
  `country` varchar(20) NOT NULL DEFAULT '',
  `rank` varchar(10) NOT NULL DEFAULT '',
  `birth` date NOT NULL DEFAULT '0000-00-00',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。


-- 导出  表 weiqi_2.post 结构
CREATE TABLE IF NOT EXISTS `post` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(250) NOT NULL,
  `text` text NOT NULL,
  `status` int(11) NOT NULL DEFAULT '0',
  `posted` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。


-- 导出  表 weiqi_2.sgf 结构
CREATE TABLE IF NOT EXISTS `sgf` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `time` date NOT NULL DEFAULT '0000-00-00',
  `place` char(50) NOT NULL DEFAULT '',
  `event` char(100) NOT NULL DEFAULT '',
  `black` char(50) NOT NULL DEFAULT 'b',
  `white` char(50) NOT NULL DEFAULT 'w',
  `rule` char(50) NOT NULL DEFAULT '',
  `result` char(50) NOT NULL DEFAULT '',
  `steps` mediumtext NOT NULL,
  `update` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。


-- 导出  表 weiqi_2.text 结构
CREATE TABLE IF NOT EXISTS `text` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `text` mediumtext NOT NULL,
  `status` int(11) NOT NULL DEFAULT '0',
  `create` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。


-- 导出  表 weiqi_2.text_player 结构
CREATE TABLE IF NOT EXISTS `text_player` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `playerid` int(11) NOT NULL,
  `textid` int(11) NOT NULL,
  `status` int(11) NOT NULL DEFAULT '0',
  `create` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。


-- 导出  表 weiqi_2.user 结构
CREATE TABLE IF NOT EXISTS `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` char(50) NOT NULL,
  `password` char(50) NOT NULL,
  `email` char(50) NOT NULL DEFAULT '',
  `ip` char(40) NOT NULL DEFAULT '127.0.0.1',
  `status` int(11) NOT NULL DEFAULT '0',
  `register` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uname` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 数据导出被取消选择。
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
