package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type samHttpProxy struct {
	host      string
	client    *samList
	transport *http.Transport
	handle    *samHttpProxy
	err       error
	c         bool
}

var hopHeaders = []string{
	"Accept",
	"Accept-Encoding",
	//"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Proxy-Connection",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
	"X-Forwarded-For",
}

func (proxy *samHttpProxy) delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		Log("si-http-proxy.go Sanitizing headers: " + h)
		header.Del(h)
	}
	header.Set("User-Agent", "MYOB/6.66 (AN/ON)")
}

func (proxy *samHttpProxy) copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			Log("si-http-proxy.go Copying headers: " + k + "," + v)
			if dst.Get(k) != "" {
				dst.Set(k, v)
			} else {
				dst.Add(k, v)
			}
		}
	}
}

func (proxy *samHttpProxy) prepare() {
	Log("si-http-proxy.go Initializing handler handle")
	if err := http.ListenAndServe(proxy.host, proxy.handle); err != nil {
		Log("si-http-proxy.go Fatal Error: proxy not started")
	}
}

func (proxy *samHttpProxy) checkURLType(rW http.ResponseWriter, rq *http.Request) bool {

	log.Println(rq.RemoteAddr, " ", rq.Method, " ", rq.URL)

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
	log.Println(rq.RemoteAddr, " ", rq.Method, " ", rq.URL)

	if !proxy.checkURLType(rW, rq) {
		return
	}

	Log(rq.URL.String())

	rq.RequestURI = ""
	proxy.delHopHeaders(rq.Header)

	client, dir := proxy.client.sendClientRequestHttp(rq)

	Log("si-http-proxy.go Retrieving client")

	if client != nil {
		Log("si-http-proxy.go Client was retrieved: ", dir)

		resp, err := client.Do(rq)
		if proxy.c, proxy.err = Warn(err, "si-http-proxy.go Encountered an oddly formed response. Skipping.", "si-http-proxy.go Processing Response"); !proxy.c {
			http.Error(rW, "Http Proxy Server Error", http.StatusInternalServerError)
		} else {

			r := proxy.client.copyRequest(rq, resp, dir)

			if r != nil {
				Log("si-http-proxy.go SAM-Provided Tunnel Address:", rq.RemoteAddr)
				Log("si-http-proxy.go Response Status:", r.Status)

				proxy.delHopHeaders(r.Header)

				proxy.copyHeader(rW.Header(), r.Header)


                if r.StatusCode >= 200 {
                    if r.StatusCode == 301 {
                        Log("si-http-proxy.go Detected redirect.")
                    }else if r.StatusCode < 301 {
                        rW.WriteHeader(r.StatusCode)
                        read, err := ioutil.ReadAll(r.Body)
                        if proxy.c, proxy.err = Warn(err, "si-http-proxy.go Response body error:", "si-http-proxy.go Read response body"); proxy.c {
                            io.Copy(rW, ioutil.NopCloser(bytes.NewBuffer(read)))
                        }
                    }else{
                        rW.WriteHeader(r.StatusCode)
                        log.Println("si-http-proxy.go Response status:", r.StatusCode)
                    }
                }else{
                    rW.WriteHeader(r.StatusCode)
                    log.Println("si-http-proxy.go Response status:", r.StatusCode)
                }
			}
		}
	} else {
		Log(dir)
	}
}

func createHttpProxy(proxAddr string, proxPort string, samStack *samList, initAddress string) *samHttpProxy {
	var samProxy samHttpProxy
	samProxy.host = proxAddr + ":" + proxPort
	log.Println("si-http-proxy.go Starting HTTP proxy on:" + samProxy.host)
	samProxy.client = samStack
	samProxy.handle = &samProxy
	log.Println("si-http-proxy.go Connected SAM isolation stack to the HTTP proxy server")
	go samProxy.prepare()
	log.Println("si-http-proxy.go HTTP Proxy prepared")
	return &samProxy
}
