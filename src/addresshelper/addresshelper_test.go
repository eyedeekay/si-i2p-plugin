package dii2pah

import (
	"testing"
	"time"
)

import (
	"github.com/eyedeekay/jumphelper/src"
	"github.com/eyedeekay/si-i2p-plugin/src/errors"
)

func CreateAddressHelperIB() {
	time.Sleep(2 * time.Second)
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	a, err := NewAddressHelperFromOptions(
		SetAddressHelperHost("127.0.0.1"),
		SetAddressHelperPort("7854"),
	)
	c, b := a.CheckAddressHelperString("i2p-projekt.i2p")
	dii2perrs.Log("TestCreateAddressHelperIB Test Complete: true", c)
	if !b {
		dii2perrs.Fatal(err, "TestCreateAddressHelperIB", "TestCreateAddressHelperIB")
	}
}

func CreateAddressHelperNIB() {
	time.Sleep(2 * time.Second)
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	a, err := NewAddressHelperFromOptions(
		SetAddressHelperHost("127.0.0.1"),
		SetAddressHelperPort("7854"),
	)
	c, b := a.CheckAddressHelperString("notarealaddress.i2p")
	dii2perrs.Log("TestCreateAddressHelperNIB", c)
	if !b {
		dii2perrs.Fatal(err, "TestCreateAddressHelperNIB Test Complete: true", "TestCreateAddressHelperNIB Test Complete: true")
	}
}

func TestCreateAddressHelper(t *testing.T) {
	jumphelper.NewService("localhost", "7854", "../../misc/addresses.csv", "127.0.0.1", "7656", []string{}, false)
	CreateAddressHelperIB()
	CreateAddressHelperNIB()
}
