package keeper_test

import (
	gocontext "context"
	"fmt"

	simtestutil "github.com/adminoid/cosmos-sdk/testutil/sims"
	sdk "github.com/adminoid/cosmos-sdk/types"
	"github.com/adminoid/cosmos-sdk/types/query"
	"github.com/adminoid/cosmos-sdk/x/staking/types"
)

func (suite *IntegrationTestSuite) TestGRPCQueryValidators() {
	queryClient, vals := suite.queryClient, suite.vals
	var req *types.QueryValidatorsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		numVals  int
		hasNext  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryValidatorsRequest{}
			},
			true,

			len(vals) + 1, // +1 validator from genesis state
			false,
		},
		{
			"empty status returns all the validators",
			func() {
				req = &types.QueryValidatorsRequest{Status: ""}
			},
			true,
			len(vals) + 1, // +1 validator from genesis state
			false,
		},
		{
			"invalid request",
			func() {
				req = &types.QueryValidatorsRequest{Status: "test"}
			},
			false,
			0,
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryValidatorsRequest{
					Status:     types.Bonded.String(),
					Pagination: &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			true,
			1,
			true,
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			valsResp, err := queryClient.Validators(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.NotNil(valsResp)
				suite.Equal(tc.numVals, len(valsResp.Validators))
				suite.Equal(uint64(len(vals))+1, valsResp.Pagination.Total) // +1 validator from genesis state

				if tc.hasNext {
					suite.NotNil(valsResp.Pagination.NextKey)
				} else {
					suite.Nil(valsResp.Pagination.NextKey)
				}
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGRPCQueryDelegatorValidators() {
	app, ctx, queryClient, addrs := suite.app, suite.ctx, suite.queryClient, suite.addrs
	params := app.StakingKeeper.GetParams(ctx)
	delValidators := app.StakingKeeper.GetDelegatorValidators(ctx, addrs[0], params.MaxValidators)
	var req *types.QueryDelegatorValidatorsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryDelegatorValidatorsRequest{}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryDelegatorValidatorsRequest{
					DelegatorAddr: addrs[0].String(),
					Pagination:    &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.DelegatorValidators(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.Equal(1, len(res.Validators))
				suite.NotNil(res.Pagination.NextKey)
				suite.Equal(uint64(len(delValidators)), res.Pagination.Total)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGRPCQueryDelegatorValidator() {
	queryClient, addrs, vals := suite.queryClient, suite.addrs, suite.vals
	addr := addrs[1]
	addrVal, addrVal1 := vals[0].OperatorAddress, vals[1].OperatorAddress
	var req *types.QueryDelegatorValidatorRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryDelegatorValidatorRequest{}
			},
			false,
		},
		{
			"invalid delegator, validator pair",
			func() {
				req = &types.QueryDelegatorValidatorRequest{
					DelegatorAddr: addr.String(),
					ValidatorAddr: addrVal,
				}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryDelegatorValidatorRequest{
					DelegatorAddr: addr.String(),
					ValidatorAddr: addrVal1,
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.DelegatorValidator(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.Equal(addrVal1, res.Validator.OperatorAddress)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGRPCQueryDelegation() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc, addrAcc1 := addrs[0], addrs[1]
	addrVal := vals[0].OperatorAddress
	valAddr, err := sdk.ValAddressFromBech32(addrVal)
	suite.NoError(err)
	delegation, found := app.StakingKeeper.GetDelegation(ctx, addrAcc, valAddr)
	suite.True(found)
	var req *types.QueryDelegationRequest

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryDelegationRequest{}
			},
			false,
		},
		{
			"invalid validator, delegator pair",
			func() {
				req = &types.QueryDelegationRequest{
					DelegatorAddr: addrAcc1.String(),
					ValidatorAddr: addrVal,
				}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryDelegationRequest{DelegatorAddr: addrAcc.String(), ValidatorAddr: addrVal}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Delegation(gocontext.Background(), req)
			if tc.expPass {
				suite.Equal(delegation.ValidatorAddress, res.DelegationResponse.Delegation.ValidatorAddress)
				suite.Equal(delegation.DelegatorAddress, res.DelegationResponse.Delegation.DelegatorAddress)
				suite.Equal(sdk.NewCoin(sdk.DefaultBondDenom, delegation.Shares.TruncateInt()), res.DelegationResponse.Balance)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGRPCQueryDelegatorDelegations() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc := addrs[0]
	addrVal1 := vals[0].OperatorAddress
	valAddr, err := sdk.ValAddressFromBech32(addrVal1)
	suite.NoError(err)
	delegation, found := app.StakingKeeper.GetDelegation(ctx, addrAcc, valAddr)
	suite.True(found)
	var req *types.QueryDelegatorDelegationsRequest

	testCases := []struct {
		msg       string
		malleate  func()
		onSuccess func(suite *IntegrationTestSuite, response *types.QueryDelegatorDelegationsResponse)
		expErr    bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryDelegatorDelegationsRequest{}
			},
			func(suite *IntegrationTestSuite, response *types.QueryDelegatorDelegationsResponse) {},
			true,
		},
		{
			"valid request with no delegations",
			func() {
				req = &types.QueryDelegatorDelegationsRequest{DelegatorAddr: addrs[4].String()}
			},
			func(suite *IntegrationTestSuite, response *types.QueryDelegatorDelegationsResponse) {
				suite.Equal(uint64(0), response.Pagination.Total)
				suite.Len(response.DelegationResponses, 0)
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryDelegatorDelegationsRequest{
					DelegatorAddr: addrAcc.String(),
					Pagination:    &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			func(suite *IntegrationTestSuite, response *types.QueryDelegatorDelegationsResponse) {
				suite.Equal(uint64(2), response.Pagination.Total)
				suite.Len(response.DelegationResponses, 1)
				suite.Equal(sdk.NewCoin(sdk.DefaultBondDenom, delegation.Shares.TruncateInt()), response.DelegationResponses[0].Balance)
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.DelegatorDelegations(gocontext.Background(), req)
			if tc.expErr {
				suite.Error(err)
			} else {
				suite.NoError(err)
				tc.onSuccess(suite, res)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGRPCQueryValidatorDelegations() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc := addrs[0]
	addrVal1 := vals[1].OperatorAddress
	valAddrs := simtestutil.ConvertAddrsToValAddrs(addrs)
	addrVal2 := valAddrs[4]
	valAddr, err := sdk.ValAddressFromBech32(addrVal1)
	suite.NoError(err)
	delegation, found := app.StakingKeeper.GetDelegation(ctx, addrAcc, valAddr)
	suite.True(found)

	var req *types.QueryValidatorDelegationsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		expErr   bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryValidatorDelegationsRequest{}
			},
			false,
			true,
		},
		{
			"invalid validator delegator pair",
			func() {
				req = &types.QueryValidatorDelegationsRequest{ValidatorAddr: addrVal2.String()}
			},
			false,
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryValidatorDelegationsRequest{
					ValidatorAddr: addrVal1,
					Pagination:    &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			true,
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.ValidatorDelegations(gocontext.Background(), req)
			if tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.Len(res.DelegationResponses, 1)
				suite.NotNil(res.Pagination.NextKey)
				suite.Equal(uint64(2), res.Pagination.Total)
				suite.Equal(addrVal1, res.DelegationResponses[0].Delegation.ValidatorAddress)
				suite.Equal(sdk.NewCoin(sdk.DefaultBondDenom, delegation.Shares.TruncateInt()), res.DelegationResponses[0].Balance)
			} else if !tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.Nil(res.DelegationResponses)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGRPCQueryUnbondingDelegation() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc2 := addrs[1]
	addrVal2 := vals[1].OperatorAddress

	unbondingTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 2)
	valAddr, err1 := sdk.ValAddressFromBech32(addrVal2)
	suite.NoError(err1)
	_, err := app.StakingKeeper.Undelegate(ctx, addrAcc2, valAddr, sdk.NewDecFromInt(unbondingTokens))
	suite.NoError(err)

	unbond, found := app.StakingKeeper.GetUnbondingDelegation(ctx, addrAcc2, valAddr)
	suite.True(found)
	var req *types.QueryUnbondingDelegationRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryUnbondingDelegationRequest{}
			},
			false,
		},
		{
			"invalid request",
			func() {
				req = &types.QueryUnbondingDelegationRequest{}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryUnbondingDelegationRequest{
					DelegatorAddr: addrAcc2.String(), ValidatorAddr: addrVal2,
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.UnbondingDelegation(gocontext.Background(), req)
			if tc.expPass {
				suite.NotNil(res)
				suite.Equal(unbond, res.Unbond)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGRPCQueryDelegatorUnbondingDelegations() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc, addrAcc1 := addrs[0], addrs[1]
	addrVal, addrVal2 := vals[0].OperatorAddress, vals[1].OperatorAddress

	unbondingTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 2)
	valAddr1, err1 := sdk.ValAddressFromBech32(addrVal)
	suite.NoError(err1)
	_, err := app.StakingKeeper.Undelegate(ctx, addrAcc, valAddr1, sdk.NewDecFromInt(unbondingTokens))
	suite.NoError(err)
	valAddr2, err1 := sdk.ValAddressFromBech32(addrVal2)
	suite.NoError(err1)
	_, err = app.StakingKeeper.Undelegate(ctx, addrAcc, valAddr2, sdk.NewDecFromInt(unbondingTokens))
	suite.NoError(err)

	unbond, found := app.StakingKeeper.GetUnbondingDelegation(ctx, addrAcc, valAddr1)
	suite.True(found)
	var req *types.QueryDelegatorUnbondingDelegationsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		expErr   bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryDelegatorUnbondingDelegationsRequest{}
			},
			false,
			true,
		},
		{
			"invalid request",
			func() {
				req = &types.QueryDelegatorUnbondingDelegationsRequest{DelegatorAddr: addrAcc1.String()}
			},
			false,
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryDelegatorUnbondingDelegationsRequest{
					DelegatorAddr: addrAcc.String(),
					Pagination:    &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			true,
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.DelegatorUnbondingDelegations(gocontext.Background(), req)
			if tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.NotNil(res.Pagination.NextKey)
				suite.Equal(uint64(2), res.Pagination.Total)
				suite.Len(res.UnbondingResponses, 1)
				suite.Equal(unbond, res.UnbondingResponses[0])
			} else if !tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.Nil(res.UnbondingResponses)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGRPCQueryPoolParameters() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	bondDenom := sdk.DefaultBondDenom

	// Query pool
	res, err := queryClient.Pool(gocontext.Background(), &types.QueryPoolRequest{})
	suite.NoError(err)
	bondedPool := app.StakingKeeper.GetBondedPool(ctx)
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
	suite.Equal(app.BankKeeper.GetBalance(ctx, notBondedPool.GetAddress(), bondDenom).Amount, res.Pool.NotBondedTokens)
	suite.Equal(app.BankKeeper.GetBalance(ctx, bondedPool.GetAddress(), bondDenom).Amount, res.Pool.BondedTokens)

	// Query Params
	resp, err := queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
	suite.NoError(err)
	suite.Equal(app.StakingKeeper.GetParams(ctx), resp.Params)
}

func (suite *IntegrationTestSuite) TestGRPCQueryHistoricalInfo() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	hi, found := app.StakingKeeper.GetHistoricalInfo(ctx, 5)
	suite.True(found)

	var req *types.QueryHistoricalInfoRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryHistoricalInfoRequest{}
			},
			false,
		},
		{
			"invalid request with negative height",
			func() {
				req = &types.QueryHistoricalInfoRequest{Height: -1}
			},
			false,
		},
		{
			"valid request with old height",
			func() {
				req = &types.QueryHistoricalInfoRequest{Height: 4}
			},
			false,
		},
		{
			"valid request with current height",
			func() {
				req = &types.QueryHistoricalInfoRequest{Height: 5}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.HistoricalInfo(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.NotNil(res)
				suite.True(hi.Equal(res.Hist))
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGRPCQueryRedelegations() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals

	addrAcc, addrAcc1 := addrs[0], addrs[1]
	valAddrs := simtestutil.ConvertAddrsToValAddrs(addrs)
	val1, val2, val3, val4 := vals[0], vals[1], valAddrs[3], valAddrs[4]
	delAmount := app.StakingKeeper.TokensFromConsensusPower(ctx, 1)
	_, err := app.StakingKeeper.Delegate(ctx, addrAcc1, delAmount, types.Unbonded, val1, true)
	suite.NoError(err)
	applyValidatorSetUpdates(suite.T(), ctx, app.StakingKeeper, -1)

	rdAmount := app.StakingKeeper.TokensFromConsensusPower(ctx, 1)
	_, err = app.StakingKeeper.BeginRedelegation(ctx, addrAcc1, val1.GetOperator(), val2.GetOperator(), sdk.NewDecFromInt(rdAmount))
	suite.NoError(err)
	applyValidatorSetUpdates(suite.T(), ctx, app.StakingKeeper, -1)

	redel, found := app.StakingKeeper.GetRedelegation(ctx, addrAcc1, val1.GetOperator(), val2.GetOperator())
	suite.True(found)

	var req *types.QueryRedelegationsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		expErr   bool
	}{
		{
			"request redelegations for non existent addr",
			func() {
				req = &types.QueryRedelegationsRequest{DelegatorAddr: addrAcc.String()}
			},
			false,
			false,
		},
		{
			"request redelegations with non existent pairs",
			func() {
				req = &types.QueryRedelegationsRequest{
					DelegatorAddr: addrAcc.String(), SrcValidatorAddr: val3.String(),
					DstValidatorAddr: val4.String(),
				}
			},
			false,
			true,
		},
		{
			"request redelegations with delegatoraddr, sourceValAddr, destValAddr",
			func() {
				req = &types.QueryRedelegationsRequest{
					DelegatorAddr: addrAcc1.String(), SrcValidatorAddr: val1.OperatorAddress,
					DstValidatorAddr: val2.OperatorAddress, Pagination: &query.PageRequest{},
				}
			},
			true,
			false,
		},
		{
			"request redelegations with delegatoraddr and sourceValAddr",
			func() {
				req = &types.QueryRedelegationsRequest{
					DelegatorAddr: addrAcc1.String(), SrcValidatorAddr: val1.OperatorAddress,
					Pagination: &query.PageRequest{},
				}
			},
			true,
			false,
		},
		{
			"query redelegations with sourceValAddr only",
			func() {
				req = &types.QueryRedelegationsRequest{
					SrcValidatorAddr: val1.GetOperator().String(),
					Pagination:       &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			true,
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Redelegations(gocontext.Background(), req)
			if tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.Len(res.RedelegationResponses, len(redel.Entries))
				suite.Equal(redel.DelegatorAddress, res.RedelegationResponses[0].Redelegation.DelegatorAddress)
				suite.Equal(redel.ValidatorSrcAddress, res.RedelegationResponses[0].Redelegation.ValidatorSrcAddress)
				suite.Equal(redel.ValidatorDstAddress, res.RedelegationResponses[0].Redelegation.ValidatorDstAddress)
				suite.Len(redel.Entries, len(res.RedelegationResponses[0].Entries))
			} else if !tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.Nil(res.RedelegationResponses)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *IntegrationTestSuite) TestGRPCQueryValidatorUnbondingDelegations() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc1, _ := addrs[0], addrs[1]
	val1 := vals[0]

	// undelegate
	undelAmount := app.StakingKeeper.TokensFromConsensusPower(ctx, 2)
	_, err := app.StakingKeeper.Undelegate(ctx, addrAcc1, val1.GetOperator(), sdk.NewDecFromInt(undelAmount))
	suite.NoError(err)
	applyValidatorSetUpdates(suite.T(), ctx, app.StakingKeeper, -1)

	var req *types.QueryValidatorUnbondingDelegationsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryValidatorUnbondingDelegationsRequest{}
			},
			false,
		},
		{
			"valid request",
			func() {
				req = &types.QueryValidatorUnbondingDelegationsRequest{
					ValidatorAddr: val1.GetOperator().String(),
					Pagination:    &query.PageRequest{Limit: 1, CountTotal: true},
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.ValidatorUnbondingDelegations(gocontext.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.Equal(uint64(1), res.Pagination.Total)
				suite.Equal(1, len(res.UnbondingResponses))
				suite.Equal(res.UnbondingResponses[0].ValidatorAddress, val1.OperatorAddress)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}
