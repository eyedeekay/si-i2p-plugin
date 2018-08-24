package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"
)

import (
	"github.com/eyedeekay/jumphelper/src"
	"github.com/eyedeekay/samrtc/src"
	"github.com/eyedeekay/si-i2p-plugin/src"
	"github.com/eyedeekay/si-i2p-plugin/src/client"
	"github.com/eyedeekay/si-i2p-plugin/src/errors"
	"github.com/eyedeekay/si-i2p-plugin/src/server"
	//"github.com/crawshaw/littleboss"
)

var exit = false

var (
	samAddrString         = flag.String("bridge-addr", "127.0.0.1", "host: of the SAM bridge")
	samPortString         = flag.String("bridge-port", "7656", ":port of the SAM bridge")
	proxAddrString        = flag.String("proxy-addr", "127.0.0.1", "host: of the HTTP proxy")
	proxPortString        = flag.String("proxy-port", "4443", ":port of the HTTP proxy")
	socksAddrString       = flag.String("socks-addr", "127.0.0.1", "host: of the SOCKS proxy")
	socksPortString       = flag.String("socks-port", "4446", ":port of the SOCKS proxy")
	addrHelperHostString  = flag.String("ah-addr", "127.0.0.1", "host: of the SAM bridge")
	addrHelperPortString  = flag.String("ah-port", "7858", ":port of the SAM bridge")
	debugConnection       = flag.Bool("conn-debug", false, "Print connection debug info")
	useHTTPProxy          = flag.Bool("http-proxy", true, "run the HTTP proxy(default true)")
	useSOCKSProxy         = flag.Bool("socks-proxy", false, "run the SOCKS proxy(default false)")
	verboseLogging        = flag.Bool("verbose", false, "Print connection debug info")
	address               = flag.String("url", "", "i2p URL you want to retrieve")
	addressHelper         = flag.String("addresshelper", "http://joajgazyztfssty4w2on5oaqksz6tqoxbduy553y34mf4byv6gpq.b32.i2p/export/alive-hosts.txt", "Jump/Addresshelper service you want to use")
	destLifespan          = flag.Int("lifespan", 12, "Lifespan of an idle i2p destination in minutes(default twelve)")
	timeoutTime           = flag.Int("timeout", 6, "Timeout duration in minutes(default six)")
	tunnelLength          = flag.Int("tunlength", 3, "Tunnel Length(default 3)")
	inboundTunnels        = flag.Int("in-tunnels", 8, "Inbound Tunnel Count(default 8)")
	outboundTunnels       = flag.Int("out-tunnels", 8, "Inbound Tunnel Count(default 8)")
	keepAlives            = flag.Bool("disable-keepalives", false, "Disable keepalives(default false)")
	idleConns             = flag.Int("idle-conns", 4, "Maximium idle connections per host(default 4)")
	inboundBackups        = flag.Int("in-backups", 3, "Inbound Backup Count(default 3)")
	outboundBackups       = flag.Int("out-backups", 3, "Inbound Backup Count(default 3)")
	internalAddressHelper = flag.Bool("internal-ah", true, "Use internal address helper")
	addressBook           = flag.String("addressbook", "/etc/si-i2p-plugin/addresses.csv", "path to local addressbook(default ./etc/si-i2p-plugin/addresses.csv) (Unused without internal-ah)")
	internalSamRTCHost    = flag.Bool("internal-rtc", false, "Use internal SamRTC(Experimenatl)(default false)")
	rtcHostString         = flag.String("rtc-addr", "127.0.0.1", "host: of the RTC over SAM bridge")
	rtcPortString         = flag.String("rtc-port", "7682", ":port of the RTC over SAM bridge")
	diskAvoidance         = flag.Bool("avoidance", true, "Disk Avoidance Mode(default true)")
)

