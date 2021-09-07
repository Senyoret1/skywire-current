package skyenv

import (
	"github.com/skycoin/dmsg/dmsgpty"
	"time"

	"github.com/skycoin/dmsg/cipher"
)

// Constants for skywire root directories.
const (
	DefaultSkywirePath = "."
	PackageSkywirePath = "/opt/skywire"
)

// Constants for old default services.
const (
	OldDefaultTpDiscAddr          = "http://transport.discovery.skywire.skycoin.com"
	OldDefaultDmsgDiscAddr        = "http://dmsg.discovery.skywire.skycoin.com"
	OldDefaultServiceDiscAddr     = "http://service.discovery.skycoin.com"
	OldDefaultRouteFinderAddr     = "http://routefinder.skywire.skycoin.com"
	OldDefaultUptimeTrackerAddr   = "http://uptime-tracker.skywire.skycoin.com"
	OldDefaultAddressResolverAddr = "http://address.resolver.skywire.skycoin.com"
)

// Constants for new default services.
const (
	DefaultTpDiscAddr          = "http://tpd.skywire.skycoin.com"
	DefaultDmsgDiscAddr        = "http://dmsgd.skywire.skycoin.com"
	DefaultServiceDiscAddr     = "http://sd.skycoin.com"
	DefaultRouteFinderAddr     = "http://rf.skywire.skycoin.com"
	DefaultUptimeTrackerAddr   = "http://ut.skywire.skycoin.com"
	DefaultAddressResolverAddr = "http://ar.skywire.skycoin.com"
	DefaultSetupPK             = "0324579f003e6b4048bae2def4365e634d8e0e3054a20fc7af49daf2a179658557"
)

// Constants for testing deployment.
const (
	TestTpDiscAddr          = "http://transport.discovery.skywire.cc"
	TestDmsgDiscAddr        = "http://dmsg.discovery.skywire.cc"
	TestServiceDiscAddr     = "http://service.discovery.skywire.cc"
	TestRouteFinderAddr     = "http://routefinder.skywire.cc"
	TestUptimeTrackerAddr   = "http://uptime.tracker.skywire.cc"
	TestAddressResolverAddr = "http://address.resolver.skywire.cc"
	TestSetupPK             = "026c5a07de617c5c488195b76e8671bf9e7ee654d0633933e202af9e111ffa358d"
)

// Dmsg port constants.
// TODO(evanlinjin): Define these properly. These are currently random.
const (
	DmsgCtrlPort           uint16 = 7   // Listening port for dmsgctrl protocol (similar to TCP Echo Protocol).
	DmsgSetupPort          uint16 = 36  // Listening port of a setup node.
	DmsgHypervisorPort     uint16 = 46  // Listening port of a hypervisor for incoming RPC visor connections over dmsg.
	DmsgTransportSetupPort uint16 = 47  // Listening port for transport setup RPC over dmsg.
	DmsgAwaitSetupPort     uint16 = 136 // Listening port of a visor for setup operations.
)

// Transport port constants.
const (
	TransportPort uint16 = 45 // Listening port of a visor for incoming transports.
)

// Default dmsgpty constants.
const (
	DmsgPtyPort          uint16 = 22
	DefaultDmsgPtyCLINet        = "unix"
)

// DefaultDmsgPtyCLIAddr determines default CLI address per each platform
func DefaultDmsgPtyCLIAddr() string {
	return dmsgpty.DefaultCLIAddr()
}

// Default STCP constants.
const (
	DefaultSTCPAddr = ":7777"
)

// Default skywire app constants.
const (
	SkychatName        = "skychat"
	SkychatPort uint16 = 1
	SkychatAddr        = ":8001"

	SkysocksName        = "skysocks"
	SkysocksPort uint16 = 3

	SkysocksClientName        = "skysocks-client"
	SkysocksClientPort uint16 = 13
	SkysocksClientAddr        = ":1080"

	VPNServerName        = "vpn-server"
	VPNServerPort uint16 = 44

	VPNClientName = "vpn-client"
	// TODO(darkrengarius): this one's not needed for the app to run but lack of it causes errors
	VPNClientPort uint16 = 43
)

// RPC constants.
const (
	DefaultRPCAddr      = "localhost:3435"
	DefaultRPCTimeout   = 20 * time.Second
	TransportRPCTimeout = 1 * time.Minute
	UpdateRPCTimeout    = 6 * time.Hour // update requires huge timeout
)

// Default skywire app server and discovery constants
const (
	DefaultAppSrvAddr     = "localhost:5505"
	AppDiscUpdateInterval = time.Minute
	DefaultAppBinPath     = DefaultSkywirePath + "/apps"
	DefaultLogLevel       = "info"
)

// Package defaults
const (
	PackageAppBinPath       = PackageSkywirePath + "/apps"
	PackageLocalPath        = PackageSkywirePath + "/local"
	PackageDmsgPtyWhiteList = PackageSkywirePath + "/dmsgpty/whitelist.json"
	PackageDmsgPtyCLIAddr   = PackageSkywirePath + "/dmsgpty/cli.sock"
	PackageTpLogStore       = PackageSkywirePath + "/transport_logs"
	PackageDBPath           = PackageSkywirePath + "/users.db"
	PackageEnableTLS        = false
	PackageTLSKey           = PackageSkywirePath + "/ssl/key.pem"
	PackageTLSCert          = PackageSkywirePath + "/ssl/cert.pem"
)

// Default routing constants
const (
	DefaultTpLogStore = DefaultSkywirePath + "/transport_logs"
)

// Skybian defaults
const (
	SkybianAppBinPath       = "/usr/bin/apps"
	SkybianDmsgPtyWhiteList = "/var/skywire-visor/dsmgpty/whitelist.json"
	SkybianDmsgPtyCLIAddr   = "/run/skywire-visor/dmsgpty/cli.sock"
	SkybianLocalPath        = "/var/skywire-visor/apps"
	SkybianTpLogStore       = "/var/skywire-visor/transports"
	SkybianEnableTLS        = false
	SkybianDBPath           = "/var/skywire-visor/users.db"
	SkybianTLSKey           = "/var/skywire-visor/ssl/key.pem"
	SkybianTLSCert          = "/var/skywire-visor/ssl/cert.pem"
)

// Default local constants
const (
	DefaultLocalPath = DefaultSkywirePath + "/local"
)

// Default hypervisor constants
const (
	DefaultHypervisorDB = ".skycoin/hypervisor/users.db"
	DefaultEnableAuth   = true
	DefaultEnableTLS    = false
	DefaultTLSKey       = DefaultSkywirePath + "/ssl/key.pem"
	DefaultTLSCert      = DefaultSkywirePath + "/ssl/cert.pem"
)

// MustPK unmarshals string PK to cipher.PubKey. It panics if unmarshaling fails.
func MustPK(pk string) cipher.PubKey {
	var sPK cipher.PubKey
	if err := sPK.UnmarshalText([]byte(pk)); err != nil {
		panic(err)
	}

	return sPK
}

// GetStunServers gives back deafault Stun Servers
func GetStunServers() []string {
	return []string{
		"45.118.133.242:3478",
		"192.53.173.68:3478",
		"192.46.228.39:3478",
		"192.53.113.106:3478",
		"192.53.117.158:3478",
		"192.53.114.142:3478",
		"139.177.189.166:3478",
		"192.46.227.227:3478",
	}
}
