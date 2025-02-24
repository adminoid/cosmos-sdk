package simulation

import (
	"bytes"
	"fmt"

	"github.com/adminoid/cosmos-sdk/codec"
	sdk "github.com/adminoid/cosmos-sdk/types"
	"github.com/adminoid/cosmos-sdk/types/kv"
	"github.com/adminoid/cosmos-sdk/x/distribution/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding distribution type.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.FeePoolKey):
			var feePoolA, feePoolB types.FeePool
			cdc.MustUnmarshal(kvA.Value, &feePoolA)
			cdc.MustUnmarshal(kvB.Value, &feePoolB)
			return fmt.Sprintf("%v\n%v", feePoolA, feePoolB)

		case bytes.Equal(kvA.Key[:1], types.ProposerKey):
			return fmt.Sprintf("%v\n%v", sdk.ConsAddress(kvA.Value), sdk.ConsAddress(kvB.Value))

		case bytes.Equal(kvA.Key[:1], types.ValidatorOutstandingRewardsPrefix):
			var rewardsA, rewardsB types.ValidatorOutstandingRewards
			cdc.MustUnmarshal(kvA.Value, &rewardsA)
			cdc.MustUnmarshal(kvB.Value, &rewardsB)
			return fmt.Sprintf("%v\n%v", rewardsA, rewardsB)

		case bytes.Equal(kvA.Key[:1], types.DelegatorWithdrawAddrPrefix):
			return fmt.Sprintf("%v\n%v", sdk.AccAddress(kvA.Value), sdk.AccAddress(kvB.Value))

		case bytes.Equal(kvA.Key[:1], types.DelegatorStartingInfoPrefix):
			var infoA, infoB types.DelegatorStartingInfo
			cdc.MustUnmarshal(kvA.Value, &infoA)
			cdc.MustUnmarshal(kvB.Value, &infoB)
			return fmt.Sprintf("%v\n%v", infoA, infoB)

		case bytes.Equal(kvA.Key[:1], types.ValidatorHistoricalRewardsPrefix):
			var rewardsA, rewardsB types.ValidatorHistoricalRewards
			cdc.MustUnmarshal(kvA.Value, &rewardsA)
			cdc.MustUnmarshal(kvB.Value, &rewardsB)
			return fmt.Sprintf("%v\n%v", rewardsA, rewardsB)

		case bytes.Equal(kvA.Key[:1], types.ValidatorCurrentRewardsPrefix):
			var rewardsA, rewardsB types.ValidatorCurrentRewards
			cdc.MustUnmarshal(kvA.Value, &rewardsA)
			cdc.MustUnmarshal(kvB.Value, &rewardsB)
			return fmt.Sprintf("%v\n%v", rewardsA, rewardsB)

		case bytes.Equal(kvA.Key[:1], types.ValidatorAccumulatedCommissionPrefix):
			var commissionA, commissionB types.ValidatorAccumulatedCommission
			cdc.MustUnmarshal(kvA.Value, &commissionA)
			cdc.MustUnmarshal(kvB.Value, &commissionB)
			return fmt.Sprintf("%v\n%v", commissionA, commissionB)

		case bytes.Equal(kvA.Key[:1], types.ValidatorSlashEventPrefix):
			var eventA, eventB types.ValidatorSlashEvent
			cdc.MustUnmarshal(kvA.Value, &eventA)
			cdc.MustUnmarshal(kvB.Value, &eventB)
			return fmt.Sprintf("%v\n%v", eventA, eventB)

		default:
			panic(fmt.Sprintf("invalid distribution key prefix %X", kvA.Key[:1]))
		}
	}
}
