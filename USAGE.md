Usage of ./bin/si-i2p-plugin:
=============================

All Options:
------------

          -addressbook string
                path to local addressbook(default ./addresses.csv) (Unused without internal-ah) (default "./addresses.csv")
          -addresshelper string
                Jump/Addresshelper service you want to use (default "http://joajgazyztfssty4w2on5oaqksz6tqoxbduy553y34mf4byv6gpq.b32.i2p/export/alive-hosts.txt")
          -ah-addr string
                host: of the SAM bridge (default "127.0.0.1")
          -ah-port string
                :port of the SAM bridge (default "7054")
          -bridge-addr string
                host: of the SAM bridge (default "127.0.0.1")
          -bridge-port string
                :port of the SAM bridge (default "7656")
          -conn-debug
                Print connection debug info
          -directory string
                The working directory you want to use, defaults to current directory (default "/home/idk/local-manifest/crypto-manifest/si-i2p-plugin")
          -disable-keepalives
                Disable keepalives(default false)
          -http-proxy
                run the HTTP proxy(default true) (default true)
          -idle-conns int
                Maximium idle connections per host(default 8) (default 8)
          -in-backups int
                Inbound Backup Count(default 3) (default 3)
          -in-tunnels int
                Inbound Tunnel Count(default 15) (default 15)
          -internal-ah
                Use internal address helper (default true)
          -lifespan int
                Lifespan of an idle i2p destination in minutes(default twelve) (default 12)
          -out-backups int
                Inbound Backup Count(default 3) (default 3)
          -out-tunnels int
                Inbound Tunnel Count(default 15) (default 15)
          -proxy-addr string
                host: of the HTTP proxy (default "127.0.0.1")
          -proxy-port string
                :port of the HTTP proxy (default "4443")
          -socks-addr string
                host: of the SOCKS proxy (default "127.0.0.1")
          -socks-port string
                :port of the SOCKS proxy (default "4446")
          -socks-proxy
                run the SOCKS proxy(default false)
          -timeout int
                Timeout duration in minutes(default six) (default 6)
          -tunlength int
                Tunnel Length(default 3) (default 3)
          -url string
                i2p URL you want to retrieve
          -verbose
                Print connection debug info
