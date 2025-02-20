package v2

import (
	"github.com/adminoid/cosmos-sdk/codec"
	"github.com/adminoid/cosmos-sdk/store/prefix"
	storetypes "github.com/adminoid/cosmos-sdk/store/types"
	"github.com/adminoid/cosmos-sdk/types"
	"github.com/adminoid/cosmos-sdk/x/feegrant"
)

func addAllowancesByExpTimeQueue(ctx types.Context, store storetypes.KVStore, cdc codec.BinaryCodec) error {
	prefixStore := prefix.NewStore(store, FeeAllowanceKeyPrefix)
	iterator := prefixStore.Iterator(nil, nil)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var grant feegrant.Grant
		bz := iterator.Value()
		if err := cdc.Unmarshal(bz, &grant); err != nil {
			return err
		}

		grantInfo, err := grant.GetGrant()
		if err != nil {
			return err
		}

		exp, err := grantInfo.ExpiresAt()
		if err != nil {
			return err
		}

		if exp != nil {
			// store key is not changed in 0.46
			key := iterator.Key()
			if exp.Before(ctx.BlockTime()) {
				prefixStore.Delete(key)
			} else {
				grantByExpTimeQueueKey := FeeAllowancePrefixQueue(exp, key)
				store.Set(grantByExpTimeQueueKey, []byte{})
			}
		}
	}

	return nil
}

func MigrateStore(ctx types.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	return addAllowancesByExpTimeQueue(ctx, store, cdc)
}
