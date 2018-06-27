package dii2p

import (
	//"io/ioutil"
	"net/http"
	//"net/url"
	"os"
	//"strings"

	"github.com/eyedeekay/jumphelper/src"
    "github.com/eyedeekay/si-i2p-plugin/src"
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
			dii2p.Warn(e, "addresshelper.go Address Lookup Error", "addresshelper.go this should never be reached")
			return url, false
		}
		if !b {
			dii2p.Warn(nil, "addresshelper.go !b"+url.URL.String()+".b32.i2p", "addresshelper.go !b"+url.URL.String()+".b32.i2p")
			return url, !b
		} else {
			s, c := addressBook.jumpClient.Request(url.URL.String())
			if c != nil {
				dii2p.Log(s + ".b32.i2p")
				url.URL.Host = s + ".b32.i2p"
			}
			return url, !b
		}

	}
	return url, false
}

// CheckAddressHelperString determines how the addresshelper will be used for an address
func (addressBook *AddressHelper) CheckAddressHelperString(url string) (string, bool) {
	if url != "" {
		b, e := addressBook.jumpClient.Check(url)
		if e != nil {
			dii2p.Warn(e, "addresshelper.go Address Lookup Error", "addresshelper.go this should never be reached")
			return "", false
		}
		if b {
			dii2p.Warn(nil, "addresshelper.go !b "+url+".b32.i2p", "addresshelper.go !b "+url+".b32.i2p")
			return url, false
		} else {
			s, c := addressBook.jumpClient.Request(url)
			if c != nil {
				dii2p.Warn(nil, "addresshelper.go b "+s+".b32.i2p", "addresshelper.go b "+s+".b32.i2p")
				url = s + ".b32.i2p"
			}
			return url, true
		}
	}
	return url, false
}

// NewAddressHelper creates a new address helper from string options
func NewAddressHelper(AddressHelperURL, jumpHost, jumpPort string) *AddressHelper {
	a, e := NewAddressHelperFromOptions(
		SetAddressHelperURL(AddressHelperURL),
		SetAddressHelperHost(jumpHost),
		SetAddressHelperPort(jumpPort),
		SetAddressBookPath("addressbook.txt"),
	)
	dii2p.Fatal(e, "addresshelper.go failed to create addresshelper from strings", "addresshelper.go created from strings")
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
	dii2p.Fatal(a.err, "addresshelper.go failed to setup standalone addresshelper.", "addresshelper.go connecting standalone addresshelper:", a.addressHelperURL, a.jumpHostString, ":", a.jumpPortString)
	a.pairs = []string{}
	a.c = false
	return &a, a.err
}
