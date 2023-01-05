package transition

import (
	"github.com/jeromelaurens/erigon/cl/clparams"
	"github.com/jeromelaurens/erigon/cmd/erigon-cl/core/state"
)

// StateTransistor takes care of state transition
type StateTransistor struct {
	state         *state.BeaconState
	beaconConfig  *clparams.BeaconChainConfig
	genesisConfig *clparams.GenesisConfig
}

func New(state *state.BeaconState, beaconConfig *clparams.BeaconChainConfig, genesisConfig *clparams.GenesisConfig) *StateTransistor {
	return &StateTransistor{
		state:         state,
		beaconConfig:  beaconConfig,
		genesisConfig: genesisConfig,
	}
}
