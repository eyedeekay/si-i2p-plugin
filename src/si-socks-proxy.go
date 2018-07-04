package dii2p

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/armon/go-socks5"
)

import (
	"github.com/eyedeekay/si-i2p-plugin/src/resolver"
	//"github.com/eyedeekay/si-i2p-plugin/src/addresshelper"
	"github.com/eyedeekay/si-i2p-plugin/src/client"
	"github.com/eyedeekay/si-i2p-plugin/src/errors"
	"github.com/eyedeekay/si-i2p-plugin/src/helpers"
)

// SamSOCKSProxy is a SOCKS proxy that automatically isolates per-destination
type SamSOCKSProxy struct {
	Addr        string
	client      *dii2pmain.SamList
	transport   *http.Transport
	newHandle   *socks5.Server
	addressbook *jumpresolver.JumpResolver
	timeoutTime time.Duration
	keepAlives  bool
	err         error
	c           bool
}

func (proxy *SamSOCKSProxy) prepare() {
	dii2perrs.Log("si-socks-proxy.go Initializing handler handle")
	if err := proxy.newHandle.ListenAndServe("tcp", proxy.Addr); err != nil {
		dii2perrs.Log("si-socks-proxy.go dii2perrs.Fatal Error: proxy not started")
	}
}

// ServeSOCKS Starts serving a SOCKS proxy
func (proxy *SamSOCKSProxy) ServeSOCKS(rW http.ResponseWriter, rq *http.Request) {
	if &rq == nil {
		return
	}

	dii2perrs.Log("si-socks-proxy.go", rq.Host, " ", rq.RemoteAddr, " ", rq.Method, " ", rq.URL.String())

	if !dii2phelper.CheckURLType(rq.URL.String()) {
		return
	}

	dii2perrs.Log("si-socks-proxy.go ", rq.URL.String())

	proxy.checkResponse(rW, rq)

}

func (proxy *SamSOCKSProxy) checkResponse(rW http.ResponseWriter, rq *http.Request) {
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

	delHopHeaders(req.Header)

	dii2perrs.Log("si-socks-proxy.go Retrieving client")

	client, dir := proxy.client.SendClientRequestHTTP(req)

	if client != nil {
		dii2perrs.Log("si-socks-proxy.go Client was retrieved: ", dir)
		resp, doerr := proxy.Do(req, client, 0)
		if proxy.c, proxy.err = dii2perrs.Warn(doerr, "si-socks-proxy.go Encountered an oddly formed response. Skipping.", "si-socks-proxy.go Processing Response"); proxy.c {
			resp := proxy.client.CopyRequest(req, resp, dir)
			proxy.printResponse(rW, resp)
			dii2perrs.Log("si-socks-proxy.go responded")
			return
		}
		if !strings.Contains(doerr.Error(), "malformed HTTP status code") && !strings.Contains(doerr.Error(), "use of closed network connection") {
			if resp != nil {
				resp := proxy.client.CopyRequest(req, resp, dir)
				proxy.printResponse(rW, resp)
			}
			dii2perrs.Log("si-socks-proxy.go status error", doerr.Error())
			return
		}
		dii2perrs.Log("si-socks-proxy.go status error", doerr.Error())
		return
	}
	log.Println("si-socks-proxy.go client retrieval error")
	return
}

// Do does a request
func (proxy *SamSOCKSProxy) Do(req *http.Request, client *http.Client, x int) (*http.Response, error) {
	req.RequestURI = ""

	resp, doerr := client.Do(req)

	if req.Close {
		dii2perrs.Log("request must be closed")
	}

	if proxy.c, proxy.err = dii2perrs.Warn(doerr, "si-socks-proxy.go Response body error:", "si-socks-proxy.go Read response body"); proxy.c {
		return resp, doerr
	}
	return resp, doerr
}

func (proxy *SamSOCKSProxy) printResponse(rW http.ResponseWriter, r *http.Response) {
	if r != nil {
		defer r.Body.Close()
		copyHeader(rW.Header(), r.Header)
		rW.WriteHeader(r.StatusCode)
		io.Copy(rW, r.Body)
		dii2perrs.Log("si-socks-proxy.go Response status:", r.Status)
	}
}

// CreateSOCKSProxy generates a SOCKS proxy conf := &socks5.Config{}
func CreateSOCKSProxy(proxAddr, proxPort, initAddress, ahAddr, ahPort, addressHelperURL string, samStack *dii2pmain.SamList, timeoutTime int, keepAlives bool) *SamSOCKSProxy {
	var samProxy SamSOCKSProxy
	samProxy.Addr = proxAddr + ":" + proxPort
	samProxy.keepAlives = keepAlives
	samProxy.addressbook, samProxy.err = jumpresolver.NewJumpResolver(ahAddr, ahPort)
	dii2perrs.Fatal(samProxy.err, "si-http-proxy.go Addresshelper not available", "si-i2p-proxy.go Connected to addresshelper")
	log.Println("si-socks-proxy.go Starting SOCKS proxy on:" + samProxy.Addr)
	samProxy.client = samStack
	samProxy.timeoutTime = time.Duration(timeoutTime) * time.Minute
	samProxy.newHandle, samProxy.err = socks5.New(&socks5.Config{})
	dii2perrs.Fatal(samProxy.err, "si-socks-proxy.go SOCKS proxy creation error", "si-socks-proxy.go SOCKS proxy created")
	log.Println("si-socks-proxy.go Connected SAM isolation stack to the SOCKS proxy server")
	go samProxy.prepare()
	log.Println("si-socks-proxy.go SOCKS Proxy prepared")
	return &samProxy
}
