package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	assetkeeper "github.com/MonCatCat/comdex/v13/x/asset/keeper"
	"github.com/MonCatCat/comdex/v13/x/market/expected"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

type Keeper struct {
	cdc              codec.BinaryCodec
	key              storetypes.StoreKey
	params           paramstypes.Subspace
	scoped           expected.ScopedKeeper
	assetKeeper      assetkeeper.Keeper
	bandoraclekeeper expected.BandOracleKeeper
}

func NewKeeper(cdc codec.BinaryCodec, key storetypes.StoreKey, params paramstypes.Subspace, scoped expected.ScopedKeeper, assetKeeper assetkeeper.Keeper, bandoraclekeeper expected.BandOracleKeeper) Keeper {
	return Keeper{
		cdc:              cdc,
		key:              key,
		params:           params,
		scoped:           scoped,
		assetKeeper:      assetKeeper,
		bandoraclekeeper: bandoraclekeeper,
	}
}

func (k Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.key)
}
