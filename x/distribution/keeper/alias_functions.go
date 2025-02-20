package keeper

import (
	sdk "github.com/adminoid/cosmos-sdk/types"
	"github.com/adminoid/cosmos-sdk/x/distribution/types"
)

// get outstanding rewards
func (k Keeper) GetValidatorOutstandingRewardsCoins(ctx sdk.Context, val sdk.ValAddress) sdk.DecCoins {
	return k.GetValidatorOutstandingRewards(ctx, val).Rewards
}

// get the community coins
func (k Keeper) GetFeePoolCommunityCoins(ctx sdk.Context) sdk.DecCoins {
	return k.GetFeePool(ctx).CommunityPool
}

// GetDistributionAccount returns the distribution ModuleAccount
func (k Keeper) GetDistributionAccount(ctx sdk.Context) sdk.ModuleAccountI {
	return k.authKeeper.GetModuleAccount(ctx, types.ModuleName)
}
