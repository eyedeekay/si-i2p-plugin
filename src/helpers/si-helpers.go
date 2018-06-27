package dii2phelper

import (
	"strings"
)

import (
	"github.com/eyedeekay/si-i2p-plugin/src/errors"
)

//CheckURLType assures a URL is intended for i2p
func CheckURLType(request string) bool {

	dii2perrs.Log(request)

	test := strings.Split(request, ".i2p")

	if len(test) < 2 {
		msg := "Non i2p domain detected. Skipping."
		dii2perrs.Log(msg) //Outproxy support? Might be cool.
		return false
	}
	n := strings.Split(strings.Replace(strings.Replace(test[0], "https://", "", -1), "http://", "", -1), "/")
	if len(n) > 1 {
		msg := "Non i2p domain detected, possible attempt to impersonate i2p domain in path. Skipping."
		dii2perrs.Log(msg) //Outproxy support? Might be cool. Riskier here.
		return false
	}
	strings.Contains(request, "http")
	if !strings.Contains(request, "http") {
		if strings.Contains(request, "https") {
			msg := "Dropping https request for now, assumed attempt to get clearnet resource."
			dii2perrs.Log(msg)
			return false
		}
		msg := "unsupported protocal scheme " + request
		dii2perrs.Log(msg)
		return false
	}
	return true
}

//CleanURL will at some point be replaced.
func CleanURL(request string) (string, string) {
	dii2perrs.Log("sam-http.go cleanURL Request " + request)
	//url := strings.Replace(request, "http://", "", -1)
	var url string
	if !strings.HasPrefix(request, "http://") {
		url = "http://" + request
	}
	url = request

	if strings.HasSuffix(url, ".i2p") {
		url = url + "/"
	}

	host := strings.Replace(
		strings.SplitAfter(url, ".i2p/")[0],
		"http://",
		"",
		-1,
	)

	if strings.HasSuffix(host, ".i2p/") {
		host = strings.TrimSuffix(host, "/")
	}
	if strings.HasSuffix(url, ".i2p/") {
		url = strings.TrimSuffix(url, "/")
	}

	dii2perrs.Log("sam-http.go cleanURL Request URL " + url)
	dii2perrs.Log("sam-http.go cleanURL Request Host ", host)

	return host, url
}
