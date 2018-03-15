package main

import (
    "bufio"
	"os"
	"path/filepath"
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
    num := 254
	bnoden := str
	if len(str) > num {
		bnoden = str[0:num] + "..."
	}
	return bnoden
}

func setupFolder(directory string) bool {
	pathConnectionExists, err := exists(truncatePath(filepath.Join(connectionDirectory)))
	if e, _ := Fatal(err, "Parent Directory Error", "Parent Directory Check", truncatePath(filepath.Join(connectionDirectory))); e {
		if !pathConnectionExists {
			Log("Creating a connection:", directory)
			os.Mkdir(truncatePath(filepath.Join(connectionDirectory, directory)), 0755)
			return true
		} else {
			os.RemoveAll(truncatePath(filepath.Join(connectionDirectory, directory)))
			Log("Creating a connection:", directory)
			os.Mkdir(truncatePath(filepath.Join(connectionDirectory, directory)), 0755)
			return true
		}
	} else {
		return false
	}
}

func checkFolder(directory string) bool {
	pathConnectionExists, err := exists(truncatePath(filepath.Join(connectionDirectory)))
	if e, _ := Fatal(err, "Parent Directory Error", "Parent Directory Check", truncatePath(filepath.Join(connectionDirectory))); e {
		if !pathConnectionExists {
			Log("Creating a connection:", directory)
			os.MkdirAll(truncatePath(filepath.Join(connectionDirectory, directory)), 0755)
			return true
		}else{
            return false
        }
	} else {
		return false
	}
}

func setupFiFo(directory, path string) (string, *os.File, error) {
	mkPath := filepath.Join(connectionDirectory, directory, truncatePath(path))
	pathExists, pathErr := exists(mkPath)
	if e, c := Fatal(pathErr, "File Check Error", "File Check", mkPath); e {
		if !pathExists {
			mkErr := syscall.Mkfifo(mkPath, 0755)
			Log("Preparing to create Pipe:", mkPath)
			if f, d := Fatal(mkErr, "Pipe Creation Error", "Creating Pipe", mkPath); f {
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
	mkPath := filepath.Join(connectionDirectory, directory, truncatePath(path))
	_, pathErr := exists(mkPath)
	if e, c := Fatal(pathErr, "File Check Error", "File Check", mkPath); e {
		Log("Opening the Named Pipe as a Scanner...")
        retScanner := bufio.NewScanner(pipe)
        retScanner.Split(bufio.ScanLines)
        Log("Created a named Pipe for sending requests:", mkPath)
		return retScanner, nil
	}else{
        return nil, c
    }
}

