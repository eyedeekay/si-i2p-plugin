package main

import (
	"github.com/cmotc/sam3"
        "io"
//      "flag"
//	"fmt"
        "net"
)

type i2pHTTPProxy struct {
        sentBytes       uint64
        recievedBytes   uint64
        localAddr       *net.TCPAddr
        remoteAddr      i2pHTTPTunnel
        lconn, rconn    io.ReadWriteCloser
        erred           bool
        errsig          chan bool

}

func (p *i2pHTTPProxy) err(s string, err error) {
	if p.erred {
		return
	}
	if err != io.EOF {
//		p.Log.Warn(s, err)
	}
	p.errsig <- true
	p.erred = true
}

func (i2proxy *i2pHTTPProxy) Starti2pHTTPProxy(){
        defer i2proxy.lconn.Close();
        var err error
//        i2proxy.remoteAddr.TCPAddr()
//        i2proxy.rconn, err = net.DialTCP("tcp", nil, )
        if err != nil {
                i2proxy.err("Initial Connection to the i2p tunnel failed %s", err)
                return
        }
        defer i2proxy.rconn.Close()
}

func Newi2pHTTPProxy(lconn *net.TCPConn, laddr *net.TCPAddr, samAddrString string) *i2pHTTPProxy{
        var temp i2pHTTPProxy
        temp.lconn              = lconn
        temp.localAddr          = laddr
        temp.remoteAddr,_       = *Newi2pHTTPTunnel(sam3.NewSAM(samAddrString), laddr)
        temp.erred              = false
        temp.errsig             = make(chan bool)
        return &temp
}





