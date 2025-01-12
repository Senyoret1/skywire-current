package handshake

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/skycoin/dmsg/pkg/dmsg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skywire-utilities/pkg/cipher"
)

const (
	port1 = 10
	port2 = 11
)

func TestHandshake(t *testing.T) {
	type hsResult struct {
		lAddr dmsg.Addr
		rAddr dmsg.Addr
		err   error
	}

	for i := byte(0); i < 64; i++ {
		initPK, initSK, err := cipher.GenerateDeterministicKeyPair(append([]byte("init"), i))
		require.NoError(t, err)

		respPK, _, err := cipher.GenerateDeterministicKeyPair(append([]byte("resp"), i))
		require.NoError(t, err)

		iAddr := dmsg.Addr{PK: initPK, Port: port1}
		rAddr := dmsg.Addr{PK: respPK, Port: port2}

		initC, respC := net.Pipe()

		deadline := time.Now().Add(Timeout)

		respCh := make(chan hsResult, 1)

		go func() {
			defer close(respCh)

			respHS := ResponderHandshake(func(f2 Frame2) error {
				if f2.SrcAddr.PK != initPK {
					return errors.New("unexpected src addr pk")
				}
				if f2.DstAddr.PK != respPK {
					return errors.New("unexpected dst addr pk")
				}
				return nil
			})

			lAddr, rAddr, err := respHS(respC, deadline)
			respCh <- hsResult{lAddr: lAddr, rAddr: rAddr, err: err}
		}()

		initHS := InitiatorHandshake(initSK, iAddr, rAddr)

		var initR hsResult
		initR.lAddr, initR.rAddr, initR.err = initHS(initC, deadline)

		assert.NoError(t, err)
		assert.Equal(t, initR.lAddr, iAddr)
		assert.Equal(t, initR.rAddr, rAddr)

		rr := <-respCh

		assert.NoError(t, rr.err)
		assert.Equal(t, rr.lAddr, rAddr)
		assert.Equal(t, rr.rAddr, iAddr)

		assert.NoError(t, initC.Close())
		assert.NoError(t, respC.Close())
	}
}
