package health

import (
	"context"

	"github.com/jeromelaurens/erigon/common/hexutil"
	"github.com/jeromelaurens/erigon/rpc"
)

type NetAPI interface {
	PeerCount(_ context.Context) (hexutil.Uint, error)
}

type EthAPI interface {
	GetBlockByNumber(_ context.Context, number rpc.BlockNumber, fullTx bool) (map[string]interface{}, error)
	Syncing(ctx context.Context) (interface{}, error)
}
