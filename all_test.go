package weiqi

import (
	"testing"
)

func TestDbDesc(t *testing.T) {
	if rows, err := Texts.ListDesc(40, 0); err == nil {
		defer rows.Close()
		for rows.Next() {
			var text Text
			err = rows.Struct(&text)
			if err != nil {
				t.Fatal(err.Error())
			}
			t.Log(text)
		}
	} else {
		t.Fatal(err.Error())
	}
}
