package jumpresolver

import (
	"fmt"
	"strconv"
)

//JumpResolverOption is a JumpResolver option
type JumpResolverOption func(*JumpResolver) error

//SetJumpResolverHost sets the host of the JumpResolver client's SAM bridge
func SetJumpResolverHost(s string) func(*JumpResolver) error {
	return func(c *JumpResolver) error {
		c.jumpHostString = s
		return nil
	}
}

//SetJumpResolverPort sets the port of the JumpResolver client's SAM bridge
func SetJumpResolverPort(s string) func(*JumpResolver) error {
	return func(c *JumpResolver) error {
		port, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		if port < 65536 && port > -1 {
			c.jumpPortString = s
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetJumpResolverPortInt sets the port of the JumpResolver client's SAM bridge
func SetJumpResolverPortInt(s int) func(*JumpResolver) error {
	return func(c *JumpResolver) error {
		if s < 65536 && s > -1 {
			c.jumpPortString = strconv.Itoa(s)
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}
