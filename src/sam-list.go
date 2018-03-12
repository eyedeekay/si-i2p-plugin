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
	samStack.Fatal(err)
	if !pathConnectionExists {
		log.Println("Creating a connection:", "parent")
		os.Mkdir(filepath.Join(connectionDirectory, "parent"), 0755)
	} else {
		os.RemoveAll(filepath.Join(connectionDirectory, "parent"))
		log.Println("Creating a connection:", "parent")
		os.Mkdir(filepath.Join(connectionDirectory, "parent"), 0755)
	}

	samStack.sendPath = filepath.Join(connectionDirectory, "parent", "send")
	pathSendExists, sendErr := exists(samStack.sendPath)
	samStack.Fatal(sendErr)
	if !pathSendExists {
		samStack.err = syscall.Mkfifo(samStack.sendPath, 0755)
		log.Println("Preparing to create Pipe:", samStack.sendPath)
		samStack.Fatal(samStack.err)
		log.Println("checking for problems...")
		samStack.sendPipe, samStack.err = os.OpenFile(samStack.sendPath, os.O_RDWR|os.O_CREATE, 0755)
		log.Println("Opening the Named Pipe as a Scanner...")
		samStack.sendScan = *bufio.NewScanner(samStack.sendPipe)
		samStack.sendScan.Split(bufio.ScanLines)
		log.Println("Opening the Named Pipe as a Scanner...")
		log.Println("Created a named Pipe for sending requests:", samStack.sendPath)
	}

	samStack.recvPath = filepath.Join(connectionDirectory, "parent", "recv")
	pathRecvExists, recvErr := exists(samStack.recvPath)
	samStack.Fatal(recvErr)
	if !pathRecvExists {
		samStack.err = syscall.Mkfifo(samStack.recvPath, 0755)
		log.Println("Preparing to create Pipe:", samStack.recvPath)
		samStack.Fatal(samStack.err)
		log.Println("checking for problems...")
		samStack.recvPipe, samStack.err = os.OpenFile(samStack.recvPath, os.O_RDWR|os.O_CREATE, 0755)
		samStack.recvPipe.WriteString("")
		log.Println("Created a named Pipe for recieving responses:", samStack.recvPath)
	}

	samStack.delPath = filepath.Join(connectionDirectory, "parent", "del")
	pathDelExists, delErr := exists(samStack.delPath)
	samStack.Fatal(delErr)
	if !pathDelExists {
		samStack.err = syscall.Mkfifo(samStack.delPath, 0755)
		log.Println("Preparing to create Pipe:", samStack.delPath)
		samStack.Fatal(samStack.err)
		log.Println("checking for problems...")
		samStack.delPipe, samStack.err = os.OpenFile(samStack.delPath, os.O_RDWR|os.O_CREATE, 0755)
		samStack.recvPipe.WriteString("")
		log.Println("Opening the Named Pipe as a File...")
		samStack.delScan = *bufio.NewScanner(samStack.delPipe)
		samStack.delScan.Split(bufio.ScanLines)
		log.Println("Opening the Named Pipe as a Scanner...")
		log.Println("Created a named Pipe for closing the connection:", samStack.delPath)
	}
	samStack.up = true
}

func (samStack *samList) createClient(request string) {
	log.Println("Appending client to SAM stack.")
	//samStack.listOfClients = append(samStack.listOfClients, newSamHttp(samStack.samAddrString, samStack.samPortString, samStack.samBridgeClient, request))
	samStack.listOfClients = append(samStack.listOfClients, newSamHttp(samStack.samAddrString, samStack.samPortString, request))
}

func (samStack *samList) createClientHttp(request *http.Request) {
	log.Println("Appending client to SAM stack.")
	samStack.listOfClients = append(samStack.listOfClients, newSamHttpHttp(samStack.samAddrString, samStack.samPortString, request))
}

func (samStack *samList) createSamList(samAddrString string, samPortString string) {
	samStack.samAddrString = samAddrString
	samStack.samPortString = samPortString
	log.Println("Requesting a new SAM-based http client")
	samStack.Fatal(samStack.err)
	log.Println("Established SAM connection")
	if !samStack.up {
		samStack.initPipes()
		log.Println("Parent proxy pipes initialized. Parent proxy set to up.")
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
			log.Println("Client pipework for %s found.", request)
			log.Println("Request sent")
			found = true
			return &client
		}
	}
	if !found {
		log.Println("Client pipework for %s not found: Creating.", request)
		samStack.createClient(request)
		for index, client := range samStack.listOfClients {
			log.Println("Checking client requests", index+1)
			log.Println("of", len(samStack.listOfClients))
			if client.hostCheck(request) {
				log.Println("Client pipework for %s found.", request)
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
	log.Println("Reading requests:")
	for samStack.sendScan.Scan() {
		if samStack.sendScan.Text() != "" {
			go samStack.sendClientRequest(samStack.sendScan.Text())
		}
	}
}

func (samStack *samList) writeResponses() {
	log.Println("Writing responses:")
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
	log.Println("Responsifying string:")
	return tmp
}

func (samStack *samList) writeRecieved(response string) bool {
	b := false
	if response != "" {
		log.Println("Got response:")
		io.Copy(samStack.recvPipe, samStack.responsify(response))
		b = true
	}
	return b
}

func (samStack *samList) readDelete() bool {
	log.Println("Managing pipes:")
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

func (samStack *samList) Warn(err error) bool {
	if err != nil {
		log.Println("WARN: ", err)
		samStack.err = err
		return true
	}
	return false
}

func (samStack *samList) Fatal(err error) bool {
	if err != nil {
		defer samStack.cleanupClient()
		log.Fatal("FATAL: ", err)
		samStack.err = err
		return true
	}
	return false
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
