package keeper

import (
	"github.com/MonCatCat/comdex/v13/x/auctionsV2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) HandleAuctionParamsProposal(ctx sdk.Context, p *types.DutchAutoBidParamsProposal) error {
	return k.AddAuctionParams(ctx, p.AuctionParams)
}
