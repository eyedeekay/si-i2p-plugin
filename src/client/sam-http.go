package dii2pmain

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/eyedeekay/gosam"
	"github.com/eyedeekay/i2pasta/addresshelper"
)

import (
	"github.com/eyedeekay/si-i2p-plugin/src/errors"
	"github.com/eyedeekay/si-i2p-plugin/src/helpers"
)

//SamHTTP is an HTTP proxy which requests resources from the i2p network using
//the same unique destination
type SamHTTP struct {
	subCache []SamURL
	err      error
	c        bool

	samBridgeClient *goSam.Client
	assistant       i2paddresshelper.I2paddresshelper
	jar             *cookiejar.Jar
	samAddrString   string
	samPortString   string
	initRequestURL  string

	useTime                time.Time
	lifeTime               time.Duration
	tunnelLength           int
	inboundQuantity        int
	outboundQuantity       int
	inboundBackupQuantity  int
	outboundBackupQuantity int
	idleConns              int

	transport *http.Transport
	subClient *http.Client

	timeoutTime      time.Duration
	otherTimeoutTime time.Duration
	keepAlives       bool

	host      string
	directory string

	sendPath string
	sendPipe *os.File
	sendScan *bufio.Scanner

	namePath string
	nameFile *os.File
	name     string

	idPath string
	idFile *os.File
	id     int32

	base64Path string
	base64File *os.File
	base64     string
}

func (samConn *SamHTTP) initPipes() {
	dii2phelper.SetupFolder(filepath.Join(dii2phelper.ConnectionDirectory, samConn.host))

	samConn.sendPath, samConn.sendPipe, samConn.err = dii2phelper.SetupFiFo(filepath.Join(dii2phelper.ConnectionDirectory, samConn.host), "send")
	if samConn.c, samConn.err = dii2perrs.Fatal(samConn.err, "sam-http.go Pipe setup error", "sam-http.go Pipe setup"); samConn.c {
		samConn.sendScan, samConn.err = dii2phelper.SetupScanner(filepath.Join(dii2phelper.ConnectionDirectory, samConn.host), "send", samConn.sendPipe)
		if samConn.c, samConn.err = dii2perrs.Fatal(samConn.err, "sam-http.go Scanner setup Error:", "sam-http.go Scanner set up successfully."); !samConn.c {
			samConn.CleanupClient()
		}
	}

	samConn.namePath, samConn.nameFile, samConn.err = dii2phelper.SetupFile(filepath.Join(dii2phelper.ConnectionDirectory, samConn.host), "name")
	if samConn.c, samConn.err = dii2perrs.Fatal(samConn.err, "sam-http.go Pipe setup error", "sam-http.go Pipe setup"); samConn.c {
		samConn.nameFile.WriteString("")
	}

	samConn.idPath, samConn.idFile, samConn.err = dii2phelper.SetupFile(filepath.Join(dii2phelper.ConnectionDirectory, samConn.host), "id")
	if samConn.c, samConn.err = dii2perrs.Fatal(samConn.err, "sam-http.go Pipe setup error", "sam-http.go Pipe setup"); samConn.c {
		samConn.idFile.WriteString("")
	}

	samConn.base64Path, samConn.base64File, samConn.err = dii2phelper.SetupFile(filepath.Join(dii2phelper.ConnectionDirectory, samConn.host), "base64")
	if samConn.c, samConn.err = dii2perrs.Fatal(samConn.err, "sam-http.go Pipe setup error", "sam-http.go Pipe setup"); samConn.c {
		samConn.idFile.WriteString("")
	}

}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

//Dial is a custom Dialer function that allows us to keep the same i2p destination
//on a per-eepSite basis
func (samConn *SamHTTP) Dial(network, addr string) (net.Conn, error) {
	samConn.samBridgeClient, samConn.err = goSam.NewClientFromOptions(
		goSam.SetHost(samConn.samAddrString),
		goSam.SetPort(samConn.samPortString),
		goSam.SetDebug(dii2perrs.DEBUG),
		goSam.SetUnpublished(true),
		goSam.SetInLength(uint(samConn.tunnelLength)),
		goSam.SetOutLength(uint(samConn.tunnelLength)),
		goSam.SetInQuantity(uint(samConn.inboundQuantity)),
		goSam.SetOutQuantity(uint(samConn.outboundQuantity)),
	)
	if samConn.c, samConn.err = dii2perrs.Warn(samConn.err, "sam-http.go SAM connection error", "sam-http.go Initializing SAM connection"); samConn.c {
		return samConn.subDial(network, addr)
	}
	return samConn.samBridgeClient.SamConn, &errorString{"SAM connection error"}
}

