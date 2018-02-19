package main

import (
    "bufio"
	"io"
	"log"
	"net/http"
    "os"
    "path/filepath"
    "strings"
    "strconv"
    "syscall"
    "net/url"

	"github.com/eyedeekay/gosam"
)

type samHttp struct{
    subCache []samUrl
    //sam *goSam.Client
    err error

    transport *http.Transport
    http *http.Client
    host string

    sendPath string
    sendPipe *os.File
    sendBuff bufio.Reader

    namePath string
    nameFile *os.File
    name string
}

var connectionDirectory string

func (samConn *samHttp) initPipes(){
    pathConnectionExists, pathErr := exists(filepath.Join(connectionDirectory, samConn.host))
    log.Println("Directory Check", filepath.Join(connectionDirectory, samConn.host))
    samConn.checkErr(pathErr)
    if ! pathConnectionExists {
        log.Println("Creating a connection:", samConn.host)
        os.Mkdir(filepath.Join(connectionDirectory, samConn.host), 0755)
    }

    samConn.sendPath = filepath.Join(connectionDirectory, samConn.host, "send")
    pathSendExists, sendPathErr := exists(samConn.sendPath)
    samConn.checkErr(sendPathErr)
    if ! pathSendExists {
        err := syscall.Mkfifo(samConn.sendPath, 0755)
        log.Println("Preparing to create Pipe:", samConn.sendPath)
        samConn.checkErr(err)
        log.Println("checking for problems...")
        samConn.sendPipe, err = os.OpenFile(samConn.sendPath , os.O_RDWR|os.O_CREATE, 0755)
        log.Println("Opening the Named Pipe as a File...")
        samConn.sendBuff = *bufio.NewReader(samConn.sendPipe)
        log.Println("Opening the Named Pipe as a Buffer...")
        log.Println("Created a named Pipe for sending requests:", samConn.sendPath)
    }

    samConn.namePath = filepath.Join(connectionDirectory, samConn.host, "name")
    pathNameExists, recvNameErr := exists(samConn.namePath)
    samConn.checkErr(recvNameErr)
    if ! pathNameExists {
        samConn.nameFile, samConn.err = os.Create(samConn.namePath)
        log.Println("Preparing to create File:", samConn.namePath)
        samConn.checkErr(samConn.err)
        log.Println("checking for problems...")
        log.Println("Opening the File...")
        samConn.nameFile, samConn.err = os.OpenFile(samConn.namePath, os.O_RDWR|os.O_CREATE, 0644)
        log.Println("Created a File for the full name:", samConn.namePath)
    }

}


func (samConn *samHttp) createClient(request string) {
    samConn.http = &http.Client{Transport: samConn.transport}
    if samConn.host == "" {
        samConn.host, _ = samConn.hostSet(request)
        samConn.initPipes()
    }
    samConn.writeName(request)
    samConn.subCache = append(samConn.subCache, newSamUrl(samConn.host))
}

func (samConn *samHttp) createClientHttp(request *http.Request) {
    samConn.http = &http.Client{Transport: samConn.transport}
    if samConn.host == "" {
        samConn.host, _ = samConn.hostSet(request.Host)
        samConn.initPipes()
    }
    samConn.writeName(request.Host)
    samConn.subCache = append(samConn.subCache, newSamUrlHttp(request))
}

