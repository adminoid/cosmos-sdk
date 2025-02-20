package params

import (
	"context"
	"encoding/json"

	govv1beta1 "github.com/adminoid/cosmos-sdk/x/gov/types/v1beta1"

	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	modulev1 "cosmossdk.io/api/cosmos/params/module/v1"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/depinject"

	"github.com/adminoid/cosmos-sdk/client"
	"github.com/adminoid/cosmos-sdk/codec"
	codectypes "github.com/adminoid/cosmos-sdk/codec/types"
	store "github.com/adminoid/cosmos-sdk/store/types"
	sdk "github.com/adminoid/cosmos-sdk/types"
	"github.com/adminoid/cosmos-sdk/types/module"
	simtypes "github.com/adminoid/cosmos-sdk/types/simulation"
	"github.com/adminoid/cosmos-sdk/x/params/client/cli"
	"github.com/adminoid/cosmos-sdk/x/params/keeper"
	"github.com/adminoid/cosmos-sdk/x/params/types"
	"github.com/adminoid/cosmos-sdk/x/params/types/proposal"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModule{}
)

// AppModuleBasic defines the basic application module used by the params module.
type AppModuleBasic struct{}

// Name returns the params module's name.
func (AppModuleBasic) Name() string {
	return proposal.ModuleName
}

// RegisterLegacyAminoCodec registers the params module's types on the given LegacyAmino codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	proposal.RegisterLegacyAminoCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the params
// module.
func (AppModuleBasic) DefaultGenesis(_ codec.JSONCodec) json.RawMessage { return nil }

// ValidateGenesis performs genesis state validation for the params module.
func (AppModuleBasic) ValidateGenesis(_ codec.JSONCodec, config client.TxEncodingConfig, _ json.RawMessage) error {
	return nil
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the params module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *gwruntime.ServeMux) {
	if err := proposal.RegisterQueryHandlerClient(context.Background(), mux, proposal.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// GetTxCmd returns no root tx command for the params module.
func (AppModuleBasic) GetTxCmd() *cobra.Command { return nil }

// GetQueryCmd returns no root query command for the params module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.NewQueryCmd()
}

func (am AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	proposal.RegisterInterfaces(registry)
}

// AppModule implements an application module for the distribution module.
type AppModule struct {
	AppModuleBasic

	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(k keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         k,
	}
}

var _ appmodule.AppModule = AppModule{}

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs a no-op.
func (am AppModule) InitGenesis(_ sdk.Context, _ codec.JSONCodec, _ json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// GenerateGenesisState performs a no-op.
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {}

// RegisterServices registers a gRPC query service to respond to the
// module-specific gRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	proposal.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

// ProposalContents returns all the params content functions used to
// simulate governance proposals.
func (am AppModule) ProposalContents(simState module.SimulationState) []simtypes.WeightedProposalContent {
	return nil
}

// RegisterStoreDecoder doesn't register any type.
func (AppModule) RegisterStoreDecoder(sdr sdk.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(_ module.SimulationState) []simtypes.WeightedOperation {
	return nil
}

// ExportGenesis performs a no-op.
func (am AppModule) ExportGenesis(_ sdk.Context, _ codec.JSONCodec) json.RawMessage {
	return nil
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

//
// App Wiring Setup
//

func init() {
	appmodule.Register(&modulev1.Module{},
		appmodule.Provide(
			ProvideModule,
			ProvideSubspace,
		))
}

//nolint:revive
type ParamsInputs struct {
	depinject.In

	KvStoreKey        *store.KVStoreKey
	TransientStoreKey *store.TransientStoreKey
	Cdc               codec.Codec
	LegacyAmino       *codec.LegacyAmino
}

//nolint:revive
type ParamsOutputs struct {
	depinject.Out

	ParamsKeeper keeper.Keeper
	Module       appmodule.AppModule
	GovHandler   govv1beta1.HandlerRoute
}

func ProvideModule(in ParamsInputs) ParamsOutputs {
	k := keeper.NewKeeper(in.Cdc, in.LegacyAmino, in.KvStoreKey, in.TransientStoreKey)

	m := NewAppModule(k)
	govHandler := govv1beta1.HandlerRoute{RouteKey: proposal.RouterKey, Handler: NewParamChangeProposalHandler(k)}

	return ParamsOutputs{ParamsKeeper: k, Module: m, GovHandler: govHandler}
}

type SubspaceInputs struct {
	depinject.In

	Key       depinject.ModuleKey
	Keeper    keeper.Keeper
	KeyTables map[string]types.KeyTable
}

func ProvideSubspace(in SubspaceInputs) types.Subspace {
	moduleName := in.Key.Name()
	kt, exists := in.KeyTables[moduleName]
	if !exists {
		return in.Keeper.Subspace(moduleName)
	}
	return in.Keeper.Subspace(moduleName).WithKeyTable(kt)
}