func (samConn *SamHTTP) subDial(network, addr string) (net.Conn, error) {
	if samConn.name != "" {
		if samConn.id != 0 {
			return samConn.connect()
		}
		return nil, &errorString{"ID error"}
	}
	return nil, &errorString{"Hostname error"}
}

func (samConn *SamHTTP) connect() (net.Conn, error) {
	if samConn.samBridgeClient != nil {
		samConn.err = samConn.samBridgeClient.StreamConnect(samConn.id, samConn.name)
		if samConn.c, samConn.err = dii2perrs.Warn(samConn.err, "sam-http.go Error connecting SAM streams", "sam-http.go Connecting SAM streams"); samConn.c {
			dii2perrs.Log("sam-http.go Stream Connection established")
			return samConn.samBridgeClient.SamConn, samConn.err
		}
		return samConn.reConnect()
	}
	return samConn.reConnect()
}

func (samConn *SamHTTP) reConnect() (net.Conn, error) {
	samConn.samBridgeClient, samConn.err = goSam.NewClientFromOptions(
		goSam.SetHost(samConn.samAddrString),
		goSam.SetPort(samConn.samPortString),
		goSam.SetDebug(dii2perrs.DEBUG),
		goSam.SetUnpublished(true),
		goSam.SetInLength(uint(samConn.tunnelLength)),
		goSam.SetOutLength(uint(samConn.tunnelLength)),
		goSam.SetInQuantity(uint(samConn.inboundQuantity)),
		goSam.SetOutQuantity(uint(samConn.outboundQuantity)),
	)
	if samConn.c, samConn.err = dii2perrs.Warn(samConn.err, "sam-http.go 133 SAM Client connection error", "sam-http.go SAM client connecting"); samConn.c {
		dii2perrs.Log("sam-http.go SAM Connection established")
		samConn.err = samConn.samBridgeClient.StreamConnect(samConn.id, samConn.name)
		if samConn.c, samConn.err = dii2perrs.Warn(samConn.err, "sam-http.go Connecting SAM streams", "sam-http.go Connecting SAM streams"); samConn.c {
			dii2perrs.Log("sam-http.go Stream Connection established")
			return samConn.samBridgeClient.SamConn, samConn.err
		}
		return samConn.reConnect()
	}
	//samConn.samBridgeClient.Close()
	return samConn.reConnect()
}

func (samConn *SamHTTP) checkRedirect(req *http.Request, via []*http.Request) error {
	return nil
}

func (samConn *SamHTTP) setupTransport() {
	dii2perrs.Log("sam-http.go Setting Transport")
	dii2perrs.Log("sam-http.go Setting Dial function")
	samConn.transport = &http.Transport{
		Dial: samConn.Dial,
		//Dial:                  samConn.samBridgeClient.Dial,
		MaxIdleConns:          0,
		MaxIdleConnsPerHost:   samConn.idleConns,
		DisableKeepAlives:     samConn.keepAlives,
		ResponseHeaderTimeout: samConn.otherTimeoutTime,
		ExpectContinueTimeout: samConn.otherTimeoutTime,
		IdleConnTimeout:       samConn.timeoutTime,
		TLSNextProto:          make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}
	dii2perrs.Log("sam-http.go Initializing sub-client")
	samConn.subClient = &http.Client{
		Timeout:       samConn.timeoutTime,
		Transport:     samConn.transport,
		Jar:           samConn.jar,
		CheckRedirect: nil,
	}
	//
}

