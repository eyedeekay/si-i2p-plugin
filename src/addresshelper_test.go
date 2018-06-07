package dii2p

import (
	"testing"
)

func TestCreateAddressHelper(t *testing.T) {
	Verbose = true
	DEBUG = true
	a, err := NewAddressHelperFromOptions(
		SetAddressHelperURL("http://inr.i2p"),
		SetAddressHelperHost("127.0.0.1"),
		SetAddressHelperPort("7656"),
		SetAddressBookPath("addressbook.txt"),
	)
	a.Lookup("i2pforum.i2p")
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("TestCreateAddressHelper Test Complete: true")
	}
}
