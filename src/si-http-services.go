package dii2p

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

import (
	"github.com/eyedeekay/si-i2p-plugin/src/errors"
	"github.com/eyedeekay/si-i2p-plugin/src/helpers"
)

// SamServices is a structure for managing SAM services
type SamServices struct {
	listOfServices []samHTTPService
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

func (samServiceStack *SamServices) initPipes() {
	dii2phelper.SetupFolder(samServiceStack.dir)

	samServiceStack.genrPath, samServiceStack.genrPipe, samServiceStack.err = dii2phelper.SetupFiFo(filepath.Join(dii2phelper.ConnectionDirectory, samServiceStack.dir), "genr")
	if samServiceStack.c, samServiceStack.err = dii2perrs.Fatal(samServiceStack.err, "Pipe setup error", "Pipe setup"); samServiceStack.c {
		samServiceStack.genrScan, samServiceStack.err = dii2phelper.SetupScanner(filepath.Join(dii2phelper.ConnectionDirectory, samServiceStack.dir), "genr", samServiceStack.genrPipe)
		if samServiceStack.c, samServiceStack.err = dii2perrs.Fatal(samServiceStack.err, "Scanner setup Error:", "Scanner set up successfully."); !samServiceStack.c {
			samServiceStack.cleanupServices()
		}
	}

	samServiceStack.lsPath, samServiceStack.lsPipe, samServiceStack.err = dii2phelper.SetupFiFo(filepath.Join(dii2phelper.ConnectionDirectory, samServiceStack.dir), "ls")
	if samServiceStack.c, samServiceStack.err = dii2perrs.Fatal(samServiceStack.err, "Pipe setup error", "Pipe setup"); samServiceStack.c {
		samServiceStack.lsPipe.WriteString("")
	}

	samServiceStack.delPath, samServiceStack.delPipe, samServiceStack.err = dii2phelper.SetupFiFo(filepath.Join(dii2phelper.ConnectionDirectory, samServiceStack.dir), "del")
	if samServiceStack.c, samServiceStack.err = dii2perrs.Fatal(samServiceStack.err, "Pipe setup error", "Pipe setup"); samServiceStack.c {
		samServiceStack.delScan, samServiceStack.err = dii2phelper.SetupScanner(filepath.Join(dii2phelper.ConnectionDirectory, samServiceStack.dir), "del", samServiceStack.delPipe)
		if samServiceStack.c, samServiceStack.err = dii2perrs.Fatal(samServiceStack.err, "Scanner setup Error:", "Scanner set up successfully."); !samServiceStack.c {
			samServiceStack.cleanupServices()
		}
	}
	samServiceStack.up = true
}

func (samServiceStack *SamServices) createService(alias string) {
	dii2perrs.Log("Appending service to SAM service stack.")
	samServiceStack.listOfServices = append(samServiceStack.listOfServices, createSamHTTPService(samServiceStack.samAddrString, samServiceStack.samPortString, alias))
}

func (samServiceStack *SamServices) findService(request string) *samHTTPService {
	found := false
	var s samHTTPService
	for index, service := range samServiceStack.listOfServices {
		log.Println("Checking client requests", index+1)
		log.Println("of", len(samServiceStack.listOfServices))
		if service.serviceCheck(request) {
			dii2perrs.Log("Client pipework for", request, "found.", request)
			dii2perrs.Log("Request sent")
			found = true
			return &service
		}
	}
	if !found {
		dii2perrs.Log("Client pipework for", request, "not found: Creating.")
		samServiceStack.createService(request)
		for index, service := range samServiceStack.listOfServices {
			log.Println("Checking client requests", index+1)
			log.Println("of", len(samServiceStack.listOfServices))
			if service.serviceCheck(request) {
				dii2perrs.Log("Client pipework for", request, "found.")
				s = service
			}
		}
	}
	return &s
}

func (samServiceStack *SamServices) createServiceList() {
	if !samServiceStack.up {
		samServiceStack.initPipes()
		dii2perrs.Log("Parent proxy pipes initialized. Parent proxy set to up.")
	}
}

func (samServiceStack *SamServices) sendServiceRequest(index string) {
	samServiceStack.findService(index).sendContent(index)
}

func (samServiceStack *SamServices) responsify(input string) io.Reader {
	tmp := strings.NewReader(input)
	dii2perrs.Log("Responsifying string:")
	return tmp
}

// ServiceRequest requests a new service interface from the SAM bridge
func (samServiceStack *SamServices) ServiceRequest() {
	dii2perrs.Log("Reading requests:")
	for samServiceStack.genrScan.Scan() {
		if samServiceStack.genrScan.Text() == "y" || samServiceStack.genrScan.Text() == "Y" || samServiceStack.genrScan.Text() == "g" || samServiceStack.genrScan.Text() == "G" || samServiceStack.genrScan.Text() == "n" || samServiceStack.genrScan.Text() == "N" || samServiceStack.genrScan.Text() == "new" {
			go samServiceStack.sendServiceRequest(samServiceStack.genrScan.Text())
		}
	}
}

func (samServiceStack *SamServices) writeDetails(details string) bool {
	b := false
	if details != "" {
		dii2perrs.Log("Got response:")
		io.Copy(samServiceStack.lsPipe, samServiceStack.responsify(details))
		b = true
	}
	return b
}

func (samServiceStack *SamServices) writeResponses() {
	dii2perrs.Log("Writing responses:")
	for i, service := range samServiceStack.listOfServices {
		dii2perrs.Log("Checking for responses: ", i+1)
		dii2perrs.Log("of: ", len(samServiceStack.listOfServices))
		if service.printDetails() != "" {
			go samServiceStack.writeDetails(service.printDetails())
		}
	}
}

// ReadDelete checks whether to shut down the service manager
func (samServiceStack *SamServices) ReadDelete() bool {
	dii2perrs.Log("Managing pipes:")
	for samServiceStack.delScan.Scan() {
		if samServiceStack.delScan.Text() == "y" || samServiceStack.delScan.Text() == "Y" {
			defer samServiceStack.cleanupServices()
			return true
		}
		return false
	}
	return false
}

func (samServiceStack *SamServices) cleanupServices() {
	samServiceStack.genrPipe.Close()
	samServiceStack.lsPipe.Close()
	for _, service := range samServiceStack.listOfServices {
		service.cleanupService()
	}
	samServiceStack.delPipe.Close()
	os.RemoveAll(filepath.Join(dii2phelper.ConnectionDirectory, "service"))
}

// CreateSamServiceList Creates a Service Manager from functional arguments
func CreateSamServiceList(opts ...func(*SamServices) error) (*SamServices, error) {
	var samServiceList SamServices
	samServiceList.dir = "services"
	for _, o := range opts {
		if err := o(&samServiceList); err != nil {
			return nil, err
		}
	}
	samServiceList.createServiceList()
	return &samServiceList, nil
}
