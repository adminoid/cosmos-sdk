package codec

import (
	"github.com/adminoid/cosmos-sdk/codec"
	cryptocodec "github.com/adminoid/cosmos-sdk/crypto/codec"
	sdk "github.com/adminoid/cosmos-sdk/types"
)

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(Amino)
)

func init() {
	cryptocodec.RegisterCrypto(Amino)
	codec.RegisterEvidences(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
}
