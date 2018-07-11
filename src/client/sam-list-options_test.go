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
		SetTimeout(600),
		SetKeepAlives(true),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("CreateSamListOptionsAll Test Complete: true")
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
		t.Log("CreateSamListOptionsInitAddress Test Complete: true")
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
		t.Log("CreateSamListOptionsHost Test Complete: true")
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
		t.Log("CreateSamListOptionsPort Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsTimeout(t *testing.T) {
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	samProxies, err := CreateSamList(
		SetTimeout(600),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("CreateSamListOptionsTimeout Test Complete: true")
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
		t.Log("CreateSamListOptionsKeepAlives Test Complete: true")
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
		t.Log("CreateSamListOptionsKeepAlives Test Complete: true")
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
		t.Log("CreateSamListOptionsKeepAlives Test Complete: true")
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
		t.Log("CreateSamListOptionsKeepAlives Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsIdleConns(t *testing.T) {
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	samProxies, err := CreateSamList(
		SetIdleConns(15),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("CreateSamListOptionsKeepAlives Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsInBackups(t *testing.T) {
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	samProxies, err := CreateSamList(
		SetInBackups(15),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("CreateSamListOptionsKeepAlives Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsOutBackups(t *testing.T) {
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	samProxies, err := CreateSamList(
		SetOutBackups(15),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("CreateSamListOptionsKeepAlives Test Complete: true")
	}
	samProxies.CleanupClient()
}

func TestCreateSamListOptionsLifespan(t *testing.T) {
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	samProxies, err := CreateSamList(
		SetLifespan(15),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("CreateSamListOptionsKeepAlives Test Complete: true")
	}
	samProxies.CleanupClient()
}

/* Stub for new tests.
func TestCreateSamListOptionsXXXXXXXXXXXX(t *testing.T) {
	dii2perrs.Verbose = true
	dii2perrs.DEBUG = true
	samProxies, err := CreateSamList(
		SetXXXXXXXXXXXX(15),
	)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("CreateSamListOptionsKeepAlives Test Complete: true")
	}
	samProxies.CleanupClient()
}
*/
