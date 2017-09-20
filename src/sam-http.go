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
        sam *goSam.Client
        err error

        transport *http.Transport
        http *http.Client
        host string

        sendPath string
        sendPipe *os.File
        sendBuff bufio.Reader

        recvPath string
        recvPipe *os.File

        namePath string
        namePipe *os.File
        name string

        delPath string
        delPipe *os.File
        delBuff bufio.Reader
}

var connectionDirectory string

func (samConn *samHttp) initPipes(){
        pathConnectionExists, err := exists(filepath.Join(connectionDirectory, samConn.host))
        samConn.checkErr(err)
        if ! pathConnectionExists {
                fmt.Println("Creating a connection:", samConn.host)
                os.Mkdir(filepath.Join(connectionDirectory, samConn.host), 0755)
        }

        samConn.sendPath = filepath.Join(connectionDirectory, samConn.host, "send")
        pathSendExists, sendErr := exists(samConn.sendPath)
        samConn.checkErr(sendErr)
        if ! pathSendExists {
                samConn.err = syscall.Mkfifo(samConn.sendPath, 0755)
                fmt.Println("Preparing to create Pipe:", samConn.sendPath)
                samConn.checkErr(samConn.err)
                fmt.Println("checking for problems...")
                samConn.sendPipe, samConn.err = os.OpenFile(samConn.sendPath , os.O_RDWR|os.O_CREATE, 0755)
                fmt.Println("Opening the Named Pipe as a File...")
                samConn.sendBuff = *bufio.NewReader(samConn.sendPipe)
                fmt.Println("Opening the Named Pipe as a Buffer...")
                fmt.Println("Created a named Pipe for sending requests:", samConn.sendPath)
        }

        samConn.recvPath = filepath.Join(connectionDirectory, samConn.host, "recv")
        pathRecvExists, recvErr := exists(samConn.recvPath)
        samConn.checkErr(recvErr)
        if ! pathRecvExists {
                samConn.err = syscall.Mkfifo(samConn.recvPath, 0755)
                fmt.Println("Preparing to create Pipe:", samConn.recvPath)
                samConn.checkErr(samConn.err)
                fmt.Println("checking for problems...")
                samConn.recvPipe, samConn.err = os.OpenFile(samConn.recvPath , os.O_RDWR|os.O_CREATE, 0755)
                fmt.Println("Created a named Pipe for recieving responses:", samConn.recvPath)
        }

        samConn.namePath = filepath.Join(connectionDirectory, samConn.host, "name")
        pathNameExists, nameErr := exists(samConn.namePath)
        samConn.checkErr(nameErr)
        if ! pathNameExists {
                samConn.err = syscall.Mkfifo(samConn.namePath, 0755)
                fmt.Println("Preparing to create Pipe:", samConn.namePath)
                samConn.checkErr(samConn.err)
                fmt.Println("checking for problems...")
                samConn.namePipe, samConn.err = os.OpenFile(samConn.namePath , os.O_RDWR|os.O_CREATE, 0755)
                fmt.Println("Created a named Pipe for the full name:", samConn.namePath)
        }

        samConn.delPath = filepath.Join(connectionDirectory, samConn.host, "del")
        pathDelExists, delErr := exists(samConn.delPath)
        samConn.checkErr(delErr)
        if ! pathDelExists{
                samConn.err = syscall.Mkfifo(samConn.delPath, 0755)
                fmt.Println("Preparing to create Pipe:", samConn.delPath)
                samConn.checkErr(samConn.err)
                fmt.Println("checking for problems...")
                samConn.delPipe, samConn.err = os.OpenFile(samConn.delPath , os.O_RDWR|os.O_CREATE, 0755)
                fmt.Println("Opening the Named Pipe as a File...")
                samConn.delBuff = *bufio.NewReader(samConn.delPipe)
                fmt.Println("Opening the Named Pipe as a Buffer...")
                fmt.Println("Created a named Pipe for closing the connection:", samConn.delPath)
        }
}


func (samConn *samHttp) createClient(samAddr string, samPort string, request string) {
        fmt.Println("Creating a new SAMv3 Client.")
        samCombined := samAddr + ":" + samPort
        samConn.sam, samConn.err = goSam.NewClient(samCombined)
        samConn.checkErr(samConn.err)
        fmt.Println("Established SAM connection")
        samConn.transport = &http.Transport{
		Dial: samConn.sam.Dial,
	}
        samConn.http = &http.Client{Transport: samConn.transport}
        samConn.host = ""
        if samConn.host == "" {
                samConn.host = samConn.hostSet(request)
                samConn.initPipes()
        }
        samConn.writeName(request)
}

func (samConn *samHttp) hostSet(request string) string{
        host := strings.SplitAfterN(request, ".i2p", 1 )[0]
        _, err := url.ParseRequestURI(host)
        if err != nil {
                host = strings.Replace(host, "http://", "", -1)
        }
        fmt.Println("Setting up micro-proxy for:", "http://" + host)
        return host
}

func (samConn *samHttp) hostGet() string{
        return "http://" + samConn.host
}

func (samConn *samHttp) hostCheck(request string) bool{
        host := strings.SplitAfterN(request, ".i2p", 1 )[0]
        _, err := url.ParseRequestURI(host)
        if err != nil {
                host = strings.Replace(host, "http://", "", -1)
        }
        if samConn.host == host {
                return true
        }else{
                return false
        }
}

func (samConn *samHttp) getRequest(request string) string{
        host := request
        _, err := url.ParseRequestURI(host)
        if err != nil {
                host = strings.Replace(request, "http://", "", -1)
        }
        return host
}

func (samConn *samHttp) sendRequest(request string) int{
        resp, err := samConn.http.Get(samConn.getRequest(request))
        samConn.checkErr(err)
        defer resp.Body.Close()
        io.Copy(samConn.recvPipe, resp.Body)
        return 0
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
        line, _, err := samConn.delBuff.ReadLine()
        samConn.checkErr(err)
        n := len(line)
        fmt.Println("Reading n bytes from exit pipe:", strconv.Itoa(n))
        if n == 0 {
                fmt.Println("Maintaining Connection:", samConn.hostGet())
                return false
        }else if n < 0 {
                fmt.Println("Something wierd happened with :", line)
                fmt.Println("end determined at index :", strconv.Itoa(n))
                return false
        }else{
                s := string( line[:n] )
                if s == "y" {
                        fmt.Println("Deleting connection: %s", samConn.host )
                        defer samConn.cleanupClient()
                        return true
                }else{
                        return false
                }
        }
}


func (samConn *samHttp) writeName(request string){
        if samConn.host == "" {
                fmt.Println("Setting hostname:" )
                samConn.host = samConn.hostSet(request)
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
        samConn.recvPipe.Close()
        samConn.namePipe.Close()
        samConn.delPipe.Close()
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
        samConn.createClient(samAddrString, samPortString, request)
        return samConn
}

