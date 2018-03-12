package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	//"time"
	"net/url"

	"github.com/eyedeekay/gosam"
	//"github.com/cryptix/goSam"
)

type samHttp struct {
	subCache []samUrl
	err      error

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
	sendScan bufio.Scanner

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
	pathConnectionExists, pathErr := exists(filepath.Join(connectionDirectory, samConn.host))
	//samConn.Log("Directory Check", filepath.Join(connectionDirectory, samConn.host))
	samConn.Fatal(pathErr, "Directory Check Error", "Directory Check", filepath.Join(connectionDirectory, samConn.host))
	if !pathConnectionExists {
		samConn.Log("Creating a connection:", samConn.host)
		os.Mkdir(filepath.Join(connectionDirectory, samConn.host), 0755)
	}

	samConn.sendPath = filepath.Join(connectionDirectory, samConn.host, "send")
	pathSendExists, sendPathErr := exists(samConn.sendPath)
	samConn.Fatal(sendPathErr, "Path Check Error", "Path Check", samConn.sendPath)
	if !pathSendExists {
		err := syscall.Mkfifo(samConn.sendPath, 0755)
		samConn.Log("Preparing to create Pipe:", samConn.sendPath)
		samConn.Fatal(err, "File Creation Error", "Creating File", samConn.sendPath)
		samConn.Log("checking for problems...")
		samConn.sendPipe, err = os.OpenFile(samConn.sendPath, os.O_RDWR|os.O_CREATE, 0755)
		samConn.Log("Opening the Named Pipe as a File...")
		samConn.sendScan = *bufio.NewScanner(samConn.sendPipe)
		samConn.sendScan.Split(bufio.ScanLines)
		samConn.Log("Opening the Named Pipe as a Buffer...")
		samConn.Log("Created a named Pipe for sending requests:", samConn.sendPath)
	}

	samConn.namePath = filepath.Join(connectionDirectory, samConn.host, "name")
	pathNameExists, recvNameErr := exists(samConn.namePath)
	samConn.Fatal(recvNameErr, "Name File Check error", "Name File Check", samConn.namePath)
	if !pathNameExists {
		samConn.nameFile, samConn.err = os.Create(samConn.namePath)
		samConn.Log("Preparing to create File:", samConn.namePath)
		samConn.Fatal(samConn.err, "File Creation Error", "Creating File", samConn.namePath)
		samConn.Log("checking for problems...")
		samConn.Log("Opening the File...")
		samConn.nameFile, samConn.err = os.OpenFile(samConn.namePath, os.O_RDWR|os.O_CREATE, 0644)
		samConn.Log("Created a File for the full name:", samConn.namePath)
	}

	samConn.idPath = filepath.Join(connectionDirectory, samConn.host, "id")
	pathIdExists, recvIdErr := exists(samConn.idPath)
	samConn.Fatal(recvIdErr, "ID File Check Error", "ID File Check", samConn.idPath)
	if !pathIdExists {
		samConn.idFile, samConn.err = os.Create(samConn.idPath)
		samConn.Log("Preparing to create File:", samConn.idPath)
		samConn.Fatal(samConn.err, "File Creation Error", "Creating File", samConn.idPath)
		samConn.Log("checking for problems...")
		samConn.Log("Opening the File...")
		samConn.idFile, samConn.err = os.OpenFile(samConn.idPath, os.O_RDWR|os.O_CREATE, 0644)
		samConn.Log("Created a File for the full id:", samConn.idPath)
	}

	samConn.base64Path = filepath.Join(connectionDirectory, samConn.host, "base64")
	pathBase64Exists, recvBase64Err := exists(samConn.base64Path)
	samConn.Fatal(recvBase64Err, "Local Base64 File Check Error", "Local Base64 File Check", samConn.base64Path)
	if !pathBase64Exists {
		samConn.base64File, samConn.err = os.Create(samConn.base64Path)
		samConn.Log("Preparing to create File:", samConn.base64Path)
		samConn.Fatal(samConn.err, "File Creation Error", "Creating File", samConn.base64Path)
		samConn.Log("checking for problems...")
		samConn.Log("Opening the File...")
		samConn.base64File, samConn.err = os.OpenFile(samConn.base64Path, os.O_RDWR|os.O_CREATE, 0644)
		samConn.Log("Created a File for the full local base64:", samConn.base64Path)
	}

}

