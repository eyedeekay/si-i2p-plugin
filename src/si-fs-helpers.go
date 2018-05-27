// +build !windows

package dii2p

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"syscall"
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
    return str
}

func safeUrlString(str string) string {
    temp := strings.SplitN(str, "/", -1)
    var r string
    for _, i := range temp {
        r += safeNames(i)
    }
    return r
}

func truncatePaths(str string) string {
	temp := strings.SplitN(str, "/", -1)
	var fixedpath string
	for x, i := range temp {
        if x < len(temp)-1 {
            i = safeNames(i)
        }
		if i != "" {
			fixedpath += truncatePath(i) + "/"
			if fixedpath != "" {
				Log("si-fs-helpers.go ", truncatePath(i))
			}
		}
	}
	if strings.HasSuffix(fixedpath, "/") {
		fixedpath = fixedpath[:len(fixedpath)-len("/")]
	}
	return fixedpath
}

func setupFolder(directory string) bool {
	pathConnectionExists, err := exists(truncatePaths(filepath.Join(connectionDirectory, directory)))
	if e, _ := Fatal(err, "si-fs-helpers.go Parent Directory Error", "si-fs-helpers.go Parent Directory Check", truncatePaths(filepath.Join(connectionDirectory))); e {
		if !pathConnectionExists {
			Log("si-fs-helpers.go Creating a connection:", directory)
			os.Mkdir(truncatePaths(filepath.Join(connectionDirectory, directory)), 0755)
			return true
		} else {
			os.RemoveAll(truncatePaths(filepath.Join(connectionDirectory, directory)))
			Log("si-fs-helpers.go Creating a connection:", directory)
			os.Mkdir(truncatePaths(filepath.Join(connectionDirectory, directory)), 0755)
			return true
		}
	} else {
		return false
	}
}

func checkFolder(directory string) bool {
	pathConnectionExists, err := exists(truncatePaths(filepath.Join(connectionDirectory, directory)))
	if e, _ := Fatal(err, "si-fs-helpers.go Child Directory Error", "si-fs-helpers.go Child Directory Check", truncatePaths(filepath.Join(connectionDirectory))); e {
		if !pathConnectionExists {
			Log("si-fs-helpers.go Creating a child directory folder:", directory)
			os.MkdirAll(truncatePaths(filepath.Join(connectionDirectory, directory)), 0755)
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func setupFile(directory, path string) (string, *os.File, error) {
	mkPath := truncatePaths(filepath.Join(connectionDirectory, directory, path))
	pathExists, pathErr := exists(mkPath)
	if e, c := Fatal(pathErr, "si-fs-helpers.go File Check Error", "si-fs-helpers.go File Check", mkPath); e {
		if !pathExists {
			Log("si-fs-helpers.go Preparing to create File:", mkPath)
			file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
			if f, d := Fatal(err, "si-fs-helpers.go File Check Error", "si-fs-helpers.go File Check", mkPath); f {
				return mkPath, file, d
			}
		} else {
			g := os.Remove(mkPath)
			if f, d := Fatal(g, "si-fs-helpers.go File Check Error", "si-fs-helpers.go File Check", mkPath); !f {
				return mkPath, nil, d
			}
			file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
			if h, i := Fatal(err, "si-fs-helpers.go File Check Error", "si-fs-helpers.go File Check", mkPath); h {
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
	if e, c := Fatal(pathErr, "si-fs-helpers.go File Check Error", "si-fs-helpers.go File Check", mkPath); e {
		if !pathExists {
			mkErr := syscall.Mkfifo(mkPath, 0755)
			Log("si-fs-helpers.go Preparing to create Pipe:", mkPath)
			if f, d := Fatal(mkErr, "si-fs-helpers.go Pipe Creation Error", "si-fs-helpers.go Creating Pipe", mkPath); f {
				file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
				return mkPath, file, err
			} else {
				return mkPath, nil, d
			}
		} else {
			file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
			return mkPath, file, err
		}
	} else {
		return mkPath, nil, c
	}
}

func setupScanner(directory, path string, pipe *os.File) (*bufio.Scanner, error) {
	mkPath := truncatePaths(filepath.Join(connectionDirectory, directory, path))
	_, pathErr := exists(mkPath)
	if e, c := Fatal(pathErr, "si-fs-helpers.go File Check Error", "si-fs-helpers.go File Check", mkPath); e {
		Log("si-fs-helpers.go Opening the Named Pipe as a Scanner...")
		retScanner := bufio.NewScanner(pipe)
		retScanner.Split(bufio.ScanLines)
		Log("si-fs-helpers.go Created a named Pipe for sending requests:", mkPath)
		return retScanner, nil
	} else {
		return nil, c
	}
}

//func setupCookieJar()

//This function does nothing on Unix-like platforms. It is only here to clear
//the contents of files that would normally be named pipes on Windows.
func clearFile(directory, path string) {

}
