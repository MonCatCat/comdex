package keeper

import (
	auctiontypes "github.com/comdex-official/comdex/x/auction/types"
	lendtypes "github.com/comdex-official/comdex/x/lend/types"
	"github.com/comdex-official/comdex/x/liquidation/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) LiquidateBorrows(ctx sdk.Context) error {
	borrows, found := k.lend.GetBorrows(ctx)
	params := k.GetParams(ctx)
	if !found {
		return nil
	}
	liquidationOffsetHolder, found := k.GetLiquidationOffsetHolder(ctx, lendtypes.AppID, types.VaultLiquidationsOffsetPrefix)
	if !found {
		liquidationOffsetHolder = types.NewLiquidationOffsetHolder(lendtypes.AppID, 0)
	}
	borrowIDs := borrows.BorrowIDs
	start, end := types.GetSliceStartEndForLiquidations(len(borrowIDs), int(liquidationOffsetHolder.CurrentOffset), int(params.LiquidationBatchSize))

	if start == end {
		liquidationOffsetHolder.CurrentOffset = 0
		start, end = types.GetSliceStartEndForLiquidations(len(borrowIDs), int(liquidationOffsetHolder.CurrentOffset), int(params.LiquidationBatchSize))
	}
	newBorrowIDs := borrowIDs[start:end]
	for l := range newBorrowIDs {
		borrowPos, found := k.lend.GetBorrow(ctx, newBorrowIDs[l])
		if !found {
			continue
		}
		lendPair, _ := k.lend.GetLendPair(ctx, borrowPos.PairID)
		lendPos, _ := k.lend.GetLend(ctx, borrowPos.LendingID)
		killSwitchParams, _ := k.esm.GetKillSwitchData(ctx, lendPos.AppID)
		if killSwitchParams.BreakerEnable {
			continue
		}
		pool, _ := k.lend.GetPool(ctx, lendPos.PoolID)
		assetIn, _ := k.asset.GetAsset(ctx, lendPair.AssetIn)
		assetOut, _ := k.asset.GetAsset(ctx, lendPair.AssetOut)
		var currentCollateralizationRatio sdk.Dec

		liqThreshold, _ := k.lend.GetAssetRatesStats(ctx, lendPair.AssetIn)
		liqThresholdBridgedAssetOne, _ := k.lend.GetAssetRatesStats(ctx, pool.FirstBridgedAssetID)
		liqThresholdBridgedAssetTwo, _ := k.lend.GetAssetRatesStats(ctx, pool.SecondBridgedAssetID)
		firstBridgedAsset, _ := k.asset.GetAsset(ctx, pool.FirstBridgedAssetID)
		if borrowPos.BridgedAssetAmount.Amount.Equal(sdk.ZeroInt()) {
			currentCollateralizationRatio, _ = k.lend.CalculateCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.UpdatedAmountOut, assetOut)
			if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold) {
				err := k.CreateLockedBorrow(ctx, borrowPos, currentCollateralizationRatio, lendPos.AppID)
				if err != nil {
					continue
				}
				k.lend.UpdateBorrowStats(ctx, lendPair, borrowPos, borrowPos.AmountOut.Amount, false)
				k.lend.DeleteBorrow(ctx, borrowIDs[l])
				err = k.lend.UpdateUserBorrowIDMapping(ctx, lendPos.Owner, borrowIDs[l], false)
				if err != nil {
					continue
				}
				err = k.lend.UpdateBorrowIDByOwnerAndPoolMapping(ctx, lendPos.Owner, borrowIDs[l], lendPair.AssetOutPoolID, false)
				if err != nil {
					continue
				}
				err = k.lend.UpdateBorrowIdsMapping(ctx, borrowIDs[l], false)
				if err != nil {
					continue
				}

			}

		} else {
			if borrowPos.BridgedAssetAmount.Denom == firstBridgedAsset.Denom {
				currentCollateralizationRatio, _ = k.lend.CalculateCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.UpdatedAmountOut, assetOut)
				if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold.Mul(liqThresholdBridgedAssetOne.LiquidationThreshold)) {
					err := k.CreateLockedBorrow(ctx, borrowPos, currentCollateralizationRatio, lendPos.AppID)
					if err != nil {
						continue
					}
					k.lend.UpdateBorrowStats(ctx, lendPair, borrowPos, borrowPos.AmountOut.Amount, false)
					k.lend.DeleteBorrow(ctx, borrowIDs[l])
					err = k.lend.UpdateUserBorrowIDMapping(ctx, lendPos.Owner, borrowIDs[l], false)
					if err != nil {
						continue
					}
					err = k.lend.UpdateBorrowIDByOwnerAndPoolMapping(ctx, lendPos.Owner, borrowIDs[l], lendPair.AssetOutPoolID, false)
					if err != nil {
						continue
					}
					err = k.lend.UpdateBorrowIdsMapping(ctx, borrowIDs[l], false)
					if err != nil {
						continue
					}

				}
			} else {
				currentCollateralizationRatio, _ = k.lend.CalculateCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.UpdatedAmountOut, assetOut)

				if sdk.Dec.GT(currentCollateralizationRatio, liqThreshold.LiquidationThreshold.Mul(liqThresholdBridgedAssetTwo.LiquidationThreshold)) {
					err := k.CreateLockedBorrow(ctx, borrowPos, currentCollateralizationRatio, lendPos.AppID)
					if err != nil {
						continue
					}
					k.lend.UpdateBorrowStats(ctx, lendPair, borrowPos, borrowPos.AmountOut.Amount, false)
					k.lend.DeleteBorrow(ctx, borrowIDs[l])
					err = k.lend.UpdateUserBorrowIDMapping(ctx, lendPos.Owner, borrowIDs[l], false)
					if err != nil {
						continue
					}
					err = k.lend.UpdateBorrowIDByOwnerAndPoolMapping(ctx, lendPos.Owner, borrowIDs[l], lendPair.AssetOutPoolID, false)
					if err != nil {
						continue
					}
					err = k.lend.UpdateBorrowIdsMapping(ctx, borrowIDs[l], false)
					if err != nil {
						continue
					}

				}
			}
		}
		lockedVaultID := k.GetLockedVaultID(ctx)
		err := k.UpdateLockedBorrows(ctx, lendPos.AppID, lockedVaultID)
		if err != nil {
			return nil
		}
	}
	liquidationOffsetHolder.CurrentOffset = uint64(end)
	k.SetLiquidationOffsetHolder(ctx, types.VaultLiquidationsOffsetPrefix, liquidationOffsetHolder)

	return nil
}

