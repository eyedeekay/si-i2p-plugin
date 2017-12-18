package main

import (
    "bufio"
    "path/filepath"
    "fmt"
    "io"
    "log"
    "os"
    "strconv"
    "strings"
    "syscall"
)

type samList struct{
    stackOfSams []samHttp
    err error
    up bool

    samAddrString string
    samPortString string

    sendPath string
    sendPipe *os.File
    sendBuff bufio.Reader
    sendScan bufio.Scanner

    recvPath string
    recvPipe *os.File

    delPath string
    delPipe *os.File
    //delBuff bufio.Reader
    delBuff bufio.Scanner
}

func (samStack * samList) initPipes(){
    pathConnectionExists, err := exists(filepath.Join(connectionDirectory, "parent"))
    samStack.checkErr(err)
    if ! pathConnectionExists {
        fmt.Println("Creating a connection:", "parent")
        os.Mkdir(filepath.Join(connectionDirectory, "parent"), 0755)
    }else{
        os.RemoveAll(filepath.Join(connectionDirectory, "parent"))
        fmt.Println("Creating a connection:", "parent")
        os.Mkdir(filepath.Join(connectionDirectory, "parent"), 0755)
    }

    samStack.sendPath = filepath.Join(connectionDirectory, "parent", "send")
    pathSendExists, sendErr := exists(samStack.sendPath)
    samStack.checkErr(sendErr)
    if ! pathSendExists {
        samStack.err = syscall.Mkfifo(samStack.sendPath, 0755)
        fmt.Println("Preparing to create Pipe:", samStack.sendPath)
        samStack.checkErr(samStack.err)
        fmt.Println("checking for problems...")
        samStack.sendPipe, samStack.err = os.OpenFile(samStack.sendPath , os.O_RDWR|os.O_CREATE, 0755)
       fmt.Println("Opening the Named Pipe as a File...")
        samStack.sendScan = *bufio.NewScanner(samStack.sendPipe)
        samStack.sendBuff = *bufio.NewReader(samStack.sendPipe)
        fmt.Println("Opening the Named Pipe as a Scanner...")
        fmt.Println("Created a named Pipe for sending requests:", samStack.sendPath)
    }

    samStack.recvPath = filepath.Join(connectionDirectory, "parent", "recv")
    pathRecvExists, recvErr := exists(samStack.recvPath)
    samStack.checkErr(recvErr)
    if ! pathRecvExists {
        samStack.err = syscall.Mkfifo(samStack.recvPath, 0755)
        fmt.Println("Preparing to create Pipe:", samStack.recvPath)
        samStack.checkErr(samStack.err)
        fmt.Println("checking for problems...")
        samStack.recvPipe, samStack.err = os.OpenFile(samStack.recvPath , os.O_RDWR|os.O_CREATE, 0755)
        samStack.recvPipe.WriteString("")
        fmt.Println("Created a named Pipe for recieving responses:", samStack.recvPath)
    }

    samStack.delPath = filepath.Join(connectionDirectory, "parent", "del")
    pathDelExists, delErr := exists(samStack.delPath)
    samStack.checkErr(delErr)
    if ! pathDelExists{
        samStack.err = syscall.Mkfifo(samStack.delPath, 0755)
        fmt.Println("Preparing to create Pipe:", samStack.delPath)
        samStack.checkErr(samStack.err)
        fmt.Println("checking for problems...")
        samStack.delPipe, samStack.err = os.OpenFile(samStack.delPath , os.O_RDWR|os.O_CREATE, 0755)
        samStack.recvPipe.WriteString("")
        fmt.Println("Opening the Named Pipe as a File...")
        samStack.delBuff = *bufio.NewScanner(samStack.delPipe)
        fmt.Println("Opening the Named Pipe as a Scanner...")
        fmt.Println("Created a named Pipe for closing the connection:", samStack.delPath)
    }
    samStack.up = true;
}

func (samStack *samList) createClient(request string){
    fmt.Println("Requesting a new SAM-based http client")
    samStack.stackOfSams = append(samStack.stackOfSams, newSamHttp(samStack.samAddrString, samStack.samPortString, request))
}

func (samStack *samList) createSamList(samAddrString string, samPortString string){
    samStack.samAddrString = samAddrString
    samStack.samPortString = samPortString
    if ! samStack.up {
        samStack.initPipes()
        fmt.Println("Parent proxy pipes initialized. Parent proxy set to up.")
    }
}

