package dii2p

import (
	"net/http"
	"testing"
)

func NewSamURLHTTPTest(t *testing.T) {
	req, err := http.NewRequest("GET", "http://i2p-projekt.i2p", nil)
	t.Error(err)
	samURL := NewSamURLHTTP(req)
	t.Log(samURL.subDirectory)
	samURL.cleanupDirectory()
}

func TestNewSamURLHTTPString(t *testing.T) {
	samURL := NewSamURL("i2p-projekt.i2p")
	t.Log(samURL.subDirectory)
	samURL.cleanupDirectory()
}
