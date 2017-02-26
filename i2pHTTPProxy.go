
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
        sam, _          *sam3.SAM
        test_keys       sam3.I2PKeys
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

func (i2proxy *i2pHTTPProxy) RequestRemoteListener(i2paddr string) (*sam3.StreamListener) {
        var x int
        for index, remote := range i2proxy.remoteAddr {
                if i2paddr == remote.stringAddr { x = index; }
        }
        return i2proxy.remoteAddr[x].listener
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
        if i2paddr != ""{
                p(i2paddr)
                var test_keys, test_err    = i2proxy.sam.NewKeys()
                i2proxy.remoteAddr = append(i2proxy.remoteAddr, *Newi2pHTTPTunnel(i2proxy.sam, i2proxy.localAddr, test_keys))
                tstream, test_err           := i2proxy.sam.NewStreamSession("Stream Isolation Test Proxy", test_keys, sam3.Options_Fat)
                if test_err != nil {
                        p("Failed to set up i2p stream session with test destination '%s'\n", test_err)
                }
                p("Set up i2p stream session with remote destination: ", tstream)
                return i2proxy.RequestRemoteListener(i2paddr)
        }else{
                p("No i2paddr, assuming a test tunnel")
                var test_keys, test_err    = i2proxy.sam.NewKeys()
                if test_err != nil {
                        p("Failed to set up new Destination Key for Test Tunnel '%s'\n", test_err)
                }
                p("Generated startup destination for test tunnel: ", test_keys)
                //tstream, test_err           := i2proxy.sam.NewStreamSession("Stream Isolation Test Proxy", test_keys, sam3.Options_Fat)
                i2proxy.remoteAddr = append(i2proxy.remoteAddr, *Newi2pHTTPTunnel(i2proxy.sam, i2proxy.localAddr, test_keys))
                if test_err != nil {
                        p("Failed to set up i2p stream session with test destination '%s'\n", test_err)
                }
                p("Set up i2p stream session with test destination: ", i2proxy.RequestRemoteStream(i2paddr))
                if test_err != nil {
                        i2proxy.err("Failed to set up i2p tunnel connection '%s'\n", test_err)
                } else { p("Started an i2p tunnel.") }
                return i2proxy.RequestRemoteListener(i2paddr)
        }
}

func (i2proxy *i2pHTTPProxy) Starti2pHTTPProxy(){
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

func (i2proxy *i2pHTTPProxy) pipe(src, dst io.ReadWriter) {
	islocal := src == i2proxy.lconn
	var dataDirection string
	if islocal {
		dataDirection = ">>> %d bytes sent%s"
	} else {
		dataDirection = "<<< %d bytes recieved%s"
	}
	var byteFormat string
	if i2proxy.OutputHex {
		byteFormat = "%x"
	} else {
		byteFormat = "%s"
	}
	//directional copy (64k buffer)
	buff := make([]byte, 0xffff)
	for {
		n, err := src.Read(buff)
		if err != nil {
			i2proxy.err("Read failed '%s'\n", err)
			return
		}
		b := buff[:n]
		//show output
		i2proxy.Log.Printf(dataDirection, n, "")
		i2proxy.Log.Panicf(byteFormat, b)
		//write out result
		n, err = dst.Write(b)
		if err != nil {
			i2proxy.err("Write failed '%s'\n", err)
			return
		}
		if islocal {
			i2proxy.sentBytes += uint64(n)
		} else {
			i2proxy.recievedBytes += uint64(n)
		}
	}
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
        //i2proxy.lconn, temp_err        = i2proxy.listener.AcceptTCP()
        return i2proxy.lconn
}

func (i2proxy *i2pHTTPProxy) SetupSAMBridge(samAddrString string) (*sam3.SAM) {
        var temp_err error
        i2proxy.sam, temp_err          = sam3.NewSAM(samAddrString)
        if temp_err != nil {
                i2proxy.err("Failed to set up i2p SAM Bridge connection '%s'\n", temp_err)
        }else{ p("Connected to the SAM bridge") }
        return i2proxy.sam
}

func Newi2pHTTPProxy(proxAddrString string, samAddrString string, logAddrWriter *bufio.Writer) *i2pHTTPProxy{
        var temp i2pHTTPProxy
        temp.SetupHTTPProxy(proxAddrString)
        temp.SetupSAMBridge(samAddrString)
        temp.RequestTunnel("")
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
