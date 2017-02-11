package main

import (
	"github.com/cmotc/sam3"
//        "io"
        "net"
        "bytes"
)

type i2pHTTPTunnel struct {
        keypair, _      sam3.I2PKeys
        stream, _       *sam3.StreamSession
        remoteI2PAddr,_ sam3.I2PAddr
        iconn, _        *sam3.SAMConn
        initialized     bool
        lconn, _        net.Conn
        listener, _     *sam3.StreamListener
        buf             bytes.Buffer
}

func Newi2pHTTPTunnel(insam *sam3.SAM, laddr *net.TCPAddr, raddr sam3.I2PKeys ) * i2pHTTPTunnel {
        var temp i2pHTTPTunnel
        temp.keypair, _         = insam.NewKeys()
        p("Per-Site Keypair: ", temp.keypair)
        temp.stream, _          = insam.NewStreamSession("clientTun", temp.keypair, sam3.Options_Fat)
        p("Started Stream Session")
        temp.remoteI2PAddr, _   = insam.Lookup(raddr.String())
        p("Connecting to this site: ", raddr.String())
        temp.iconn, _           = temp.stream.DialI2P(temp.remoteI2PAddr)
        p("Dialing this connection.")
        temp.listener, _        = temp.stream.Listen()
        p("Setting up the per-site listener", temp.listener)
        temp.lconn, _            = temp.listener.Accept()
        p("Setting up the connection", temp.lconn)
//        b                       := make([]byte, 4096)
//        buf                     := bytes.NewBuffer(b)
        go temp.Write([]byte("Hello i2p!"))
        return &temp
}

func Newi2pHTTPTunnelFromString(insam *sam3.SAM, laddr *net.TCPAddr, raddr string ) * i2pHTTPTunnel {
        var temp i2pHTTPTunnel
        temp.keypair, _         = insam.NewKeys()
        p("Per-Site Keypair: ", temp.keypair)
        temp.stream, _          = insam.NewStreamSession("clientTun", temp.keypair, sam3.Options_Fat)
        p("Started Stream Session")
        temp.remoteI2PAddr, _   = insam.Lookup(raddr)
        p("Connecting to this site: ", raddr)
        temp.iconn, _           = temp.stream.DialI2P(temp.remoteI2PAddr)
        p("Dialing this connection: ", temp.iconn)
        temp.listener, _        = temp.stream.Listen()
        p("Setting up the per-site listener", temp.listener)
        temp.lconn, _            = temp.listener.Accept()
        p("Setting up the connection", temp.lconn)
//        b                       := make([]byte, 4096)
//        buf                     := bytes.NewBuffer(b)
        go temp.Write([]byte("Hello i2p!"))
        return &temp
}

func (i2ptun *i2pHTTPTunnel) String() string{
        return i2ptun.keypair.String()
}

func (i2ptun *i2pHTTPTunnel) Write(stream []byte) (int, error){
//        buf     := bytes.NewBuffer(stream)
        return i2ptun.iconn.Write(stream)
}