func (samStack *samList) sendClientRequest(request string){
    found := false
    for index, client := range samStack.stackOfSams {
        fmt.Println("Checking client requests", index + 1)
        fmt.Println("of", len(samStack.stackOfSams))
        if client.hostCheck(request){
            fmt.Println("Client pipework for %s found.", request)
            go client.sendRequest(request)
            found = true
        }
    }
    if ! found {
        fmt.Println("Client pipework for %s not found: Creating.", request)
        samStack.createClient(request)
        for index := len(samStack.stackOfSams)-1; index >= 0 ; index-- {
            fmt.Println("Checking client requests", index + 1)
            fmt.Println("of", len(samStack.stackOfSams))
            if samStack.stackOfSams[index].hostCheck(request){
                fmt.Println("Client pipework for %s found.", request)
                go samStack.stackOfSams[index].sendRequest(request)
                found = true
            }
        }
    }
}

func (samStack *samList) sendText() (string, int) {
    s := ""
    //for samStack.sendBuff.Scan() {
    samStack.sendScan.Scan()
    s += samStack.sendScan.Text()
    //}
    fmt.Println(s)
    if s != "" {
        return s, len(s)
    }else{
        return "", 0
    }
}

func (samStack *samList) readRequest() string{
    //s, n := samStack.sendText()
    line, _, err := samStack.sendBuff.ReadLine()
    samStack.checkErr(err)
    n := len(line)
    //fmt.Println("Reading n bytes from Parent send pipe:", strconv.Itoa(n))
    if n == 0 {
        fmt.Println("Flush the pipe maybe?:")
    }else if n < 0 {
        fmt.Println("Something wierd happened with the Parent Send pipe." )
        fmt.Println("end determined at index :", strconv.Itoa(n))
    }else{
        s := string( line[:n] )
        fmt.Println("Sending request:", s)
        samStack.sendClientRequest(s)
        return s
    }
    return ""
}

func (samStack *samList) writeResponses(){
    for i, client := range samStack.stackOfSams {
        fmt.Println("Checking for responses: %s", i+1)
        fmt.Println("of: ", len(samStack.stackOfSams))
        samStack.writeRecieved(client.printResponse())
    }
}

func (samStack *samList) responsify(input string) io.Reader {
    tmp := strings.NewReader(input)
    fmt.Println("Responsifying string: ")
    return tmp
}

func (samStack *samList) writeRecieved(response string){
    fmt.Println("Response test", response)
    if response != "" {
        fmt.Println("Got response: %s", response )
        io.Copy(samStack.recvPipe, samStack.responsify(response))
    }
}

func (samStack *samList) delText() (string, int) {
    s := ""
    samStack.delBuff.Scan()
    s += samStack.delBuff.Text()
    fmt.Println(s)
    if s != "" {
        return s, len(s)
    }else{
        return "", 0
    }
}

func (samStack *samList) readDelete() bool {
    s, n := samStack.delText()
    fmt.Println("Reading n bytes from exit pipe:", strconv.Itoa(n))
    if n == 0 {
        return false
    }else if n < 0 {
        fmt.Println("Something wierd happened with :", s)
        fmt.Println("end determined at index :", strconv.Itoa(n))
        return false
    }else{
        if s == "y" {
            fmt.Println("Closing proxy.")
            defer samStack.cleanupClient()
            return true
        }else{
            return false
        }
    }
}

func (samStack *samList) cleanupClient(){
    samStack.sendPipe.Close()
    samStack.recvPipe.Close()
    for _, client := range samStack.stackOfSams {
        client.cleanupClient()
    }
    samStack.delPipe.Close()
    os.RemoveAll(filepath.Join(connectionDirectory, "parent"))
}

func (samStack *samList) checkErr(err error) {
	if err != nil {
        samStack.cleanupClient()
		log.Fatal(err)
	}
}

func createSamList(samAddr string, samPort string, initAddress string) samList{
    var samStack samList
    fmt.Println("Generating parent proxy structure.")
    samStack.up = false
    fmt.Println("Parent proxy set to down.")
    samStack.createSamList(samAddr, samPort)
    fmt.Println("SAM list created")
    samStack.sendPipe.WriteString(initAddress + "\n")
    return samStack
}
