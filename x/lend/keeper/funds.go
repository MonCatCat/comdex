package keeper

import (
	"github.com/comdex-official/comdex/x/lend/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) GetReserveFunds(_ sdk.Context, pool types.Pool) sdk.Int {
	return sdk.NewInt(int64(pool.ReserveFunds))
}

func (k Keeper) UpdateReserveBalances(ctx sdk.Context, assetID uint64, moduleName string, payment sdk.Coin, inc bool) error {
	newAmount := payment.Amount.Quo(sdk.NewIntFromUint64(types.Uint64Two))
	buyBackStats, _ := k.GetBuyBackDepositStats(ctx)
	reserveStats, _ := k.GetReserveDepositStats(ctx)
	//var reserveBalanceStats []types.BalanceStats

	//var balanceStats []types.BalanceStats
	if inc {
		for _, v := range buyBackStats.BalanceStats {
			if v.AssetID == assetID {
				v.Amount = v.Amount.Add(newAmount)
				k.SetBuyBackDepositStats(ctx, buyBackStats)

			}
			//balanceStats = append(balanceStats, v)
			//newDepositStats := types.DepositStats{BalanceStats: balanceStats}

		}
		for _, v := range reserveStats.BalanceStats {
			if v.AssetID == assetID {
				v.Amount = v.Amount.Add(newAmount)
				k.SetReserveDepositStats(ctx, reserveStats)

			}
			//reserveBalanceStats = append(reserveBalanceStats, v)
			//newUserDepositStats := types.DepositStats{BalanceStats: reserveBalanceStats}
		}
		if err := k.bank.SendCoinsFromModuleToModule(ctx, moduleName, types.ModuleName, sdk.NewCoins(payment)); err != nil {
			return err
		}
	} else {
		for _, v := range buyBackStats.BalanceStats {
			if v.AssetID == assetID {
				v.Amount = v.Amount.Sub(newAmount)
				k.SetBuyBackDepositStats(ctx, buyBackStats)

			}
			//balanceStats = append(balanceStats, v)
			//newDepositStats := types.DepositStats{BalanceStats: balanceStats}
		}
		for _, v := range reserveStats.BalanceStats {
			if v.AssetID == assetID {
				v.Amount = v.Amount.Sub(newAmount)
				k.SetReserveDepositStats(ctx, reserveStats)

			}
			//reserveBalanceStats = append(reserveBalanceStats, v)
			//newUserDepositStats := types.DepositStats{BalanceStats: reserveBalanceStats}
		}
		if err := k.bank.SendCoinsFromModuleToModule(ctx, types.ModuleName, moduleName, sdk.NewCoins(payment)); err != nil {
			return err
		}
	}
	return nil
}

func (k Keeper) UpdateLendStats(ctx sdk.Context, AssetID, PoolID uint64, amount sdk.Int, inc bool) {
	assetStats, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, AssetID, PoolID)
	depositStats, _ := k.GetDepositStats(ctx)
	userDepositStats, _ := k.GetUserDepositStats(ctx)
	//var balanceStats []types.BalanceStats
	//var userBalanceStats []types.BalanceStats

	if inc {

		assetStats.TotalLend = assetStats.TotalLend.Add(amount)
		k.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)

		for _, v := range depositStats.BalanceStats {
			if v.AssetID == AssetID {
				v.Amount = v.Amount.Add(amount)
				k.SetDepositStats(ctx, depositStats)
				break
			}
			//balanceStats = append(balanceStats, v)
			//newDepositStats := types.DepositStats{BalanceStats: balanceStats}
		}
		for _, v := range userDepositStats.BalanceStats {
			if v.AssetID == AssetID {
				v.Amount = v.Amount.Add(amount)
				k.SetUserDepositStats(ctx, userDepositStats)
				break
			}
			//userBalanceStats = append(userBalanceStats, v)
			//newUserDepositStats := types.DepositStats{BalanceStats: userBalanceStats}
		}
	} else {
		assetStats.TotalLend = assetStats.TotalLend.Sub(amount)
		k.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)
		for _, v := range depositStats.BalanceStats {
			if v.AssetID == AssetID {
				v.Amount = v.Amount.Sub(amount)
				k.SetDepositStats(ctx, depositStats)
				break
			}
			//balanceStats = append(balanceStats, v)
			//newDepositStats := types.DepositStats{BalanceStats: balanceStats}
		}
		for _, v := range userDepositStats.BalanceStats {
			if v.AssetID == AssetID {
				v.Amount = v.Amount.Sub(amount)
				k.SetUserDepositStats(ctx, userDepositStats)
				break
			}
			//userBalanceStats = append(userBalanceStats, v)
			//newUserDepositStats := types.DepositStats{BalanceStats: userBalanceStats}
		}
	}
}

func (k Keeper) UpdateBorrowStats(ctx sdk.Context, pair types.Extended_Pair, borrowPos types.BorrowAsset, amount sdk.Int, inc bool) {
	if inc {
		borrowStats, _ := k.GetBorrowStats(ctx)
		//var userBalanceStats []types.BalanceStats
		for _, v := range borrowStats.BalanceStats {
			if v.AssetID == pair.AssetOut {
				v.Amount = v.Amount.Add(amount)
				k.SetBorrowStats(ctx, borrowStats)
				break
			}
			//userBalanceStats = append(userBalanceStats, v)
			//newUserDepositStats := types.DepositStats{BalanceStats: userBalanceStats}
		}

		assetStats, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, pair.AssetOut, pair.AssetOutPoolID)
		if borrowPos.IsStableBorrow {
			assetStats.TotalStableBorrowed = assetStats.TotalStableBorrowed.Add(amount)
			k.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)
		} else {
			assetStats.TotalBorrowed = assetStats.TotalBorrowed.Add(amount)
			k.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)
		}

	} else {
		borrowStats, _ := k.GetBorrowStats(ctx)
		//var userBalanceStats []types.BalanceStats
		for _, v := range borrowStats.BalanceStats {
			if v.AssetID == pair.AssetOut {
				v.Amount = v.Amount.Sub(amount)
				k.SetBorrowStats(ctx, borrowStats)
				break
			}
			//userBalanceStats = append(userBalanceStats, v)
			//newUserDepositStats := types.DepositStats{BalanceStats: userBalanceStats}
		}

		assetStats, _ := k.GetAssetStatsByPoolIDAndAssetID(ctx, pair.AssetOut, pair.AssetOutPoolID)
		if borrowPos.IsStableBorrow {
			assetStats.TotalStableBorrowed = assetStats.TotalStableBorrowed.Sub(amount)
			k.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)
		} else {
			assetStats.TotalBorrowed = assetStats.TotalBorrowed.Sub(amount)
			k.SetAssetStatsByPoolIDAndAssetID(ctx, assetStats)
		}
	}
}
