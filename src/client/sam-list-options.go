package dii2pmain

import (
	"fmt"
	"strconv"
)

import (
	"github.com/eyedeekay/si-i2p-plugin/src/helpers"
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
func SetPort(v string) func(*SamList) error {
	return func(c *SamList) error {
		port, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("Invalid port; non-number.")
		}
		if port < 65536 && port > -1 {
			c.samPortString = v
			return nil
		}
		return fmt.Errorf("Invalid port.")
	}
}

//SetPortInt sets the port of the client's SAM bridge
func SetPortInt(v int) func(*SamList) error {
	return func(c *SamList) error {
		if v < 65536 && v > -1 {
			c.samPortString = strconv.Itoa(v)
			return nil
		}
		return fmt.Errorf("Invalid port.")
	}
}

//SetTimeout set's the client timeout
func SetTimeout(s int) func(*SamList) error {
	return func(c *SamList) error {
		if s > 5 {
			if s <= c.lifeTime {
				c.timeoutTime = s
				return nil
			}
			return fmt.Errorf("A specified lifetime must be greater than a specified timeout.")
		}
		return fmt.Errorf("Timeout must be greater than 5 minutes.")
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
		if !dii2phelper.CheckURLType(s) {
			c.lastAddress = ""
			return fmt.Errorf("Init Address was not an i2p url.")
		}
		c.lastAddress = s
		return nil
	}
}

//SetLifespan set's the time before an inactive client is torn down
func SetLifespan(s int) func(*SamList) error {
	return func(c *SamList) error {
		if c.timeoutTime <= s {
			c.lifeTime = s
			return nil
		}
		return fmt.Errorf("A specified lifetime must be greater than a specified timeout.")
	}
}

//SetTunLength set's the symmetric inbound and outbound tunnel lengths
func SetTunLength(s int) func(*SamList) error {
	return func(c *SamList) error {
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

//SetInQuantity set's the inbound tunnel quantity
func SetInQuantity(s int) func(*SamList) error {
	return func(c *SamList) error {
		if s > 0 {
			if s < 16 {
				c.inboundQuantity = s
				return nil
			}
			return fmt.Errorf("Tunnel quantity must be less than 16.")
		}
		return fmt.Errorf("Tunnel quantity must be greater than 0.")
	}
}

//SetOutQuantity set's the outbound tunnel quantity
func SetOutQuantity(s int) func(*SamList) error {
	return func(c *SamList) error {
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

//SetInBackups set's the inbound backup tunnel quantity
func SetInBackups(s int) func(*SamList) error {
	return func(c *SamList) error {
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

//SetOutBackups set's the outbound backup tunnel quantity
func SetOutBackups(s int) func(*SamList) error {
	return func(c *SamList) error {
		if s >= 0 {
			if s < 6 {
				c.outboundBackupQuantity = s
				return nil
			}
			return fmt.Errorf("Outbound backup tunnel quantity cannot be negative.")
		}
		return fmt.Errorf("Outbound backup tunnel quantity must be less than 6")
	}
}

//SetIdleConns set's the max idle connections per host
func SetIdleConns(s int) func(*SamList) error {
	return func(c *SamList) error {
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
