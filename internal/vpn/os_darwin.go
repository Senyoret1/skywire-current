//go:build darwin
// +build darwin

package vpn

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/skycoin/skywire/pkg/util/osutil"
)

// SetupTUN sets the allocated TUN interface up, setting its IP, gateway, netmask and MTU.
func (c *Client) SetupTUN(ifcName, ipCIDR, gateway string, mtu int) error {
	ip, netmask, err := parseCIDR(ipCIDR)
	if err != nil {
		return fmt.Errorf("error parsing IP CIDR: %w", err)
	}
	if err := c.setSysPrivileges(); err != nil {
		print(fmt.Sprintf("Failed to setup system privileges for SetupTUN: %v\n", err))
		return err
	}
	defer c.releaseSysPrivileges()
	return osutil.Run("ifconfig", ifcName, ip, gateway, "mtu", strconv.Itoa(mtu), "netmask", netmask, "up")

}

// ChangeRoute changes current route to `ipCIDR` to go through the `gateway`
// in the OS routing table.
func (c *Client) ChangeRoute(ipCIDR, gateway string) error {
	return c.modifyRoutingTable("change", ipCIDR, gateway)
}

// AddRoute adds route to `ipCIDR` through the `gateway` to the OS routing table.
func (c *Client) AddRoute(ipCIDR, gateway string) error {
	if err := c.modifyRoutingTable("add", ipCIDR, gateway); err != nil {
		var e *osutil.ErrorWithStderr
		if errors.As(err, &e) {
			if strings.Contains(string(e.Stderr), "File exists") {
				return nil
			}
		}
		return err
	}
	return nil
}

// DeleteRoute removes route to `ipCIDR` through the `gateway` from the OS routing table.
func (c *Client) DeleteRoute(ipCIDR, gateway string) error {
	return c.modifyRoutingTable("delete", ipCIDR, gateway)
}

func (c *Client) modifyRoutingTable(action, ipCIDR, gateway string) error {
	ip, netmask, err := parseCIDR(ipCIDR)
	if err != nil {
		return fmt.Errorf("error parsing IP CIDR: %w", err)
	}

	if err := c.setSysPrivileges(); err != nil {
		print(fmt.Sprintf("Failed to setup system privileges for %s: %v\n", action, err))
		return err
	}
	defer c.releaseSysPrivileges()
	return osutil.Run("route", action, "-net", ip, gateway, netmask)
}

// SetupTUN sets the allocated TUN interface up, setting its IP, gateway, netmask and MTU.
func (s *Server) SetupTUN(ifcName, ipCIDR, gateway string, mtu int) error {
	ip, netmask, err := parseCIDR(ipCIDR)
	if err != nil {
		return fmt.Errorf("error parsing IP CIDR: %w", err)
	}

	return osutil.Run("ifconfig", ifcName, ip, gateway, "mtu", strconv.Itoa(mtu), "netmask", netmask, "up")
}
