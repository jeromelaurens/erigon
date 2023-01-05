package handshake

import (
	"bytes"
	"context"
	"sync"

	"github.com/jeromelaurens/erigon/cl/clparams"
	"github.com/jeromelaurens/erigon/cl/cltypes"
	"github.com/jeromelaurens/erigon/cmd/sentinel/sentinel/communication"
	"github.com/jeromelaurens/erigon/cmd/sentinel/sentinel/communication/ssz_snappy"
	"github.com/jeromelaurens/erigon/common"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/core/host"
	"go.uber.org/zap/buffer"
)

// HandShaker is the data type which will handle handshakes and determine if
// The peer is worth connecting to or not.
type HandShaker struct {
	ctx context.Context
	// Status object to send over.
	status        *cltypes.Status // Contains status object for handshakes
	set           bool
	rule          RuleFunc // Method that determine if peer is worth connecting to or not.
	host          host.Host
	genesisConfig *clparams.GenesisConfig
	beaconConfig  *clparams.BeaconChainConfig

	mu sync.Mutex
}

func New(ctx context.Context, genesisConfig *clparams.GenesisConfig, beaconConfig *clparams.BeaconChainConfig, host host.Host, rule RuleFunc) *HandShaker {
	return &HandShaker{
		ctx:           ctx,
		rule:          rule,
		host:          host,
		genesisConfig: genesisConfig,
		beaconConfig:  beaconConfig,
		status:        &cltypes.Status{},
	}
}

// SetStatus sets the current network status against which we can validate peers.
func (h *HandShaker) SetStatus(status *cltypes.Status) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.set = true
	h.status = status
}

// Status returns the underlying status (only for giving out responses)
func (h *HandShaker) Status() *cltypes.Status {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.status
}

// Set returns the underlying status (only for giving out responses)
func (h *HandShaker) IsSet() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.set
}

func (h *HandShaker) ValidatePeer(id peer.ID) bool {
	// Unprotected if it is not set
	if !h.IsSet() {
		return true
	}
	status := h.Status()
	// Encode our status
	var buffer buffer.Buffer
	if err := ssz_snappy.EncodeAndWrite(&buffer, status); err != nil {
		return false
	}

	data := common.CopyBytes(buffer.Bytes())

	response, errResponse, err := communication.SendRequestRawToPeer(h.ctx, h.host, data, communication.StatusProtocolV1, id)
	if err != nil || errResponse {
		return false
	}
	responseStatus := &cltypes.Status{}

	if err := ssz_snappy.DecodeAndReadNoForkDigest(bytes.NewReader(response), responseStatus); err != nil {
		return false
	}

	return h.rule(responseStatus, status, h.genesisConfig, h.beaconConfig)
}
