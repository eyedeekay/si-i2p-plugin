package main

import (
    "flag"
	"log"
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
    useHttpProxy := *flag.Bool("http-proxy", false,
        "run the HTTP proxy" )
    Defwd, _ := os.Getwd()
    workDirectory   := *flag.String("directory", Defwd,
        "The working directory you want to use, defaults to current directory")
    address   := *flag.String("url", "http://i2p-projekt.i2p",
        "i2p URL you want to retrieve")

    log.SetOutput(os.Stdout)
    log.SetFlags(log.Lshortfile)

    log.Println( "Sam Address:", samAddrString )
    log.Println( "Sam Port:", samPortString )
    log.Println( "Proxy Address:", proxAddrString )
    log.Println( "Proxy Port:", proxPortString )
    log.Println( "Working Directory:", workDirectory )
    log.Println( "Debug mode:", debugConnection)
    log.Println( "Debug mode:", useHttpProxy)
    log.Println( "Initial URL:", address)

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

    if useHttpProxy {
        samProxy := createHttpProxy(proxAddrString, proxPortString, samStack)
        log.Println("Sam Proxy Started:" + samProxy.host)
    }

    time.Sleep(3000 * time.Millisecond)

    log.Println("Created client, starting loop...")
    for exit != true{
        samStack.readRequest()
        go samStack.writeResponses()
        go closeProxy(samStack)
        time.Sleep(10 * time.Millisecond)
    }

    samStack.cleanupClient()
}

func closeProxy(samStack samList){
    exit = samStack.readDelete()
}
