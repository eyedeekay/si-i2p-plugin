package dii2p

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type samHTTPProxy struct {
	Host        string
	client      *SamList
	transport   *http.Transport
	newHandle   *http.Server
	addressbook *addressHelper
	timeoutTime time.Duration
	keepAlives  bool
	err         error
	c           bool
}

func (proxy *samHTTPProxy) delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		Log("si-http-proxy.go Sanitizing headers: ", h, header.Get(h))
		header.Del(h)
	}
	if header.Get("User-Agent") != "MYOB/6.66 (AN/ON)" {
		header.Set("User-Agent", "MYOB/6.66 (AN/ON)")
	}
}

func (proxy *samHTTPProxy) copyHeader(dst, src http.Header) {
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

func (proxy *samHTTPProxy) prepare() {
	Log("si-http-proxy.go Initializing handler handle")
	if err := proxy.newHandle.ListenAndServe(); err != nil {
		Log("si-http-proxy.go Fatal Error: proxy not started")
	}
}

//export ServeHTTP
func (proxy *samHTTPProxy) ServeHTTP(rW http.ResponseWriter, rq *http.Request) {
	if &rq == nil {
		return
	}

	Log("si-http-proxy.go", rq.Host, " ", rq.RemoteAddr, " ", rq.Method, " ", rq.URL.String())

	if !CheckURLType(rq.URL.String()) {
		return
	}

	Log("si-http-proxy.go ", rq.URL.String())

	proxy.checkResponse(rW, rq)

}

func (proxy *samHTTPProxy) checkResponse(rW http.ResponseWriter, rq *http.Request) {
	if rq == nil {
		return
	}

	rq.RequestURI = ""

	req, ah := proxy.addressbook.checkAddressHelper(rq)

	//req.RequestURI = ""
	if proxy.keepAlives {
		req.Close = proxy.keepAlives
	}

	proxy.delHopHeaders(req.Header)

	Log("si-http-proxy.go Retrieving client")

	client, dir := proxy.client.sendClientRequestHTTP(req)

	if client != nil {
		Log("si-http-proxy.go Client was retrieved: ", dir)
		resp, doerr := proxy.Do(req, client, 0, ah)
		if proxy.c, proxy.err = Warn(doerr, "si-http-proxy.go Encountered an oddly formed response. Skipping.", "si-http-proxy.go Processing Response"); proxy.c {
			resp := proxy.client.copyRequest(req, resp, dir)
			proxy.printResponse(rW, resp)
			Log("si-http-proxy.go responded")
			return
		} else {
			if !strings.Contains(doerr.Error(), "malformed HTTP status code") && !strings.Contains(doerr.Error(), "use of closed network connection") {
				if resp != nil {
					resp := proxy.client.copyRequest(req, resp, dir)
					proxy.printResponse(rW, resp)
					return
				}
				Log("si-http-proxy.go status error", doerr.Error())
				return
			}
			Log("si-http-proxy.go status error", doerr.Error())
			return
		}
	} else {
		log.Println("si-http-proxy.go client retrieval error")
		return
	}
}

//export Do
func (proxy *samHTTPProxy) Do(req *http.Request, client *http.Client, x int, useah bool) (*http.Response, error) {
	req.RequestURI = ""

	resp, doerr := client.Do(req)

	if req.Close {
		Log("request must be closed")
	}

	if proxy.c, proxy.err = Warn(doerr, "si-http-proxy.go Response body error:", "si-http-proxy.go Read response body"); proxy.c {
		return resp, doerr
	} else {
		if useah {
			if strings.Contains(doerr.Error(), "Hostname error") {
				log.Println("Unknown Hostname")
				proxy.addressbook.Lookup(req.Host)
				requ, stage2 := proxy.addressbook.checkAddressHelper(req)
				if stage2 {
					log.Println("Redirecting", req.Host, "to", requ.Host)
					requ.RequestURI = ""
					return client.Do(requ)
				}
			} else {
				return client.Do(req)
			}
		} else {
			return client.Do(req)
		}
	}
	return resp, doerr
}

func (proxy *samHTTPProxy) printResponse(rW http.ResponseWriter, r *http.Response) {
	if r != nil {
		defer r.Body.Close()
		proxy.copyHeader(rW.Header(), r.Header)
		rW.WriteHeader(r.StatusCode)
		io.Copy(rW, r.Body)
		Log("si-http-proxy.go Response status:", r.Status)
	}
}

//export CreateHttpProxy
func CreateHttpProxy(proxAddr, proxPort, initAddress, addressHelperUrl string, samStack *SamList, timeoutTime int, keepAlives bool) *samHTTPProxy {
	var samProxy samHTTPProxy
	samProxy.Host = proxAddr + ":" + proxPort
	samProxy.keepAlives = keepAlives
	samProxy.addressbook = newAddressHelper(addressHelperUrl, samStack.samAddrString, samStack.samPortString)
	log.Println("si-http-proxy.go Starting HTTP proxy on:" + samProxy.Host)
	samProxy.client = samStack
	samProxy.timeoutTime = time.Duration(timeoutTime) * time.Minute
	samProxy.newHandle = &http.Server{
		Addr:         samProxy.Host,
		Handler:      &samProxy,
		ReadTimeout:  samProxy.timeoutTime,
		WriteTimeout: samProxy.timeoutTime,
	}
	log.Println("si-http-proxy.go Connected SAM isolation stack to the HTTP proxy server")
	go samProxy.prepare()
	log.Println("si-http-proxy.go HTTP Proxy prepared")
	return &samProxy
}
