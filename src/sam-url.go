package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type samUrl struct {
	err          error
	subDirectory string

	recvPath string
	recvFile *os.File

	timePath string
	timeFile *os.File

	delPath string
	delPipe *os.File
	delBuff bufio.Reader
}

func (subUrl *samUrl) initPipes() {
	pathConnectionExists, pathErr := exists(filepath.Join(connectionDirectory, subUrl.subDirectory))
	subUrl.Log("Starting Directory Creation", filepath.Join(connectionDirectory, subUrl.subDirectory))
	subUrl.Fatal(pathErr, "Directory Creation Error", "Directory Check", filepath.Join(connectionDirectory, subUrl.subDirectory))
	if !pathConnectionExists {
		subUrl.Log("Creating a connection:", subUrl.subDirectory)
		os.MkdirAll(filepath.Join(connectionDirectory, subUrl.subDirectory), 0755)
	}

	subUrl.recvPath = filepath.Join(connectionDirectory, subUrl.subDirectory, "recv")
	pathRecvExists, recvPathErr := exists(subUrl.recvPath)
	subUrl.Fatal(recvPathErr, "File Check Error", "Checking file", subUrl.recvPath)
	if !pathRecvExists {
		subUrl.recvFile, subUrl.err = os.Create(subUrl.recvPath)
		//subUrl.Log("Preparing to create File:", subUrl.recvPath)
		subUrl.Fatal(subUrl.err, "File Creation Error", "Creating file", subUrl.recvPath)
		subUrl.Log("checking for problems...")
		subUrl.Log("Opening the File...")
		subUrl.recvFile, subUrl.err = os.OpenFile(subUrl.recvPath, os.O_RDWR|os.O_CREATE, 0644)
		subUrl.Log("Created a File for recieving responses:", subUrl.recvPath)
	}

	subUrl.timePath = filepath.Join(connectionDirectory, subUrl.subDirectory, "time")
	pathTimeExists, recvTimeErr := exists(subUrl.timePath)
	subUrl.Fatal(recvTimeErr, "File Check Error", "Checking file", subUrl.timePath)
	if !pathTimeExists {
		subUrl.timeFile, subUrl.err = os.Create(subUrl.timePath)
		//subUrl.Log("Preparing to create File:", subUrl.timePath)
		subUrl.Fatal(subUrl.err, "File Creation Error", "Creating file", subUrl.timePath)
		subUrl.Log("checking for problems...")
		subUrl.Log("Opening the File...")
		subUrl.timeFile, subUrl.err = os.OpenFile(subUrl.timePath, os.O_RDWR|os.O_CREATE, 0644)
		subUrl.Log("Created a File for timing responses:", subUrl.timePath)
	}

	subUrl.delPath = filepath.Join(connectionDirectory, subUrl.subDirectory, "del")
	pathDelExists, delPathErr := exists(subUrl.delPath)
	subUrl.Fatal(delPathErr, "File Check Error", "Checking file", subUrl.delPath)
	if !pathDelExists {
		err := syscall.Mkfifo(subUrl.delPath, 0755)
		//subUrl.Log("Preparing to create Pipe:", subUrl.delPath)
		subUrl.Fatal(err, "File Creation Error", "Creating file", subUrl.delPath)
		subUrl.Log("checking for problems...")
		subUrl.delPipe, err = os.OpenFile(subUrl.delPath, os.O_RDWR|os.O_CREATE, 0755)
		subUrl.Log("Opening the Named Pipe as a File...")
		subUrl.delBuff = *bufio.NewReader(subUrl.delPipe)
		subUrl.Log("Opening the Named Pipe as a Buffer...")
		subUrl.Log("Created a named Pipe for closing the connection:", subUrl.delPath)
	}
}

func (subUrl *samUrl) createDirectory(requestdir string) {
	subUrl.subDirectory = subUrl.dirSet(requestdir)
	subUrl.initPipes()
}

func (subUrl *samUrl) scannerText() (string, error) {
	d, err := ioutil.ReadFile(subUrl.recvPath)
	subUrl.Fatal(err, "Scanner error", "Scanning recv")
	s := string(d)
	if s != "" {
		subUrl.Log("Read file", s)
		return s, err
	}
	return "", err
}

func (subUrl *samUrl) dirSet(requestdir string) string {
	subUrl.Log("Requesting directory: ", requestdir+"/")
	d1 := requestdir
	d2 := strings.Replace(d1, "//", "/", -1)
	return d2
}

