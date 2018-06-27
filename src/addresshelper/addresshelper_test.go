package dii2pah

import (
	"testing"
)

func TestCreateAddressHelperIB(t *testing.T) {
	Verbose = true
	DEBUG = true
	a, err := NewAddressHelperFromOptions(
		SetAddressHelperURL("http://inr.i2p"),
		SetAddressHelperHost("127.0.0.1"),
		SetAddressHelperPort("7054"),
		SetAddressBookPath("addressbook.txt"),
	)
	b32, b := a.CheckAddressHelperString("i2p-projekt.i2p")
	if b {
		t.dii2perrs.Fatal(err)
	} else {
		t.dii2perrs.Log("TestCreateAddressHelper Test Complete: true", b32)
	}
}

func TestCreateAddressHelperNIB(t *testing.T) {
	Verbose = true
	DEBUG = true
	a, err := NewAddressHelperFromOptions(
		SetAddressHelperURL("http://inr.i2p"),
		SetAddressHelperHost("127.0.0.1"),
		SetAddressHelperPort("7054"),
		SetAddressBookPath("addressbook.txt"),
	)
	b32, b := a.CheckAddressHelperString("i2pforum.i2p")
	t.dii2perrs.Log(b32, b, err)
	if b {
		t.dii2perrs.Fatal("TestCreateAddressHelper Test Complete: false")
	}
}
