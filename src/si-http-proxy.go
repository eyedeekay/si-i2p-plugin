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
		proxy.Log("Sanitizing headers: " + h)
		header.Del(h)
	}
	header.Set("User-Agent", "MYOB/6.66 (AN/ON)")
}

func (proxy *samHttpProxy) copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			proxy.Log("Copying headers: " + k + "," + v)
            if dst.Get(k) != "" {
                dst.Set(k, v)
            }else{
                dst.Add(k, v)
            }
		}
	}
}

func (proxy *samHttpProxy) prepare() {
	proxy.Log("Initializing handler handle")
	if err := http.ListenAndServe(proxy.host, proxy.handle); err != nil {
		proxy.Log("Fatal Error: proxy not started")
	}
}

func (proxy *samHttpProxy) checkURLType(rW http.ResponseWriter, rq *http.Request) bool {

	log.Println(rq.RemoteAddr, " ", rq.Method, " ", rq.URL)

	test := strings.Split(rq.URL.String(), ".i2p")

	if len(test) < 2 {
		msg := "Non i2p domain detected. Skipping."
		proxy.Log(msg) //Outproxy support? Might be cool.
		http.Error(rW, "Http Proxy Server Error", http.StatusInternalServerError)
		return false
	} else {
		n := strings.Split(strings.Replace(strings.Replace(test[0], "https://", "", -1), "http://", "", -1), "/")
		if len(n) > 1 {
			msg := "Non i2p domain detected, possible attempt to impersonate i2p domain in path. Skipping."
			proxy.Log(msg) //Outproxy support? Might be cool. Riskier here.
			http.Error(rW, "Http Proxy Server Error", http.StatusInternalServerError)
			return false
		}
	}
	if rq.URL.Scheme != "http" {
		if rq.URL.Scheme == "https" {
			msg := "Dropping https request for now, assumed attempt to get clearnet resource." + rq.URL.Scheme
			proxy.Log(msg)
			http.Error(rW, "Http Proxy Server Error", http.StatusInternalServerError)
			return false
		} else {
			msg := "unsupported protocal scheme " + rq.URL.Scheme
			proxy.Log(msg)
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

	proxy.Log(rq.URL.String())

	rq.RequestURI = ""
	proxy.delHopHeaders(rq.Header)

	client, dir := proxy.client.sendClientRequestHttp(rq)

	proxy.Log("Retrieving client")

	if client != nil {
		proxy.Log("Client was retrieved: ", dir)

		resp, err := client.Do(rq)
		if err != nil {
			proxy.Warn(err, "Encountered an oddly formed response. Skipping.", "Processing Response")
			http.Error(rW, "Http Proxy Server Error", http.StatusInternalServerError)
		} else {

			r := proxy.client.copyRequest(rq, resp, dir)

			if r != nil {
				proxy.Log("SAM-Provided Tunnel Address:", rq.RemoteAddr)
				proxy.Log("Response Status:", r.Status)

				proxy.delHopHeaders(r.Header)

				proxy.copyHeader(rW.Header(), r.Header)


                log.Println("Response status:", r.StatusCode)

                rW.WriteHeader(r.StatusCode)
                read, err := ioutil.ReadAll(r.Body)

                if proxy.Warn(err, "Response body error:", "Read response body"){
                    io.Copy(rW, ioutil.NopCloser(bytes.NewBuffer(read)))
                }
			}
		}
	} else {
		proxy.Log(dir)
	}
}

func (proxy *samHttpProxy) Log(msg ...string) {
	if verbose {
		log.Println("LOG: ", msg)
	}
}

func (proxy *samHttpProxy) Warn(err error, errmsg string, msg ...string) bool {
	log.Println(msg)
	if err != nil {
		log.Println("WARN: ", err)
		return false
	}
	proxy.err = nil
	return true
}

func (proxy *samHttpProxy) Fatal(err error, errmsg string, msg ...string) {
	if err != nil {
		proxy.err = err
		defer proxy.client.cleanupClient()
		log.Fatal("Fatal: ", err)
	}
}

func createHttpProxy(proxAddr string, proxPort string, samStack *samList, initAddress string) *samHttpProxy {
	var samProxy samHttpProxy
	samProxy.host = proxAddr + ":" + proxPort
	log.Println("Starting HTTP proxy on:" + samProxy.host)
	samProxy.client = samStack
	samProxy.handle = &samProxy
	log.Println("Connected SAM isolation stack to the HTTP proxy server")
	go samProxy.prepare()
	log.Println("HTTP Proxy prepared")
	return &samProxy
}
