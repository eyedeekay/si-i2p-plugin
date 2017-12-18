package main

import (
    "bufio"
    //"bytes"
    "fmt"
	"io/ioutil"
	"log"
	"net/http"
    "os"
    "path/filepath"
    //"strings"
    "strconv"
    "syscall"
    //"net/url"

)

type samUrl struct{
    err error

    transport *http.Transport
    http *http.Client
    subdirectory string

    recvPath string
    recvPipe *os.File
    //recvWriter bufio.Writer
    //recvReader bufio.Scanner

    delPath string
    delPipe *os.File
    delBuff bufio.Reader
}

func (subUrl *samUrl) initPipes(){
    pathConnectionExists, pathErr := exists(filepath.Join(connectionDirectory, subUrl.subdirectory))
    fmt.Println("Directory Check", filepath.Join(connectionDirectory, subUrl.subdirectory))
    subUrl.checkErr(pathErr)
    if ! pathConnectionExists {
        fmt.Println("Creating a connection:", subUrl.subdirectory)
        os.Mkdir(filepath.Join(connectionDirectory, subUrl.subdirectory), 0755)
    }

    subUrl.recvPath = filepath.Join(connectionDirectory, subUrl.subdirectory, "recv")
    pathRecvExists, recvPathErr := exists(subUrl.recvPath)
    subUrl.checkErr(recvPathErr)
    if ! pathRecvExists {
        subUrl.recvPipe, subUrl.err = os.Create(subUrl.recvPath)
        fmt.Println("Preparing to create File:", subUrl.recvPath)
        subUrl.checkErr(subUrl.err)
        fmt.Println("checking for problems...")
        fmt.Println("Opening the File...")
        subUrl.recvPipe, subUrl.err = os.OpenFile(subUrl.recvPath, os.O_RDWR|os.O_CREATE, 0755)
        fmt.Println("Created a File for recieving responses:", subUrl.recvPath)
    }

    subUrl.delPath = filepath.Join(connectionDirectory, subUrl.subdirectory, "del")
    pathDelExists, delPathErr := exists(subUrl.delPath)
    subUrl.checkErr(delPathErr)
    if ! pathDelExists{
        err := syscall.Mkfifo(subUrl.delPath, 0755)
        fmt.Println("Preparing to create Pipe:", subUrl.delPath)
        subUrl.checkErr(err)
        fmt.Println("checking for problems...")
        subUrl.delPipe, err = os.OpenFile(subUrl.delPath , os.O_RDWR|os.O_CREATE, 0755)
        fmt.Println("Opening the Named Pipe as a File...")
        subUrl.delBuff = *bufio.NewReader(subUrl.delPipe)
        fmt.Println("Opening the Named Pipe as a Buffer...")
        fmt.Println("Created a named Pipe for closing the connection:", subUrl.delPath)
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
    /*for subUrl.recvReader.Scan() {
        s += subUrl.recvReader.Text()
    }*/
    if s != "" {
        return s, len(s)
    }else{
       return "", 0
    }
}

func (subUrl *samUrl) dirSet(requestdir string) string {
    return requestdir
}

func (subUrl *samUrl) copyDirectory(response *http.Response, directory string) bool{
    b := false
    if directory == subUrl.subdirectory {
        if response.StatusCode == http.StatusOK {
            defer response.Body.Close()
            body, _ := ioutil.ReadAll(response.Body)
            subUrl.recvPipe.Write(body)
        }
        b = true
    }
    return b
}

func (subUrl *samUrl) cleanupDirectory(){
    subUrl.recvPipe.Close()
    subUrl.delPipe.Close()
    os.RemoveAll(filepath.Join(connectionDirectory, subUrl.subdirectory))
}

func (subUrl *samUrl) readDelete() int {
    line, _, err := subUrl.delBuff.ReadLine()
    subUrl.checkErr(err)
    n := len(line)
    fmt.Println("Reading n bytes from exit pipe:", strconv.Itoa(n))
    if n < 0 {
        fmt.Println("Something wierd happened with :", line)
        fmt.Println("end determined at index :", strconv.Itoa(n))
        return n
    }else{
        s := string( line[:n] )
        if s == "y" {
            fmt.Println("Deleting connection: %s", subUrl.subdirectory )
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
    fmt.Println("Creating a new cache directory.")
    var subUrl samUrl
    subUrl.subdirectory = requestdir
    subUrl.createDirectory(requestdir)
    return subUrl
}

