package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/MonCatCat/comdex/v13/x/asset/types"
)

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.params.SetParamSet(ctx, &params)
}

func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var params types.Params
	k.params.GetParamSet(ctx, &params)
	return params
}
