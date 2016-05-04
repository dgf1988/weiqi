package h

import (
	"testing"
)

func TestA(t *testing.T) {
	if parseMethod("POST") != POST {
		t.Error(parseMethod("POST"), POST)
	}
	if parseMethod("GET") != GET {
		t.Error(parseMethod("GET"), GET)
	}

	var a = []string{"POST"}
	var b = formatMethods(POST)
	for i := range a {
		if a[i] != b[i] {
			t.Error(a[i], b[i])
		}
	}
}
