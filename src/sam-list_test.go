package dii2p

import (
	"testing"
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
		t.Fatal(err)
	} else {
		t.Log("CreateSamList Test Complete: true")
	}
	samProxies.CleanupClient()
}
