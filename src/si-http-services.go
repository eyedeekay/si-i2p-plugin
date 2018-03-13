package main

import (
	"bufio"
	"log"
	"path/filepath"
	//"net/http"
	"io"
	"os"
	"strings"
	"syscall"
)

type samServices struct {
	listOfServices []samHttpService
	samAddrString  string
	samPortString  string
	err            error
	up             bool

	genrPath string
	genrPipe *os.File
	genrScan bufio.Scanner

	lsPath string
	lsPipe *os.File

	delPath string
	delPipe *os.File
	delScan bufio.Scanner
}

func (samServiceStack *samServices) initPipes() {
	pathConnectionExists, err := exists(filepath.Join(connectionDirectory, "service"))
	samServiceStack.Fatal(err)
	if !pathConnectionExists {
		samServiceStack.Log("Creating a connection:", "service")
		os.Mkdir(filepath.Join(connectionDirectory, "service"), 0755)
	} else {
		os.RemoveAll(filepath.Join(connectionDirectory, "service"))
		samServiceStack.Log("Creating a connection:", "service")
		os.Mkdir(filepath.Join(connectionDirectory, "service"), 0755)
	}

	samServiceStack.genrPath = filepath.Join(connectionDirectory, "service", "genr")
	pathgenrExists, genrErr := exists(samServiceStack.genrPath)
	samServiceStack.Fatal(genrErr)
	if !pathgenrExists {
		samServiceStack.err = syscall.Mkfifo(samServiceStack.genrPath, 0755)
		samServiceStack.Log("Preparing to create Pipe:", samServiceStack.genrPath)
		samServiceStack.Fatal(samServiceStack.err)
		samServiceStack.Log("checking for problems...")
		samServiceStack.genrPipe, samServiceStack.err = os.OpenFile(samServiceStack.genrPath, os.O_RDWR|os.O_CREATE, 0755)
		samServiceStack.Log("Opening the Named Pipe as a Scanner...")
		samServiceStack.genrScan = *bufio.NewScanner(samServiceStack.genrPipe)
		samServiceStack.genrScan.Split(bufio.ScanLines)
		samServiceStack.Log("Opening the Named Pipe as a Scanner...")
		samServiceStack.Log("Created a named Pipe for generating new i2p http services:", samServiceStack.genrPath)
	}

	samServiceStack.lsPath = filepath.Join(connectionDirectory, "service", "ls")
	pathlsExists, lsErr := exists(samServiceStack.lsPath)
	samServiceStack.Fatal(lsErr)
	if !pathlsExists {
		samServiceStack.err = syscall.Mkfifo(samServiceStack.lsPath, 0755)
		samServiceStack.Log("Preparing to create Pipe:", samServiceStack.lsPath)
		samServiceStack.Fatal(samServiceStack.err)
		samServiceStack.Log("checking for problems...")
		samServiceStack.lsPipe, samServiceStack.err = os.OpenFile(samServiceStack.lsPath, os.O_RDWR|os.O_CREATE, 0755)
		samServiceStack.lsPipe.WriteString("")
		samServiceStack.Log("Created a named Pipe for monitoring service information:", samServiceStack.lsPath)
	}

	samServiceStack.delPath = filepath.Join(connectionDirectory, "service", "del")
	pathDelExists, delErr := exists(samServiceStack.delPath)
	samServiceStack.Fatal(delErr)
	if !pathDelExists {
		samServiceStack.err = syscall.Mkfifo(samServiceStack.delPath, 0755)
		samServiceStack.Log("Preparing to create Pipe:", samServiceStack.delPath)
		samServiceStack.Fatal(samServiceStack.err)
		samServiceStack.Log("checking for problems...")
		samServiceStack.delPipe, samServiceStack.err = os.OpenFile(samServiceStack.delPath, os.O_RDWR|os.O_CREATE, 0755)
		samServiceStack.lsPipe.WriteString("")
		samServiceStack.Log("Opening the Named Pipe as a File...")
		samServiceStack.delScan = *bufio.NewScanner(samServiceStack.delPipe)
		samServiceStack.delScan.Split(bufio.ScanLines)
		samServiceStack.Log("Opening the Named Pipe as a Scanner...")
		samServiceStack.Log("Created a named Pipe for closing all i2p http services:", samServiceStack.delPath)
	}
	samServiceStack.up = true
}

