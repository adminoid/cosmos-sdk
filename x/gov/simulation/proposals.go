package simulation

import (
	"math/rand"

	sdk "github.com/adminoid/cosmos-sdk/types"
	simtypes "github.com/adminoid/cosmos-sdk/types/simulation"
	"github.com/adminoid/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/adminoid/cosmos-sdk/x/simulation"
)

// OpWeightSubmitTextProposal app params key for text proposal
const OpWeightSubmitTextProposal = "op_weight_submit_text_proposal"

// ProposalContents defines the module weighted proposals' contents
func ProposalContents() []simtypes.WeightedProposalContent {
	return []simtypes.WeightedProposalContent{
		simulation.NewWeightedProposalContent(
			OpWeightMsgDeposit,
			DefaultWeightTextProposal,
			SimulateTextProposalContent,
		),
	}
}

// SimulateTextProposalContent returns a random text proposal content.
func SimulateTextProposalContent(r *rand.Rand, _ sdk.Context, _ []simtypes.Account) simtypes.Content {
	return v1beta1.NewTextProposal(
		simtypes.RandStringOfLength(r, 140),
		simtypes.RandStringOfLength(r, 5000),
	)
}
