package tx

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/adminoid/cosmos-sdk/codec"
	codectypes "github.com/adminoid/cosmos-sdk/codec/types"
	"github.com/adminoid/cosmos-sdk/std"
	"github.com/adminoid/cosmos-sdk/testutil/testdata"
	sdk "github.com/adminoid/cosmos-sdk/types"
	txtestutil "github.com/adminoid/cosmos-sdk/x/auth/tx/testutil"
)

func TestGenerator(t *testing.T) {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	interfaceRegistry.RegisterImplementations((*sdk.Msg)(nil), &testdata.TestMsg{})
	protoCodec := codec.NewProtoCodec(interfaceRegistry)
	suite.Run(t, txtestutil.NewTxConfigTestSuite(NewTxConfig(protoCodec, DefaultSignModes)))
}
