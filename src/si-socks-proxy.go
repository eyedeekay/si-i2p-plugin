package dii2p

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/armon/go-socks5"
)

type samSOCKSProxy struct {
	Addr        string
	client      *SamList
	transport   *http.Transport
	newHandle   *socks5.Server
	addressbook *AddressHelper
	timeoutTime time.Duration
	keepAlives  bool
	err         error
	c           bool
}

func (proxy *samSOCKSProxy) delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		Log("si-socks-proxy.go Sanitizing headers: ", h, header.Get(h))
		header.Del(h)
	}
	if header.Get("User-Agent") != "MYOB/6.66 (AN/ON)" {
		header.Set("User-Agent", "MYOB/6.66 (AN/ON)")
	}
}

func (proxy *samSOCKSProxy) copyHeader(dst, src http.Header) {
	if dst != nil && src != nil {
		for k, vv := range src {
			if vv != nil {
				for _, v := range vv {
					if v != "" {
						Log("si-socks-proxy.go Copying headers: " + k + "," + v)
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

func (proxy *samSOCKSProxy) prepare() {
	Log("si-socks-proxy.go Initializing handler handle")
	if err := proxy.newHandle.ListenAndServe("tcp", proxy.Addr); err != nil {
		Log("si-socks-proxy.go Fatal Error: proxy not started")
	}
}

func (proxy *samSOCKSProxy) ServeSOCKS(rW http.ResponseWriter, rq *http.Request) {
	if &rq == nil {
		return
	}

	Log("si-socks-proxy.go", rq.Host, " ", rq.RemoteAddr, " ", rq.Method, " ", rq.URL.String())

	if !CheckURLType(rq.URL.String()) {
		return
	}

	Log("si-socks-proxy.go ", rq.URL.String())

	proxy.checkResponse(rW, rq)

}

func (proxy *samSOCKSProxy) checkResponse(rW http.ResponseWriter, rq *http.Request) {
	if rq == nil {
		return
	}

	rq.RequestURI = ""

	req, _ := proxy.addressbook.checkAddressHelper(rq)

	req.RequestURI = ""
	if proxy.keepAlives {
		req.Close = proxy.keepAlives
	}

	proxy.delHopHeaders(req.Header)

	Log("si-socks-proxy.go Retrieving client")

	client, dir := proxy.client.sendClientRequestHTTP(req)

	time.Sleep(1 * time.Second)

	if client != nil {
		Log("si-socks-proxy.go Client was retrieved: ", dir)
		resp, doerr := proxy.Do(req, client, 0)
		if proxy.c, proxy.err = Warn(doerr, "si-socks-proxy.go Encountered an oddly formed response. Skipping.", "si-socks-proxy.go Processing Response"); proxy.c {
			resp := proxy.client.copyRequest(req, resp, dir)
			proxy.printResponse(rW, resp)
			Log("si-socks-proxy.go responded")
			return
		} else {
			if !strings.Contains(doerr.Error(), "malformed HTTP status code") && !strings.Contains(doerr.Error(), "use of closed network connection") {
				if resp != nil {
					resp := proxy.client.copyRequest(req, resp, dir)
					proxy.printResponse(rW, resp)
				}
				Log("si-socks-proxy.go status error", doerr.Error())
				return
			}
			Log("si-socks-proxy.go status error", doerr.Error())
			return
		}
	} else {
		log.Println("si-socks-proxy.go client retrieval error")
		return
	}
}

func (proxy *samSOCKSProxy) Do(req *http.Request, client *http.Client, x int) (*http.Response, error) {
	req.RequestURI = ""

	resp, doerr := client.Do(req)

	if req.Close {
		Log("request must be closed")
	}

	if proxy.c, proxy.err = Warn(doerr, "si-socks-proxy.go Response body error:", "si-socks-proxy.go Read response body"); proxy.c {
		return resp, doerr
	} else {
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
	}
	return resp, doerr
}

func (proxy *samSOCKSProxy) printResponse(rW http.ResponseWriter, r *http.Response) {
	if r != nil {
		defer r.Body.Close()
		proxy.copyHeader(rW.Header(), r.Header)
		rW.WriteHeader(r.StatusCode)
		io.Copy(rW, r.Body)
		//r.Body.Close()
		Log("si-socks-proxy.go Response status:", r.Status)
	}
}

func CreateSOCKSProxy(proxAddr, proxPort, initAddress, addressHelperURL string, samStack *SamList, timeoutTime int, keepAlives bool) *samSOCKSProxy {
	var samProxy samSOCKSProxy
	samProxy.Addr = proxAddr + ":" + proxPort
	samProxy.keepAlives = keepAlives
	samProxy.addressbook = NewAddressHelper(addressHelperURL, samStack.samAddrString, samStack.samPortString)
	log.Println("si-socks-proxy.go Starting SOCKS proxy on:" + samProxy.Addr)
	samProxy.client = samStack
	samProxy.timeoutTime = time.Duration(timeoutTime) * time.Minute
	conf := &socks5.Config{}
	samProxy.newHandle, samProxy.err = socks5.New(conf)
	Fatal(samProxy.err, "si-socks-proxy.go SOCKS proxy creation error", "si-socks-proxy.go SOCKS proxy created")
	log.Println("si-socks-proxy.go Connected SAM isolation stack to the SOCKS proxy server")
	go samProxy.prepare()
	log.Println("si-socks-proxy.go SOCKS Proxy prepared")
	return &samProxy
}
