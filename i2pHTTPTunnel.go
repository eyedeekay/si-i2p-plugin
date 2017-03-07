package main

import (
	"github.com/cmotc/sam3"
        "io"
        "log"
        "net"
        "bytes"
)

/*sam is the working SAM bridge*/
var sam         *sam3.SAM

/*SamAddr is the string used to find the SAM bridge initially*/
var SamAddr     string
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
                p("" + fail, err)
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
        err("Generated per-site Keypair:", "Failed to generate Per-Site Keypair", tempErr)
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
		n, pipeErr := i2proxy.localConnection.Read(buff)
                err("Read Succeeded", "Read failed '%s'\n", pipeErr)
		b := buff[:n]
		//show output
		i2proxy.Log.Printf(dataDirection, n, "")
		i2proxy.Log.Panicf(byteFormat, b)
		//write out result
		n, pipeErr = i2ptun.dst.Write(b)
                err("Write Succeeded", "Write failed '%s'\n", pipeErr)
		if islocal {
			i2proxy.sentBytes += uint64(n)
		} else {
			i2proxy.recievedBytes += uint64(n)
		}
	}
}


/*SetupSAMBridge assures that variables related to the SAM bridge are set*/
func SetupSAMBridge(samAddrString string) (*sam3.SAM, string) {
        var tempErr error
        if( SamAddr == "" ) {
                SamAddr = samAddrString
                sam, tempErr          = sam3.NewSAM(samAddrString)
                err("Failed to set up i2p SAM Bridge connection '%s'\n",
                        "Connected to the SAM bridge",
                        tempErr)
                defer sam.Close()
        }else{
                if(sam == nil){
                        SamAddr = samAddrString
                        sam, tempErr          = sam3.NewSAM(samAddrString)
                        err("Failed to set up i2p SAM Bridge connection '%s'\n",
                                "Connected to the SAM bridge",
                                tempErr)
                }else{
                        p("SamAddr: is already set:", SamAddr)
                }
        }
        return sam, SamAddr
}

/*Newi2pHTTPTunnel Create a new half-open i2p tunnel to a non-specific destination*/
func Newi2pHTTPTunnel(laddr *net.TCPAddr, samAddrString string) * i2pHTTPTunnel {
        var temp i2pHTTPTunnel
        temp.stringAddr = ""
        sam, SamAddr = SetupSAMBridge(samAddrString)
        temp.keypair = temp.NewKeyPair()
        temp.remoteConnection = temp.LookupDestination()
        defer temp.remoteConnection.Close()
        p("Setting up the connection", temp.remoteConnection)
//        b                       := make([]byte, 4096)
//        buf                     := bytes.NewBuffer(b)
        //go temp.Write([]byte("Hello i2p!"))
        return &temp
}

/*Newi2pHTTPTunnelFromString Create a new destination for a specific site*/
func Newi2pHTTPTunnelFromString( laddr *net.TCPAddr, samAddrString string, lookupI2PAddr string ) * i2pHTTPTunnel {
        var temp i2pHTTPTunnel
        temp.stringAddr = lookupI2PAddr
        sam, SamAddr = SetupSAMBridge(samAddrString);
        temp.keypair = temp.NewKeyPair()
        temp.remoteConnection = temp.LookupDestination()
        defer temp.remoteConnection.Close()
        p("Setting up the connection", temp.remoteConnection)
//        b                       := make([]byte, 4096)
//        buf                     := bytes.NewBuffer(b)
        //go temp.Write([]byte("Hello i2p!"))
        return &temp
}

func (i2ptun *i2pHTTPTunnel) String() string{
        return i2ptun.keypair.String()
}

func (i2ptun *i2pHTTPTunnel) Write(stream []byte) (int, error){
        return i2ptun.remoteConnection.Write(stream)
}

func (i2ptun *i2pHTTPTunnel) Read(stream []byte) (int, error){
        return i2ptun.remoteConnection.Read(stream)
}

