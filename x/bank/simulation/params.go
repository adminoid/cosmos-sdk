package simulation

import (
	"fmt"
	"math/rand"

	"github.com/adminoid/cosmos-sdk/x/simulation"

	simtypes "github.com/adminoid/cosmos-sdk/types/simulation"
	"github.com/adminoid/cosmos-sdk/x/bank/types"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyDefaultSendEnabled),
			func(r *rand.Rand) string {
				return fmt.Sprintf("%v", RandomGenesisDefaultSendEnabledParam(r))
			},
		),
	}
}