func (k Keeper) CreateLockedBorrow(ctx sdk.Context, borrow lendtypes.BorrowAsset, collateralizationRatio sdk.Dec, appID uint64) error {
	lockedVaultID := k.GetLockedVaultID(ctx)
	lendPos, _ := k.lend.GetLend(ctx, borrow.LendingID)
	kind := &types.LockedVault_BorrowMetaData{
		BorrowMetaData: &types.BorrowMetaData{
			LendingId:          borrow.LendingID,
			IsStableBorrow:     borrow.IsStableBorrow,
			StableBorrowRate:   borrow.StableBorrowRate,
			BridgedAssetAmount: borrow.BridgedAssetAmount,
		},
	}
	var value = types.LockedVault{
		LockedVaultId:                lockedVaultID + 1,
		AppId:                        appID,
		OriginalVaultId:              borrow.ID,
		ExtendedPairId:               borrow.PairID,
		Owner:                        lendPos.Owner,
		AmountIn:                     borrow.AmountIn.Amount,
		AmountOut:                    borrow.AmountOut.Amount,
		UpdatedAmountOut:             borrow.AmountOut.Amount.Add(borrow.Interest_Accumulated),
		Initiator:                    types.ModuleName,
		IsAuctionComplete:            false,
		IsAuctionInProgress:          false,
		CrAtLiquidation:              collateralizationRatio,
		CurrentCollaterlisationRatio: collateralizationRatio,
		CollateralToBeAuctioned:      sdk.ZeroDec(),
		LiquidationTimestamp:         ctx.BlockTime(),
		SellOffHistory:               nil,
		InterestAccumulated:          borrow.Interest_Accumulated,
		Kind:                         kind,
	}
	k.SetLockedVault(ctx, value)
	k.SetLockedVaultID(ctx, value.LockedVaultId)
	return nil
}

