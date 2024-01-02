package tokenmint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/MonCatCat/comdex/v13/x/tokenmint/keeper"
	"github.com/MonCatCat/comdex/v13/x/tokenmint/types"
)

func InitGenesis(ctx sdk.Context, k keeper.Keeper, state *types.GenesisState) {
	k.SetParams(ctx, state.Params)

	for _, item := range state.TokenMint {
		k.SetTokenMint(ctx, item)
	}
}

func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.NewGenesisState(
		k.GetTotalTokenMinted(ctx),
		k.GetParams(ctx),
	)
}
