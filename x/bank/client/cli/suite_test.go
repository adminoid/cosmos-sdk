package cli_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/suite"
	rpcclientmock "github.com/tendermint/tendermint/rpc/client/mock"

	"github.com/adminoid/cosmos-sdk/client"
	"github.com/adminoid/cosmos-sdk/crypto/keyring"
	clitestutil "github.com/adminoid/cosmos-sdk/testutil/cli"
	testutilmod "github.com/adminoid/cosmos-sdk/types/module/testutil"
	"github.com/adminoid/cosmos-sdk/x/bank"
)

type CLITestSuite struct {
	suite.Suite

	kr      keyring.Keyring
	encCfg  testutilmod.TestEncodingConfig
	baseCtx client.Context
}

func TestMigrateTestSuite(t *testing.T) {
	suite.Run(t, new(CLITestSuite))
}

func (s *CLITestSuite) SetupSuite() {
	s.encCfg = testutilmod.MakeTestEncodingConfig(bank.AppModuleBasic{})
	s.kr = keyring.NewInMemory(s.encCfg.Codec)
	s.baseCtx = client.Context{}.
		WithKeyring(s.kr).
		WithTxConfig(s.encCfg.TxConfig).
		WithCodec(s.encCfg.Codec).
		WithClient(clitestutil.MockTendermintRPC{Client: rpcclientmock.Client{}}).
		WithAccountRetriever(client.MockAccountRetriever{}).
		WithOutput(io.Discard)
}
