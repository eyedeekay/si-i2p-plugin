package dii2pah

import (
	"testing"
)

import (
	"github.com/eyedeekay/si-i2p-plugin/src/errors"
)

func TestCreateAddressHelperIB(t *testing.T) {
	dii2perrs.Verbose = true
	//dii2perrs.DEBUG = true
	a, err := NewAddressHelperFromOptions(
		SetAddressHelperURL("http://inr.i2p"),
		SetAddressHelperHost("127.0.0.1"),
		SetAddressHelperPort("7054"),
		SetAddressBookPath("addressbook.txt"),
	)
	b32, b := a.CheckAddressHelperString("i2p-projekt.i2p")
	if b {
		t.Fatal(err)
	} else {
		t.Log("TestCreateAddressHelper Test Complete: true", b32)
	}
}

func TestCreateAddressHelperNIB(t *testing.T) {
	dii2perrs.Verbose = true
	//dii2perrs.DEBUG = true
	a, err := NewAddressHelperFromOptions(
		SetAddressHelperURL("http://inr.i2p"),
		SetAddressHelperHost("127.0.0.1"),
		SetAddressHelperPort("7054"),
		SetAddressBookPath("addressbook.txt"),
	)
	b32, b := a.CheckAddressHelperString("i2pforum.i2p")
	t.Log(b32, b, err)
	if b {
		t.Fatal("TestCreateAddressHelper Test Complete: false")
	}
}
