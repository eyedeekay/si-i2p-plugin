
package main

import (
	"github.com/cmotc/sam3"
        "bufio"
        "io"
//        "io/ioutil"
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
        lconn, rconn    io.ReadWriteCloser

        erred           bool
        errsig          chan bool
	Log       log.Logger
	OutputHex bool
}

func (i2proxy *i2pHTTPProxy) err(s string, err error) {
	if i2proxy.erred {
		return
	}
	if err != io.EOF {
		i2proxy.Log.Panicf(s, err)
	}
        if err != nil {
                p("" + s, err)
        }else{
                p(s)
        }
	i2proxy.errsig <- true
	i2proxy.erred = true
}

type setNoDelayer interface {
	SetNoDelay(bool) error
}

func (i2proxy *i2pHTTPProxy) RequestRemoteStream(i2paddr string) (*sam3.StreamSession) {
        var x int
        for index, remote := range i2proxy.remoteAddr {
                if i2paddr == remote.stringAddr { x = index; }
        }
        return i2proxy.remoteAddr[x].stream
}

func (i2proxy *i2pHTTPProxy) RequestLastRemoteStream() (*sam3.StreamSession) {
        return i2proxy.remoteAddr[len(i2proxy.remoteAddr)].stream
}

func (i2proxy *i2pHTTPProxy) RequestRemoteListener(i2paddr string) (*sam3.StreamListener) {
        var x int
        for index, remote := range i2proxy.remoteAddr {
                if i2paddr == remote.stringAddr { x = index; }
        }
        return i2proxy.remoteAddr[x].listener
}

func (i2proxy *i2pHTTPProxy) RequestLastRemoteListener() (*sam3.StreamListener) {
        return i2proxy.remoteAddr[len(i2proxy.remoteAddr)].listener
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

func (i2proxy *i2pHTTPProxy) RequestTunnel(i2paddr string) (*sam3.StreamListener){
        i2proxy.remoteAddr = append(i2proxy.remoteAddr, *Newi2pHTTPTunnelFromString(i2proxy.localAddr, SamAddr, i2paddr))
        p("Set up i2p stream session with remote destination: ", i2paddr)
        return i2proxy.RequestRemoteListener(i2paddr)
}

func (i2proxy *i2pHTTPProxy) RequestHalfOpenTunnel() (*sam3.StreamListener){
        i2proxy.remoteAddr = append(i2proxy.remoteAddr, *Newi2pHTTPTunnel(i2proxy.localAddr, SamAddr))
        return i2proxy.RequestLastRemoteListener()
}

func (i2proxy *i2pHTTPProxy) RequestSomeTunnel(i2paddr string) (*sam3.StreamListener){
        if( i2paddr == "" ){
                p("No i2paddr, assuming a half-open tunnel")
                return i2proxy.RequestHalfOpenTunnel()
        }else{
                p("Connecting to i2paddr: ", i2paddr)
                return i2proxy.RequestTunnel(i2paddr)
        }
}

func (i2proxy *i2pHTTPProxy) Starti2pHTTPProxy(){
        var temp_err error
        i2proxy.lconn, temp_err        = i2proxy.listener.AcceptTCP()
        if temp_err != nil {
                i2proxy.err("Failed not accepting local connections '%s'\n", temp_err)
        } else{ p("Accepting local connections on: localhost:4443\n") }
        defer i2proxy.lconn.Close();
        var err error
//        i2proxy.RequestTunnel()
//        i2proxy.remoteAddr[0].TCPAddr()
//        i2proxy.rconn, err = net.DialTCP("tcp", nil, )
        if err != nil {
                i2proxy.err("Initial Connection to the i2p tunnel failed %s", err)
                return
        }
        p("Finally connected to i2p for this web site:")
        defer i2proxy.rconn.Close()
	//display both ends
	i2proxy.Log.Printf("Opened %s >>> %s", i2proxy.localAddr.String(),
                i2proxy.remoteAddr[0].String())
	//bidirectional copy
	go i2proxy.pipe(i2proxy.lconn, i2proxy.rconn)
	go i2proxy.pipe(i2proxy.rconn, i2proxy.lconn)
	//wait for close...
	<-i2proxy.errsig
	i2proxy.Log.Printf("Closed (%d bytes sent, %d bytes recieved)",
                i2proxy.sentBytes, i2proxy.recievedBytes)
}

func (i2proxy *i2pHTTPProxy) SetupHTTPProxy(proxAddrString string) (io.ReadWriteCloser) {
        i2proxy.String             = proxAddrString
        var temp_err error
        i2proxy.localAddr, temp_err    = net.ResolveTCPAddr("tcp", proxAddrString)
        if temp_err != nil {
                i2proxy.err("Failed to resolve address for local proxy'%s'\n", temp_err)
        } else { p("Started an http proxy.") }
        i2proxy.listener, temp_err     = net.ListenTCP("tcp", i2proxy.localAddr)
        if temp_err != nil {
                i2proxy.err("Failed to set up TCP listener '%s'\n", temp_err)
        } else { p("Started a tcp listener.") }
        return i2proxy.lconn
}

func Newi2pHTTPProxy(proxAddrString string, samAddrString string, logAddrWriter *bufio.Writer) *i2pHTTPProxy{
        var temp i2pHTTPProxy
        temp.SetupHTTPProxy(proxAddrString)
        //temp.SetupSAMBridge(samAddrString)
        SamAddr = samAddrString
        temp.RequestHalfOpenTunnel()
	tbuf                    := make([]byte, 4096)
	tn, test_err := temp.TestRemoteRead(tbuf)
        if test_err != nil {
                temp.err("Failed to read from pipe '%s'\n", test_err)
        } else{ p("Server received: " + string(tbuf[:tn])) }
        temp.erred              = false
        temp.errsig             = make(chan bool)
	temp.Log                = *log.New(logAddrWriter,
                "Stream Isolating Parent Proxy Reported an Error", 0)
        return &temp
}
