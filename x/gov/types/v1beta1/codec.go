package v1beta1

import (
	"github.com/adminoid/cosmos-sdk/codec"
	"github.com/adminoid/cosmos-sdk/codec/legacy"
	codectypes "github.com/adminoid/cosmos-sdk/codec/types"
	sdk "github.com/adminoid/cosmos-sdk/types"
	"github.com/adminoid/cosmos-sdk/types/msgservice"
	authzcodec "github.com/adminoid/cosmos-sdk/x/authz/codec"
	govcodec "github.com/adminoid/cosmos-sdk/x/gov/codec"
	groupcodec "github.com/adminoid/cosmos-sdk/x/group/codec"
)

// RegisterLegacyAminoCodec registers all the necessary types and interfaces for the
// governance module.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*Content)(nil), nil)
	legacy.RegisterAminoMsg(cdc, &MsgSubmitProposal{}, "cosmos-sdk/MsgSubmitProposal")
	legacy.RegisterAminoMsg(cdc, &MsgDeposit{}, "cosmos-sdk/MsgDeposit")
	legacy.RegisterAminoMsg(cdc, &MsgVote{}, "cosmos-sdk/MsgVote")
	legacy.RegisterAminoMsg(cdc, &MsgVoteWeighted{}, "cosmos-sdk/MsgVoteWeighted")
	cdc.RegisterConcrete(&TextProposal{}, "cosmos-sdk/TextProposal", nil)
}

// RegisterInterfaces registers the interfaces types with the Interface Registry.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSubmitProposal{},
		&MsgVote{},
		&MsgVoteWeighted{},
		&MsgDeposit{},
	)
	registry.RegisterInterface(
		"cosmos.gov.v1beta1.Content",
		(*Content)(nil),
		&TextProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

func init() {
	// Register all Amino interfaces and concrete types on the authz  and gov Amino codec so that this can later be
	// used to properly serialize MsgGrant, MsgExec and MsgSubmitProposal instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
	RegisterLegacyAminoCodec(govcodec.Amino)
	RegisterLegacyAminoCodec(groupcodec.Amino)
}
