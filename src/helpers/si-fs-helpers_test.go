package dii2phelper

import (
	"testing"
)

func TestSetupFile(t *testing.T) {
	SetupFolder("test")
	path, file, err := SetupFile("test", "file")
	t.Log(path, file)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetupFiFo(t *testing.T) {
	SetupFolder("test")
	path, pipe, err := SetupFiFo("test", "fifo")
	t.Log(path, pipe)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetupScanner(t *testing.T) {
	SetupFolder("test")
	path, pipe, err := SetupFiFo("test", "fifo")
	t.Log(path, pipe)
	if err != nil {
		t.Fatal(err)
	}
	_, err = SetupScanner("test", "scanner", pipe)
	if err != nil {
		t.Fatal(err)
	}
}
