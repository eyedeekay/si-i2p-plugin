package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	".."
	"github.com/eyedeekay/gosam"
	//"github.com/cryptix/goSam"
)

var exit bool = false

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
		"run the HTTP proxy(default true)")
	verboseLogging := flag.Bool("verbose", false,
		"Print connection debug info")
	Defwd, _ := os.Getwd()
	workDirectory := flag.String("directory", Defwd,
		"The working directory you want to use, defaults to current directory")
	address := flag.String("url", "",
		"i2p URL you want to retrieve")
	addressHelper := flag.String("addresshelper", "http://inr.i2p",
		"Jump/Addresshelper service you want to use")
	timeoutTime := flag.Int("timeout", 6,
		"Timeout duration in minutes(default six)")
	keepAlives := flag.Bool("disable-keepalives", false,
		"Disable keepalives(default false)")

	flag.Parse()

	log.SetOutput(os.Stdout)

	dii2p.Log("si-i2p-plugin.go Sam Address:", *samAddrString)
	dii2p.Log("si-i2p-plugin.go Sam Port:", *samPortString)
	dii2p.Log("si-i2p-plugin.go Proxy Address:", *proxAddrString)
	dii2p.Log("si-i2p-plugin.go Proxy Port:", *proxPortString)
	dii2p.Log("si-i2p-plugin.go Working Directory:", *workDirectory)
	dii2p.Log("si-i2p-plugin.go Addresshelper Services:", *addressHelper)
	log.Println("si-i2p-plugin.go Timeout Time:", *timeoutTime, "minutes")

	if !*keepAlives {
		dii2p.Log("si-i2p-plugin.go Keepalives Enabled")
	} else {
		dii2p.Log("si-i2p-plugin.go Keepalives Disabled")
	}

	if *debugConnection {
		dii2p.Log("si-i2p-plugin.go Debug mode: true")
	}
	if *verboseLogging {
		dii2p.Log("si-i2p-plugin.go Verbose mode: true")
	}
	if *useHttpProxy {
		dii2p.Log("si-i2p-plugin.go Using HTTP proxy: true")
	}
	dii2p.Log("si-i2p-plugin.go Initial URL:", *address)

	dii2p.Verbose = *verboseLogging

	goSam.ConnDebug = *debugConnection

	var samProxies *dii2p.SamList
	var samService *dii2p.SamServices

	samProxies = dii2p.CreateSamList(*samAddrString, *samPortString, *address, *timeoutTime, *keepAlives)
	samService = dii2p.CreateSamServiceList(*samAddrString, *samPortString)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			if sig == os.Interrupt {
				samProxies.CleanupClient()
			}
		}
	}()

	httpUp := false

	if *useHttpProxy {
		if !httpUp {
			samProxy := dii2p.CreateHttpProxy(*proxAddrString, *proxPortString, *address, *addressHelper, samProxies, *timeoutTime, *keepAlives)
			dii2p.Log("si-i2p-plugin.go HTTP Proxy Started:" + samProxy.Host)
			httpUp = true
		}
	}

	dii2p.Log("si-i2p-plugin.go Created client, starting loop...")

	for exit != true {
		go closeProxy(samProxies)
		go closeServices(samService)
		go samProxies.WriteResponses()
		//go samService.writeContents()
		go samService.ServiceRequest()
		samProxies.ReadRequest()

		time.Sleep(1 * time.Second)
	}

	samProxies.CleanupClient()
}

func closeProxy(samProxies *dii2p.SamList) {
	exit = samProxies.ReadDelete()
}

func closeServices(samServiceList *dii2p.SamServices) {
	exit = samServiceList.ReadDelete()
}