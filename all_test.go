package weiqi

import (
	"testing"
	"github.com/dgf1988/weiqi/db"
)

func TestSgf(t *testing.T) {
	var sgfs = make([]Sgf, 0)
	var err error
	var rows db.Rows
	if rows, err = Db.Sgf.List(400,0); err != nil {
		t.Fatal(err.Error())
	} else {
		defer rows.Close()
		for rows.Next() {
			var sgf Sgf
			if err = rows.Struct(&sgf); err != nil {
				t.Fatal(err.Error())
			} else {
				sgfs = append(sgfs, sgf)
			}
		}
		if err = rows.Err(); err != nil {
			t.Fatal(err.Error())
		}
	}
	for i := range sgfs {
		if _, err = Db.Sgf.Update(sgfs[i].Id).Values(nil,  nil, nil,  nil, nil, nil, nil, nil, sgfs[i].ToSgf()); err != nil {
			t.Fatal(sgfs[i].Id, err.Error())
		}
	}
}