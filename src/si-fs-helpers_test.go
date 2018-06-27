package dii2p

import (
	"testing"
)

func TestSetupFile(t *testing.T) {
	checkFolder("test")
	path, file, err := setupFile("test", "file")
	t.Log(path, file)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetupFiFo(t *testing.T) {
	checkFolder("test")
	path, pipe, err := setupFiFo("test", "fifo")
	t.Log(path, pipe)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetupScanner(t *testing.T) {
	checkFolder("test")
	path, pipe, err := setupFiFo("test", "fifo")
	t.Log(path, pipe)
	if err != nil {
		t.Fatal(err)
	}
	_, err = setupScanner("test", "scanner", pipe)
	if err != nil {
		t.Fatal(err)
	}
}
