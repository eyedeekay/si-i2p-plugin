package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type samHttpProxy struct {
	host        string
	client      *samList
	transport   *http.Transport
	newHandle   *http.Server
	addressbook *addressHelper
	timeoutTime time.Duration
	err         error
	c           bool
}

var hopHeaders = []string{
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Proxy-Connection",
	"X-Forwarded-For",
}

func (proxy *samHttpProxy) delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		Log("si-http-proxy.go Sanitizing headers: ", h, header.Get(h))
		header.Del(h)
	}
	if header.Get("User-Agent") != "MYOB/6.66 (AN/ON)" {
		header.Set("User-Agent", "MYOB/6.66 (AN/ON)")
	}
}

func (proxy *samHttpProxy) copyHeader(dst, src http.Header) {
	if dst != nil && src != nil {
		for k, vv := range src {
			if vv != nil {
				for _, v := range vv {
					if v != "" {
						Log("si-http-proxy.go Copying headers: " + k + "," + v)
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

func (proxy *samHttpProxy) prepare() {
	Log("si-http-proxy.go Initializing handler handle")
	if err := proxy.newHandle.ListenAndServe(); err != nil {
		//if err:= http.ListenAndServe(proxy.host, proxy.handle); err != nil
		Log("si-http-proxy.go Fatal Error: proxy not started")
	}
}

func (proxy *samHttpProxy) checkURLType(rW http.ResponseWriter, rq *http.Request) bool {

	Log("si-http-proxy.go", rq.Host, " ", rq.RemoteAddr, " ", rq.Method, " ", rq.URL.String())

	test := strings.Split(rq.URL.String(), ".i2p")

	if len(test) < 2 {
		msg := "Non i2p domain detected. Skipping."
		Log(msg) //Outproxy support? Might be cool.
		http.Error(rW, "Http Proxy Server Error", http.StatusInternalServerError)
		return false
	} else {
		n := strings.Split(strings.Replace(strings.Replace(test[0], "https://", "", -1), "http://", "", -1), "/")
		if len(n) > 1 {
			msg := "Non i2p domain detected, possible attempt to impersonate i2p domain in path. Skipping."
			Log(msg) //Outproxy support? Might be cool. Riskier here.
			http.Error(rW, "Http Proxy Server Error", http.StatusInternalServerError)
			return false
		}
	}
	if rq.URL.Scheme != "http" {
		if rq.URL.Scheme == "https" {
			msg := "Dropping https request for now, assumed attempt to get clearnet resource." + rq.URL.Scheme
			Log(msg)
			http.Error(rW, "Http Proxy Server Error", http.StatusInternalServerError)
			return false
		} else {
			msg := "unsupported protocal scheme " + rq.URL.Scheme
			Log(msg)
			http.Error(rW, "Http Proxy Server Error", http.StatusInternalServerError)
			return false
		}
	} else {
		return true
	}
}

func (proxy *samHttpProxy) ServeHTTP(rW http.ResponseWriter, rq *http.Request) {
	if &rq == nil {
		return
	}

	Log("si-http-proxy.go", rq.Host, " ", rq.RemoteAddr, " ", rq.Method, " ", rq.URL.String())

	if !proxy.checkURLType(rW, rq) {
		return
	}

	Log("si-http-proxy.go ", rq.URL.String())
	rq.RequestURI = ""

	req, _ := proxy.addressbook.checkAddressHelper(rq)

	req.RequestURI = ""
	//req.Close = true

	proxy.checkResponse(rW, req)

}

func (proxy *samHttpProxy) checkResponse(rW http.ResponseWriter, req *http.Request) {
	if &req == nil {
		return
	}

	req.RequestURI = ""

	proxy.delHopHeaders(req.Header)

	Log("si-http-proxy.go Retrieving client")

	client, dir := proxy.client.sendClientRequestHttp(req)

	if client != nil {
		Log("si-http-proxy.go Client was retrieved: ", dir)
		resp, doerr := proxy.Do(req, client, 0)
		if proxy.c, proxy.err = Warn(doerr, "si-http-proxy.go Encountered an oddly formed response. Skipping.", "si-http-proxy.go Processing Response"); !proxy.c {
			if !strings.Contains(doerr.Error(), "malformed HTTP status code") && !strings.Contains(doerr.Error(), "use of closed network connection") {
				proxy.printResponse(rW, resp)
			}
			return
		} else {
			r := proxy.client.copyRequest(req, resp, dir)
			proxy.printResponse(rW, r)
			return
		}
	} else {
		log.Println("si-http-proxy.go client retrieval error")
		return
	}
}

func (proxy *samHttpProxy) Do(req *http.Request, client *http.Client, x int) (*http.Response, error) {
	resp, doerr := client.Do(req)

	if req.Close {
		log.Println("request must be closed")
	}

	if resp != nil {
		if proxy.c, proxy.err = Warn(doerr, "si-http-proxy.go Response body error:", "si-http-proxy.go Read response body"); proxy.c {
			return resp, doerr
		} else {
			if strings.Contains(doerr.Error(), "Hostname error") {
				proxy.addressbook.Lookup(req.Host)
				//return proxy.reDo(req, client, x)
                return client.Do(req)
			}
			return resp, doerr
		}
	}
	return client.Do(req) // proxy.reDo(req, client, x+1)
}

func (proxy *samHttpProxy) reDo(req *http.Request, client *http.Client, x int) (*http.Response, error) {
	y := x * 1
	if x < 3 {
		time.Sleep(time.Duration(y) * time.Second)
		log.Println("si-http-proxy.go retrying attempt:", x, "for", req.URL.String(), "after", y, "seconds")
		return proxy.Do(req, client, x)
	} else {
		resp, err := client.Do(req)
		return resp, err
	}
}

func (proxy *samHttpProxy) printResponse(rW http.ResponseWriter, r *http.Response) {
	if r != nil {
		rW.WriteHeader(r.StatusCode)
		readstring, readerr := ioutil.ReadAll(r.Body)
		if proxy.c, proxy.err = Warn(readerr, "si-http-proxy.go Response body error:", "si-http-proxy.go Read response body"); proxy.c {
			proxy.copyHeader(rW.Header(), r.Header)
			io.Copy(rW, ioutil.NopCloser(bytes.NewBuffer(readstring)))
		}
		Log("si-http-proxy.go Response status:", r.Status)
		return
	}
}

func createHttpProxy(proxAddr, proxPort string, samStack *samList, addressHelperUrl, initAddress string, timeoutTime int) *samHttpProxy {
	var samProxy samHttpProxy
	samProxy.host = proxAddr + ":" + proxPort
	samProxy.addressbook = newAddressHelper(addressHelperUrl, samStack.samAddrString, samStack.samPortString)
	log.Println("si-http-proxy.go Starting HTTP proxy on:" + samProxy.host)
	samProxy.client = samStack
	samProxy.timeoutTime = time.Duration(timeoutTime) * time.Minute
	samProxy.newHandle = &http.Server{
		Addr:         samProxy.host,
		Handler:      &samProxy,
		ReadTimeout:  samProxy.timeoutTime,
		WriteTimeout: samProxy.timeoutTime,
	}
	log.Println("si-http-proxy.go Connected SAM isolation stack to the HTTP proxy server")
	go samProxy.prepare()
	log.Println("si-http-proxy.go HTTP Proxy prepared")
	return &samProxy
}
