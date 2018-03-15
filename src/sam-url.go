package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type samUrl struct {
	err          error
	c            bool
	subDirectory string

	recvPath string
	recvFile *os.File

	timePath string
	timeFile *os.File

	delPath string
	delPipe *os.File
	//delBuff bufio.Reader
	delScan *bufio.Scanner
}

func (subUrl *samUrl) initPipes() {
    checkFolder(filepath.Join(connectionDirectory, subUrl.subDirectory))

    subUrl.recvPath, subUrl.recvFile, subUrl.err = setupFiFo(filepath.Join(connectionDirectory, subUrl.subDirectory), "recv")
    if subUrl.c, subUrl.err = Fatal(subUrl.err, "Pipe setup error", "Pipe setup"); subUrl.c {
        subUrl.recvFile.WriteString("")
    }

    subUrl.timePath, subUrl.timeFile, subUrl.err = setupFiFo(filepath.Join(connectionDirectory, subUrl.subDirectory), "time")
    if subUrl.c, subUrl.err = Fatal(subUrl.err, "Pipe setup error", "Pipe setup"); subUrl.c {
        subUrl.timeFile.WriteString("")
    }

    subUrl.delPath, subUrl.delPipe, subUrl.err = setupFiFo(filepath.Join(connectionDirectory, subUrl.subDirectory), "del")
    if subUrl.c, subUrl.err = Fatal(subUrl.err, "Pipe setup error", "Pipe setup"); subUrl.c {
        subUrl.delScan, subUrl.err = setupScanner(filepath.Join(connectionDirectory, subUrl.subDirectory), "del", subUrl.delPipe)
        if subUrl.c, subUrl.err = Fatal(subUrl.err, "Scanner setup Error:", "Scanner set up successfully."); !subUrl.c {
            subUrl.cleanupDirectory()
        }
    }

}

func (subUrl *samUrl) createDirectory(requestdir string) {
	subUrl.subDirectory = subUrl.dirSet(requestdir)
	subUrl.initPipes()
}

func (subUrl *samUrl) scannerText() (string, error) {
	d, err := ioutil.ReadFile(subUrl.recvPath)
	if subUrl.c, subUrl.err = Fatal(err, "Scanner error", "Scanning recv"); subUrl.c {
        return "", subUrl.err
    }
	s := string(d)
	if s != "" {
		Log("Read file", s)
		return s, err
	}
	return "", err
}

func (subUrl *samUrl) dirSet(requestdir string) string {
	Log("Requesting directory: ", requestdir+"/")
	d1 := requestdir
	d2 := strings.Replace(d1, "//", "/", -1)
	return d2
}

func (subUrl *samUrl) checkDirectory(directory string) bool {
	b := false
	if directory == subUrl.subDirectory {
		Log("Directory / ", directory+" : equals : "+subUrl.subDirectory)
		b = true
	} else {
		Log("Directory / ", directory+" : does not equal : "+subUrl.subDirectory)
	}
	return b
}

func (subUrl *samUrl) copyDirectory(response *http.Response, directory string) bool {
	b := false
	if subUrl.checkDirectory(directory) {
		if response != nil {
			log.Println("Response Status ", response.StatusCode)
			if response.StatusCode == http.StatusOK {
				Log("Setting file in cache")
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
				Log("Setting file in cache")
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
	if subUrl.c, subUrl.err = Warn(err, "Response Write Error", "Writing responses"); subUrl.c {
		Log("Writing files.")
		subUrl.recvFile.Write(body)
		Log("Retrieval time: ", time.Now().String())
		subUrl.timeFile.WriteString(time.Now().String())
	}
}

func (subUrl *samUrl) dealResponseHttp(request *http.Request, response *http.Response) *http.Response {
	defer response.Body.Close()
	header := response.Header
	trailer := response.Trailer
	status := response.Status
	statusCode := response.StatusCode
	proto := response.Proto
	protoMajor := response.ProtoMajor
	protoMinor := response.ProtoMinor
	body, err := ioutil.ReadAll(response.Body)
	if subUrl.c, subUrl.err = Warn(err, "Response Read Error", "Reading response from proxy"); subUrl.c {
		Log("Writing files.")
		subUrl.recvFile.Write(body)
		r := &http.Response{
			Status:        status,
			StatusCode:    statusCode,
			Proto:         proto,
			ProtoMajor:    protoMajor,
			ProtoMinor:    protoMinor,
			Body:          ioutil.NopCloser(bytes.NewBuffer(body)),
			ContentLength: int64(len(body)),
			Request:       request,
			Header:        header,
			Trailer:       trailer,
		}
		Log("Retrieval time: ", time.Now().String())
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

func (subUrl *samUrl) readDelete() bool {
	Log("Managing pipes:")
	for subUrl.delScan.Scan() {
		if subUrl.delScan.Text() == "y" || subUrl.delScan.Text() == "Y" {
			defer subUrl.cleanupDirectory()
			return true
		} else {
			return false
		}
	}
	return false
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
