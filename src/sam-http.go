package main

import (
    "bufio"
	"io"
	"log"
	"net/http"
    "os"
    "path/filepath"
    "strings"
    "syscall"
    //"time"
    "net/url"

	"github.com/eyedeekay/gosam"
)

type samHttp struct{
    subCache []samUrl
    err error

    transport *http.Transport
    subClient *http.Client
    host string
    directory string

    sendPath string
    sendPipe *os.File
    sendScan bufio.Scanner

    namePath string
    nameFile *os.File
    name string
}

var connectionDirectory string

func (samConn *samHttp) initPipes(){
    pathConnectionExists, pathErr := exists(filepath.Join(connectionDirectory, samConn.host))
    log.Println("Directory Check", filepath.Join(connectionDirectory, samConn.host))
    samConn.Fatal(pathErr)
    if ! pathConnectionExists {
        log.Println("Creating a connection:", samConn.host)
        os.Mkdir(filepath.Join(connectionDirectory, samConn.host), 0755)
    }

    samConn.sendPath = filepath.Join(connectionDirectory, samConn.host, "send")
    pathSendExists, sendPathErr := exists(samConn.sendPath)
    samConn.Fatal(sendPathErr)
    if ! pathSendExists {
        err := syscall.Mkfifo(samConn.sendPath, 0755)
        log.Println("Preparing to create Pipe:", samConn.sendPath)
        samConn.Fatal(err)
        log.Println("checking for problems...")
        samConn.sendPipe, err = os.OpenFile(samConn.sendPath , os.O_RDWR|os.O_CREATE, 0755)
        log.Println("Opening the Named Pipe as a File...")
        samConn.sendScan = *bufio.NewScanner(samConn.sendPipe)
        samConn.sendScan.Split(bufio.ScanLines)
        log.Println("Opening the Named Pipe as a Buffer...")
        log.Println("Created a named Pipe for sending requests:", samConn.sendPath)
    }

    samConn.namePath = filepath.Join(connectionDirectory, samConn.host, "name")
    pathNameExists, recvNameErr := exists(samConn.namePath)
    samConn.Fatal(recvNameErr)
    if ! pathNameExists {
        samConn.nameFile, samConn.err = os.Create(samConn.namePath)
        log.Println("Preparing to create File:", samConn.namePath)
        samConn.Fatal(samConn.err)
        log.Println("checking for problems...")
        log.Println("Opening the File...")
        samConn.nameFile, samConn.err = os.OpenFile(samConn.namePath, os.O_RDWR|os.O_CREATE, 0644)
        log.Println("Created a File for the full name:", samConn.namePath)
    }

}


func (samConn *samHttp) createClient(request string, sam *goSam.Client) {
    log.Println("Setting Transport")
    log.Println("Setting Dial function")
    samConn.transport = &http.Transport{
		Dial: sam.Dial,
	}
    log.Println("Initializing sub-client")
    samConn.subClient = &http.Client{
        Transport: samConn.transport    }

    if samConn.host == "" {
        samConn.host, samConn.directory = samConn.hostSet(request)
        samConn.initPipes()
    }
    samConn.writeName(request, sam)
    samConn.subCache = append(samConn.subCache, newSamUrl(samConn.directory))
}

func (samConn *samHttp) createClientHttp(request *http.Request, sam *goSam.Client) {
    log.Println("Setting Transport")
    log.Println("Setting Dial function")
    samConn.transport = &http.Transport{
		Dial: sam.Dial,
	}
    log.Println("Initializing sub-client")
    samConn.subClient = &http.Client{
        //Timeout: time.Second * 10,
        Transport: samConn.transport    }

    if samConn.host == "" {
        samConn.host, samConn.directory = samConn.hostSet(request.URL.String())
        samConn.initPipes()
    }
    samConn.writeName(request.URL.String(), sam)
    samConn.subCache = append(samConn.subCache, newSamUrlHttp(request))
}

func (samConn *samHttp) cleanURL(request string) (string, string){
    log.Println("cleanURL Trim 0 " + request)
    //http://i2p-projekt.i2p/en/downloads
    url := strings.Replace(request, "http://", "", -1)
    log.Println("cleanURL Request URL " + url)
    //i2p-projekt.i2p/en/downloads
    host := strings.SplitAfter(url, ".i2p")[0]
    log.Println("cleanURL Trim 2 " + host)
    //i2p-projekt.i2p
    return host, url
}

