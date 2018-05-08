package dii2p

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type samUrl struct {
	err          error
	c            bool
	subDirectory string
	mutex        *sync.Mutex

	recvPath string
	recvFile *os.File

	timePath string
	timeFile *os.File

	delPath string
	delPipe *os.File
	delScan *bufio.Scanner
}

func (subUrl *samUrl) initPipes() {
	checkFolder(filepath.Join(connectionDirectory, subUrl.subDirectory))

	subUrl.recvPath, subUrl.recvFile, subUrl.err = setupFile(filepath.Join(connectionDirectory, subUrl.subDirectory), "recv")
	if subUrl.c, subUrl.err = Fatal(subUrl.err, "sam-url.go Pipe setup error", "sam-url.go Pipe setup"); subUrl.c {
		subUrl.recvFile.WriteString("")
	}

	subUrl.timePath, subUrl.timeFile, subUrl.err = setupFiFo(filepath.Join(connectionDirectory, subUrl.subDirectory), "time")
	if subUrl.c, subUrl.err = Fatal(subUrl.err, "Pipe setup error", "sam-url.go Pipe setup"); subUrl.c {
		subUrl.timeFile.WriteString("")
	}

	subUrl.delPath, subUrl.delPipe, subUrl.err = setupFiFo(filepath.Join(connectionDirectory, subUrl.subDirectory), "del")
	if subUrl.c, subUrl.err = Fatal(subUrl.err, "sam-url.go Pipe setup error", "sam-url.go Pipe setup"); subUrl.c {
		subUrl.delScan, subUrl.err = setupScanner(filepath.Join(connectionDirectory, subUrl.subDirectory), "del", subUrl.delPipe)
		if subUrl.c, subUrl.err = Fatal(subUrl.err, "sam-url.go Scanner setup Error:", "sam-url.go Scanner set up successfully."); !subUrl.c {
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
	if subUrl.c, subUrl.err = Fatal(err, "sam-url.go Scanner error", "sam-url.go Scanning recv"); subUrl.c {
		return "", subUrl.err
	}
	s := string(d)
	if s != "" {
		Log("sam-url.go Read file", s)
		return s, err
	}
	return "", err
}

func (subUrl *samUrl) dirSet(requestdir string) string {
	Log("sam-url.go Requesting directory: ", requestdir+"/")
	d1 := requestdir
	d2 := strings.Replace(d1, "//", "/", -1)
	return d2
}

func (subUrl *samUrl) checkDirectory(directory string) bool {
	b := false
	if directory == subUrl.subDirectory {
		Log("sam-url.go Directory / ", directory+" : equals : "+subUrl.subDirectory)
		b = true
	} else {
		Log("sam-url.go Directory / ", directory+" : does not equal : "+subUrl.subDirectory)
	}
	return b
}

func (subUrl *samUrl) copyDirectory(response *http.Response, directory string) bool {
	b := false
	subUrl.mutex.Lock()
	if subUrl.checkDirectory(directory) {
		if response != nil {
			Log("sam-url.go Response Status ", response.Status)
			if response.StatusCode == http.StatusOK {
				Log("sam-url.go Setting file in cache")
				subUrl.dealResponse(response)
			}
		}
		b = true
	}
	subUrl.mutex.Unlock()
	return b
}

func (subUrl *samUrl) copyDirectoryHttp(request *http.Request, response *http.Response, directory string) *http.Response {
	subUrl.mutex.Lock()
	if subUrl.checkDirectory(directory) {
		if response != nil {
			Log("sam-url.go Response Status ", response.Status)
			if response.StatusCode == http.StatusOK {
				Log("sam-url.go Setting file in cache")
				resp := subUrl.dealResponseHttp(request, response)
				return resp
			}
		}
	}
	subUrl.mutex.Unlock()
	return response
}

func (subUrl *samUrl) dealResponse(response *http.Response) {
	//defer
	body, err := ioutil.ReadAll(response.Body)
	//defer response.Body.Close()
	if subUrl.c, subUrl.err = Warn(err, "sam-url.go Response Write Error", "sam-url.go Writing responses"); subUrl.c {
		Log("sam-url.go Writing files.")
		subUrl.recvFile.Write(body)
		Log("sam-url.go Retrieval time: ", time.Now().String())
		subUrl.timeFile.WriteString(time.Now().String())
	}
}

func (subUrl *samUrl) printHeader(src http.Header) {
	if src != nil {
		for k, vv := range src {
			if vv != nil {
				for _, v := range vv {
					if v != "" {
						Log("sam-url.go Copying headers: " + k + "," + v)
					}
				}
			}
		}
	}
}

//func (subUrl *samUrl) dealResponseHttp(request *http.Request, response *http.Response) *http.Response {
func (subUrl *samUrl) dealResponseHttp(request *http.Request, response *http.Response) *http.Response {
	defer response.Body.Close()
	transferEncoding := response.TransferEncoding
	unCompressed := response.Uncompressed
	header := response.Header
	trailer := response.Trailer
	status := response.Status
	statusCode := response.StatusCode
	proto := response.Proto
	protoMajor := response.ProtoMajor
	protoMinor := response.ProtoMinor
	doClose := response.Close
	//responseBody := response.Body
	//contentLength := response.ContentLength
	//doClose := false
	body, err := ioutil.ReadAll(response.Body)
	//response.Body.Close()
	if subUrl.c, subUrl.err = Warn(err, "sam-url.go Response read error", "sam-url.go Reading response from proxy"); subUrl.c {
		Log("sam-url.go Writing files.")
		_, e := subUrl.recvFile.Write(body)
		contentLength := int64(len(body))
		if subUrl.c, subUrl.err = Warn(e, "sam-url.go File writing error", "sam-url.go Wrote response to file"); subUrl.c {
			r := &http.Response{
				Status:           status,
				StatusCode:       statusCode,
				Proto:            proto,
				ProtoMajor:       protoMajor,
				ProtoMinor:       protoMinor,
				Body:             ioutil.NopCloser(bytes.NewBuffer(body)),
				ContentLength:    contentLength,
				Request:          request,
				Header:           header,
				Trailer:          trailer,
				TransferEncoding: transferEncoding,
				Uncompressed:     unCompressed,
				Close:            doClose,
			}
			subUrl.printHeader(header)
			Log("sam-url.go Retrieval time: ", time.Now().String())
			subUrl.timeFile.WriteString(time.Now().String())
			return r
		} else {
			return nil
		}
	}
	return nil
}

func (subUrl *samUrl) cleanupDirectory() {
	subUrl.recvFile.Close()
	subUrl.timeFile.Close()
	subUrl.delPipe.Close()
	os.RemoveAll(filepath.Join(connectionDirectory, subUrl.subDirectory))
}

func (subUrl *samUrl) readDelete() bool {
	Log("sam-url.go Managing pipes:")
	for subUrl.delScan.Scan() {
		if subUrl.delScan.Text() == "y" || subUrl.delScan.Text() == "Y" {
			defer subUrl.cleanupDirectory()
			return true
		} else {
			return false
		}
	}
	clearFile(filepath.Join(connectionDirectory, subUrl.subDirectory), "del")
	return false
}

func NewSamUrl(requestdir string) samUrl {
	Log("sam-url.go Creating a new cache directory.")
	var subUrl samUrl
	subUrl.createDirectory(requestdir)
	subUrl.mutex = &sync.Mutex{}
	return subUrl
}


func NewSamUrlHttp(request *http.Request) samUrl {
	Log("sam-url.go Creating a new cache directory.")
	var subUrl samUrl
	log.Println(subUrl.subDirectory)
	subUrl.createDirectory(request.Host + request.URL.Path)
	subUrl.mutex = &sync.Mutex{}
	return subUrl
}

