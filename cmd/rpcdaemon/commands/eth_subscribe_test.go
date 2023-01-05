package commands

import (
	"context"
	"fmt"
	"testing"

	"github.com/jeromelaurens/erigon/cmd/rpcdaemon/rpcservices"
	"github.com/jeromelaurens/erigon/common"
	"github.com/jeromelaurens/erigon/core"
	"github.com/jeromelaurens/erigon/eth/protocols/eth"
	"github.com/jeromelaurens/erigon/ethdb/privateapi"
	"github.com/jeromelaurens/erigon/rlp"
	"github.com/jeromelaurens/erigon/turbo/rpchelper"
	"github.com/jeromelaurens/erigon/turbo/snapshotsync"
	"github.com/jeromelaurens/erigon/turbo/stages"
	"github.com/ledgerwatch/erigon-lib/direct"
	"github.com/ledgerwatch/erigon-lib/gointerfaces/sentry"
	"github.com/stretchr/testify/require"
)

func TestEthSubscribe(t *testing.T) {
	m, require := stages.Mock(t), require.New(t)
	br := snapshotsync.NewBlockReaderWithSnapshots(m.BlockSnapshots)
	chain, err := core.GenerateChain(m.ChainConfig, m.Genesis, m.Engine, m.DB, 7, func(i int, b *core.BlockGen) {
		b.SetCoinbase(common.Address{1})
	}, false /* intermediateHashes */)
	require.NoError(err)

	b, err := rlp.EncodeToBytes(&eth.BlockHeadersPacket66{
		RequestId:          1,
		BlockHeadersPacket: chain.Headers,
	})
	require.NoError(err)

	m.ReceiveWg.Add(1)
	for _, err = range m.Send(&sentry.InboundMessage{Id: sentry.MessageId_BLOCK_HEADERS_66, Data: b, PeerId: m.PeerId}) {
		require.NoError(err)
	}
	m.ReceiveWg.Wait() // Wait for all messages to be processed before we proceeed

	ctx := context.Background()
	backendServer := privateapi.NewEthBackendServer(ctx, nil, m.DB, m.Notifications.Events, br, nil, nil, nil, false)
	backendClient := direct.NewEthBackendClientDirect(backendServer)
	backend := rpcservices.NewRemoteBackend(backendClient, m.DB, br)
	ff := rpchelper.New(ctx, backend, nil, nil, func() {})

	newHeads, id := ff.SubscribeNewHeads(16)
	defer ff.UnsubscribeHeads(id)

	initialCycle := true
	highestSeenHeader := chain.TopBlock.NumberU64()
	if _, err := stages.StageLoopStep(m.Ctx, m.ChainConfig, m.DB, m.Sync, m.Notifications, initialCycle, m.UpdateHead); err != nil {
		t.Fatal(err)
	}

	for i := uint64(1); i <= highestSeenHeader; i++ {
		header := <-newHeads
		fmt.Printf("Got header %d\n", header.Number.Uint64())
		require.Equal(i, header.Number.Uint64())
	}
}
