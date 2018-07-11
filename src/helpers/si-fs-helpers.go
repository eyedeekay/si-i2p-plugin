// +build !windows

package dii2phelper

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

import (
	"github.com/eyedeekay/si-i2p-plugin/src/errors"
)

var DiskAvoidance = false

// ConnectionDirectory is the global working directory of the service
var ConnectionDirectory string

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func truncatePath(str string) string {
	num := 90
	bnoden := str
	if len(str) > num {
		bnoden = str[0:num]
	}
	return bnoden
}

func safeNames(str string) string {
	switch d := str; d {
	case "send":
		return "_send"
	case "recv":
		return "_recv"
	case "del":
		return "_del"
	case "time":
		return "_time"
	case "name":
		return "_name"
	case "id":
		return "_id"
	case "base64":
		return "_base64"
	default:
		return str
	}
}

// SafeURLString removes illegal characters from URL strings
func SafeURLString(str string) string {
	r := strings.Replace(
		strings.Replace(
			str,
			"http:",
			"",
			-1,
		),
		"//",
		"/",
		-1,
	)

	if strings.HasSuffix(r, "/") {
		r = strings.TrimSuffix(r, "/")
	}
	if strings.HasPrefix(r, "/") {
		r = strings.TrimPrefix(r, "/")
	}
	dii2perrs.Log("si-fs-helpers.go final directory", r)
	return r
}

func truncatePaths(str string) string {
	temp := strings.SplitN(
		SafeURLString(str),
		"/",
		-1,
	)
	var fixedpath string
	for x, i := range temp {
		if i != "" {
			if x < len(temp)-1 {
				fixedpath += safeNames(truncatePath(i)) + "/"
			} else {
				fixedpath += truncatePath(i)
			}
			if fixedpath != "" {
				dii2perrs.Log("si-fs-helpers.go fixedpath", fixedpath)
			}
		}
	}
	if strings.HasSuffix(fixedpath, "/") {
		fixedpath = strings.TrimSuffix(fixedpath, "/")
	}
	return fixedpath
}

// SetupFolder Creates a folder for a site or directory control interface
func SetupFolder(directory string) bool {
	pathConnectionExists, err := exists(truncatePaths(filepath.Join(ConnectionDirectory, directory)))
	if e, _ := dii2perrs.Fatal(err, "si-fs-helpers.go Parent Directory Error", "si-fs-helpers.go Parent Directory Check", truncatePaths(filepath.Join(ConnectionDirectory))); e {
		if !pathConnectionExists {
			dii2perrs.Log("si-fs-helpers.go Creating a connection:", directory)
			os.MkdirAll(truncatePaths(filepath.Join(ConnectionDirectory, directory)), 0755)
			return true
		}
		os.RemoveAll(truncatePaths(filepath.Join(ConnectionDirectory, directory)))
		dii2perrs.Log("si-fs-helpers.go Creating a connection:", directory)
		os.MkdirAll(truncatePaths(filepath.Join(ConnectionDirectory, directory)), 0755)
		return true
	}
	return false
}

// CheckFolder ensures a folder exists. It's probably redundant.
func CheckFolder(directory string) bool {
	pathConnectionExists, err := exists(truncatePaths(filepath.Join(ConnectionDirectory, directory)))
	if e, _ := dii2perrs.Fatal(err, "si-fs-helpers.go Child Directory Error", "si-fs-helpers.go Child Directory Check", truncatePaths(filepath.Join(ConnectionDirectory))); e {
		if !pathConnectionExists {
			dii2perrs.Log("si-fs-helpers.go Creating a child directory folder:", directory)
			os.MkdirAll(truncatePaths(filepath.Join(ConnectionDirectory, directory)), 0755)
			return true
		}
		return false
	}
	return false
}

// SetupFile creates a regular file
func SetupFile(directory, path string) (string, *os.File, error) {
	mkPath := truncatePaths(filepath.Join(ConnectionDirectory, directory, path))
	pathExists, pathErr := exists(mkPath)
	if e, c := dii2perrs.Fatal(pathErr, "si-fs-helpers.go File Check Error", "si-fs-helpers.go File Check", mkPath); e {
		if !pathExists {
			dii2perrs.Log("si-fs-helpers.go Preparing to create File:", mkPath)
			file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
			if f, d := dii2perrs.Fatal(err, "si-fs-helpers.go File Check Error", "si-fs-helpers.go File Check", mkPath); f {
				return mkPath, file, d
			}
		} else {
			g := os.Remove(mkPath)
			if f, d := dii2perrs.Fatal(g, "si-fs-helpers.go File Check Error", "si-fs-helpers.go File Check", mkPath); !f {
				return mkPath, nil, d
			}
			file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
			if h, i := dii2perrs.Fatal(err, "si-fs-helpers.go File Check Error", "si-fs-helpers.go File Check", mkPath); h {
				return mkPath, file, i
			}
			return mkPath, nil, err
		}
	} else {
		return mkPath, nil, c
	}
	return mkPath, nil, pathErr
}

// SetupFiFo creates a named pipe
func SetupFiFo(directory, path string) (string, *os.File, error) {
	mkPath := truncatePaths(filepath.Join(ConnectionDirectory, directory, path))
	pathExists, pathErr := exists(mkPath)
	if e, _ := dii2perrs.Fatal(pathErr, "si-fs-helpers.go File Check Error", "si-fs-helpers.go File Check", mkPath); e {
		if !pathExists {
			mkErr := syscall.Mkfifo(mkPath, 0755)
			dii2perrs.Log("si-fs-helpers.go Preparing to create Pipe:", mkPath)
			if f, _ := dii2perrs.Warn(mkErr, "si-fs-helpers.go Pipe Creation Error", "si-fs-helpers.go Creating Pipe", mkPath); f {
				file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
				return mkPath, file, err
			}
            CheckFolder(filepath.Join(ConnectionDirectory, directory))
            file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
            return mkPath, file, err
			//return mkPath, nil, c
		}
		file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
		return mkPath, file, err
	}
	return mkPath, nil, nil
}

// SetupScanner sets up a scanner
func SetupScanner(directory, path string, pipe *os.File) (*bufio.Scanner, error) {
	mkPath := truncatePaths(filepath.Join(ConnectionDirectory, directory, path))
	_, pathErr := exists(mkPath)
	if e, c := dii2perrs.Fatal(pathErr, "si-fs-helpers.go File Check Error", "si-fs-helpers.go File Check", mkPath); e {
		dii2perrs.Log("si-fs-helpers.go Opening the Named Pipe as a Scanner...")
		retScanner := bufio.NewScanner(pipe)
		retScanner.Split(bufio.ScanLines)
		dii2perrs.Log("si-fs-helpers.go Created a named Pipe for sending requests:", mkPath)
		return retScanner, c
	}
	return nil, pathErr
}

//func SetupCookieJar()

//This function does nothing on Unix-like platforms. It is only here to clear
//the contents of files that would normally be named pipes on Windows.

// ClearFile does nothing here.
func ClearFile(directory, path string) {

}
