package main

import (
    "bufio"
    "bytes"
	"io/ioutil"
	"log"
	"net/http"
    "os"
    "path/filepath"
    "strings"
    "strconv"
    "syscall"
    "time"
)

type samUrl struct{
    err error
    subDirectory string

    recvPath string
    recvFile *os.File

    timePath string
    timeFile *os.File

    delPath string
    delPipe *os.File
    delBuff bufio.Reader
}

func (subUrl *samUrl) initPipes(){
    pathConnectionExists, pathErr := exists(filepath.Join(connectionDirectory, subUrl.subDirectory))
    log.Println("Directory Check", filepath.Join(connectionDirectory, subUrl.subDirectory))
    subUrl.Fatal(pathErr)
    if ! pathConnectionExists {
        log.Println("Creating a connection:", subUrl.subDirectory)
        os.MkdirAll(filepath.Join(connectionDirectory, subUrl.subDirectory), 0755)
    }

    subUrl.recvPath = filepath.Join(connectionDirectory, subUrl.subDirectory, "recv")
    pathRecvExists, recvPathErr := exists(subUrl.recvPath)
    subUrl.Fatal(recvPathErr)
    if ! pathRecvExists {
        subUrl.recvFile, subUrl.err = os.Create(subUrl.recvPath)
        log.Println("Preparing to create File:", subUrl.recvPath)
        subUrl.Fatal(subUrl.err)
        log.Println("checking for problems...")
        log.Println("Opening the File...")
        subUrl.recvFile, subUrl.err = os.OpenFile(subUrl.recvPath, os.O_RDWR|os.O_CREATE, 0644)
        log.Println("Created a File for recieving responses:", subUrl.recvPath)
    }

    subUrl.timePath = filepath.Join(connectionDirectory, subUrl.subDirectory, "time")
    pathTimeExists, recvTimeErr := exists(subUrl.timePath)
    subUrl.Fatal(recvTimeErr)
    if ! pathTimeExists {
        subUrl.timeFile, subUrl.err = os.Create(subUrl.timePath)
        log.Println("Preparing to create File:", subUrl.timePath)
        subUrl.Fatal(subUrl.err)
        log.Println("checking for problems...")
        log.Println("Opening the File...")
        subUrl.timeFile, subUrl.err = os.OpenFile(subUrl.timePath, os.O_RDWR|os.O_CREATE, 0644)
        log.Println("Created a File for timing responses:", subUrl.timePath)
    }

    subUrl.delPath = filepath.Join(connectionDirectory, subUrl.subDirectory, "del")
    pathDelExists, delPathErr := exists(subUrl.delPath)
    subUrl.Fatal(delPathErr)
    if ! pathDelExists{
        err := syscall.Mkfifo(subUrl.delPath, 0755)
        log.Println("Preparing to create Pipe:", subUrl.delPath)
        subUrl.Fatal(err)
        log.Println("checking for problems...")
        subUrl.delPipe, err = os.OpenFile(subUrl.delPath , os.O_RDWR|os.O_CREATE, 0755)
        log.Println("Opening the Named Pipe as a File...")
        subUrl.delBuff = *bufio.NewReader(subUrl.delPipe)
        log.Println("Opening the Named Pipe as a Buffer...")
        log.Println("Created a named Pipe for closing the connection:", subUrl.delPath)
    }
}

func (subUrl *samUrl) createDirectory(requestdir string) {
    subUrl.subDirectory = subUrl.dirSet(requestdir)
    subUrl.initPipes()
}

func (subUrl *samUrl) scannerText() (string, error) {
    d, err := ioutil.ReadFile(subUrl.recvPath)
    subUrl.Fatal(err)
    s := string(d)
    if s != "" {
        log.Println("Read file", s)
        return s, err
    }
    return "", err
}

func (subUrl *samUrl) dirSet(requestdir string) string {
    log.Println("Requesting directory: ", requestdir + "/")
    d1 := requestdir + "/"
    d2 := strings.Replace(d1, "//", "/", -1)
    return d2
}

