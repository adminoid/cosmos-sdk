package baseapp

import (
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/adminoid/cosmos-sdk/types"
)

// ParamStore defines the interface the parameter store used by the BaseApp must
// fulfill.
type ParamStore interface {
	Get(ctx sdk.Context) (*tmproto.ConsensusParams, error)
	Has(ctx sdk.Context) bool
	Set(ctx sdk.Context, cp *tmproto.ConsensusParams)
}
