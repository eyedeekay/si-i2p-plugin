package dii2p

import (
	"testing"
)

import (
	"github.com/eyedeekay/si-i2p-plugin/src/errors"
)

func TestCreateSamList(t *testing.T) {
	dii2perrs.Verbose = true
	DEBUG = true
	samProxies, err := CreateSamList(
		SetInitAddress("http://i2p-projekt.i2p"),
		SetHost("localhost"),
		SetPort("7656"),
		SetTimeout(600),
		SetKeepAlives(true),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("CreateSamList Test Complete: true")
	}
	samProxies.CleanupClient()
}
