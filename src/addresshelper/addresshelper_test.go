package dii2pah

import (
	"testing"
	"time"
)

import (
	"github.com/eyedeekay/si-i2p-plugin/src/errors"
)

func TestCreateAddressHelperIB(t *testing.T) {
	time.Sleep(2 * time.Second)
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	a, err := NewAddressHelperFromOptions(
		SetAddressHelperURL("http://inr.i2p"),
		SetAddressHelperHost("127.0.0.1"),
		SetAddressHelperPort("7054"),
		SetAddressBookPath("addressbook.txt"),
	)
	c, b := a.CheckAddressHelperString("i2p-projekt.i2p")
	t.Log("TestCreateAddressHelperIB Test Complete: true", c)
	if !b {
		t.Fatal("TestCreateAddressHelperIB", err)
	}
}

func TestCreateAddressHelperNIB(t *testing.T) {
	time.Sleep(2 * time.Second)
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	a, err := NewAddressHelperFromOptions(
		SetAddressHelperURL("http://inr.i2p"),
		SetAddressHelperHost("127.0.0.1"),
		SetAddressHelperPort("7054"),
		SetAddressBookPath("addressbook.txt"),
	)
	c, b := a.CheckAddressHelperString("fireaxe.i2p")
	t.Log("TestCreateAddressHelperNIB", c)
	if !b {
		t.Fatal("TestCreateAddressHelperNIB Test Complete: true", err)
	}
}
