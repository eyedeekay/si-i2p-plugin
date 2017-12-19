package main

import (
    "bufio"
    "bytes"
    "fmt"
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
    sam *goSam.Client
    err error

    transport *http.Transport
    http *http.Client
    host string

    sendPath string
    sendPipe *os.File
    sendBuff bufio.Reader

    namePath string
    namePipe *os.File
    name string
}

var connectionDirectory string

func (samConn *samHttp) initPipes(){
    pathConnectionExists, pathErr := exists(filepath.Join(connectionDirectory, samConn.host))
    fmt.Println("Directory Check", filepath.Join(connectionDirectory, samConn.host))
    samConn.checkErr(pathErr)
    if ! pathConnectionExists {
        fmt.Println("Creating a connection:", samConn.host)
        os.Mkdir(filepath.Join(connectionDirectory, samConn.host), 0755)
    }

    samConn.sendPath = filepath.Join(connectionDirectory, samConn.host, "send")
    pathSendExists, sendPathErr := exists(samConn.sendPath)
    samConn.checkErr(sendPathErr)
    if ! pathSendExists {
        err := syscall.Mkfifo(samConn.sendPath, 0755)
        fmt.Println("Preparing to create Pipe:", samConn.sendPath)
        samConn.checkErr(err)
        fmt.Println("checking for problems...")
        samConn.sendPipe, err = os.OpenFile(samConn.sendPath , os.O_RDWR|os.O_CREATE, 0755)
        fmt.Println("Opening the Named Pipe as a File...")
        samConn.sendBuff = *bufio.NewReader(samConn.sendPipe)
        fmt.Println("Opening the Named Pipe as a Buffer...")
        fmt.Println("Created a named Pipe for sending requests:", samConn.sendPath)
    }

    samConn.namePath = filepath.Join(connectionDirectory, samConn.host, "name")
    pathNameExists, namePathErr := exists(samConn.namePath)
    samConn.checkErr(namePathErr)
    if ! pathNameExists {
        err := syscall.Mkfifo(samConn.namePath, 0755)
        fmt.Println("Preparing to create Pipe:", samConn.namePath)
        samConn.checkErr(err)
        fmt.Println("checking for problems...")
        samConn.namePipe, err = os.OpenFile(samConn.namePath , os.O_RDWR|os.O_CREATE, 0755)
        fmt.Println("Created a named Pipe for the full name:", samConn.namePath)
    }

}


func (samConn *samHttp) createClient(samAddr string, samPort string, request string) {
    samCombined := samAddr + ":" + samPort
    samConn.sam, samConn.err = goSam.NewClient(samCombined)
    samConn.checkErr(samConn.err)
    fmt.Println("Established SAM connection")
    samConn.transport = &http.Transport{
		Dial: samConn.sam.Dial,
	}
    samConn.http = &http.Client{Transport: samConn.transport}
    if samConn.host == "" {
        samConn.host, _ = samConn.hostSet(request)
        samConn.initPipes()
    }
    samConn.subCache = append(samConn.subCache, newSamUrl(samConn.host))
    samConn.writeName(request)
}

func (samConn *samHttp) hostSet(request string) (string, string){
    tmp := strings.Replace(request, "http://", "", -1)
    host := strings.SplitAfterN(tmp, ".i2p", 1 )[0]
    _, err := url.ParseRequestURI("http://" + host)
    if err != nil {
        host = strings.Replace(host, "http://", "", -1)
    }
    directory := strings.Replace(tmp, host, "", -1)
    fmt.Println("Setting up micro-proxy for:", "http://" + host)
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
            fmt.Println("Request host ", comphost)
            fmt.Println("Is equal to client host", samConn.host)
            return true
        }else{
            fmt.Println("Request host ", comphost)
            fmt.Println("Is not equal to client host", samConn.host)
            return false
        }
    }else{
        host = strings.Replace(host, "http://", "", -1)
        host = strings.SplitAfterN(request, ".i2p", -1 )[0]
        host = strings.Replace(host, "http://", "", -1)
        if samConn.host == host {
            fmt.Println("Request host ", host)
            fmt.Println("Is equal to client host", samConn.host)
            return true
        }else{
            fmt.Println("Request host ", host)
            fmt.Println("Is not equal to client host", samConn.host)
            return false
        }
    }
}

