package main

import (
	"strconv"
)

func atoi(num string) int {
	n, err := strconv.Atoi(num)
	if err != nil {
		return 0
	}
	return n
}

func atoi64(num string) int64 {
	n, err := strconv.ParseInt(num, 10, 64)
	if err != nil {
		return 0
	}
	return n
}

func itoa(num int) string {
	return strconv.Itoa(num)
}
