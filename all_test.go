package weiqi

import (
	"testing"
)

func TestDbDesc(t *testing.T) {
	players, err := Players.List(100, 0)
	if err != nil {
		t.Error(err.Error())
		return
	}
	for i := range players {
		t.Log(i, players[i])
	}

	posts, err := Posts.List(100, 0)
	if err != nil {
		t.Error(err.Error())
		return
	}
	for i := range posts {
		t.Log(i, posts[i][1], len(posts[i][2].(string)))
	}
}
