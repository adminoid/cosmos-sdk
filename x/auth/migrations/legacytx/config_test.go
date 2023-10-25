package legacytx_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/adminoid/cosmos-sdk/codec"
	cryptoAmino "github.com/adminoid/cosmos-sdk/crypto/codec"
	"github.com/adminoid/cosmos-sdk/testutil/testdata"
	sdk "github.com/adminoid/cosmos-sdk/types"
	"github.com/adminoid/cosmos-sdk/x/auth/migrations/legacytx"
	txtestutil "github.com/adminoid/cosmos-sdk/x/auth/tx/testutil"
)

func testCodec() *codec.LegacyAmino {
	cdc := codec.NewLegacyAmino()
	sdk.RegisterLegacyAminoCodec(cdc)
	cryptoAmino.RegisterCrypto(cdc)
	cdc.RegisterConcrete(&testdata.TestMsg{}, "cosmos-sdk/Test", nil)
	return cdc
}

func TestStdTxConfig(t *testing.T) {
	cdc := testCodec()
	txGen := legacytx.StdTxConfig{Cdc: cdc}
	suite.Run(t, txtestutil.NewTxConfigTestSuite(txGen))
}