func (samConn *samHttp) Dial(network, addr string) (net.Conn, error) {
	samCombined := samConn.samAddrString + ":" + samConn.samPortString
	samConn.samBridgeClient, samConn.err = goSam.NewClient(samCombined)
	samConn.Fatal(samConn.err, "SAM connection error", "Initializing SAM connection")
	samConn.err = samConn.samBridgeClient.StreamConnect(samConn.id, samConn.name)
	samConn.Warn(samConn.err, "Stream connection error", "Connecting SAM streams")
	return samConn.samBridgeClient.SamConn, nil
	//return samConn.samBridgeClient.Dial(network, addr)
}

func (samConn *samHttp) createClient(request string, samAddrString string, samPortString string) {
	samConn.samAddrString = samAddrString
	samConn.samPortString = samPortString
	samCombined := samConn.samAddrString + ":" + samConn.samPortString
	samConn.samBridgeClient, samConn.err = goSam.NewClient(samCombined)
	samConn.Fatal(samConn.err, "SAM Client Connection Error", "SAM client connecting", samCombined)
	samConn.Log("Setting Transport")
	samConn.Log("Setting Dial function")
	samConn.transport = &http.Transport{
		Dial: samConn.Dial,
	}
	samConn.Log("Initializing sub-client")
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

func (samConn *samHttp) createClientHttp(request *http.Request, samAddrString string, samPortString string) {
	samConn.samAddrString = samAddrString
	samConn.samPortString = samPortString
	samCombined := samConn.samAddrString + ":" + samConn.samPortString
	samConn.samBridgeClient, samConn.err = goSam.NewClient(samCombined)
	samConn.Fatal(samConn.err, "SAM Connection Error", "Connecting to SAM Bridge", samCombined)
	samConn.Log("Setting Transport")
	samConn.Log("Setting Dial function")
	samConn.transport = &http.Transport{
		Dial: samConn.Dial,
	}
	samConn.Log("Initializing sub-client")
	samConn.subClient = &http.Client{
		//Timeout: client.Timeout,
		Transport: samConn.transport}

	if samConn.host == "" {
		samConn.host, samConn.directory = samConn.hostSet(request.URL.String())
		samConn.initPipes()
	}
	samConn.writeName(request.URL.String()) //, samConn.samBridgeClient)
	samConn.subCache = append(samConn.subCache, newSamUrlHttp(request))
}

func (samConn *samHttp) cleanURL(request string) (string, string) {
	samConn.Log("cleanURL Trim 0 " + request)
	//http://i2p-projekt.i2p/en/downloads
	url := strings.Replace(request, "http://", "", -1)
	samConn.Log("cleanURL Request URL " + url)
	//i2p-projekt.i2p/en/downloads
	host := strings.SplitAfter(url, ".i2p")[0]
	samConn.Log("cleanURL Trim 2 " + host)
	return host, url
}

func (samConn *samHttp) hostSet(request string) (string, string) {
	host, req := samConn.cleanURL(request)
	samConn.Log("Setting up micro-proxy for:", "http://"+host)
	samConn.Log("in Directory", req)
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
			samConn.Log("Request host ", host)
			samConn.Log("Is equal to client host", samConn.host)
			return true
		} else {
			samConn.Log("Request host ", host)
			samConn.Log("Is not equal to client host", samConn.host)
			return false
		}
	} else {
		if samConn.host == host {
			samConn.Log("Request host ", host)
			samConn.Log("Is equal to client host", samConn.host)
			return true
		} else {
			samConn.Log("Request host ", host)
			samConn.Log("Is not equal to client host", samConn.host)
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
		samConn.Log("URL failed validation, correcting to:", r)
	} else {
		samConn.Log("URL passed validation:", request)
	}
	samConn.Log("Request will be managed in:", directory)
	return r, directory
}

func (samConn *samHttp) sendRequest(request string) (*http.Response, error) {
	r, dir := samConn.getURL(request)
	samConn.Log("Getting resource", request)
	resp, err := samConn.subClient.Get(r)
	samConn.Warn(err, "Response Error", "Getting Response")
	samConn.Log("Pumping result to top of parent pipe")
	samConn.copyRequest(resp, dir)
	return resp, err
}

func (samConn *samHttp) sendRequestHttp(request *http.Request) (*http.Client, string) {
	r, dir := samConn.getURL(request.URL.String())
	samConn.Log("Getting resource", r, "In ", dir)
	return samConn.subClient, dir
}