func (samServiceStack *samServices) createService(alias string) {
	samServiceStack.Log("Appending service to SAM service stack.")
	samServiceStack.listOfServices = append(samServiceStack.listOfServices, createSamHttpService(samServiceStack.samAddrString, samServiceStack.samPortString, alias))
}

func (samServiceStack *samServices) findService(request string) *samHttpService {
	found := false
	var s samHttpService
	for index, service := range samServiceStack.listOfServices {
		log.Println("Checking client requests", index+1)
		log.Println("of", len(samServiceStack.listOfServices))
		if service.serviceCheck(request) {
			samServiceStack.Log("Client pipework for %s found.", request)
			samServiceStack.Log("Request sent")
			found = true
			return &service
		}
	}
	if !found {
		samServiceStack.Log("Client pipework for %s not found: Creating.", request)
		samServiceStack.createService(request)
		for index, service := range samServiceStack.listOfServices {
			log.Println("Checking client requests", index+1)
			log.Println("of", len(samServiceStack.listOfServices))
			if service.serviceCheck(request) {
				samServiceStack.Log("Client pipework for %s found.", request)
				s = service
			}
		}
	}
	return &s
}

func (samServiceStack *samServices) createServiceList(samAddr string, samPort string) {
	samServiceStack.samAddrString = samAddr
	samServiceStack.samPortString = samPort
	//samServiceStack.
	samServiceStack.Log("Established SAM connection")
	if !samServiceStack.up {
		samServiceStack.initPipes()
		samServiceStack.Log("Parent proxy pipes initialized. Parent proxy set to up.")
	}
}

func (samServiceStack *samServices) sendServiceRequest(index string) {
	samServiceStack.findService(index).sendContent(index)
}

func (samServiceStack *samServices) responsify(input string) io.Reader {
	tmp := strings.NewReader(input)
	samServiceStack.Log("Responsifying string:")
	return tmp
}

func (samServiceStack *samServices) readRequest() {
	samServiceStack.Log("Reading requests:")
	for samServiceStack.genrScan.Scan() {
		if samServiceStack.genrScan.Text() != "" {
			go samServiceStack.sendServiceRequest(samServiceStack.genrScan.Text())
		}
	}
}

func (samServiceStack *samServices) writeDetails(details string) bool {
	b := false
	if details != "" {
		samServiceStack.Log("Got response:")
		io.Copy(samServiceStack.lsPipe, samServiceStack.responsify(details))
		b = true
	}
	return b
}

func (samServiceStack *samServices) writeResponses() {
	samServiceStack.Log("Writing responses:")
	for i, service := range samServiceStack.listOfServices {
		log.Println("Checking for responses: %s", i+1)
		log.Println("of: ", len(samServiceStack.listOfServices))
		if service.printDetails() != "" {
			go samServiceStack.writeDetails(service.printDetails())
		}
	}
}

func (samServiceStack *samServices) readDelete() bool {
	samServiceStack.Log("Managing pipes:")
	for samServiceStack.delScan.Scan() {
		if samServiceStack.delScan.Text() == "y" || samServiceStack.delScan.Text() == "Y" {
			defer samServiceStack.cleanupServices()
			return true
		} else {
			return false
		}
	}
	return false
}

//func (samServiceStack *samServices) Blank() {}

func (samServiceStack *samServices) Log(msg ...string) {
	if verbose {
		log.Println("LOG: ", msg)
	}
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

func (samServiceStack *samServices) cleanupServices() {
	samServiceStack.genrPipe.Close()
	samServiceStack.lsPipe.Close()
	for _, client := range samServiceStack.listOfServices {
		client.cleanupClient()
	}
	samServiceStack.delPipe.Close()
	os.RemoveAll(filepath.Join(connectionDirectory, "service"))
}

func createSamServiceList(samAddr string, samPort string) *samServices {
	var samServiceList samServices
	samServiceList.createServiceList(samAddr, samPort)
	return &samServiceList
}
