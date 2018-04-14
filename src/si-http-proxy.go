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
	//handle      *samHttpProxy
    newHandle   *http.Server
	addressbook *addressHelper
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
	Log("si-http-proxy.go", rq.Host, " ", rq.RemoteAddr, " ", rq.Method, " ", rq.URL.String())

	if !proxy.checkURLType(rW, rq) {
		return
	}

	Log("si-http-proxy.go ", rq.URL.String())
	rq.RequestURI = ""
    //rq.Close = true

	req, need := proxy.addressbook.checkAddressHelper(*rq)

	if req == nil {
		return
	}
	proxy.delHopHeaders(req.Header)

	var client *http.Client
	var dir, base64 string

	Log("si-http-proxy.go Retrieving client")
	if need {
		_, base64 = proxy.addressbook.getPair(req.URL)
		client, dir = proxy.client.sendClientRequestBase64Http(req, base64)
	} else {
		client, dir = proxy.client.sendClientRequestHttp(req)
	}

	if client != nil {
		Log("si-http-proxy.go Client was retrieved: ", dir)
		resp, err := client.Do(req)
		if proxy.c, proxy.err = Warn(err, "si-http-proxy.go Encountered an oddly formed response. Skipping.", "si-http-proxy.go Processing Response"); !proxy.c {
            if resp != nil {
                proxy.copyHeader(rW.Header(), resp.Header)
                read, err := ioutil.ReadAll(resp.Body)
                if proxy.c, proxy.err = Warn(err, "si-http-proxy.go Response body error:", "si-http-proxy.go Read response body"); proxy.c {
                    resp.Body.Close()
                    io.Copy(rW, ioutil.NopCloser(bytes.NewBuffer(read)))
                }
            }
            rW.WriteHeader(resp.StatusCode)
			return
		} else {
			r := proxy.client.copyRequest(req, resp, dir, base64)
            //r.Body.Close()
			if r != nil {
				Log("si-http-proxy.go SAM-Provided Tunnel Address:", req.RemoteAddr)
				Log("si-http-proxy.go Response Status:", r.Status)
				proxy.copyHeader(rW.Header(), r.Header)
				if r.StatusCode >= 200 {
					//if r.StatusCode == 301 {
						//Log("si-http-proxy.go Detected redirect.")
                        //return
					//}
					if r.StatusCode < 309 {
						rW.WriteHeader(r.StatusCode)
						read, err := ioutil.ReadAll(r.Body)
						if proxy.c, proxy.err = Warn(err, "si-http-proxy.go Response body error:", "si-http-proxy.go Read response body"); proxy.c {
                            r.Body.Close()
							io.Copy(rW, ioutil.NopCloser(bytes.NewBuffer(read)))
						}
						return
					}
					rW.WriteHeader(r.StatusCode)
					log.Println("si-http-proxy.go Response status:", r.StatusCode)
                    return
				} else {
					rW.WriteHeader(r.StatusCode)
					log.Println("si-http-proxy.go Response status:", r.StatusCode)
					return
				}
                rW.WriteHeader(r.StatusCode)
				log.Println("si-http-proxy.go Response status:", r.StatusCode)
				return
			}
		}
	} else {
		Log(dir)
	}
}

func createHttpProxy(proxAddr string, proxPort string, samStack *samList, initAddress string) *samHttpProxy {
	var samProxy samHttpProxy
	samProxy.host = proxAddr + ":" + proxPort
	samProxy.addressbook = newAddressHelper()
	log.Println("si-http-proxy.go Starting HTTP proxy on:" + samProxy.host)
	samProxy.client = samStack
	samProxy.newHandle = &http.Server {
        Addr: samProxy.host,
        Handler: &samProxy,
        ReadTimeout: time.Duration(10 * time.Second),
        WriteTimeout: time.Duration(10 * time.Second),
    }
	log.Println("si-http-proxy.go Connected SAM isolation stack to the HTTP proxy server")
	go samProxy.prepare()
	log.Println("si-http-proxy.go HTTP Proxy prepared")
	return &samProxy
}
