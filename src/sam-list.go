package main

import (
        "bufio"
        "path/filepath"
        "fmt"
        "io"
        "log"
        "os"
        "strconv"
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

        recvPath string
        recvPipe *os.File

        delPath string
        delPipe *os.File
        delBuff bufio.Reader
}

func (samStack * samList) initPipes(){
        pathConnectionExists, err := exists(filepath.Join(connectionDirectory, "parent"))
        samStack.checkErr(err)
        if ! pathConnectionExists {
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
                samStack.sendBuff = *bufio.NewReader(samStack.sendPipe)
                fmt.Println("Opening the Named Pipe as a Buffer...")
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
                fmt.Println("Opening the Named Pipe as a File...")
                samStack.delBuff = *bufio.NewReader(samStack.delPipe)
                fmt.Println("Opening the Named Pipe as a Buffer...")
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
                fmt.Println("Checking client", index)
                fmt.Println("of", len(samStack.stackOfSams))
                if client.hostCheck(request){
                        fmt.Println("Client pipework for %s found.", request)
                        client.sendRequest(request)
                        found = true
                }
        }
        if ! found {
                fmt.Println("Client pipework for %s not found: Creating.", request)
                samStack.createClient(request)
                for index := len(samStack.stackOfSams)-1; index >= 0 ; index-- {
                        fmt.Println("Checking client", index)
                        fmt.Println("of", len(samStack.stackOfSams))
                        if samStack.stackOfSams[index].hostCheck(request){
                                fmt.Println("Client pipework for %s found.", request)
                                samStack.stackOfSams[index].sendRequest(request)
                                found = true
                        }
                }
        }
}

func (samStack *samList) readRequest() string{
        line, _, err := samStack.sendBuff.ReadLine()
        samStack.checkErr(err)
        n := len(line)
        fmt.Println("Reading n bytes from Parent send pipe:", strconv.Itoa(n))
        if n == 0 {
                fmt.Println("Flush the pipe maybe?:")
        }else if n < 0 {
                fmt.Println("Something wierd happened with :", line)
                fmt.Println("end determined at index :", strconv.Itoa(n))
        }else{
                s := string( line[:n] )
                fmt.Println("Sending request:", s)
                samStack.sendClientRequest(s)
                return s
        }
        return ""
}

func (samStack *samList) httpResponse(request string) (io.Reader, error){
        found := false
        var response io.Reader
        var err error
        for _, client := range samStack.stackOfSams {
                if client.hostCheck(request){
                        response, err = os.OpenFile(client.recvPath, os.O_RDWR|os.O_CREATE, 0755)
                        samStack.checkErr(err)
                        found = true
                }
        }
        if ! found {
                fmt.Println("Child proxy not found for: ", request)
        }
        return response, err
}

func (samStack *samList) writeResponse(request string){
        if request != "" {
                resp, err := samStack.httpResponse(request)
                samStack.checkErr(err)
                io.Copy(samStack.recvPipe, resp)
        }
}

func (samStack *samList) readDelete() bool {
        line, _, err := samStack.delBuff.ReadLine()
        samStack.checkErr(err)
        n := len(line)
        fmt.Println("Reading n bytes from exit pipe:", strconv.Itoa(n))
        if n == 0 {
                fmt.Println("Maintaining Connection.")
                return false
        }else if n < 0 {
                fmt.Println("Something wierd happened with :", line)
                fmt.Println("end determined at index :", strconv.Itoa(n))
                return false
        }else{
                s := string( line[:n] )
                if s == "y" {
                        fmt.Println("Closing proxy.")
                        defer samStack.cleanupClient()
                        return true
                }else{
                        return false
                }
        }
}

func (samStack *samList) clientLoop(){
        for true {
                s := samStack.readRequest()
                samStack.writeResponse(s)
        }
}

func (samStack *samList) cleanupClient(){
        samStack.sendPipe.Close()
        samStack.recvPipe.Close()
        for _, client := range samStack.stackOfSams {
                client.cleanupClient()
        }
        os.RemoveAll(filepath.Join(connectionDirectory))
}

func (samStack *samList) checkErr(err error) {
	if err != nil {
                samStack.cleanupClient()
		log.Fatal(err)
	}
}

func createSamList(samAddr string, samPort string) samList{
        var samStack samList
        fmt.Println("Generating parent proxy structure.")
        samStack.up = false
        fmt.Println("Parent proxy set to down.")
        samStack.createSamList(samAddr, samPort)
        fmt.Println("SAM list created")
        return samStack
}