func (samConn *SamHTTP) createClient() {
	samConn.jar, samConn.err = cookiejar.New(nil)
	if samConn.c, samConn.err = dii2perrs.Fatal(samConn.err, "sam-http.go Cookie Jar creation error", "sam-http.go Cookie Jar creating", samConn.samAddrString, samConn.samPortString); samConn.c {
		dii2perrs.Log("sam-http.go Cookie Jar created")
	}
	samConn.samBridgeClient, samConn.err = goSam.NewClientFromOptions(
		goSam.SetHost(samConn.samAddrString),
		goSam.SetPort(samConn.samPortString),
		goSam.SetDebug(dii2perrs.DEBUG),
		goSam.SetUnpublished(true),
		goSam.SetInLength(uint(samConn.tunnelLength)),
		goSam.SetOutLength(uint(samConn.tunnelLength)),
		goSam.SetInQuantity(uint(samConn.inboundQuantity)),
		goSam.SetOutQuantity(uint(samConn.outboundQuantity)),
		goSam.SetInBackups(uint(samConn.inboundBackupQuantity)),
		goSam.SetOutBackups(uint(samConn.outboundBackupQuantity)),
	)
	if samConn.c, samConn.err = dii2perrs.Fatal(samConn.err, "sam-http.go SAM Client Connection Error", "sam-http.go SAM client connecting", samConn.samAddrString, samConn.samPortString); samConn.c {
		dii2perrs.Log("sam-http.go Setting Transport")
		dii2perrs.Log("sam-http.go Setting Dial function")
		samConn.setupTransport()
		if samConn.host == "" {
			samConn.host, samConn.directory = samConn.hostSet(samConn.initRequestURL)
			samConn.initPipes()
		}
		samConn.setName(samConn.initRequestURL)
		samConn.subCache = append(samConn.subCache, NewSamURL(samConn.directory))
	}
}

func (samConn *SamHTTP) hostSet(request string) (string, string) {
	host, req := dii2phelper.CleanURL(request)
	dii2perrs.Log("sam-http.go Setting up micro-proxy for:", "http://"+host)
	dii2perrs.Log("sam-http.go in Directory", req)
	return host, req
}

func (samConn *SamHTTP) hostGet() string {
	return "http://" + samConn.host
}

func (samConn *SamHTTP) hostCheck(request string) bool {
	host, u := dii2phelper.CleanURL(request)
	_, err := url.ParseRequestURI(u)
	dii2perrs.Log("sam-http.go keeping client alive")
	samConn.useTime = time.Now()
	if err == nil {
		if samConn.host == host {
			dii2perrs.Log("sam-http.go Request host ", host, "is equal to client host", samConn.host)
			return true
		}
		dii2perrs.Log("sam-http.go Request host ", host, "is not equal to client host", samConn.host)
		return false
	}
	return false
}

func (samConn *SamHTTP) lifetimeCheck(request string) bool {
	if samConn.lifeTime < time.Now().Sub(samConn.useTime) {
		dii2perrs.Warn(nil, "sam-http.go Error Removing inactive client after ", "sam-http.go Removing inactive client", samConn.host, "after", samConn.lifeTime, "minutes.")
		samConn.useTime = time.Now().Add(time.Duration(120) * time.Second)
		return true
	}
	samConn.useTime = time.Now()
	return false
}

func (samConn *SamHTTP) getURL(request string) (string, string) {
	r := request
	//directory := strings.Replace(dii2phelper.SafeURLString(request), "http://", "", -1)
	directory := dii2phelper.SafeURLString(request)
	_, err := url.ParseRequestURI(r)
	if err != nil {
		r = "http://" + request
		dii2perrs.Log("sam-http.go URL failed validation, correcting to:", r)
	} else {
		dii2perrs.Log("sam-http.go URL passed validation:", request)
	}
	dii2perrs.Log("sam-http.go Request will be managed in:", directory)
	return r, directory
}

func (samConn *SamHTTP) sendRequest(request string) (*http.Response, error) {
	r, dir := samConn.getURL(request)
	dii2perrs.Log("sam-http.go Getting resource", request)
	if samConn.subClient != nil {
		resp, err := samConn.subClient.Get(r)
		dii2perrs.Warn(err, "sam-http.go Response Error", "sam-http.go Getting Response")
		dii2perrs.Log("sam-http.go Pumping result to top of parent pipe")
		samConn.copyRequest(resp, dir)
		return resp, err
	}
	return nil, nil
}

func (samConn *SamHTTP) getURLHTTP(request *http.Request) (string, string) {
	directory := dii2phelper.SafeURLString(request.URL.String())
	return request.URL.String(), directory
}

func (samConn *SamHTTP) sendRequestHTTP(request *http.Request) (*http.Client, string) {
	r, dir := samConn.getURLHTTP(request)
	dii2perrs.Log("sam-http.go Getting resource", r, "In ", dir)
	return samConn.subClient, dir
}

func (samConn *SamHTTP) sendRequestBase64HTTP(request *http.Request, base64helper string) (*http.Client, string) {
	r, dir := samConn.getURL(request.URL.String())
	dii2perrs.Log("sam-http.go Getting resource", r, "In ", dir)
	return samConn.subClient, dir
}

