package dii2p


import (
	"testing"
)

func TestCreateSamListt(t *testing.T){
    samlst := CreateSamList("localhost", "7656", "http://i2p-projekt.i2p", 600, true)
    samlst.CleanupClient()
}
