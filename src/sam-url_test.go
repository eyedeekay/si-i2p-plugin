package dii2p

import (
	"net/http"
	"testing"
)

func NewSamURLHTTPTest(t *testing.T) {
	req, err := http.NewRequest("GET", "http://i2p-projekt.i2p", nil)
	t.Error(err)
	SamURL := NewSamURLHTTP(req)
	t.Log(SamURL.subDirectory)
	SamURL.cleanupDirectory()
}

func TestNewSamURLHTTPString(t *testing.T) {
	SamURL := NewSamURL("i2p-projekt.i2p")
	t.Log(SamURL.subDirectory)
	SamURL.cleanupDirectory()
}
