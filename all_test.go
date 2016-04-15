package weiqi

import (
	"os"
	"testing"
)

func TestDbDesc(t *testing.T) {
	html_head := defHtmlHead()
	err := html_head.Execute(os.Stdout, map[string]interface{}{
		"a": Head{"title", "desc", "keywords"},
	})
	if err != nil {
		t.Error(err.Error())
	}
}
