package dii2p

import (
    //"io/ioutil"
	"net/http"
	//"net/url"
	"os"
	//"strings"

    "github.com/eyedeekay/jumphelper/src"
)

// AddressHelper prioritizes the address sources
type AddressHelper struct {
    jumpClient *jumphelper.Client
	jumpHostString    string
	jumpPortString    string

    addressHelperURL string

	bookPath string
	bookFile *os.File
	pairs    []string

	err error
	c   bool
}

// CheckAddressHelper determines how the addresshelper will be used for an address
func (addressBook *AddressHelper) CheckAddressHelper(url *http.Request) (*http.Request, bool) {
	return url, false
}

// NewAddressHelper creates a new address helper from string options
func NewAddressHelper(AddressHelperURL, jumpHost, jumpPort string) *AddressHelper {
	a, _ := NewAddressHelperFromOptions(
		SetAddressHelperURL(AddressHelperURL),
		SetAddressHelperHost(jumpHost),
		SetAddressHelperPort(jumpPort),
		SetAddressBookPath("addressbook.txt"),
	)
	return a
}

// NewAddressHelperFromOptions creates a new address helper from functional arguments
func NewAddressHelperFromOptions(opts ...func(*AddressHelper) error) (*AddressHelper, error) {
	var a AddressHelper
	a.addressHelperURL = "inr.i2p"
	a.jumpHostString = "127.0.0.1"
	a.jumpPortString = "7054"
	a.bookPath = "addressbook.txt"
	for _, o := range opts {
		if err := o(&a); err != nil {
			return nil, err
		}
	}
	Fatal(a.err, "addresshelper.go failed to setup SAM bridge for addresshelper.", "addresshelper.go connecting to SAM bridge on:", a.addressHelperURL, a.jumpHostString, ":", a.jumpPortString)
	a.pairs = []string{}
	a.c = false
	return &a, a.err
}
