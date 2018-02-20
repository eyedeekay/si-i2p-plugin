package main

import (
    "bufio"
    "path/filepath"
    "log"
    "net/http"
    "io"
    "os"
    "strings"
    "syscall"

    "github.com/eyedeekay/gosam"
)

type samList struct{
    stackOfSams []samHttp
    sam *goSam.Client
    err error
    up bool

    samAddrString string
    samPortString string

    sendPath string
    sendPipe *os.File
    sendScan bufio.Scanner

    recvPath string
    recvPipe *os.File

    delPath string
    delPipe *os.File
    delScan bufio.Scanner
}

func (samStack * samList) initPipes(){
    pathConnectionExists, err := exists(filepath.Join(connectionDirectory, "parent"))
    samStack.checkErr(err)
    if ! pathConnectionExists {
        log.Println("Creating a connection:", "parent")
        os.Mkdir(filepath.Join(connectionDirectory, "parent"), 0755)
    }else{
        os.RemoveAll(filepath.Join(connectionDirectory, "parent"))
        log.Println("Creating a connection:", "parent")
        os.Mkdir(filepath.Join(connectionDirectory, "parent"), 0755)
    }

    samStack.sendPath = filepath.Join(connectionDirectory, "parent", "send")
    pathSendExists, sendErr := exists(samStack.sendPath)
    samStack.checkErr(sendErr)
    if ! pathSendExists {
        samStack.err = syscall.Mkfifo(samStack.sendPath, 0755)
        log.Println("Preparing to create Pipe:", samStack.sendPath)
        samStack.checkErr(samStack.err)
        log.Println("checking for problems...")
        samStack.sendPipe, samStack.err = os.OpenFile(samStack.sendPath , os.O_RDWR|os.O_CREATE, 0755)
        log.Println("Opening the Named Pipe as a Scanner...")
        samStack.sendScan = *bufio.NewScanner(samStack.sendPipe)
        samStack.sendScan.Split(bufio.ScanLines)
        log.Println("Opening the Named Pipe as a Scanner...")
        log.Println("Created a named Pipe for sending requests:", samStack.sendPath)
    }

    samStack.recvPath = filepath.Join(connectionDirectory, "parent", "recv")
    pathRecvExists, recvErr := exists(samStack.recvPath)
    samStack.checkErr(recvErr)
    if ! pathRecvExists {
        samStack.err = syscall.Mkfifo(samStack.recvPath, 0755)
        log.Println("Preparing to create Pipe:", samStack.recvPath)
        samStack.checkErr(samStack.err)
        log.Println("checking for problems...")
        samStack.recvPipe, samStack.err = os.OpenFile(samStack.recvPath , os.O_RDWR|os.O_CREATE, 0755)
        samStack.recvPipe.WriteString("")
        log.Println("Created a named Pipe for recieving responses:", samStack.recvPath)
    }

    samStack.delPath = filepath.Join(connectionDirectory, "parent", "del")
    pathDelExists, delErr := exists(samStack.delPath)
    samStack.checkErr(delErr)
    if ! pathDelExists{
        samStack.err = syscall.Mkfifo(samStack.delPath, 0755)
        log.Println("Preparing to create Pipe:", samStack.delPath)
        samStack.checkErr(samStack.err)
        log.Println("checking for problems...")
        samStack.delPipe, samStack.err = os.OpenFile(samStack.delPath , os.O_RDWR|os.O_CREATE, 0755)
        samStack.recvPipe.WriteString("")
        log.Println("Opening the Named Pipe as a File...")
        samStack.delScan = *bufio.NewScanner(samStack.delPipe)
        samStack.delScan.Split(bufio.ScanLines)
        log.Println("Opening the Named Pipe as a Scanner...")
        log.Println("Created a named Pipe for closing the connection:", samStack.delPath)
    }
    samStack.up = true;
}

func (samStack *samList) createClient(request string){
    log.Println("Appending client to SAM stack.")
    samStack.stackOfSams = append(samStack.stackOfSams, newSamHttp(samStack.samAddrString, samStack.samPortString, samStack.sam, request))
}

func (samStack *samList) createClientHttp(request *http.Request){
    log.Println("Appending client to SAM stack.")
    samStack.stackOfSams = append(samStack.stackOfSams, newSamHttpHttp(samStack.samAddrString, samStack.samPortString, samStack.sam, request))
}

