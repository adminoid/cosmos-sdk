package cli_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	rpcclientmock "github.com/tendermint/tendermint/rpc/client/mock"

	"github.com/adminoid/cosmos-sdk/client"
	"github.com/adminoid/cosmos-sdk/client/flags"
	"github.com/adminoid/cosmos-sdk/crypto/keyring"
	svrcmd "github.com/adminoid/cosmos-sdk/server/cmd"
	clitestutil "github.com/adminoid/cosmos-sdk/testutil/cli"
	sdk "github.com/adminoid/cosmos-sdk/types"
	testutilmod "github.com/adminoid/cosmos-sdk/types/module/testutil"
	"github.com/adminoid/cosmos-sdk/x/evidence"
	"github.com/adminoid/cosmos-sdk/x/evidence/client/cli"
)

func TestGetQueryCmd(t *testing.T) {
	encCfg := testutilmod.MakeTestEncodingConfig(evidence.AppModuleBasic{})
	kr := keyring.NewInMemory(encCfg.Codec)
	baseCtx := client.Context{}.
		WithKeyring(kr).
		WithTxConfig(encCfg.TxConfig).
		WithCodec(encCfg.Codec).
		WithClient(clitestutil.MockTendermintRPC{Client: rpcclientmock.Client{}}).
		WithAccountRetriever(client.MockAccountRetriever{}).
		WithOutput(io.Discard).
		WithChainID("test-chain")

	testCases := map[string]struct {
		args           []string
		ctxGen         func() client.Context
		expCmdOutput   string
		expectedOutput string
		expectErr      bool
	}{
		"non-existent evidence": {
			[]string{"DF0C23E8634E480F84B9D5674A7CDC9816466DEC28A3358F73260F68D28D7660"},
			func() client.Context {
				bz, _ := encCfg.Codec.Marshal(&sdk.TxResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return baseCtx.WithClient(c)
			},
			"DF0C23E8634E480F84B9D5674A7CDC9816466DEC28A3358F73260F68D28D7660",
			"",
			true,
		},
		"all evidence (default pagination)": {
			[]string{},
			func() client.Context {
				bz, _ := encCfg.Codec.Marshal(&sdk.TxResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return baseCtx.WithClient(c)
			},
			"",
			"evidence: []\npagination: null",
			false,
		},
		"all evidence (json output)": {
			[]string{
				fmt.Sprintf("--%s=json", flags.FlagOutput),
			},
			func() client.Context {
				bz, _ := encCfg.Codec.Marshal(&sdk.TxResponse{})
				c := clitestutil.NewMockTendermintRPC(abci.ResponseQuery{
					Value: bz,
				})
				return baseCtx.WithClient(c)
			},
			"",
			`{"evidence":[],"pagination":null}`,
			false,
		},
	}

	for name, tc := range testCases {
		tc := tc

		t.Run(name, func(t *testing.T) {
			var outBuf bytes.Buffer

			clientCtx := tc.ctxGen().WithOutput(&outBuf)
			ctx := svrcmd.CreateExecuteContext(context.Background())

			cmd := cli.GetQueryCmd()
			cmd.SetOut(io.Discard)
			require.NotNil(t, cmd)

			cmd.SetContext(ctx)
			cmd.SetArgs(tc.args)

			require.NoError(t, client.SetCmdClientContextHandler(clientCtx, cmd))

			if len(tc.args) != 0 {
				require.Contains(t, fmt.Sprint(cmd), tc.expCmdOutput)
			}

			out, err := clitestutil.ExecTestCLICmd(baseCtx, cmd, tc.args)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Contains(t, fmt.Sprint(cmd), "evidence [] [] Query for evidence by hash or for all (paginated) submitted evidence")
			require.Contains(t, strings.TrimSpace(out.String()), tc.expectedOutput)
		})
	}
}
