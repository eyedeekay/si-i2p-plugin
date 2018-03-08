package main

import (
    "log"
	"io"
	"net/http"
)

type samHttpProxy struct {
    host string
    client *samList
    transport *http.Transport
    handle *samHttpProxy
    err error
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
        log.Println("Sanitizing headers: " + h)
		header.Del(h)
	}
    header.Set("User-Agent", "MYOB/6.66 (AN/ON)")
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
    log.Println("Initializing handler handle")
    if err := http.ListenAndServe(proxy.host, proxy.handle); err != nil {
        log.Println("Fatal Error: proxy not started")
    }
}

func (proxy *samHttpProxy) checkURLType(rW http.ResponseWriter, rq *http.Request) bool {
    log.Println(rq.RemoteAddr, " ", rq.Method, " ", rq.URL)
    if rq.URL.Scheme != "http" && rq.URL.Scheme != "https" {
	  	msg := "unsupported protocal scheme " + rq.URL.Scheme
		http.Error(rW, msg, http.StatusBadRequest)
		log.Println(msg)
		return false
	}else{
        return true
    }
}

func (proxy *samHttpProxy) ServeHTTP(rW http.ResponseWriter, rq *http.Request){
    log.Println(rq.RemoteAddr, " ", rq.Method, " ", rq.URL)

    if ! proxy.checkURLType(rW, rq) {
        return
    }
    log.Println(rq.URL.String())

    rq.RequestURI = ""
    proxy.delHopHeaders(rq.Header)

    client, dir := proxy.client.sendClientRequestHttp(rq)

    log.Println("Client was retrieved: ", dir)

    resp, err := client.Do(rq)
    if err != nil {
        log.Println("Encountered an oddly formed response. Skipping.")
        //http.Error(rW, "Http Proxy Server Error", http.StatusInternalServerError)
    }else{

        r := proxy.client.copyRequest(rq, resp, dir)


        if r != nil {
            log.Println("SAM-Provided Tunnel Address:", rq.RemoteAddr)
            log.Println("Response Status:", r.Status)

            proxy.delHopHeaders(r.Header)

            proxy.copyHeader(rW.Header(), r.Header)
            rW.WriteHeader(r.StatusCode)
            io.Copy(rW, r.Body)
        }
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