func (samStack *samList) createSamList(samAddrString string, samPortString string){
    samStack.samAddrString = samAddrString
    samStack.samPortString = samPortString
    log.Println("Requesting a new SAM-based http client")
    samCombined := samStack.samAddrString + ":" + samStack.samPortString
    samStack.sam, samStack.err = goSam.NewClient(samCombined)
    samStack.checkErr(samStack.err)
    log.Println("Established SAM connection")
    if ! samStack.up {
        samStack.initPipes()
        log.Println("Parent proxy pipes initialized. Parent proxy set to up.")
    }
}

func (samStack *samList) sendClientRequest(request string) string{
    found := false
    for index, client := range samStack.stackOfSams {
        log.Println("Checking client requests", index + 1)
        log.Println("of", len(samStack.stackOfSams))
        if client.hostCheck(request){
            log.Println("Client pipework for %s found.", request)
            client.sendRequest(request)
            log.Println("Request sent")
            found = true
        }
    }
    if ! found {
        log.Println("Client pipework for %s not found: Creating.", request)
        samStack.createClient(request)
        for index, client := range samStack.stackOfSams {
            log.Println("Checking client requests", index + 1)
            log.Println("of", len(samStack.stackOfSams))
            if client.hostCheck(request){
                log.Println("Client pipework for %s found.", request)
                client.sendRequest(request)
                log.Println("Request sent")
                found = true
            }
        }
    }
    return request
}

func (samStack *samList) sendClientRequestHttp(request *http.Request) *http.Client {
    found := false
    log.Println("The SAM stack me exists?")
    for index, client := range samStack.stackOfSams {
        log.Println("Checking client requests", index + 1)
        log.Println("of", len(samStack.stackOfSams))
        if client.hostCheck(request.Host){
            log.Println("Client pipework for %s found.", request.Host)
            log.Println("URL scheme", request.URL.Scheme)
            found = true
            log.Println("Request sent")
            return client.sendRequestHttp(request)
        }
    }
    if ! found {
        log.Println("Client pipework for %s not found: Creating.", request.Host)
        samStack.createClientHttp(request)
        for index, client := range samStack.stackOfSams {
            log.Println("Checking client requests", index + 1)
            log.Println("of", len(samStack.stackOfSams))
            if client.hostCheckHttp(request){
                log.Println("Client pipework for %s found.", request.URL.String() )
                log.Println("URL scheme", request.URL.Scheme)
                found = true
                log.Println("Request sent")
                return client.sendRequestHttp(request)
            }
        }
    }
    return nil
}

func (samStack *samList) readRequest() string{
    for samStack.sendScan.Scan(){
        return samStack.sendClientRequest(samStack.sendScan.Text())
    }
    return ""
}

func (samStack *samList) writeResponses(){
    for i, client := range samStack.stackOfSams {
        log.Println("Checking for responses: %s", i+1)
        log.Println("of: ", len(samStack.stackOfSams))
        b := samStack.writeRecieved(client.printResponse())
        if b == true {
            break
        }
    }
}

func (samStack *samList) responsify(input string) io.Reader {
    tmp := strings.NewReader(input)
    log.Println("Responsifying string:")
    return tmp
}

func (samStack *samList) writeRecieved(response string) bool {
    b := false
    if response != "" {
        log.Println("Got response:")
        io.Copy(samStack.recvPipe, samStack.responsify(response))
        b = true
    }
    return b
}

func (samStack *samList) readDelete() bool {
    for samStack.delScan.Scan(){
        if samStack.delScan.Text() == "y" || samStack.delScan.Text() == "Y" {
            defer samStack.cleanupClient()
            return true
        }else{
            return false
        }
    }
    return false
}

func (samStack *samList) cleanupClient(){
    samStack.sendPipe.Close()
    samStack.recvPipe.Close()
    for _, client := range samStack.stackOfSams {
        client.cleanupClient()
    }
    samStack.delPipe.Close()
    err := samStack.sam.Close()
    samStack.checkErr(err)
    os.RemoveAll(filepath.Join(connectionDirectory, "parent"))
}

func (samStack *samList) checkErr(err error) {
	if err != nil {
        log.Fatal(err)
        samStack.cleanupClient()
	}
}

func createSamList(samAddr string, samPort string, initAddress string) samList{
    var samStack samList
    log.Println("Generating parent proxy structure.")
    samStack.up = false
    log.Println("Parent proxy set to down.")
    samStack.createSamList(samAddr, samPort)
    log.Println("SAM list created")
    samStack.sendPipe.WriteString(initAddress + "\n")
    return samStack
}
