package main

import (
	"github.com/cmotc/sam3"
        "bufio"
        "io"
        "log"
        "net"
)

type i2pHTTPProxy struct {
        localListener        *net.TCPListener

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

/*sam is the working SAM bridge*/
var sam         *sam3.SAM

/*SamAddr is the string used to find the SAM bridge initially*/
var SamAddr     string

func (i2proxy *i2pHTTPProxy) RequestRemoteConnection(i2paddr string) (int){
        var x = -1
        if(len(i2proxy.remoteAddr) > 0){
                for index, remote := range i2proxy.remoteAddr {
                        if i2paddr == remote.stringAddr { x = index; }
                }
        }
        return x
}

func (i2proxy *i2pHTTPProxy) RequestTunnel(i2paddr string) (i2pHTTPTunnel){
        searched := i2proxy.RequestRemoteConnection(i2paddr);
        if(searched < 0){
                i2proxy.remoteAddr = append(i2proxy.remoteAddr,
                        *Newi2pHTTPTunnelFromString(i2proxy.localAddr,
                        SamAddr,
                        i2paddr))
                searched = i2proxy.RequestRemoteConnection(i2paddr);
        }
        p("Set up i2p stream session with remote destination: ", i2paddr)
        r := i2proxy.remoteAddr[searched]
        return r
}

func (i2proxy *i2pHTTPProxy) RequestHalfOpenTunnel() (i2pHTTPTunnel){
        if(&i2proxy.remoteAddr != nil){
                p(len(i2proxy.remoteAddr))
                if(len(i2proxy.remoteAddr) > 0){
                        if(i2proxy.remoteAddr[len(i2proxy.remoteAddr)].stringAddr != ""){
                                i2proxy.remoteAddr = append(i2proxy.remoteAddr,
                                        *Newi2pHTTPTunnel(i2proxy.localAddr,
                                        SamAddr))
                        }
                }else{
                        i2proxy.remoteAddr = append(i2proxy.remoteAddr,
                                        *Newi2pHTTPTunnel(i2proxy.localAddr,
                                        SamAddr))
                }
        }
        r := i2proxy.remoteAddr[i2proxy.RequestRemoteConnection("")]
        return r
}

func (i2proxy *i2pHTTPProxy) RequestDestination(i2paddr string) (i2pHTTPTunnel){
        var r i2pHTTPTunnel
        //RequestRemoteConnection(i2paddr)
        if (i2paddr != ""){
                r = i2proxy.RequestTunnel(i2paddr)
        }else{
                r = i2proxy.RequestHalfOpenTunnel()
        }
        return r
}

func (i2proxy *i2pHTTPProxy) RequestPipe(src io.ReadWriter, i2paddr string) {
        remote := i2proxy.RequestDestination(i2paddr)
        remote.pipe(*i2proxy)
}

func (i2proxy *i2pHTTPProxy) Starti2pHTTPProxy(){
        var tempErr error
	//bidirectional copy
        i2proxy.localConnection, tempErr        = i2proxy.localListener.AcceptTCP()
        err("Failed not accepting local connections\n",
                "Accepting local connections on:",
                tempErr)
        defer i2proxy.localConnection.Close()
        //i2proxy.localConnection, _        = i2proxy.localListener.AcceptTCP()
        //go i2proxy.RequestPipe(i2proxy.localConnection, "")
        i2proxy.RequestPipe(i2proxy.localConnection, "")
        //go i2proxy.RequestPipe(i2proxy.localConnection, "i2p-projekt.i2p")
	//wait for close...
	<-i2proxy.errsig
	i2proxy.Log.Printf("Closed (%d bytes sent, %d bytes recieved)",
                i2proxy.sentBytes, i2proxy.recievedBytes)
}

func (i2proxy *i2pHTTPProxy) SetupHTTPListener(proxAddrString string) (*net.TCPListener) {
        i2proxy.String             = proxAddrString
        var tempErr error
        tempErr = nil
        i2proxy.localAddr, tempErr    = net.ResolveTCPAddr("tcp", i2proxy.String)
        err("Failed to resolve address for local proxy\n",
                "Resolved address for local proxy: " + i2proxy.String,
                tempErr)
        tempErr = nil
        //i2proxy.localListener, tempErr     = net.ListenTCP("tcp", i2proxy.localAddr)
        i2proxy.localListener, _     = net.ListenTCP("tcp", i2proxy.localAddr)
        err("Failed to set up TCP listener.\n",
                "Started a tcp listener: " + i2proxy.String,
                tempErr)
        tempErr = nil
        return i2proxy.localListener
}

/*SetupSAMBridge assures that variables related to the SAM bridge are set*/
func SetupSAMBridge(samAddrString string) (*sam3.SAM, string) {
        var tempErr error
        if( SamAddr == "" ) {
                SamAddr = samAddrString
                sam, tempErr          = sam3.NewSAM(samAddrString)
                err("Failed to set up i2p SAM Bridge connection\n",
                        "Connected to the SAM bridge: " + SamAddr,
                        tempErr)
                //defer sam.Close()
        }else{
                if(sam == nil){
                        p("SamAddr: is already set: ", SamAddr)
                        SamAddr = samAddrString
                        sam, tempErr          = sam3.NewSAM(samAddrString)
                        err("Failed to set up i2p SAM Bridge connection\n",
                                "Connected to the SAM bridge: " + SamAddr,
                                tempErr)
                        //defer sam.Close()
                }else{
                        p("SamAddr: is already set: ", SamAddr)
                        p("SAM Bridge is already connected.\n")
                }
        }
        return sam, SamAddr
}

/*Newi2pHTTPProxy Create a new local HTTP proxy to request sites and destinations from*/
func Newi2pHTTPProxy(proxAddrString string, samAddrString string, logAddrWriter *bufio.Writer) *i2pHTTPProxy{
        var temp i2pHTTPProxy
        sam, SamAddr = SetupSAMBridge(samAddrString)
        temp.localListener = temp.SetupHTTPListener(proxAddrString)
        temp.RequestPipe(temp.localConnection, "zzz.i2p");
        temp.erred              = false
        temp.errsig             = make(chan bool)
	Log                = *log.New(logAddrWriter,
                "Stream Isolating Parent Proxy Reported an Error", 0)
        return &temp
}
