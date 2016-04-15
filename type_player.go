package main

import (
	"database/sql"
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
	return formatSex(p.Sex)
}

const (
	C_SEX_BOY  = 1
	C_SEX_GIRL = 2
)

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

func dbCountPlayer(where string) int64 {
	var n int64
	if where == "" {
		n, _ = dbCount("player")
	} else {
		n, _ = dbCountBy("player", where)
	}
	return n
}

func dbGetPlayer(id int64) (*Player, error) {
	if id <= 0 {
		return nil, ErrPrimaryKey
	}
	var p Player
	row := db.QueryRow("select * from player where id = ? limit 1", id)
	err := row.Scan(&p.Id, &p.Name, &p.Sex, &p.Country, &p.Rank, &p.Birth)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func dbFindPlayer(name string) (*Player, error) {
	var p Player
	row := db.QueryRow("select * from player where pname = ? limit 1", name)
	err := row.Scan(&p.Id, &p.Name, &p.Sex, &p.Country, &p.Rank, &p.Birth)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func dbWherePlayer(where string) (*Player, error) {
	if where == "" {
		return nil, sql.ErrNoRows
	}
	var p Player
	row := db.QueryRow("select * from player where " + where + " limit 1")
	err := row.Scan(&p.Id, &p.Name, &p.Sex, &p.Country, &p.Rank, &p.Birth)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func dbListPlayer(take, skip int) ([]Player, error) {
	rows, err := db.Query("select * from player order by id desc limit ?,?", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ps = make([]Player, 0)
	for rows.Next() {
		var p Player
		err := rows.Scan(&p.Id, &p.Name, &p.Sex, &p.Country, &p.Rank, &p.Birth)
		if err != nil {
			return ps, err
		}
		ps = append(ps, p)
	}
	return ps, rows.Err()
}

func dbAddPlayer(p *Player) (int64, error) {
	if p.Id > 0 {
		return -1, ErrPrimaryKey
	}
	res, err := db.Exec("insert into player (pname, psex, pcountry, prank, pbirth) values (?,?,?,?,?)", p.Name, p.Sex, p.Country, p.Rank, p.Birth)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

func dbUpdatePlayer(p *Player) (int64, error) {
	if p.Id <= 0 {
		return -1, ErrPrimaryKey
	}
	res, err := db.Exec("update player set pname=?,psex=?,pcountry=?,prank=?,pbirth=? where id = ? limit 1", p.Name, p.Sex, p.Country, p.Rank, p.Birth, p.Id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func dbDeletePlayer(id int64) (int64, error) {
	if id <= 0 {
		return -1, ErrPrimaryKey
	}
	res, err := db.Exec("delete from player where id = ? limit 1", id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}
