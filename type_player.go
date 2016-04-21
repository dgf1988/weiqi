package weiqi

import (
	"time"
)

/*
CREATE TABLE `player` (
	`id` INT(11) NOT NULL AUTO_INCREMENT,
	`pname` CHAR(50) NOT NULL,
	`psex` INT(11) NOT NULL DEFAULT '0',
	`pcountry` VARCHAR(20) NOT NULL DEFAULT '',
	`prank` VARCHAR(10) NOT NULL DEFAULT '',
	`pbirth` DATE NOT NULL DEFAULT '0000-00-00',
	PRIMARY KEY (`id`)
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB
;
*/

type Player struct {
	Id      int64
	Name    string
	Sex     int64
	Country string
	Rank    string
	Birth   time.Time
}

func (p Player) StrSex() string {
	return sexToChinese(p.Sex)
}

const (
	C_SEX_BOY  = 1
	C_SEX_GIRL = 2
)

func sexToChinese(sex int64) string {
	switch sex {
	case 1:
		return "男"
	case 2:
		return "女"
	default:
		return ""
	}
}

func chineseToSex(sex string) int64 {
	switch sex {
	case "男":
		return 1
	case "女":
		return 2
	default:
		return 0
	}
}