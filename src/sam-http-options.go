package dii2p

import (
	"fmt"
	"strconv"
    "time"
)

//ClientOption is a SamHTTP option
type ConnectOption func(*SamHTTP) error

func SetSamHTTPHost(s string) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		c.samAddrString = s
		return nil
	}
}

//SetPort sets the port of the client's SAM bridge
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

func SetSamHTTPRequest(s string) func(*SamHTTP) error {
	return func(c *SamHTTP) error {
		c.initRequestURL = s
		return nil
	}
}

func SetSamHTTPTimeout(s int) func(*SamHTTP) error {
    return func(c *SamHTTP) error {
		c.timeoutTime = time.Duration(s) * time.Minute
        c.otherTimeoutTime = time.Duration(s / 3) * time.Minute
		return nil
	}
}

func SetSamHTTPKeepAlives(s bool) func(*SamHTTP) error {
    return func(c *SamHTTP) error {
		c.keepAlives = s
		return nil
	}
}
