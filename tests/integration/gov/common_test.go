package gov_test

import (
	"testing"

	"cosmossdk.io/math"
	"gotest.tools/v3/assert"

	sdk "github.com/adminoid/cosmos-sdk/types"
	authtypes "github.com/adminoid/cosmos-sdk/x/auth/types"
	"github.com/adminoid/cosmos-sdk/x/gov/types"
	v1 "github.com/adminoid/cosmos-sdk/x/gov/types/v1"
	"github.com/adminoid/cosmos-sdk/x/gov/types/v1beta1"
	stakingtypes "github.com/adminoid/cosmos-sdk/x/staking/types"
)

var (
	valTokens           = sdk.TokensFromConsensusPower(42, sdk.DefaultPowerReduction)
	TestProposal        = v1beta1.NewTextProposal("Test", "description")
	TestDescription     = stakingtypes.NewDescription("T", "E", "S", "T", "Z")
	TestCommissionRates = stakingtypes.NewCommissionRates(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec())
)

// mkTestLegacyContent creates a MsgExecLegacyContent for testing purposes.
func mkTestLegacyContent(t *testing.T) *v1.MsgExecLegacyContent {
	msgContent, err := v1.NewLegacyContent(TestProposal, authtypes.NewModuleAddress(types.ModuleName).String())
	assert.NilError(t, err)

	return msgContent
}
