package market

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/MonCatCat/comdex/v13/x/market/keeper"
	"github.com/MonCatCat/comdex/v13/x/market/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(_ sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		default:
			return nil, errors.Wrapf(types.ErrorUnknownMsgType, "%T", msg)
		}
	}
}
