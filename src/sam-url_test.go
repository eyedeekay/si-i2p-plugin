package dii2p

import (
	"net/http"
	"testing"
)

import (
    "github.com/eyedeekay/si-i2p-plugin/src/errors"
)

func NewSamURLHTTPTest(t *testing.T) {
	req, err := http.NewRequest("GET", "http://i2p-projekt.i2p", nil)
	t.Error(err)
	SamURL := NewSamURLHTTP(req)
	t.dii2perrs.Log(SamURL.subDirectory)
	SamURL.cleanupDirectory()
}

func TestNewSamURLHTTPString(t *testing.T) {
	SamURL := NewSamURL("i2p-projekt.i2p")
	t.dii2perrs.Log(SamURL.subDirectory)
	SamURL.cleanupDirectory()
}
