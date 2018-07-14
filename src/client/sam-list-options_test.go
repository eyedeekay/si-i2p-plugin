package dii2pmain

import (
	"testing"
)

import (
	"github.com/eyedeekay/si-i2p-plugin/src/errors"
)

func TestCreateSamListOptionsAll(t *testing.T) {
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	samProxies, err := CreateSamList(
		SetInitAddress("http://i2p-projekt.i2p"),
		SetHost("localhost"),
		SetPort("7656"),
		SetTimeout(6),
		SetKeepAlives(true),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("sam-list-options_test.go CreateSamListOptionsAll Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsHost(t *testing.T) {
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	samProxies, err := CreateSamList(
		SetHost("localhost"),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("sam-list-options_test.go CreateSamListOptionsHost Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsPort(t *testing.T) {
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	samProxies, err := CreateSamList(
		SetPort("7656"),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("sam-list-options_test.go CreateSamListOptionsPort Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsPortInt(t *testing.T) {
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	samProxies, err := CreateSamList(
		SetPortInt(7656),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("sam-list-options_test.go CreateSamListOptionsPort Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsTimeout(t *testing.T) {
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	samProxies, err := CreateSamList(
		SetTimeout(6),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("sam-list-options_test.go CreateSamListOptionsTimeout Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsKeepAlives(t *testing.T) {
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	samProxies, err := CreateSamList(
		SetKeepAlives(true),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("sam-list-options_test.go CreateSamListOptionsKeepAlives Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsInitAddress(t *testing.T) {
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	samProxies, err := CreateSamList(
		SetInitAddress("http://i2p-projekt.i2p"),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("sam-list-options_test.go CreateSamListOptionsInitAddress Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsLifespan(t *testing.T) {
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	samProxies, err := CreateSamList(
		SetLifespan(12),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("sam-list-options_test.go CreateSamListOptionsLifeSpan Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsTunLength(t *testing.T) {
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	samProxies, err := CreateSamList(
		SetTunLength(3),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("sam-list-options_test.go CreateSamListOptionsTunLength Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsInQuantity(t *testing.T) {
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	samProxies, err := CreateSamList(
		SetInQuantity(15),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("sam-list-options_test.go CreateSamListOptionsInQuantity Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsOutQuantity(t *testing.T) {
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	samProxies, err := CreateSamList(
		SetOutQuantity(15),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("sam-list-options_test.go CreateSamListOptionsOutQuantity Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsIdleConns(t *testing.T) {
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	samProxies, err := CreateSamList(
		SetIdleConns(4),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("sam-list-options_test.go CreateSamListOptionsIdleConns Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsSetInBackups(t *testing.T) {
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	samProxies, err := CreateSamList(
		SetInBackups(3),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("sam-list-options_test.go CreateSamListOptionsInBackups Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsOutBackups(t *testing.T) {
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	samProxies, err := CreateSamList(
		SetOutBackups(3),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("sam-list-options_test.go CreateSamListOptionsOutBackups Test Complete: true")
	}
	samProxies.CleanupClient()
}