func (subUrl *samUrl) copyDirectory(response *http.Response, directory string) bool{
    b := false
    d1 := strings.Replace(subUrl.subDirectory, "/", "", -1)
    d2 := strings.Replace(directory, "/", "", -1)
    if d2 == d1 {
        log.Println("Directory / ", directory + " : compare : " + subUrl.subDirectory )
        if response != nil {
            log.Println("Response Status ", response.StatusCode)
            if response.StatusCode == http.StatusOK {
                log.Println("Setting file in cache")
                subUrl.dealResponse(response)
            }
        }
        b = true
    }
    return b
}

func (subUrl *samUrl) copyDirectoryHttp(request *http.Request, response *http.Response, directory string) (bool, *http.Response){
    b := false
    d1 := strings.Replace(subUrl.subDirectory, "/", "", -1)
    d2 := strings.Replace(directory, "/", "", -1)
    if d2 == d1 {
        log.Println("Directory / ", directory + " : compare : " + subUrl.subDirectory )
        if response != nil {
            log.Println("Response Status ", response.StatusCode)
            if response.StatusCode == http.StatusOK {
                log.Println("Setting file in cache")
                b = true
                resp := subUrl.dealResponseHttp(request ,response)
                return b, resp
            }
        }
        b = true
    }
    return b, response
}

func (subUrl *samUrl) dealResponse(response *http.Response){
    defer response.Body.Close()
    body, err := ioutil.ReadAll(response.Body)
    subUrl.Fatal(err)
    log.Println("Writing files.")
    subUrl.recvFile.Write(body)
    log.Println("Retrieval time: ", time.Now().String())
    subUrl.timeFile.WriteString(time.Now().String())
}

func (subUrl *samUrl) dealResponseHttp(request *http.Request, response *http.Response)(*http.Response){
    defer response.Body.Close()
    body, err := ioutil.ReadAll(response.Body)
    subUrl.Fatal(err)
    log.Println("Writing files.")
    subUrl.recvFile.Write(body)
    r := &http.Response{
        Status:        "200 OK",
        StatusCode:    200,
        Proto:         "HTTP/1.1",
        ProtoMajor:    1,
        ProtoMinor:    1,
        Body:          ioutil.NopCloser(bytes.NewBuffer(body)),
        ContentLength: int64(len(body)),
        Request:       request,
        Header:        make(http.Header, 0),
    }
    log.Println("Retrieval time: ", time.Now().String())
    subUrl.timeFile.WriteString(time.Now().String())
    return r
}

func (subUrl *samUrl) cleanupDirectory(){
    subUrl.recvFile.Close()
    subUrl.timeFile.Close()
    subUrl.delPipe.Close()
    os.RemoveAll(filepath.Join(connectionDirectory, subUrl.subDirectory))
}

func (subUrl *samUrl) readDelete() int {
    line, _, err := subUrl.delBuff.ReadLine()
    subUrl.Fatal(err)
    n := len(line)
    log.Println("Reading n bytes from exit pipe:", strconv.Itoa(n))
    if n < 0 {
        log.Println("Something wierd happened with :", line)
        log.Println("end determined at index :", strconv.Itoa(n))
        return n
    }else{
        s := string( line[:n] )
        if s == "y" {
            log.Println("Deleting connection: %s", subUrl.subDirectory )
            defer subUrl.cleanupDirectory()
            return n
        }else{
            return n
        }
    }
}

func (subUrl *samUrl) Warn(err error) {
	if err != nil {
		log.Println("Warning: ", err)
	}
}

func (subUrl *samUrl) Fatal(err error) {
	if err != nil {
        defer subUrl.cleanupDirectory()
		log.Fatal("Fatal: ", err)
	}
}

func newSamUrl(requestdir string) (samUrl){
    log.Println("Creating a new cache directory.")
    var subUrl samUrl
    subUrl.createDirectory(requestdir)
    return subUrl
}

func newSamUrlHttp(request *http.Request) (samUrl){
    log.Println("Creating a new cache directory.")
    var subUrl samUrl
    log.Println(subUrl.subDirectory)
    subUrl.createDirectory(request.Host + request.URL.Path)
    return subUrl
}
