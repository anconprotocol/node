package anconsync

import (
	"github.com/0xPolygon/polygon-sdk/state"
)

func GetHooks(s Storage) func(t *state.Transition) {
	return func(t *state.Transition) {
		// logs := t.Txn().Logs()

		// // for _, log := range logs {
		// // 	log.Topics
		// // }
	}
}
