package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	//"github.com/eyedeekay/gosam"
	//"github.com/cryptix/goSam"
)

type samList struct {
	listOfClients []samHttp
	samAddrString string
	samPortString string
	err           error
	up            bool

	sendPath string
	sendPipe *os.File
	sendScan bufio.Scanner

	recvPath string
	recvPipe *os.File

	delPath string
	delPipe *os.File
	delScan bufio.Scanner
}

func (samStack *samList) initPipes() {
	pathConnectionExists, err := exists(filepath.Join(connectionDirectory, "parent"))
	samStack.Fatal(err, "Parent Directory Error", "Parent Directory Check", filepath.Join(connectionDirectory, "parent"))
	if !pathConnectionExists {
		samStack.Log("Creating a connection:", "parent")
		os.Mkdir(filepath.Join(connectionDirectory, "parent"), 0755)
	} else {
		os.RemoveAll(filepath.Join(connectionDirectory, "parent"))
		samStack.Log("Creating a connection:", "parent")
		os.Mkdir(filepath.Join(connectionDirectory, "parent"), 0755)
	}

	samStack.sendPath = filepath.Join(connectionDirectory, "parent", "send")
	pathSendExists, sendErr := exists(samStack.sendPath)
	samStack.Fatal(sendErr, "Send File Check Error", "Send File Check", samStack.sendPath)
	if !pathSendExists {
		samStack.err = syscall.Mkfifo(samStack.sendPath, 0755)
		samStack.Log("Preparing to create Pipe:", samStack.sendPath)
		samStack.Fatal(samStack.err, "Pipe Creation Error", "Creating Pipe", samStack.sendPath)
		samStack.Log("checking for problems...")
		samStack.sendPipe, samStack.err = os.OpenFile(samStack.sendPath, os.O_RDWR|os.O_CREATE, 0755)
		samStack.Log("Opening the Named Pipe as a Scanner...")
		samStack.sendScan = *bufio.NewScanner(samStack.sendPipe)
		samStack.sendScan.Split(bufio.ScanLines)
		samStack.Log("Opening the Named Pipe as a Scanner...")
		samStack.Log("Created a named Pipe for sending requests:", samStack.sendPath)
	}

	samStack.recvPath = filepath.Join(connectionDirectory, "parent", "recv")
	pathRecvExists, recvErr := exists(samStack.recvPath)
	samStack.Fatal(recvErr, "Recv File Check Error", "Recv File Check", samStack.recvPath)
	if !pathRecvExists {
		samStack.err = syscall.Mkfifo(samStack.recvPath, 0755)
		samStack.Log("Preparing to create Pipe:", samStack.recvPath)
		samStack.Fatal(samStack.err, "Pipe Creation Error", "Creating Pipe", samStack.recvPath)
		samStack.Log("checking for problems...")
		samStack.recvPipe, samStack.err = os.OpenFile(samStack.recvPath, os.O_RDWR|os.O_CREATE, 0755)
		samStack.recvPipe.WriteString("")
		samStack.Log("Created a named Pipe for recieving responses:", samStack.recvPath)
	}

	samStack.delPath = filepath.Join(connectionDirectory, "parent", "del")
	pathDelExists, delErr := exists(samStack.delPath)
	samStack.Fatal(delErr, "Del File Check Error", "Del File Check", samStack.delPath)
	if !pathDelExists {
		samStack.err = syscall.Mkfifo(samStack.delPath, 0755)
		samStack.Log("Preparing to create Pipe:", samStack.delPath)
		samStack.Fatal(samStack.err, "Pipe Creation Error", "Creating Pipe", samStack.delPath)
		samStack.Log("checking for problems...")
		samStack.delPipe, samStack.err = os.OpenFile(samStack.delPath, os.O_RDWR|os.O_CREATE, 0755)
		samStack.recvPipe.WriteString("")
		samStack.Log("Opening the Named Pipe as a File...")
		samStack.delScan = *bufio.NewScanner(samStack.delPipe)
		samStack.delScan.Split(bufio.ScanLines)
		samStack.Log("Opening the Named Pipe as a Scanner...")
		samStack.Log("Created a named Pipe for closing the connection:", samStack.delPath)
	}
	samStack.up = true
}

func (samStack *samList) createClient(request string) {
	samStack.Log("Appending client to SAM stack.")
	//samStack.listOfClients = append(samStack.listOfClients, newSamHttp(samStack.samAddrString, samStack.samPortString, samStack.samBridgeClient, request))
	samStack.listOfClients = append(samStack.listOfClients, newSamHttp(samStack.samAddrString, samStack.samPortString, request))
}

