package serenity

import (
	"math/big"
	"testing"

	"github.com/jeromelaurens/erigon/common"
	"github.com/jeromelaurens/erigon/consensus"
	"github.com/jeromelaurens/erigon/core/types"
	"github.com/jeromelaurens/erigon/params"
)

type readerMock struct{}

func (r readerMock) Config() *params.ChainConfig {
	return nil
}

func (r readerMock) CurrentHeader() *types.Header {
	return nil
}

func (r readerMock) GetHeader(common.Hash, uint64) *types.Header {
	return nil
}

func (r readerMock) GetHeaderByNumber(uint64) *types.Header {
	return nil
}

func (r readerMock) GetHeaderByHash(common.Hash) *types.Header {
	return nil
}

func (r readerMock) GetTd(common.Hash, uint64) *big.Int {
	return nil
}

// The thing only that changes beetwen normal ethash checks other than POW, is difficulty
// and nonce so we are gonna test those
func TestVerifyHeaderDifficulty(t *testing.T) {
	header := &types.Header{
		Difficulty: big.NewInt(1),
		Time:       1,
	}

	parent := &types.Header{}

	var eth1Engine consensus.Engine
	serenity := New(eth1Engine)

	err := serenity.verifyHeader(readerMock{}, header, parent)
	if err != errInvalidDifficulty {
		if err != nil {
			t.Fatalf("Serenity should not accept non-zero difficulty, got %s", err.Error())
		} else {
			t.Fatalf("Serenity should not accept non-zero difficulty")
		}
	}
}

func TestVerifyHeaderNonce(t *testing.T) {
	header := &types.Header{
		Nonce:      types.BlockNonce{1, 0, 0, 0, 0, 0, 0, 0},
		Difficulty: big.NewInt(0),
		Time:       1,
	}

	parent := &types.Header{}

	var eth1Engine consensus.Engine
	serenity := New(eth1Engine)

	err := serenity.verifyHeader(readerMock{}, header, parent)
	if err != errInvalidNonce {
		if err != nil {
			t.Fatalf("Serenity should not accept non-zero difficulty, got %s", err.Error())
		} else {
			t.Fatalf("Serenity should not accept non-zero difficulty")
		}
	}
}
