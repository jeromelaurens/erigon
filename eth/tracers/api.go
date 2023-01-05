package tracers

import (
	"github.com/jeromelaurens/erigon/eth/tracers/logger"
	"github.com/jeromelaurens/erigon/turbo/adapter/ethapi"
)

// TraceConfig holds extra parameters to trace functions.
type TraceConfig struct {
	*logger.LogConfig
	Tracer         *string
	Timeout        *string
	Reexec         *uint64
	NoRefunds      *bool // Turns off gas refunds when tracing
	StateOverrides *ethapi.StateOverrides
}
