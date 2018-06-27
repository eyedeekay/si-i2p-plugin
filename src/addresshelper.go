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
	jumpClient     *jumphelper.Client
	jumpHostString string
	jumpPortString string

	addressHelperURL string

	bookPath string
	bookFile *os.File
	pairs    []string

	err error
	c   bool
}

// CheckAddressHelper determines how the addresshelper will be used for an address
func (addressBook *AddressHelper) CheckAddressHelper(url *http.Request) (*http.Request, bool) {
	if url != nil {
		b, e := addressBook.jumpClient.Check(url.URL.String())
		if e != nil {
			return url, false
		}
		if !b {
			return url, !b
		} else {
			s, c := addressBook.jumpClient.Request(url.URL.String())
			if c != nil {
				url.URL.Host = s
			}
			return url, !b
		}

	}
	return url, false
}

// CheckAddressHelperString determines how the addresshelper will be used for an address
func (addressBook *AddressHelper) CheckAddressHelperString(url string) (string, bool) {
	b, e := addressBook.jumpClient.Check(url)
	if e != nil {
		return "", false
	}
	return url, !b
}

// NewAddressHelper creates a new address helper from string options
func NewAddressHelper(AddressHelperURL, jumpHost, jumpPort string) *AddressHelper {
	a, e := NewAddressHelperFromOptions(
		SetAddressHelperURL(AddressHelperURL),
		SetAddressHelperHost(jumpHost),
		SetAddressHelperPort(jumpPort),
		SetAddressBookPath("addressbook.txt"),
	)
	Fatal(e, "addresshelper.go failed to create addresshelper from strings", "addresshelper.go created from strings")
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
	a.jumpClient, a.err = jumphelper.NewClient(a.jumpHostString, a.jumpPortString)
	Fatal(a.err, "addresshelper.go failed to setup standalone addresshelper.", "addresshelper.go connecting standalone addresshelper:", a.addressHelperURL, a.jumpHostString, ":", a.jumpPortString)
	a.pairs = []string{}
	a.c = false
	return &a, a.err
}
