package dii2p

import (
	"testing"
)

import (
    "github.com/eyedeekay/si-i2p-plugin/src/errors"
)

func TestCreateSamListOptionsAll(t *testing.T) {
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
		t.dii2perrs.Log("CreateSamListOptionsAll Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsInitAddress(t *testing.T) {
	Verbose = true
	DEBUG = true
	samProxies, err := CreateSamList(
		SetInitAddress("http://i2p-projekt.i2p"),
	)
	if err != nil {
		t.dii2perrs.Fatal(err)
	} else {
		t.dii2perrs.Log("CreateSamListOptionsInitAddress Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsHost(t *testing.T) {
	Verbose = true
	DEBUG = true
	samProxies, err := CreateSamList(
		SetHost("localhost"),
	)
	if err != nil {
		t.dii2perrs.Fatal(err)
	} else {
		t.dii2perrs.Log("CreateSamListOptionsHost Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsPort(t *testing.T) {
	Verbose = true
	DEBUG = true
	samProxies, err := CreateSamList(
		SetPort("7656"),
	)
	if err != nil {
		t.dii2perrs.Fatal(err)
	} else {
		t.dii2perrs.Log("CreateSamListOptionsPort Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsTimeout(t *testing.T) {
	Verbose = true
	DEBUG = true
	samProxies, err := CreateSamList(
		SetTimeout(600),
	)
	if err != nil {
		t.dii2perrs.Fatal(err)
	} else {
		t.dii2perrs.Log("CreateSamListOptionsTimeout Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsKeepAlives(t *testing.T) {
	Verbose = true
	DEBUG = true
	samProxies, err := CreateSamList(
		SetKeepAlives(true),
	)
	if err != nil {
		t.dii2perrs.Fatal(err)
	} else {
		t.dii2perrs.Log("CreateSamListOptionsKeepAlives Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsTunLength(t *testing.T) {
	Verbose = true
	DEBUG = true
	samProxies, err := CreateSamList(
		SetTunLength(3),
	)
	if err != nil {
		t.dii2perrs.Fatal(err)
	} else {
		t.dii2perrs.Log("CreateSamListOptionsKeepAlives Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsInQuantity(t *testing.T) {
	Verbose = true
	DEBUG = true
	samProxies, err := CreateSamList(
		SetInQuantity(15),
	)
	if err != nil {
		t.dii2perrs.Fatal(err)
	} else {
		t.dii2perrs.Log("CreateSamListOptionsKeepAlives Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsOutQuantity(t *testing.T) {
	Verbose = true
	DEBUG = true
	samProxies, err := CreateSamList(
		SetOutQuantity(15),
	)
	if err != nil {
		t.dii2perrs.Fatal(err)
	} else {
		t.dii2perrs.Log("CreateSamListOptionsKeepAlives Test Complete: true")
	}
	samProxies.CleanupClient()
}
