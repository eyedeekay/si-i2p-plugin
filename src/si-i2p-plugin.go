package main

import (
    "flag"
	"log"
    "os"
    "os/signal"
    "time"

    "github.com/eyedeekay/gosam"
    //"github.com/cryptix/goSam"
)

var exit bool = false

func main(){
	samAddrString   := flag.String("bridge-addr", "127.0.0.1",
        "host: of the SAM bridge")
    samPortString   := flag.String("bridge-port", "7656",
        ":port of the SAM bridge")
    proxAddrString  := flag.String("proxy-addr", "127.0.0.1",
        "host: of the HTTP proxy")
    proxPortString  := flag.String("proxy-port", "4443",
        ":port of the HTTP proxy")
    debugConnection := flag.Bool("conn-debug", true,
        "Print connection debug info" )
    useHttpProxy := flag.Bool("http-proxy", true,
        "run the HTTP proxy" )
    Defwd, _ := os.Getwd()
    workDirectory   := flag.String("directory", Defwd,
        "The working directory you want to use, defaults to current directory")
    address   := flag.String("url", "",
        "i2p URL you want to retrieve")

    flag.Parse()

    log.SetOutput(os.Stdout)
    log.SetFlags(log.Lshortfile)

    log.Println( "Sam Address:", *samAddrString )
    log.Println( "Sam Port:", *samPortString )
    log.Println( "Proxy Address:", *proxAddrString )
    log.Println( "Proxy Port:", *proxPortString )
    log.Println( "Working Directory:", *workDirectory )
    log.Println( "Debug mode:", *debugConnection)
    log.Println( "Using HTTP proxy:", *useHttpProxy)
    log.Println( "Initial URL:", *address)

    goSam.ConnDebug = *debugConnection

    var samProxies *samList
    samProxies = createSamList(*samAddrString, *samPortString, *address)

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    go func(){
        for sig := range c {
            if sig == os.Interrupt {
                samProxies.cleanupClient()
            }
        }
    }()

    httpUp := false

    if *useHttpProxy {
        if ! httpUp {
            samProxy := createHttpProxy(*proxAddrString, *proxPortString, samProxies, *address)
            log.Println("HTTP Proxy Started:" + samProxy.host)
            httpUp = true
        }
    }

    log.Println("Created client, starting loop...")

    for exit != true{
        samProxies.readRequest()
        go samProxies.writeResponses()
        go closeProxy(samProxies)

        time.Sleep(100 * time.Millisecond)
    }

    samProxies.cleanupClient()
}

func closeProxy(samProxies *samList){
    exit = samProxies.readDelete()
}
