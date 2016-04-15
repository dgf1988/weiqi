package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

/*
CREATE TABLE `sgf` (
`id` INT(11) NOT NULL AUTO_INCREMENT,
`stime` DATE NOT NULL DEFAULT '0000-00-00',
`splace` CHAR(50) NOT NULL DEFAULT '',
`sevent` CHAR(100) NOT NULL DEFAULT '',
`sblack` CHAR(50) NOT NULL DEFAULT 'b',
`swhite` CHAR(50) NOT NULL DEFAULT 'w',
`srule` CHAR(50) NOT NULL DEFAULT '',
`sresult` CHAR(50) NOT NULL DEFAULT '',
`ssteps` MEDIUMTEXT NOT NULL,
`supdate` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
PRIMARY KEY (`id`)
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB
;
*/

const (
	SGF_TABLENAME = "sgf"
	SGF_CHARSET   = "utf-8"
)

type Sgf struct {
	Id     int64
	Time   time.Time
	Place  string
	Event  string
	Black  string
	White  string
	Rule   string
	Result string
	Steps  string
	Update time.Time
}

func (this Sgf) ToSgf() string {
	if this.Steps == "" {
		return ""
	}
	items_sgf := make([]string, 0)
	items_sgf = append(items_sgf, "(;")
	items_sgf = append(items_sgf, fmt.Sprintf("CA[%s]EV[%s]", SGF_CHARSET, this.Event))
	if this.Time.IsZero() {
		items_sgf = append(items_sgf, "DT[0000-00-00]")
	} else {
		items_sgf = append(items_sgf, fmt.Sprintf("DT[%s]", this.Time.Format("2006-01-02")))
	}
	items_sgf = append(items_sgf, fmt.Sprintf("PC[%s]", this.Place))
	items_sgf = append(items_sgf, fmt.Sprintf("PB[%s]", this.Black))
	items_sgf = append(items_sgf, fmt.Sprintf("PW[%s]", this.White))
	items_sgf = append(items_sgf, fmt.Sprintf("RL[%s]", this.Rule))
	items_sgf = append(items_sgf, fmt.Sprintf("RE[%s]", this.Result))
	items_sgf = append(items_sgf, "\n")
	items_sgf = append(items_sgf, this.Steps)
	items_sgf = append(items_sgf, ")")
	return strings.Join(items_sgf, "")
}

func dbCountSgf(where string) (int64, error) {
	if where == "" {
		return dbCount(SGF_TABLENAME)
	}
	return dbCountBy(SGF_TABLENAME, where)
}

func dbAddSgf(s *Sgf) (int64, error) {
	insertsql := "insert into sgf (stime, splace, sevent, sblack, swhite, srule, sresult, ssteps) values (?,?,?,?,?,?,?,?)"
	res, err := db.Exec(insertsql, s.Time, s.Place, s.Event, s.Black, s.White, s.Rule, s.Result, s.Steps)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}

func dbUpdateSgf(s *Sgf) (int64, error) {
	if s.Id <= 0 {
		return -1, ErrPrimaryKey
	}
	updatesql := "update sgf set stime=?, splace=?, sevent=?, sblack=?, swhite=?, srule=?, sresult=?, ssteps=? where id = ?  limit 1"
	res, err := db.Exec(updatesql, s.Time, s.Place, s.Event, s.Black, s.White, s.Rule, s.Result, s.Steps, s.Id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func dbGetSgf(id int64) (*Sgf, error) {
	getsql := "select * from sgf where id = ? limit 1"
	row := db.QueryRow(getsql, id)
	var s Sgf
	err := row.Scan(&s.Id, &s.Time, &s.Place, &s.Event, &s.Black, &s.White, &s.Rule, &s.Result, &s.Steps, &s.Update)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func dbDelSgf(id int64) (int64, error) {
	res, err := db.Exec("delete from sgf where id = ? limit 1", id)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func dbFindSgf(where string) (*Sgf, error) {
	if where == "" {
		return nil, sql.ErrNoRows
	}
	var s Sgf
	row := db.QueryRow("select * from sgf where " + where + " limit 1")
	err := row.Scan(&s.Id, &s.Time, &s.Place, &s.Event, &s.Black, &s.White, &s.Rule, &s.Result, &s.Steps, &s.Update)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func dbListSgf(take, skip int) ([]Sgf, error) {
	rows, err := db.Query("select * from sgf  order by sgf.id desc limit ?,?", skip, take)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	sgflist := make([]Sgf, 0)
	for rows.Next() {
		var s Sgf
		err := rows.Scan(&s.Id, &s.Time, &s.Place, &s.Event, &s.Black, &s.White, &s.Rule, &s.Result, &s.Steps, &s.Update)
		if err != nil {
			return sgflist, err
		}
		sgflist = append(sgflist, s)
	}
	return sgflist, rows.Err()
}
