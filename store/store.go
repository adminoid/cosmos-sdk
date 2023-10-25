package store

import (
	dbm "github.com/cosmos/cosmos-db"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/adminoid/cosmos-sdk/store/cache"
	"github.com/adminoid/cosmos-sdk/store/metrics"
	"github.com/adminoid/cosmos-sdk/store/rootmulti"
	"github.com/adminoid/cosmos-sdk/store/types"
)

func NewCommitMultiStore(db dbm.DB, logger log.Logger, metricGatherer metrics.StoreMetrics) types.CommitMultiStore {
	return rootmulti.NewStore(db, logger, metricGatherer)
}

func NewCommitKVStoreCacheManager() types.MultiStorePersistentCache {
	return cache.NewCommitKVStoreCacheManager(cache.DefaultCommitKVStoreCacheSize)
}
