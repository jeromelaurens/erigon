package commands

import (
	"context"
	"encoding/json"

	"github.com/jeromelaurens/erigon/cmd/rpcdaemon/cli/httpcfg"
	"github.com/jeromelaurens/erigon/common"
	"github.com/jeromelaurens/erigon/common/hexutil"
	"github.com/jeromelaurens/erigon/rpc"
	jsoniter "github.com/json-iterator/go"
	"github.com/ledgerwatch/erigon-lib/kv"
)

// TraceAPI RPC interface into tracing API
type TraceAPI interface {
	// Ad-hoc (see ./trace_adhoc.go)
	ReplayBlockTransactions(ctx context.Context, blockNr rpc.BlockNumberOrHash, traceTypes []string) ([]*TraceCallResult, error)
	ReplayTransaction(ctx context.Context, txHash common.Hash, traceTypes []string) (*TraceCallResult, error)
	Call(ctx context.Context, call TraceCallParam, types []string, blockNr *rpc.BlockNumberOrHash) (*TraceCallResult, error)
	CallMany(ctx context.Context, calls json.RawMessage, blockNr *rpc.BlockNumberOrHash) ([]*TraceCallResult, error)
	RawTransaction(ctx context.Context, txHash common.Hash, traceTypes []string) ([]interface{}, error)

	// Filtering (see ./trace_filtering.go)
	Transaction(ctx context.Context, txHash common.Hash) (ParityTraces, error)
	Get(ctx context.Context, txHash common.Hash, txIndicies []hexutil.Uint64) (*ParityTrace, error)
	Block(ctx context.Context, blockNr rpc.BlockNumber) (ParityTraces, error)
	Filter(ctx context.Context, req TraceFilterRequest, stream *jsoniter.Stream) error
}

// TraceAPIImpl is implementation of the TraceAPI interface based on remote Db access
type TraceAPIImpl struct {
	*BaseAPI
	kv            kv.RoDB
	maxTraces     uint64
	gasCap        uint64
	compatibility bool // Bug for bug compatiblity with OpenEthereum
}

// NewTraceAPI returns NewTraceAPI instance
func NewTraceAPI(base *BaseAPI, kv kv.RoDB, cfg *httpcfg.HttpCfg) *TraceAPIImpl {
	return &TraceAPIImpl{
		BaseAPI:       base,
		kv:            kv,
		maxTraces:     cfg.MaxTraces,
		gasCap:        cfg.Gascap,
		compatibility: cfg.TraceCompatibility,
	}
}
