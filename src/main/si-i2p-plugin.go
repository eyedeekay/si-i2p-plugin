package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	//"github.com/eyedeekay/si-i2p-plugin/src"
    "github.com/eyedeekay/jumphelper/src"
	".."
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
	socksAddrString := flag.String("socks-addr", "127.0.0.1",
		"host: of the SOCKS proxy")
	socksPortString := flag.String("socks-port", "4446",
		":port of the SOCKS proxy")
	addrHelperHostString := flag.String("ah-addr", "127.0.0.1",
		"host: of the SAM bridge")
	addrHelperPortString := flag.String("ah-port", "7054",
		":port of the SAM bridge")
	debugConnection := flag.Bool("conn-debug", false,
		"Print connection debug info")
	useHTTPProxy := flag.Bool("http-proxy", true,
		"run the HTTP proxy(default true)")
	useSOCKSProxy := flag.Bool("socks-proxy", false,
		"run the SOCKS proxy(default true)")
	verboseLogging := flag.Bool("verbose", false,
		"Print connection debug info")
	Defwd, _ := os.Getwd()
	workDirectory := flag.String("directory", Defwd,
		"The working directory you want to use, defaults to current directory")
	address := flag.String("url", "",
		"i2p URL you want to retrieve")
	addressHelper := flag.String("addresshelper", "http://inr.i2p",
		"Jump/Addresshelper service you want to use")
	destLifespan := flag.Int("lifespan", 12,
		"Lifespan of an idle i2p destination in minutes(default twelve)")
	timeoutTime := flag.Int("timeout", 6,
		"Timeout duration in minutes(default six)")
	tunnelLength := flag.Int("tunlength", 3,
		"Tunnel Length(default 3)")
	inboundTunnels := flag.Int("in-tunnels", 15,
		"Inbound Tunnel Count(default 15)")
	outboundTunnels := flag.Int("out-tunnels", 15,
		"Inbound Tunnel Count(default 15)")
	keepAlives := flag.Bool("disable-keepalives", false,
		"Disable keepalives(default false)")
	idleConns := flag.Int("idle-conns", 15,
		"Maximium idle connections per host(default 4)")
	inboundBackups := flag.Int("in-backups", 3,
		"Inbound Backup Count(default 3)")
	outboundBackups := flag.Int("out-backups", 3,
		"Inbound Backup Count(default 3)")
	internalAddressHelper := flag.Bool("internal-ah", true,
        "Use internal address helper")
    addressBook := flag.String("addressbook", "./addresses.csv",
		"path to local addressbook(default ./addresses.csv) (Unused without internal-ah)")
	//diskAvoidance := flag.Bool("avoidance", true,
	//  "Disk Avoidance Mode(default true)")

	flag.Parse()

	log.SetOutput(os.Stdout)

	dii2p.Log("si-i2p-plugin.go Sam Address:", *samAddrString)
	dii2p.Log("si-i2p-plugin.go Sam Port:", *samPortString)
	dii2p.Log("si-i2p-plugin.go HTTP Proxy Address:", *proxAddrString)
	dii2p.Log("si-i2p-plugin.go HTTP Proxy Port:", *proxPortString)
	dii2p.Log("si-i2p-plugin.go SOCKS Proxy Address:", *socksAddrString)
	dii2p.Log("si-i2p-plugin.go SOCKS Proxy Port:", *socksPortString)
	dii2p.Log("si-i2p-plugin.go Addresshelper Address:", *addrHelperHostString)
	dii2p.Log("si-i2p-plugin.go Addresshelper Port:", *addrHelperPortString)
	dii2p.Log("si-i2p-plugin.go Working Directory:", *workDirectory)
	dii2p.Log("si-i2p-plugin.go Addresshelper Services:", *addressHelper)
	dii2p.Log("si-i2p-plugin.go Timeout Time:", *timeoutTime, "minutes")
	dii2p.Log("si-i2p-plugin.go Tunnel Length:", *tunnelLength)
	dii2p.Log("si-i2p-plugin.go Inbound Tunnel Quantity:", *inboundTunnels)
	dii2p.Log("si-i2p-plugin.go Outbound Tunnel Quantity", *outboundTunnels)
	dii2p.Log("si-i2p-plugin.go Idle Tunnel Count:", *idleConns)
	dii2p.Log("si-i2p-plugin.go Inbound Backup Quantity:", *inboundBackups)
	dii2p.Log("si-i2p-plugin.go Outbound Backup Quantity", *outboundBackups)

	*useSOCKSProxy = false

    if *internalAddressHelper {
        dii2p.Log("si-i2p-plugin.go starting internal addresshelper with")
        jumphelper.NewService(*addrHelperHostString, *addrHelperPortString, *addressBook, *samAddrString, *samPortString)
    }

	if !*keepAlives {
		dii2p.Log("si-i2p-plugin.go Keepalives Enabled")
	} else {
		dii2p.Log("si-i2p-plugin.go Keepalives Disabled")
	}
	if *debugConnection {
		dii2p.DEBUG = *debugConnection
		dii2p.Log("si-i2p-plugin.go Debug mode: true")
	}
	if *verboseLogging {
		dii2p.Verbose = *verboseLogging
		dii2p.Log("si-i2p-plugin.go Verbose mode: true")
	}
	if *useHTTPProxy {
		dii2p.Log("si-i2p-plugin.go Using HTTP proxy: true")
	}
	dii2p.Log("si-i2p-plugin.go Initial URL:", *address)

	samProxies, err := dii2p.CreateSamList(
		dii2p.SetInitAddress(*address),
		dii2p.SetHost(*samAddrString),
		dii2p.SetPort(*samPortString),
		dii2p.SetTimeout(*timeoutTime),
		dii2p.SetKeepAlives(*keepAlives),
		dii2p.SetLifespan(*destLifespan),
		dii2p.SetTunLength(*tunnelLength),
		dii2p.SetInQuantity(*inboundTunnels),
		dii2p.SetOutQuantity(*outboundTunnels),
		dii2p.SetIdleConns(*idleConns),
		dii2p.SetInBackups(*inboundBackups),
		dii2p.SetOutBackups(*outboundBackups),
	)

	if err != nil {
		log.Fatal(err)
	}

	samService, err := dii2p.CreateSamServiceList(
		dii2p.SetServHost(*samAddrString),
		dii2p.SetServPort(*samPortString),
	)
	if err != nil {
		log.Fatal(err)
	}

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
	socksUp := false

	if *useHTTPProxy {
		if !httpUp {
			samProxy := dii2p.CreateHTTPProxy(*proxAddrString, *proxPortString, *address, *addrHelperHostString, *addrHelperPortString, *addressHelper, samProxies, *timeoutTime, *keepAlives)
			dii2p.Log("si-i2p-plugin.go HTTP Proxy Started:" + samProxy.Addr)
			httpUp = true
		}
	}

	if *useSOCKSProxy {
		if !socksUp {
			samProxy := dii2p.CreateSOCKSProxy(*proxAddrString, *proxPortString, *address, *addrHelperHostString, *addrHelperPortString, *addressHelper, samProxies, *timeoutTime, *keepAlives)
			dii2p.Log("si-i2p-plugin.go Socks Proxy Started:" + samProxy.Addr)
			socksUp = true
		}
	}

	dii2p.Log("si-i2p-plugin.go Created client, starting loop...")

	for exit != true {
		go closeProxy(samProxies)
		go closeServices(samService)
		go samProxies.WriteResponses()
		//go samService.writeContents()
		//go samService.ServiceRequest()
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