func (samConn *samHttp) findSubCache(response *http.Response, directory string) *samUrl {
	b := false
	var u samUrl
	for _, url := range samConn.subCache {
		samConn.Log("Seeking Subdirectory", url.subDirectory)
		if url.checkDirectory(directory) {
			return &url
		}
	}
	if b == false {
		samConn.Log("has not been retrieved yet. Setting up:", directory)
		samConn.subCache = append(samConn.subCache, newSamUrl(directory))
		for _, url := range samConn.subCache {
			samConn.Log("Seeking Subdirectory", url.subDirectory)
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
	samConn.Log("Turning string into a response", input)
	return tmp
}

/**/
func (samConn *samHttp) printResponse() string {
	s, e := samConn.scannerText()
	samConn.Fatal(e, "Response Retrieval Error", "Retrieving Responses")
	return s
}

func (samConn *samHttp) readRequest() string {
	text := samConn.sendScan.Text()
	for samConn.sendScan.Scan() {
		samConn.sendRequest(text)
	}
	return text
}

func (samConn *samHttp) readDelete() bool {
	b := false
	for _, dir := range samConn.subCache {
		n := dir.readDelete()
		if n == 0 {
			samConn.Log("Maintaining Connection:", samConn.hostGet())
		} else if n > 0 {
			b = true
		}
	}
	return b
}

func (samConn *samHttp) writeName(request string) {
	if samConn.checkName() {
		samConn.host, samConn.directory = samConn.hostSet(request)
		samConn.Log("Setting hostname:", samConn.host)
		samConn.Log("Looking up hostname:", samConn.host)
		samConn.name, samConn.err = samConn.samBridgeClient.Lookup(samConn.host)
		samConn.nameFile.WriteString(samConn.name)
		samConn.Log("Caching base64 address of:", samConn.host+" "+samConn.name)
		samConn.id, samConn.base64, samConn.err = samConn.samBridgeClient.CreateStreamSession("")
		samConn.idFile.WriteString(fmt.Sprint(samConn.id))
		samConn.Warn(samConn.err, "Local Base64 Caching error", "Cachine Base64 Address of:", request)
		log.Println("Tunnel id: ", samConn.id)
		samConn.Log("Tunnel dest: ", samConn.base64)
		samConn.base64File.WriteString(samConn.base64)
		samConn.Log("New Connection Name: ", samConn.base64)
	} else {
		samConn.host, samConn.directory = samConn.hostSet(request)
		samConn.Log("Setting hostname:", samConn.host)
		samConn.initPipes()
		samConn.Log("Looking up hostname:", samConn.host)
		samConn.name, samConn.err = samConn.samBridgeClient.Lookup(samConn.host)
		samConn.Log("Caching base64 address of:", samConn.host+" "+samConn.name)
		samConn.nameFile.WriteString(samConn.name)
		samConn.id, samConn.base64, samConn.err = samConn.samBridgeClient.CreateStreamSession("")
		samConn.idFile.WriteString(fmt.Sprint(samConn.id))
		samConn.Warn(samConn.err, "Local Base64 Caching error", "Cachine Base64 Address of:", request)
		log.Println("Tunnel id: ", samConn.id)
		samConn.Log("Tunnel dest: ", samConn.base64)
		samConn.base64File.WriteString(samConn.base64)
		samConn.Log("New Connection Name: ", samConn.base64)
	}
}

func (samConn *samHttp) checkName() bool {
	samConn.Log("seeing if the connection needs a name:")
	if samConn.name != "" {
		samConn.Log("Naming connection: Connection name was empty.")
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
	samConn.Fatal(err, "Closing SAM Bridge Error", "Closing SAM Bridge")
	os.RemoveAll(filepath.Join(connectionDirectory, samConn.host))
}

func (samConn *samHttp) Log(msg ...string) {
	if verbose {
		log.Println("LOG: ", msg)
	}
}

func (samConn *samHttp) Warn(err error, errmsg string, msg ...string) bool {
	log.Println(msg)
	if err != nil {
		log.Println("WARN: ", err)
		return false
	}
	samConn.err = nil
	return true
}

func (samConn *samHttp) Fatal(err error, errmsg string, msg ...string) {
	log.Println(msg)
	if err != nil {
		samConn.err = err
		defer samConn.cleanupClient()
		log.Fatal("FATAL: ", errmsg, err)
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
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
