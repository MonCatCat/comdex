package wasm

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"

	assetkeeper "github.com/MonCatCat/comdex/v9/x/asset/keeper"
	auctionKeeper "github.com/MonCatCat/comdex/v9/x/auction/keeper"
	collectorKeeper "github.com/MonCatCat/comdex/v9/x/collector/keeper"
	esmKeeper "github.com/MonCatCat/comdex/v9/x/esm/keeper"
	lendKeeper "github.com/MonCatCat/comdex/v9/x/lend/keeper"
	liquidationKeeper "github.com/MonCatCat/comdex/v9/x/liquidation/keeper"
	liquidityKeeper "github.com/MonCatCat/comdex/v9/x/liquidity/keeper"
	lockerkeeper "github.com/MonCatCat/comdex/v9/x/locker/keeper"
	rewardsKeeper "github.com/MonCatCat/comdex/v9/x/rewards/keeper"
	tokenMintkeeper "github.com/MonCatCat/comdex/v9/x/tokenmint/keeper"
	vaultKeeper "github.com/MonCatCat/comdex/v9/x/vault/keeper"
)

func RegisterCustomPlugins(
	locker *lockerkeeper.Keeper,
	tokenMint *tokenMintkeeper.Keeper,
	asset *assetkeeper.Keeper,
	rewards *rewardsKeeper.Keeper,
	collector *collectorKeeper.Keeper,
	liquidation *liquidationKeeper.Keeper,
	auction *auctionKeeper.Keeper,
	esm *esmKeeper.Keeper,
	vault *vaultKeeper.Keeper,
	lend *lendKeeper.Keeper,
	liquidity *liquidityKeeper.Keeper,
) []wasmkeeper.Option {
	comdexQueryPlugin := NewQueryPlugin(asset, locker, tokenMint, rewards, collector, liquidation, esm, vault, lend, liquidity)

	appDataQueryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(comdexQueryPlugin),
	})
	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(*locker, *rewards, *asset, *collector, *liquidation, *auction, *tokenMint, *esm, *vault, *liquidity),
	)

	return []wasm.Option{
		appDataQueryPluginOpt,
		messengerDecoratorOpt,
	}
}
