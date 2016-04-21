package weiqi

import (
	"testing"
)

func TestDbDesc(t *testing.T) {
	t.Log(Players.Set(12, nil, "xiaoli"))
	t.Log(Players.GetSlice(12))
}
