package dii2p

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/eyedeekay/i2pasta/addresshelper"
	"github.com/eyedeekay/i2pasta/convert"
)

type AddressHelper struct {
	assistant *i2paddresshelper.I2paddresshelper
	converter i2pconv.I2pconv

	addressHelperURL string
	samHostString    string
	samPortString    string

	bookPath string
	bookFile *os.File
	pairs    []string

	err error
	c   bool
}

func (addressBook *AddressHelper) base32ify(url *http.Request) (*http.Request, bool) {
	_, b32 := addressBook.getBase32(url.URL)
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

		rq, err := http.NewRequest(url.Method, url.URL.Scheme+"://"+b32+newpath, strings.NewReader(newBody))
		if addressBook.c, addressBook.err = Fatal(err, "addresshelper.go New request formation error", "addresshelper.go New request generated"); addressBook.c {
			Log("addresshelper.go rewrote request")
		}
		return rq, true
	}
	return url, false
}

func (addressBook *AddressHelper) checkAddressHelper(url *http.Request) (*http.Request, bool) {
	if strings.Contains(url.URL.String(), "?i2paddresshelper=") {
		Log("addresshelper.go ?i2paddresshelper detected")
		addressBook.addPair(url.URL)
		return addressBook.base32ify(url)
	} else if !addressBook.checkAddPair(url.URL.Host) {
		Log("addresshelper.go addressBook URL detected")
		return addressBook.base32ify(url)
	} else if strings.Contains(url.URL.String(), ".b32.i2p") {
		rq, err := http.NewRequest(url.Method, strings.TrimRight(url.URL.String(), "/"), url.Body)
		if addressBook.c, addressBook.err = Fatal(err, "addresshelper.go Request return error", "addresshelper.go Returning same request"); addressBook.c {
			Log("addresshelper.go base32 URL detected")
			Log("addresshelper.go no rewrite required")
			return rq, false
		}
		Log("addresshelper.go wierd base32 error you need to debug when you're not violently ill. Oh my god you forgot what this was now and you just noticed this log message.")
		return url, false
	} else if addressBook.assistant.CheckIsUnknownAddress(url.URL.String()) {
		if addressBook.Lookup(url.URL.String()) {
			return url, true
		}
	}
	rq, err := http.NewRequest(url.Method, url.URL.String(), url.Body)
	if addressBook.c, addressBook.err = Fatal(err, "addresshelper.go Request return error", "addresshelper.go Returning same request"); addressBook.c {
		Log("addresshelper.go no rewrite required")
		return rq, false
	}
	return url, false
}

func (addressBook *AddressHelper) checkAddPair(arg string) bool {
	addressBook.updateAh()
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

func (addressBook *AddressHelper) checkHasPair(arg string) bool {
	for _, pair := range addressBook.pairs {
		kvPair := strings.SplitN(pair, "=", 2)
		if kvPair != nil {
			if len(kvPair) == 2 {
				if kvPair[0] == arg {
					return true
				}
			}
		}
	}
	return false
}

func (addressBook *AddressHelper) Lookup(req string) bool {
	rv, jerr := addressBook.assistant.QueryHelper(req)
	if jerr != "jumperror" {
		addressBook.addPairString(rv)
		return true
	}
	Log("addressbook.go Jump URL not found")
	return false
}

func (addressBook *AddressHelper) addPairString(url string) {
	segments := strings.Split(strings.Replace(url, "http://", "", -1), "/")
	host := segments[0]
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

func (addressBook *AddressHelper) addPair(url *url.URL) {
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

func (addressBook *AddressHelper) getPair(url *url.URL) (string, string) {
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

func (addressBook *AddressHelper) getBase32(url *url.URL) (string, string) {
	for _, p := range addressBook.pairs {
		kv := strings.SplitN(p, "=", 2)
		if kv != nil {
			if len(kv) == 2 {
				if kv[0] == url.Host {
					b32, err := addressBook.converter.I2p64to32(kv[1])
					if addressBook.c, addressBook.err = Warn(err, "addresshelper.go Base32 conversion failure", "addresshelper.go Base32 converted"); addressBook.c {
						Log("addresshelper.go b32:", b32)
						return kv[0], b32
					}
				}
			}
		}
	}
	return "", ""
}

func (addressBook *AddressHelper) fileCheck(line string) bool {
	temp, err := ioutil.ReadFile(addressBook.bookPath)
	if addressBook.c, addressBook.err = Warn(err, "addresshelper.go File check error, handling:", "addresshelper.go Checking Addressbook file", addressBook.bookPath); addressBook.c {
		return !strings.Contains(string(temp), line)
	}
	return true
}

func (addressBook *AddressHelper) updateAh() {
	if addressBook.c, addressBook.err = exists(addressBook.bookPath); addressBook.c {
		addressBook.bookFile, addressBook.err = os.OpenFile(addressBook.bookPath, os.O_APPEND|os.O_WRONLY, 0755)
	} else {
		addressBook.bookFile, addressBook.err = os.OpenFile(addressBook.bookPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	}
	if addressBook.c, addressBook.err = Fatal(addressBook.err, "addresshelper.go File I/O errors", "addresshelper.go Addressbook file written"); addressBook.c {
		defer addressBook.bookFile.Close()
		if len(addressBook.pairs) > 1 {
			line := addressBook.pairs[len(addressBook.pairs)-1] + "\n"
			if addressBook.fileCheck(line) {
				addressBook.bookFile.WriteString(line)
			} else {
				Log("addresshelper.go Address already in Address Book")
			}
		}
	}
}

func NewAddressHelper(AddressHelperURL, samHost, samPort string) *AddressHelper {
	a, _ := NewAddressHelperFromOptions(
		SetAddressHelperURL(AddressHelperURL),
		SetAddressHelperHost(samHost),
		SetAddressHelperPort(samPort),
		SetAddressBookPath("addressbook.txt"),
	)
	return a
}

func NewAddressHelperFromOptions(opts ...func(*AddressHelper) error) (*AddressHelper, error) {
	var a AddressHelper
	a.addressHelperURL = "inr.i2p"
	a.samHostString = "127.0.0.1"
	a.samPortString = "7656"
	a.bookPath = "addressbook.txt"
	for _, o := range opts {
		if err := o(&a); err != nil {
			return nil, err
		}
	}
	a.assistant, a.err = i2paddresshelper.NewI2pAddressHelper(a.addressHelperURL, a.samHostString, a.samPortString, DEBUG)
	Fatal(a.err, "addresshelper.go failed to setup SAM bridge for addresshelper.", "addresshelper.go connecting to SAM bridge on:", a.addressHelperURL, a.samHostString, ":", a.samPortString)
	a.pairs = []string{}
	a.c = false
	return &a, a.err
}
