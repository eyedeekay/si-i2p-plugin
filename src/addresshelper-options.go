package dii2p

import (
	"fmt"
	"strconv"
)

//AddressHelperConnectOption is a AddressHelper option
type AddressHelperConnectOption func(*AddressHelper) error

//SetAddressHelperURL sets the host of the addresshelper client's SAM bridge
func SetAddressHelperURL(s string) func(*AddressHelper) error {
	return func(c *AddressHelper) error {
		c.addressHelperURL = s
		return nil
	}
}

//SetAddressHelperHost sets the host of the addresshelper client's SAM bridge
func SetAddressHelperHost(s string) func(*AddressHelper) error {
	return func(c *AddressHelper) error {
		c.jumpHostString = s
		return nil
	}
}

//SetAddressHelperPort sets the port of the addresshelper client's SAM bridge
func SetAddressHelperPort(s interface{}) func(*AddressHelper) error {
	return func(c *AddressHelper) error {
		switch v := s.(type) {
		case string:
			port, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("Invalid port; non-number")
			}
			if port < 65536 && port > -1 {
				c.jumpPortString = v
				return nil
			}
			return fmt.Errorf("Invalid port")
		case int:
			if v < 65536 && v > -1 {
				c.jumpPortString = strconv.Itoa(v)
				return nil
			}
			return fmt.Errorf("Invalid port")
		default:
			return fmt.Errorf("Invalid port")
		}
	}
}

//SetAddressBookPath sets the address book path
func SetAddressBookPath(s string) func(*AddressHelper) error {
	return func(c *AddressHelper) error {
		c.bookPath = s
		return nil
	}
}
