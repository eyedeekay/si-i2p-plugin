package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	//	"path/filepath"
	"strconv"
	"strings"

	"github.com/eyedeekay/i2pasta/addresshelper"
	"github.com/eyedeekay/i2pasta/convert"
)

type addressHelper struct {
	assistant i2paddresshelper.I2paddresshelper
	converter i2pconv.I2pconv

	helperUrls []string

	rq       *http.Request
	bookPath string
	bookFile *os.File
	pairs    []string

	err error
	c   bool
}

func (addressBook *addressHelper) checkAddressHelper(url http.Request) (*http.Request, bool) {
	if strings.Contains(url.URL.String(), "?i2paddresshelper=") {
		addressBook.addPair(url.URL)
		_, b32 := addressBook.getBase32(url.URL)
		Log("addresshelper.go ?i2paddresshelper detected")
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
		body, err := ioutil.ReadAll(url.Body)
		if addressBook.c, addressBook.err = Fatal(err, "addresshelper.go Body rewrite error", "addresshelper.go Body rewriting"); addressBook.c {
			newBody := strings.Replace(string(body), url.Host, b32, -1)
			Log("addresshelper.go request body", url.Host, url.URL.Scheme+"://"+b32+newpath, string(newBody))

			addressBook.rq, addressBook.err = http.NewRequest(url.Method, url.URL.Scheme+"://"+b32+newpath, strings.NewReader(newBody))
			if addressBook.c, addressBook.err = Fatal(addressBook.err, "addresshelper.go New request formation error", "addresshelper.go New request generated"); addressBook.c {
				Log("addresshelper.go rewrote request")
			}
			return addressBook.rq, true
		}
	} else {
		addressBook.rq, addressBook.err = http.NewRequest(url.Method, url.URL.String(), url.Body)
		if addressBook.c, addressBook.err = Fatal(addressBook.err, "addresshelper.go Request return error", "addresshelper.go Returning same request"); addressBook.c {
			Log("addresshelper.go no rewrite required")
		}
		return addressBook.rq, false
	}
	return &url, false
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

func (addressBook *addressHelper) getBase32(url *url.URL) (string, string) {
	for _, p := range addressBook.pairs {
		kv := strings.SplitN(p, "=", 2)
		if kv != nil {
			if len(kv) == 2 {
				if kv[0] == url.Host {
					b32, err := addressBook.converter.I2p64to32(kv[1])
					if addressBook.c, addressBook.err = Warn(err, "addresshelper.go Base32 conversion failure", "Base32 converted"); addressBook.c {
                        Log("addresshelper.go b32:", b32)
						return kv[0], b32
					}
				}
			}
		}
	}
	return "", ""
}

func (addressBook *addressHelper) fileCheck(line string) bool {
	temp, err := ioutil.ReadFile(addressBook.bookPath)
	if addressBook.c, addressBook.err = Warn(err, "addresshelper.go File check error, handling:", "addresshelper.go Checking Addressbook file", addressBook.bookPath); addressBook.c {
		return !strings.Contains(string(temp), line)
	} else {
		return true
	}
}

func (addressBook *addressHelper) updateAh() {
	if addressBook.c, addressBook.err = exists(addressBook.bookPath); addressBook.c {
		addressBook.bookFile, addressBook.err = os.OpenFile(addressBook.bookPath, os.O_APPEND|os.O_WRONLY, 0755)
	} else {
		addressBook.bookFile, addressBook.err = os.OpenFile(addressBook.bookPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	}
	if addressBook.c, addressBook.err = Fatal(addressBook.err, "addresshelper.go File I/O errors", "addresshelper.go Addressbook file written"); addressBook.c {
		defer addressBook.bookFile.Close()
		line := addressBook.pairs[len(addressBook.pairs)-1] + "\n"
		if addressBook.fileCheck(line) {
			addressBook.bookFile.WriteString(line)
		} else {
			Log("addresshelper.go Address already in Address Book")
		}
	}
}

func newAddressHelper(addressHelperUrl string) *addressHelper {
	var a addressHelper
	a.helperUrls = make([]string, 0)
	a.helperUrls = append(a.helperUrls, strings.SplitN(addressHelperUrl, ",", -1)...)
	for index, address := range a.helperUrls {
		Log("addresshelper.go address:", address, "index:", strconv.Itoa(index))
	}
	a.pairs = []string{}
	a.rq = &http.Request{}
	a.err = nil
	a.c = false
	a.bookPath = "addressbook.txt"
	return &a
}
