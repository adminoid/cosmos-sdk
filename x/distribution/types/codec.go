package types

import (
	"github.com/adminoid/cosmos-sdk/codec"
	"github.com/adminoid/cosmos-sdk/codec/legacy"
	"github.com/adminoid/cosmos-sdk/codec/types"
	cryptocodec "github.com/adminoid/cosmos-sdk/crypto/codec"
	sdk "github.com/adminoid/cosmos-sdk/types"
	"github.com/adminoid/cosmos-sdk/types/msgservice"
	authzcodec "github.com/adminoid/cosmos-sdk/x/authz/codec"
	govcodec "github.com/adminoid/cosmos-sdk/x/gov/codec"
	govtypes "github.com/adminoid/cosmos-sdk/x/gov/types/v1beta1"
	groupcodec "github.com/adminoid/cosmos-sdk/x/group/codec"
)

// RegisterLegacyAminoCodec registers the necessary x/distribution interfaces
// and concrete types on the provided LegacyAmino codec. These types are used
// for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawDelegatorReward{}, "cosmos-sdk/MsgWithdrawDelegationReward")
	legacy.RegisterAminoMsg(cdc, &MsgWithdrawValidatorCommission{}, "cosmos-sdk/MsgWithdrawValCommission")
	legacy.RegisterAminoMsg(cdc, &MsgSetWithdrawAddress{}, "cosmos-sdk/MsgModifyWithdrawAddress")
	legacy.RegisterAminoMsg(cdc, &MsgFundCommunityPool{}, "cosmos-sdk/MsgFundCommunityPool")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateParams{}, "cosmos-sdk/distribution/MsgUpdateParams")
	legacy.RegisterAminoMsg(cdc, &MsgCommunityPoolSpend{}, "cosmos-sdk/distr/MsgCommunityPoolSpend")
	legacy.RegisterAminoMsg(cdc, &MsgDepositValidatorRewardsPool{}, "cosmos-sdk/distr/MsgDepositValRewards")

	cdc.RegisterConcrete(Params{}, "cosmos-sdk/x/distribution/Params", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgWithdrawDelegatorReward{},
		&MsgWithdrawValidatorCommission{},
		&MsgSetWithdrawAddress{},
		&MsgFundCommunityPool{},
		&MsgUpdateParams{},
		&MsgCommunityPoolSpend{},
		&MsgDepositValidatorRewardsPool{},
	)

	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&CommunityPoolSpendProposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)

	// Register all Amino interfaces and concrete types on the authz  and gov Amino codec
	// so that this can later be used to properly serialize MsgGrant and MsgExec
	// instances.
	RegisterLegacyAminoCodec(authzcodec.Amino)
	RegisterLegacyAminoCodec(govcodec.Amino)
	RegisterLegacyAminoCodec(groupcodec.Amino)
}
