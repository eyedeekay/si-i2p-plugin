package main

import (
    "bufio"
    /*"bytes"*/
    "fmt"
	"io"
	"log"
	"net/http"
    "os"
    "path/filepath"
    /*"strings"
    "strconv"*/
    "syscall"
    /*"net/url"*/

	//"github.com/eyedeekay/gosam"
)

type samUrl struct{
    err error

    transport *http.Transport
    http *http.Client
    subdirectory string

    recvPath string
    recvPipe *os.File
    recvBuff bufio.Scanner

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
    }else{
        os.RemoveAll(filepath.Join(connectionDirectory, subUrl.subdirectory))
        fmt.Println("Creating a connection:", subUrl.subdirectory)
        os.Mkdir(filepath.Join(connectionDirectory, subUrl.subdirectory), 0755)
    }

    subUrl.recvPath = filepath.Join(connectionDirectory, subUrl.subdirectory, "recv")
    pathRecvExists, recvPathErr := exists(subUrl.recvPath)
    subUrl.checkErr(recvPathErr)
    if ! pathRecvExists {
        subUrl.recvPipe, subUrl.err = os.Create(subUrl.recvPath)
        fmt.Println("Preparing to create Pipe:", subUrl.recvPath)
        subUrl.checkErr(subUrl.err)
        fmt.Println("checking for problems...")
        subUrl.recvBuff = *bufio.NewScanner(subUrl.recvPipe)
        fmt.Println("Opening the Named Pipe as a Buffer...")
        fmt.Println("Created a named Pipe for recieving responses:", subUrl.recvPath)
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
    if subUrl.subdirectory == "" {
        subUrl.subdirectory = subUrl.dirSet(requestdir)
        subUrl.initPipes()
    }
}

func (subUrl *samUrl) scannerText() (string, int) {
    s := ""
    for subUrl.recvBuff.Scan() {
        //samConn.recvBuff.Scan()
        s += subUrl.recvBuff.Text()
    }
    if s != "" {
        return s, len(s)
    }else{
       return "no response", 0
    }
}

func (subUrl *samUrl) dirSet(requestdir string) string {
    return requestdir
}

func (subUrl *samUrl) copyDirectory(response io.Reader, directory string){
    if directory == subUrl.subdirectory {
        io.Copy(subUrl.recvPipe, response)
    }
}

func (subUrl *samUrl) cleanupDirectory(){
    subUrl.recvPipe.Close()
    subUrl.delPipe.Close()
    os.RemoveAll(filepath.Join(connectionDirectory, subUrl.subdirectory))
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