func (k Keeper) UpdateLockedBorrows(ctx sdk.Context, appID, id uint64) error {
	lockedVault, _ := k.GetLockedVault(ctx, appID, id)
	pair, _ := k.lend.GetLendPair(ctx, lockedVault.ExtendedPairId)
	borrowMetaData := lockedVault.GetBorrowMetaData()
	if borrowMetaData != nil {
		lendPos, _ := k.lend.GetLend(ctx, borrowMetaData.LendingId)
		pool, _ := k.lend.GetPool(ctx, lendPos.PoolID)
		var unliquidatePointPercentage sdk.Dec
		firstBridgeAsset, _ := k.asset.GetAsset(ctx, pool.FirstBridgedAssetID)
		firstBridgeAssetStats, _ := k.lend.GetAssetRatesStats(ctx, pool.FirstBridgedAssetID)
		secondBridgeAssetStats, _ := k.lend.GetAssetRatesStats(ctx, pool.SecondBridgedAssetID)
		liqThreshold, _ := k.lend.GetAssetRatesStats(ctx, pair.AssetIn)

		if !borrowMetaData.BridgedAssetAmount.Amount.Equal(sdk.ZeroInt()) {
			if borrowMetaData.BridgedAssetAmount.Denom == firstBridgeAsset.Denom {
				unliquidatePointPercentage = liqThreshold.LiquidationThreshold.Mul(firstBridgeAssetStats.LiquidationThreshold)
			} else {
				unliquidatePointPercentage = liqThreshold.LiquidationThreshold.Mul(secondBridgeAssetStats.LiquidationThreshold)
			}
		} else {
			unliquidatePointPercentage = liqThreshold.LiquidationThreshold
		}

		assetRatesStats, _ := k.lend.GetAssetRatesStats(ctx, pair.AssetIn)

		if (!lockedVault.IsAuctionInProgress && !lockedVault.IsAuctionComplete) || (lockedVault.IsAuctionComplete && lockedVault.CurrentCollaterlisationRatio.GTE(unliquidatePointPercentage)) {

			assetIn, _ := k.asset.GetAsset(ctx, pair.AssetIn)

			assetOut, _ := k.asset.GetAsset(ctx, pair.AssetOut)

			collateralizationRatio, err := k.lend.CalculateCollateralizationRatio(ctx, lockedVault.AmountIn, assetIn, lockedVault.UpdatedAmountOut, assetOut)
			if err != nil {
				return nil
			}

			assetInTotal, _ := k.market.CalcAssetPrice(ctx, assetIn.Id, lockedVault.AmountIn)
			assetOutTotal, _ := k.market.CalcAssetPrice(ctx, assetOut.Id, lockedVault.AmountOut)
			deductionPercentage, _ := sdk.NewDecFromStr("1.0")

			var c sdk.Dec
			if !borrowMetaData.BridgedAssetAmount.Amount.Equal(sdk.ZeroInt()) {
				if borrowMetaData.BridgedAssetAmount.Denom == firstBridgeAsset.Denom {
					c = assetRatesStats.LiquidationThreshold.Mul(firstBridgeAssetStats.Ltv)
				} else {
					c = assetRatesStats.LiquidationThreshold.Mul(secondBridgeAssetStats.Ltv)
				}

			} else {
				c = assetRatesStats.LiquidationThreshold
			}
			penalty := assetRatesStats.LiquidationPenalty.Add(assetRatesStats.LiquidationBonus)
			b := deductionPercentage.Add(penalty)
			totalIn := assetInTotal
			totalOut := assetOutTotal
			factor1 := c.Mul(totalIn)
			factor2 := b.Mul(c)
			numerator := totalOut.Sub(factor1)
			denominator := deductionPercentage.Sub(factor2)
			selloffAmount := numerator.Quo(denominator)
			updatedLockedVault := lockedVault
			if lockedVault.SellOffHistory == nil {
				twaData, _ := k.market.GetTwa(ctx, assetIn.Id)
				aip := sdk.NewDec(int64(twaData.Twa))
				liquidationDeductionAmt := selloffAmount.Mul(penalty)
				liquidationDeductionAmount := liquidationDeductionAmt.Quo(aip)
				bonusToBidderAmt := selloffAmount.Mul(assetRatesStats.LiquidationBonus)

				bonusToBidderAmount := bonusToBidderAmt.Quo(aip)
				penaltyToReserveAmt := selloffAmount.Mul(assetRatesStats.LiquidationPenalty)
				penaltyToReserveAmount := penaltyToReserveAmt.Quo(aip)
				err = k.bank.SendCoinsFromModuleToModule(ctx, pool.ModuleName, auctiontypes.ModuleName, sdk.NewCoins(sdk.NewCoin(assetIn.Denom, sdk.NewInt(bonusToBidderAmount.TruncateInt64()))))
				if err != nil {
					return err
				}
				err = k.lend.UpdateReserveBalances(ctx, pair.AssetIn, pool.ModuleName, sdk.NewCoin(assetIn.Denom, sdk.NewInt(penaltyToReserveAmount.TruncateInt64())), true)
				if err != nil {
					return err
				}
				cAsset, _ := k.asset.GetAsset(ctx, assetRatesStats.CAssetID)
				updatedLockedVault.AmountIn = updatedLockedVault.AmountIn.Sub(sdk.NewInt(liquidationDeductionAmount.TruncateInt64()))
				lendPos.AmountIn.Amount = lendPos.AmountIn.Amount.Sub(sdk.NewInt(liquidationDeductionAmount.TruncateInt64()))
				lendPos.UpdatedAmountIn = lendPos.UpdatedAmountIn.Sub(sdk.NewInt(liquidationDeductionAmount.TruncateInt64()))
				err = k.bank.BurnCoins(ctx, pool.ModuleName, sdk.NewCoins(sdk.NewCoin(cAsset.Denom, sdk.NewInt(penaltyToReserveAmount.TruncateInt64()))))
				if err != nil {
					return err
				}
				k.lend.SetLend(ctx, lendPos)
			}
			var collateralToBeAuctioned sdk.Dec

			if selloffAmount.GTE(totalIn) {
				collateralToBeAuctioned = totalIn
			} else {

				collateralToBeAuctioned = selloffAmount
			}
			updatedLockedVault.CurrentCollaterlisationRatio = collateralizationRatio
			updatedLockedVault.CollateralToBeAuctioned = collateralToBeAuctioned
			k.SetLockedVault(ctx, updatedLockedVault)
		}
		newUpdatedLockedVault, _ := k.GetLockedVault(ctx, appID, id)
		err := k.auction.LendDutchActivator(ctx, newUpdatedLockedVault)
		if err != nil {
			ctx.Logger().Error("error in dutch activator")
		}
	}
	return nil
}

