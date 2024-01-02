package lend

import (
	utils "github.com/MonCatCat/comdex/v13/types"
	"github.com/MonCatCat/comdex/v13/x/lend/keeper"
	"github.com/MonCatCat/comdex/v13/x/lend/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock, k keeper.Keeper) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, ctx.BlockTime(), telemetry.MetricKeyBeginBlocker)

	_ = utils.ApplyFuncIfNoError(ctx, func(ctx sdk.Context) error {
		if ctx.BlockHeight()%14400 == 0 {
			err := k.DeletePoolAndTransferInterest(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
