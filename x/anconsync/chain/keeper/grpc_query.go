package keeper

import (
	"github.com/anconprotocol/node/x/anconsync/chain/types"
)

var _ types.QueryServer = Keeper{}
