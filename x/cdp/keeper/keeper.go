package keeper

import (
	"fmt"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/comdex-official/comdex/x/cdp/types"
)

type (
	Keeper struct {
		cdc           codec.BinaryMarshaler
		storeKey      sdk.StoreKey
		memKey        sdk.StoreKey
		accountKeeper types.AccountKeeper
		bankKeeper    types.BankKeeper
		paramSpace    types.ParamSubspace
	}
)

func NewKeeper(
	cdc codec.BinaryMarshaler,
	storeKey,
	memKey sdk.StoreKey,

) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