func (samConn *samHttp) getRequest(request string) (string, string){
    host := request
    //tmp := strings.SplitAfterN(request, ".i2p", -1)
    directory := strings.Replace(request, "http://", "", -1)
    _, err := url.ParseRequestURI(host)
    if err != nil {
        host = "http://" + request
        fmt.Println("URL failed validation, correcting to:", host)
    }else{
        fmt.Println("URL passed validation:", request)
    }
    return host, directory
}

func (samConn *samHttp) sendRequest(request string) int{
    r, dir := samConn.getRequest(request)
    resp, err := samConn.http.Get(r)
    samConn.checkErr(err)
    samConn.copyRequest(resp, dir)
    return 0
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
        samConn.subCache = append(samConn.subCache, newSamUrl(directory))
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
    fmt.Println("Turning string %s into a response", input)
    return tmp
}

func (samConn *samHttp) printResponse() string{
    s, n := samConn.scannerText()
    if n == 0 {
        fmt.Println("Maintaining Connection:", samConn.hostGet())
        return ""
    }else if n < 0 {
        fmt.Println("Something wierd happened with :" , s)
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
    fmt.Println("Reading n bytes from send pipe:", strconv.Itoa(n))
    if n == 0 {
        fmt.Println("Maintaining Connection:", samConn.hostGet())
    }else if n < 0 {
        fmt.Println("Something wierd happened with :", line)
        fmt.Println("end determined at index :", strconv.Itoa(n))
    }else{
        s := string( line[:n] )
        fmt.Println("Sending request:", s)
        samConn.sendRequest(s)
    }
}

func (samConn *samHttp) readDelete() bool {
    b := false
    for _, dir := range samConn.subCache {
        n := dir.readDelete()
        if n == 0 {
            fmt.Println("Maintaining Connection:", samConn.hostGet())
        }else if n > 0 {
            b = true
        }
    }
    return b
}


func (samConn *samHttp) writeName(request string){
    if samConn.host == "" {
        fmt.Println("Setting hostname:" )
        //directory := ""
        samConn.host, _ = samConn.hostSet(request)
        samConn.initPipes()
    }
    fmt.Println("Attempting to write-out connection name:")
    if samConn.checkName() {
        samConn.name, samConn.err = samConn.sam.Lookup(samConn.host)
        fmt.Println("New Connection Name: %s", samConn.name)
        samConn.checkErr(samConn.err)
        io.Copy(samConn.namePipe, bytes.NewBufferString(samConn.name))
    }
}

func (samConn *samHttp) checkName() bool{
    fmt.Println("seeing if the connection needs a name:")
    if samConn.name == "" {
        fmt.Println("Naming connection:")
        return true
    }else{
        return false
    }
}

func (samConn *samHttp) cleanupClient(){
    samConn.sendPipe.Close()
    samConn.namePipe.Close()
    for _, url := range samConn.subCache {
        url.cleanupDirectory()
    }
    samConn.sam.Close()
    os.RemoveAll(filepath.Join(connectionDirectory, samConn.host))
}

func (samConn *samHttp) checkErr(err error) {
	if err != nil {
        samConn.cleanupClient()
		log.Fatal(err)
	}
}

func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return true, err
}

func newSamHttp(samAddrString string, samPortString string, request string) (samHttp){
    fmt.Println("Creating a new SAMv3 Client.")
    var samConn samHttp
    samConn.name = ""
    samConn.host = ""
    samConn.createClient(samAddrString, samPortString, request)
    return samConn
}

