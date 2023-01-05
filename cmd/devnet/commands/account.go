package commands

import (
	"fmt"

	"github.com/jeromelaurens/erigon/cmd/devnet/models"
	"github.com/jeromelaurens/erigon/cmd/devnet/requests"

	"github.com/jeromelaurens/erigon/common"
)

const (
	addr = "0x71562b71999873DB5b286dF957af199Ec94617F7"
)

func callGetBalance(addr string, blockNum models.BlockNumber, checkBal uint64) {
	fmt.Printf("Getting balance for address: %q...\n", addr)
	address := common.HexToAddress(addr)
	bal, err := requests.GetBalance(models.ReqId, address, blockNum)
	if err != nil {
		fmt.Printf("FAILURE => %v\n", err)
		return
	}

	if checkBal > 0 && checkBal != bal {
		fmt.Printf("FAILURE => Balance should be %d, got %d\n", checkBal, bal)
		return
	}

	fmt.Printf("SUCCESS => Balance: %d\n", bal)
}
