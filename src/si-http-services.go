package main

import (
	"bufio"
	"log"
	"path/filepath"
	//"net/http"
	"io"
	"os"
	"strings"
)

type samServices struct {
	listOfServices []samHttpService
	samAddrString  string
	samPortString  string
	err            error
    c              bool
	up             bool
    dir            string

	genrPath string
	genrPipe *os.File
	genrScan *bufio.Scanner

	lsPath string
	lsPipe *os.File

	delPath string
	delPipe *os.File
	delScan *bufio.Scanner
}

func (samServiceStack *samServices) initPipes() {
	setupFolder(samServiceStack.dir)

    samServiceStack.genrPath, samServiceStack.genrPipe, samServiceStack.err = setupFiFo(filepath.Join(connectionDirectory, samServiceStack.dir), "genr")
    if samServiceStack.c, samServiceStack.err = Fatal(samServiceStack.err, "Pipe setup error", "Pipe setup"); samServiceStack.c {
        samServiceStack.genrScan, samServiceStack.err = setupScanner(filepath.Join(connectionDirectory, samServiceStack.dir), "genr", samServiceStack.genrPipe)
        if samServiceStack.c, samServiceStack.err = Fatal(samServiceStack.err, "Scanner setup Error:", "Scanner set up successfully."); !samServiceStack.c {
            samServiceStack.cleanupServices()
        }
    }

    samServiceStack.lsPath, samServiceStack.lsPipe, samServiceStack.err = setupFiFo(filepath.Join(connectionDirectory, samServiceStack.dir), "ls")
    if samServiceStack.c, samServiceStack.err = Fatal(samServiceStack.err, "Pipe setup error", "Pipe setup"); samServiceStack.c {
        samServiceStack.lsPipe.WriteString("")
    }

	samServiceStack.delPath, samServiceStack.delPipe, samServiceStack.err = setupFiFo(filepath.Join(connectionDirectory, samServiceStack.dir), "del")
    if samServiceStack.c, samServiceStack.err = Fatal(samServiceStack.err, "Pipe setup error", "Pipe setup"); samServiceStack.c {
        samServiceStack.delScan, samServiceStack.err = setupScanner(filepath.Join(connectionDirectory, samServiceStack.dir), "del", samServiceStack.delPipe)
        if samServiceStack.c, samServiceStack.err = Fatal(samServiceStack.err, "Scanner setup Error:", "Scanner set up successfully."); !samServiceStack.c {
            samServiceStack.cleanupServices()
        }
    }
	samServiceStack.up = true
}

func (samServiceStack *samServices) createService(alias string) {
	Log("Appending service to SAM service stack.")
	samServiceStack.listOfServices = append(samServiceStack.listOfServices, createSamHttpService(samServiceStack.samAddrString, samServiceStack.samPortString, alias))
}

func (samServiceStack *samServices) findService(request string) *samHttpService {
	found := false
	var s samHttpService
	for index, service := range samServiceStack.listOfServices {
		log.Println("Checking client requests", index+1)
		log.Println("of", len(samServiceStack.listOfServices))
		if service.serviceCheck(request) {
			Log("Client pipework for %s found.", request)
			Log("Request sent")
			found = true
			return &service
		}
	}
	if !found {
		Log("Client pipework for %s not found: Creating.", request)
		samServiceStack.createService(request)
		for index, service := range samServiceStack.listOfServices {
			log.Println("Checking client requests", index+1)
			log.Println("of", len(samServiceStack.listOfServices))
			if service.serviceCheck(request) {
				Log("Client pipework for %s found.", request)
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
	Log("Established SAM connection")
	if !samServiceStack.up {
		samServiceStack.initPipes()
		Log("Parent proxy pipes initialized. Parent proxy set to up.")
	}
}

func (samServiceStack *samServices) sendServiceRequest(index string) {
	samServiceStack.findService(index).sendContent(index)
}

func (samServiceStack *samServices) responsify(input string) io.Reader {
	tmp := strings.NewReader(input)
	Log("Responsifying string:")
	return tmp
}

func (samServiceStack *samServices) readRequest() {
	Log("Reading requests:")
	for samServiceStack.genrScan.Scan() {
		if samServiceStack.genrScan.Text() != "" {
			go samServiceStack.sendServiceRequest(samServiceStack.genrScan.Text())
		}
	}
}

func (samServiceStack *samServices) writeDetails(details string) bool {
	b := false
	if details != "" {
		Log("Got response:")
		io.Copy(samServiceStack.lsPipe, samServiceStack.responsify(details))
		b = true
	}
	return b
}

func (samServiceStack *samServices) writeResponses() {
	Log("Writing responses:")
	for i, service := range samServiceStack.listOfServices {
		log.Println("Checking for responses: %s", i+1)
		log.Println("of: ", len(samServiceStack.listOfServices))
		if service.printDetails() != "" {
			go samServiceStack.writeDetails(service.printDetails())
		}
	}
}

func (samServiceStack *samServices) readDelete() bool {
	Log("Managing pipes:")
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

func (samServiceStack *samServices) cleanupServices() {
	samServiceStack.genrPipe.Close()
	samServiceStack.lsPipe.Close()
	for _, service := range samServiceStack.listOfServices {
		service.cleanupService()
	}
	samServiceStack.delPipe.Close()
	os.RemoveAll(filepath.Join(connectionDirectory, "service"))
}

func (samServiceStack *samServices) run(){

}

func createSamServiceList(samAddr string, samPort string) *samServices {
	var samServiceList samServices
    samServiceList.dir = "services"
	samServiceList.createServiceList(samAddr, samPort)
	return &samServiceList
}
