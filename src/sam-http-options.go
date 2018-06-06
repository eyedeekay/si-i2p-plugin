package dii2p

import (
	"fmt"
	"strconv"
	"time"
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
func SetSamHTTPPort(s interface{}) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		switch v := s.(type) {
		case string:
			port, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("Invalid port; non-number")
			}
			if port < 65536 && port > -1 {
				c.samPortString = v
				return nil
			}
			return fmt.Errorf("Invalid port")
		case int:
			if v < 65536 && v > -1 {
				c.samPortString = strconv.Itoa(v)
				return nil
			}
			return fmt.Errorf("Invalid port")
		default:
			return fmt.Errorf("Invalid port")
		}
	}
}

//SetSamHTTPRequest sets the initial request URL for the SamHTTP connection
func SetSamHTTPRequest(s string) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		c.initRequestURL = s
		return nil
	}
}

//SetSamHTTPTimeout sets the timeout of the SamHTTP connection
func SetSamHTTPTimeout(s int) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		c.timeoutTime = time.Duration(s) * time.Minute
		c.otherTimeoutTime = time.Duration(s/3) * time.Minute
		return nil
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
		c.lifeTime = time.Duration(s) * time.Minute
		c.useTime = time.Now()
		return nil
	}
}

//SetSamHTTPTunLength set's the symmetric inbound and outbound tunnel lengths
func SetSamHTTPTunLength(s int) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		c.tunnelLength = s
		return nil
	}
}

//SetSamHTTPInQuantity set's the inbound tunnel quantity
func SetSamHTTPInQuantity(s int) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		c.inboundQuantity = s
		return nil
	}
}

//SetSamHTTPOutQuantity set's the outbound tunnel quantity
func SetSamHTTPOutQuantity(s int) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		c.outboundQuantity = s
		return nil
	}
}
