package testutil

import (
	"cosmossdk.io/core/appconfig"
	_ "github.com/adminoid/cosmos-sdk/x/auth"           // import as blank for app wiring
	_ "github.com/adminoid/cosmos-sdk/x/auth/tx/config" // import as blank for app wiring
	_ "github.com/adminoid/cosmos-sdk/x/bank"           // import as blank for app wiring
	_ "github.com/adminoid/cosmos-sdk/x/consensus"      // import as blank for app wiring
	_ "github.com/adminoid/cosmos-sdk/x/genutil"        // import as blank for app wiring
	_ "github.com/adminoid/cosmos-sdk/x/mint"           // import as blank for app wiring
	_ "github.com/adminoid/cosmos-sdk/x/nft/module"     // import as blank for app wiring
	_ "github.com/adminoid/cosmos-sdk/x/params"         // import as blank for app wiring
	_ "github.com/adminoid/cosmos-sdk/x/staking"        // import as blank for app wiring

	authtypes "github.com/adminoid/cosmos-sdk/x/auth/types"
	banktypes "github.com/adminoid/cosmos-sdk/x/bank/types"
	consensustypes "github.com/adminoid/cosmos-sdk/x/consensus/types"
	genutiltypes "github.com/adminoid/cosmos-sdk/x/genutil/types"
	minttypes "github.com/adminoid/cosmos-sdk/x/mint/types"
	"github.com/adminoid/cosmos-sdk/x/nft"
	paramstypes "github.com/adminoid/cosmos-sdk/x/params/types"
	stakingtypes "github.com/adminoid/cosmos-sdk/x/staking/types"

	runtimev1alpha1 "cosmossdk.io/api/cosmos/app/runtime/v1alpha1"
	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	authmodulev1 "cosmossdk.io/api/cosmos/auth/module/v1"
	bankmodulev1 "cosmossdk.io/api/cosmos/bank/module/v1"
	consensusmodulev1 "cosmossdk.io/api/cosmos/consensus/module/v1"
	genutilmodulev1 "cosmossdk.io/api/cosmos/genutil/module/v1"
	mintmodulev1 "cosmossdk.io/api/cosmos/mint/module/v1"
	nftmodulev1 "cosmossdk.io/api/cosmos/nft/module/v1"
	paramsmodulev1 "cosmossdk.io/api/cosmos/params/module/v1"
	stakingmodulev1 "cosmossdk.io/api/cosmos/staking/module/v1"
	txconfigv1 "cosmossdk.io/api/cosmos/tx/config/v1"
)

var AppConfig = appconfig.Compose(&appv1alpha1.Config{
	Modules: []*appv1alpha1.ModuleConfig{
		{
			Name: "runtime",
			Config: appconfig.WrapAny(&runtimev1alpha1.Module{
				AppName: "NFTApp",
				BeginBlockers: []string{
					minttypes.ModuleName,
					stakingtypes.ModuleName,
					genutiltypes.ModuleName,
				},
				EndBlockers: []string{
					stakingtypes.ModuleName,
					genutiltypes.ModuleName,
				},
				InitGenesis: []string{
					authtypes.ModuleName,
					banktypes.ModuleName,
					stakingtypes.ModuleName,
					minttypes.ModuleName,
					genutiltypes.ModuleName,
					nft.ModuleName,
					paramstypes.ModuleName,
					consensustypes.ModuleName,
				},
			}),
		},
		{
			Name: authtypes.ModuleName,
			Config: appconfig.WrapAny(&authmodulev1.Module{
				Bech32Prefix: "cosmos",
				ModuleAccountPermissions: []*authmodulev1.ModuleAccountPermission{
					{Account: authtypes.FeeCollectorName},
					{Account: minttypes.ModuleName, Permissions: []string{authtypes.Minter}},
					{Account: stakingtypes.BondedPoolName, Permissions: []string{authtypes.Burner, stakingtypes.ModuleName}},
					{Account: stakingtypes.NotBondedPoolName, Permissions: []string{authtypes.Burner, stakingtypes.ModuleName}},
					{Account: nft.ModuleName},
				},
			}),
		},
		{
			Name:   banktypes.ModuleName,
			Config: appconfig.WrapAny(&bankmodulev1.Module{}),
		},
		{
			Name:   stakingtypes.ModuleName,
			Config: appconfig.WrapAny(&stakingmodulev1.Module{}),
		},
		{
			Name:   paramstypes.ModuleName,
			Config: appconfig.WrapAny(&paramsmodulev1.Module{}),
		},
		{
			Name:   consensustypes.ModuleName,
			Config: appconfig.WrapAny(&consensusmodulev1.Module{}),
		},
		{
			Name:   "tx",
			Config: appconfig.WrapAny(&txconfigv1.Config{}),
		},
		{
			Name:   genutiltypes.ModuleName,
			Config: appconfig.WrapAny(&genutilmodulev1.Module{}),
		},
		{
			Name:   minttypes.ModuleName,
			Config: appconfig.WrapAny(&mintmodulev1.Module{}),
		},
		{
			Name:   nft.ModuleName,
			Config: appconfig.WrapAny(&nftmodulev1.Module{}),
		},
	},
})
