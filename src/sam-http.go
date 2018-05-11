package dii2p

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"io"
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

// Remove this when you get the options laid in properly.
var DEBUG bool

type SamHttp struct {
	subCache []samUrl
	err      error
	c        bool

	samBridgeClient *goSam.Client
	assistant       i2paddresshelper.I2paddresshelper
	jar             *cookiejar.Jar
	samAddrString   string
	samPortString   string

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

func (samConn *SamHttp) initPipes() {
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

	samConn.base64Path, samConn.base64File, samConn.err = setupFile(filepath.Join(connectionDirectory, samConn.host), "id")
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

func (samConn *SamHttp) Dial(network, addr string) (net.Conn, error) {
	samConn.samBridgeClient, samConn.err = goSam.NewClientFromOptions(goSam.SetHost(samConn.samAddrString), goSam.SetPort(samConn.samPortString), goSam.SetDebug(DEBUG), goSam.SetUnpublished(true), goSam.SetInQuantity(15), goSam.SetOutQuantity(15))
	if samConn.c, samConn.err = Warn(samConn.err, "sam-http.go SAM connection error", "sam-http.go Initializing SAM connection"); samConn.c {
		return samConn.subDial(network, addr)
	}
	return samConn.samBridgeClient.SamConn, &errorString{"SAM connection error"}
}

func (samConn *SamHttp) subDial(network, addr string) (net.Conn, error) {
	if samConn.name != "" {
		if samConn.id != 0 {
			return samConn.Connect()
		} else {
			return nil, &errorString{"ID error"}
		}
	} else {
		return nil, &errorString{"Hostname error"}
	}
}

func (samConn *SamHttp) Connect() (net.Conn, error) {
	if samConn.samBridgeClient != nil {
		samConn.err = samConn.samBridgeClient.StreamConnect(samConn.id, samConn.name)
		if samConn.c, samConn.err = Warn(samConn.err, "sam-http.go Error connecting SAM streams", "sam-http.go Connecting SAM streams"); samConn.c {
			Log("sam-http.go Stream Connection established")
			return samConn.samBridgeClient.SamConn, samConn.err
		} else {
			return samConn.reConnect()
		}
	} else {
		return samConn.reConnect()
	}
}

func (samConn *SamHttp) reConnect() (net.Conn, error) {
	samConn.samBridgeClient, samConn.err = goSam.NewClientFromOptions(goSam.SetHost(samConn.samAddrString), goSam.SetPort(samConn.samPortString), goSam.SetDebug(DEBUG), goSam.SetUnpublished(true), goSam.SetInQuantity(15), goSam.SetOutQuantity(15))
	if samConn.c, samConn.err = Warn(samConn.err, "sam-http.go 133 SAM Client connection error", "sam-http.go SAM client connecting"); samConn.c {
		Log("sam-http.go SAM Connection established")
		samConn.err = samConn.samBridgeClient.StreamConnect(samConn.id, samConn.name)
		if samConn.c, samConn.err = Warn(samConn.err, "sam-http.go Connecting SAM streams", "sam-http.go Connecting SAM streams"); samConn.c {
			Log("sam-http.go Stream Connection established")
			return samConn.samBridgeClient.SamConn, samConn.err
		} else {
			return samConn.reConnect()
		}
	} else {
		//samConn.samBridgeClient.Close()
		return samConn.reConnect()
	}
}

func (samConn *SamHttp) checkRedirect(req *http.Request, via []*http.Request) error {
	return nil
}

func (samConn *SamHttp) setupTransport() {
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

func (samConn *SamHttp) createClient(request string, samAddrString string, samPortString string) {
	samConn.samAddrString = samAddrString
	samConn.samPortString = samPortString
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
			samConn.host, samConn.directory = samConn.hostSet(request)
			samConn.initPipes()
		}
		samConn.setName(request)
		samConn.subCache = append(samConn.subCache, NewSamUrl(samConn.directory))
	}
}

func (samConn *SamHttp) createClientHttp(request *http.Request, samAddrString string, samPortString string) {
	samConn.samAddrString = samAddrString
	samConn.samPortString = samPortString
	samConn.samBridgeClient, samConn.err = goSam.NewClientFromOptions(goSam.SetHost(samConn.samAddrString), goSam.SetPort(samConn.samPortString), goSam.SetDebug(DEBUG), goSam.SetUnpublished(true), goSam.SetInQuantity(15), goSam.SetOutQuantity(15))
	if samConn.c, samConn.err = Fatal(samConn.err, "sam-http.go 205 SAM Client Connection Error", "sam-http.go SAM client connecting", samConn.samAddrString, samConn.samPortString); samConn.c {
		Log("sam-http.go Setting Transport")
		Log("sam-http.go Setting Dial function")
		samConn.setupTransport()
		if samConn.host == "" {
			samConn.host, samConn.directory = samConn.hostSet(request.URL.String())
			samConn.initPipes()
		}
		samConn.setName(request.URL.String()) //, samConn.samBridgeClient)
		samConn.subCache = append(samConn.subCache, NewSamUrlHttp(request))
	}
}

func (samConn *SamHttp) cleanURL(request string) (string, string) {
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

func (samConn *SamHttp) hostSet(request string) (string, string) {
	host, req := samConn.cleanURL(request)
	Log("sam-http.go Setting up micro-proxy for:", "http://"+host)
	Log("sam-http.go in Directory", req)
	return host, req
}

func (samConn *SamHttp) hostGet() string {
	return "http://" + samConn.host
}

func (samConn *SamHttp) hostCheck(request string) bool {
	host, _ := samConn.cleanURL(request)
	_, err := url.ParseRequestURI(host)
	if err == nil {
		if samConn.host == host {
			Log("sam-http.go Request host ", host, "is equal to client host", samConn.host)
			return true
		} else {
			Log("sam-http.go Request host ", host, "is not equal to client host", samConn.host)
			return false
		}
	} else {
		if samConn.host == host {
			Log("sam-http.go Request host ", host, "is equal to client host", samConn.host)
			return true
		} else {
			Log("sam-http.go Request host ", host, "is not equal to client host", samConn.host)
			return false
		}
	}
}

func (samConn *SamHttp) getURL(request string) (string, string) {
	r := request
	directory := strings.Replace(request, "http://", "", -1)
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

func (samConn *SamHttp) sendRequest(request string) (*http.Response, error) {
	r, dir := samConn.getURL(request)
	Log("sam-http.go Getting resource", request)
	if samConn.subClient != nil {
		resp, err := samConn.subClient.Get(r)
		Warn(err, "sam-http.go Response Error", "sam-http.go Getting Response")
		Log("sam-http.go Pumping result to top of parent pipe")
		samConn.copyRequest(resp, dir)
		return resp, err
	} else {
		return nil, nil
	}
}

func (samConn *SamHttp) getURLHttp(request *http.Request) (string, string) {
	directory := strings.Replace(request.URL.String(), "http://", "", -1)
	return request.URL.String(), directory
}

func (samConn *SamHttp) sendRequestHttp(request *http.Request) (*http.Client, string) {
	r, dir := samConn.getURLHttp(request)
	Log("sam-http.go Getting resource", r, "In ", dir)
	return samConn.subClient, dir
}

func (samConn *SamHttp) sendRequestBase64Http(request *http.Request, base64helper string) (*http.Client, string) {
	r, dir := samConn.getURL(request.URL.String())
	Log("sam-http.go Getting resource", r, "In ", dir)
	return samConn.subClient, dir
}

func (samConn *SamHttp) findSubCache(response *http.Response, directory string) *samUrl {
	b := false
	var u samUrl
	for _, url := range samConn.subCache {
		Log("sam-http.go Seeking Subdirectory", url.subDirectory)
		if url.checkDirectory(directory) {
			return &url
		}
	}
	if b == false {
		Log("sam-http.go has not been retrieved yet. Setting up:", directory)
		samConn.subCache = append(samConn.subCache, NewSamUrl(directory))
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

func (samConn *SamHttp) copyRequest(response *http.Response, directory string) {
	samConn.findSubCache(response, directory).copyDirectory(response, directory)
}

func (samConn *SamHttp) copyRequestHttp(request *http.Request, response *http.Response, directory string) *http.Response {
	return samConn.findSubCache(response, directory).copyDirectoryHttp(request, response, directory)
}

func (samConn *SamHttp) scannerText() (string, error) {
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

/**/
func (samConn *SamHttp) responsify(input string) io.Reader {
	tmp := strings.NewReader(input)
	Log("sam-http.go Turning string into a response", input)
	return tmp
}

/**/
func (samConn *SamHttp) printResponse() string {
	s, e := samConn.scannerText()
	if samConn.c, samConn.err = Fatal(e, "sam-http.go Response Retrieval Error", "sam-http.go Retrieving Responses"); !samConn.c {
		Log("sam-http.go Response Panic")
		samConn.CleanupClient()
	} else {
		Log("sam-http.go Response Retrieved")
	}
	return s
}

func (samConn *SamHttp) readRequest() string {
	text := samConn.sendScan.Text()
	for samConn.sendScan.Scan() {
		samConn.sendRequest(text)
	}
	clearFile(filepath.Join(connectionDirectory, samConn.directory), "send")
	return text
}

/*
func (samConn *SamHttp) readDelete() bool {
	b := false
	for _, dir := range samConn.subCache {
		n := dir.readDelete()
		if n == 0 {
			Log("sam-http.go Maintaining Connection:", samConn.hostGet())
		} else if n > 0 {
			b = true
		}
	}
	return b
}
*/

func (samConn *SamHttp) writeName() {
	Log("sam-http.go Looking up hostname:", samConn.host)
	samConn.name, samConn.err = samConn.samBridgeClient.Lookup(samConn.host)
	samConn.nameFile.WriteString(samConn.name)
}

func (samConn *SamHttp) writeSession(request string) {
	Log("sam-http.go Caching base64 address of:", samConn.host+" "+samConn.name)
	samConn.id, samConn.base64, samConn.err = samConn.samBridgeClient.CreateStreamSession("")
	samConn.idFile.WriteString(fmt.Sprint(samConn.id))
	Warn(samConn.err, "sam-http.go Local Base64 Caching error", "sam-http.go Cachine Base64 Address of:", request)
	log.Println("sam-http.go Tunnel id: ", samConn.id)
	Log("sam-http.go Tunnel dest: ", samConn.base64)
	samConn.base64File.WriteString(samConn.base64)
	Log("sam-http.go New Connection Name: ", samConn.base64)
}

func (samConn *SamHttp) setName(request string) {
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

func (samConn *SamHttp) checkName() bool {
	Log("sam-http.go seeing if the connection needs a name:")
	if samConn.name != "" {
		Log("sam-http.go Naming connection: Connection name was empty.")
		return true
	} else {
		return false
	}
}

//export CleanupClient
func (samConn *SamHttp) CleanupClient() {
	samConn.sendPipe.Close()
	samConn.nameFile.Close()
	for _, url := range samConn.subCache {
		url.cleanupDirectory()
	}
	err := samConn.samBridgeClient.Close()
	if samConn.c, samConn.err = Fatal(err, "sam-http.go Closing SAM bridge error, retrying.", "sam-http.go Closing SAM bridge"); !samConn.c {
		samConn.samBridgeClient.Close()
	}
	os.RemoveAll(filepath.Join(connectionDirectory, samConn.host))
}

func newSamHttp(samAddrString, samPortString, request string, timeoutTime int, keepAlives bool) SamHttp {
	Log("sam-http.go Creating a new SAMv3 Client: ", request)
	var samConn SamHttp
	samConn.timeoutTime = time.Duration(timeoutTime) * time.Minute
	samConn.otherTimeoutTime = time.Duration(timeoutTime/3) * time.Minute
	samConn.keepAlives = keepAlives
	Log(request)
	samConn.createClient(request, samAddrString, samPortString)
	return samConn
}

func newSamHttpHttp(samAddrString, samPortString string, request *http.Request, timeoutTime int, keepAlives bool) SamHttp {
	Log("sam-http.go Creating a new SAMv3 Client.")
	var samConn SamHttp
	samConn.timeoutTime = time.Duration(timeoutTime) * time.Minute
	samConn.otherTimeoutTime = time.Duration(timeoutTime/3) * time.Minute
	samConn.keepAlives = keepAlives
	Log(request.Host + request.URL.Path)
	samConn.createClientHttp(request, samAddrString, samPortString)
	return samConn
}
