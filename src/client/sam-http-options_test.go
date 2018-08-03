package dii2pmain

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func request() *http.Request {
	u, e := url.Parse("http://i2p-projekt.i2p")
	if e != nil {
		return nil
	}
	body := ""
	contentLength := int64(len(body))
	return &http.Request{
		URL:           u,
		Body:          ioutil.NopCloser(strings.NewReader(string(body))),
		ContentLength: contentLength,
		RequestURI:    "",
	}
}

func TestCreateSamHTTPOptionsAllOld(t *testing.T) {
	length := 1
	quant := 15
	timeout := 6
	lifeout := 12
	backup := 3
	idles := 4
	h := newSamHTTPHTTP("127.0.0.1",
		"7656",
		request(),
		timeout,
		lifeout,
		true,
		length,
		quant,
		quant,
		idles,
		backup,
		backup)
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPHost(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPHost("127.0.0.1"),
	)
	if e != nil {
		t.Fatal("sam-http-options_test.go Host setting error")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPPort(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPPort("7656"),
	)
	if e != nil {
		t.Fatal("sam-http-options_test.go Port setting error from String")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPPortInt(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPPortInt(7656),
	)
	if e != nil {
		t.Fatal("sam-http-options_test.go Port setting error from Int")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPRequest(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPRequest("http://i2p-projekt.i2p"),
	)
	if e != nil {
		t.Fatal("sam-http-options_test.go Error with test HTTP request over SAM")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPTimeout(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPTimeout(6),
	)
	if e != nil {
		t.Fatal("sam-http-options_test.go Timeout setting error")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPKeepAlives(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPKeepAlives(true),
	)
	if e != nil {
		t.Fatal("sam-http-options_test.go Keep-Alive setting error")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPLifespan(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPLifespan(12),
	)
	if e != nil {
		t.Fatal("sam-http-options_test.go Lifespan setting error")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPTunLength(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPTunLength(1),
	)
	if e != nil {
		t.Fatal("sam-http-options_test.go Error setting tunnel length")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPInQuantity(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPInQuantity(1),
	)
	if e != nil {
		t.Fatal("sam-http-options_test.go Error setting inbound tunnel quantity")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPOutQuantity(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPOutQuantity(1),
	)
	if e != nil {
		t.Fatal("sam-http-options_test.go Error setting outbound tunnel quantity")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPInBackupQuantity(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPInBackupQuantity(1),
	)
	if e != nil {
		t.Fatal("sam-http-options_test.go Error setting inbound backup tunnel quantity")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPOutBackupQuantity(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPOutBackupQuantity(1),
	)
	if e != nil {
		t.Fatal("sam-http-options_test.go Error setting outbound backup tunnel quantity")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPIdleQuantity(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPIdleQuantity(2),
	)
	if e != nil {
		t.Fatal("sam-http-options_test.go Error setting idle tunnel quantity")
	}
	h.CleanupClient()
}
