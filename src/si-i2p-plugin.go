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
var verbose bool = false

func main() {
	samAddrString := flag.String("bridge-addr", "127.0.0.1",
		"host: of the SAM bridge")
	samPortString := flag.String("bridge-port", "7656",
		":port of the SAM bridge")
	proxAddrString := flag.String("proxy-addr", "127.0.0.1",
		"host: of the HTTP proxy")
	proxPortString := flag.String("proxy-port", "4443",
		":port of the HTTP proxy")
	debugConnection := flag.Bool("conn-debug", false,
		"Print connection debug info")
	useHttpProxy := flag.Bool("http-proxy", true,
		"run the HTTP proxy")
	verboseLogging := flag.Bool("verbose", true,
		"Print connection debug info")
	Defwd, _ := os.Getwd()
	workDirectory := flag.String("directory", Defwd,
		"The working directory you want to use, defaults to current directory")
	address := flag.String("url", "",
		"i2p URL you want to retrieve")

	flag.Parse()

	log.SetOutput(os.Stdout)
	//log.SetFlags(log.Lshortfile)

	Log("si-i2p-plugin.go Sam Address:", *samAddrString)
	Log("si-i2p-plugin.go Sam Port:", *samPortString)
	Log("si-i2p-plugin.go Proxy Address:", *proxAddrString)
	Log("si-i2p-plugin.go Proxy Port:", *proxPortString)
	Log("si-i2p-plugin.go Working Directory:", *workDirectory)
	log.Println("si-i2p-plugin.go Debug mode:", *debugConnection)
	log.Println("si-i2p-plugin.go Verbose mode:", *verboseLogging)
	log.Println("si-i2p-plugin.go Using HTTP proxy:", *useHttpProxy)
	Log("si-i2p-plugin.go Initial URL:", *address)

	verbose = *verboseLogging

	goSam.ConnDebug = *debugConnection

	var samProxies *samList
    var samService *samServices

	samProxies = createSamList(*samAddrString, *samPortString, *address)
    samService = createSamServiceList(*samAddrString, *samPortString)
    samService.run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			if sig == os.Interrupt {
				samProxies.cleanupClient()
			}
		}
	}()

	httpUp := false

	if *useHttpProxy {
		if !httpUp {
			samProxy := createHttpProxy(*proxAddrString, *proxPortString, samProxies, *address)
			Log("si-i2p-plugin.go HTTP Proxy Started:" + samProxy.host)
			httpUp = true
		}
	}

	Log("si-i2p-plugin.go Created client, starting loop...")

	for exit != true {
		go closeProxy(samProxies)
		go samProxies.writeResponses()
		samProxies.readRequest()

		time.Sleep(100 * time.Millisecond)
	}

	samProxies.cleanupClient()
}

func closeProxy(samProxies *samList) {
	exit = samProxies.readDelete()
}

func Log(msg ...string) {
	if verbose {
		log.Println("LOG: ", msg)
	}
}

func LogA(msg []string) {
	if verbose {
		log.Println("LOG: ", msg)
	}
}

func Warn(err error, errmsg string, msg ...string) (bool, error) {
	LogA(msg)
	if err != nil {
		log.Println("WARN: ", err)
		return false, nil
	}
	return true, nil
}

func Fatal(err error, errmsg string, msg ...string) (bool, error) {
	LogA(msg)
	if err != nil {
		log.Fatal("FATAL: ", errmsg, err)
		return false, err
	}
	return true, nil
}
