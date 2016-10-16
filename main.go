package main

import (
//	"github.com/cmotc/sam3"
//        "io"
        "flag"
//	"fmt"
//        "net"
)

func main(){
	samAddrPtr      := flag.String("bridge", "127.0.0.1:7656",
                "host:port of the SAM bridge")
        samAddrString   := *samAddrPtr
	proxAddrPtr     := flag.String("proxy", "127.0.0.1:4443",
                "host:port of the SAM bridge")
        proxAddrString  := *proxAddrPtr
        siProxy         := Newi2pHTTPProxy(proxAddrString, samAddrString)
        siProxy.Starti2pHTTPProxy()
}
