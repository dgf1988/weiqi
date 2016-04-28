package weiqi

import (
	"fmt"
	"strings"
	"time"
	"sort"
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
	c_SgfCharset = "utf-8"
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
	items_sgf = append(items_sgf, fmt.Sprintf("CA[%s]EV[%s]", c_SgfCharset, this.Event))
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

type sgfOrderByTimeDesc []Sgf

func (arr sgfOrderByTimeDesc) Len() int { return len(arr)}
func (arr sgfOrderByTimeDesc) Swap(i, j int) { arr[i], arr[j] = arr[j], arr[i]}
func (arr sgfOrderByTimeDesc) Less(i, j int) bool { return arr[i].Time.After(arr[j].Time)}

func listSgfOrderByTimeDesc(take, skip int) ([]Sgf, error) {
	var sgfs = make([]Sgf, 0)
	if rows, err := Db.Sgf.Query("order by sgf.time desc limit ?, ?", skip, take); err != nil {
		return nil, err
	} else {
		defer rows.Close()
		for rows.Next() {
			var sgf Sgf
			if err = rows.Struct(&sgf); err != nil {
				return nil, err
			} else {
				sgfs = append(sgfs, sgf)
			}
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	return sgfs, nil
}

func listSgfByNameOrderByTimeDesc(name string) ([]Sgf, error) {
	var sgfs = make([]Sgf, 0)
	if rows, err := Db.Sgf.Any(nil, nil, nil, nil, name, name); err != nil {
		return nil, err
	} else {
		defer rows.Close()
		for rows.Next() {
			var sgf Sgf
			if err = rows.Struct(&sgf); err != nil {
				return nil, err
			} else {
				sgfs = append(sgfs, sgf)
			}
		}
		if err = rows.Err(); err != nil {
			return nil, err
		}
	}
	sort.Sort(sgfOrderByTimeDesc(sgfs))
	return sgfs, nil
}

func listSgfByNamesOrderByTimeDesc(names ...string) ([]Sgf, error) {
	var sgfs = make([]Sgf, 0)
	for _, name := range names {
		if name == "" {
			continue
		}
		if rows, err := Db.Sgf.Query("where sgf.black = ? or sgf.white = ? order by sgf.time desc", name, name); err != nil {
			return nil, err
		} else {
			defer rows.Close()
			for rows.Next() {
				var sgf Sgf
				if err = rows.Struct(&sgf); err != nil {
					return nil, err
				} else {
					sgfs = append(sgfs, sgf)
				}
			}
			if err = rows.Err(); err != nil {
				return nil, err
			}
		}
	}
	sort.Sort(sgfOrderByTimeDesc(sgfs))
	return sgfs, nil
}

func countSgfNumberByPlayerName(name string) (int64, error) {
	return Db.Sgf.Count("where sgf.black = ? or sgf.white = ?", name, name)
}