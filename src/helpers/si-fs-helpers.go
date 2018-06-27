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

func safeURLString(str string) string {
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
	dii2phelpererrs.Log("si-fs-helpers.go final directory", r)
	return r
}

func truncatePaths(str string) string {
	temp := strings.SplitN(
		safeURLString(str),
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
				dii2phelpererrs.Log("si-fs-helpers.go fixedpath", fixedpath)
			}
		}
	}
	if strings.HasSuffix(fixedpath, "/") {
		fixedpath = strings.TrimSuffix(fixedpath, "/")
	}
	return fixedpath
}

func setupFolder(directory string) bool {
	pathConnectionExists, err := exists(truncatePaths(filepath.Join(connectionDirectory, directory)))
	if e, _ := dii2phelpererrs.Fatal(err, "si-fs-helpers.go Parent Directory Error", "si-fs-helpers.go Parent Directory Check", truncatePaths(filepath.Join(connectionDirectory))); e {
		if !pathConnectionExists {
			dii2phelpererrs.Log("si-fs-helpers.go Creating a connection:", directory)
			os.Mkdir(truncatePaths(filepath.Join(connectionDirectory, directory)), 0755)
			return true
		}
		os.RemoveAll(truncatePaths(filepath.Join(connectionDirectory, directory)))
		dii2phelpererrs.Log("si-fs-helpers.go Creating a connection:", directory)
		os.Mkdir(truncatePaths(filepath.Join(connectionDirectory, directory)), 0755)
		return true
	}
	return false
}

func checkFolder(directory string) bool {
	pathConnectionExists, err := exists(truncatePaths(filepath.Join(connectionDirectory, directory)))
	if e, _ := dii2phelpererrs.Fatal(err, "si-fs-helpers.go Child Directory Error", "si-fs-helpers.go Child Directory Check", truncatePaths(filepath.Join(connectionDirectory))); e {
		if !pathConnectionExists {
			dii2phelpererrs.Log("si-fs-helpers.go Creating a child directory folder:", directory)
			os.MkdirAll(truncatePaths(filepath.Join(connectionDirectory, directory)), 0755)
			return true
		}
		return false
	}
	return false
}

func setupFile(directory, path string) (string, *os.File, error) {
	mkPath := truncatePaths(filepath.Join(connectionDirectory, directory, path))
	pathExists, pathErr := exists(mkPath)
	if e, c := dii2phelpererrs.Fatal(pathErr, "si-fs-helpers.go File Check Error", "si-fs-helpers.go File Check", mkPath); e {
		if !pathExists {
			dii2phelpererrs.Log("si-fs-helpers.go Preparing to create File:", mkPath)
			file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
			if f, d := dii2phelpererrs.Fatal(err, "si-fs-helpers.go File Check Error", "si-fs-helpers.go File Check", mkPath); f {
				return mkPath, file, d
			}
		} else {
			g := os.Remove(mkPath)
			if f, d := dii2phelpererrs.Fatal(g, "si-fs-helpers.go File Check Error", "si-fs-helpers.go File Check", mkPath); !f {
				return mkPath, nil, d
			}
			file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
			if h, i := dii2phelpererrs.Fatal(err, "si-fs-helpers.go File Check Error", "si-fs-helpers.go File Check", mkPath); h {
				return mkPath, file, i
			}
			return mkPath, nil, err
		}
	} else {
		return mkPath, nil, c
	}
	return mkPath, nil, pathErr
}

func setupFiFo(directory, path string) (string, *os.File, error) {
	mkPath := truncatePaths(filepath.Join(connectionDirectory, directory, path))
	pathExists, pathErr := exists(mkPath)
	if e, c := dii2phelpererrs.Fatal(pathErr, "si-fs-helpers.go File Check Error", "si-fs-helpers.go File Check", mkPath); e {
		if !pathExists {
			mkErr := syscall.Mkfifo(mkPath, 0755)
			dii2phelpererrs.Log("si-fs-helpers.go Preparing to create Pipe:", mkPath)
			if f, _ := dii2phelpererrs.Fatal(mkErr, "si-fs-helpers.go Pipe Creation Error", "si-fs-helpers.go Creating Pipe", mkPath); f {
				file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
				return mkPath, file, err
			}
			return mkPath, nil, c
		}
		file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
		return mkPath, file, err
	}
	return mkPath, nil, nil
}

func setupScanner(directory, path string, pipe *os.File) (*bufio.Scanner, error) {
	mkPath := truncatePaths(filepath.Join(connectionDirectory, directory, path))
	_, pathErr := exists(mkPath)
	if e, c := dii2phelpererrs.Fatal(pathErr, "si-fs-helpers.go File Check Error", "si-fs-helpers.go File Check", mkPath); e {
		dii2phelpererrs.Log("si-fs-helpers.go Opening the Named Pipe as a Scanner...")
		retScanner := bufio.NewScanner(pipe)
		retScanner.Split(bufio.ScanLines)
		dii2phelpererrs.Log("si-fs-helpers.go Created a named Pipe for sending requests:", mkPath)
		return retScanner, c
	}
	return nil, pathErr
}

//func setupCookieJar()

//This function does nothing on Unix-like platforms. It is only here to clear
//the contents of files that would normally be named pipes on Windows.
func clearFile(directory, path string) {

}