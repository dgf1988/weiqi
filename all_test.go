package weiqi

import (
	"testing"
)

func TestDbDesc(t *testing.T) {
	t.Log(Texts.Count(""))
	t.Log(PlayerText.Count(""))
}
