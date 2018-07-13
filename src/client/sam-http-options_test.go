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
	timeout := 600
	lifeout := 1200
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
		t.Fatal("Host setting error")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPPort(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPPort("7656"),
	)
	if e != nil {
		t.Fatal("Port setting error from String")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPPortInt(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPPortInt(7656),
	)
	if e != nil {
		t.Fatal("Port setting error from Int")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPRequest(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPRequest("http://i2p-projekt.i2p"),
	)
	if e != nil {
		t.Fatal("")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPTimeout(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPTimeout(6),
	)
	if e != nil {
		t.Fatal("")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPKeepAlives(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPKeepAlives(true),
	)
	if e != nil {
		t.Fatal("")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPLifespan(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPLifespan(12),
	)
	if e != nil {
		t.Fatal("")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPTunLength(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPTunLength(1),
	)
	if e != nil {
		t.Fatal("")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPInQuantity(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPInQuantity(1),
	)
	if e != nil {
		t.Fatal("")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPOutQuantity(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPOutQuantity(1),
	)
	if e != nil {
		t.Fatal("")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPInBackupQuantity(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPInBackupQuantity(1),
	)
	if e != nil {
		t.Fatal("")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPOutBackupQuantity(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPOutBackupQuantity(1),
	)
	if e != nil {
		t.Fatal("")
	}
	h.CleanupClient()
}

func TestCreateSamHTTPOptionsSetSamHTTPIdleQuantity(t *testing.T) {
	h, e := NewSamHTTPFromOptions(
		SetSamHTTPIdleQuantity(2),
	)
	if e != nil {
		t.Fatal("")
	}
	h.CleanupClient()
}
