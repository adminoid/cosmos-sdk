package types

import (
	"fmt"

	sdk "github.com/adminoid/cosmos-sdk/types"
	paramtypes "github.com/adminoid/cosmos-sdk/x/params/types"
)

// ParamStoreKeyConstantFee is the constant fee parameter
var ParamStoreKeyConstantFee = []byte("ConstantFee")

// type declaration for parameters
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable(
		paramtypes.NewParamSetPair(ParamStoreKeyConstantFee, sdk.Coin{}, validateConstantFee),
	)
}

func validateConstantFee(i interface{}) error {
	v, ok := i.(sdk.Coin)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if !v.IsValid() {
		return fmt.Errorf("invalid constant fee: %s", v)
	}

	return nil
}
