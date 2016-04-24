package weiqi

import (
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
