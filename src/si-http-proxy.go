package dii2p

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

import (
	"github.com/eyedeekay/si-i2p-plugin/src/addresshelper"
	"github.com/eyedeekay/si-i2p-plugin/src/errors"
	"github.com/eyedeekay/si-i2p-plugin/src/helpers"
)

//SamHTTPProxy is an http proxy for making isolated SAM requests
type SamHTTPProxy struct {
	Addr        string
	client      *SamList
	transport   *http.Transport
	newHandle   *http.Server
	addressbook *dii2pah.AddressHelper
	timeoutTime time.Duration
	keepAlives  bool
	err         error
	c           bool
}

func (proxy *SamHTTPProxy) delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		dii2perrs.Log("si-http-proxy.go Sanitizing headers: ", h, header.Get(h))
		header.Del(h)
	}
	if header.Get("User-Agent") != "MYOB/6.66 (AN/ON)" {
		header.Set("User-Agent", "MYOB/6.66 (AN/ON)")
	}
}

func (proxy *SamHTTPProxy) copyHeader(dst, src http.Header) {
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

func (proxy *SamHTTPProxy) prepare() {
	dii2perrs.Log("si-http-proxy.go Initializing handler handle")
	if err := proxy.newHandle.ListenAndServe(); err != nil {
		dii2perrs.Log("si-http-proxy.go dii2perrs.Fatal Error: proxy not started")
	}
}

//export ServeHTTP
func (proxy *SamHTTPProxy) ServeHTTP(rW http.ResponseWriter, rq *http.Request) {
	if &rq == nil {
		return
	}

	dii2perrs.Log("si-http-proxy.go", rq.Host, " ", rq.RemoteAddr, " ", rq.Method, " ", rq.URL.String())

	if !dii2phelper.CheckURLType(rq.URL.String()) {
		return
	}

	dii2perrs.Log("si-http-proxy.go ", rq.URL.String())

	proxy.checkResponse(rW, rq)

}

func (proxy *SamHTTPProxy) checkResponse(rW http.ResponseWriter, rq *http.Request) {
	if rq == nil {
		return
	}

	rq.RequestURI = ""

	req, useAddressHelper := proxy.addressbook.CheckAddressHelper(rq)
	if useAddressHelper {
		dii2perrs.Log("si-http-proxy.go using jump helper")
	}

	req.RequestURI = ""
	if proxy.keepAlives {
		req.Close = proxy.keepAlives
	}

	proxy.delHopHeaders(req.Header)

	dii2perrs.Log("si-http-proxy.go Retrieving client")

	client, dir := proxy.client.sendClientRequestHTTP(req)

	if client != nil {
		dii2perrs.Log("si-http-proxy.go Client was retrieved: ", dir)
		resp, doerr := proxy.Do(req, client, 0)
		if proxy.c, proxy.err = dii2perrs.Warn(doerr, "si-http-proxy.go Encountered an oddly formed response. Skipping.", "si-http-proxy.go Processing Response"); proxy.c {
			resp := proxy.client.copyRequest(req, resp, dir)
			proxy.printResponse(rW, resp)
			dii2perrs.Log("si-http-proxy.go responded")
			return
		}
		if !strings.Contains(doerr.Error(), "malformed HTTP status code") && !strings.Contains(doerr.Error(), "use of closed network connection") {
			if resp != nil {
				resp := proxy.client.copyRequest(req, resp, dir)
				proxy.printResponse(rW, resp)
				return
			}
			dii2perrs.Log("si-http-proxy.go status error", doerr.Error())
			return
		}
		dii2perrs.Log("si-http-proxy.go status error", doerr.Error())
		return
	}
	log.Println("si-http-proxy.go client retrieval error")
	return
}

//Do does a request
func (proxy *SamHTTPProxy) Do(req *http.Request, client *http.Client, x int) (*http.Response, error) {
	req.RequestURI = ""

	resp, doerr := client.Do(req)

	if req.Close {
		dii2perrs.Log("request must be closed")
	}

	if proxy.c, proxy.err = dii2perrs.Warn(doerr, "si-http-proxy.go Response body error:", "si-http-proxy.go Read response body"); proxy.c {
		return resp, doerr
	}

	return resp, doerr
}

func (proxy *SamHTTPProxy) printResponse(rW http.ResponseWriter, r *http.Response) {
	if r != nil {
		defer r.Body.Close()
		proxy.copyHeader(rW.Header(), r.Header)
		rW.WriteHeader(r.StatusCode)
		io.Copy(rW, r.Body)
		dii2perrs.Log("si-http-proxy.go Response status:", r.Status)
	}
}

//CreateHTTPProxy creates a SamHTTPProxy
func CreateHTTPProxy(proxAddr, proxPort, initAddress, ahAddr, ahPort, addressHelperURL string, samStack *SamList, timeoutTime int, keepAlives bool) *SamHTTPProxy {
	var samProxy SamHTTPProxy
	samProxy.Addr = proxAddr + ":" + proxPort
	samProxy.keepAlives = keepAlives
	samProxy.addressbook = dii2pah.NewAddressHelper(addressHelperURL, ahAddr, ahPort)
	log.Println("si-http-proxy.go Starting HTTP proxy on:" + samProxy.Addr)
	samProxy.client = samStack
	samProxy.timeoutTime = time.Duration(timeoutTime) * time.Minute
	samProxy.newHandle = &http.Server{
		Addr:         samProxy.Addr,
		Handler:      &samProxy,
		ReadTimeout:  samProxy.timeoutTime,
		WriteTimeout: samProxy.timeoutTime,
	}
	log.Println("si-http-proxy.go Connected SAM isolation stack to the HTTP proxy server")
	go samProxy.prepare()
	log.Println("si-http-proxy.go HTTP Proxy prepared")
	return &samProxy
}
