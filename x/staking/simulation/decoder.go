package simulation

import (
	"bytes"
	"fmt"

	"github.com/adminoid/cosmos-sdk/codec"
	sdk "github.com/adminoid/cosmos-sdk/types"
	"github.com/adminoid/cosmos-sdk/types/kv"
	"github.com/adminoid/cosmos-sdk/x/staking/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding staking type.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.LastTotalPowerKey):
			var powerA, powerB sdk.IntProto

			cdc.MustUnmarshal(kvA.Value, &powerA)
			cdc.MustUnmarshal(kvB.Value, &powerB)

			return fmt.Sprintf("%v\n%v", powerA, powerB)
		case bytes.Equal(kvA.Key[:1], types.ValidatorsKey):
			var validatorA, validatorB types.Validator

			cdc.MustUnmarshal(kvA.Value, &validatorA)
			cdc.MustUnmarshal(kvB.Value, &validatorB)

			return fmt.Sprintf("%v\n%v", validatorA, validatorB)
		case bytes.Equal(kvA.Key[:1], types.LastValidatorPowerKey),
			bytes.Equal(kvA.Key[:1], types.ValidatorsByConsAddrKey),
			bytes.Equal(kvA.Key[:1], types.ValidatorsByPowerIndexKey):
			return fmt.Sprintf("%v\n%v", sdk.ValAddress(kvA.Value), sdk.ValAddress(kvB.Value))

		case bytes.Equal(kvA.Key[:1], types.DelegationKey):
			var delegationA, delegationB types.Delegation

			cdc.MustUnmarshal(kvA.Value, &delegationA)
			cdc.MustUnmarshal(kvB.Value, &delegationB)

			return fmt.Sprintf("%v\n%v", delegationA, delegationB)
		case bytes.Equal(kvA.Key[:1], types.UnbondingDelegationKey),
			bytes.Equal(kvA.Key[:1], types.UnbondingDelegationByValIndexKey):
			var ubdA, ubdB types.UnbondingDelegation

			cdc.MustUnmarshal(kvA.Value, &ubdA)
			cdc.MustUnmarshal(kvB.Value, &ubdB)

			return fmt.Sprintf("%v\n%v", ubdA, ubdB)
		case bytes.Equal(kvA.Key[:1], types.RedelegationKey),
			bytes.Equal(kvA.Key[:1], types.RedelegationByValSrcIndexKey):
			var redA, redB types.Redelegation

			cdc.MustUnmarshal(kvA.Value, &redA)
			cdc.MustUnmarshal(kvB.Value, &redB)

			return fmt.Sprintf("%v\n%v", redA, redB)
		case bytes.Equal(kvA.Key[:1], types.ParamsKey):
			var paramsA, paramsB types.Params

			cdc.MustUnmarshal(kvA.Value, &paramsA)
			cdc.MustUnmarshal(kvB.Value, &paramsB)

			return fmt.Sprintf("%v\n%v", paramsA, paramsB)
		default:
			panic(fmt.Sprintf("invalid staking key prefix %X", kvA.Key[:1]))
		}
	}
}
