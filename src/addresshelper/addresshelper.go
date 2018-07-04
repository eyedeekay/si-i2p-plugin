package dii2pah

import (
	//"io/ioutil"
	"net/http"
	//"net/url"
	"os"
	//"strings"
)

import (
	"github.com/eyedeekay/jumphelper/src"
	"github.com/eyedeekay/si-i2p-plugin/src/errors"
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

func (addressBook *AddressHelper) request(url *http.Request) *http.Request {
	return &http.Request{
		Method: url.Method,
		//URL:
		Proto:      url.Proto,
		ProtoMajor: url.ProtoMajor,
		ProtoMinor: url.ProtoMinor,
		Header:     url.Header,
		//Body:
		ContentLength:    url.ContentLength,
		TransferEncoding: url.TransferEncoding,
		Close:            url.Close,
		//Host:
		Form:     url.Form,
		PostForm: url.PostForm,
		//MultiPartForm: url.MultiPartForm,
		Trailer:    url.Trailer,
		RemoteAddr: url.RemoteAddr,
		//RequestURI: '',
	}
}

// CheckAddressHelper determines how the addresshelper will be used for an address
func (addressBook *AddressHelper) CheckAddressHelper(url *http.Request) (*http.Request, bool) {
	/*
	   	u, b := addressBook.CheckAddressHelperString(url.URL.String())
	   	if !b {
	   		dii2perrs.Warn(nil, "addresshelper.go !b"+u, "addresshelper.go !b"+u)
	   		url.URL.Host = u
	           //nurl := addressBook.request(url)
	   		//return nurl, true
	           return url, true
	   	}
	*/
	return url, false
}

// CheckAddressHelperString determines how the addresshelper will be used for an address
func (addressBook *AddressHelper) CheckAddressHelperString(url string) (string, bool) {
	if url != "" {
		b, e := addressBook.jumpClient.Check(url)
		if e != nil {
			dii2perrs.Warn(e, "addresshelper.go Address Lookup Error", "addresshelper.go this should never be reached")
			return "", false
		}
		if !b {
			s, c := addressBook.jumpClient.Request(url)
			if c == nil {
				dii2perrs.Warn(nil, "addresshelper.go !b "+s+".b32.i2p", "addresshelper.go !b "+s+".b32.i2p")
				url = s + ".b32.i2p"
				return url, true
			}
		}
		dii2perrs.Warn(nil, "addresshelper.go b "+url, "addresshelper.go b "+url)
		return url, false
	}
	return url, false
}

// NewAddressHelper creates a new address helper from string options
func NewAddressHelper(AddressHelperURL, jumpHost, jumpPort string) (*AddressHelper, error) {
	a, e := NewAddressHelperFromOptions(
		SetAddressHelperURL(AddressHelperURL),
		SetAddressHelperHost(jumpHost),
		SetAddressHelperPort(jumpPort),
		SetAddressBookPath("addressbook.txt"),
	)
	return a, e
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
	dii2perrs.Fatal(a.err, "addresshelper.go failed to setup standalone addresshelper.", "addresshelper.go connecting standalone addresshelper:", a.addressHelperURL, a.jumpHostString, ":", a.jumpPortString)
	a.pairs = []string{}
	a.c = false
	return &a, a.err
}
