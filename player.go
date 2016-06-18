package weiqi

import (
	"html/template"
	"time"
    "database/sql"
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
	rankMap = map[int64]string{
		109: "九段",
		108: "八段",
		107: "七段",
		106: "六段",
		105: "五段",
		104: "四段",
		103: "三段",
		102: "二段",
		101: "初段",
	}
)

func formatRank(rank int64) string {
	ch, ok := rankMap[rank]
	if ok {
		return ch
	} else {
		return ""
	}
}

func parseRank(chRank string) int64 {
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
	case "一段", "初段":
		return 101
	default:
		return 0
	}
}

type PlayerTable struct {
	Id      int64
	Name    string
	Sex     int64
	Country string
	Rank    int64
	Birth   time.Time
}

func (t Text) HtmlText() template.HTML {
	return template.HTML(t.Text)
}

func (p PlayerTable) StrSex() string {
	return formatSex(p.Sex)
}

func (p PlayerTable) StrRank() string {
	return formatRank(p.Rank)
}

func formatSex(sex int64) string {
	switch sex {
	case 1:
		return "男"
	case 2:
		return "女"
	default:
		return ""
	}
}

func parseSex(sex string) int64 {
	switch sex {
	case "男":
		return 1
	case "女":
		return 2
	default:
		return 0
	}
}

func listPlayerOrderByRankDesc(take, skip int) ([]PlayerTable, error) {
	var players []PlayerTable
	if rows, err := Db.Player.Query("order by player.rank desc, player.birth desc limit ?, ?", skip, take); err != nil {
		return nil, err
	} else {
		defer rows.Close()
		for rows.Next() {
			var player PlayerTable
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

type Player struct {
	Id int64
	Name string
	Sex string
	Country string
	Rank string
	Birth time.Time
	Text string
}

func (player *Player) StrBirth() string {
    return player.Birth.Format("2006年1月2日")
}

func (player *Player) HtmlText() template.HTML {
    return template.HTML(parseTextToHtml(player.Text))
}

func (player *Player) Save() (id int64, err error) {
    if player.Id > 0 {
        _, err = Db.Player.Update(player.Id).Values(nil, player.Name, parseSex(player.Sex), player.Country, parseRank(player.Rank), player.Birth)
        if err != nil {
            return
        }
        var textid int64
        var playertextid int64
        err = Db.PlayerText.Get(nil, player.Id).Scan(&playertextid, nil, &textid)
        if err == nil {
            err = Db.Text.Get(textid).Scan()
            if err == nil {
                _, err = Db.Text.Update(textid).Values(nil, player.Text)
                if err != nil {
                    return
                }
            } else if err == sql.ErrNoRows {
                textid, err = Db.Text.Add(nil, player.Text)
                if err != nil {
                    return
                }
                _, err = Db.PlayerText.Update(playertextid).Values(nil, nil, textid)
                if err != nil {
                    return
                }
            } else {
                return
            }
        } else if err == sql.ErrNoRows {
            textid, err = Db.Text.Add(nil, player.Text)
            if err != nil {
                return
            }
            _, err = Db.PlayerText.Add(nil, player.Id, textid)
            if err != nil {
                return
            }
        } else {
            return
        }

    } else {
        err = Db.Player.Get(nil, player.Name).Scan(&id)
        if err == nil {
            player.Id = id
            return player.Save()
        } else if err != sql.ErrNoRows {
            return
        }
        id, err = Db.Player.Add(nil, player.Name, parseSex(player.Sex), player.Country, parseRank(player.Rank), player.Birth)
        if err != nil {
            return
        }

        var textid int64
        textid, err = Db.Text.Add(nil, player.Text)
        if err != nil {
            return
        }
        _, err = Db.PlayerText.Add(nil, id, textid)
        if err != nil {
            return
        }
    }
    return
}

func DelPlayer(id int64) (err error) {
    _, err = Db.Player.Del(id)
    if err != nil {
        return
    }
    var textid int64
    var playertextid int64
    err = Db.PlayerText.Get(nil, id).Scan(&playertextid, nil, &textid)
    if err == nil {
        _, err = Db.Text.Del(textid)
        if err != nil {
            return
        }
        _, err = Db.PlayerText.Del(playertextid)
        if err != nil {
            return
        }
    } else if err == sql.ErrNoRows {

    } else {
        return
    }
    return
}

func GetPlayer(id int64) (player *Player, err error) {
    var (
        name string
        sex int64
        country string
        rank int64
        birth time.Time
        text string
    )
    err = Db.Player.Get(id).Scan(nil, &name, &sex, &country, &rank, &birth)
    if err != nil {
        return nil, err
    }
    var textid int64
    var playertextid int64
    err = Db.PlayerText.Get(nil, id).Scan(&playertextid, nil, &textid)
    if err == nil {
        err = Db.Text.Get(textid).Scan(nil, &text)
        if err == nil {

        } else if err == sql.ErrNoRows {
            _, err = Db.PlayerText.Del(playertextid)
            if err != nil {
                return nil, err
            }
        } else {
            return nil, err
        }
    } else if err != sql.ErrNoRows {
        return nil, err
    }
    player = new(Player)
    player.Id = id
    player.Name = name
    player.Sex = formatSex(sex)
    player.Country = country
    player.Rank = formatRank(rank)
    player.Birth = birth
    player.Text = text
    return player, nil
}
