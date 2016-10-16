package main

import (
	"github.com/cmotc/sam3"
//        "io"
        "net"
        "bytes"
)

type i2pHTTPTunnel struct {
        sam, _          *sam3.SAM
        remoteAddr      string
        keypair, _      sam3.I2PKeys
        stream, _       *sam3.StreamSession
        listener, _     *sam3.StreamListener   
        conn, _         net.Conn     
        buf             *bytes.Buffer
}

func Newi2pHTTPTunnel(samAddrString string, laddr *net.TCPAddr) * i2pHTTPTunnel {
        var temp i2pHTTPTunnel
        temp.sam, _             = sam3.NewSAM(samAddrString)
        temp.remoteAddr         = samAddrString
        temp.keypair, _         = temp.sam.NewKeys()
        temp.stream, _          = temp.sam.NewStreamSession("clientTun", temp.keypair, sam3.Options_Medium)
        temp.listener, _        = temp.stream.Listen()
        temp.conn, _            = temp.listener.Accept()
        b                       := make([]byte, 4096)
        buf                     = bytes.NewBuffer(b)
        return &temp
}

func (i2ptun *i2pHTTPTunnel) String() string{
        return i2ptun.remoteAddr
}

//func (i2ptun *i2pHTTPTunnel) i2pTunnel(){
        
//}