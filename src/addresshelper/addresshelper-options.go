package dii2pah

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
func SetAddressHelperPort(s string) func(*AddressHelper) error {
	return func(c *AddressHelper) error {
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

//SetAddressHelperPortInt sets the port of the addresshelper client's SAM bridge
func SetAddressHelperPortInt(s int) func(*AddressHelper) error {
	return func(c *AddressHelper) error {
		if s < 65536 && s > -1 {
			c.jumpPortString = strconv.Itoa(s)
			return nil
		}
		return fmt.Errorf("Invalid port")
	}
}

//SetAddressBookPath sets the address book path
func SetAddressBookPath(s string) func(*AddressHelper) error {
	return func(c *AddressHelper) error {
		c.bookPath = s
		return nil
	}
}