func (samStack *samList) createClientHttp(request *http.Request) {
	samStack.Log("Appending client to SAM stack.")
	samStack.listOfClients = append(samStack.listOfClients, newSamHttpHttp(samStack.samAddrString, samStack.samPortString, request))
}

func (samStack *samList) createSamList(samAddrString string, samPortString string) {
	samStack.samAddrString = samAddrString
	samStack.samPortString = samPortString
	//samStack.Log("Requesting a new SAM-based http client")
	//samStack.Fatal(samStack.err, "", )
	samStack.Log("Established SAM connection")
	if !samStack.up {
		samStack.initPipes()
		samStack.Log("Parent proxy pipes initialized. Parent proxy set to up.")
	}
}

func (samStack *samList) sendClientRequest(request string) {
	samStack.findClient(request).sendRequest(request)
}

func (samStack *samList) sendClientRequestHttp(request *http.Request) (*http.Client, string) {
	return samStack.findClient(request.URL.String()).sendRequestHttp(request)
}

func (samStack *samList) findClient(request string) *samHttp {
	found := false
	var c samHttp
	for index, client := range samStack.listOfClients {
		log.Println("Checking client requests", index+1)
		log.Println("of", len(samStack.listOfClients))
		if client.hostCheck(request) {
			samStack.Log("Client pipework for %s found.", request)
			samStack.Log("Request sent")
			found = true
			return &client
		}
	}
	if !found {
		samStack.Log("Client pipework for %s not found: Creating.", request)
		samStack.createClient(request)
		for index, client := range samStack.listOfClients {
			log.Println("Checking client requests", index+1)
			log.Println("of", len(samStack.listOfClients))
			if client.hostCheck(request) {
				samStack.Log("Client pipework for %s found.", request)
				c = client
			}
		}
	}
	return &c
}

func (samStack *samList) copyRequest(request *http.Request, response *http.Response, directory string) *http.Response {
	return samStack.findClient(request.URL.String()).copyRequestHttp(request, response, directory)
}

func (samStack *samList) readRequest() {
	samStack.Log("Reading requests:")
	for samStack.sendScan.Scan() {
		if samStack.sendScan.Text() != "" {
			go samStack.sendClientRequest(samStack.sendScan.Text())
		}
	}
}

func (samStack *samList) writeResponses() {
	samStack.Log("Writing responses:")
	for i, client := range samStack.listOfClients {
		log.Println("Checking for responses: %s", i+1)
		log.Println("of: ", len(samStack.listOfClients))
		//b :=
		if client.printResponse() != "" {
			go samStack.writeRecieved(client.printResponse())
		}
		//if b == true {
		//break
		//}
	}
}

func (samStack *samList) responsify(input string) io.Reader {
	tmp := strings.NewReader(input)
	samStack.Log("Responsifying string:")
	return tmp
}

func (samStack *samList) writeRecieved(response string) bool {
	b := false
	if response != "" {
		samStack.Log("Got response:")
		io.Copy(samStack.recvPipe, samStack.responsify(response))
		b = true
	}
	return b
}

func (samStack *samList) readDelete() bool {
	samStack.Log("Managing pipes:")
	for samStack.delScan.Scan() {
		if samStack.delScan.Text() == "y" || samStack.delScan.Text() == "Y" {
			defer samStack.cleanupClient()
			return true
		} else {
			return false
		}
	}
	return false
}

func (samStack *samList) cleanupClient() {
	samStack.sendPipe.Close()
	samStack.recvPipe.Close()
	for _, client := range samStack.listOfClients {
		client.cleanupClient()
	}
	samStack.delPipe.Close()
	os.RemoveAll(filepath.Join(connectionDirectory, "parent"))
}

func (samStack *samList) Log(msg ...string) {
	if verbose {
		log.Println("LOG: ", msg)
	}
}

func (samStack *samList) Warn(err error, errmsg string, msg ...string) bool {
	log.Println(msg)
	if err != nil {
		log.Println("WARN: ", err)
		return false
	}
	samStack.err = nil
	return true
}

func (samStack *samList) Fatal(err error, errmsg string, msg ...string) {
	log.Println(msg)
	if err != nil {
		defer samStack.cleanupClient()
		log.Fatal("FATAL: ", err)
		samStack.err = err
	}
}

func createSamList(samAddr string, samPort string, initAddress string) *samList {
	var samStack samList
	log.Println("Generating parent proxy structure.")
	samStack.up = false
	log.Println("Parent proxy set to down.")
	samStack.createSamList(samAddr, samPort)
	log.Println("SAM list created")
	if initAddress != "" {
		samStack.sendPipe.WriteString(initAddress + "\n")
	}
	return &samStack
}
