package main

import (
    "bufio"
    "path/filepath"
    "log"
    //"net/http"
    //"io"
    "os"
    //"strings"
    "syscall"

    "github.com/eyedeekay/gosam"
)


type samServices struct{
    listOfServices []samHttpService
    samBridgeClient *goSam.Client
    err error
    up bool

    samAddrString string
    samPortString string

    genrPath string
    genrPipe *os.File
    genrScan bufio.Scanner

    lsPath string
    lsPipe *os.File

    delPath string
    delPipe *os.File
    delScan bufio.Scanner
}

func (samServiceStack * samServices) initPipes(){
    pathConnectionExists, err := exists(filepath.Join(connectionDirectory, "service"))
    samServiceStack.Fatal(err)
    if ! pathConnectionExists {
        log.Println("Creating a connection:", "service")
        os.Mkdir(filepath.Join(connectionDirectory, "service"), 0755)
    }else{
        os.RemoveAll(filepath.Join(connectionDirectory, "service"))
        log.Println("Creating a connection:", "service")
        os.Mkdir(filepath.Join(connectionDirectory, "service"), 0755)
    }

    samServiceStack.genrPath = filepath.Join(connectionDirectory, "service", "genr")
    pathgenrExists, genrErr := exists(samServiceStack.genrPath)
    samServiceStack.Fatal(genrErr)
    if ! pathgenrExists {
        samServiceStack.err = syscall.Mkfifo(samServiceStack.genrPath, 0755)
        log.Println("Preparing to create Pipe:", samServiceStack.genrPath)
        samServiceStack.Fatal(samServiceStack.err)
        log.Println("checking for problems...")
        samServiceStack.genrPipe, samServiceStack.err = os.OpenFile(samServiceStack.genrPath , os.O_RDWR|os.O_CREATE, 0755)
        log.Println("Opening the Named Pipe as a Scanner...")
        samServiceStack.genrScan = *bufio.NewScanner(samServiceStack.genrPipe)
        samServiceStack.genrScan.Split(bufio.ScanLines)
        log.Println("Opening the Named Pipe as a Scanner...")
        log.Println("Created a named Pipe for generating new i2p http services:", samServiceStack.genrPath)
    }

    samServiceStack.lsPath = filepath.Join(connectionDirectory, "service", "ls")
    pathlsExists, lsErr := exists(samServiceStack.lsPath)
    samServiceStack.Fatal(lsErr)
    if ! pathlsExists {
        samServiceStack.err = syscall.Mkfifo(samServiceStack.lsPath, 0755)
        log.Println("Preparing to create Pipe:", samServiceStack.lsPath)
        samServiceStack.Fatal(samServiceStack.err)
        log.Println("checking for problems...")
        samServiceStack.lsPipe, samServiceStack.err = os.OpenFile(samServiceStack.lsPath , os.O_RDWR|os.O_CREATE, 0755)
        samServiceStack.lsPipe.WriteString("")
        log.Println("Created a named Pipe for monitoring service information:", samServiceStack.lsPath)
    }

    samServiceStack.delPath = filepath.Join(connectionDirectory, "service", "del")
    pathDelExists, delErr := exists(samServiceStack.delPath)
    samServiceStack.Fatal(delErr)
    if ! pathDelExists{
        samServiceStack.err = syscall.Mkfifo(samServiceStack.delPath, 0755)
        log.Println("Preparing to create Pipe:", samServiceStack.delPath)
        samServiceStack.Fatal(samServiceStack.err)
        log.Println("checking for problems...")
        samServiceStack.delPipe, samServiceStack.err = os.OpenFile(samServiceStack.delPath , os.O_RDWR|os.O_CREATE, 0755)
        samServiceStack.lsPipe.WriteString("")
        log.Println("Opening the Named Pipe as a File...")
        samServiceStack.delScan = *bufio.NewScanner(samServiceStack.delPipe)
        samServiceStack.delScan.Split(bufio.ScanLines)
        log.Println("Opening the Named Pipe as a Scanner...")
        log.Println("Created a named Pipe for closing all i2p http services:", samServiceStack.delPath)
    }
    samServiceStack.up = true;
}

func (samServiceStack *samServices) Warn(err error) {
	if err != nil {
        log.Println("Warning: ", err)
	}
}

func (samServiceStack *samServices) Fatal(err error) {
	if err != nil {
        defer samServiceStack.cleanupServices()
        log.Fatal("Fatal: ", err)
	}
}

func (samServiceStack *samServices) cleanupServices(){
    samServiceStack.genrPipe.Close()
    samServiceStack.lsPipe.Close()
    for _, client := range samServiceStack.listOfServices {
        client.cleanupClient()
    }
    samServiceStack.delPipe.Close()
    err := samServiceStack.samBridgeClient.Close()
    samServiceStack.Fatal(err)
    os.RemoveAll(filepath.Join(connectionDirectory, "service"))
}

func newSamServiceList(samStack *samList) *samServices {
    var samServiceList samServices
    return &samServiceList
}
