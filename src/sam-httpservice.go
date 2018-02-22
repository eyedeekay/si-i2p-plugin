package main

import (
    "bufio"
	//"io"
	"log"
	"net/http"
    "os"
    "path/filepath"
    //"strings"
    "syscall"
    //"net/url"

	//"github.com/eyedeekay/gosam"
)

type samHttpService struct{
    subCache []samUrl
    err error

    transport *http.Transport
    subClient *http.Client
    host string
    directory string

    servPath string
    servPipe *os.File
    servScan bufio.Scanner

    namePath string
    nameFile *os.File
    name string
}

func (samService *samHttpService) initPipes(){
    pathConnectionExists, pathErr := exists(filepath.Join(connectionDirectory, samService.host))
    log.Println("Directory Check", filepath.Join(connectionDirectory, samService.host))
    samService.Fatal(pathErr)
    if ! pathConnectionExists {
        log.Println("Creating a connection:", samService.host)
        os.Mkdir(filepath.Join(connectionDirectory, samService.host), 0755)
    }

    samService.servPath = filepath.Join(connectionDirectory, samService.host, "serv")
    pathservExists, servPathErr := exists(samService.servPath)
    samService.Fatal(servPathErr)
    if ! pathservExists {
        err := syscall.Mkfifo(samService.servPath, 0755)
        log.Println("Preparing to create Pipe:", samService.servPath)
        samService.Fatal(err)
        log.Println("checking for problems...")
        samService.servPipe, err = os.OpenFile(samService.servPath , os.O_RDWR|os.O_CREATE, 0755)
        log.Println("Opening the Named Pipe as a File...")
        samService.servScan = *bufio.NewScanner(samService.servPipe)
        samService.servScan.Split(bufio.ScanLines)
        log.Println("Opening the Named Pipe as a Buffer...")
        log.Println("Created a named Pipe for connecting to an http server:", samService.servPath)
    }

    samService.namePath = filepath.Join(connectionDirectory, samService.host, "name")
    pathNameExists, recvNameErr := exists(samService.namePath)
    samService.Fatal(recvNameErr)
    if ! pathNameExists {
        samService.nameFile, samService.err = os.Create(samService.namePath)
        log.Println("Preparing to create File:", samService.namePath)
        samService.Fatal(samService.err)
        log.Println("checking for problems...")
        log.Println("Opening the File...")
        samService.nameFile, samService.err = os.OpenFile(samService.namePath, os.O_RDWR|os.O_CREATE, 0644)
        log.Println("Created a File for the full name:", samService.namePath)
    }

}

func (samService *samHttpService) Warn(err error) {
	if err != nil {
        log.Println("Warning: ", err)
	}
}

func (samService *samHttpService) Fatal(err error) {
	if err != nil {
        defer samService.cleanupClient()
        log.Fatal("Fatal: ", err)
	}
}
func (samService *samHttpService) cleanupClient(){
    samService.servPipe.Close()
    samService.nameFile.Close()
    for _, url := range samService.subCache {
        url.cleanupDirectory()
    }
    os.RemoveAll(filepath.Join(connectionDirectory, samService.host))
}
