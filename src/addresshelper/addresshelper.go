package dii2pah

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
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

	err error
	c   bool
}

func (addressBook *AddressHelper) request(req *http.Request, addr string) *http.Request {
	u, e := url.Parse(addr)
	if e != nil {
		return req
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		dii2perrs.Warn(err, "", "")
		return req
	}
	contentLength := int64(len(body))
	return &http.Request{
		Method:           req.Method,
		URL:              u,
		Proto:            req.Proto,
		ProtoMajor:       req.ProtoMajor,
		ProtoMinor:       req.ProtoMinor,
		Header:           req.Header,
		Body:             ioutil.NopCloser(strings.NewReader(string(body))),
		ContentLength:    contentLength,
		TransferEncoding: req.TransferEncoding,
		Close:            req.Close,
		Form:             req.Form,
		PostForm:         req.PostForm,
		MultipartForm:    req.MultipartForm,
		Trailer:          req.Trailer,
		RequestURI:       "",
		Response:         req.Response,
	}
}

// CheckAddressHelper determines how the addresshelper will be used for an address
func (addressBook *AddressHelper) CheckAddressHelper(req *http.Request) (*http.Request, bool) {
	u, b := addressBook.CheckAddressHelperString(req.URL.String())
	if !b {
		dii2perrs.Warn(nil, "addresshelper.go !b"+u, "addresshelper.go !b"+u)
		return addressBook.request(req, u), true
	}
	return req, false
}

// CheckAddressHelperString determines how the addresshelper will be used for an address
func (addressBook *AddressHelper) CheckAddressHelperString(req string) (string, bool) {
	if req != "" {
        if strings.HasSuffix(strings.Split(req,".b32.i2p")[0], ".b32.i2p") {
            return req, false
        }
		b, e := addressBook.jumpClient.Check(req)
		if e != nil {
			dii2perrs.Warn(e, "addresshelper.go Address Lookup Error", "addresshelper.go this should never be reached")
			return "", false
		}
		if !b {
			s, c := addressBook.jumpClient.Request(req)
			if c == nil {
				dii2perrs.Warn(nil, "addresshelper.go !b "+s+".b32.i2p", "addresshelper.go !b "+s+".b32.i2p")
				req = s + ".b32.i2p"
				return req, true
			}
		}
		dii2perrs.Warn(nil, "addresshelper.go b "+req, "addresshelper.go b "+req)
		return req, false
	}
	return req, false
}

// NewAddressHelper creates a new address helper from string options
func NewAddressHelper(AddressHelperURL, jumpHost, jumpPort string) (*AddressHelper, error) {
	a, e := NewAddressHelperFromOptions(
		SetAddressHelperHost(jumpHost),
		SetAddressHelperPort(jumpPort),
	)
	return a, e
}

// NewAddressHelperFromOptions creates a new address helper from functional arguments
func NewAddressHelperFromOptions(opts ...func(*AddressHelper) error) (*AddressHelper, error) {
	var a AddressHelper
	a.jumpHostString = "127.0.0.1"
	a.jumpPortString = "7854"
	for _, o := range opts {
		if err := o(&a); err != nil {
			return nil, err
		}
	}
	a.jumpClient, a.err = jumphelper.NewClient(a.jumpHostString, a.jumpPortString, false)
	dii2perrs.Fatal(a.err, "addresshelper.go failed to setup standalone addresshelper.", "addresshelper.go connecting standalone addresshelper:", a.jumpHostString, ":", a.jumpPortString)
	a.c = false
	return &a, a.err
}
