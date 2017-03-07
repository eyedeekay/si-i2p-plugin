package main

import (
	"github.com/cmotc/sam3"
        "bufio"
        "io"
        "log"
        "net"
)

type i2pHTTPProxy struct {
        listener        *net.TCPListener

        sentBytes       uint64
        recievedBytes   uint64

        localAddr       *net.TCPAddr
        String          string
        remoteAddr      []i2pHTTPTunnel
        localConnection io.ReadWriteCloser

        erred           bool
        errsig          chan bool
	Log             log.Logger
	OutputHex       bool
}

func (i2proxy *i2pHTTPProxy) RequestRemoteStream(i2paddr string) (*sam3.StreamSession) {
        var x int
        for index, remote := range i2proxy.remoteAddr {
                if i2paddr == remote.stringAddr { x = index; }
        }
        return i2proxy.remoteAddr[x].streamSession
}

func (i2proxy *i2pHTTPProxy) RequestLastRemoteStream() (*sam3.StreamSession) {
        return i2proxy.remoteAddr[len(i2proxy.remoteAddr)].streamSession
}

func (i2proxy *i2pHTTPProxy) RequestRemoteConnection(i2paddr string) (*sam3.SAMConn) {
        var x int
        for index, remote := range i2proxy.remoteAddr {
                if i2paddr == remote.stringAddr { x = index; }
        }
        return i2proxy.remoteAddr[x].remoteConnection
}

func (i2proxy *i2pHTTPProxy) RequestLastRemoteConnection() (*sam3.SAMConn) {
        return i2proxy.remoteAddr[len(i2proxy.remoteAddr)].remoteConnection
}

func (i2proxy *i2pHTTPProxy) RequestRemoteRead(i2paddr string, buf []byte) (int, error) {
        var x int
        for index, remote := range i2proxy.remoteAddr {
                if i2paddr == remote.stringAddr { x = index; }
        }
        return i2proxy.remoteAddr[x].Read(buf)
}

func (i2proxy *i2pHTTPProxy) TestRemoteRead(buf []byte) (int, error) {
        return i2proxy.remoteAddr[0].Read(buf)
}

func (i2proxy *i2pHTTPProxy) RequestTunnel(i2paddr string) (*sam3.SAMConn){
        i2proxy.remoteAddr = append(i2proxy.remoteAddr,
                *Newi2pHTTPTunnelFromString(i2proxy.localAddr,
                SamAddr,
                i2paddr))
        p("Set up i2p stream session with remote destination: ", i2paddr)
        return i2proxy.RequestRemoteConnection(i2paddr)
}

func (i2proxy *i2pHTTPProxy) RequestHalfOpenTunnel() (*sam3.SAMConn){
        i2proxy.remoteAddr = append(i2proxy.remoteAddr,
                *Newi2pHTTPTunnel(i2proxy.localAddr,
                SamAddr))
        return i2proxy.RequestLastRemoteConnection()
}

func (i2proxy *i2pHTTPProxy) RequestSomeTunnel(i2paddr string) (*sam3.SAMConn){
        var r *sam3.SAMConn
        if( i2paddr == "" ){
                p("No i2paddr, assuming a half-open tunnel")
                r = i2proxy.RequestHalfOpenTunnel()
        }else{
                p("Connecting to i2paddr: ", i2paddr)
                r = i2proxy.RequestTunnel(i2paddr)
        }
        return r
}

func (i2proxy *i2pHTTPProxy) RequestDestination(i2paddr string) (*sam3.SAMConn){
        var r *sam3.SAMConn
        var f = false
        if (i2paddr != ""){
                for _, remote := range i2proxy.remoteAddr {
                        if i2paddr == remote.stringAddr {
                                r = i2proxy.RequestRemoteConnection(i2paddr)
                                f = true
                        }
                }
                if(!f){
                        r = i2proxy.RequestSomeTunnel(i2paddr)
                }
        }else{
                r = i2proxy.RequestHalfOpenTunnel()
        }
        return r
}

func (i2proxy *i2pHTTPProxy) RequestPipe(src io.ReadWriter, i2paddr string) {
        var f = false
        if (i2paddr != ""){
                for _, remote := range i2proxy.remoteAddr {
                        if i2paddr == remote.stringAddr {
                                remote.pipe(*i2proxy)
                                f = true
                        }
                }
                if(!f){
                        i2proxy.RequestSomeTunnel(i2paddr)
                        i2proxy.RequestPipe(i2proxy.localConnection, i2paddr)
                }
        }else{
                i2proxy.RequestHalfOpenTunnel()
        }
}

func (i2proxy *i2pHTTPProxy) Starti2pHTTPProxy(){
        var tempErr error
        i2proxy.localConnection, tempErr        = i2proxy.listener.AcceptTCP()
        err("Failed not accepting local connections '%s'\n",
                "Accepting local connections on: localhost:4443\n",
                 tempErr)
        defer i2proxy.localConnection.Close();
        //i2proxy.RequestTunnel()
        //i2proxy.localConnection, tempErr = net.DialTCP("tcp", nil, )
        err("Initial Connection to the i2p tunnel failed %s",
                "Finally connected to i2p for this web site:",
                tempErr)
        defer i2proxy.localConnection.Close()
	//display both ends
	i2proxy.Log.Printf("Opened %s >>> %s", i2proxy.localAddr.String(),
                i2proxy.remoteAddr[0].String())
	//bidirectional copy
        go i2proxy.RequestPipe(i2proxy.localConnection, "")
        go i2proxy.RequestPipe(i2proxy.localConnection, "")
	//wait for close...
	<-i2proxy.errsig
	i2proxy.Log.Printf("Closed (%d bytes sent, %d bytes recieved)",
                i2proxy.sentBytes, i2proxy.recievedBytes)
}

func (i2proxy *i2pHTTPProxy) SetupHTTPProxy(proxAddrString string) (io.ReadWriteCloser) {
        i2proxy.String             = proxAddrString
        var tempErr error
        i2proxy.localAddr, tempErr    = net.ResolveTCPAddr("tcp", proxAddrString)
        err("Failed to resolve address for local proxy'%s'\n",
                "Started an http proxy.\n",
                tempErr)
        i2proxy.listener, tempErr     = net.ListenTCP("tcp", i2proxy.localAddr)
        err("Failed to set up TCP listener '%s'\n",
                "Started a tcp listener.\n",
                tempErr)
        return i2proxy.localConnection
}

/*Newi2pHTTPProxy Create a new local HTTP proxy to request sites and destinations from*/
func Newi2pHTTPProxy(proxAddrString string, samAddrString string, logAddrWriter *bufio.Writer) *i2pHTTPProxy{
        var temp i2pHTTPProxy
        temp.SetupHTTPProxy(proxAddrString)
        SamAddr = samAddrString
        //temp.RequestHalfOpenTunnel()
        temp.RequestDestination("zzz.i2p")
	tbuf                    := make([]byte, 4096)
	_, tempErr := temp.TestRemoteRead(tbuf)
        err("Failed to read from pipe '%s'\n", "Server received message from pipe", tempErr)
        temp.erred              = false
        temp.errsig             = make(chan bool)
	Log                = *log.New(logAddrWriter,
                "Stream Isolating Parent Proxy Reported an Error", 0)
        return &temp
}
