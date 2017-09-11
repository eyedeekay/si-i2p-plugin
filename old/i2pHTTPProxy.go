package main

import (
	"github.com/eyedeekay/sam3"
        "bufio"
        "io"
        "log"
        "net"
        "strings"
        "unicode"
)

type i2pHTTPProxy struct {
        localListener        *net.TCPListener

        sentBytes       uint64
        recievedBytes   uint64

        localAddr       *net.TCPAddr
        String          string
        curAddrString   string
        requestBytes     []byte
        remoteAddr      []i2pHTTPTunnel
        localConnection io.ReadWriteCloser

        erred           bool
        errsig          chan bool
	Log             log.Logger

}

/*sam is the working SAM bridge*/
var sam         *sam3.SAM

/*SamAddr is the string used to find the SAM bridge initially*/
var SamAddr     string

func (i2proxy *i2pHTTPProxy) RequestRemoteConnection(i2paddr string) (int){
        var x = -1
        if(len(i2proxy.remoteAddr) > 0){
                for index, remote := range i2proxy.remoteAddr {
                        if remote.StringCheck(i2paddr) {
                                x = index;
                        }
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
                p("Set up i2p stream session with remote destination: ", i2paddr)
        }
        r := i2proxy.remoteAddr[searched]
        return r
}

func (i2proxy *i2pHTTPProxy) RequestHalfOpenTunnel() (i2pHTTPTunnel){
        if(&i2proxy.remoteAddr != nil){
                if(len(i2proxy.remoteAddr) > 0){
                        if(i2proxy.remoteAddr[len(i2proxy.remoteAddr)].StringCheck("")){
                                i2proxy.remoteAddr = append(i2proxy.remoteAddr,
                                        *Newi2pHTTPTunnel(i2proxy.localAddr,
                                        SamAddr))
                        }
                }else{
                        i2proxy.remoteAddr = append(i2proxy.remoteAddr,
                                        *Newi2pHTTPTunnel(i2proxy.localAddr,
                                        SamAddr))
                }
        }else{
                i2proxy.remoteAddr = append(i2proxy.remoteAddr,
                        *Newi2pHTTPTunnel(i2proxy.localAddr,
                        SamAddr))
        }
        r := i2proxy.remoteAddr[len(i2proxy.remoteAddr) - 1]
        return r
}

func (i2proxy *i2pHTTPProxy) RequestDestination(ii2paddr string) (i2pHTTPTunnel){
        var r i2pHTTPTunnel
        //RequestRemoteConnection(i2paddr)
        i2paddr := ii2paddr
        p("Requesting Destination: ", i2paddr)
        if (i2paddr != ""){
                r = i2proxy.RequestTunnel(i2paddr)
        }else{
                r = i2proxy.RequestHalfOpenTunnel()
        }
        return r
}

func (i2proxy *i2pHTTPProxy) stripSpaces(str string) string {
    return strings.Map(func(r rune) rune {
        if unicode.IsSpace(r) {
            // if the character is a space, drop it
            return -1
        }
        // else keep it in the string
        return r
    }, str)
}

func (i2proxy *i2pHTTPProxy) AddrPipe() (){
	//directional copy (64k buffer)
        defer i2proxy.localConnection.Close()
        var request string
        var addr string
	buff := make([]byte, 0xffff)
	for {
		n, pipeErr := i2proxy.localConnection.Read(buff)
                err("Read failed '%s'\n", "Read Succeeded", pipeErr)
		b := buff[:n]
                request = string(b)
                preAddr := strings.SplitAfter(request, "\n")
                for _,subline := range preAddr{
                        if( strings.Contains(subline, "Host") ){
                                addr = strings.TrimLeft(
                                                strings.TrimLeft(
                                                        strings.TrimRight(subline, " "),
                                        " " ),
                                 "Host: ")
                        }
                }

                i2proxy.requestBytes = []byte(request)
                i2proxy.curAddrString = i2proxy.stripSpaces(addr)
                p(i2proxy.curAddrString)
                p(string(i2proxy.requestBytes))
                //i2proxy.WriteReadPipe()
                i2proxy.WritePipe(i2proxy.curAddrString, i2proxy.requestBytes)
                p("---")
                i2proxy.ReadPipe(i2proxy.curAddrString)
                p(n)
	}
        //return addr, []byte(request)
}

func (i2proxy *i2pHTTPProxy) WritePipe(i2paddr string, request []byte) {
        remote := i2proxy.RequestDestination(i2paddr)
        p("Opened ", i2proxy.localAddr.String(), " >>> ", remote.remoteI2PAddr.Base32())
        remote.wpipe(request, *i2proxy)
}

func (i2proxy *i2pHTTPProxy) ReadPipe(i2paddr string) {
        remote := i2proxy.RequestDestination(i2paddr)
        p("Opened ", remote.remoteI2PAddr.Base32(), " >>> ", i2proxy.localAddr.String())
        remote.rpipe(*i2proxy)
}

func (i2proxy *i2pHTTPProxy) WriteReadPipe(){
        i2proxy.WritePipe(i2proxy.curAddrString, i2proxy.requestBytes)
        p("----")
        i2proxy.ReadPipe(i2proxy.curAddrString)
}
func (i2proxy *i2pHTTPProxy) Starti2pHTTPProxy(){
//        defer i2proxy.localConnection.Close()
	//bidirectional copy
        go i2proxy.AddrPipe()
        //go i2proxy.ReadPipe(i2proxy.curAddrString)
        //go i2proxy.WriteReadPipe()
        //go i2proxy.bWritePipe(i2proxy.AddrPipe())
        //go i2proxy.WritePipe(get_me)
        //go i2proxy.ReadPipe(get_me)
	//wait for close...
	//<-i2proxy.errsig
	//i2proxy.Log.Printf("Closed (%d bytes sent, %d bytes recieved)",
                //i2proxy.sentBytes, i2proxy.recievedBytes)
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
        i2proxy.localListener, tempErr     = net.ListenTCP("tcp", i2proxy.localAddr)
        err("Failed to set up TCP listener.\n",
                "Started a tcp listener: " + i2proxy.String,
                tempErr)
        tempErr = nil
        i2proxy.localConnection, tempErr        = i2proxy.localListener.AcceptTCP()
        err("Failed not accepting local connections\n",
                "Accepting local connections on:",
                tempErr)
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
        }else{
                if(sam == nil){
                        p("SamAddr: is already set: ", SamAddr)
                        SamAddr = samAddrString
                        sam, tempErr          = sam3.NewSAM(samAddrString)
                        err("Failed to set up i2p SAM Bridge connection\n",
                                "Connected to the SAM bridge: " + SamAddr,
                                tempErr)
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
        temp.RequestDestination("");
        temp.erred              = false
        temp.errsig             = make(chan bool)
	Log                = *log.New(logAddrWriter,
                "Stream Isolating Parent Proxy Reported an Error", 0)
        return &temp
}
