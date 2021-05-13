package test

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	m := make(map[string]int)
	value, ok := m["1"]
	t.Log(value)
	t.Log(ok)
	fmt.Scan()
}