func (subUrl *samUrl) checkDirectory(directory string) bool {
	b := false
	if directory == subUrl.subDirectory {
		subUrl.Log("Directory / ", directory+" : equals : "+subUrl.subDirectory)
		b = true
	} else {
		subUrl.Log("Directory / ", directory+" : does not equal : "+subUrl.subDirectory)
	}
	return b
}

func (subUrl *samUrl) copyDirectory(response *http.Response, directory string) bool {
	b := false
	if subUrl.checkDirectory(directory) {
		if response != nil {
			log.Println("Response Status ", response.StatusCode)
			if response.StatusCode == http.StatusOK {
				subUrl.Log("Setting file in cache")
				subUrl.dealResponse(response)
			}
		}
		b = true
	}
	return b
}

func (subUrl *samUrl) copyDirectoryHttp(request *http.Request, response *http.Response, directory string) *http.Response {
	if subUrl.checkDirectory(directory) {
		if response != nil {
			log.Println("Response Status ", response.StatusCode)
			if response.StatusCode == http.StatusOK {
				subUrl.Log("Setting file in cache")
				resp := subUrl.dealResponseHttp(request, response)
				return resp
			}
		}
	}
	return response
}

func (subUrl *samUrl) dealResponse(response *http.Response) {
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if subUrl.Warn(err, "Response Write Error", "Writing responses") {
		subUrl.Log("Writing files.")
		subUrl.recvFile.Write(body)
		subUrl.Log("Retrieval time: ", time.Now().String())
		subUrl.timeFile.WriteString(time.Now().String())
	}
}

func (subUrl *samUrl) dealResponseHttp(request *http.Request, response *http.Response) *http.Response {
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if subUrl.Warn(err, "Response Read Error", "Reading response from proxy") {
		subUrl.Log("Writing files.")
		subUrl.recvFile.Write(body)
		r := &http.Response{
			Status:        "200 OK",
			StatusCode:    200,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			Body:          ioutil.NopCloser(bytes.NewBuffer(body)),
			ContentLength: int64(len(body)),
			Request:       request,
			Header:        request.Header,
		}
		subUrl.Log("Retrieval time: ", time.Now().String())
		subUrl.timeFile.WriteString(time.Now().String())
		return r
	}
	return response
}

func (subUrl *samUrl) cleanupDirectory() {
	subUrl.recvFile.Close()
	subUrl.timeFile.Close()
	subUrl.delPipe.Close()
	os.RemoveAll(filepath.Join(connectionDirectory, subUrl.subDirectory))
}

func (subUrl *samUrl) readDelete() int {
	line, _, err := subUrl.delBuff.ReadLine()
	subUrl.Fatal(err, "Reading from exit pipe error", "Reading exit pipe")
	n := len(line)
	subUrl.Log("Reading n bytes from exit pipe:", strconv.Itoa(n))
	if n < 0 {
		//subUrl.Log("Something wierd happened with :", line)
		subUrl.Log("end determined at index :", strconv.Itoa(n))
		return n
	} else {
		s := string(line[:n])
		if s == "y" {
			subUrl.Log("Deleting connection: %s", subUrl.subDirectory)
			defer subUrl.cleanupDirectory()
			return n
		} else {
			return n
		}
	}
}

func (subUrl *samUrl) Log(msg ...string) {
	if verbose {
		log.Println("LOG: ", msg)
	}
}

func (subUrl *samUrl) Warn(err error, errmsg string, msg ...string) bool {
	log.Println(msg)
	if err != nil {
		log.Println("WARN: ", err)
		return false
	}
	subUrl.err = nil
	return true
}

func (subUrl *samUrl) Fatal(err error, errmsg string, msg ...string) {
	log.Println(msg)
	if err != nil {
		subUrl.err = err
		defer subUrl.cleanupDirectory()
		log.Fatal("FATAL: ", errmsg, err)
	}
}

func newSamUrl(requestdir string) samUrl {
	log.Println("Creating a new cache directory.")
	var subUrl samUrl
	subUrl.createDirectory(requestdir)
	return subUrl
}

func newSamUrlHttp(request *http.Request) samUrl {
	log.Println("Creating a new cache directory.")
	var subUrl samUrl
	log.Println(subUrl.subDirectory)
	subUrl.createDirectory(request.Host + request.URL.Path)
	return subUrl
}
