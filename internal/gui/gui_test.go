//go:build systray
// +build systray

package gui

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skywire-utilities/pkg/cipher"
	"github.com/skycoin/skywire/pkg/visor/visorconfig"
)

func TestGetAvailPublicVPNServers(t *testing.T) {
	pk, sk := cipher.GenerateKeyPair()
	common := &visorconfig.Common{
		Version: "v1.1.0",
		SK:      sk,
		PK:      pk,
	}
	config := visorconfig.MakeBaseConfig(common)
	servers := GetAvailPublicVPNServers(config)
	require.NotEqual(t, nil, servers)
	require.NotEqual(t, []string{}, servers)
	t.Logf("Servers: %v", servers)
}

func TestReadEmbeddedIcon(t *testing.T) {
	b, err := ReadSysTrayIcon()
	require.NoError(t, err)
	require.NotEqual(t, 0, len(b))
}
