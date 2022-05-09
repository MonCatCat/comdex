package bonding

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/comdex-official/comdex/x/bonding/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker is called on every block.
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
}

// Called every block to automatically unlock matured locks.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	// disable automatic withdraw before specific block height
	// it is actually for testing with legacy
	// MinBlockHeightToBeginAutoWithdrawing := int64(6)
	// if ctx.BlockHeight() < MinBlockHeightToBeginAutoWithdrawing {
	// 	return []abci.ValidatorUpdate{}
	// }

	// // delete synthetic locks matured before bonding deletion
	// k.DeleteAllMaturedSyntheticLocks(ctx)

	// // withdraw and delete locks
	// k.WithdrawAllMaturedLocks(ctx)
	// return []abci.ValidatorUpdate{}
}

// TODO: add invariant that no native bonding existent synthetic bonding exists by calling GetAllSyntheticBondings
// TODO: if superfluid does not delete synthetic bonding before native bonding deletion, it won't be able to be deleted
