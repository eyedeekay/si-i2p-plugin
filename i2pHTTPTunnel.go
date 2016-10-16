package main

import (
	"github.com/cmotc/sam3"
        "net"
)

type i2pHTTPTunnel struct {
        sam, _          sam3.SAM
        localAddr       *net.TCPAddr
        keypair, _      sam3.I2PKeys
}

func Newi2pHTTPTunnel(samb,_ *sam3.SAM, laddr *net.TCPAddr) * i2pHTTPTunnel {
        var temp i2pHTTPTunnel
        temp.sam        = *samb
        temp.localAddr  = laddr
        temp.keypair, _    = temp.sam.NewKeys()
        return &temp
}

//func (i2ptun *i2pHTTPTunnel) i2pTunnel(){
        
//}