func main() {
	//lb := littleboss.New("si-i2p-plugin")
	//lb.Run(func(ctx context.Context){
	Defwd, _ := os.Getwd()
	workDirectory := flag.String("directory", Defwd, "The working directory you want to use, defaults to current directory")

	flag.Parse()

	log.SetOutput(os.Stdout)

	dii2perrs.Log("si-i2p-plugin.go Sam Address:", *samAddrString)
	dii2perrs.Log("si-i2p-plugin.go Sam Port:", *samPortString)
	dii2perrs.Log("si-i2p-plugin.go HTTP Proxy Address:", *proxAddrString)
	dii2perrs.Log("si-i2p-plugin.go HTTP Proxy Port:", *proxPortString)
	dii2perrs.Log("si-i2p-plugin.go SOCKS Proxy Address:", *socksAddrString)
	dii2perrs.Log("si-i2p-plugin.go SOCKS Proxy Port:", *socksPortString)
	dii2perrs.Log("si-i2p-plugin.go Addresshelper Address:", *addrHelperHostString)
	dii2perrs.Log("si-i2p-plugin.go Addresshelper Port:", *addrHelperPortString)
	dii2perrs.Log("si-i2p-plugin.go Working Directory:", *workDirectory)
	dii2perrs.Log("si-i2p-plugin.go Addresshelper Services:", *addressHelper)
	dii2perrs.Log("si-i2p-plugin.go Timeout Time:", *timeoutTime, "minutes")
	dii2perrs.Log("si-i2p-plugin.go Tunnel Length:", *tunnelLength)
	dii2perrs.Log("si-i2p-plugin.go Inbound Tunnel Quantity:", *inboundTunnels)
	dii2perrs.Log("si-i2p-plugin.go Outbound Tunnel Quantity", *outboundTunnels)
	dii2perrs.Log("si-i2p-plugin.go Idle Tunnel Count:", *idleConns)
	dii2perrs.Log("si-i2p-plugin.go Inbound Backup Quantity:", *inboundBackups)
	dii2perrs.Log("si-i2p-plugin.go Outbound Backup Quantity", *outboundBackups)
	if !*keepAlives {
		dii2perrs.Log("si-i2p-plugin.go Keepalives Enabled")
	} else {
		dii2perrs.Log("si-i2p-plugin.go Keepalives Disabled")
	}
	if *debugConnection {
		dii2perrs.DEBUG = *debugConnection
		dii2perrs.Log("si-i2p-plugin.go Debug mode: true")
	}
	if *verboseLogging {
		dii2perrs.Verbose = *verboseLogging
		dii2perrs.Log("si-i2p-plugin.go Verbose mode: true")
	}

	if *internalAddressHelper {
		dii2perrs.Log("si-i2p-plugin.go starting internal addresshelper with")
		jumphelper.NewService(
			*addrHelperHostString,
			*addrHelperPortString,
			*addressBook,
			*samAddrString,
			*samPortString,
			[]string{*addressHelper},
			*internalAddressHelper,
		)
	}
	if *internalSamRTCHost {
		dii2perrs.Log("si-i2p-plugin.go starting internal RTC forwarder")
		if err := samrtc.NewEmbedSamRTCHostFromOptions(
			samrtc.SetHostSamHost(*samAddrString),
			samrtc.SetHostSamPort(*samPortString),
			samrtc.SetHostLocalHost(*rtcHostString),
			samrtc.SetHostLocalPort(*rtcPortString),
			samrtc.SetHostSamTunName("siSAMRTC"),
			samrtc.SetHostSamVerbose(*verboseLogging),
		); err != nil {
			dii2perrs.Fatal(err, "si-i2p-plugin.go failed to start internal RTC forwarder")
		}
	}

	samProxies, err := dii2pmain.CreateSamList(
		dii2pmain.SetInitAddress(*address),
		dii2pmain.SetHost(*samAddrString),
		dii2pmain.SetPort(*samPortString),
		dii2pmain.SetTimeout(*timeoutTime),
		dii2pmain.SetKeepAlives(*keepAlives),
		dii2pmain.SetLifespan(*destLifespan),
		dii2pmain.SetTunLength(*tunnelLength),
		dii2pmain.SetInQuantity(*inboundTunnels),
		dii2pmain.SetOutQuantity(*outboundTunnels),
		dii2pmain.SetIdleConns(*idleConns),
		dii2pmain.SetInBackups(*inboundBackups),
		dii2pmain.SetOutBackups(*outboundBackups),
	)

	if err != nil {
		log.Fatal(err)
	}

	samService, err := dii2pserv.CreateSamServiceList(
		dii2pserv.SetServHost(*samAddrString),
		dii2pserv.SetServPort(*samPortString),
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
			samProxy := dii2p.CreateHTTPProxy(
				*proxAddrString,
				*proxPortString,
				*address,
				*addrHelperHostString,
				*addrHelperPortString,
				*addressHelper,
				samProxies,
				*timeoutTime,
				*keepAlives,
			)
			dii2perrs.Log("si-i2p-plugin.go HTTP Proxy Started:" + samProxy.Addr)
			httpUp = true
		}
	}

	if *useSOCKSProxy {
		if !socksUp {
			samProxy := dii2p.CreateSOCKSProxy(
				*proxAddrString,
				*proxPortString,
				*address,
				*addrHelperHostString,
				*addrHelperPortString,
				*addressHelper,
				samProxies,
				*timeoutTime,
				*keepAlives,
			)
			dii2perrs.Log("si-i2p-plugin.go Socks Proxy Started:" + samProxy.Addr)
			socksUp = true
		}
	}

	dii2perrs.Log("si-i2p-plugin.go Created client, starting loop...")

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
	//})
}

func closeProxy(samProxies *dii2pmain.SamList) {
	exit = samProxies.ReadDelete()
}

func closeServices(samServiceList *dii2pserv.SamServices) {
	exit = samServiceList.ReadDelete()
}
