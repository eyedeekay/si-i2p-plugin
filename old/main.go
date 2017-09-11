package main

import (
        "bufio"
        "flag"
	"fmt"
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
        proxAddrStraightPtr             := flag.String("proxy-addr", "127.0.0.1",
                "host: of the HTTP proxy")
                proxAddrStraightString          := *proxAddrStraightPtr
        proxPortStraightPtr             := flag.String("proxy-port", "4443",
                ":port of the HTTP proxy")


                proxPortStraightString          := *proxPortStraightPtr
        proxAddrString := proxAddrStraightString + ":" + proxPortStraightString

        p("Proxy addr:port = ", proxAddrString)
        usr, err      := user.Current()
        if err != nil { check(err) }else{ p(usr) }
        logPath                 := usr.HomeDir
        logPath         += "/.i2pstreams.log"
        logPathPtr              := flag.String("log", logPath,
                "path to save log files")
        logPathString           := *logPathPtr
        p("Log Path", logPathString)
        logPathPath, err        := os.Create(logPathString)
        if err != nil { check(err) } else { defer logPathPath.Close() }
        logPathWriter           := bufio.NewWriter(logPathPath)
        siProxy                 := Newi2pHTTPProxy(proxAddrString,
        samAddrString, logPathWriter)
        for {
                go siProxy.Starti2pHTTPProxy()
        }
}
