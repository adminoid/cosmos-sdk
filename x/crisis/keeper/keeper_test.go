package keeper_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/adminoid/cosmos-sdk/testutil"
	sdk "github.com/adminoid/cosmos-sdk/types"
	moduletestutil "github.com/adminoid/cosmos-sdk/types/module/testutil"
	"github.com/adminoid/cosmos-sdk/x/crisis"
	"github.com/adminoid/cosmos-sdk/x/crisis/keeper"
	crisistestutil "github.com/adminoid/cosmos-sdk/x/crisis/testutil"
	"github.com/adminoid/cosmos-sdk/x/crisis/types"
)

func TestLogger(t *testing.T) {
	ctrl := gomock.NewController(t)
	supplyKeeper := crisistestutil.NewMockSupplyKeeper(ctrl)

	key := sdk.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(t, key, sdk.NewTransientStoreKey("transient_test"))
	encCfg := moduletestutil.MakeTestEncodingConfig(crisis.AppModuleBasic{})
	keeper := keeper.NewKeeper(encCfg.Codec, key, 5, supplyKeeper, "", "")

	require.Equal(t,
		testCtx.Ctx.Logger().With("module", "x/"+types.ModuleName),
		keeper.Logger(testCtx.Ctx))
}

func TestInvariants(t *testing.T) {
	ctrl := gomock.NewController(t)
	supplyKeeper := crisistestutil.NewMockSupplyKeeper(ctrl)

	key := sdk.NewKVStoreKey(types.StoreKey)
	encCfg := moduletestutil.MakeTestEncodingConfig(crisis.AppModuleBasic{})
	keeper := keeper.NewKeeper(encCfg.Codec, key, 5, supplyKeeper, "", "")
	require.Equal(t, keeper.InvCheckPeriod(), uint(5))

	orgInvRoutes := keeper.Routes()
	keeper.RegisterRoute("testModule", "testRoute", func(sdk.Context) (string, bool) { return "", false })
	invar := keeper.Invariants()
	require.Equal(t, len(invar), len(orgInvRoutes)+1)
}

func TestAssertInvariants(t *testing.T) {
	ctrl := gomock.NewController(t)
	supplyKeeper := crisistestutil.NewMockSupplyKeeper(ctrl)

	key := sdk.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(t, key, sdk.NewTransientStoreKey("transient_test"))
	encCfg := moduletestutil.MakeTestEncodingConfig(crisis.AppModuleBasic{})
	keeper := keeper.NewKeeper(encCfg.Codec, key, 5, supplyKeeper, "", "")

	keeper.RegisterRoute("testModule", "testRoute1", func(sdk.Context) (string, bool) { return "", false })
	require.NotPanics(t, func() { keeper.AssertInvariants(testCtx.Ctx) })

	keeper.RegisterRoute("testModule", "testRoute2", func(sdk.Context) (string, bool) { return "", true })
	require.Panics(t, func() { keeper.AssertInvariants(testCtx.Ctx) })
}
