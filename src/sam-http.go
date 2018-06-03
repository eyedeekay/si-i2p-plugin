package dii2p

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
	"strings"
	"time"

	"github.com/eyedeekay/gosam"
	"github.com/eyedeekay/i2pasta/addresshelper"
)

//DEBUG Remove this when you get the options laid in properly.
var DEBUG bool

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

	useTime  time.Time
	lifeTime time.Duration

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

var connectionDirectory string

func (samConn *SamHTTP) initPipes() {
	checkFolder(filepath.Join(connectionDirectory, samConn.host))

	samConn.sendPath, samConn.sendPipe, samConn.err = setupFiFo(filepath.Join(connectionDirectory, samConn.host), "send")
	if samConn.c, samConn.err = Fatal(samConn.err, "sam-http.go Pipe setup error", "sam-http.go Pipe setup"); samConn.c {
		samConn.sendScan, samConn.err = setupScanner(filepath.Join(connectionDirectory, samConn.host), "send", samConn.sendPipe)
		if samConn.c, samConn.err = Fatal(samConn.err, "sam-http.go Scanner setup Error:", "sam-http.go Scanner set up successfully."); !samConn.c {
			samConn.CleanupClient()
		}
	}

	samConn.namePath, samConn.nameFile, samConn.err = setupFile(filepath.Join(connectionDirectory, samConn.host), "name")
	if samConn.c, samConn.err = Fatal(samConn.err, "sam-http.go Pipe setup error", "sam-http.go Pipe setup"); samConn.c {
		samConn.nameFile.WriteString("")
	}

	samConn.idPath, samConn.idFile, samConn.err = setupFile(filepath.Join(connectionDirectory, samConn.host), "id")
	if samConn.c, samConn.err = Fatal(samConn.err, "sam-http.go Pipe setup error", "sam-http.go Pipe setup"); samConn.c {
		samConn.idFile.WriteString("")
	}

	samConn.base64Path, samConn.base64File, samConn.err = setupFile(filepath.Join(connectionDirectory, samConn.host), "base64")
	if samConn.c, samConn.err = Fatal(samConn.err, "sam-http.go Pipe setup error", "sam-http.go Pipe setup"); samConn.c {
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
	samConn.samBridgeClient, samConn.err = goSam.NewClientFromOptions(goSam.SetHost(samConn.samAddrString), goSam.SetPort(samConn.samPortString), goSam.SetDebug(DEBUG), goSam.SetUnpublished(true), goSam.SetInQuantity(15), goSam.SetOutQuantity(15))
	if samConn.c, samConn.err = Warn(samConn.err, "sam-http.go SAM connection error", "sam-http.go Initializing SAM connection"); samConn.c {
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
		if samConn.c, samConn.err = Warn(samConn.err, "sam-http.go Error connecting SAM streams", "sam-http.go Connecting SAM streams"); samConn.c {
			Log("sam-http.go Stream Connection established")
			return samConn.samBridgeClient.SamConn, samConn.err
		}
		return samConn.reConnect()
	}
	return samConn.reConnect()
}

func (samConn *SamHTTP) reConnect() (net.Conn, error) {
	samConn.samBridgeClient, samConn.err = goSam.NewClientFromOptions(goSam.SetHost(samConn.samAddrString), goSam.SetPort(samConn.samPortString), goSam.SetDebug(DEBUG), goSam.SetUnpublished(true), goSam.SetInQuantity(15), goSam.SetOutQuantity(15))
	if samConn.c, samConn.err = Warn(samConn.err, "sam-http.go 133 SAM Client connection error", "sam-http.go SAM client connecting"); samConn.c {
		Log("sam-http.go SAM Connection established")
		samConn.err = samConn.samBridgeClient.StreamConnect(samConn.id, samConn.name)
		if samConn.c, samConn.err = Warn(samConn.err, "sam-http.go Connecting SAM streams", "sam-http.go Connecting SAM streams"); samConn.c {
			Log("sam-http.go Stream Connection established")
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
	Log("sam-http.go Setting Transport")
	Log("sam-http.go Setting Dial function")
	samConn.transport = &http.Transport{
		Dial: samConn.Dial,
		//Dial:                  samConn.samBridgeClient.Dial,
		MaxIdleConns:          0,
		MaxIdleConnsPerHost:   4,
		DisableKeepAlives:     samConn.keepAlives,
		ResponseHeaderTimeout: samConn.otherTimeoutTime,
		ExpectContinueTimeout: samConn.otherTimeoutTime,
		IdleConnTimeout:       samConn.timeoutTime,
		TLSNextProto:          make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	}
	Log("sam-http.go Initializing sub-client")
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
	if samConn.c, samConn.err = Fatal(samConn.err, "sam-http.go Cookie Jar creation error", "sam-http.go Cookie Jar creating", samConn.samAddrString, samConn.samPortString); samConn.c {
		Log("sam-http.go Cookie Jar created")
	}
	samConn.samBridgeClient, samConn.err = goSam.NewClientFromOptions(goSam.SetHost(samConn.samAddrString), goSam.SetPort(samConn.samPortString), goSam.SetDebug(DEBUG), goSam.SetUnpublished(true), goSam.SetInQuantity(15), goSam.SetOutQuantity(15))
	if samConn.c, samConn.err = Fatal(samConn.err, "sam-http.go SAM Client Connection Error", "sam-http.go SAM client connecting", samConn.samAddrString, samConn.samPortString); samConn.c {
		Log("sam-http.go Setting Transport")
		Log("sam-http.go Setting Dial function")
		samConn.setupTransport()
		if samConn.host == "" {
			samConn.host, samConn.directory = samConn.hostSet(samConn.initRequestURL)
			samConn.initPipes()
		}
		samConn.setName(samConn.initRequestURL)
		samConn.subCache = append(samConn.subCache, NewSamURL(samConn.directory))
	}
}

func (samConn *SamHTTP) cleanURL(request string) (string, string) {
	Log("sam-http.go cleanURL Trim 0 " + request)
	//http://i2p-projekt.i2p/en/downloads
	url := strings.Replace(request, "http://", "", -1)
	Log("sam-http.go cleanURL Request URL " + url)
	//i2p-projekt.i2p/en/downloads
	if strings.HasSuffix(url, ".i2p") {
		url = url + "/"
	}
	host := strings.SplitAfter(url, ".i2p/")[0]
	if strings.HasSuffix(host, ".i2p/") {
		host = host[:len(host)-len("/")]
	}
	if strings.HasSuffix(url, ".i2p/") {
		url = url[:len(url)-len("/")]
	}
	Log("sam-http.go cleanURL Trim 2 " + host)
	return host, url
}

func (samConn *SamHTTP) hostSet(request string) (string, string) {
	host, req := samConn.cleanURL(request)
	Log("sam-http.go Setting up micro-proxy for:", "http://"+host)
	Log("sam-http.go in Directory", req)
	return host, req
}

func (samConn *SamHTTP) hostGet() string {
	return "http://" + samConn.host
}

func (samConn *SamHTTP) hostCheck(request string) int {
	host, _ := samConn.cleanURL(request)
	_, err := url.ParseRequestURI(host)
    if samConn.lifeTime < time.Now().Sub(samConn.useTime) {
        Warn(nil, "sam-http.go Removing inactive client", "sam-http.go Removing inactive client", samConn.host)
		samConn.CleanupClient()
		return -1
	}
	Log("sam-http.go keeping client alive")
	samConn.useTime = time.Now()
	if err == nil {
		if samConn.host == host {
			Log("sam-http.go Request host ", host, "is equal to client host", samConn.host)
			return 1
		}
		Log("sam-http.go Request host ", host, "is not equal to client host", samConn.host)
		return 0
	}
	if samConn.host == host {
		Log("sam-http.go Request host ", host, "is equal to client host", samConn.host)
		return 1
	}
	Log("sam-http.go Request host ", host, "is not equal to client host", samConn.host)
	return 0
}

func (samConn *SamHTTP) getURL(request string) (string, string) {
	r := request
	directory := strings.Replace(safeURLString(request), "http://", "", -1)
	_, err := url.ParseRequestURI(r)
	if err != nil {
		r = "http://" + request
		Log("sam-http.go URL failed validation, correcting to:", r)
	} else {
		Log("sam-http.go URL passed validation:", request)
	}
	Log("sam-http.go Request will be managed in:", directory)
	return r, directory
}

func (samConn *SamHTTP) sendRequest(request string) (*http.Response, error) {
	r, dir := samConn.getURL(request)
	Log("sam-http.go Getting resource", request)
	if samConn.subClient != nil {
		resp, err := samConn.subClient.Get(r)
		Warn(err, "sam-http.go Response Error", "sam-http.go Getting Response")
		Log("sam-http.go Pumping result to top of parent pipe")
		samConn.copyRequest(resp, dir)
		return resp, err
	}
	return nil, nil
}

func (samConn *SamHTTP) getURLHTTP(request *http.Request) (string, string) {
	directory := strings.Replace(safeURLString(request.URL.String()), "http://", "", -1)
	return request.URL.String(), directory
}

func (samConn *SamHTTP) sendRequestHTTP(request *http.Request) (*http.Client, string) {
	r, dir := samConn.getURLHTTP(request)
	Log("sam-http.go Getting resource", r, "In ", dir)
	return samConn.subClient, dir
}

func (samConn *SamHTTP) sendRequestBase64HTTP(request *http.Request, base64helper string) (*http.Client, string) {
	r, dir := samConn.getURL(request.URL.String())
	Log("sam-http.go Getting resource", r, "In ", dir)
	return samConn.subClient, dir
}

func (samConn *SamHTTP) findSubCache(response *http.Response, directory string) *SamURL {
	b := false
	var u SamURL
	for _, url := range samConn.subCache {
		Log("sam-http.go Seeking Subdirectory", url.subDirectory)
		if url.checkDirectory(directory) {
			return &url
		}
	}
	if b == false {
		Log("sam-http.go has not been retrieved yet. Setting up:", directory)
		samConn.subCache = append(samConn.subCache, NewSamURL(directory))
		for _, url := range samConn.subCache {
			Log("sam-http.go Seeking Subdirectory", url.subDirectory)
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

func (samConn *SamHTTP) scannerText() (string, error) {
	text := ""
	var err error
	for _, url := range samConn.subCache {
		text, err = url.scannerText()
		if len(text) > 0 {
			break
		}
	}
	return text, err
}

func (samConn *SamHTTP) printResponse() string {
	s, e := samConn.scannerText()
	if samConn.c, samConn.err = Fatal(e, "sam-http.go Response Retrieval Error", "sam-http.go Retrieving Responses"); !samConn.c {
		Log("sam-http.go Response Panic")
		samConn.CleanupClient()
	} else {
		Log("sam-http.go Response Retrieved")
	}
	return s
}

func (samConn *SamHTTP) readRequest() string {
	text := samConn.sendScan.Text()
	for samConn.sendScan.Scan() {
		samConn.sendRequest(text)
	}
	clearFile(filepath.Join(connectionDirectory, samConn.directory), "send")
	return text
}

func (samConn *SamHTTP) readDelete() bool {
	b := false
	for _, dir := range samConn.subCache {
		n := dir.readDelete()
		if !n {
			Log("sam-http.go Maintaining Connection:", samConn.hostGet())
		} else {
			b = n
		}
	}
	return b
}

func (samConn *SamHTTP) writeName() {
	Log("sam-http.go Looking up hostname:", samConn.host)
	samConn.name, samConn.err = samConn.samBridgeClient.Lookup(samConn.host)
	samConn.nameFile.WriteString(samConn.name)
}

func (samConn *SamHTTP) writeSession(request string) {
	Log("sam-http.go Caching base64 address of:", samConn.host+" "+samConn.name)
	samConn.id, samConn.base64, samConn.err = samConn.samBridgeClient.CreateStreamSession("")
	samConn.idFile.WriteString(fmt.Sprint(samConn.id))
	Warn(samConn.err, "sam-http.go Local Base64 Caching error", "sam-http.go Cachine Base64 Address of:", request)
	log.Println("sam-http.go Tunnel id: ", samConn.id)
	Log("sam-http.go Tunnel dest: ", samConn.base64)
	samConn.base64File.WriteString(samConn.base64)
	Log("sam-http.go New Connection Name: ", samConn.base64)
}

func (samConn *SamHTTP) setName(request string) {
	if samConn.checkName() {
		samConn.host, samConn.directory = samConn.hostSet(request)
		Log("sam-http.go Setting hostname:", samConn.host)
		samConn.writeName()
		samConn.writeSession(request)
	} else {
		samConn.host, samConn.directory = samConn.hostSet(request)
		Log("sam-http.go Setting hostname:", samConn.host)
		samConn.initPipes()
		samConn.writeName()
		samConn.writeSession(request)
	}
}

func (samConn *SamHTTP) checkName() bool {
	Log("sam-http.go seeing if the connection needs a name:")
	if samConn.name != "" {
		Log("sam-http.go Naming connection: Connection name was empty.")
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
		url.cleanupDirectory()
	}
	err := samConn.samBridgeClient.Close()
	if samConn.c, samConn.err = Warn(err, "sam-http.go Closing SAM bridge error, retrying.", "sam-http.go Closing SAM bridge"); !samConn.c {
		samConn.samBridgeClient.Close()
	}
	os.RemoveAll(filepath.Join(connectionDirectory, samConn.host))
}

func newSamHTTP(samAddrString, samPortString, request string, timeoutTime, lifeTime int, keepAlives bool) SamHTTP {
	Log("sam-http.go Creating a new SAMv3 Client.")
	samConn, err := NewSamHTTPFromOptions(
		SetSamHTTPHost(samAddrString),
		SetSamHTTPPort(samPortString),
		SetSamHTTPRequest(request),
		SetSamHTTPTimeout(timeoutTime),
		SetSamHTTPKeepAlives(keepAlives),
		SetSamHTTPLifespan(lifeTime),
	)
	Fatal(err, "sam-http.go Pipe setup error", "sam-http.go Pipe setup")
	return samConn
}

func newSamHTTPHTTP(samAddrString, samPortString string, request *http.Request, timeoutTime, lifeTime int, keepAlives bool) SamHTTP {
	Log("sam-http.go Creating a new SAMv3 Client.")
	samConn, err := NewSamHTTPFromOptions(
		SetSamHTTPHost(samAddrString),
		SetSamHTTPPort(samPortString),
		SetSamHTTPRequest(request.URL.String()),
		SetSamHTTPTimeout(timeoutTime),
		SetSamHTTPKeepAlives(keepAlives),
		SetSamHTTPLifespan(lifeTime),
	)
	Fatal(err, "sam-http.go Pipe setup error", "sam-http.go Pipe setup")
	return samConn
}

//NewSamHTTPFromOptions creates a new SamHTTP connection manager for a single eepSite
func NewSamHTTPFromOptions(opts ...func(*SamHTTP) error) (SamHTTP, error) {
	Log("sam-http.go Creating a new SAMv3 Client.")
	var samConn SamHTTP
	for _, o := range opts {
		if err := o(&samConn); err != nil {
			return samConn, err
		}
	}
	samConn.createClient()
	return samConn, nil
}