func (samConn *SamHTTP) findSubCache(response *http.Response, directory string) *SamURL {
	b := false
	var u SamURL
	for _, url := range samConn.subCache {
		dii2perrs.Log("sam-http.go Seeking Subdirectory", url.subDirectory)
		if url.checkDirectory(directory) {
			return &url
		}
	}
	if b == false {
		dii2perrs.Log("sam-http.go has not been retrieved yet. Setting up:", directory)
		samConn.subCache = append(samConn.subCache, NewSamURL(directory))
		for _, url := range samConn.subCache {
			dii2perrs.Log("sam-http.go Seeking Subdirectory", url.subDirectory)
			if url.checkDirectory(directory) {
				u = url
				return &u
			}
		}
	}
	return &u
}

func (samConn *SamHTTP) copyRequest(response *http.Response, directory string) {
	samConn.findSubCache(response, directory).copyDirectory(response, directory)
}

func (samConn *SamHTTP) copyRequestHTTP(request *http.Request, response *http.Response, directory string) *http.Response {
	return samConn.findSubCache(response, directory).copyDirectoryHTTP(request, response, directory)
}

func (samConn *SamHTTP) ScannerText() (string, error) {
	text := ""
	var err error
	for _, url := range samConn.subCache {
		text, err = url.ScannerText()
		if len(text) > 0 {
			break
		}
	}
	return text, err
}

func (samConn *SamHTTP) printResponse() string {
	s, e := samConn.ScannerText()
	if samConn.c, samConn.err = dii2perrs.Fatal(e, "sam-http.go Response Retrieval Error", "sam-http.go Retrieving Responses"); !samConn.c {
		dii2perrs.Log("sam-http.go Response Panic")
		samConn.CleanupClient()
	} else {
		dii2perrs.Log("sam-http.go Response Retrieved")
	}
	return s
}

func (samConn *SamHTTP) readRequest() string {
	text := samConn.sendScan.Text()
	for samConn.sendScan.Scan() {
		samConn.sendRequest(text)
	}
	dii2phelper.ClearFile(filepath.Join(dii2phelper.ConnectionDirectory, samConn.directory), "send")
	return text
}

func (samConn *SamHTTP) readDelete() bool {
	b := false
	for _, dir := range samConn.subCache {
		n := dir.readDelete()
		if !n {
			dii2perrs.Log("sam-http.go Maintaining Connection:", samConn.hostGet())
		} else {
			b = n
		}
	}
	return b
}

func (samConn *SamHTTP) writeName() {
	dii2perrs.Log("sam-http.go Looking up hostname:", samConn.host)
	samConn.name, samConn.err = samConn.samBridgeClient.Lookup(samConn.host)
	samConn.nameFile.WriteString(samConn.name)
}

func (samConn *SamHTTP) writeSession(request string) {
	dii2perrs.Log("sam-http.go Caching base64 address of:", samConn.host+" "+samConn.name)
	samConn.id, samConn.base64, samConn.err = samConn.samBridgeClient.CreateStreamSession("")
	samConn.idFile.WriteString(fmt.Sprint(samConn.id))
	dii2perrs.Warn(samConn.err, "sam-http.go Local Base64 Caching error", "sam-http.go Cachine Base64 Address of:", request)
	log.Println("sam-http.go Tunnel id: ", samConn.id)
	dii2perrs.Log("sam-http.go Tunnel dest: ", samConn.base64)
	samConn.base64File.WriteString(samConn.base64)
	dii2perrs.Log("sam-http.go New Connection Name: ", samConn.base64)
}

func (samConn *SamHTTP) setName(request string) {
	if samConn.checkName() {
		samConn.host, samConn.directory = samConn.hostSet(request)
		dii2perrs.Log("sam-http.go Setting hostname:", samConn.host)
		samConn.writeName()
		samConn.writeSession(request)
	} else {
		samConn.host, samConn.directory = samConn.hostSet(request)
		dii2perrs.Log("sam-http.go Setting hostname:", samConn.host)
		samConn.initPipes()
		samConn.writeName()
		samConn.writeSession(request)
	}
}

func (samConn *SamHTTP) checkName() bool {
	dii2perrs.Log("sam-http.go seeing if the connection needs a name:")
	if samConn.name != "" {
		dii2perrs.Log("sam-http.go Naming connection: Connection name was empty.")
		return true
	}
	return false
}

