package weiqi

import (
	"testing"
)

func TestSgf(t *testing.T) {
    ps, err := ListProject(10, 0)
    if err != nil {
        t.Fatal(err.Error())
    }
    for _, p := range ps {
        t.Log(p.Id, p.Name, p.Text, p.Items)
    }
}