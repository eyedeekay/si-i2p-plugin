package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/eyedeekay/gosam"
	//"github.com/cryptix/goSam"
)

type samHttp struct {
	subCache []samUrl
	err      error
	c        bool

	samBridgeClient *goSam.Client
	samAddrString   string
	samPortString   string

	transport *http.Transport
	subClient *http.Client

	//Timeout time.Duration
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

func (samConn *samHttp) initPipes() {
    checkFolder(filepath.Join(connectionDirectory,samConn.host))

    samConn.sendPath, samConn.sendPipe, samConn.err = setupFiFo(filepath.Join(connectionDirectory, samConn.host), "send")
    if samConn.c, samConn.err = Fatal(samConn.err, "Pipe setup error", "Pipe setup"); samConn.c {
        samConn.sendScan, samConn.err = setupScanner(filepath.Join(connectionDirectory, samConn.host), "send", samConn.sendPipe)
        if samConn.c, samConn.err = Fatal(samConn.err, "Scanner setup Error:", "Scanner set up successfully."); !samConn.c {
            samConn.cleanupClient()
        }
    }

    samConn.namePath, samConn.nameFile, samConn.err = setupFiFo(filepath.Join(connectionDirectory, samConn.host), "name")
    if samConn.c, samConn.err = Fatal(samConn.err, "Pipe setup error", "Pipe setup"); samConn.c {
        samConn.nameFile.WriteString("")
    }

    samConn.idPath, samConn.idFile, samConn.err = setupFiFo(filepath.Join(connectionDirectory, samConn.host), "id")
    if samConn.c, samConn.err = Fatal(samConn.err, "Pipe setup error", "Pipe setup"); samConn.c {
        samConn.idFile.WriteString("")
    }

    samConn.base64Path, samConn.base64File, samConn.err = setupFiFo(filepath.Join(connectionDirectory, samConn.host), "id")
    if samConn.c, samConn.err = Fatal(samConn.err, "Pipe setup error", "Pipe setup"); samConn.c {
        samConn.idFile.WriteString("")
    }

}

func (samConn *samHttp) Dial(network, addr string) (net.Conn, error) {
	samCombined := samConn.samAddrString + ":" + samConn.samPortString
	samConn.samBridgeClient, samConn.err = goSam.NewClient(samCombined)
	if samConn.c, samConn.err = Warn(samConn.err, "SAM connection error", "Initializing SAM connection"); samConn.c {
		if samConn.name != "" {
			if samConn.id != 0 {
				samConn.err = samConn.samBridgeClient.StreamConnect(samConn.id, samConn.name)
				Warn(samConn.err, "Stream connection error", "Connecting SAM streams")
			}
		}
	}
	return samConn.samBridgeClient.SamConn, nil
}

func (samConn *samHttp) createClient(request string, samAddrString string, samPortString string) {
	samConn.samAddrString = samAddrString
	samConn.samPortString = samPortString
	samCombined := samConn.samAddrString + ":" + samConn.samPortString
	samConn.samBridgeClient, samConn.err = goSam.NewClient(samCombined)
	if samConn.c, samConn.err = Fatal(samConn.err, "SAM Client Connection Error", "SAM client connecting", samCombined); samConn.c {
        Log("Setting Transport")
        Log("Setting Dial function")
        samConn.transport = &http.Transport{
            Dial: samConn.Dial,
        }
        Log("Initializing sub-client")
        samConn.subClient = &http.Client{
            //Timeout: client.Timeout,
            Transport: samConn.transport}

        if samConn.host == "" {
            samConn.host, samConn.directory = samConn.hostSet(request)
            samConn.initPipes()
        }
        samConn.writeName(request)
        samConn.subCache = append(samConn.subCache, newSamUrl(samConn.directory))
    }
}

func (samConn *samHttp) createClientHttp(request *http.Request, samAddrString string, samPortString string) {
	samConn.samAddrString = samAddrString
	samConn.samPortString = samPortString
	samCombined := samConn.samAddrString + ":" + samConn.samPortString
	samConn.samBridgeClient, samConn.err = goSam.NewClient(samCombined)
	if samConn.c, samConn.err = Fatal(samConn.err, "SAM Client Connection Error", "SAM client connecting", samCombined); samConn.c {
        Log("Setting Transport")
        Log("Setting Dial function")
        samConn.transport = &http.Transport{
            Dial: samConn.Dial,
        }
        Log("Initializing sub-client")
        samConn.subClient = &http.Client{
            //Timeout: client.Timeout,
            Timeout:   time.Duration(600 * time.Second),
            Transport: samConn.transport}

        if samConn.host == "" {
            samConn.host, samConn.directory = samConn.hostSet(request.URL.String())
            samConn.initPipes()
        }
        samConn.writeName(request.URL.String()) //, samConn.samBridgeClient)
        samConn.subCache = append(samConn.subCache, newSamUrlHttp(request))
    }
}

func (samConn *samHttp) cleanURL(request string) (string, string) {
	Log("cleanURL Trim 0 " + request)
	//http://i2p-projekt.i2p/en/downloads
	url := strings.Replace(request, "http://", "", -1)
	Log("cleanURL Request URL " + url)
	//i2p-projekt.i2p/en/downloads
	host := strings.SplitAfter(url, ".i2p")[0]
	Log("cleanURL Trim 2 " + host)
	return host, url
}

func (samConn *samHttp) hostSet(request string) (string, string) {
	host, req := samConn.cleanURL(request)
	Log("Setting up micro-proxy for:", "http://"+host)
	Log("in Directory", req)
	return host, req
}

func (samConn *samHttp) hostGet() string {
	return "http://" + samConn.host
}

func (samConn *samHttp) hostCheck(request string) bool {
	host, _ := samConn.cleanURL(request)
	_, err := url.ParseRequestURI(host)
	if err == nil {
		if samConn.host == host {
			Log("Request host ", host)
			Log("Is equal to client host", samConn.host)
			return true
		} else {
			Log("Request host ", host)
			Log("Is not equal to client host", samConn.host)
			return false
		}
	} else {
		if samConn.host == host {
			Log("Request host ", host)
			Log("Is equal to client host", samConn.host)
			return true
		} else {
			Log("Request host ", host)
			Log("Is not equal to client host", samConn.host)
			return false
		}
	}
}

func (samConn *samHttp) getURL(request string) (string, string) {
	r := request
	directory := strings.Replace(request, "http://", "", -1)
	_, err := url.ParseRequestURI(r)
	if err != nil {
		r = "http://" + request
		Log("URL failed validation, correcting to:", r)
	} else {
		Log("URL passed validation:", request)
	}
	Log("Request will be managed in:", directory)
	return r, directory
}

func (samConn *samHttp) sendRequest(request string) (*http.Response, error) {
	r, dir := samConn.getURL(request)
	Log("Getting resource", request)
	resp, err := samConn.subClient.Get(r)
	Warn(err, "Response Error", "Getting Response")
	Log("Pumping result to top of parent pipe")
	samConn.copyRequest(resp, dir)
	return resp, err
}

func (samConn *samHttp) sendRequestHttp(request *http.Request) (*http.Client, string) {
	r, dir := samConn.getURL(request.URL.String())
	Log("Getting resource", r, "In ", dir)
	return samConn.subClient, dir
}

func (samConn *samHttp) findSubCache(response *http.Response, directory string) *samUrl {
	b := false
	var u samUrl
	for _, url := range samConn.subCache {
		Log("Seeking Subdirectory", url.subDirectory)
		if url.checkDirectory(directory) {
			return &url
		}
	}
	if b == false {
		Log("has not been retrieved yet. Setting up:", directory)
		samConn.subCache = append(samConn.subCache, newSamUrl(directory))
		for _, url := range samConn.subCache {
			Log("Seeking Subdirectory", url.subDirectory)
			if url.checkDirectory(directory) {
				u = url
				return &u
			}
		}
	}
	return &u
}

func (samConn *samHttp) copyRequest(response *http.Response, directory string) {
	samConn.findSubCache(response, directory).copyDirectory(response, directory)
}

func (samConn *samHttp) copyRequestHttp(request *http.Request, response *http.Response, directory string) *http.Response {
	return samConn.findSubCache(response, directory).copyDirectoryHttp(request, response, directory)
}

func (samConn *samHttp) scannerText() (string, error) {
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
func (samConn *samHttp) responsify(input string) io.Reader {
	tmp := strings.NewReader(input)
	Log("Turning string into a response", input)
	return tmp
}

/**/
func (samConn *samHttp) printResponse() string {
	s, e := samConn.scannerText()
	if samConn.c, samConn.err = Fatal(e, "Response Retrieval Error", "Retrieving Responses"); !samConn.c{
        samConn.cleanupClient()
    }
	return s
}

func (samConn *samHttp) readRequest() string {
	text := samConn.sendScan.Text()
	for samConn.sendScan.Scan() {
		samConn.sendRequest(text)
	}
	return text
}
/*
func (samConn *samHttp) readDelete() bool {
	b := false
	for _, dir := range samConn.subCache {
		n := dir.readDelete()
		if n == 0 {
			Log("Maintaining Connection:", samConn.hostGet())
		} else if n > 0 {
			b = true
		}
	}
	return b
}
*/
func (samConn *samHttp) writeName(request string) {
	if samConn.checkName() {
		samConn.host, samConn.directory = samConn.hostSet(request)
		Log("Setting hostname:", samConn.host)
		Log("Looking up hostname:", samConn.host)
		samConn.name, samConn.err = samConn.samBridgeClient.Lookup(samConn.host)
		samConn.nameFile.WriteString(samConn.name)
		Log("Caching base64 address of:", samConn.host+" "+samConn.name)
		samConn.id, samConn.base64, samConn.err = samConn.samBridgeClient.CreateStreamSession("")
		samConn.idFile.WriteString(fmt.Sprint(samConn.id))
		Warn(samConn.err, "Local Base64 Caching error", "Cachine Base64 Address of:", request)
		log.Println("Tunnel id: ", samConn.id)
		Log("Tunnel dest: ", samConn.base64)
		samConn.base64File.WriteString(samConn.base64)
		Log("New Connection Name: ", samConn.base64)
	} else {
		samConn.host, samConn.directory = samConn.hostSet(request)
		Log("Setting hostname:", samConn.host)
		samConn.initPipes()
		Log("Looking up hostname:", samConn.host)
		samConn.name, samConn.err = samConn.samBridgeClient.Lookup(samConn.host)
		Log("Caching base64 address of:", samConn.host+" "+samConn.name)
		samConn.nameFile.WriteString(samConn.name)
		samConn.id, samConn.base64, samConn.err = samConn.samBridgeClient.CreateStreamSession("")
		samConn.idFile.WriteString(fmt.Sprint(samConn.id))
		Warn(samConn.err, "Local Base64 Caching error", "Cachine Base64 Address of:", request)
		log.Println("Tunnel id: ", samConn.id)
		Log("Tunnel dest: ", samConn.base64)
		samConn.base64File.WriteString(samConn.base64)
		Log("New Connection Name: ", samConn.base64)
	}
}

func (samConn *samHttp) checkName() bool {
	Log("seeing if the connection needs a name:")
	if samConn.name != "" {
		Log("Naming connection: Connection name was empty.")
		return true
	} else {
		return false
	}
}

func (samConn *samHttp) cleanupClient() {
	samConn.sendPipe.Close()
	samConn.nameFile.Close()
	for _, url := range samConn.subCache {
		url.cleanupDirectory()
	}
	err := samConn.samBridgeClient.Close()
	if samConn.c, samConn.err = Fatal(err, "Closing SAM bridge error, retrying.", "Closing SAM bridge"); !samConn.c{
        samConn.samBridgeClient.Close()
    }
	os.RemoveAll(filepath.Join(connectionDirectory, samConn.host))
}

func newSamHttp(samAddrString string, samPortString string, request string) samHttp {
	log.Println("Creating a new SAMv3 Client: ", request)
	var samConn samHttp
	samConn.createClient(request, samAddrString, samPortString)
	return samConn
}

func newSamHttpHttp(samAddrString string, samPortString string, request *http.Request) samHttp {
	log.Println("Creating a new SAMv3 Client.")
	var samConn samHttp
	log.Println(request.Host + request.URL.Path)
	samConn.createClientHttp(request, samAddrString, samPortString)
	return samConn
}
