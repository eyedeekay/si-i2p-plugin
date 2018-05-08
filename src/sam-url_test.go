package dii2p

import (
	"net/http"
	"testing"
)

func NewSamUrlHttpTest(t *testing.T) {
	req, err := http.NewRequest("GET", "http://i2p-projekt.i2p", nil)
    t.Error(err)
	samUrl := NewSamUrlHttp(req)
    t.Log(samUrl.subDirectory)
    samUrl.cleanupDirectory()
}

func TestNewSamUrlHttpString(t *testing.T) {
	samUrl := NewSamUrl("i2p-projekt.i2p")
    t.Log(samUrl.subDirectory)
    samUrl.cleanupDirectory()
}
