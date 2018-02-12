package main

import (
    //"bufio"
    "fmt"
	"io"
	//"log"
	"net/http"
    //"os"
    //"path/filepath"
    //"strings"
    //"strconv"
    //"syscall"
    //"net/url"

	//"github.com/eyedeekay/gosam"
)

type samHttpProxy struct {
    host string

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
		header.Del(h)
	}
}

func (proxy *samHttpProxy) copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func (proxy *samHttpProxy) prepare(){
    handle := &samHttpProxy{}
    if err := http.ListenAndServe(proxy.host, handle); err == nil {
        fmt.Println("Fatal Error: proxy not started")
    }
}

func (proxy *samHttpProxy) checkURLType(rW http.ResponseWriter, rq *http.Request) bool {
    fmt.Println(rq.RemoteAddr, " ", rq.Method, " ", rq.URL)
    if rq.URL.Scheme != "http" && rq.URL.Scheme != "https" {
	  	msg := "unsupported protocal scheme "+rq.URL.Scheme
		http.Error(rW, msg, http.StatusBadRequest)
		fmt.Println(msg)
		return false
	}else{
        return true
    }
}

func (proxy *samHttpProxy) ServeHTTP(rW http.ResponseWriter, rq *http.Request){
    fmt.Println("")
    if proxy.checkURLType(rW, rq) {
        client := &http.Client{}
        rq.RequestURI = ""
        //proxy.delHopHeaders(rq.Header)
        resp, err := client.Do(rq)
        if err != nil {
            http.Error(rW, "Server Error", http.StatusInternalServerError)
            fmt.Println("Fatal: ServeHTTP:", err)
        }
        defer resp.Body.Close()

        fmt.Println(rq.RemoteAddr, " ", resp.Status)

        proxy.delHopHeaders(resp.Header)

        proxy.copyHeader(rW.Header(), resp.Header)
        rW.WriteHeader(resp.StatusCode)
        io.Copy(rW, resp.Body)
    }
}

func createHttpProxy(proxAddr string, proxPort string, samStack samList) samHttpProxy {
    var samProxy samHttpProxy
    samProxy.host = proxAddr + ":" + proxPort
    return samProxy
}
