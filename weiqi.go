package weiqi

var statusMap = map[int]string{
	-1: "删除",
	0: "草稿",
	1: "发布",
}

func statusToString(status int) string {
	s, ok := statusMap[status]
	if ok {
		return s
	} else {
		return ""
	}
}

func stringToStatus(s string) int {
	switch s {
	case "删除":
		return -1
	case "草稿":
		return 0
	case "发布":
		return 1
	default:
		panic("weiqi: the status not found in status map")
	}
}
