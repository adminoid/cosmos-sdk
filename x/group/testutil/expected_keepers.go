// This file only used to generate mocks

package testutil

import (
	sdk "github.com/adminoid/cosmos-sdk/types"
	bank "github.com/adminoid/cosmos-sdk/x/bank/types"
	"github.com/adminoid/cosmos-sdk/x/group"
)

// AccountKeeper extends `AccountKeeper` from expected_keepers.
type AccountKeeper interface {
	group.AccountKeeper
}

// BankKeeper extends bank `MsgServer` to mock `Send` and to register handlers in MsgServiceRouter
type BankKeeper interface {
	group.BankKeeper
	bank.MsgServer

	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}
