package exported

import (
	sdk "github.com/adminoid/cosmos-sdk/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

type (
	// ParamStore defines an interface that implements the legacy x/params Subspace
	// type.
	//
	// NOTE: This is used solely for migration of x/params managed parameters.
	ParamStore interface {
		Get(ctx sdk.Context, key []byte, ptr interface{})
	}

	// ConsensusParamSetter defines the interface fulfilled by BaseApp's
	// ParamStore which allows setting its appVersion field.
	ConsensusParamSetter interface {
		Get(ctx sdk.Context) (*tmproto.ConsensusParams, error)
		Has(ctx sdk.Context) bool
		Set(ctx sdk.Context, cp *tmproto.ConsensusParams)
	}
)
