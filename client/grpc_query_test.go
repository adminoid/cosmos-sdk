package client_test

import (
	"context"
	"testing"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmjson "github.com/tendermint/tendermint/libs/json"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"cosmossdk.io/depinject"
	"github.com/adminoid/cosmos-sdk/baseapp"
	"github.com/adminoid/cosmos-sdk/codec"
	codectypes "github.com/adminoid/cosmos-sdk/codec/types"
	"github.com/adminoid/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/adminoid/cosmos-sdk/runtime"
	"github.com/adminoid/cosmos-sdk/testutil/sims"
	"github.com/adminoid/cosmos-sdk/testutil/testdata"
	sdk "github.com/adminoid/cosmos-sdk/types"
	"github.com/adminoid/cosmos-sdk/x/auth/testutil"
	authtypes "github.com/adminoid/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/adminoid/cosmos-sdk/x/bank/keeper"
	"github.com/adminoid/cosmos-sdk/x/bank/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	ctx                   sdk.Context
	genesisAccount        *authtypes.BaseAccount
	bankClient            types.QueryClient
	testClient            testdata.QueryClient
	genesisAccountBalance int64
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")
	var (
		interfaceRegistry codectypes.InterfaceRegistry
		bankKeeper        bankkeeper.BaseKeeper
		appBuilder        *runtime.AppBuilder
		cdc               codec.Codec
	)

	// TODO duplicated from testutils/sims/app_helpers.go
	// need more composable startup options for simapp, this test needed a handle to the closed over genesis account
	// to query balances
	err := depinject.Inject(testutil.AppConfig, &interfaceRegistry, &bankKeeper, &appBuilder, &cdc)
	s.NoError(err)

	app := appBuilder.Build(log.NewNopLogger(), dbm.NewMemDB(), nil)
	err = app.Load(true)
	s.NoError(err)

	valSet, err := sims.CreateRandomValidatorSet()
	s.NoError(err)

	// generate genesis account
	s.genesisAccountBalance = 100000000000000
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := types.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(s.genesisAccountBalance))),
	}

	genesisState, err := sims.GenesisStateWithValSet(cdc, app.DefaultGenesis(), valSet, []authtypes.GenesisAccount{acc}, balance)
	s.NoError(err)

	stateBytes, err := tmjson.MarshalIndent(genesisState, "", " ")
	s.NoError(err)

	// init chain will set the validator set and initialize the genesis accounts
	app.InitChain(
		abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: sims.DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	app.Commit()
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
		Height:             app.LastBlockHeight() + 1,
		AppHash:            app.LastCommitID().Hash,
		ValidatorsHash:     valSet.Hash(),
		NextValidatorsHash: valSet.Hash(),
	}})

	// end of app init

	s.ctx = app.BaseApp.NewContext(false, tmproto.Header{})
	queryHelper := baseapp.NewQueryServerTestHelper(s.ctx, interfaceRegistry)
	types.RegisterQueryServer(queryHelper, bankKeeper)
	testdata.RegisterQueryServer(queryHelper, testdata.QueryImpl{})
	s.bankClient = types.NewQueryClient(queryHelper)
	s.testClient = testdata.NewQueryClient(queryHelper)
	s.genesisAccount = acc
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
}

func (s *IntegrationTestSuite) TestGRPCQuery() {
	denom := sdk.DefaultBondDenom

	// gRPC query to test service should work
	testRes, err := s.testClient.Echo(context.Background(), &testdata.EchoRequest{Message: "hello"})
	s.Require().NoError(err)
	s.Require().Equal("hello", testRes.Message)

	// gRPC query to bank service should work
	var header metadata.MD
	res, err := s.bankClient.Balance(
		context.Background(),
		&types.QueryBalanceRequest{Address: s.genesisAccount.GetAddress().String(), Denom: denom},
		grpc.Header(&header), // Also fetch grpc header
	)
	s.Require().NoError(err)
	bal := res.GetBalance()
	s.Equal(sdk.NewCoin(denom, sdk.NewInt(s.genesisAccountBalance)), *bal)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
