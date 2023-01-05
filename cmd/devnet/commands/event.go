package commands

import (
	"fmt"

	"github.com/jeromelaurens/erigon/cmd/devnet/services"
	"github.com/jeromelaurens/erigon/common"
)

func callSubscribeToNewHeads(hash common.Hash) {
	_, err := services.SearchBlockForTransactionHash(hash)
	if err != nil {
		fmt.Printf("FAILURE => %v\n", err)
		return
	}
}
