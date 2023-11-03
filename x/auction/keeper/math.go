package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Multiply(a, b sdk.Dec) sdk.Dec {
	return a.Mul(b)
}

func (k Keeper) getOutflowTokenInitialPrice(price sdk.Int, buffer sdk.Dec) sdk.Dec {
	result := buffer.Mul(sdk.NewDec(price.Int64()))
	return result
}

func (k Keeper) getOutflowTokenEndPrice(price, cusp sdk.Dec) sdk.Dec {
	result := Multiply(price, cusp)
	return result
}

func (k Keeper) getPriceFromLinearDecreaseFunction(top sdk.Dec, tau, dur sdk.Int) sdk.Dec {
	result1 := tau.Sub(dur)
	result2 := top.Mul(sdk.NewDec(result1.Int64()))
	result3 := result2.Quo(sdk.NewDec(tau.Int64()))
	return result3
}
