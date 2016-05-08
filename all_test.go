package weiqi

import (
	"testing"
)

func TestSgf(t *testing.T) {
	var project Project
	project.Name = "人机大战"
	project.Text = "李世石打电脑"
	project.Items = make([]Item, 0)
	project.AddItem("资金", "100万美元")
	project.AddItem("地点", "韩国")
	project.AddItem("主办方", "谷歌公司")
	project.AddItem("协办方", "韩国棋院")

	var id int64
	var err error
	if id, err = addProject(project); err != nil {
		t.Fatal(err.Error())
	} else {
		p, err := getProject(id)
		if err != nil {
			t.Fatal(err.Error())
		} else {
			t.Log(*p)
		}
	}
}