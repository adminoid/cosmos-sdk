package v2_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/adminoid/cosmos-sdk/testutil"
	sdk "github.com/adminoid/cosmos-sdk/types"
	moduletestutil "github.com/adminoid/cosmos-sdk/types/module/testutil"
	"github.com/adminoid/cosmos-sdk/x/mint"
	"github.com/adminoid/cosmos-sdk/x/mint/exported"
	v2 "github.com/adminoid/cosmos-sdk/x/mint/migrations/v2"
	"github.com/adminoid/cosmos-sdk/x/mint/types"
)

type mockSubspace struct {
	ps types.Params
}

func newMockSubspace(ps types.Params) mockSubspace {
	return mockSubspace{ps: ps}
}

func (ms mockSubspace) GetParamSet(ctx sdk.Context, ps exported.ParamSet) {
	*ps.(*types.Params) = ms.ps
}

func TestMigrate(t *testing.T) {
	encCfg := moduletestutil.MakeTestEncodingConfig(mint.AppModuleBasic{})
	cdc := encCfg.Codec

	storeKey := sdk.NewKVStoreKey(v2.ModuleName)
	tKey := sdk.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)
	store := ctx.KVStore(storeKey)

	legacySubspace := newMockSubspace(types.DefaultParams())
	require.NoError(t, v2.Migrate(ctx, store, legacySubspace, cdc))

	var res types.Params
	bz := store.Get(v2.ParamsKey)
	require.NoError(t, cdc.Unmarshal(bz, &res))
	require.Equal(t, legacySubspace.ps, res)
}
