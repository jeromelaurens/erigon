package stagedsync

import (
	"fmt"

	"github.com/VictoriaMetrics/metrics"
	"github.com/huandu/xstrings"
	"github.com/jeromelaurens/erigon/eth/stagedsync/stages"
	"github.com/ledgerwatch/erigon-lib/kv"
)

var syncMetrics = map[stages.SyncStage]*metrics.Counter{}

func init() {
	for _, v := range stages.AllStages {
		syncMetrics[v] = metrics.GetOrCreateCounter(
			fmt.Sprintf(
				`sync{stage="%s"}`,
				xstrings.ToSnakeCase(string(v)),
			),
		)
	}
}

// UpdateMetrics - need update metrics manually because current "metrics" package doesn't support labels
// need to fix it in future
func UpdateMetrics(tx kv.Tx) error {
	for id, m := range syncMetrics {
		progress, err := stages.GetStageProgress(tx, id)
		if err != nil {
			return err
		}
		m.Set(progress)
	}
	return nil
}
