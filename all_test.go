package weiqi

import (
	"os"
	"testing"
)

func TestDbDesc(t *testing.T) {
	htmlHead := defHtmlHead()
	err := htmlHead.Execute(os.Stdout, map[string]interface{}{
		"Head": Head{"title", "desc", []string{"1", "2"}}}, nil)
	if err != nil {
		t.Error(err.Error())
	}
}
