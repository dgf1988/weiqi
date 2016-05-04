package weiqi

import "fmt"

const (
	constStatusDraft = iota
	constStatusRelease
	constStatusDelete
)

//Status 状态定义
type Status struct {
	Value int64
	Name  string
}

var weiqiStatus = []Status{
	Status{constStatusDraft, "草稿"}, Status{constStatusRelease, "发布"}, Status{constStatusDelete, "删除"},
}

func formatStatus(statusValue int64) string {
	for _, s := range weiqiStatus {
		if s.Value == statusValue {
			return s.Name
		}
	}
	panic(fmt.Sprint("status:", statusValue, " no found"))
}

func parseStatus(statusName string) int64 {
	for _, s := range weiqiStatus {
		if s.Name == statusName {
			return s.Value
		}
	}
	panic(fmt.Sprint("status:", statusName, " no found"))
}