func (k Keeper) UnLiquidateLockedBorrows(ctx sdk.Context, appID, id uint64, dutchAuction auctiontypes.DutchAuction) error {
	lockedVault, _ := k.GetLockedVault(ctx, appID, id)
	borrowMetadata := lockedVault.GetBorrowMetaData()
	if borrowMetadata != nil {
		lendPos, _ := k.lend.GetLend(ctx, borrowMetadata.LendingId)
		assetInPool, _ := k.lend.GetPool(ctx, lendPos.PoolID)
		firstBridgedAsset, _ := k.asset.GetAsset(ctx, assetInPool.FirstBridgedAssetID)
		userAddress, _ := sdk.AccAddressFromBech32(lockedVault.Owner)
		pair, _ := k.lend.GetLendPair(ctx, lockedVault.ExtendedPairId)
		assetStats, _ := k.lend.GetAssetRatesStats(ctx, pair.AssetIn)
		assetIn, _ := k.asset.GetAsset(ctx, pair.AssetIn)
		assetOut, _ := k.asset.GetAsset(ctx, pair.AssetOut)
		cAssetIn, _ := k.asset.GetAsset(ctx, assetStats.CAssetID)

		if lockedVault.IsAuctionComplete {
			if borrowMetadata.BridgedAssetAmount.IsZero() {
				//also calculate the current collaterlization ratio to ensure there is no sudden changes
				liqThreshold, _ := k.lend.GetAssetRatesStats(ctx, pair.AssetIn)
				unliquidatePointPercentage := liqThreshold.LiquidationThreshold

				if lockedVault.AmountOut.IsZero() {
					err := k.CreateLockedVaultHistory(ctx, lockedVault)
					if err != nil {
						return err
					}
					k.lend.DeleteBorrowForAddressByPair(ctx, userAddress, lockedVault.ExtendedPairId)
					k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
					if err = k.bank.SendCoinsFromModuleToAccount(ctx, assetInPool.ModuleName, userAddress, sdk.NewCoins(sdk.NewCoin(cAssetIn.Denom, lockedVault.AmountIn))); err != nil {
						return err
					}
					lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Add(lockedVault.AmountIn)
					k.lend.SetLend(ctx, lendPos)
				}
				newCalculatedCollateralizationRatio, _ := k.lend.CalculateCollateralizationRatio(ctx, lockedVault.AmountIn, assetIn, lockedVault.UpdatedAmountOut, assetOut)
				if newCalculatedCollateralizationRatio.GT(unliquidatePointPercentage) {
					updatedLockedVault := lockedVault
					updatedLockedVault.CurrentCollaterlisationRatio = newCalculatedCollateralizationRatio
					updatedLockedVault.SellOffHistory = append(updatedLockedVault.SellOffHistory, dutchAuction.String())
					k.SetLockedVault(ctx, updatedLockedVault)
					err := k.UpdateLockedBorrows(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
					if err != nil {
						return nil
					}
				}
				if newCalculatedCollateralizationRatio.LTE(unliquidatePointPercentage) {
					err := k.CreateLockedVaultHistory(ctx, lockedVault)
					if err != nil {
						return err
					}
					k.lend.DeleteBorrowForAddressByPair(ctx, userAddress, lockedVault.ExtendedPairId)
					k.lend.CreteNewBorrow(ctx, lockedVault)
					k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
				}
			} else {
				if borrowMetadata.BridgedAssetAmount.Denom == firstBridgedAsset.Denom {
					liqThresholdAssetIn, _ := k.lend.GetAssetRatesStats(ctx, pair.AssetIn)
					liqThresholdFirstBridgedAsset, _ := k.lend.GetAssetRatesStats(ctx, assetInPool.FirstBridgedAssetID)
					liqThreshold := liqThresholdAssetIn.LiquidationThreshold.Mul(liqThresholdFirstBridgedAsset.LiquidationThreshold)
					unliquidatePointPercentage := liqThreshold

					if lockedVault.AmountOut.IsZero() {
						err := k.CreateLockedVaultHistory(ctx, lockedVault)
						if err != nil {
							return err
						}
						k.lend.DeleteBorrowForAddressByPair(ctx, userAddress, lockedVault.ExtendedPairId)
						k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
						if err = k.bank.SendCoinsFromModuleToAccount(ctx, assetInPool.ModuleName, userAddress, sdk.NewCoins(sdk.NewCoin(cAssetIn.Denom, lockedVault.AmountIn))); err != nil {
							return err
						}
						lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Add(lockedVault.AmountIn)
						k.lend.SetLend(ctx, lendPos)
					}
					newCalculatedCollateralizationRatio, _ := k.lend.CalculateCollateralizationRatio(ctx, lockedVault.AmountIn, assetIn, lockedVault.UpdatedAmountOut, assetOut)
					if newCalculatedCollateralizationRatio.GT(unliquidatePointPercentage) {
						updatedLockedVault := lockedVault
						updatedLockedVault.CurrentCollaterlisationRatio = newCalculatedCollateralizationRatio
						updatedLockedVault.SellOffHistory = append(updatedLockedVault.SellOffHistory, dutchAuction.String())
						k.SetLockedVault(ctx, updatedLockedVault)
						err := k.UpdateLockedBorrows(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
						if err != nil {
							return nil
						}
					}
					if newCalculatedCollateralizationRatio.LTE(unliquidatePointPercentage) {
						err := k.CreateLockedVaultHistory(ctx, lockedVault)
						if err != nil {
							return err
						}
						k.lend.DeleteBorrowForAddressByPair(ctx, userAddress, lockedVault.ExtendedPairId)
						k.lend.CreteNewBorrow(ctx, lockedVault)
						k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
					}
				} else {
					liqThresholdAssetIn, _ := k.lend.GetAssetRatesStats(ctx, pair.AssetIn)
					liqThresholdSecondBridgedAsset, _ := k.lend.GetAssetRatesStats(ctx, assetInPool.SecondBridgedAssetID)
					liqThreshold := liqThresholdAssetIn.LiquidationThreshold.Mul(liqThresholdSecondBridgedAsset.LiquidationThreshold)
					unliquidatePointPercentage := liqThreshold

					if lockedVault.AmountOut.IsZero() {
						err := k.CreateLockedVaultHistory(ctx, lockedVault)
						if err != nil {
							return err
						}
						k.lend.DeleteBorrowForAddressByPair(ctx, userAddress, lockedVault.ExtendedPairId)
						k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
						if err = k.bank.SendCoinsFromModuleToAccount(ctx, assetInPool.ModuleName, userAddress, sdk.NewCoins(sdk.NewCoin(cAssetIn.Denom, lockedVault.AmountIn))); err != nil {
							return err
						}
						lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Add(lockedVault.AmountIn)
						k.lend.SetLend(ctx, lendPos)
					}
					newCalculatedCollateralizationRatio, _ := k.lend.CalculateCollateralizationRatio(ctx, lockedVault.AmountIn, assetIn, lockedVault.UpdatedAmountOut, assetOut)
					if newCalculatedCollateralizationRatio.GT(unliquidatePointPercentage) {
						updatedLockedVault := lockedVault
						updatedLockedVault.CurrentCollaterlisationRatio = newCalculatedCollateralizationRatio
						updatedLockedVault.SellOffHistory = append(updatedLockedVault.SellOffHistory, dutchAuction.String())
						k.SetLockedVault(ctx, updatedLockedVault)
						err := k.UpdateLockedBorrows(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
						if err != nil {
							return nil
						}
					}
					if newCalculatedCollateralizationRatio.LTE(unliquidatePointPercentage) {
						err := k.CreateLockedVaultHistory(ctx, lockedVault)
						if err != nil {
							return err
						}
						k.lend.DeleteBorrowForAddressByPair(ctx, userAddress, lockedVault.ExtendedPairId)
						k.lend.CreteNewBorrow(ctx, lockedVault)
						k.DeleteLockedVault(ctx, lockedVault.AppId, lockedVault.LockedVaultId)
					}
				}
			}
		}
	}
	return nil
}