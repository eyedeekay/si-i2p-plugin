package main

import (
	"github.com/cmotc/sam3"
        "io"
        "log"
        "net"
        "bytes"
)

var sam         *sam3.SAM
var SamAddr     string
var Log         log.Logger
var erred           bool
var errsig          chan bool

type i2pHTTPTunnel struct {
        keypair, _      sam3.I2PKeys
        stream, _       *sam3.StreamSession
        remoteI2PAddr,_ sam3.I2PAddr
        iconn, _        *sam3.SAMConn
        initialized     bool
        lconn, _        net.Conn
        listener, _     *sam3.StreamListener
        buf             bytes.Buffer
        stringAddr      string


}

func err(fail string, succeed string, the_err error) {
        if(the_err != nil){
                p("" + fail, err)
                if erred {
                        return
                }
                if the_err != io.EOF {
                        Log.Panicf(fail, err)
                }
                errsig <- true
                erred = true
        }else{
                p(succeed)
        }
}

func SetupSAMBridge(samAddrString string) (*sam3.SAM, string) {
        var temp_err error
        if( SamAddr == "" ) {
                sam, temp_err          = sam3.NewSAM(samAddrString)
                err("Failed to set up i2p SAM Bridge connection '%s'\n",
                        "Connected to the SAM bridge",
                        temp_err)
        }else{
                p("SamAddr: %s is already set\n", SamAddr)
        }
        return sam, SamAddr
}

func Newi2pHTTPTunnel(laddr *net.TCPAddr, samAddrString string) * i2pHTTPTunnel {
        var temp i2pHTTPTunnel
        temp.keypair, _         = sam.NewKeys()
        p("Per-Site Keypair: ", temp.keypair)
        temp.stream, _          = sam.NewStreamSession("clientTun", temp.keypair, sam3.Options_Fat)
        p("Started Stream Session")
        temp.stringAddr         = ""
        //p("Connecting to this address: ", temp.stringAddr)
        //temp.remoteI2PAddr, _   = sam.Lookup(raddr.String())
        //p("Connecting to this site: ", raddr.String())
        //temp.iconn, _           = temp.stream.DialI2P(temp.remoteI2PAddr)
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

func Newi2pHTTPTunnelFromString( laddr *net.TCPAddr, samAddrString string, raddr string ) * i2pHTTPTunnel {
        var temp i2pHTTPTunnel
        sam, SamAddr = SetupSAMBridge(samAddrString);
        temp.stringAddr           = raddr
        temp.keypair, _         = sam.NewKeys()
        p("Per-Site Keypair: ", temp.keypair)
        temp.stream, _          = sam.NewStreamSession("clientTun", temp.keypair, sam3.Options_Fat)
        p("Started Stream Session")
        temp.stringAddr         = raddr
        p("Connecting to this address: ", temp.stringAddr)
        temp.remoteI2PAddr, _   = sam.Lookup(raddr)
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

func (i2ptun *i2pHTTPTunnel) Read(stream []byte) (int, error){
//        buf     := bytes.NewBuffer(stream)
        return i2ptun.iconn.Read(stream)
}

