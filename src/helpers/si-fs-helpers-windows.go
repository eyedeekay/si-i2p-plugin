// +build windows

package dii2phelper

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

import (
	"github.com/eyedeekay/si-i2p-plugin/src/errors"
)

var DiskAvoidance = false

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

func truncatePaths(str string) string {
	temp := strings.SplitN(str, "/", -1)
	var fixedpath string
	for _, i := range temp {
		if i != "" {
			fixedpath += truncatePath(i) + "/"
			if fixedpath != "" {
				dii2perrs.Log("si-fs-helpers-windows.go ", truncatePath(i))
			}
		}
	}
	if strings.HasSuffix(fixedpath, "/") {
		fixedpath = fixedpath[:len(fixedpath)-len("/")]
	}
	return fixedpath
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
	return str
}

func SafeURLString(str string) string {
	temp := strings.SplitN(str, "/", -1)
	last := safeNames(temp[len(temp)-1])
	var r string
	for x, i := range temp {
		if x != len(temp)-1 {
			r += i
		} else {
			r += last
		}
	}
	return r
}

func SetupFolder(directory string) bool {
	pathConnectionExists, err := exists(truncatePaths(filepath.Join(connectionDirectory, directory)))
	if e, _ := dii2perrs.Fatal(err, "si-fs-helpers-windows.go Parent Directory Error", "si-fs-helpers-windows.go Parent Directory Check", truncatePaths(filepath.Join(connectionDirectory))); e {
		if !pathConnectionExists {
			dii2perrs.Log("si-fs-helpers-windows.go Creating a connection:", directory)
			os.Mkdir(truncatePaths(filepath.Join(connectionDirectory, directory)), 0755)
			return true
		}
		os.RemoveAll(truncatePaths(filepath.Join(connectionDirectory, directory)))
		dii2perrs.Log("si-fs-helpers-windows.go Creating a connection:", directory)
		os.Mkdir(truncatePaths(filepath.Join(connectionDirectory, directory)), 0755)
		return true

	}
	return false
}

func SetupFile(directory, path string) (string, *os.File, error) {
	mkPath := truncatePaths(filepath.Join(connectionDirectory, directory, path))
	pathExists, pathErr := exists(mkPath)
	if e, c := dii2perrs.Fatal(pathErr, "si-fs-helpers-windows.go File Check Error", "si-fs-helpers-windows.go File Check", mkPath); e {
		if !pathExists {
			dii2perrs.Log("si-fs-helpers-windows.go Preparing to create File:", mkPath)
			file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
			if f, d := dii2perrs.Fatal(err, "si-fs-helpers-windows.go File Check Error", "si-fs-helpers-windows.go File Check", mkPath); f {
				return mkPath, file, d
			}
		} else {
			g := os.Remove(mkPath)
			if f, d := dii2perrs.Fatal(g, "si-fs-helpers-windows.go File Check Error", "si-fs-helpers-windows.go File Check", mkPath); !f {
				return mkPath, nil, d
			}
			file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
			if h, i := dii2perrs.Fatal(err, "si-fs-helpers-windows.go File Check Error", "si-fs-helpers-windows.go File Check", mkPath); h {
				return mkPath, file, i
			}
			return mkPath, nil, err
		}
	} else {
		return mkPath, nil, c
	}
	return mkPath, nil, pathErr
}

func SetupFiFo(directory, path string) (string, *os.File, error) {
	mkPath := truncatePaths(filepath.Join(connectionDirectory, directory, path))
	pathExists, pathErr := exists(mkPath)
	if e, c := dii2perrs.Fatal(pathErr, "si-fs-helpers-windows.go File Check Error", "si-fs-helpers-windows.go File Check", mkPath); e {
		if !pathExists {
			dii2perrs.Log("si-fs-helpers-windows.go Preparing to create File:", mkPath)
			file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
			if f, d := dii2perrs.Fatal(err, "si-fs-helpers-windows.go File Check Error", "si-fs-helpers-windows.go File Check", mkPath); f {
				return mkPath, file, d
			}
		} else {
			g := os.Remove(mkPath)
			if f, d := dii2perrs.Fatal(g, "si-fs-helpers-windows.go File Check Error", "si-fs-helpers-windows.go File Check", mkPath); !f {
				return mkPath, nil, d
			}
			file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
			if h, i := dii2perrs.Fatal(err, "si-fs-helpers-windows.go File Check Error", "si-fs-helpers-windows.go File Check", mkPath); h {
				return mkPath, file, i
			}
			return mkPath, nil, err
		}
	} else {
		return mkPath, nil, c
	}
	return mkPath, nil, pathErr
}

func SetupScanner(directory, path string, pipe *os.File) (*bufio.Scanner, error) {
	mkPath := truncatePaths(filepath.Join(connectionDirectory, directory, path))
	_, pathErr := exists(mkPath)
	if e, c := dii2perrs.Fatal(pathErr, "si-fs-helpers-windows.go File Check Error", "si-fs-helpers-windows.go File Check", mkPath); e {
		dii2perrs.Log("si-fs-helpers-windows.go Opening the Named Pipe as a Scanner...")
		retScanner := bufio.NewScanner(pipe)
		retScanner.Split(bufio.ScanLines)
		dii2perrs.Log("si-fs-helpers-windows.go Created a named Pipe for sending requests:", mkPath)
		return retScanner, nil
	}
	return nil, c
}

//func SetupCookieJar()

func ClearFile(directory, path string) {
	mkPath := filepath.Join(connectionDirectory, directory, path)
	clearErr := os.Truncate(mkPath, 0)
	if e, c := dii2perrs.Fatal(clearErr, "si-fs-helpers-windows.go File Clear Error", "si-fs-helpers-windows.go File Cleared", mkPath); e {
		dii2perrs.Log("si-fs-helpers-windows.go Input file cleared.")
	} else {
		dii2perrs.Log("si-fs-helpers-windows.go Input file cleared.", c.Error())
	}
}
