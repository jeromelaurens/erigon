package commands

import (
	"context"
	"fmt"

	"github.com/jeromelaurens/erigon/common"
	"github.com/jeromelaurens/erigon/core/state"
	"github.com/jeromelaurens/erigon/core/systemcontracts"
	"github.com/jeromelaurens/erigon/rpc"
	"github.com/jeromelaurens/erigon/turbo/rpchelper"
)

func (api *OtterscanAPIImpl) HasCode(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (bool, error) {
	tx, err := api.db.BeginRo(ctx)
	if err != nil {
		return false, fmt.Errorf("hasCode cannot open tx: %w", err)
	}
	defer tx.Rollback()

	blockNumber, _, _, err := rpchelper.GetBlockNumber(blockNrOrHash, tx, api.filters)
	if err != nil {
		return false, err
	}
	chainConfig, err := api.chainConfig(tx)
	if err != nil {
		return false, err
	}

	reader := state.NewPlainState(tx, blockNumber, systemcontracts.SystemContractCodeLookup[chainConfig.ChainName])
	acc, err := reader.ReadAccountData(address)
	if acc == nil || err != nil {
		return false, err
	}
	return !acc.IsEmptyCodeHash(), nil
}
