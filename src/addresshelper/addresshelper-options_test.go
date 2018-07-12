package dii2pah

import (
	"testing"
    "time"
)

import (
	"github.com/eyedeekay/si-i2p-plugin/src/errors"
)

func TestAddressHelperSetHost(t *testing.T) {
    time.Sleep(2 * time.Second)
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	_, err := NewAddressHelperFromOptions(
		SetAddressHelperHost("127.0.0.1"),
	)
    if err != nil {
        dii2perrs.Fatal(err,"set host error","")
    }
}
func TestAddressHelperSetPort(t *testing.T) {
    time.Sleep(2 * time.Second)
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	_, err := NewAddressHelperFromOptions(
		SetAddressHelperPort("7054"),
	)
    if err != nil {
        dii2perrs.Fatal(err,"set port error","")
    }
}

func TestAddressHelperSetPortInt(t *testing.T) {
    time.Sleep(2 * time.Second)
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	_, err := NewAddressHelperFromOptions(
		SetAddressHelperPortInt(7054),
	)
    if err != nil {
        dii2perrs.Fatal(err,"set port int error","")
    }
}
