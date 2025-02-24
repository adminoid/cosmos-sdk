// This file only used to generate mocks

package testutil

import (
	math "cosmossdk.io/math"

	sdk "github.com/adminoid/cosmos-sdk/types"
	bankkeeper "github.com/adminoid/cosmos-sdk/x/bank/keeper"
	"github.com/adminoid/cosmos-sdk/x/gov/types"
)

// AccountKeeper extends gov's actual expected AccountKeeper with additional
// methods used in tests.
type AccountKeeper interface {
	types.AccountKeeper

	IterateAccounts(ctx sdk.Context, cb func(account sdk.AccountI) (stop bool))
}

// BankKeeper extends gov's actual expected BankKeeper with additional
// methods used in tests.
type BankKeeper interface {
	bankkeeper.Keeper
}

// StakingKeeper extends gov's actual expected StakingKeeper with additional
// methods used in tests.
type StakingKeeper interface {
	types.StakingKeeper

	BondDenom(ctx sdk.Context) string
	TokensFromConsensusPower(ctx sdk.Context, power int64) math.Int
}
