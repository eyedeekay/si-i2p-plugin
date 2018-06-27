package dii2pmain

import (
	"bufio"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

import (
	"github.com/eyedeekay/si-i2p-plugin/src/errors"
	"github.com/eyedeekay/si-i2p-plugin/src/helpers"
)

//SamList is a manager which guarantee's unique destinations for websites
//retrieved over the SAM bridge
type SamList struct {
	listOfClients []SamHTTP
	samAddrString string
	samPortString string
	keepAlives    bool
	err           error
	c             bool
	up            bool
	dir           string

	timeoutTime            int
	lifeTime               int
	tunnelLength           int
	inboundQuantity        int
	outboundQuantity       int
	inboundBackupQuantity  int
	outboundBackupQuantity int
	idleConns              int

	lastAddress string

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
	dii2phelper.SetupFolder(samStack.dir)

	samStack.sendPath, samStack.sendPipe, samStack.err = dii2phelper.SetupFiFo(filepath.Join(dii2phelper.ConnectionDirectory, samStack.dir), "send")
	if samStack.c, samStack.err = dii2perrs.Fatal(samStack.err, "sam-list.go Pipe setup error", "sam-list.go Pipe setup"); samStack.c {
		samStack.sendScan, samStack.err = dii2phelper.SetupScanner(filepath.Join(dii2phelper.ConnectionDirectory, samStack.dir), "send", samStack.sendPipe)
		if samStack.c, samStack.err = dii2perrs.Fatal(samStack.err, "sam-list.go Scanner setup Error:", "sam-list.go Scanner set up successfully."); !samStack.c {
			samStack.CleanupClient()
		}
	}

	samStack.recvPath, samStack.recvPipe, samStack.err = dii2phelper.SetupFiFo(filepath.Join(dii2phelper.ConnectionDirectory, samStack.dir), "recv")
	if samStack.c, samStack.err = dii2perrs.Fatal(samStack.err, "sam-list.go Pipe setup error", "sam-list.go Pipe setup"); samStack.c {
		samStack.recvPipe.WriteString("")
	}

	samStack.delPath, samStack.delPipe, samStack.err = dii2phelper.SetupFiFo(filepath.Join(dii2phelper.ConnectionDirectory, samStack.dir), "del")
	if samStack.c, samStack.err = dii2perrs.Fatal(samStack.err, "sam-list.go Pipe setup error", "sam-list.go Pipe setup"); samStack.c {
		samStack.delScan, samStack.err = dii2phelper.SetupScanner(filepath.Join(dii2phelper.ConnectionDirectory, samStack.dir), "del", samStack.delPipe)
		if samStack.c, samStack.err = dii2perrs.Fatal(samStack.err, "sam-list.go Scanner setup Error:", "sam-list.go Scanner set up successfully."); !samStack.c {
			samStack.CleanupClient()
		}
	}

	samStack.up = true
}

func (samStack *SamList) createClient(request string) {
	dii2perrs.Log("sam-list.go Appending client to SAM stack.")
	samStack.listOfClients = append(samStack.listOfClients,
		newSamHTTP(samStack.samAddrString,
			samStack.samPortString,
			request,
			samStack.timeoutTime,
			samStack.lifeTime,
			samStack.keepAlives,
			samStack.tunnelLength,
			samStack.inboundQuantity,
			samStack.outboundQuantity,
			samStack.idleConns,
			samStack.inboundBackupQuantity,
			samStack.outboundBackupQuantity,
		),
	)
}

func (samStack *SamList) createClientHTTP(request *http.Request) {
	dii2perrs.Log("sam-list.go Appending client to SAM stack.")
	samStack.listOfClients = append(samStack.listOfClients,
		newSamHTTPHTTP(samStack.samAddrString,
			samStack.samPortString,
			request,
			samStack.timeoutTime,
			samStack.lifeTime,
			samStack.keepAlives,
			samStack.tunnelLength,
			samStack.inboundQuantity,
			samStack.outboundQuantity,
			samStack.idleConns,
			samStack.inboundBackupQuantity,
			samStack.outboundBackupQuantity,
		),
	)
}

func (samStack *SamList) createSamList() {
	if !samStack.up {
		samStack.initPipes()
		dii2perrs.Log("sam-list.go Parent proxy pipes initialized. Parent proxy set to up.")
	}
}

func (samStack *SamList) sendClientRequest(request string) {
	client := samStack.findClient(request)
	if client != nil {
		client.sendRequest(request)
	}
}

func (samStack *SamList) sendClientRequestHTTP(request *http.Request) (*http.Client, string) {
	client := samStack.findClient(request.URL.String())
	if client != nil {
		return client.sendRequestHTTP(request)
	}
	return nil, "nil client"
}

func (samStack *SamList) hostCheck(request string) (bool, *SamHTTP) {
	if !dii2phelper.CheckURLType(request) {
		return false, nil
	}
	for index, client := range samStack.listOfClients {
		dii2perrs.Log("sam-list.go Checking client requests", strconv.Itoa(index+1), client.host)
		dii2perrs.Log("sam-list.go of", strconv.Itoa(len(samStack.listOfClients)))
		if client.hostCheck(request) {
			dii2perrs.Log("sam-list.go Client pipework for", request, "found.", client.host, "at", strconv.Itoa(index+1))
			return true, &client
		}
	}
	return false, nil
}

func (samStack *SamList) lifetimeCheck(request string) bool {
	if !dii2phelper.CheckURLType(request) {
		return false
	}
	for index, client := range samStack.listOfClients {
		if client.lifetimeCheck(request) {
			dii2perrs.Log("sam-list.go Removing inactive client after", samStack.lifeTime, "minutes.")
			samStack.listOfClients = samStack.deleteClient(samStack.listOfClients, index)
			return true
		}
	}
	return false
}

func (samStack *SamList) deleteClient(s []SamHTTP, index int) []SamHTTP {
	s[index].CleanupClient()
	return append(s[:index], s[index+1])
}

func (samStack *SamList) findClient(request string) *SamHTTP {
	if !dii2phelper.CheckURLType(request) {
		return nil
	}
	if !samStack.lifetimeCheck(request) {
		found, c := samStack.hostCheck(request)
		if found {
			return c
		}
	}
	dii2perrs.Log("sam-list.go Client pipework for", request, "not found: Creating.")
	samStack.createClient(request)
	found, c := samStack.hostCheck(request)
	if found {
		return c
	}
	return c
}

func (samStack *SamList) copyRequest(request *http.Request, response *http.Response, directory string) *http.Response {
	return samStack.findClient(request.URL.String()).copyRequestHTTP(request, response, directory)
}

//ReadRequest checks the pipes for new URLs to request
func (samStack *SamList) ReadRequest() {
	dii2perrs.Log("sam-list.go Reading requests:")
	for samStack.sendScan.Scan() {
		if samStack.sendScan.Text() != "" {
			go samStack.sendClientRequest(samStack.sendScan.Text())
		}
	}
	dii2phelper.ClearFile(filepath.Join(dii2phelper.ConnectionDirectory, samStack.dir), "send")
}

//WriteResponses writes the responses to the pipes
func (samStack *SamList) WriteResponses() {
	dii2perrs.Log("sam-list.go Writing responses:")
	for i, client := range samStack.listOfClients {
		dii2perrs.Log("sam-list.go Checking for responses: ", i+1)
		dii2perrs.Log("sam-list.go of: ", len(samStack.listOfClients))
		if client.printResponse() != "" {
			go samStack.writeRecieved(client.printResponse())
		}
	}
}

func (samStack *SamList) responsify(input string) io.ReadCloser {
	tmp := ioutil.NopCloser(strings.NewReader(input))
	defer tmp.Close()
	dii2perrs.Log("sam-list.go Responsifying string:")
	return tmp
}

func (samStack *SamList) writeRecieved(response string) bool {
	b := false
	if response != "" {
		dii2perrs.Log("sam-list.go Got response:")
		io.Copy(samStack.recvPipe, samStack.responsify(response))
		b = true
	}
	return b
}

//ReadDelete closes the SamList
func (samStack *SamList) ReadDelete() bool {
	dii2perrs.Log("sam-list.go Managing pipes:")
	for samStack.delScan.Scan() {
		if samStack.delScan.Text() == "y" || samStack.delScan.Text() == "Y" {
			defer samStack.CleanupClient()
			return true
		}
		return false
	}
	for _, client := range samStack.listOfClients {
		client.readDelete()
	}
	dii2phelper.ClearFile(filepath.Join(dii2phelper.ConnectionDirectory, samStack.dir), "del")
	return false
}

//CleanupClient tears down all SamList members
func (samStack *SamList) CleanupClient() {
	samStack.sendPipe.Close()
	samStack.recvPipe.Close()
	for _, client := range samStack.listOfClients {
		client.CleanupClient()
	}
	samStack.delPipe.Close()
	os.RemoveAll(filepath.Join(dii2phelper.ConnectionDirectory, samStack.dir))
}

//CreateSamList initializes a SamList
func CreateSamList(opts ...func(*SamList) error) (*SamList, error) {
	var samStack SamList
	samStack.dir = "parent"
	samStack.up = false
	samStack.timeoutTime = 6
	samStack.lifeTime = 2
	samStack.tunnelLength = 3
	samStack.inboundQuantity = 15
	samStack.outboundQuantity = 15
	samStack.inboundBackupQuantity = 3
	samStack.outboundBackupQuantity = 3
	samStack.idleConns = 4
	samStack.lastAddress = ""
	dii2perrs.Log("sam-list.go Parent proxy set to down.")
	dii2perrs.Log("sam-list.go Generating parent proxy structure.")
	for _, o := range opts {
		if err := o(&samStack); err != nil {
			return nil, err
		}
	}
	samStack.createSamList()
	dii2perrs.Log("sam-list.go SAM list created")
	if samStack.lastAddress != "" {
		samStack.sendPipe.WriteString(samStack.lastAddress + "\n")
	}
	return &samStack, nil
}
