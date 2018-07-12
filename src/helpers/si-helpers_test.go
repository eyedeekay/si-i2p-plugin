package dii2phelper

import (
	"fmt"
	"testing"
)

import (
	"github.com/eyedeekay/si-i2p-plugin/src/errors"
)

func TestHelperCheckURLType(t *testing.T) {
	if CheckURLType("i2p-projekt.i2p") {
		dii2perrs.Fatal(fmt.Errorf("i2p-projekt.i2p"), "This should be a failing URL", "")
	}
	if !CheckURLType("http://i2p-projekt.i2p") {
		dii2perrs.Fatal(fmt.Errorf("http://i2p-projekt.i2p"), "This should be a passing URL", "")
	}
}

func TestHelperCleanURL(t *testing.T) {
	fmt.Println(CleanURL("i2p-projekt.i2p"))
}
