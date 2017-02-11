package main

import (
//	"github.com/cmotc/sam3"
//        "io"
//        "io/ioutil"
        "bufio"
        "flag"
	"fmt"
//        "net"
        "os"
        "os/user"
)

var p = fmt.Println

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main(){
	samAddrPtr              := flag.String("bridge", "127.0.0.1:7656",
                "host:port of the SAM bridge")
        samAddrString           := *samAddrPtr
        p("Sam Bridge addr:port = ", samAddrString)
	proxAddrPtr             := flag.String("proxy", "127.0.0.1:4443",
                "host:port of the HTTP proxy")
        proxAddrString          := *proxAddrPtr
        p("Proxy addr:port = ", proxAddrString)
        usr, err                := user.Current()
        if err != nil {
                check(err)
        }
        p(usr)
        logPath                 := usr.HomeDir
        logPath         += "/.i2pstreams.log"
        logPathPtr              := flag.String("log", logPath,
                "path to save log files")
        logPathString           := *logPathPtr
        p("Log Path", logPathString)
        logPathPath, err        := os.Create(logPathString)
        if err != nil {
                check(err)
        }
        defer logPathPath.Close()
        logPathWriter           := bufio.NewWriter(logPathPath)
        for {
                siProxy                 := Newi2pHTTPProxy(proxAddrString,
                        samAddrString, logPathWriter)
                go siProxy.Starti2pHTTPProxy()
        }
}
