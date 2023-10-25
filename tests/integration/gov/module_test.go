package gov_test

import (
	"testing"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"gotest.tools/v3/assert"

	"github.com/adminoid/cosmos-sdk/testutil/configurator"
	simtestutil "github.com/adminoid/cosmos-sdk/testutil/sims"
	authkeeper "github.com/adminoid/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/adminoid/cosmos-sdk/x/auth/types"
	"github.com/adminoid/cosmos-sdk/x/gov/types"
	_ "github.com/adminoid/cosmos-sdk/x/mint"
)

func TestItCreatesModuleAccountOnInitBlock(t *testing.T) {
	var accountKeeper authkeeper.AccountKeeper
	app, err := simtestutil.SetupAtGenesis(
		configurator.NewAppConfig(
			configurator.ParamsModule(),
			configurator.AuthModule(),
			configurator.StakingModule(),
			configurator.BankModule(),
			configurator.GovModule(),
			configurator.ConsensusModule(),
		),
		&accountKeeper,
	)
	assert.NilError(t, err)

	ctx := app.BaseApp.NewContext(false, tmproto.Header{})
	acc := accountKeeper.GetAccount(ctx, authtypes.NewModuleAddress(types.ModuleName))
	assert.Assert(t, acc != nil)
}
