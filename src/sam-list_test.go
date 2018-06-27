package dii2p

import (
	"testing"
)

import (
    "github.com/eyedeekay/si-i2p-plugin/src/errors"
)

func TestCreateSamList(t *testing.T) {
	Verbose = true
	DEBUG = true
	samProxies, err := CreateSamList(
		SetInitAddress("http://i2p-projekt.i2p"),
		SetHost("localhost"),
		SetPort("7656"),
		SetTimeout(600),
		SetKeepAlives(true),
	)
	if err != nil {
		t.dii2perrs.Fatal(err)
	} else {
		t.dii2perrs.Log("CreateSamList Test Complete: true")
	}
	samProxies.CleanupClient()
}
