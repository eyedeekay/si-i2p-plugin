package dii2pmain

import (
	"testing"
)

func TestCreateSamHTTPOptionsAll(t *testing.T) {
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
