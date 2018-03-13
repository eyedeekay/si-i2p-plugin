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
	"github.com/eyedeekay/gosam"
)

type samHttpService struct {
	subCache        []samUrl
	samBridgeClient *goSam.Client
	err             error

	transport *http.Transport
	subClient *http.Client
	host      string
	directory string

	servPath string
	servPipe *os.File
	servScan bufio.Scanner

	namePath string
	nameFile *os.File
	name     string

	idPath string
	idFile *os.File
	id     int32

	base64Path string
	base64File *os.File
	base64     string
}

func (samService *samHttpService) initPipes() {
	pathConnectionExists, pathErr := exists(filepath.Join(connectionDirectory, samService.host))
	log.Println("Directory Check", filepath.Join(connectionDirectory, samService.host))
	samService.Fatal(pathErr, "Directory Check Error", "Directory Check", filepath.Join(connectionDirectory, samService.host))
	if !pathConnectionExists {
		log.Println("Creating a connection:", samService.host)
		os.Mkdir(filepath.Join(connectionDirectory, samService.host), 0755)
	}

	samService.servPath = filepath.Join(connectionDirectory, samService.host, "serv")
	pathservExists, servPathErr := exists(samService.servPath)
	samService.Fatal(servPathErr, "Serve File Check Error", "Serve file check", samService.servPath)
	if !pathservExists {
		err := syscall.Mkfifo(samService.servPath, 0755)
		log.Println("Preparing to create Pipe:", samService.servPath)
		samService.Fatal(err, "File Creation Error", "Creating File", samService.servPath)
		log.Println("checking for problems...")
		samService.servPipe, err = os.OpenFile(samService.servPath, os.O_RDWR|os.O_CREATE, 0755)
		log.Println("Opening the Named Pipe as a File...")
		samService.servScan = *bufio.NewScanner(samService.servPipe)
		samService.servScan.Split(bufio.ScanLines)
		log.Println("Opening the Named Pipe as a Buffer...")
		log.Println("Created a named Pipe for connecting to an http server:", samService.servPath)
	}

	samService.namePath = filepath.Join(connectionDirectory, samService.host, "name")
	pathNameExists, recvNameErr := exists(samService.namePath)
	samService.Fatal(recvNameErr, "Name File Check error", "Name File Check", samService.namePath)
	if !pathNameExists {
		samService.nameFile, samService.err = os.Create(samService.namePath)
		log.Println("Preparing to create File:", samService.namePath)
		samService.Fatal(samService.err, "File Creation Error", "Creating File", samService.namePath)
		log.Println("checking for problems...")
		log.Println("Opening the File...")
		samService.nameFile, samService.err = os.OpenFile(samService.namePath, os.O_RDWR|os.O_CREATE, 0644)
		log.Println("Created a File for the full name:", samService.namePath)
	}

}

func (samService *samHttpService) sendContent(index string) (*http.Response, error) {
	/*r, dir := samService.getURL(index)
	samService.Log("Getting resource", index)
	resp, err := samService.subClient.Get(r)
	samService.Warn(err, "Response Error", "Getting Response")
	samService.Log("Pumping result to top of parent pipe")
	samService.copyRequest(resp, dir)
	return resp, err*/
	return nil, nil
}

func (samService *samHttpService) serviceCheck(alias string) bool {
	return false
}

func (samService *samHttpService) scannerText() (string, error) {
	text := ""
	var err error
	for _, url := range samService.subCache {
		text, err = url.scannerText()
		if len(text) > 0 {
			break
		}
	}
	return text, err
}

/*
func (samConn *samHttp) writeName(request string) {
	if samConn.checkName() {
		samConn.host, samConn.directory = samConn.hostSet(request)
		samConn.Log("Setting hostname:", samConn.host)
		samConn.Log("Looking up hostname:", samConn.host)
		samConn.name, samConn.err = samConn.samBridgeClient.Lookup(samConn.host)
		samConn.nameFile.WriteString(samConn.name)
		samConn.Log("Caching base64 address of:", samConn.host+" "+samConn.name)
		samConn.id, samConn.base64, samConn.err = samConn.samBridgeClient.CreateStreamSession("")
		samConn.idFile.WriteString(fmt.Sprint(samConn.id))
		samConn.Warn(samConn.err, "Local Base64 Caching error", "Cachine Base64 Address of:", request)
		log.Println("Tunnel id: ", samConn.id)
		samConn.Log("Tunnel dest: ", samConn.base64)
		samConn.base64File.WriteString(samConn.base64)
		samConn.Log("New Connection Name: ", samConn.base64)
	} else {
		samConn.host, samConn.directory = samConn.hostSet(request)
		samConn.Log("Setting hostname:", samConn.host)
		samConn.initPipes()
		samConn.Log("Looking up hostname:", samConn.host)
		samConn.name, samConn.err = samConn.samBridgeClient.Lookup(samConn.host)
		samConn.Log("Caching base64 address of:", samConn.host+" "+samConn.name)
		samConn.nameFile.WriteString(samConn.name)
		samConn.id, samConn.base64, samConn.err = samConn.samBridgeClient.CreateStreamSession("")
		samConn.idFile.WriteString(fmt.Sprint(samConn.id))
		samConn.Warn(samConn.err, "Local Base64 Caching error", "Cachine Base64 Address of:", request)
		log.Println("Tunnel id: ", samConn.id)
		samConn.Log("Tunnel dest: ", samConn.base64)
		samConn.base64File.WriteString(samConn.base64)
		samConn.Log("New Connection Name: ", samConn.base64)
	}
}
*/

func (samService *samHttpService) printDetails() string {
	s, e := samService.scannerText()
	samService.Fatal(e, "Response Retrieval Error", "Retrieving Responses")
	return s
}

func (samService *samHttpService) Log(msg ...string) {
	if verbose {
		log.Println("LOG: ", msg)
	}
}

func (samService *samHttpService) Warn(err error, errmsg string, msg ...string) {
	if err != nil {
		log.Println("Warning: ", err)
	}
}

func (samService *samHttpService) Fatal(err error, errmsg string, msg ...string) {
	if err != nil {
		defer samService.cleanupClient()
		log.Fatal("Fatal: ", err)
	}
}
func (samService *samHttpService) cleanupClient() {
	samService.servPipe.Close()
	samService.nameFile.Close()
	for _, url := range samService.subCache {
		url.cleanupDirectory()
	}
	err := samService.samBridgeClient.Close()
	samService.Fatal(err, "SAM Service Connection Closing Error", "Closing SAM service Connection")
	os.RemoveAll(filepath.Join(connectionDirectory, samService.host))
}

func createSamHttpService(samAddr string, samPort string, alias string) samHttpService {
	var samService samHttpService
	return samService
}
