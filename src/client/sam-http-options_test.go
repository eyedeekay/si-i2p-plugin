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
