package main

import (
	"net/http"
	"strings"
)

type addressHelper struct {
	rq  *http.Request
	err error
	c   bool
}

func (addressBook *addressHelper) checkAddressHelper(url *http.Request) *http.Request {
	if strings.Contains(addressBook.rq.URL.String(), "?i2paddresshelper=") {
		Log("si-http-proxy.go ?i2paddresshelper detected")
		temp := strings.Split(addressBook.rq.URL.Path, "/")
		var newpath string
		for _, s := range temp {
			if !strings.Contains(addressBook.rq.URL.String(), "?i2paddresshelper=") {
				newpath += s
			}
		}
		if strings.HasSuffix(newpath, "/") {
			newpath = newpath[:len(newpath)-len("/")]
		}
		addressBook.rq.URL.Path = newpath
		addressBook.rq, addressBook.err = http.NewRequest(addressBook.rq.Method, addressBook.rq.URL.Scheme+"://"+addressBook.rq.URL.Host+newpath, addressBook.rq.Body)
		if addressBook.c, addressBook.err = Fatal(addressBook.err, "addresshelper.go ", "addresshelper.go "); addressBook.c {
			Log("addresshelper.go rewrote request")
		}
		return addressBook.rq
	} else {
		addressBook.rq, addressBook.err = http.NewRequest(url.Method, url.URL.String(), url.Body)
		if addressBook.c, addressBook.err = Fatal(addressBook.err, "addresshelper.go ", "addresshelper.go "); addressBook.c {
			Log("addresshelper.go rewrote request")
		}
		return addressBook.rq
	}
}

func newAddressHelper() *addressHelper {
	var a addressHelper
	return &a
}
