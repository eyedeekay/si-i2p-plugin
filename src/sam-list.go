package dii2p

import (
	"bufio"
	"io"
    "io/ioutil"
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

    lastAddress   string

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

func (samStack *SamList) createSamList() {
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
	if !CheckURLType(request) {
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

func (samStack *SamList) responsify(input string) io.ReadCloser {
	tmp := ioutil.NopCloser(strings.NewReader(input))
    defer tmp.Close()
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
    for _, client := range samStack.listOfClients {
        client.ReadDelete()
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

//export CreateSamList
func CreateSamList(opts ...func(*SamList) error) (*SamList, error) {
	var samStack SamList
	samStack.dir = "parent"
	samStack.up = false
	Log("sam-list.go Parent proxy set to down.")
	Log("sam-list.go Generating parent proxy structure.")
	for _, o := range opts {
		if err := o(&samStack); err != nil {
			return nil, err
		}
	}
	samStack.createSamList()
	Log("sam-list.go SAM list created")
	if samStack.lastAddress != "" {
		samStack.sendPipe.WriteString(samStack.lastAddress + "\n")
	}
	return &samStack, nil
}
