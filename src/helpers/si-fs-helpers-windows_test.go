// +build windows

package dii2phelper

import (
	"testing"
)

func TestSetupFile(t *testing.T) {
	checkFolder("test")
	path, file, err := setupFile("test", "file")
	t.dii2perrs.Log(path, file)
	if err != nil {
		t.dii2perrs.Fatal(err)
	}
}

func TestSetupFiFo(t *testing.T) {
	checkFolder("test")
	path, pipe, err := setupFiFo("test", "fifo")
	t.dii2perrs.Log(path, pipe)
	if err != nil {
		t.dii2perrs.Fatal(err)
	}
}

func TestSetupScanner(t *testing.T) {
	checkFolder("test")
	path, pipe, err := setupFiFo("test", "fifo")
	t.dii2perrs.Log(path, pipe)
	if err != nil {
		t.dii2perrs.Fatal(err)
	}
	_, err = setupScanner("test", "scanner", pipe)
	if err != nil {
		t.dii2perrs.Fatal(err)
	}
}