func (samConn *samHttp) hostSet(request string) (string, string){
    host, req := samConn.cleanURL(request)
    _, err := url.ParseRequestURI("http://" + host)
    if err != nil {
        host = strings.Replace(host, "http://", "", -1)
    }
    //directory := strings.Replace(req, host + "/", "", -1) + "/"
    //directory := strings.Replace(req, host + "/", "", -1) + "/"
    directory := req
    log.Println("Setting up micro-proxy for:", "http://" + host)
    log.Println("in Directory", directory)
    return host, directory
}

func (samConn *samHttp) hostGet() string{
    return "http://" + samConn.host
}

func (samConn *samHttp) hostCheck(request string) bool{
    host := strings.SplitAfterN(request, ".i2p", -1 )[0]
    _, err := url.ParseRequestURI(host)
    if err == nil {
            comphost := strings.Replace(host, "http://", "", -1)
            comphost = strings.SplitAfterN(request, ".i2p", -1 )[0]
            comphost = strings.Replace(host, "http://", "", -1)
        if samConn.host == comphost {
            log.Println("Request host ", comphost)
            log.Println("Is equal to client host", samConn.host)
            return true
        }else{
            log.Println("Request host ", comphost)
            log.Println("Is not equal to client host", samConn.host)
            return false
        }
    }else{
        host = strings.Replace(host, "http://", "", -1)
        host = strings.SplitAfterN(request, ".i2p", -1 )[0]
        host = strings.Replace(host, "http://", "", -1)
        if samConn.host == host {
            log.Println("Request host ", host)
            log.Println("Is equal to client host", samConn.host)
            return true
        }else{
            log.Println("Request host ", host)
            log.Println("Is not equal to client host", samConn.host)
            return false
        }
    }
}

func (samConn *samHttp) hostCheckHttp(req *http.Request) bool{
    request := req.Host
    host := strings.SplitAfterN(request, ".i2p", -1 )[0]
    _, err := url.ParseRequestURI(host)
    if err == nil {
            comphost := strings.Replace(host, "http://", "", -1)
            comphost = strings.SplitAfterN(request, ".i2p", -1 )[0]
            comphost = strings.Replace(host, "http://", "", -1)
        if samConn.host == comphost {
            log.Println("Request host ", comphost)
            log.Println("Is equal to client host", samConn.host)
            return true
        }else{
            log.Println("Request host ", comphost)
            log.Println("Is not equal to client host", samConn.host)
            return false
        }
    }else{
        host = strings.Replace(host, "http://", "", -1)
        host = strings.SplitAfterN(request, ".i2p", -1 )[0]
        host = strings.Replace(host, "http://", "", -1)
        if samConn.host == host {
            log.Println("Request host ", host)
            log.Println("Is equal to client host", samConn.host)
            return true
        }else{
            log.Println("Request host ", host)
            log.Println("Is not equal to client host", samConn.host)
            return false
        }
    }
}

func (samConn *samHttp) getURL(request string) (string, string){
    host := request
    //tmp := strings.SplitAfterN(request, ".i2p", -1)
    directory := strings.Replace(request, "http://", "", -1)
    _, err := url.ParseRequestURI(host)
    if err != nil {
        host = "http://" + request
        log.Println("URL failed validation, correcting to:", host)
    }else{
        log.Println("URL passed validation:", request)
    }
    return host, directory
}

func (samConn *samHttp) getURLHttp(req *http.Request) (*http.Request, string){
    request := req.URL.String()
    //tmp := strings.SplitAfterN(request, ".i2p", -1)
    directory := strings.Replace(request, "http://", "", -1)
    _, err := url.ParseRequestURI(req.URL.String())
    if err != nil {
        log.Println("URL failed validation, correcting to:", request)
    }else{
        log.Println("URL passed validation:", request)
    }
    return req, directory
}

func (samConn *samHttp) sendRequest(request string) (*http.Response, error ){
    r, dir := samConn.getURL(request)
    log.Println("Getting resource", request)
    resp, err := samConn.subClient.Get(r)
    samConn.Warn(err)
    log.Println("Pumping result to top of parent pipe")
    samConn.copyRequest(resp, dir)
    return resp, err
}

func (samConn *samHttp) sendRequestHttp(request *http.Request) (*http.Client, string){
    r, dir := samConn.getURLHttp(request)
    log.Println("Getting resource", r.URL.String())
    log.Println("In ", dir)
    return samConn.subClient, dir
}

