
package main

import (
//	"github.com/cmotc/sam3"
        "bufio"
        "io"
//        "io/ioutil"
//      "flag"
//	"fmt"
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

	Matcher  func([]byte)
	Replacer func([]byte) []byte

	// Settings
	Nagles    bool
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
	//nagles?
	if i2proxy.Nagles {
		if conn, ok := i2proxy.lconn.(setNoDelayer); ok {
			conn.SetNoDelay(true)
		}
		if conn, ok := i2proxy.rconn.(setNoDelayer); ok {
			conn.SetNoDelay(true)
		}
	}

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

		//execute match
		if i2proxy.Matcher != nil {
			i2proxy.Matcher(b)
		}

		//execute replace
		if i2proxy.Replacer != nil {
			b = i2proxy.Replacer(b)
		}

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

func Newi2pHTTPProxy(proxAddrString string, samAddrString string, logAddrWriter *bufio.Writer) *i2pHTTPProxy{
        var temp i2pHTTPProxy
        temp.String             = proxAddrString
        var berr error
        temp.localAddr, berr   = net.ResolveTCPAddr("tcp", proxAddrString)
        if berr != nil {
                temp.err("Failed to resolve address for local proxy'%s'\n", berr)
        }
        temp.listener, berr    = net.ListenTCP("tcp", temp.localAddr)
        if berr != nil {
                temp.err("Failed to set up TCP listener '%s'\n", berr)
        }
        temp.lconn, berr      = temp.listener.AcceptTCP()
        if berr != nil {
                temp.err("Failed to set up connection '%s'\n", berr)
        }
        temp.remoteAddr         = append(temp.remoteAddr, *Newi2pHTTPTunnel(samAddrString, temp.localAddr))
        temp.erred              = false
        temp.errsig             = make(chan bool)
	temp.Log                = *log.New(logAddrWriter,
                "Stream Isolating Parent Proxy Reported an Error", 0)
        return &temp
}