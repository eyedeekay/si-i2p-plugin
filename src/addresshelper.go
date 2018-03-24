package main

import (
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type addressHelper struct {
	rq       *http.Request
	bookPath string
	bookFile *os.File
	pairs    []string

	err error
	c   bool
}

func (addressBook *addressHelper) initFiles() {
	addressBook.bookPath, addressBook.bookFile, addressBook.err = setupFile(filepath.Join(connectionDirectory), addressBook.bookPath)
	if addressBook.c, addressBook.err = Fatal(addressBook.err, "sam-http.go Pipe setup error", "sam-http.go Pipe setup"); addressBook.c {
		addressBook.bookFile.WriteString("")
	}
}

func (addressBook *addressHelper) checkAddressHelper(url http.Request) (*http.Request, bool) {
	if strings.Contains(url.URL.String(), "?i2paddresshelper=") {
		addressBook.addPair(url.URL)
		Log("si-http-proxy.go ?i2paddresshelper detected")
		temp := strings.Split(url.URL.Path, "/")
		var newpath string
		for _, s := range temp {
			if !strings.Contains(url.URL.String(), "?i2paddresshelper=") {
				newpath += s
			}
		}
		if strings.HasSuffix(newpath, "/") {
			newpath = newpath[:len(newpath)-len("/")]
		}
		addressBook.rq, addressBook.err = http.NewRequest(url.Method, url.URL.Scheme+"://"+url.URL.Host+newpath, url.Body)
		if addressBook.c, addressBook.err = Fatal(addressBook.err, "addresshelper.go ", "addresshelper.go "); addressBook.c {
			Log("addresshelper.go rewrote request")
		}
		return addressBook.rq, true
	} else {
		addressBook.rq, addressBook.err = http.NewRequest(url.Method, url.URL.String(), url.Body)
		if addressBook.c, addressBook.err = Fatal(addressBook.err, "addresshelper.go ", "addresshelper.go "); addressBook.c {
			Log("addresshelper.go no rewrite required")
		}
		return addressBook.rq, false
	}
}

func (addressBook *addressHelper) checkAddPair(arg string) bool {
	for _, pair := range addressBook.pairs {
		kvPair := strings.SplitN(pair, "=", 2)
		if kvPair != nil {
			if len(kvPair) == 2 {
				if kvPair[0] == arg {
					return false
				}
			}
		}
	}
	return true
}

func (addressBook *addressHelper) addPair(url *url.URL) {
	segments := strings.Split(url.String(), "/")
	host := url.Host
	for _, s := range segments {
		if strings.Contains(s, "?i2paddresshelper=") {
			if addressBook.checkAddPair(host) {
				base64 := strings.Replace(strings.Split(s, "/")[0], "?i2paddresshelper=", "", -1)
				addressBook.pairs = append(addressBook.pairs, host+"="+base64)
			}
		}
	}
	addressBook.updateAh()
}

func (addressBook *addressHelper) getPair(url *url.URL) (string, string) {
	for _, p := range addressBook.pairs {
		kv := strings.SplitN(p, "=", 2)
		if kv != nil {
			if len(kv) == 2 {
				if kv[0] == url.Host {
					return kv[0], kv[1]
				}
			}
		}
	}
	return "", ""
}

func (addressBook *addressHelper) updateAh() {
	exist, _ := exists(addressBook.bookPath)
	if exist {
		os.Remove(addressBook.bookPath)
	}
	addressBook.bookFile, addressBook.err = os.Create(addressBook.bookPath)
	if addressBook.c, addressBook.err = Fatal(addressBook.err, "File I/O errors"); addressBook.c {
		defer addressBook.bookFile.Close()
		for _, t := range addressBook.pairs {
			line := t + "\n"
			addressBook.bookFile.WriteString(line)
		}
	}
}

func newAddressHelper() *addressHelper {
	var a addressHelper
	a.pairs = []string{}
	a.rq = &http.Request{}
	a.err = nil
	a.c = false
	a.bookPath = "addressbook.txt"
	a.initFiles()
	return &a
}
