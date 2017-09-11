package main

import (
	"github.com/eyedeekay/sam3"
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

        remoteI2PAddr, _         sam3.I2PAddr
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
        err("Failed to generate Per-Site Keypair",
                "Generated per-site Keypair:",
                tempErr)
        p(i2ptun.keypair)
        return i2ptun.keypair
}

func (i2ptun *i2pHTTPTunnel) NewStreamSession(lookupI2PAddr string) (*sam3.StreamSession){
        i2ptun.streamSession, _          = sam.NewStreamSession("clientTun",
                i2ptun.keypair,
                sam3.Options_Humongous)
        p("Started Stream Session", i2ptun.remoteI2PAddr.Base64())
        //i2ptun.stringAddr           = lookupI2PAddr
        //lookupDestination
        return i2ptun.streamSession
}

func (i2ptun *i2pHTTPTunnel) LookupDestination(stringAddr string) (*sam3.SAMConn, bool) {
        var tempErr error
        i2ptun.remoteI2PAddr, tempErr  = sam.Lookup(stringAddr)
        err("Failed to lookup address", "Looked up address", tempErr)
        p("Connecting to this address: ", stringAddr)
        i2ptun.streamSession = i2ptun.NewStreamSession(i2ptun.remoteI2PAddr.Base64())
        p("Dialing this connection. ", i2ptun.remoteI2PAddr.Base32())
        i2ptun.remoteConnection, tempErr           = i2ptun.streamSession.DialI2P(i2ptun.remoteI2PAddr)
        err("Failed to Dial Connection","Dialed Connection", tempErr)
        return i2ptun.remoteConnection, true
}

func (i2ptun *i2pHTTPTunnel) wpipe(buff []byte, i2proxy i2pHTTPProxy) {
        defer i2ptun.remoteConnection.Close()
        //defer i2proxy.localConnection.Close()
	//directional copy (64k buffer)
	//for {
		//write out result
                if(buff != nil){
                p("making i2p request:", string(buff))
                n, pipeErr := i2ptun.dst.Write(buff)
                err("Write failed '%s'\n", "Write Succeeded", pipeErr)
                i2proxy.sentBytes += uint64(n)
                }
	//}
}

func (i2ptun *i2pHTTPTunnel) rpipe(i2proxy i2pHTTPProxy) {
        defer i2ptun.remoteConnection.Close()
        defer i2proxy.localConnection.Close()
	var dataDirection = "<<< %d bytes recieved%s"
	//directional copy (64k buffer)
	buff := make([]byte, 0xffff)
	//for {
		n, pipeErr := i2ptun.remoteConnection.Read(buff)
                err("Read failed '%s'\n", "Read Succeeded", pipeErr)
		b := buff[:n]
		//show output
		//i2proxy.Log.Printf(dataDirection, n, "")
                p(dataDirection)
                p(b)
                p(n)
		//write out result
		n, pipeErr = i2proxy.localConnection.Write(b)
                err("Write failed '%s'\n", "Write Succeeded", pipeErr)
                i2proxy.recievedBytes += uint64(n)
	//}
}

/*Newi2pHTTPTunnel Create a new half-open i2p tunnel to a non-specific destination*/
func Newi2pHTTPTunnel(laddr *net.TCPAddr, samAddrString string) * i2pHTTPTunnel {
        var temp i2pHTTPTunnel
        temp.keypair = temp.NewKeyPair()
        p("Setting up the Half-Open connection")
        return &temp
}

/*Newi2pHTTPTunnelFromString Create a new destination for a specific site*/
func Newi2pHTTPTunnelFromString( laddr *net.TCPAddr, samAddrString string, lookupI2PAddr string ) * i2pHTTPTunnel {
        var temp i2pHTTPTunnel
        temp.keypair = temp.NewKeyPair()
        temp.remoteConnection, temp.initialized = temp.LookupDestination(lookupI2PAddr)
        p("Setting up the connection", lookupI2PAddr)
        return &temp
}

func (i2ptun *i2pHTTPTunnel) StringCheck(lookupI2PAddr string) bool{
        var t = false
        if(lookupI2PAddr != ""){
                temp, _ := sam.Lookup(lookupI2PAddr)
                //err("Failed to Look Up Destination: ", "Looked Up Destination: ", tempErr)
                p(lookupI2PAddr)
                if(!i2ptun.initialized){
                        if(i2ptun.remoteI2PAddr.String()==temp.String()){
                                t = true
                                p("Found Existing Destination: ", lookupI2PAddr)
                        }
                }else{
                        p("Found half-open Destination, Connecting:", lookupI2PAddr)
                        i2ptun.remoteConnection, i2ptun.initialized = i2ptun.LookupDestination(lookupI2PAddr)
                }
        }
        return t
}

func (i2ptun *i2pHTTPTunnel) Write(stream []byte) (int, error){
        p(stream)
        return i2ptun.remoteConnection.Write(stream)
}

func (i2ptun *i2pHTTPTunnel) Read(stream []byte) (int, error){
        p(stream)
        return i2ptun.remoteConnection.Read(stream)
}

