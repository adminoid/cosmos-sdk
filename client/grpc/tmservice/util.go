package tmservice

import (
	sdk "github.com/adminoid/cosmos-sdk/types"
	tmprototypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

// convertHeader converts tendermint header to sdk header
func convertHeader(h tmprototypes.Header) Header {
	return Header{
		Version:            h.Version,
		ChainID:            h.ChainID,
		Height:             h.Height,
		Time:               h.Time,
		LastBlockId:        h.LastBlockId,
		ValidatorsHash:     h.ValidatorsHash,
		NextValidatorsHash: h.NextValidatorsHash,
		ConsensusHash:      h.ConsensusHash,
		AppHash:            h.AppHash,
		DataHash:           h.DataHash,
		EvidenceHash:       h.EvidenceHash,
		LastResultsHash:    h.LastResultsHash,
		LastCommitHash:     h.LastCommitHash,
		ProposerAddress:    sdk.ValAddress(h.ProposerAddress).String(),
	}
}

// convertBlock converts tendermint block to sdk block
func convertBlock(tmblock *tmprototypes.Block) *Block {
	b := new(Block)

	b.Header = convertHeader(tmblock.Header)
	b.LastCommit = tmblock.LastCommit
	b.Data = tmblock.Data
	b.Evidence = tmblock.Evidence

	return b
}