//CleanupClient completely tears down a SamHTTP client
func (samConn *SamHTTP) CleanupClient() {
	samConn.sendPipe.Close()
	samConn.nameFile.Close()
	samConn.idFile.Close()
	samConn.base64File.Close()
	for _, url := range samConn.subCache {
		url.CleanupDirectory()
	}
	err := samConn.samBridgeClient.Close()
	if samConn.c, samConn.err = dii2perrs.Warn(err, "sam-http.go Closing SAM bridge error, retrying.", "sam-http.go Closing SAM bridge"); !samConn.c {
		samConn.samBridgeClient.Close()
	}
	os.RemoveAll(filepath.Join(dii2phelper.ConnectionDirectory, samConn.host))
}

func newSamHTTP(samAddrString, samPortString, request string, timeoutTime, lifeTime int, keepAlives bool, tunnelLength, inboundQuantity, outboundQuantity, idleConns, inboundBackups, outboundBackups int) SamHTTP {
	dii2perrs.Log("sam-http.go Creating a new SAMv3 Client.")
	samConn, err := NewSamHTTPFromOptions(
		SetSamHTTPHost(samAddrString),
		SetSamHTTPPort(samPortString),
		SetSamHTTPRequest(request),
		SetSamHTTPTimeout(timeoutTime),
		SetSamHTTPKeepAlives(keepAlives),
		SetSamHTTPLifespan(lifeTime),
		SetSamHTTPTunLength(tunnelLength),
		SetSamHTTPInQuantity(inboundQuantity),
		SetSamHTTPOutQuantity(outboundQuantity),
		SetSamHTTPIdleQuantity(idleConns),
		SetSamHTTPInBackupQuantity(inboundBackups),
		SetSamHTTPOutBackupQuantity(outboundBackups),
	)
	dii2perrs.Fatal(err, "sam-http.go Pipe setup error", "sam-http.go Pipe setup")
	return samConn
}

func newSamHTTPHTTP(
	samAddrString, samPortString string,
	request *http.Request,
	timeoutTime,
	lifeTime int,
	keepAlives bool,
	tunnelLength, inboundQuantity, outboundQuantity, idleConns, inboundBackups, outboundBackups int) SamHTTP {
	dii2perrs.Log("sam-http.go Creating a new SAMv3 Client.")
	samConn, err := NewSamHTTPFromOptions(
		SetSamHTTPHost(samAddrString),
		SetSamHTTPPort(samPortString),
		SetSamHTTPRequest(request.URL.String()),
		SetSamHTTPTimeout(timeoutTime),
		SetSamHTTPKeepAlives(keepAlives),
		SetSamHTTPLifespan(lifeTime),
		SetSamHTTPTunLength(tunnelLength),
		SetSamHTTPInQuantity(inboundQuantity),
		SetSamHTTPOutQuantity(outboundQuantity),
		SetSamHTTPIdleQuantity(idleConns),
		SetSamHTTPInBackupQuantity(inboundBackups),
		SetSamHTTPOutBackupQuantity(outboundBackups),
	)
	dii2perrs.Fatal(err, "sam-http.go Pipe setup error", "sam-http.go Pipe setup")
	return samConn
}

//NewSamHTTPFromOptions creates a new SamHTTP connection manager for a single eepSite
func NewSamHTTPFromOptions(opts ...func(*SamHTTP) error) (SamHTTP, error) {
	dii2perrs.Log("sam-http.go Creating a new SAMv3 Client.")
	var samConn SamHTTP
	samConn.samAddrString = "127.0.0.1"
	samConn.samPortString = "7656"
	samConn.initRequestURL = ""
	samConn.timeoutTime = time.Duration(6) * time.Minute
	samConn.otherTimeoutTime = time.Duration(2) * time.Minute
	samConn.keepAlives = true
	samConn.lifeTime = time.Duration(12) * time.Minute
	samConn.useTime = time.Now()
	samConn.tunnelLength = 3
	samConn.inboundQuantity = 15
	samConn.outboundQuantity = 15
	samConn.inboundBackupQuantity = 3
	samConn.outboundBackupQuantity = 3
	samConn.idleConns = 4
	for _, o := range opts {
		if err := o(&samConn); err != nil {
			return samConn, err
		}
	}
	samConn.createClient()
	return samConn, nil
}
