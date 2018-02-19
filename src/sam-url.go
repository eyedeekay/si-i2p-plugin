package main

import (
    "bufio"
    //"bytes"
	"io/ioutil"
	"log"
	"net/http"
    "os"
    "path/filepath"
    //"strings"
    "strconv"
    "syscall"
    //"net/url"
    "time"

)

type samUrl struct{
    err error

    transport *http.Transport
    http *http.Client
    subdirectory string

    recvPath string
    recvFile *os.File

    timePath string
    timeFile *os.File

    delPath string
    delPipe *os.File
    delBuff bufio.Reader
}

func (subUrl *samUrl) initPipes(){
    pathConnectionExists, pathErr := exists(filepath.Join(connectionDirectory, subUrl.subdirectory))
    log.Println("Directory Check", filepath.Join(connectionDirectory, subUrl.subdirectory))
    subUrl.checkErr(pathErr)
    if ! pathConnectionExists {
        log.Println("Creating a connection:", subUrl.subdirectory)
        os.MkdirAll(filepath.Join(connectionDirectory, subUrl.subdirectory), 0755)
    }

    subUrl.recvPath = filepath.Join(connectionDirectory, subUrl.subdirectory, "recv")
    pathRecvExists, recvPathErr := exists(subUrl.recvPath)
    subUrl.checkErr(recvPathErr)
    if ! pathRecvExists {
        subUrl.recvFile, subUrl.err = os.Create(subUrl.recvPath)
        log.Println("Preparing to create File:", subUrl.recvPath)
        subUrl.checkErr(subUrl.err)
        log.Println("checking for problems...")
        log.Println("Opening the File...")
        subUrl.recvFile, subUrl.err = os.OpenFile(subUrl.recvPath, os.O_RDWR|os.O_CREATE, 0644)
        log.Println("Created a File for recieving responses:", subUrl.recvPath)
    }

    subUrl.timePath = filepath.Join(connectionDirectory, subUrl.subdirectory, "time")
    pathTimeExists, recvTimeErr := exists(subUrl.timePath)
    subUrl.checkErr(recvTimeErr)
    if ! pathTimeExists {
        subUrl.timeFile, subUrl.err = os.Create(subUrl.timePath)
        log.Println("Preparing to create File:", subUrl.timePath)
        subUrl.checkErr(subUrl.err)
        log.Println("checking for problems...")
        log.Println("Opening the File...")
        subUrl.timeFile, subUrl.err = os.OpenFile(subUrl.timePath, os.O_RDWR|os.O_CREATE, 0644)
        log.Println("Created a File for timing responses:", subUrl.timePath)
    }

    subUrl.delPath = filepath.Join(connectionDirectory, subUrl.subdirectory, "del")
    pathDelExists, delPathErr := exists(subUrl.delPath)
    subUrl.checkErr(delPathErr)
    if ! pathDelExists{
        err := syscall.Mkfifo(subUrl.delPath, 0755)
        log.Println("Preparing to create Pipe:", subUrl.delPath)
        subUrl.checkErr(err)
        log.Println("checking for problems...")
        subUrl.delPipe, err = os.OpenFile(subUrl.delPath , os.O_RDWR|os.O_CREATE, 0755)
        log.Println("Opening the Named Pipe as a File...")
        subUrl.delBuff = *bufio.NewReader(subUrl.delPipe)
        log.Println("Opening the Named Pipe as a Buffer...")
        log.Println("Created a named Pipe for closing the connection:", subUrl.delPath)
    }
}

func (subUrl *samUrl) createDirectory(requestdir string) {
    subUrl.http = &http.Client{Transport: subUrl.transport}
    subUrl.subdirectory = subUrl.dirSet(requestdir)
    subUrl.initPipes()
}

func (subUrl *samUrl) scannerText() (string, int) {
    d, _ := ioutil.ReadFile(subUrl.recvPath)
    s := string(d)
    if s != "" {
        return s, len(s)
    }else{
       return "", 0
    }
}

func (subUrl *samUrl) dirSet(requestdir string) string {
    log.Println("Requesting directory: ", requestdir)
    return requestdir
}

func (subUrl *samUrl) copyDirectory(response *http.Response, directory string) bool{
    b := false
    if directory == subUrl.subdirectory {
        if response.StatusCode == http.StatusOK {
            subUrl.dealResponse(response)
        }
        b = true
    }
    return b
}

func (subUrl *samUrl) dealResponse(response *http.Response){
    defer response.Body.Close()
    body, _ := ioutil.ReadAll(response.Body)
    subUrl.recvFile.Write(body)
    subUrl.timeFile.WriteString(time.Now().String())
}

func (subUrl *samUrl) cleanupDirectory(){
    subUrl.recvFile.Close()
    subUrl.timeFile.Close()
    subUrl.delPipe.Close()
    os.RemoveAll(filepath.Join(connectionDirectory, subUrl.subdirectory))
}

func (subUrl *samUrl) readDelete() int {
    line, _, err := subUrl.delBuff.ReadLine()
    subUrl.checkErr(err)
    n := len(line)
    log.Println("Reading n bytes from exit pipe:", strconv.Itoa(n))
    if n < 0 {
        log.Println("Something wierd happened with :", line)
        log.Println("end determined at index :", strconv.Itoa(n))
        return n
    }else{
        s := string( line[:n] )
        if s == "y" {
            log.Println("Deleting connection: %s", subUrl.subdirectory )
            defer subUrl.cleanupDirectory()
            return n
        }else{
            return n
        }
    }
}

func (subUrl *samUrl) checkErr(err error) {
	if err != nil {
        subUrl.cleanupDirectory()
		log.Fatal(err)
	}
}

func newSamUrl(requestdir string) (samUrl){
    log.Println("Creating a new cache directory.")
    var subUrl samUrl
    subUrl.subdirectory = requestdir
    subUrl.createDirectory(requestdir)
    return subUrl
}

func newSamUrlHttp(request *http.Request) (samUrl){
    log.Println("Creating a new cache directory.")
    var subUrl samUrl
    subUrl.subdirectory = request.Host + request.URL.Path
    log.Println(subUrl.subdirectory)
    subUrl.createDirectory(subUrl.subdirectory)
    return subUrl
}
