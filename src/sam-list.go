package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type samList struct {
	listOfClients []samHttp
	samAddrString string
	samPortString string
	err           error
    c             bool
	up            bool
	dir           string

	sendPath string
	sendPipe *os.File
	sendScan *bufio.Scanner

	recvPath string
	recvPipe *os.File

	delPath string
	delPipe *os.File
	delScan *bufio.Scanner
}

func (samStack *samList) initPipes() {
    setupFolder(samStack.dir)

    samStack.sendPath, samStack.sendPipe, samStack.err = setupFiFo(filepath.Join(connectionDirectory, samStack.dir), "send")
    if samStack.c, samStack.err = Fatal(samStack.err, "sam-list.go Pipe setup error", "sam-list.go Pipe setup"); samStack.c {
        samStack.sendScan, samStack.err = setupScanner(filepath.Join(connectionDirectory, samStack.dir), "send", samStack.sendPipe)
        if samStack.c, samStack.err = Fatal(samStack.err, "sam-list.go Scanner setup Error:", "sam-list.go Scanner set up successfully."); !samStack.c {
            samStack.cleanupClient()
        }
    }

    samStack.recvPath, samStack.recvPipe, samStack.err = setupFiFo(filepath.Join(connectionDirectory, samStack.dir), "recv")
    if samStack.c, samStack.err = Fatal(samStack.err, "sam-list.go Pipe setup error", "sam-list.go Pipe setup"); samStack.c {
        samStack.recvPipe.WriteString("")
    }

    samStack.delPath, samStack.delPipe, samStack.err = setupFiFo(filepath.Join(connectionDirectory, samStack.dir), "del")
    if samStack.c, samStack.err = Fatal(samStack.err, "sam-list.go Pipe setup error", "sam-list.go Pipe setup"); samStack.c {
        samStack.delScan, samStack.err = setupScanner(filepath.Join(connectionDirectory, samStack.dir), "del", samStack.delPipe)
        if samStack.c, samStack.err = Fatal(samStack.err, "sam-list.go Scanner setup Error:", "sam-list.go Scanner set up successfully."); !samStack.c {
            samStack.cleanupClient()
        }
    }

	samStack.up = true
}

func (samStack *samList) createClient(request string) {
	Log("sam-list.go Appending client to SAM stack.")
	samStack.listOfClients = append(samStack.listOfClients, newSamHttp(samStack.samAddrString, samStack.samPortString, request))
}

func (samStack *samList) createClientHttp(request *http.Request) {
	Log("sam-list.go Appending client to SAM stack.")
	samStack.listOfClients = append(samStack.listOfClients, newSamHttpHttp(samStack.samAddrString, samStack.samPortString, request))
}

func (samStack *samList) createSamList(samAddrString string, samPortString string) {
	samStack.samAddrString = samAddrString
	samStack.samPortString = samPortString
	Log("sam-list.go Established SAM connection")
	if !samStack.up {
		samStack.initPipes()
		Log("sam-list.go Parent proxy pipes initialized. Parent proxy set to up.")
	}
}

func (samStack *samList) sendClientRequest(request string) {
	client := samStack.findClient(request)
	if client != nil {
		client.sendRequest(request)
	}
}

func (samStack *samList) sendClientRequestHttp(request *http.Request) (*http.Client, string) {
	client := samStack.findClient(request.URL.String())
	if client != nil {
		return client.sendRequestHttp(request)
	} else {
		return nil, "nil client"
	}
}

func (samStack *samList) checkURLType(request string) bool {

	Log(request)

	test := strings.Split(request, ".i2p")

	if len(test) < 2 {
		msg := "Non i2p domain detected. Skipping."
		Log(msg) //Outproxy support? Might be cool.
		return false
	} else {
		n := strings.Split(strings.Replace(strings.Replace(test[0], "https://", "", -1), "http://", "", -1), "/")
		if len(n) > 1 {
			msg := "Non i2p domain detected, possible attempt to impersonate i2p domain in path. Skipping."
			Log(msg) //Outproxy support? Might be cool. Riskier here.
			return false
		}
	}
	strings.Contains(request, "http")
	if !strings.Contains(request, "http") {
		if strings.Contains(request, "https") {
			msg := "Dropping https request for now, assumed attempt to get clearnet resource."
			Log(msg)
			return false
		} else {
			msg := "unsupported protocal scheme " + request
			Log(msg)
			return false
		}
	} else {
		return true
	}
}

func (samStack *samList) findClient(request string) *samHttp {
	found := false
	var c samHttp
	if !samStack.checkURLType(request) {
		return nil
	}
	for index, client := range samStack.listOfClients {
		log.Println("sam-list.go Checking client requests", index+1)
		log.Println("sam-list.go of", len(samStack.listOfClients))
		if client.hostCheck(request) {
			Log("sam-list.go Client pipework for %s found.", request)
			Log("sam-list.go Request sent")
			found = true
			return &client
		}
	}
	if !found {
		Log("sam-list.go Client pipework for %s not found: Creating.", request)
		samStack.createClient(request)
		for index, client := range samStack.listOfClients {
			log.Println("sam-list.go Checking client requests", index+1)
			log.Println("sam-list.go of", len(samStack.listOfClients))
			if client.hostCheck(request) {
				Log("sam-list.go Client pipework for %s found.", request)
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
	Log("sam-list.go Reading requests:")
	for samStack.sendScan.Scan() {
		if samStack.sendScan.Text() != "" {
			go samStack.sendClientRequest(samStack.sendScan.Text())
		}
	}
}

func (samStack *samList) writeResponses() {
	Log("sam-list.go Writing responses:")
	for i, client := range samStack.listOfClients {
		log.Println("sam-list.go Checking for responses: %s", i+1)
		log.Println("sam-list.go of: ", len(samStack.listOfClients))
		if client.printResponse() != "" {
			go samStack.writeRecieved(client.printResponse())
		}
	}
}

func (samStack *samList) responsify(input string) io.Reader {
	tmp := strings.NewReader(input)
	Log("sam-list.go Responsifying string:")
	return tmp
}

func (samStack *samList) writeRecieved(response string) bool {
	b := false
	if response != "" {
		Log("sam-list.go Got response:")
		io.Copy(samStack.recvPipe, samStack.responsify(response))
		b = true
	}
	return b
}

func (samStack *samList) readDelete() bool {
	Log("sam-list.go Managing pipes:")
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
	os.RemoveAll(filepath.Join(connectionDirectory, samStack.dir))
}

func createSamList(samAddr string, samPort string, initAddress string) *samList {
	var samStack samList
	samStack.dir = "parent"
	Log("sam-list.go Generating parent proxy structure.")
	samStack.up = false
	Log("sam-list.go Parent proxy set to down.")
	samStack.createSamList(samAddr, samPort)
	Log("sam-list.go SAM list created")
	if initAddress != "" {
		samStack.sendPipe.WriteString(initAddress + "\n")
	}
	return &samStack
}
