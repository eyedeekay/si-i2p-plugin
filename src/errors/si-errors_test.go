package dii2perrs

import (
	"fmt"
	"testing"
)

func TestLogFuncs(t *testing.T) {
	Log("Testing log")
	Warn(fmt.Errorf("Testing log"), "Testing log", "Testing log")
	Fatal(nil, "Testing log", "Testing log")
}
