package keeper_test

import (
	"testing"

	simtestutil "github.com/adminoid/cosmos-sdk/testutil/sims"
	sdk "github.com/adminoid/cosmos-sdk/types"
	"github.com/adminoid/cosmos-sdk/x/auth/keeper"
	"github.com/adminoid/cosmos-sdk/x/auth/testutil"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func BenchmarkAccountMapperGetAccountFound(b *testing.B) {
	b.ReportAllocs()
	var accountKeeper keeper.AccountKeeper
	app, err := simtestutil.Setup(testutil.AppConfig, &accountKeeper)
	require.NoError(b, err)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	// assumes b.N < 2**24
	for i := 0; i < b.N; i++ {
		arr := []byte{byte((i & 0xFF0000) >> 16), byte((i & 0xFF00) >> 8), byte(i & 0xFF)}
		addr := sdk.AccAddress(arr)
		acc := accountKeeper.NewAccountWithAddress(ctx, addr)
		accountKeeper.SetAccount(ctx, acc)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		arr := []byte{byte((i & 0xFF0000) >> 16), byte((i & 0xFF00) >> 8), byte(i & 0xFF)}
		accountKeeper.GetAccount(ctx, sdk.AccAddress(arr))
	}
}

func BenchmarkAccountMapperSetAccount(b *testing.B) {
	b.ReportAllocs()
	var accountKeeper keeper.AccountKeeper
	app, err := simtestutil.Setup(testutil.AppConfig, &accountKeeper)
	require.NoError(b, err)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	b.ResetTimer()

	// assumes b.N < 2**24
	for i := 0; i < b.N; i++ {
		arr := []byte{byte((i & 0xFF0000) >> 16), byte((i & 0xFF00) >> 8), byte(i & 0xFF)}
		addr := sdk.AccAddress(arr)
		acc := accountKeeper.NewAccountWithAddress(ctx, addr)
		accountKeeper.SetAccount(ctx, acc)
	}
}
