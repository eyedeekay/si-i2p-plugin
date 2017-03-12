package main

import (
	"github.com/cmotc/sam3"
        "io"
        "log"
        "net"
        "bytes"
)

/*Log does logging*/
var Log         log.Logger
/*erred ...*/
var erred       bool
/*errsig ...*/
var errsig      chan bool

type i2pHTTPTunnel struct {
        keypair, _              sam3.I2PKeys
        streamSession, _        *sam3.StreamSession

        stringAddr              string
        remoteI2PAddr,_         sam3.I2PAddr
        remoteConnection        *sam3.SAMConn
        initialized             bool

        buf                     bytes.Buffer
        dst                     io.ReadWriter
}

func err(fail string, succeed string, tempErr error) {
        if(tempErr != nil){
                p("ERROR OCCURRED: " + fail + ". ", err)
                if erred {
                        return
                }
                if tempErr != io.EOF {
                        Log.Panicf(fail, err)
                }
                errsig <- true
                erred = true
        }else{
                p(succeed)
        }
}

func (i2ptun *i2pHTTPTunnel) NewKeyPair() sam3.I2PKeys{
        var tempErr error
        i2ptun.keypair, tempErr         = sam.NewKeys()
        p("keypair:", i2ptun.keypair)
        err("Failed to generate Per-Site Keypair",
                "Generated per-site Keypair:",
                tempErr)
        return i2ptun.keypair
}

func (i2ptun *i2pHTTPTunnel) NewStreamSession(lookupI2PAddr string) (*sam3.StreamSession, string){
        i2ptun.streamSession, _          = sam.NewStreamSession("clientTun",
                i2ptun.keypair,
                sam3.Options_Fat)
        p("Started Stream Session")
        i2ptun.stringAddr           = lookupI2PAddr
        return i2ptun.streamSession, i2ptun.stringAddr
}

func (i2ptun *i2pHTTPTunnel) LookupDestination() *sam3.SAMConn {
        p("Connecting to this address: %s\n", i2ptun.stringAddr)
        i2ptun.streamSession, i2ptun.stringAddr = i2ptun.NewStreamSession(i2ptun.stringAddr)
        i2ptun.remoteI2PAddr, _  = sam.Lookup(i2ptun.stringAddr)
        p("Dialing this connection.")
        i2ptun.remoteConnection, _           = i2ptun.streamSession.DialI2P(i2ptun.remoteI2PAddr)
        return i2ptun.remoteConnection
}

func (i2ptun *i2pHTTPTunnel) pipe(i2proxy i2pHTTPProxy) {
	islocal := i2proxy.localConnection == i2proxy.localConnection
        defer i2ptun.remoteConnection.Close()
	var dataDirection string
	if islocal {
		dataDirection = ">>> %d bytes sent%s"
	} else {
		dataDirection = "<<< %d bytes recieved%s"
	}
	//directional copy (64k buffer)
	buff := make([]byte, 0xffff)
	for {
		n, pipeErr := i2proxy.localConnection.Read(buff)
                err("Read Succeeded", "Read failed '%s'\n", pipeErr)
		b := buff[:n]
		//show output
		i2proxy.Log.Printf(dataDirection, n, "")
		//write out result
		n, pipeErr = i2ptun.dst.Write(b)
                err("Write Succeeded", "Write failed '%s'\n", pipeErr)
                i2proxy.sentBytes += uint64(n)
	}
}

func (i2ptun *i2pHTTPTunnel) rpipe(i2proxy i2pHTTPProxy) {
	islocal := i2proxy.localConnection == i2proxy.localConnection
        defer i2ptun.remoteConnection.Close()
	var dataDirection = "<<< %d bytes recieved%s"
	//directional copy (64k buffer)
	buff := make([]byte, 0xffff)
	for {
		n, pipeErr := i2ptun.remoteConnection.Read(buff)
                err("Read Succeeded", "Read failed '%s'\n", pipeErr)
		b := buff[:n]
		//show output
		i2proxy.Log.Printf(dataDirection, n, "")
		//write out result
		n, pipeErr = i2ptun.dst.Write(b)
                err("Write Succeeded", "Write failed '%s'\n", pipeErr)
                i2proxy.recievedBytes += uint64(n)
	}
}

/*Newi2pHTTPTunnel Create a new half-open i2p tunnel to a non-specific destination*/
func Newi2pHTTPTunnel(laddr *net.TCPAddr, samAddrString string) * i2pHTTPTunnel {
        var temp i2pHTTPTunnel
        temp.stringAddr = ""
        temp.keypair = temp.NewKeyPair()
        //temp.remoteConnection = temp.LookupDestination()
        //defer temp.remoteConnection.Close()
        //p("Setting up the connection", temp.remoteConnection)
        //b                       := make([]byte, 4096)
        //buf                     := bytes.NewBuffer(b)
        //go temp.Write([]byte("Hello i2p!"))
        return &temp
}

/*Newi2pHTTPTunnelFromString Create a new destination for a specific site*/
func Newi2pHTTPTunnelFromString( laddr *net.TCPAddr, samAddrString string, lookupI2PAddr string ) * i2pHTTPTunnel {
        var temp i2pHTTPTunnel
        temp.stringAddr = lookupI2PAddr
        temp.keypair = temp.NewKeyPair()
        temp.remoteConnection = temp.LookupDestination()
        //defer temp.remoteConnection.Close()
        p("Setting up the connection", temp.remoteConnection)
        //b                       := make([]byte, 4096)
        //buf                     := bytes.NewBuffer(b)
        //go temp.Write([]byte("Hello i2p!"))
        return &temp
}

func (i2ptun *i2pHTTPTunnel) String() string{
        return i2ptun.remoteI2PAddr.String()
}

func (i2ptun *i2pHTTPTunnel) Write(stream []byte) (int, error){
        p(stream)
        return i2ptun.remoteConnection.Write(stream)
}

func (i2ptun *i2pHTTPTunnel) Read(stream []byte) (int, error){
        p(stream)
        return i2ptun.remoteConnection.Read(stream)
}

