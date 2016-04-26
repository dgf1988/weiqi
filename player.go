package weiqi

import (
	"html/template"
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

var (
	rankMap = map[int]string{
		109: "九段",
		108: "八段",
		107: "七段",
		106: "六段",
		105: "五段",
		104: "四段",
		103: "三段",
		102: "二段",
		101: "一段",
	}
)

func rankToChinese(rank int) string {
	ch, ok := rankMap[rank]
	if ok {
		return ch
	} else {
		return ""
	}
}

func chineseToRank( chRank string) int {
	switch chRank {
	case "九段":
		return 109
	case "八段":
		return 108
	case "七段":
		return 107
	case "六段":
		return 106
	case "五段":
		return 105
	case "四段":
		return 104
	case "三段":
		return 103
	case "二段":
		return 102
	case "一段":
		return 101
	default:
		return 0
	}
}

type Player struct {
	Id      int64
	Name    string
	Sex     int64
	Country string
	Rank    int
	Birth   time.Time
}

func (t Text) HtmlText() template.HTML {
	return template.HTML(t.Text)
}

func (p Player) StrSex() string {
	return sexToChinese(p.Sex)
}

func (p Player) StrRank() string {
	return rankToChinese(p.Rank)
}

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

func listPlayerOrderByRankDesc(take, skip int) ([]Player, error) {
	var players []Player
	if rows, err := Db.Player.Query("order by player.rank desc, player.birth desc limit ?, ?", skip, take); err != nil {
		return nil, err
	} else {
		defer rows.Close()
		for rows.Next() {
			var player Player
			if err = rows.Struct(&player); err != nil {
				return nil, err
			} else {
				players = append(players, player)
			}
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	return players, nil
}
