package dii2p

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type SamList struct {
	listOfClients []SamHttp
	samAddrString string
	samPortString string
	keepAlives    bool
	err           error
	c             bool
	up            bool
	dir           string
	timeoutTime   int

	sendPath string
	sendPipe *os.File
	sendScan *bufio.Scanner

	recvPath string
	recvPipe *os.File

	delPath string
	delPipe *os.File
	delScan *bufio.Scanner
}

func (samStack *SamList) initPipes() {
	setupFolder(samStack.dir)

	samStack.sendPath, samStack.sendPipe, samStack.err = setupFiFo(filepath.Join(connectionDirectory, samStack.dir), "send")
	if samStack.c, samStack.err = Fatal(samStack.err, "sam-list.go Pipe setup error", "sam-list.go Pipe setup"); samStack.c {
		samStack.sendScan, samStack.err = setupScanner(filepath.Join(connectionDirectory, samStack.dir), "send", samStack.sendPipe)
		if samStack.c, samStack.err = Fatal(samStack.err, "sam-list.go Scanner setup Error:", "sam-list.go Scanner set up successfully."); !samStack.c {
			samStack.CleanupClient()
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
			samStack.CleanupClient()
		}
	}

	samStack.up = true
}

func (samStack *SamList) createClient(request string) {
	Log("sam-list.go Appending client to SAM stack.")
	samStack.listOfClients = append(samStack.listOfClients, newSamHttp(samStack.samAddrString, samStack.samPortString, request, samStack.timeoutTime, samStack.keepAlives))
}

func (samStack *SamList) createClientHttp(request *http.Request) {
	Log("sam-list.go Appending client to SAM stack.")
	samStack.listOfClients = append(samStack.listOfClients, newSamHttpHttp(samStack.samAddrString, samStack.samPortString, request, samStack.timeoutTime, samStack.keepAlives))
}

func (samStack *SamList) createSamList(samAddrString string, samPortString string) {
	samStack.samAddrString = samAddrString
	samStack.samPortString = samPortString
	Log("sam-list.go Established SAM connection")
	if !samStack.up {
		samStack.initPipes()
		Log("sam-list.go Parent proxy pipes initialized. Parent proxy set to up.")
	}
}

func (samStack *SamList) sendClientRequest(request string) {
	client := samStack.findClient(request)
	if client != nil {
		client.sendRequest(request)
	}
}

func (samStack *SamList) sendClientRequestHttp(request *http.Request) (*http.Client, string) {
	client := samStack.findClient(request.URL.String())
	if client != nil {
		return client.sendRequestHttp(request)
	} else {
		return nil, "nil client"
	}
}

func (samStack *SamList) findClient(request string) *SamHttp {
	found := false
	var c SamHttp
	if !samStack.checkURLType(request) {
		return nil
	}
	for index, client := range samStack.listOfClients {
		Log("sam-list.go Checking client requests", strconv.Itoa(index+1), client.host)
		Log("sam-list.go of", strconv.Itoa(len(samStack.listOfClients)))
		if client.hostCheck(request) {
			Log("sam-list.go Client pipework for", request, "found.", client.host, "at", strconv.Itoa(index+1))
			found = true
			c = client
			return &c
		}
	}
	if !found {
		Log("sam-list.go Client pipework for %s not found: Creating.", request)
		samStack.createClient(request)
		for index, client := range samStack.listOfClients {
			Log("sam-list.go Checking client requests", strconv.Itoa(index+1), client.host)
			Log("sam-list.go of", strconv.Itoa(len(samStack.listOfClients)))
			if client.hostCheck(request) {
				Log("sam-list.go Client pipework for", request, "found.", client.host, "at", strconv.Itoa(index+1))
				c = client
				return &c
			}
		}
	}
	return &c
}

func (samStack *SamList) copyRequest(request *http.Request, response *http.Response, directory string) *http.Response {
	return samStack.findClient(request.URL.String()).copyRequestHttp(request, response, directory)
}

//export ReadRequest
func (samStack *SamList) ReadRequest() {
	Log("sam-list.go Reading requests:")
	for samStack.sendScan.Scan() {
		if samStack.sendScan.Text() != "" {
			go samStack.sendClientRequest(samStack.sendScan.Text())
		}
	}
	clearFile(filepath.Join(connectionDirectory, samStack.dir), "send")
}

//export WriteResponses
func (samStack *SamList) WriteResponses() {
	Log("sam-list.go Writing responses:")
	for i, client := range samStack.listOfClients {
		log.Println("sam-list.go Checking for responses: %s", i+1)
		log.Println("sam-list.go of: ", len(samStack.listOfClients))
		if client.printResponse() != "" {
			go samStack.writeRecieved(client.printResponse())
		}
	}
}

func (samStack *SamList) responsify(input string) io.Reader {
	tmp := strings.NewReader(input)
	Log("sam-list.go Responsifying string:")
	return tmp
}

func (samStack *SamList) writeRecieved(response string) bool {
	b := false
	if response != "" {
		Log("sam-list.go Got response:")
		io.Copy(samStack.recvPipe, samStack.responsify(response))
		b = true
	}
	return b
}

//export ReadDelete
func (samStack *SamList) ReadDelete() bool {
	Log("sam-list.go Managing pipes:")
	for samStack.delScan.Scan() {
		if samStack.delScan.Text() == "y" || samStack.delScan.Text() == "Y" {
			defer samStack.CleanupClient()
			return true
		} else {
			return false
		}
	}
	clearFile(filepath.Join(connectionDirectory, samStack.dir), "del")
	return false
}

//export CleanupClient
func (samStack *SamList) CleanupClient() {
	samStack.sendPipe.Close()
	samStack.recvPipe.Close()
	for _, client := range samStack.listOfClients {
		client.CleanupClient()
	}
	samStack.delPipe.Close()
	os.RemoveAll(filepath.Join(connectionDirectory, samStack.dir))
}

func (samStack *SamList) checkURLType(request string) bool {

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

//export CreateSamList
func CreateSamList(samAddr, samPort, initAddress string, timeoutTime int, keepAlives bool) *SamList {
	var samStack SamList
	samStack.timeoutTime = timeoutTime
	samStack.dir = "parent"
	Log("sam-list.go Generating parent proxy structure.")
	samStack.up = false
	samStack.keepAlives = keepAlives
	Log("sam-list.go Parent proxy set to down.")
	samStack.createSamList(samAddr, samPort)
	Log("sam-list.go SAM list created")
	if initAddress != "" {
		samStack.sendPipe.WriteString(initAddress + "\n")
	}
	return &samStack
}
