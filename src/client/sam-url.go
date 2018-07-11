package dii2pmain

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

import (
	"github.com/eyedeekay/si-i2p-plugin/src/errors"
	"github.com/eyedeekay/si-i2p-plugin/src/helpers"
)

//SamURL manages the recieve pipes for SamHTTP request targets
type SamURL struct {
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

func (subURL *SamURL) initPipes() {
	dii2phelper.SetupFolder(filepath.Join(dii2phelper.ConnectionDirectory, subURL.subDirectory))

	subURL.recvPath, subURL.recvFile, subURL.err = dii2phelper.SetupFile(filepath.Join(dii2phelper.ConnectionDirectory, subURL.subDirectory), "recv")
	if subURL.c, subURL.err = dii2perrs.Fatal(subURL.err, "sam-url.go Pipe setup error", "sam-url.go Pipe setup"); subURL.c {
		subURL.recvFile.WriteString("")
	}

	subURL.timePath, subURL.timeFile, subURL.err = dii2phelper.SetupFiFo(filepath.Join(dii2phelper.ConnectionDirectory, subURL.subDirectory), "time")
	if subURL.c, subURL.err = dii2perrs.Fatal(subURL.err, "Pipe setup error", "sam-url.go Pipe setup"); subURL.c {
		subURL.timeFile.WriteString("")
	}

	subURL.delPath, subURL.delPipe, subURL.err = dii2phelper.SetupFiFo(filepath.Join(dii2phelper.ConnectionDirectory, subURL.subDirectory), "del")
	if subURL.c, subURL.err = dii2perrs.Fatal(subURL.err, "sam-url.go Pipe setup error", "sam-url.go Pipe setup"); subURL.c {
		subURL.delScan, subURL.err = dii2phelper.SetupScanner(filepath.Join(dii2phelper.ConnectionDirectory, subURL.subDirectory), "del", subURL.delPipe)
		if subURL.c, subURL.err = dii2perrs.Fatal(subURL.err, "sam-url.go Scanner setup Error:", "sam-url.go Scanner set up successfully."); !subURL.c {
			subURL.CleanupDirectory()
		}
	}

}

func (subURL *SamURL) CreateDirectory(requestdir string) {
	subURL.subDirectory = subURL.dirSet(requestdir)
	subURL.initPipes()
}

func (subURL *SamURL) ScannerText() (string, error) {
	d, err := ioutil.ReadFile(subURL.recvPath)
	if subURL.c, subURL.err = dii2perrs.Fatal(err, "sam-url.go Scanner error", "sam-url.go Scanning recv"); subURL.c {
		return "", subURL.err
	}
	s := string(d)
	if s != "" {
		dii2perrs.Log("sam-url.go Read file", s)
		return s, err
	}
	return "", err
}

func (subURL *SamURL) dirSet(requestdir string) string {
	dii2perrs.Log("sam-url.go Requesting directory: ", requestdir+"/")
	d1 := requestdir
	d2 := strings.Replace(d1, "//", "/", -1)
	return d2
}

func (subURL *SamURL) checkDirectory(directory string) bool {
	b := false
	if directory == subURL.subDirectory {
		dii2perrs.Log("sam-url.go Directory / ", directory+" : equals : "+subURL.subDirectory)
		b = true
	} else {
		dii2perrs.Log("sam-url.go Directory / ", directory+" : does not equal : "+subURL.subDirectory)
	}
	return b
}

func (subURL *SamURL) copyDirectory(response *http.Response, directory string) bool {
	b := false
	subURL.mutex.Lock()
	if subURL.checkDirectory(directory) {
		if response != nil {
			dii2perrs.Log("sam-url.go Response Status ", response.Status)
			if response.StatusCode == http.StatusOK {
				dii2perrs.Log("sam-url.go Setting file in cache")
				subURL.dealResponse(response)
			}
		}
		b = true
	}
	subURL.mutex.Unlock()
	return b
}

func (subURL *SamURL) copyDirectoryHTTP(request *http.Request, response *http.Response, directory string) *http.Response {
	subURL.mutex.Lock()
	if subURL.checkDirectory(directory) {
		if response != nil {
			dii2perrs.Log("sam-url.go Response Status ", response.Status)
			if response.StatusCode == http.StatusOK {
				dii2perrs.Log("sam-url.go Setting file in cache")
				resp := subURL.dealResponseHTTP(request, response)
				return resp
			}
		}
	}
	subURL.mutex.Unlock()
	return response
}

func (subURL *SamURL) dealResponse(response *http.Response) {
	//defer
	body, err := ioutil.ReadAll(response.Body)
	//defer response.Body.Close()
	if subURL.c, subURL.err = dii2perrs.Warn(err, "sam-url.go Response Write Error", "sam-url.go Writing responses"); subURL.c {
		dii2perrs.Log("sam-url.go Writing files.")
		subURL.recvFile.Write(body)
		dii2perrs.Log("sam-url.go Retrieval time: ", time.Now().String())
		subURL.timeFile.WriteString(time.Now().String())
	}
}

func (subURL *SamURL) printHeader(src http.Header) {
	if src != nil {
		for k, vv := range src {
			if vv != nil {
				for _, v := range vv {
					if v != "" {
						dii2perrs.Log("sam-url.go Copying headers: " + k + "," + v)
					}
				}
			}
		}
	}
}

//func (subURL *SamURL) dealResponseHTTP(request *http.Request, response *http.Response) *http.Response {
func (subURL *SamURL) dealResponseHTTP(request *http.Request, response *http.Response) *http.Response {
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
	if subURL.c, subURL.err = dii2perrs.Warn(err, "sam-url.go Response read error", "sam-url.go Reading response from proxy"); subURL.c {
		dii2perrs.Log("sam-url.go Writing files.")
		_, e := subURL.recvFile.Write(body)
		contentLength := int64(len(body))
		if subURL.c, subURL.err = dii2perrs.Warn(e, "sam-url.go File writing error", "sam-url.go Wrote response to file"); subURL.c {
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
			dii2perrs.Log("sam-url.go Retrieval time: ", time.Now().String())
			subURL.timeFile.WriteString(time.Now().String())
			return r
		}
		return nil
	}
	return nil
}

func (subURL *SamURL) CleanupDirectory() {
	subURL.recvFile.Close()
	subURL.timeFile.Close()
	subURL.delPipe.Close()
	os.RemoveAll(filepath.Join(dii2phelper.ConnectionDirectory, subURL.subDirectory))
}

func (subURL *SamURL) readDelete() bool {
	dii2perrs.Log("sam-url.go Managing pipes:")
	for subURL.delScan.Scan() {
		if subURL.delScan.Text() == "y" || subURL.delScan.Text() == "Y" {
			defer subURL.CleanupDirectory()
			return true
		}
		return false
	}
	dii2phelper.ClearFile(filepath.Join(dii2phelper.ConnectionDirectory, subURL.subDirectory), "del")
	return false
}

//NewSamURL instantiates a SamURL
func NewSamURL(requestdir string) SamURL {
	dii2perrs.Log("sam-url.go Creating a new cache directory.")
	var subURL SamURL
	subURL.CreateDirectory(requestdir)
	subURL.mutex = &sync.Mutex{}
	return subURL
}

//NewSamURLHTTP instantiates a SamURL
func NewSamURLHTTP(request *http.Request) SamURL {
	dii2perrs.Log("sam-url.go Creating a new cache directory.")
	var subURL SamURL
	log.Println(subURL.subDirectory)
	subURL.CreateDirectory(request.Host + request.URL.Path)
	subURL.mutex = &sync.Mutex{}
	return subURL
}
