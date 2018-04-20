// +build windows

package main

import (
	"bufio"
	//"byte"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

func truncatePaths(str string) string {
	temp := strings.SplitN(str, "/", -1)
	var fixedpath string
	for _, i := range temp {
		if i != "" {
			fixedpath += truncatePath(i) + "/"
			if fixedpath != "" {
				Log("si-fs-helpers-windows.go ", truncatePath(i))
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
	if e, _ := Fatal(err, "si-fs-helpers-windows.go Parent Directory Error", "si-fs-helpers-windows.go Parent Directory Check", truncatePaths(filepath.Join(connectionDirectory))); e {
		if !pathConnectionExists {
			Log("si-fs-helpers-windows.go Creating a connection:", directory)
			os.Mkdir(truncatePaths(filepath.Join(connectionDirectory, directory)), 0755)
			return true
		} else {
			os.RemoveAll(truncatePaths(filepath.Join(connectionDirectory, directory)))
			Log("si-fs-helpers-windows.go Creating a connection:", directory)
			os.Mkdir(truncatePaths(filepath.Join(connectionDirectory, directory)), 0755)
			return true
		}
	} else {
		return false
	}
}

func checkFolder(directory string) bool {
	pathConnectionExists, err := exists(truncatePaths(filepath.Join(connectionDirectory, directory)))
	if e, _ := Fatal(err, "si-fs-helpers-windows.go Child Directory Error", "si-fs-helpers-windows.go Child Directory Check", truncatePaths(filepath.Join(connectionDirectory))); e {
		if !pathConnectionExists {
			Log("si-fs-helpers-windows.go Creating a child directory folder:", directory)
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
	if e, c := Fatal(pathErr, "si-fs-helpers-windows.go File Check Error", "si-fs-helpers-windows.go File Check", mkPath); e {
		if !pathExists {
			Log("si-fs-helpers-windows.go Preparing to create File:", mkPath)
			file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
			if f, d := Fatal(err, "si-fs-helpers-windows.go File Check Error", "si-fs-helpers-windows.go File Check", mkPath); f {
				return mkPath, file, d
			}
		} else {
			g := os.Remove(mkPath)
			if f, d := Fatal(g, "si-fs-helpers-windows.go File Check Error", "si-fs-helpers-windows.go File Check", mkPath); !f {
				return mkPath, nil, d
			}
			file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
			if h, i := Fatal(err, "si-fs-helpers-windows.go File Check Error", "si-fs-helpers-windows.go File Check", mkPath); h {
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
	if e, c := Fatal(pathErr, "si-fs-helpers-windows.go File Check Error", "si-fs-helpers-windows.go File Check", mkPath); e {
		if !pathExists {
			Log("si-fs-helpers-windows.go Preparing to create File:", mkPath)
			file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
			if f, d := Fatal(err, "si-fs-helpers-windows.go File Check Error", "si-fs-helpers-windows.go File Check", mkPath); f {
				return mkPath, file, d
			}
		} else {
			g := os.Remove(mkPath)
			if f, d := Fatal(g, "si-fs-helpers-windows.go File Check Error", "si-fs-helpers-windows.go File Check", mkPath); !f {
				return mkPath, nil, d
			}
			file, err := os.OpenFile(mkPath, os.O_RDWR|os.O_CREATE, 0755)
			if h, i := Fatal(err, "si-fs-helpers-windows.go File Check Error", "si-fs-helpers-windows.go File Check", mkPath); h {
				return mkPath, file, i
			}
			return mkPath, nil, err
		}
	} else {
		return mkPath, nil, c
	}
	return mkPath, nil, pathErr
}

func setupScanner(directory, path string, pipe *os.File) (*bufio.Scanner, error) {
	mkPath := truncatePaths(filepath.Join(connectionDirectory, directory, path))
	_, pathErr := exists(mkPath)
	if e, c := Fatal(pathErr, "si-fs-helpers-windows.go File Check Error", "si-fs-helpers-windows.go File Check", mkPath); e {
		Log("si-fs-helpers-windows.go Opening the Named Pipe as a Scanner...")
		retScanner := bufio.NewScanner(pipe)
		retScanner.Split(bufio.ScanLines)
		Log("si-fs-helpers-windows.go Created a named Pipe for sending requests:", mkPath)
		return retScanner, nil
	} else {
		return nil, c
	}
}

//func setupCookieJar()

func clearFile(directory, path string) {
	mkPath := filepath.Join(connectionDirectory, directory, path)
	info := []byte{0}
	clearErr := ioutil.WriteFile(mkPath, info, 0755)
	if e, c := Fatal(clearErr, "si-fs-helpers-windows.go File Clear Error", "si-fs-helpers-windows.go File Cleared", mkPath); e {
		Log("si-fs-helpers-windows.go Input file cleared.")
	}else{
        Log("si-fs-helpers-windows.go Input file cleared.", c.Error())
    }
}
