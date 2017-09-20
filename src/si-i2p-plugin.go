package main

import (
        //"bufio"
        "flag"
	"fmt"
        "os"
        "os/signal"
        "github.com/eyedeekay/gosam"
)

func main(){
	samAddrString   := *flag.String("bridge-addr", "127.0.0.1",
                "host: of the SAM bridge")
        samPortString   := *flag.String("bridge-port", "7656",
                ":port of the SAM bridge")
        proxAddrString  := *flag.String("proxy-addr", "127.0.0.1",
                "host: of the HTTP proxy")
        proxPortString  := *flag.String("proxy-port", "4443",
                ":port of the HTTP proxy")
        debugConnection := *flag.Bool("conn-debug", true,
                "Print connection debug info" )
        Defwd, _ := os.Getwd()
        workDirectory   := *flag.String("directory", Defwd,
                "The working directory you want to use, defaults to current directory")
        address   := *flag.String("url", "i2p-projekt.i2p",
                "i2p URL you want to retrieve")


        fmt.Println( "Sam Address:", samAddrString )
        fmt.Println( "Sam Port:", samPortString )
        fmt.Println( "Proxy Address:", proxAddrString )
        fmt.Println( "Proxy Port:", proxPortString )
        fmt.Println( "Working Directory:", workDirectory )
        fmt.Println( "Debug mode:", debugConnection)
        fmt.Println( "Initial URL:", address)

        goSam.ConnDebug = debugConnection
        var test samHttp

        c := make(chan os.Signal, 1)
        signal.Notify(c, os.Interrupt)
        go func(){
                for sig := range c {
                        if sig == os.Interrupt {
                                test.cleanupClient()
                        }
                }
        }()

        //samStack := createSamList(samAddrString, samPortString)

        test.createClient(samAddrString, samPortString, address)
        defer test.cleanupClient()

        fmt.Println("Created client, starting loop...")
        exit := false
        for exit != true{
                fmt.Println("checking for requests to send...")
                test.readRequest()
                //fmt.Println("checking for service address")
                //test.writeName()
                fmt.Println("checking contents of delete pipe")
                exit = test.readDelete()
        }
}
