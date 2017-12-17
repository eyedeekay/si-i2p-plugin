package main

import (
        "flag"
	"fmt"
        "os"
        "os/signal"
        "time"
        "github.com/eyedeekay/gosam"
)

var exit bool = false

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
        address   := *flag.String("url", "http://i2p-projekt.i2p",
                "i2p URL you want to retrieve")


        fmt.Println( "Sam Address:", samAddrString )
        fmt.Println( "Sam Port:", samPortString )
        fmt.Println( "Proxy Address:", proxAddrString )
        fmt.Println( "Proxy Port:", proxPortString )
        fmt.Println( "Working Directory:", workDirectory )
        fmt.Println( "Debug mode:", debugConnection)
        fmt.Println( "Initial URL:", address)

        goSam.ConnDebug = debugConnection

        samStack := createSamList(samAddrString, samPortString, address)

        c := make(chan os.Signal, 1)
        signal.Notify(c, os.Interrupt)
        go func(){
                for sig := range c {
                        if sig == os.Interrupt {
                                samStack.cleanupClient()
                        }
                }
        }()

        fmt.Println("Created client, starting loop...")
        for exit != true{
                go samStack.readRequest()
                samStack.writeResponses()
                go closeProxy(samStack)
                time.Sleep(10 * time.Millisecond)
        }

        samStack.cleanupClient()
}

func closeProxy(samStack samList){
        exit = samStack.readDelete()
}
