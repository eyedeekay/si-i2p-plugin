package dii2p

import (
	"net/http"
)

import (
	"github.com/eyedeekay/si-i2p-plugin/src/errors"
)

var hopHeaders = []string{
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Proxy-Connection",
	"X-Forwarded-For",
	"Accept-Language",
}

func delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		dii2perrs.Log("si-http-proxy.go Sanitizing headers: ", h, header.Get(h))
		header.Del(h)
	}
	if header.Get("User-Agent") != "MYOB/6.66 (AN/ON)" {
		header.Set("User-Agent", "MYOB/6.66 (AN/ON)")
	}
}

func copyHeader(dst, src http.Header) {
	if dst != nil && src != nil {
		for k, vv := range src {
			if vv != nil {
				for _, v := range vv {
					if v != "" {
						dii2perrs.Log("si-http-proxy.go Copying headers: " + k + "," + v)
						if dst.Get(k) != "" {
							dst.Set(k, v)
						} else {
							dst.Add(k, v)
						}
					}
				}
			}
		}
	}
}
