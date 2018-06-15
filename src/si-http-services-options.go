package dii2p

import (
	"fmt"
	"strconv"
)

// ServiceOption is a functional argument for creating an http service
type ServiceOption func(*SamServices) error

// SetServHost sets the ip address of the SAM bridge
func SetServHost(s string) func(*SamServices) error {
	return func(c *SamServices) error {
		c.samAddrString = s
		return nil
	}
}

// SetServPort sets the port of the SAM bridge
func SetServPort(s interface{}) func(*SamServices) error {
	return func(c *SamServices) error {
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
