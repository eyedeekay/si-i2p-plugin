package dii2pmain

import (
	"fmt"
	"strconv"
	"time"
)

import (
	"github.com/eyedeekay/si-i2p-plugin/src/helpers"
)

//ConnectOption is a SamHTTP option
type ConnectOption func(*SamHTTP) error

//SetSamHTTPHost sets the host of the client's SAM bridge
func SetSamHTTPHost(s string) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		c.samAddrString = s
		return nil
	}
}

//SetSamHTTPPort sets the port of the client's SAM bridge
func SetSamHTTPPort(v string) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		port, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("Invalid port; non-number")
		}
		if port < 65536 && port > -1 {
			c.samPortString = v
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetSamHTTPPortInt sets the port of the client's SAM bridge
func SetSamHTTPPortInt(v int) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		if v < 65536 && v > -1 {
			c.samPortString = strconv.Itoa(v)
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetSamHTTPRequest sets the initial request URL for the SamHTTP connection
func SetSamHTTPRequest(s string) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		if dii2phelper.CheckURLType(s) {
			c.initRequestURL = s
			return nil
		}
		return fmt.Errorf("Invalid initiate URL %s", s)
	}
}

//SetSamHTTPTimeout sets the timeout of the SamHTTP connection
func SetSamHTTPTimeout(s int) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		if s > 5 {
			if time.Duration(s)*time.Minute <= c.lifeTime {
				c.timeoutTime = time.Duration(s) * time.Minute
				return nil
			}
            tmp := time.Duration(s) * time.Minute
			return fmt.Errorf("A specified timeout must be less than a specified lifetime. %s %s %s", tmp.String(), ">" , c.lifeTime.String())
		}
		return fmt.Errorf("Timeout must be greater than 5 minutes.")
	}
}

//SetSamHTTPKeepAlives tells the SamHTTP connection whether to accept keepAlives
func SetSamHTTPKeepAlives(s bool) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		c.keepAlives = s
		return nil
	}
}

//SetSamHTTPLifespan set's the time before an inactive SamHTTP client is torn down
func SetSamHTTPLifespan(s int) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		if time.Duration(s)*time.Minute >= c.timeoutTime {
			c.lifeTime = time.Duration(s) * time.Minute
			c.useTime = time.Now()
			return nil
		}
        tmp := time.Duration(s) * time.Minute
        return fmt.Errorf("A specified lifetime must be greater than a specified timeout. %s %s %s", tmp.String(), "<", c.timeoutTime.String())
	}
}

//SetSamHTTPTunLength set's the symmetric inbound and outbound tunnel lengths
func SetSamHTTPTunLength(s int) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		if s >= 0 {
			if s <= 7 {
				c.tunnelLength = s
				return nil
			}
			return fmt.Errorf("Tunnel length must be less than seven.")
		}
		return fmt.Errorf("Tunnel length must be greater than or equal to 0.")
	}
}

//SetSamHTTPInQuantity set's the inbound tunnel quantity
func SetSamHTTPInQuantity(s int) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		if s > 0 {
			if s < 16 {
				c.inboundQuantity = s
				c.inboundQuantity = s
				return nil
			}
			return fmt.Errorf("Tunnel quantity must be less than 16.")
		}
		return fmt.Errorf("Tunnel quantity must be greater than 0.")
	}
}

//SetSamHTTPOutQuantity set's the outbound tunnel quantity
func SetSamHTTPOutQuantity(s int) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		if s > 0 {
			if s < 16 {
				c.outboundQuantity = s
				return nil
			}
			return fmt.Errorf("Tunnel quantity must be less than 16.")
		}
		return fmt.Errorf("Tunnel quantity must be greater than 0.")
	}
}

//SetSamHTTPInBackupQuantity set's the inbound tunnel quantity
func SetSamHTTPInBackupQuantity(s int) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		if s >= 0 {
			if s < 6 {
				c.inboundBackupQuantity = s
				return nil
			}
			return fmt.Errorf("Inbound backup tunnel quantity cannot be negative.")
		}
		return fmt.Errorf("Inbound backup tunnel quantity must be less than 6")
	}
}

//SetSamHTTPOutBackupQuantity set's the outbound tunnel quantity
func SetSamHTTPOutBackupQuantity(s int) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		if s >= 0 {
			if s < 6 {
				c.outboundBackupQuantity = s
				c.outboundBackupQuantity = s
				return nil
			}
			return fmt.Errorf("Outbound backup tunnel quantity cannot be negative.")
		}
		return fmt.Errorf("Outbound backup tunnel quantity must be less than 6")
	}
}

//SetSamHTTPIdleQuantity set's the outbound tunnel quantity
func SetSamHTTPIdleQuantity(s int) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		if s > 0 {
			if s < 11 {
				c.idleConns = s
				return nil
			}
			return fmt.Errorf("Idle connection quantity must less than than 11.")
		}
		return fmt.Errorf("Idle connection quantity must be greater than 0.")
	}
}
