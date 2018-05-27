package dii2p

import (
	"net/http"
	"testing"
)

func NewSamURLHttpTest(t *testing.T) {
	req, err := http.NewRequest("GET", "http://i2p-projekt.i2p", nil)
	t.Error(err)
	samUrl := NewSamURLHttp(req)
	t.Log(samUrl.subDirectory)
	samUrl.cleanupDirectory()
}

func TestNewSamURLHttpString(t *testing.T) {
	samUrl := NewSamURL("i2p-projekt.i2p")
	t.Log(samUrl.subDirectory)
	samUrl.cleanupDirectory()
}
