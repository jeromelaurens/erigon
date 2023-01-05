package observer

import (
	"context"
	"testing"

	"github.com/jeromelaurens/erigon/crypto"
	"github.com/jeromelaurens/erigon/eth/protocols/eth"
	"github.com/jeromelaurens/erigon/p2p/enode"
	"github.com/jeromelaurens/erigon/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandshake(t *testing.T) {
	t.Skip("only for dev")

	// grep 'self=enode' the log, and paste it here
	// url := "enode://..."
	url := params.MainnetBootnodes[0]
	node := enode.MustParseV4(url)
	myPrivateKey, _ := crypto.GenerateKey()

	ctx := context.Background()
	hello, status, err := Handshake(ctx, node.IP(), node.TCP(), node.Pubkey(), myPrivateKey)

	require.Nil(t, err)
	require.NotNil(t, hello)
	assert.Equal(t, uint64(5), hello.Version)
	assert.NotEmpty(t, hello.ClientID)
	assert.Contains(t, hello.ClientID, "erigon")

	require.NotNil(t, status)
	assert.Equal(t, uint32(eth.ETH66), status.ProtocolVersion)
	assert.Equal(t, uint64(1), status.NetworkID)
}
