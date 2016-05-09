package weiqi

import (
	"testing"
)

func TestSgf(t *testing.T) {
    var p *Player
    var err error

    p, err = GetPlayer(1)
    if err != nil {
        t.Fatal(err.Error())
    } else {
        p.Id = 0
        p.Name = "丁国锋"
        p.Rank = "七段"
        p.Text = "我是丁国锋"
        var id int64
        id, err = p.Save()
        if id > 0 {
            t.Log("是插入")
        } else {
            t.Log("是更新")
        }
        if err != nil {
            t.Fatal(err.Error())
        }
        t.Log(GetPlayer(1))
    }
}