package dii2p

import (
	"fmt"
	"strconv"
)

//ClientOption is a SamList option
type ClientOption func(*SamList) error

//SetHost sets the host of the client's SAM bridge
func SetHost(s string) func(*SamList) error {
	return func(c *SamList) error {
		c.samAddrString = s
		return nil
	}
}

//SetPort sets the port of the client's SAM bridge
func SetPort(s interface{}) func(*SamList) error {
	return func(c *SamList) error {
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

//SetTimeout set's the client timeout
func SetTimeout(s int) func(*SamList) error {
	return func(c *SamList) error {
		c.timeoutTime = s
		return nil
	}
}

//SetKeepAlives tells the client whether or not to allow keepalives
func SetKeepAlives(s bool) func(*SamList) error {
	return func(c *SamList) error {
		c.keepAlives = s
		return nil
	}
}

//SetInitAddress tells the client to retrieve an URL before any other URL
func SetInitAddress(s string) func(*SamList) error {
	return func(c *SamList) error {
		if s == "" {
			c.lastAddress = s
			return nil
		}
		if !CheckURLType(s) {
			c.lastAddress = ""
			return fmt.Errorf("Init Address was not an i2p url")
		}
		c.lastAddress = s
		return nil
	}
}