func (samConn *samHttp) hostSet(request string) (string, string){
    tmp := strings.Replace(request, "http://", "", -1)
    host := strings.SplitAfterN(tmp, ".i2p", 1 )[0]
    _, err := url.ParseRequestURI("http://" + host)
    if err != nil {
        host = strings.Replace(host, "http://", "", -1)
    }
    directory := strings.Replace(tmp, host, "", -1)
    log.Println("Setting up micro-proxy for:", "http://" + host)
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

func (samConn *samHttp) getURLHttp(req *http.Request) (string, string){
    request := req.URL.String()
    //tmp := strings.SplitAfterN(request, ".i2p", -1)
    directory := strings.Replace(request, "http://", "", -1)
    _, err := url.ParseRequestURI(req.Host)
    if err != nil {
        log.Println("URL failed validation, correcting to:", request)
    }else{
        log.Println("URL passed validation:", request)
    }
    return req.Host, directory
}

func (samConn *samHttp) sendRequest(request string) (*http.Response, error ){
    r, dir := samConn.getURL(request)
    resp, err := samConn.http.Get(r)
    samConn.checkErr(err)
    samConn.copyRequest(resp, dir)
    return resp, err
}

func (samConn *samHttp) sendRequestHttp(request *http.Request) (*http.Response, error ){
    r, dir := samConn.getURLHttp(request)
    resp, err := samConn.http.Get(r)
    samConn.checkErr(err)
    samConn.copyRequest(resp, dir)
    return resp, err
}

func (samConn *samHttp) copyRequest(response *http.Response, directory string){
    b := false
    for _, url := range samConn.subCache {
        b = url.copyDirectory(response, directory)
        if b == true {
            break
        }
    }
    if b == false {
        log.Println("%s has not been retrieved yet. Setting up:", directory)
        samConn.subCache = append(samConn.subCache, newSamUrl(directory))
        for _, url := range samConn.subCache {
            b = url.copyDirectory(response, directory)
            if b == true {
                break
            }
        }
    }
}

func (samConn *samHttp) scannerText() (r string, l int) {
   text := ""
   length := 0
    for _, url := range samConn.subCache {
        text, length = url.scannerText()
        if length > 0 {
            break
        }
    }
    return text, length
}

func (samConn *samHttp) responsify(input string) io.Reader {
    tmp := strings.NewReader(input)
    log.Println("Turning string %s into a response", input)
    return tmp
}

func (samConn *samHttp) printResponse() string{
    s, n := samConn.scannerText()
    if n == 0 {
        log.Println("Maintaining Connection:", samConn.hostGet())
        return ""
    }else if n < 0 {
        log.Println("Something wierd happened with :" , s)
        return ""
    }else{
        //io.Copy(samConn.recvPipe, samConn.responsify(s))
        return s
    }
}

func (samConn *samHttp) readRequest(){
    line, _, err := samConn.sendBuff.ReadLine()
    samConn.checkErr(err)
    n := len(line)
    log.Println("Reading n bytes from send pipe:", strconv.Itoa(n))
    if n == 0 {
        log.Println("Maintaining Connection:", samConn.hostGet())
    }else if n < 0 {
        log.Println("Something wierd happened with :", line)
        log.Println("end determined at index :", strconv.Itoa(n))
    }else{
        s := string( line[:n] )
        log.Println("Sending request:", s)
        samConn.sendRequest(s)
    }
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


func (samConn *samHttp) writeName(request string){
    if samConn.host == "" {
        log.Println("Setting hostname:" )
        //directory := ""
        samConn.host, _ = samConn.hostSet(request)
        samConn.initPipes()
    }
    log.Println("Attempting to write-out connection name:")
    if samConn.checkName() {
        //samConn.name, samConn.err = samConn.sam.Lookup(samConn.host)
        //samConn.name, samConn.err = samConn.sam.Lookup(samConn.host)
        log.Println("New Connection Name: %s", samConn.name)
        samConn.checkErr(samConn.err)
        samConn.nameFile.WriteString(samConn.name)
    }
}

func (samConn *samHttp) checkName() bool{
    log.Println("seeing if the connection needs a name:")
    if samConn.name == "" {
        log.Println("Naming connection:")
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
    /*err := samConn.sam.Close()
    if err != nil {
        log.Println(err)
    }*/
    os.RemoveAll(filepath.Join(connectionDirectory, samConn.host))
}

func (samConn *samHttp) checkErr(err error) {
	if err != nil {
        log.Println(err)
        samConn.cleanupClient()
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
    samConn.name = ""
    samConn.host = ""
    log.Println(request)
    log.Println("Setting Dial function")
    samConn.transport = &http.Transport{
		Dial: sam.Dial,
	}
    samConn.createClient(request)
    return samConn
}

func newSamHttpHttp(samAddrString string, samPortString string, sam *goSam.Client, request *http.Request) (samHttp){
    log.Println("Creating a new SAMv3 Client.")
    var samConn samHttp
    samConn.name = ""
    samConn.host = ""
    log.Println(request.Host + request.URL.Path)
    log.Println("Setting Dial function")
    samConn.transport = &http.Transport{
		Dial: sam.Dial,
	}
    samConn.createClientHttp(request)
    return samConn
}
