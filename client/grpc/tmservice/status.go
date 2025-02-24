package tmservice

import (
	"context"

	coretypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/adminoid/cosmos-sdk/client"
)

func getNodeStatus(ctx context.Context, clientCtx client.Context) (*coretypes.ResultStatus, error) {
	node, err := clientCtx.GetNode()
	if err != nil {
		return &coretypes.ResultStatus{}, err
	}
	return node.Status(ctx)
}