func (samConn *samHttp) copyRequest(response *http.Response, directory string){
    b := false
    for _, url := range samConn.subCache {
        log.Println("Seeking Subdirectory", url.subDirectory)
        b = url.copyDirectory(response, directory)
        if b == true {
            log.Println("Found Subdirectory", url.subDirectory)
            break
        }
    }
    if b == false {
        log.Println("has not been retrieved yet. Setting up:", directory)
        samConn.subCache = append(samConn.subCache, newSamUrl(directory))
        for _, url := range samConn.subCache {
            log.Println("Seeking Subdirectory", url.subDirectory)
            b = url.copyDirectory(response, directory)
            if b == true {
                log.Println("Found Subdirectory", url.subDirectory)
                break
            }
        }
    }
}

func (samConn *samHttp) copyRequestHttp(request *http.Request, response *http.Response, directory string)(*http.Response){
    b := false
    for _, url := range samConn.subCache {
        log.Println("Seeking Subdirectory", url.subDirectory)
        b, d := url.copyDirectoryHttp(request ,response, directory)
        if b == true {
            log.Println("Found Subdirectory", url.subDirectory)
            return d
        }
    }
    if b == false {
        log.Println("has not been retrieved yet. Setting up:", directory)
        samConn.subCache = append(samConn.subCache, newSamUrl(directory))
        for _, url := range samConn.subCache {
            log.Println("Seeking Subdirectory", url.subDirectory)
            b, d := url.copyDirectoryHttp(request ,response, directory)
            if b == true {
                log.Println("Found Subdirectory", url.subDirectory)
                return d
            }
        }
    }
    return response
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
    log.Println("Turning string into a response", input)
    return tmp
}
/**/
func (samConn *samHttp) printResponse() string{
    s, e := samConn.scannerText()
    samConn.Fatal(e)
    return s
}

func (samConn *samHttp) readRequest() string{
    text := samConn.sendScan.Text()
    for samConn.sendScan.Scan(){
        samConn.sendRequest(text)
    }
    return text
}

func (samConn *samHttp) readDelete() bool {
    b := false
    for _, dir := range samConn.subCache {
        n := dir.readDelete()
        if n == 0 {
            log.Println("Maintaining Connection:", samConn.hostGet())
        }else if n > 0 {
            b = true
        }
    }
    return b
}

func (samConn *samHttp) writeName(request string, sam *goSam.Client){
    if samConn.checkName() {
        log.Println("Looking up hostname:", samConn.host )
        samConn.name, samConn.err = sam.Lookup(samConn.host)
        log.Println("New Connection Name: ", samConn.host)
        log.Println("Caching base64 address of:", samConn.host )
        samConn.Warn(samConn.err)
        samConn.nameFile.WriteString(samConn.name)
    }else{
        samConn.host, samConn.directory = samConn.hostSet(request)
        log.Println("Setting hostname:", samConn.host )
        samConn.initPipes()
        samConn.name, samConn.err = sam.Lookup(samConn.host)
        log.Println("New Connection Name: ", samConn.host)
        log.Println("Caching base64 address of:", samConn.host )
        samConn.Warn(samConn.err)
        samConn.nameFile.WriteString(samConn.name)
    }
}

func (samConn *samHttp) checkName() bool{
    log.Println("seeing if the connection needs a name:")
    if samConn.name != "" {
        log.Println("Naming connection: Connection name was empty.")
        return true
    }else{
        return false
    }
}

func (samConn *samHttp) cleanupClient(){
    samConn.sendPipe.Close()
    samConn.nameFile.Close()
    for _, url := range samConn.subCache {
        url.cleanupDirectory()
    }
    os.RemoveAll(filepath.Join(connectionDirectory, samConn.host))
}

func (samConn *samHttp) Warn(err error) {
	if err != nil {
        log.Println("Warning: ", err)
	}
}

func (samConn *samHttp) Fatal(err error) {
	if err != nil {
        defer samConn.cleanupClient()
        log.Fatal("Fatal: ", err)
	}
}

func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

func newSamHttp(samAddrString string, samPortString string, sam *goSam.Client, request string) (samHttp){
    log.Println("Creating a new SAMv3 Client.")
    var samConn samHttp
    log.Println(request)
    samConn.createClient(request, sam)
    return samConn
}

func newSamHttpHttp(samAddrString string, samPortString string, sam *goSam.Client, request *http.Request) (samHttp){
    log.Println("Creating a new SAMv3 Client.")
    var samConn samHttp
    log.Println(request.Host + request.URL.Path)
    samConn.createClientHttp(request, sam)
    //client.Do(request)
    //samConn.createClient(request.URL.String(), sam)
    return samConn
}
