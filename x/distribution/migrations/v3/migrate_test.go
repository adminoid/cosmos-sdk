package v3_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/adminoid/cosmos-sdk/testutil"
	sdk "github.com/adminoid/cosmos-sdk/types"
	moduletestutil "github.com/adminoid/cosmos-sdk/types/module/testutil"
	"github.com/adminoid/cosmos-sdk/x/distribution"
	"github.com/adminoid/cosmos-sdk/x/distribution/exported"
	v3 "github.com/adminoid/cosmos-sdk/x/distribution/migrations/v3"
	"github.com/adminoid/cosmos-sdk/x/distribution/types"
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
	cdc := moduletestutil.MakeTestEncodingConfig(distribution.AppModuleBasic{}).Codec
	storeKey := sdk.NewKVStoreKey(v3.ModuleName)
	tKey := sdk.NewTransientStoreKey("transient_test")
	ctx := testutil.DefaultContext(storeKey, tKey)
	store := ctx.KVStore(storeKey)

	legacySubspace := newMockSubspace(types.DefaultParams())
	require.NoError(t, v3.MigrateStore(ctx, storeKey, legacySubspace, cdc))

	var res types.Params
	bz := store.Get(v3.ParamsKey)
	require.NoError(t, cdc.Unmarshal(bz, &res))
	require.Equal(t, legacySubspace.ps, res)
}
