package main

import (
    //"bufio"
    //"fmt"
	//"io"
	//"log"
	"net/http"
    //"os"
    //"path/filepath"
    //"strings"
    //"strconv"
    //"syscall"
    //"net/url"

	//"github.com/eyedeekay/gosam"
)

type samHttpProxy struct {
    http *http.Client
    err error

    transport *http.Transport
    host string
}


func (proxy *samHttpProxy) prepare(){

}

func newHttpProxy() samHttpProxy {
    var samProxy samHttpProxy
    //samProxy.
    return samProxy
}
