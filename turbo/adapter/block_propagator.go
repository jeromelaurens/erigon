package adapter

import (
	"context"
	"math/big"

	"github.com/jeromelaurens/erigon/core/types"
)

type BlockPropagator func(ctx context.Context, block *types.Block, td *big.Int)
