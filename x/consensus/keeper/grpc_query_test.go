package keeper_test

import (
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/adminoid/cosmos-sdk/x/consensus/types"
)

func (s *KeeperTestSuite) TestGRPCQueryConsensusParams() {
	defaultConsensusParams := tmtypes.DefaultConsensusParams().ToProto()

	testCases := []struct {
		msg      string
		req      types.QueryParamsRequest
		malleate func()
		response types.QueryParamsResponse
		expPass  bool
	}{
		{
			"success",
			types.QueryParamsRequest{},
			func() {
				input := &types.MsgUpdateParams{
					Authority: s.consensusParamsKeeper.GetAuthority(),
					Block:     defaultConsensusParams.Block,
					Validator: defaultConsensusParams.Validator,
					Evidence:  defaultConsensusParams.Evidence,
				}
				s.msgServer.UpdateParams(s.ctx, input)
			},
			types.QueryParamsResponse{
				Params: &tmproto.ConsensusParams{
					Block:     defaultConsensusParams.Block,
					Validator: defaultConsensusParams.Validator,
					Evidence:  defaultConsensusParams.Evidence,
					Version:   defaultConsensusParams.Version,
				},
			},
			true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.msg, func() {
			s.SetupTest() // reset

			tc.malleate()
			res, err := s.queryClient.Params(s.ctx, &tc.req)

			if tc.expPass {
				s.Require().NoError(err)
				s.Require().NotNil(res)
				s.Require().Equal(tc.response.Params, res.Params)
			} else {
				s.Require().Error(err)
				s.Require().Nil(res)
			}
		})
	}
}
