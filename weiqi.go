package weiqi

import "fmt"

const (
	c_statusDraft = iota
	c_statusRelease
	c_statusDelete
)

type State struct {
	Value int
	Name string
}

var weiqiStatus = []State{
	State{0, "草稿"}, State{1, "发布"}, State{2, "删除"},
}

func statusToString(statusValue int) string {
	for _, s := range weiqiStatus {
		if s.Value == statusValue {
			return s.Name
		}
	}
	panic(fmt.Sprint("status:", statusValue, " no found"))
}

func stringToStatus(statusName string) int {
	for _, s := range weiqiStatus {
		if s.Name == statusName {
			return s.Value
		}
	}
	panic(fmt.Sprint("status:", statusName, " no found"))
}
