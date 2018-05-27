package dii2p

import (
    "strings"
)

func CheckURLType(request string) bool {

	Log(request)

	test := strings.Split(request, ".i2p")

	if len(test) < 2 {
		msg := "Non i2p domain detected. Skipping."
		Log(msg) //Outproxy support? Might be cool.
		return false
	} else {
		n := strings.Split(strings.Replace(strings.Replace(test[0], "https://", "", -1), "http://", "", -1), "/")
		if len(n) > 1 {
			msg := "Non i2p domain detected, possible attempt to impersonate i2p domain in path. Skipping."
			Log(msg) //Outproxy support? Might be cool. Riskier here.
			return false
		}
	}
	strings.Contains(request, "http")
	if !strings.Contains(request, "http") {
		if strings.Contains(request, "https") {
			msg := "Dropping https request for now, assumed attempt to get clearnet resource."
			Log(msg)
			return false
		} else {
			msg := "unsupported protocal scheme " + request
			Log(msg)
			return false
		}
	} else {
		return true
	}
}
