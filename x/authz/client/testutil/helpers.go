package authz

import (
	"github.com/adminoid/cosmos-sdk/client"
	"github.com/adminoid/cosmos-sdk/testutil"
	clitestutil "github.com/adminoid/cosmos-sdk/testutil/cli"
	"github.com/adminoid/cosmos-sdk/x/authz/client/cli"
)

func CreateGrant(clientCtx client.Context, args []string) (testutil.BufferWriter, error) {
	cmd := cli.NewCmdGrantAuthorization()
	return clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
}
