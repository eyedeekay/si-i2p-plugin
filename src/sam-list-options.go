package dii2p

import (
    "fmt"
    "strconv"
)

type Option func(*SamList) error

func SetHost(s string) func(*SamList) error {
    return func(c *SamList) error {
        c.samAddrString = s
        return nil
    }
}

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
			} else {
				return fmt.Errorf("Invalid port")
			}
		case int:
			if v < 65536 && v > -1 {
				c.samPortString = strconv.Itoa(v)
				return nil
			} else {
				return fmt.Errorf("Invalid port")
			}
		default:
			return fmt.Errorf("Invalid port")
		}
    }
}

func SetTimeout(s int) func(*SamList) error {
    return func(c *SamList) error {
        c.timeoutTime = s
        return nil
    }
}

func SetKeepAlives(s bool) func(*SamList) error {
    return func(c *SamList) error {
        c.keepAlives = s
        return nil
    }
}
