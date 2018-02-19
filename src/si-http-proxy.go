package main

import (
    "log"
	"io"
	"net/http"
)

type samHttpProxy struct {
    host string
    client samList
    err error
}

var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
    "X-Forwarded-For",
}

func (proxy *samHttpProxy) delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
        log.Println("Sanitizing headers: " + h)
		header.Del(h)
	}
}

func (proxy *samHttpProxy) copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
            log.Println("Copying headers: " + k + "," + v )
			dst.Add(k, v)
		}
	}
}

func (proxy *samHttpProxy) prepare(){
    handle := &samHttpProxy{}
    log.Println("Initializing handler handle")
    if err := http.ListenAndServe(proxy.host, handle); err == nil {
        log.Println("Fatal Error: proxy not started")
    }
}

func (proxy *samHttpProxy) checkURLType(rW http.ResponseWriter, rq *http.Request) bool {
    log.Println(rq.RemoteAddr, " ", rq.Method, " ", rq.URL)
    if rq.URL.Scheme != "http" && rq.URL.Scheme != "https" {
	  	msg := "unsupported protocal scheme "+rq.URL.Scheme
		http.Error(rW, msg, http.StatusBadRequest)
		log.Println(msg)
		return false
	}else{
        return true
    }
}

func (proxy *samHttpProxy) ServeHTTP(rW http.ResponseWriter, rq *http.Request){
    if proxy.checkURLType(rW, rq) {
        log.Println("")
        //client := &http.Client{}
        rq.RequestURI = ""
        proxy.delHopHeaders(rq.Header)
        resp, err := proxy.client.sendClientRequestHttp(rq)
        if err != nil {
            http.Error(rW, "Http Proxy Server Error", http.StatusInternalServerError)
            log.Fatal("Fatal: ServeHTTP:", err)
        }

        log.Println(rq.RemoteAddr, " ", resp.Status)

        proxy.delHopHeaders(resp.Header)

        proxy.copyHeader(rW.Header(), resp.Header)
        rW.WriteHeader(resp.StatusCode)
        io.Copy(rW, resp.Body)
    }
}

func createHttpProxy(proxAddr string, proxPort string, samStack samList) samHttpProxy {
    var samProxy samHttpProxy
    samProxy.host = proxAddr + ":" + proxPort
    log.Println("Starting HTTP proxy on:" + samProxy.host)
    samProxy.client = samStack
    log.Println("Connected SAM isolation stack to the HTTP proxy server")
    go samProxy.prepare()
    log.Println("HTTP Proxy prepared")
    return samProxy
}
