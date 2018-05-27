package dii2p

import (
	"bufio"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type samURL struct {
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

func (subURL *samURL) initPipes() {
	checkFolder(filepath.Join(connectionDirectory, subURL.subDirectory))

	subURL.recvPath, subURL.recvFile, subURL.err = setupFile(filepath.Join(connectionDirectory, subURL.subDirectory), "recv")
	if subURL.c, subURL.err = Fatal(subURL.err, "sam-url.go Pipe setup error", "sam-url.go Pipe setup"); subURL.c {
		subURL.recvFile.WriteString("")
	}

	subURL.timePath, subURL.timeFile, subURL.err = setupFiFo(filepath.Join(connectionDirectory, subURL.subDirectory), "time")
	if subURL.c, subURL.err = Fatal(subURL.err, "Pipe setup error", "sam-url.go Pipe setup"); subURL.c {
		subURL.timeFile.WriteString("")
	}

	subURL.delPath, subURL.delPipe, subURL.err = setupFiFo(filepath.Join(connectionDirectory, subURL.subDirectory), "del")
	if subURL.c, subURL.err = Fatal(subURL.err, "sam-url.go Pipe setup error", "sam-url.go Pipe setup"); subURL.c {
		subURL.delScan, subURL.err = setupScanner(filepath.Join(connectionDirectory, subURL.subDirectory), "del", subURL.delPipe)
		if subURL.c, subURL.err = Fatal(subURL.err, "sam-url.go Scanner setup Error:", "sam-url.go Scanner set up successfully."); !subURL.c {
			subURL.cleanupDirectory()
		}
	}

}

func (subURL *samURL) createDirectory(requestdir string) {
	subURL.subDirectory = subURL.dirSet(requestdir)
	subURL.initPipes()
}

func (subURL *samURL) scannerText() (string, error) {
	d, err := ioutil.ReadFile(subURL.recvPath)
	if subURL.c, subURL.err = Fatal(err, "sam-url.go Scanner error", "sam-url.go Scanning recv"); subURL.c {
		return "", subURL.err
	}
	s := string(d)
	if s != "" {
		Log("sam-url.go Read file", s)
		return s, err
	}
	return "", err
}

func (subURL *samURL) dirSet(requestdir string) string {
	Log("sam-url.go Requesting directory: ", requestdir+"/")
	d1 := requestdir
	d2 := strings.Replace(d1, "//", "/", -1)
	return d2
}

func (subURL *samURL) checkDirectory(directory string) bool {
	b := false
	if directory == subURL.subDirectory {
		Log("sam-url.go Directory / ", directory+" : equals : "+subURL.subDirectory)
		b = true
	} else {
		Log("sam-url.go Directory / ", directory+" : does not equal : "+subURL.subDirectory)
	}
	return b
}

func (subURL *samURL) copyDirectory(response *http.Response, directory string) bool {
	b := false
	subURL.mutex.Lock()
	if subURL.checkDirectory(directory) {
		if response != nil {
			Log("sam-url.go Response Status ", response.Status)
			if response.StatusCode == http.StatusOK {
				Log("sam-url.go Setting file in cache")
				subURL.dealResponse(response)
			}
		}
		b = true
	}
	subURL.mutex.Unlock()
	return b
}

func (subURL *samURL) copyDirectoryHTTP(request *http.Request, response *http.Response, directory string) *http.Response {
	subURL.mutex.Lock()
	if subURL.checkDirectory(directory) {
		if response != nil {
			Log("sam-url.go Response Status ", response.Status)
			if response.StatusCode == http.StatusOK {
				Log("sam-url.go Setting file in cache")
				resp := subURL.dealResponseHTTP(request, response)
				return resp
			}
		}
	}
	subURL.mutex.Unlock()
	return response
}

func (subURL *samURL) dealResponse(response *http.Response) {
	//defer
	body, err := ioutil.ReadAll(response.Body)
	//defer response.Body.Close()
	if subURL.c, subURL.err = Warn(err, "sam-url.go Response Write Error", "sam-url.go Writing responses"); subURL.c {
		Log("sam-url.go Writing files.")
		subURL.recvFile.Write(body)
		Log("sam-url.go Retrieval time: ", time.Now().String())
		subURL.timeFile.WriteString(time.Now().String())
	}
}

func (subURL *samURL) printHeader(src http.Header) {
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

//func (subURL *samURL) dealResponseHTTP(request *http.Request, response *http.Response) *http.Response {
func (subURL *samURL) dealResponseHTTP(request *http.Request, response *http.Response) *http.Response {
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
	if subURL.c, subURL.err = Warn(err, "sam-url.go Response read error", "sam-url.go Reading response from proxy"); subURL.c {
		Log("sam-url.go Writing files.")
		_, e := subURL.recvFile.Write(body)
		contentLength := int64(len(body))
		if subURL.c, subURL.err = Warn(e, "sam-url.go File writing error", "sam-url.go Wrote response to file"); subURL.c {
			r := &http.Response{
				Status:           status,
				StatusCode:       statusCode,
				Proto:            proto,
				ProtoMajor:       protoMajor,
				ProtoMinor:       protoMinor,
				Body:             ioutil.NopCloser(strings.NewReader(string(body))),
				ContentLength:    contentLength,
				Request:          request,
				Header:           header,
				Trailer:          trailer,
				TransferEncoding: transferEncoding,
				Uncompressed:     unCompressed,
				Close:            doClose,
			}
			subURL.printHeader(header)
			Log("sam-url.go Retrieval time: ", time.Now().String())
			subURL.timeFile.WriteString(time.Now().String())
			return r
		} else {
			return nil
		}
	}
	return nil
}

func (subURL *samURL) cleanupDirectory() {
	subURL.recvFile.Close()
	subURL.timeFile.Close()
	subURL.delPipe.Close()
	os.RemoveAll(filepath.Join(connectionDirectory, subURL.subDirectory))
}

func (subURL *samURL) readDelete() bool {
	Log("sam-url.go Managing pipes:")
	for subURL.delScan.Scan() {
		if subURL.delScan.Text() == "y" || subURL.delScan.Text() == "Y" {
			defer subURL.cleanupDirectory()
			return true
		} else {
			return false
		}
	}
	clearFile(filepath.Join(connectionDirectory, subURL.subDirectory), "del")
	return false
}

func NewSamURL(requestdir string) samURL {
	Log("sam-url.go Creating a new cache directory.")
	var subURL samURL
	subURL.createDirectory(requestdir)
	subURL.mutex = &sync.Mutex{}
	return subURL
}

func NewSamURLHTTP(request *http.Request) samURL {
	Log("sam-url.go Creating a new cache directory.")
	var subURL samURL
	log.Println(subURL.subDirectory)
	subURL.createDirectory(request.Host + request.URL.Path)
	subURL.mutex = &sync.Mutex{}
	return subURL
}
