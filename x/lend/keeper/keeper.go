package keeper

import (
	"fmt"
	assettypes "github.com/comdex-official/comdex/x/asset/types"
	esmtypes "github.com/comdex-official/comdex/x/esm/types"
	"github.com/comdex-official/comdex/x/lend/expected"
	"github.com/comdex-official/comdex/x/lend/types"
	liquidationtypes "github.com/comdex-official/comdex/x/liquidation/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"
	"strconv"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   sdk.StoreKey
		memKey     sdk.StoreKey
		paramstore paramtypes.Subspace
		bank       expected.BankKeeper
		account    expected.AccountKeeper
		asset      expected.AssetKeeper
		market     expected.MarketKeeper
		esm        expected.EsmKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey sdk.StoreKey,
	ps paramtypes.Subspace,
	bank expected.BankKeeper,
	account expected.AccountKeeper,
	asset expected.AssetKeeper,
	market expected.MarketKeeper,
	esm expected.EsmKeeper,

) Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{

		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		bank:       bank,
		account:    account,
		asset:      asset,
		market:     market,
		esm:        esm,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) ModuleBalance(ctx sdk.Context, moduleName string, denom string) sdk.Int {
	return k.bank.GetBalance(ctx, authtypes.NewModuleAddress(moduleName), denom).Amount
}

func uint64InAssetData(a uint64, list []types.AssetDataPoolMapping) bool {
	for _, b := range list {
		if b.AssetID == a {
			return true
		}
	}
	return false
}

func (k Keeper) LendAsset(ctx sdk.Context, lenderAddr string, AssetID uint64, Amount sdk.Coin, PoolID, AppID uint64) error {
	killSwitchParams, _ := k.GetKillSwitchData(ctx, AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}
	asset, found := k.GetAsset(ctx, AssetID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	pool, found := k.GetPool(ctx, PoolID)
	if !found {
		return types.ErrPoolNotFound
	}
	appMapping, found := k.GetApp(ctx, AppID)
	if !found {
		return types.ErrorAppMappingDoesNotExist
	}
	if appMapping.Name != types.AppName {
		return types.ErrorAppMappingIDMismatch
	}

	if Amount.Denom != asset.Denom {
		return sdkerrors.Wrap(types.ErrBadOfferCoinAmount, Amount.Denom)
	}

	found = uint64InAssetData(AssetID, pool.AssetData)
	if !found {
		return sdkerrors.Wrap(types.ErrInvalidAssetIDForPool, strconv.FormatUint(AssetID, 10))
	}

	addr, _ := sdk.AccAddressFromBech32(lenderAddr)

	if k.HasLendForAddressByAsset(ctx, addr, AssetID, PoolID) {
		return types.ErrorDuplicateLend
	}

	loanTokens := sdk.NewCoins(Amount)

	assetRatesStat, found := k.GetAssetRatesStats(ctx, AssetID)
	if !found {
		return sdkerrors.Wrap(types.ErrorAssetRatesStatsNotFound, strconv.FormatUint(AssetID, 10))
	}
	cAsset, _ := k.GetAsset(ctx, assetRatesStat.CAssetID)
	cToken := sdk.NewCoin(cAsset.Denom, Amount.Amount)

	if err := k.bank.SendCoinsFromAccountToModule(ctx, addr, pool.ModuleName, loanTokens); err != nil {
		return err
	}
	// mint c/Token and set new total cToken supply

	cTokens := sdk.NewCoins(cToken)
	if err := k.bank.MintCoins(ctx, pool.ModuleName, cTokens); err != nil {
		return err
	}

	err := k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, addr, cTokens)
	if err != nil {
		return err
	}

	lendID := k.GetUserLendIDHistory(ctx)

	var globalIndex sdk.Dec
	assetStats, _ := k.AssetStatsByPoolIDAndAssetID(ctx, AssetID, PoolID)
	if assetStats.LendApr.IsZero() {
		globalIndex = sdk.OneDec()
	} else {
		globalIndex = assetStats.LendApr
	}

	lendPos := types.LendAsset{
		ID:                  lendID + 1,
		AssetID:             AssetID,
		PoolID:              PoolID,
		Owner:               lenderAddr,
		AmountIn:            Amount,
		LendingTime:         ctx.BlockTime(),
		UpdatedAmountIn:     Amount.Amount,
		AvailableToBorrow:   Amount.Amount,
		Reward_Accumulated:  sdk.ZeroInt(),
		AppID:               AppID,
		GlobalIndex:         globalIndex,
		LastInteractionTime: ctx.BlockTime(),
		CPoolName:           pool.CPoolName,
	}
	k.UpdateLendStats(ctx, AssetID, PoolID, Amount.Amount, true)
	k.SetUserLendIDHistory(ctx, lendPos.ID)
	k.SetLend(ctx, lendPos)
	k.SetLendForAddressByAsset(ctx, addr, lendPos.AssetID, lendPos.ID, lendPos.PoolID)
	err = k.UpdateUserLendIDMapping(ctx, lenderAddr, lendPos.ID, true)
	if err != nil {
		return err
	}
	err = k.UpdateLendIDByOwnerAndPoolMapping(ctx, lenderAddr, lendPos.ID, lendPos.PoolID, true)
	if err != nil {
		return err
	}
	err = k.UpdateLendIDsMapping(ctx, lendPos.ID, true)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) WithdrawAsset(ctx sdk.Context, addr string, lendID uint64, withdrawal sdk.Coin) error {
	lenderAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return err
	}

	lendPos, found := k.GetLend(ctx, lendID)
	if !found {
		return types.ErrLendNotFound
	}
	killSwitchParams, _ := k.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}
	indexGlobalCurrent, err := k.IterateLends(ctx, lendID)
	if err != nil {
		return err
	}
	lendPos, _ = k.GetLend(ctx, lendID)
	lendPos.GlobalIndex = indexGlobalCurrent
	lendPos.LastInteractionTime = ctx.BlockTime()

	getAsset, _ := k.GetAsset(ctx, lendPos.AssetID)
	pool, _ := k.GetPool(ctx, lendPos.PoolID)

	if lendPos.Owner != addr {
		return types.ErrLendAccessUnauthorised
	}

	if withdrawal.Amount.GT(lendPos.AvailableToBorrow) {
		return types.ErrWithdrawAmountLimitExceeds
	}

	if withdrawal.Denom != getAsset.Denom {
		return sdkerrors.Wrap(types.ErrBadOfferCoinAmount, withdrawal.Denom)
	}

	reservedAmount := k.GetReserveFunds(ctx, pool)
	availableAmount := k.ModuleBalance(ctx, pool.ModuleName, withdrawal.Denom)

	if withdrawal.Amount.GT(lendPos.AmountIn.Amount) {
		return sdkerrors.Wrap(types.ErrWithdrawalAmountExceeds, withdrawal.String())
	}

	if withdrawal.Amount.GT(availableAmount.Sub(reservedAmount)) {
		return sdkerrors.Wrap(types.ErrLendingPoolInsufficient, withdrawal.String())
	}

	assetRatesStat, found := k.GetAssetRatesStats(ctx, lendPos.AssetID)
	if !found {
		return sdkerrors.Wrap(types.ErrorAssetRatesStatsNotFound, strconv.FormatUint(lendPos.AssetID, 10))
	}
	cAsset, _ := k.GetAsset(ctx, assetRatesStat.CAssetID)
	cToken := sdk.NewCoin(cAsset.Denom, withdrawal.Amount)
	cTokens := sdk.NewCoins(cToken)

	tokens := sdk.NewCoins(withdrawal)
	if err != nil {
		return err
	}

	if withdrawal.Amount.LT(lendPos.UpdatedAmountIn) {
		if withdrawal.Amount.GTE(lendPos.AmountIn.Amount) {
			if err = k.SendCoinFromAccountToModule(ctx, lenderAddr, pool.ModuleName, cToken); err != nil {
				return err
			}

			//burn c/Token
			err = k.bank.BurnCoins(ctx, pool.ModuleName, cTokens)
			if err != nil {
				return err
			}

			if err = k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, lenderAddr, tokens); err != nil {
				return err
			}
			subtractionFactor := lendPos.UpdatedAmountIn.Sub(withdrawal.Amount)
			subtractionAmount := lendPos.AmountIn.Amount.Sub(subtractionFactor)
			k.UpdateLendStats(ctx, lendPos.AssetID, lendPos.PoolID, subtractionAmount, false)
			lendPos.Reward_Accumulated = sdk.ZeroInt()
			lendPos.AmountIn.Amount = lendPos.AmountIn.Amount.Sub(subtractionAmount)
			lendPos.UpdatedAmountIn = lendPos.UpdatedAmountIn.Sub(withdrawal.Amount)
			lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(withdrawal.Amount)
			k.SetLend(ctx, lendPos)

		} else {
			if err = k.SendCoinFromAccountToModule(ctx, lenderAddr, pool.ModuleName, cToken); err != nil {
				return err
			}
			//burn c/Token
			err = k.bank.BurnCoins(ctx, pool.ModuleName, cTokens)
			if err != nil {
				return err
			}

			if err = k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, lenderAddr, tokens); err != nil {
				return err
			}
			k.UpdateLendStats(ctx, lendPos.AssetID, lendPos.PoolID, withdrawal.Amount, false)
			if withdrawal.Amount.GTE(lendPos.Reward_Accumulated) {
				lendPos.Reward_Accumulated = sdk.ZeroInt()
			} else {
				lendPos.Reward_Accumulated = lendPos.Reward_Accumulated.Sub(withdrawal.Amount)
			}
			lendPos.AmountIn.Amount = lendPos.AmountIn.Amount.Sub(withdrawal.Amount)
			lendPos.UpdatedAmountIn = lendPos.UpdatedAmountIn.Sub(withdrawal.Amount)
			lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(withdrawal.Amount)
			k.SetLend(ctx, lendPos)

		}
	} else {
		return types.ErrWithdrawAmountLimitExceeds
	}
	return nil
}

func (k Keeper) DepositAsset(ctx sdk.Context, addr string, lendID uint64, deposit sdk.Coin) error {
	lenderAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return err
	}

	lendPos, found := k.GetLend(ctx, lendID)
	if !found {
		return types.ErrLendNotFound
	}

	killSwitchParams, _ := k.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}
	indexGlobalCurrent, err := k.IterateLends(ctx, lendID)
	if err != nil {
		return err
	}
	lendPos, _ = k.GetLend(ctx, lendID)
	lendPos.GlobalIndex = indexGlobalCurrent
	lendPos.LastInteractionTime = ctx.BlockTime()
	getAsset, _ := k.GetAsset(ctx, lendPos.AssetID)
	pool, _ := k.GetPool(ctx, lendPos.PoolID)

	if deposit.Denom != getAsset.Denom {
		return sdkerrors.Wrap(types.ErrBadOfferCoinAmount, deposit.Denom)
	}

	assetRatesStat, found := k.GetAssetRatesStats(ctx, lendPos.AssetID)
	if !found {
		return sdkerrors.Wrap(types.ErrorAssetRatesStatsNotFound, strconv.FormatUint(lendPos.AssetID, 10))
	}
	cAsset, _ := k.GetAsset(ctx, assetRatesStat.CAssetID)
	cToken := sdk.NewCoin(cAsset.Denom, deposit.Amount)

	cTokens := sdk.NewCoins(cToken)

	if err = k.bank.MintCoins(ctx, pool.ModuleName, cTokens); err != nil {
		return err
	}

	if err = k.bank.SendCoinsFromAccountToModule(ctx, lenderAddr, pool.ModuleName, sdk.NewCoins(deposit)); err != nil {
		return err
	}

	err = k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, lenderAddr, cTokens)
	if err != nil {
		return err
	}

	lendPos.AmountIn = lendPos.AmountIn.Add(deposit)
	lendPos.UpdatedAmountIn = lendPos.UpdatedAmountIn.Add(deposit.Amount)
	lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Add(deposit.Amount)

	k.UpdateLendStats(ctx, lendPos.AssetID, lendPos.PoolID, deposit.Amount, true)
	k.SetLend(ctx, lendPos)
	return nil
}

func (k Keeper) CloseLend(ctx sdk.Context, addr string, lendID uint64) error {
	lenderAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return err
	}

	lendPos, found := k.GetLend(ctx, lendID)
	if !found {
		return types.ErrLendNotFound
	}

	killSwitchParams, _ := k.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}
	indexGlobalCurrent, err := k.IterateLends(ctx, lendID)
	if err != nil {
		return err
	}
	lendPos, _ = k.GetLend(ctx, lendID)
	lendPos.GlobalIndex = indexGlobalCurrent
	lendPos.LastInteractionTime = ctx.BlockTime()

	pool, _ := k.GetPool(ctx, lendPos.PoolID)

	if lendPos.Owner != addr {
		return types.ErrLendAccessUnauthorised
	}

	lendIDToBorrowIDMapping, _ := k.GetLendIDToBorrowIDMapping(ctx, lendID)
	if lendIDToBorrowIDMapping.BorrowingID != nil {
		return types.ErrBorrowingPositionOpen
	}
	reservedAmount := k.GetReserveFunds(ctx, pool)
	availableAmount := k.ModuleBalance(ctx, pool.ModuleName, lendPos.AmountIn.Denom)

	if lendPos.UpdatedAmountIn.GT(availableAmount.Sub(reservedAmount)) {
		return sdkerrors.Wrap(types.ErrLendingPoolInsufficient, lendPos.UpdatedAmountIn.String())
	}

	tokens := sdk.NewCoins(sdk.NewCoin(lendPos.AmountIn.Denom, lendPos.UpdatedAmountIn))
	assetRatesStat, found := k.GetAssetRatesStats(ctx, lendPos.AssetID)
	if !found {
		return sdkerrors.Wrap(types.ErrorAssetRatesStatsNotFound, strconv.FormatUint(lendPos.AssetID, 10))
	}
	cAsset, _ := k.GetAsset(ctx, assetRatesStat.CAssetID)
	cToken := sdk.NewCoin(cAsset.Denom, lendPos.UpdatedAmountIn)

	if err = k.SendCoinFromAccountToModule(ctx, lenderAddr, pool.ModuleName, cToken); err != nil {
		return err
	}

	cTokens := sdk.NewCoins(cToken)
	err = k.bank.BurnCoins(ctx, pool.ModuleName, cTokens)
	if err != nil {
		return err
	}

	if err = k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, lenderAddr, tokens); err != nil {
		return err
	}

	k.UpdateLendStats(ctx, lendPos.AssetID, lendPos.PoolID, lendPos.AmountIn.Amount, false)
	k.DeleteLendForAddressByAsset(ctx, lenderAddr, lendPos.AssetID, lendPos.PoolID)

	err = k.UpdateUserLendIDMapping(ctx, addr, lendPos.ID, false)
	if err != nil {
		return err
	}
	err = k.UpdateLendIDByOwnerAndPoolMapping(ctx, addr, lendPos.ID, lendPos.PoolID, false)
	if err != nil {
		return err
	}
	err = k.UpdateLendIDsMapping(ctx, lendPos.ID, false)
	if err != nil {
		return err
	}
	k.DeleteLend(ctx, lendPos.ID)

	return nil
}

func uint64InSlice(a uint64, list []uint64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

//nolint:funlen
func (k Keeper) BorrowAsset(ctx sdk.Context, addr string, lendID, pairID uint64, IsStableBorrow bool, AmountIn, loan sdk.Coin) error {
	lenderAddr, _ := sdk.AccAddressFromBech32(addr)

	lendPos, found := k.GetLend(ctx, lendID)
	if !found {
		return types.ErrLendNotFound
	}

	killSwitchParams, _ := k.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}

	if lendPos.Owner != addr {
		return types.ErrLendAccessUnauthorised
	}

	pair, found := k.GetLendPair(ctx, pairID)
	if !found {
		return types.ErrorPairNotFound
	}
	pairMapping, _ := k.GetAssetToPair(ctx, pair.AssetIn, lendPos.PoolID)
	found = uint64InSlice(pairID, pairMapping.PairID)
	if !found {
		return types.ErrorPairNotFound
	}
	//check cr ratio
	assetIn, found := k.GetAsset(ctx, lendPos.AssetID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	assetOut, found := k.GetAsset(ctx, pair.AssetOut)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	assetInRatesStats, found := k.GetAssetRatesStats(ctx, pair.AssetIn)
	if !found {
		return types.ErrAssetStatsNotFound
	}

	cAsset, found := k.GetAsset(ctx, assetInRatesStats.CAssetID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	if AmountIn.Denom != cAsset.Denom {
		return types.ErrBadOfferCoinType
	}

	if k.HasBorrowForAddressByPair(ctx, lenderAddr, pairID) {
		return types.ErrorDuplicateBorrow
	}

	if AmountIn.Amount.GT(lendPos.AvailableToBorrow) {
		return types.ErrAvailableToBorrowInsufficient
	}

	AssetInPool, found := k.GetPool(ctx, lendPos.PoolID)
	if !found {
		return types.ErrPoolNotFound
	}
	AssetOutPool, found := k.GetPool(ctx, pair.AssetOutPoolID)
	if !found {
		return types.ErrPoolNotFound
	}

	if IsStableBorrow && !assetInRatesStats.EnableStableBorrow {
		return sdkerrors.Wrap(types.ErrStableBorrowDisabled, loan.String())
	}

	err := k.VerifyCollateralizationRatio(ctx, AmountIn.Amount, assetIn, loan.Amount, assetOut, assetInRatesStats.Ltv)
	if err != nil {
		return err
	}
	borrowID := k.GetUserBorrowIDHistory(ctx)

	reservedAmount := k.GetReserveFunds(ctx, AssetOutPool)
	availableAmount := k.ModuleBalance(ctx, AssetOutPool.ModuleName, loan.Denom)
	// check sufficient amt in pool to borrow
	if loan.Amount.GT(availableAmount.Sub(reservedAmount)) {
		return sdkerrors.Wrap(types.ErrBorrowingPoolInsufficient, loan.String())
	}
	assetStats, _ := k.AssetStatsByPoolIDAndAssetID(ctx, pair.AssetOut, pair.AssetOutPoolID)
	reserveGlobalIndex, err := k.GetReserveRate(ctx, pair.AssetOutPoolID, pair.AssetOut)
	if err != nil {
		reserveGlobalIndex = sdk.OneDec()
	}
	globalIndex := assetStats.BorrowApr
	if !pair.IsInterPool {
		AmountOut := loan
		// take c/Tokens from the user
		if err = k.SendCoinFromAccountToModule(ctx, lenderAddr, AssetInPool.ModuleName, AmountIn); err != nil {
			return err
		}

		if err = k.SendCoinFromModuleToAccount(ctx, AssetOutPool.ModuleName, lenderAddr, loan); err != nil {
			return err
		}

		var StableBorrowRate sdk.Dec
		if assetInRatesStats.EnableStableBorrow && IsStableBorrow {
			StableBorrowRate, err = k.GetBorrowAPRByAssetID(ctx, AssetOutPool.PoolID, pair.AssetOut, IsStableBorrow)
			if err != nil {
				return err
			}
		} else {
			StableBorrowRate = sdk.ZeroDec()
		}

		borrowPos := types.BorrowAsset{
			ID:                   borrowID + 1,
			LendingID:            lendID,
			PairID:               pairID,
			AmountIn:             AmountIn,
			AmountOut:            AmountOut,
			BridgedAssetAmount:   sdk.NewCoin(loan.Denom, sdk.NewInt(0)),
			IsStableBorrow:       IsStableBorrow,
			StableBorrowRate:     StableBorrowRate,
			BorrowingTime:        ctx.BlockTime(),
			UpdatedAmountOut:     AmountOut.Amount,
			Interest_Accumulated: sdk.ZeroInt(),
			GlobalIndex:          globalIndex,
			ReserveGlobalIndex:   reserveGlobalIndex,
			LastInteractionTime:  ctx.BlockTime(),
			CPoolName:            AssetOutPool.CPoolName,
		}
		k.UpdateBorrowStats(ctx, pair, borrowPos, AmountOut.Amount, true)
		err = k.UpdateBorrowIdsMapping(ctx, borrowPos.ID, true)
		if err != nil {
			return err
		}
		lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(AmountIn.Amount)
		k.SetLend(ctx, lendPos)
		k.SetUserBorrowIDHistory(ctx, borrowPos.ID)
		k.SetBorrow(ctx, borrowPos)
		k.SetBorrowForAddressByPair(ctx, lenderAddr, pairID, borrowPos.ID)
		err = k.UpdateUserBorrowIDMapping(ctx, lendPos.Owner, borrowPos.ID, true)
		if err != nil {
			return err
		}
		err = k.UpdateBorrowIDByOwnerAndPoolMapping(ctx, lendPos.Owner, borrowPos.ID, pair.AssetOutPoolID, true)
		if err != nil {
			return err
		}
		err = k.UpdateLendIDToBorrowIDMapping(ctx, borrowPos.LendingID, borrowPos.ID, true)
		if err != nil {
			return err
		}

	} else {
		updatedAmtIn := AmountIn.Amount.ToDec().Mul(assetInRatesStats.Ltv)
		priceAssetIn, found := k.market.GetTwa(ctx, pair.AssetIn)
		if !found {
			return types.ErrorPriceDoesNotExist
		}
		amtIn := updatedAmtIn.TruncateInt().Mul(sdk.NewIntFromUint64(priceAssetIn.Twa))

		priceFirstBridgedAsset, found := k.market.GetTwa(ctx, AssetInPool.FirstBridgedAssetID)
		if !found {
			return types.ErrorPriceDoesNotExist
		}
		priceSecondBridgedAsset, found := k.market.GetTwa(ctx, AssetInPool.SecondBridgedAssetID)
		if !found {
			return types.ErrorPriceDoesNotExist
		}
		firstBridgedAsset, found := k.GetAsset(ctx, AssetInPool.FirstBridgedAssetID)
		if !found {
			return assettypes.ErrorAssetDoesNotExist
		}
		secondBridgedAsset, found := k.GetAsset(ctx, AssetInPool.SecondBridgedAssetID)
		if !found {
			return assettypes.ErrorAssetDoesNotExist
		}
		// qty of first and second bridged asset to be sent over different pool according to the borrow Pool

		if priceFirstBridgedAsset.Twa == types.Uint64Zero || priceSecondBridgedAsset.Twa == types.Uint64Zero {
			return types.ErrorPriceDoesNotExist
		}

		firstBridgedAssetQty := amtIn.Quo(sdk.NewIntFromUint64(priceFirstBridgedAsset.Twa))
		firstBridgedAssetBal := k.ModuleBalance(ctx, AssetInPool.ModuleName, firstBridgedAsset.Denom)
		secondBridgedAssetQty := amtIn.Quo(sdk.NewIntFromUint64(priceSecondBridgedAsset.Twa))
		secondBridgedAssetBal := k.ModuleBalance(ctx, AssetInPool.ModuleName, secondBridgedAsset.Denom)

		firstBridgedAssetRatesStats, found := k.GetAssetRatesStats(ctx, AssetInPool.FirstBridgedAssetID)
		if !found {
			return types.ErrAssetStatsNotFound
		}
		secondBridgedAssetRatesStats, found := k.GetAssetRatesStats(ctx, AssetInPool.SecondBridgedAssetID)
		if !found {
			return types.ErrAssetStatsNotFound
		}
		if firstBridgedAssetQty.LT(firstBridgedAssetBal) {
			err = k.VerifyCollateralizationRatio(ctx, firstBridgedAssetQty, firstBridgedAsset, loan.Amount, assetOut, firstBridgedAssetRatesStats.Ltv)
			if err != nil {
				return err
			}
			// take c/Tokens from the user
			if err = k.SendCoinFromAccountToModule(ctx, lenderAddr, AssetInPool.ModuleName, AmountIn); err != nil {
				return err
			}
			bridgedAssetAmount := sdk.NewCoin(firstBridgedAsset.Denom, firstBridgedAssetQty)
			if err = k.SendCoinFromModuleToModule(ctx, AssetInPool.ModuleName, AssetOutPool.ModuleName, sdk.NewCoins(bridgedAssetAmount)); err != nil {
				return err
			}

			if err = k.SendCoinFromModuleToAccount(ctx, AssetOutPool.ModuleName, lenderAddr, loan); err != nil {
				return err
			}

			AmountOut := loan

			var StableBorrowRate sdk.Dec
			if assetInRatesStats.EnableStableBorrow && IsStableBorrow {
				StableBorrowRate, err = k.GetBorrowAPRByAssetID(ctx, AssetOutPool.PoolID, pair.AssetOut, IsStableBorrow)
				if err != nil {
					return err
				}
			} else {
				StableBorrowRate = sdk.ZeroDec()
			}

			borrowPos := types.BorrowAsset{
				ID:                   borrowID + 1,
				LendingID:            lendID,
				PairID:               pairID,
				AmountIn:             AmountIn,
				AmountOut:            AmountOut,
				BridgedAssetAmount:   bridgedAssetAmount,
				IsStableBorrow:       IsStableBorrow,
				StableBorrowRate:     StableBorrowRate,
				BorrowingTime:        ctx.BlockTime(),
				UpdatedAmountOut:     AmountOut.Amount,
				Interest_Accumulated: sdk.ZeroInt(),
				GlobalIndex:          globalIndex,
				ReserveGlobalIndex:   reserveGlobalIndex,
				LastInteractionTime:  ctx.BlockTime(),
				CPoolName:            AssetOutPool.CPoolName,
			}
			k.UpdateBorrowStats(ctx, pair, borrowPos, AmountOut.Amount, true)

			err = k.UpdateBorrowIdsMapping(ctx, borrowPos.ID, true)
			if err != nil {
				return err
			}
			lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(AmountIn.Amount)
			k.SetLend(ctx, lendPos)
			k.SetUserBorrowIDHistory(ctx, borrowPos.ID)
			k.SetBorrow(ctx, borrowPos)
			k.SetBorrowForAddressByPair(ctx, lenderAddr, pairID, borrowPos.ID)
			err = k.UpdateUserBorrowIDMapping(ctx, lendPos.Owner, borrowPos.ID, true)
			if err != nil {
				return err
			}
			err = k.UpdateBorrowIDByOwnerAndPoolMapping(ctx, lendPos.Owner, borrowPos.ID, pair.AssetOutPoolID, true)
			if err != nil {
				return err
			}
			err = k.UpdateLendIDToBorrowIDMapping(ctx, borrowPos.LendingID, borrowPos.ID, true)
			if err != nil {
				return err
			}

		} else if secondBridgedAssetQty.LT(secondBridgedAssetBal) {
			err = k.VerifyCollateralizationRatio(ctx, secondBridgedAssetQty, secondBridgedAsset, loan.Amount, assetOut, secondBridgedAssetRatesStats.Ltv)
			if err != nil {
				return err
			}
			// take c/Tokens from the user
			if err = k.SendCoinFromAccountToModule(ctx, lenderAddr, AssetInPool.ModuleName, AmountIn); err != nil {
				return err
			}

			bridgedAssetAmount := sdk.NewCoin(secondBridgedAsset.Denom, secondBridgedAssetQty)
			if err = k.SendCoinFromModuleToModule(ctx, AssetInPool.ModuleName, AssetOutPool.ModuleName, sdk.NewCoins(bridgedAssetAmount)); err != nil {
				return err
			}

			if err = k.SendCoinFromModuleToAccount(ctx, AssetOutPool.ModuleName, lenderAddr, loan); err != nil {
				return err
			}

			AmountOut := loan

			var StableBorrowRate sdk.Dec
			if assetInRatesStats.EnableStableBorrow && IsStableBorrow {
				StableBorrowRate, err = k.GetBorrowAPRByAssetID(ctx, AssetOutPool.PoolID, pair.AssetOut, IsStableBorrow)
				if err != nil {
					return err
				}
			} else {
				StableBorrowRate = sdk.ZeroDec()
			}

			borrowPos := types.BorrowAsset{
				ID:                   borrowID + 1,
				LendingID:            lendID,
				PairID:               pairID,
				AmountIn:             AmountIn,
				AmountOut:            AmountOut,
				BridgedAssetAmount:   bridgedAssetAmount,
				IsStableBorrow:       IsStableBorrow,
				StableBorrowRate:     StableBorrowRate,
				BorrowingTime:        ctx.BlockTime(),
				UpdatedAmountOut:     AmountOut.Amount,
				Interest_Accumulated: sdk.ZeroInt(),
				GlobalIndex:          globalIndex,
				ReserveGlobalIndex:   reserveGlobalIndex,
				LastInteractionTime:  ctx.BlockTime(),
				CPoolName:            AssetOutPool.CPoolName,
			}
			k.UpdateBorrowStats(ctx, pair, borrowPos, AmountOut.Amount, true)
			err = k.UpdateBorrowIdsMapping(ctx, borrowPos.ID, true)
			if err != nil {
				return err
			}
			lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(AmountIn.Amount)
			k.SetLend(ctx, lendPos)
			k.SetUserBorrowIDHistory(ctx, borrowPos.ID)
			k.SetBorrow(ctx, borrowPos)
			k.SetBorrowForAddressByPair(ctx, lenderAddr, pairID, borrowPos.ID)
			err = k.UpdateUserBorrowIDMapping(ctx, lendPos.Owner, borrowPos.ID, true)
			if err != nil {
				return err
			}
			err = k.UpdateBorrowIDByOwnerAndPoolMapping(ctx, lendPos.Owner, borrowPos.ID, pair.AssetOutPoolID, true)
			if err != nil {
				return err
			}
			err = k.UpdateLendIDToBorrowIDMapping(ctx, borrowPos.LendingID, borrowPos.ID, true)
			if err != nil {
				return err
			}

		} else {
			return types.ErrBorrowingPoolInsufficient
		}
	}
	return nil
}

func (k Keeper) RepayAsset(ctx sdk.Context, borrowID uint64, borrowerAddr string, payment sdk.Coin) error {
	borrowPos, found := k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}
	addr, _ := sdk.AccAddressFromBech32(borrowerAddr)
	pair, found := k.GetLendPair(ctx, borrowPos.PairID)
	if !found {
		return types.ErrorPairNotFound
	}
	assetStats, found := k.GetAssetRatesStats(ctx, pair.AssetOut)
	if !found {
		return types.ErrAssetStatsNotFound
	}
	cAsset, found := k.GetAsset(ctx, assetStats.CAssetID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	pool, found := k.GetPool(ctx, pair.AssetOutPoolID)
	if !found {
		return types.ErrPoolNotFound
	}
	lendPos, found := k.GetLend(ctx, borrowPos.LendingID)
	if !found {
		return types.ErrLendNotFound
	}

	killSwitchParams, _ := k.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}

	if lendPos.Owner != borrowerAddr {
		return types.ErrLendAccessUnauthorised
	}
	indexGlobalCurrent, reserveGlobalIndex, err := k.IterateBorrow(ctx, borrowID)
	if err != nil {
		return err
	}
	borrowPos, found = k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}
	if borrowPos.AmountOut.Denom != payment.Denom {
		return types.ErrBadOfferCoinAmount
	}
	borrowPos.GlobalIndex = indexGlobalCurrent
	borrowPos.ReserveGlobalIndex = reserveGlobalIndex
	borrowPos.LastInteractionTime = ctx.BlockTime()

	if payment.Amount.LT(borrowPos.UpdatedAmountOut) {
		if payment.Amount.LTE(borrowPos.Interest_Accumulated) {
			// sending repayment to moduleAcc from borrower
			if err := k.bank.SendCoinsFromAccountToModule(ctx, addr, pool.ModuleName, sdk.NewCoins(payment)); err != nil {
				return err
			}
			borrowPos.UpdatedAmountOut = borrowPos.UpdatedAmountOut.Sub(payment.Amount)
			reservePoolRecords, _ := k.GetReservePoolRecordsForBorrow(ctx, borrowID)
			amtToReservePool := reservePoolRecords.InterestAccumulated
			if amtToReservePool.TruncateInt().LTE(payment.Amount) {
				if amtToReservePool.TruncateInt().LT(sdk.ZeroInt()) {
					return types.ErrReserveRatesNotFound
				}
				if amtToReservePool.TruncateInt().GT(sdk.ZeroInt()) {
					amount := sdk.NewCoin(payment.Denom, amtToReservePool.TruncateInt())
					err := k.SetReserveBalances(ctx, pool.ModuleName, pair.AssetOut, amount)
					if err != nil {
						return err
					}
				}
				if borrowPos.Interest_Accumulated.GT(amtToReservePool.TruncateInt()) {
					err := k.MintCoin(ctx, pool.ModuleName, sdk.NewCoin(cAsset.Denom, payment.Amount.Sub(amtToReservePool.TruncateInt())))
					if err != nil {
						return err
					}
				}
				reservePoolRecords.InterestAccumulated = sdk.ZeroDec()
				k.SetReservePoolRecordsForBorrow(ctx, reservePoolRecords)
			} else {
				if amtToReservePool.TruncateInt().LT(sdk.ZeroInt()) {
					return types.ErrReserveRatesNotFound
				}
				if amtToReservePool.TruncateInt().GT(sdk.ZeroInt()) {
					amount := sdk.NewCoin(payment.Denom, amtToReservePool.TruncateInt())
					err := k.SetReserveBalances(ctx, pool.ModuleName, pair.AssetOut, amount)
					if err != nil {
						return err
					}
				}
				reservePoolRecords.InterestAccumulated = sdk.Dec(reservePoolRecords.InterestAccumulated.TruncateInt().Sub(payment.Amount))
				k.SetReservePoolRecordsForBorrow(ctx, reservePoolRecords)
			}
			borrowPos.Interest_Accumulated = borrowPos.Interest_Accumulated.Sub(payment.Amount)
			k.SetBorrow(ctx, borrowPos)

		} else {
			// sending repayment to moduleAcc from borrower
			if err := k.bank.SendCoinsFromAccountToModule(ctx, addr, pool.ModuleName, sdk.NewCoins(payment)); err != nil {
				return err
			}
			borrowPos.UpdatedAmountOut = borrowPos.UpdatedAmountOut.Sub(payment.Amount)
			amtOut := borrowPos.AmountOut.Amount.Sub(payment.Amount)
			amountOut := amtOut.Add(borrowPos.Interest_Accumulated)
			borrowPos.AmountOut.Amount = amountOut

			reservePoolRecords, _ := k.GetReservePoolRecordsForBorrow(ctx, borrowID)
			amtToReservePool := reservePoolRecords.InterestAccumulated
			if amtToReservePool.TruncateInt().LT(sdk.ZeroInt()) {
				return types.ErrReserveRatesNotFound
			}
			if amtToReservePool.TruncateInt().GT(sdk.ZeroInt()) {
				amount := sdk.NewCoin(payment.Denom, amtToReservePool.TruncateInt())
				err := k.SetReserveBalances(ctx, pool.ModuleName, pair.AssetOut, amount)
				if err != nil {
					return err
				}
			}
			if borrowPos.Interest_Accumulated.GT(amtToReservePool.TruncateInt()) {
				err := k.MintCoin(ctx, pool.ModuleName, sdk.NewCoin(cAsset.Denom, borrowPos.Interest_Accumulated.Sub(amtToReservePool.TruncateInt())))
				if err != nil {
					return err
				}
			}
			amtBorrowStats := payment.Amount.Sub(borrowPos.Interest_Accumulated)
			borrowPos.Interest_Accumulated = sdk.ZeroInt()
			reservePoolRecords.InterestAccumulated = sdk.ZeroDec()
			k.SetReservePoolRecordsForBorrow(ctx, reservePoolRecords)
			k.UpdateBorrowStats(ctx, pair, borrowPos, amtBorrowStats, false)
			k.SetBorrow(ctx, borrowPos)

		}
	} else {
		return types.ErrInvalidRepayment
	}
	return nil
}

func (k Keeper) DepositBorrowAsset(ctx sdk.Context, borrowID uint64, addr string, AmountIn sdk.Coin) error {
	borrowPos, found := k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}
	lendID := borrowPos.LendingID
	pairID := borrowPos.PairID
	lenderAddr, _ := sdk.AccAddressFromBech32(addr)

	lendPos, found := k.GetLend(ctx, lendID)
	if !found {
		return types.ErrLendNotFound
	}

	killSwitchParams, _ := k.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}

	if lendPos.Owner != addr {
		return types.ErrLendAccessUnauthorised
	}
	indexGlobalCurrent, reserveGlobalIndex, err := k.IterateBorrow(ctx, borrowID)
	if err != nil {
		return err
	}
	borrowPos, found = k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}
	borrowPos.GlobalIndex = indexGlobalCurrent
	borrowPos.ReserveGlobalIndex = reserveGlobalIndex
	borrowPos.LastInteractionTime = ctx.BlockTime()
	assetRatesStat, found := k.GetAssetRatesStats(ctx, lendPos.AssetID)
	if !found {
		return sdkerrors.Wrap(types.ErrorAssetRatesStatsNotFound, strconv.FormatUint(lendPos.AssetID, 10))
	}
	cAsset, found := k.GetAsset(ctx, assetRatesStat.CAssetID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	if AmountIn.Denom != cAsset.Denom {
		return sdkerrors.Wrap(types.ErrBadOfferCoinAmount, AmountIn.Denom)
	}
	if AmountIn.Amount.GT(lendPos.AvailableToBorrow) {
		return types.ErrAvailableToBorrowInsufficient
	}
	pair, found := k.GetLendPair(ctx, pairID)
	if !found {
		return types.ErrorPairNotFound
	}
	AssetInPool, found := k.GetPool(ctx, lendPos.PoolID)
	if !found {
		return types.ErrPoolNotFound
	}
	AssetOutPool, found := k.GetPool(ctx, pair.AssetOutPoolID)
	if !found {
		return types.ErrPoolNotFound
	}

	if !pair.IsInterPool {
		if err := k.SendCoinFromAccountToModule(ctx, lenderAddr, AssetInPool.ModuleName, AmountIn); err != nil {
			return err
		}
		lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(AmountIn.Amount)
		k.SetLend(ctx, lendPos)
		borrowPos.AmountIn = borrowPos.AmountIn.Add(AmountIn)
		k.SetBorrow(ctx, borrowPos)

	} else {
		assetIn := lendPos.UpdatedAmountIn
		priceAssetIn, found := k.market.GetTwa(ctx, pair.AssetIn)
		if !found {
			return types.ErrorPriceDoesNotExist
		}
		amtIn := assetIn.Mul(sdk.NewIntFromUint64(priceAssetIn.Twa))

		priceFirstBridgedAsset, found := k.market.GetTwa(ctx, AssetInPool.FirstBridgedAssetID)
		if !found {
			return types.ErrorPriceDoesNotExist
		}
		priceSecondBridgedAsset, found := k.market.GetTwa(ctx, AssetInPool.SecondBridgedAssetID)
		if !found {
			return types.ErrorPriceDoesNotExist
		}
		firstBridgedAsset, found := k.GetAsset(ctx, AssetInPool.FirstBridgedAssetID)
		if !found {
			return assettypes.ErrorAssetDoesNotExist
		}
		secondBridgedAsset, found := k.GetAsset(ctx, AssetInPool.SecondBridgedAssetID)
		if !found {
			return assettypes.ErrorAssetDoesNotExist
		}

		// qty of first and second bridged asset to be sent over different pool according to the borrow Pool

		firstBridgedAssetq := amtIn.Quo(sdk.NewIntFromUint64(priceFirstBridgedAsset.Twa))
		firstBridgedAssetQty := firstBridgedAssetq.ToDec().Mul(assetRatesStat.Ltv)
		firstBridgedAssetBal := k.ModuleBalance(ctx, AssetInPool.ModuleName, firstBridgedAsset.Denom)
		secondBridgedAssetq := amtIn.Quo(sdk.NewIntFromUint64(priceSecondBridgedAsset.Twa))
		secondBridgedAssetQty := secondBridgedAssetq.ToDec().Mul(assetRatesStat.Ltv)
		secondBridgedAssetBal := k.ModuleBalance(ctx, AssetInPool.ModuleName, secondBridgedAsset.Denom)

		if borrowPos.BridgedAssetAmount.Denom == firstBridgedAsset.Denom && firstBridgedAssetQty.LT(firstBridgedAssetBal.ToDec()) {
			// take c/Tokens from the user
			if err := k.SendCoinFromAccountToModule(ctx, lenderAddr, AssetInPool.ModuleName, AmountIn); err != nil {
				return err
			}

			if err := k.SendCoinFromModuleToModule(ctx, AssetInPool.ModuleName, AssetOutPool.ModuleName, sdk.NewCoins(sdk.NewCoin(firstBridgedAsset.Denom, firstBridgedAssetQty.TruncateInt()))); err != nil {
				return err
			}
			lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(AmountIn.Amount)
			k.SetLend(ctx, lendPos)
			borrowPos.AmountIn = borrowPos.AmountIn.Add(AmountIn)
			borrowPos.BridgedAssetAmount.Amount = borrowPos.BridgedAssetAmount.Amount.Add(firstBridgedAssetQty.TruncateInt())
			k.SetBorrow(ctx, borrowPos)

		} else if secondBridgedAssetQty.LT(secondBridgedAssetBal.ToDec()) {
			// take c/Tokens from the user
			if err := k.SendCoinFromAccountToModule(ctx, lenderAddr, AssetInPool.ModuleName, AmountIn); err != nil {
				return err
			}

			if err := k.SendCoinFromModuleToModule(ctx, AssetInPool.ModuleName, AssetOutPool.ModuleName, sdk.NewCoins(sdk.NewCoin(secondBridgedAsset.Denom, secondBridgedAssetQty.TruncateInt()))); err != nil {
				return err
			}
			lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Sub(AmountIn.Amount)
			k.SetLend(ctx, lendPos)
			borrowPos.AmountIn = borrowPos.AmountIn.Add(AmountIn)
			borrowPos.BridgedAssetAmount.Amount = borrowPos.BridgedAssetAmount.Amount.Add(secondBridgedAssetQty.TruncateInt())
			k.SetBorrow(ctx, borrowPos)

		} else {
			return types.ErrBridgeAssetQtyInsufficient
		}
	}
	return nil
}

func (k Keeper) DrawAsset(ctx sdk.Context, borrowID uint64, borrowerAddr string, amount sdk.Coin) error {
	borrowPos, found := k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}
	addr, _ := sdk.AccAddressFromBech32(borrowerAddr)
	pair, found := k.GetLendPair(ctx, borrowPos.PairID)
	if !found {
		return types.ErrorPairNotFound
	}
	pool, found := k.GetPool(ctx, pair.AssetOutPoolID)
	if !found {
		return types.ErrPoolNotFound
	}
	lendPos, found := k.GetLend(ctx, borrowPos.LendingID)
	if !found {
		return types.ErrLendNotFound
	}

	killSwitchParams, _ := k.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}

	if lendPos.Owner != borrowerAddr {
		return types.ErrLendAccessUnauthorised
	}
	indexGlobalCurrent, reserveGlobalIndex, err := k.IterateBorrow(ctx, borrowID)
	if err != nil {
		return err
	}
	borrowPos, found = k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}
	borrowPos.GlobalIndex = indexGlobalCurrent
	borrowPos.ReserveGlobalIndex = reserveGlobalIndex
	borrowPos.LastInteractionTime = ctx.BlockTime()
	if borrowPos.AmountOut.Denom != amount.Denom {
		return types.ErrBadOfferCoinAmount
	}
	assetIn, found := k.GetAsset(ctx, lendPos.AssetID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	assetOut, found := k.GetAsset(ctx, pair.AssetOut)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	assetRatesStats, found := k.GetAssetRatesStats(ctx, pair.AssetIn)
	if !found {
		return types.ErrorAssetStatsNotFound
	}
	err = k.VerifyCollateralizationRatio(ctx, borrowPos.AmountIn.Amount, assetIn, borrowPos.UpdatedAmountOut.Add(amount.Amount), assetOut, assetRatesStats.Ltv)
	if err != nil {
		return err
	}
	if err = k.SendCoinFromModuleToAccount(ctx, pool.ModuleName, addr, amount); err != nil {
		return err
	}
	borrowPos.UpdatedAmountOut = borrowPos.UpdatedAmountOut.Add(amount.Amount)
	borrowPos.AmountOut = borrowPos.AmountOut.Add(amount)
	k.SetBorrow(ctx, borrowPos)
	k.UpdateBorrowStats(ctx, pair, borrowPos, amount.Amount, true)

	return nil
}

func (k Keeper) CloseBorrow(ctx sdk.Context, borrowerAddr string, borrowID uint64) error {
	borrowPos, found := k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}
	addr, _ := sdk.AccAddressFromBech32(borrowerAddr)
	pair, found := k.GetLendPair(ctx, borrowPos.PairID)
	if !found {
		return types.ErrorPairNotFound
	}
	assetStats, found := k.GetAssetRatesStats(ctx, pair.AssetOut)
	if !found {
		return types.ErrAssetStatsNotFound
	}
	cAsset, found := k.GetAsset(ctx, assetStats.CAssetID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	pool, found := k.GetPool(ctx, pair.AssetOutPoolID)
	if !found {
		return types.ErrPoolNotFound
	}
	lendPos, found := k.GetLend(ctx, borrowPos.LendingID)
	if !found {
		return types.ErrLendNotFound
	}
	killSwitchParams, _ := k.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}
	if lendPos.Owner != borrowerAddr {
		return types.ErrLendAccessUnauthorised
	}
	indexGlobalCurrent, reserveGlobalIndex, err := k.IterateBorrow(ctx, borrowID)
	if err != nil {
		return err
	}
	borrowPos, found = k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}
	borrowPos.GlobalIndex = indexGlobalCurrent
	borrowPos.ReserveGlobalIndex = reserveGlobalIndex
	borrowPos.LastInteractionTime = ctx.BlockTime()
	assetInPool, found := k.GetPool(ctx, lendPos.PoolID)
	if !found {
		return types.ErrPoolNotFound
	}
	assetOut, found := k.GetAsset(ctx, pair.AssetOut)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	lenderAddr, _ := sdk.AccAddressFromBech32(lendPos.Owner)
	if borrowPos.UpdatedAmountOut.LTE(sdk.ZeroInt()) {
		return types.ErrInsufficientFunds
	}
	amt := sdk.NewCoins(sdk.NewCoin(assetOut.Denom, borrowPos.UpdatedAmountOut))
	if err = k.bank.SendCoinsFromAccountToModule(ctx, addr, pool.ModuleName, amt); err != nil {
		return err
	}
	if err = k.bank.SendCoinsFromModuleToAccount(ctx, assetInPool.ModuleName, lenderAddr, sdk.NewCoins(borrowPos.AmountIn)); err != nil {
		return err
	}

	reservePoolRecords, _ := k.GetReservePoolRecordsForBorrow(ctx, borrowID)
	amtToReservePool := reservePoolRecords.InterestAccumulated
	if amtToReservePool.TruncateInt().LT(sdk.ZeroInt()) {
		return types.ErrReserveRatesNotFound
	}
	if amtToReservePool.TruncateInt().GT(sdk.ZeroInt()) {
		amount := sdk.NewCoin(assetOut.Denom, amtToReservePool.TruncateInt())
		err = k.SetReserveBalances(ctx, pool.ModuleName, pair.AssetOut, amount)
		if err != nil {
			return err
		}
	}
	if borrowPos.Interest_Accumulated.GT(amtToReservePool.TruncateInt()) {
		err = k.MintCoin(ctx, pool.ModuleName, sdk.NewCoin(cAsset.Denom, borrowPos.Interest_Accumulated.Sub(amtToReservePool.TruncateInt())))
		if err != nil {
			return err
		}
	}
	err = k.UpdateUserBorrowIDMapping(ctx, lendPos.Owner, borrowPos.ID, false)
	if err != nil {
		return err
	}
	err = k.UpdateBorrowIDByOwnerAndPoolMapping(ctx, lendPos.Owner, borrowPos.ID, pair.AssetOutPoolID, false)
	if err != nil {
		return err
	}
	err = k.UpdateBorrowIdsMapping(ctx, borrowPos.ID, false)
	if err != nil {
		return err
	}
	k.DeleteBorrowForAddressByPair(ctx, addr, borrowPos.PairID)
	err = k.UpdateLendIDToBorrowIDMapping(ctx, borrowPos.LendingID, borrowPos.ID, false)
	if err != nil {
		return err
	}

	if pair.IsInterPool {
		if err = k.bank.SendCoinsFromModuleToModule(ctx, pool.ModuleName, assetInPool.ModuleName, sdk.NewCoins(borrowPos.BridgedAssetAmount)); err != nil {
			return err
		}
	}

	k.UpdateBorrowStats(ctx, pair, borrowPos, borrowPos.AmountOut.Amount, false)

	lendPos.AvailableToBorrow = lendPos.AvailableToBorrow.Add(borrowPos.AmountIn.Amount)
	k.SetLend(ctx, lendPos)
	k.DeleteBorrow(ctx, borrowID)

	return nil
}

func (k Keeper) BorrowAlternate(ctx sdk.Context, lenderAddr string, AssetID, PoolID uint64, AmountIn sdk.Coin, PairID uint64, IsStableBorrow bool, AmountOut sdk.Coin, AppID uint64) error {
	killSwitchParams, _ := k.GetKillSwitchData(ctx, AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}
	asset, found := k.GetAsset(ctx, AssetID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	pool, found := k.GetPool(ctx, PoolID)
	if !found {
		return types.ErrPoolNotFound
	}
	_, found = k.GetApp(ctx, AppID)
	if !found {
		return types.ErrorAppMappingDoesNotExist
	}

	if AmountIn.Denom != asset.Denom {
		return sdkerrors.Wrap(types.ErrBadOfferCoinAmount, AmountIn.Denom)
	}

	found = uint64InAssetData(AssetID, pool.AssetData)
	if !found {
		return sdkerrors.Wrap(types.ErrInvalidAssetIDForPool, strconv.FormatUint(AssetID, 10))
	}

	addr, _ := sdk.AccAddressFromBech32(lenderAddr)

	if k.HasLendForAddressByAsset(ctx, addr, AssetID, PoolID) {
		return types.ErrorDuplicateLend
	}

	loanTokens := sdk.NewCoins(AmountIn)

	assetRatesStat, found := k.GetAssetRatesStats(ctx, AssetID)
	if !found {
		return sdkerrors.Wrap(types.ErrorAssetRatesStatsNotFound, strconv.FormatUint(AssetID, 10))
	}
	cAsset, found := k.GetAsset(ctx, assetRatesStat.CAssetID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	cToken := sdk.NewCoin(cAsset.Denom, AmountIn.Amount)

	if err := k.bank.SendCoinsFromAccountToModule(ctx, addr, pool.ModuleName, loanTokens); err != nil {
		return err
	}
	// mint c/Token and set new total cToken supply

	cTokens := sdk.NewCoins(cToken)
	if err := k.bank.MintCoins(ctx, pool.ModuleName, cTokens); err != nil {
		return err
	}

	err := k.bank.SendCoinsFromModuleToAccount(ctx, pool.ModuleName, addr, cTokens)
	if err != nil {
		return err
	}

	lendID := k.GetUserLendIDHistory(ctx)

	var globalIndex sdk.Dec
	assetStats, _ := k.AssetStatsByPoolIDAndAssetID(ctx, AssetID, PoolID)
	if assetStats.LendApr.IsZero() {
		globalIndex = sdk.OneDec()
	} else {
		globalIndex = assetStats.LendApr
	}

	lendPos := types.LendAsset{
		ID:                  lendID + 1,
		AssetID:             AssetID,
		PoolID:              PoolID,
		Owner:               lenderAddr,
		AmountIn:            AmountIn,
		LendingTime:         ctx.BlockTime(),
		UpdatedAmountIn:     AmountIn.Amount,
		AvailableToBorrow:   AmountIn.Amount,
		Reward_Accumulated:  sdk.ZeroInt(),
		AppID:               AppID,
		GlobalIndex:         globalIndex,
		LastInteractionTime: ctx.BlockTime(),
		CPoolName:           pool.CPoolName,
	}
	k.UpdateLendStats(ctx, AssetID, PoolID, AmountIn.Amount, true)
	k.SetUserLendIDHistory(ctx, lendPos.ID)
	k.SetLend(ctx, lendPos)
	k.SetLendForAddressByAsset(ctx, addr, lendPos.AssetID, lendPos.ID, lendPos.PoolID)
	err = k.UpdateUserLendIDMapping(ctx, lenderAddr, lendPos.ID, true)
	if err != nil {
		return err
	}
	err = k.UpdateLendIDByOwnerAndPoolMapping(ctx, lenderAddr, lendPos.ID, lendPos.PoolID, true)
	if err != nil {
		return err
	}
	err = k.UpdateLendIDsMapping(ctx, lendPos.ID, true)
	if err != nil {
		return err
	}

	err = k.BorrowAsset(ctx, lenderAddr, lendPos.ID, PairID, IsStableBorrow, cToken, AmountOut)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) FundModAcc(ctx sdk.Context, moduleName string, assetID uint64, lenderAddr sdk.AccAddress, payment sdk.Coin) error {
	loanTokens := sdk.NewCoins(payment)
	if err := k.bank.SendCoinsFromAccountToModule(ctx, lenderAddr, moduleName, loanTokens); err != nil {
		return err
	}

	asset, found := k.GetAsset(ctx, assetID)
	if !found {
		return types.ErrLendNotFound
	}

	if asset.Denom != payment.Denom {
		return types.ErrBadOfferCoinType
	}

	assetRatesStat, found := k.GetAssetRatesStats(ctx, assetID)
	if !found {
		return sdkerrors.Wrap(types.ErrorAssetRatesStatsNotFound, strconv.FormatUint(assetID, 10))
	}
	cAsset, found := k.GetAsset(ctx, assetRatesStat.CAssetID)
	if !found {
		return assettypes.ErrorAssetDoesNotExist
	}
	cToken := sdk.NewCoin(cAsset.Denom, payment.Amount)

	err := k.MintCoin(ctx, moduleName, cToken)
	if err != nil {
		return err
	}
	depositStats, _ := k.GetDepositStats(ctx)
	var balanceStats []types.BalanceStats
	for _, v := range depositStats.BalanceStats {
		if v.AssetID == assetID {
			v.Amount = v.Amount.Add(payment.Amount)
		}
		balanceStats = append(balanceStats, v)
		newDepositStats := types.DepositStats{BalanceStats: balanceStats}
		k.SetDepositStats(ctx, newDepositStats)
	}

	return nil
}

func (k Keeper) SetReserveBalances(ctx sdk.Context, moduleName string, assetID uint64, payment sdk.Coin) error {
	err := k.UpdateReserveBalances(ctx, assetID, moduleName, payment, true)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) Store(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

func (k Keeper) CreteNewBorrow(ctx sdk.Context, liqBorrow liquidationtypes.LockedVault) {
	kind := liqBorrow.GetBorrowMetaData()
	borrowID := k.GetUserBorrowIDHistory(ctx)
	pair, _ := k.GetLendPair(ctx, liqBorrow.ExtendedPairId)
	AssetOut, _ := k.GetAsset(ctx, pair.AssetOut)
	assetInRatesStats, _ := k.GetAssetRatesStats(ctx, pair.AssetIn)
	AssetRatesStats, _ := k.GetAssetRatesStats(ctx, pair.AssetIn)
	cAssetIn, _ := k.GetAsset(ctx, AssetRatesStats.CAssetID)
	AssetOutPool, _ := k.GetPool(ctx, pair.AssetOutPoolID)
	lendPos, _ := k.GetLend(ctx, kind.LendingId)
	AssetInPool, _ := k.GetPool(ctx, lendPos.PoolID)
	assetStats, _ := k.AssetStatsByPoolIDAndAssetID(ctx, pair.AssetOut, pair.AssetOutPoolID)
	reserveGlobalIndex, err := k.GetReserveRate(ctx, pair.AssetOutPoolID, pair.AssetOut)
	if err != nil {
		reserveGlobalIndex = sdk.OneDec()
	}
	globalIndex := assetStats.BorrowApr

	borrowPos := types.BorrowAsset{
		ID:                   borrowID + 1,
		LendingID:            kind.LendingId,
		PairID:               liqBorrow.ExtendedPairId,
		AmountIn:             sdk.NewCoin(cAssetIn.Denom, liqBorrow.AmountIn),
		AmountOut:            sdk.NewCoin(AssetOut.Denom, liqBorrow.AmountOut),
		BridgedAssetAmount:   kind.BridgedAssetAmount,
		IsStableBorrow:       kind.IsStableBorrow,
		StableBorrowRate:     kind.StableBorrowRate,
		BorrowingTime:        ctx.BlockTime(),
		UpdatedAmountOut:     liqBorrow.AmountOut,
		Interest_Accumulated: sdk.ZeroInt(),
		GlobalIndex:          globalIndex,
		ReserveGlobalIndex:   reserveGlobalIndex,
		LastInteractionTime:  ctx.BlockTime(),
		CPoolName:            AssetOutPool.CPoolName,
	}
	// Adjusting bridged asset qty after auctions
	if kind.BridgedAssetAmount.Amount != sdk.ZeroInt() {
		priceAssetIn, _ := k.market.GetTwa(ctx, pair.AssetIn)
		adjustedBridgedAssetAmt := borrowPos.AmountIn.Amount.ToDec().Mul(assetInRatesStats.Ltv)
		amtIn := adjustedBridgedAssetAmt.TruncateInt().Mul(sdk.NewIntFromUint64(priceAssetIn.Twa))
		priceFirstBridgedAsset, _ := k.market.GetTwa(ctx, AssetInPool.FirstBridgedAssetID)
		priceSecondBridgedAsset, _ := k.market.GetTwa(ctx, AssetInPool.SecondBridgedAssetID)
		firstBridgedAsset, _ := k.GetAsset(ctx, AssetInPool.FirstBridgedAssetID)

		if kind.BridgedAssetAmount.Denom == firstBridgedAsset.Denom {
			firstBridgedAssetQty := amtIn.Quo(sdk.NewIntFromUint64(priceFirstBridgedAsset.Twa))
			diff := borrowPos.BridgedAssetAmount.Amount.Sub(firstBridgedAssetQty)
			if diff.GT(sdk.ZeroInt()) {
				err := k.SendCoinFromModuleToModule(ctx, AssetOutPool.ModuleName, AssetInPool.ModuleName, sdk.NewCoins(sdk.NewCoin(borrowPos.BridgedAssetAmount.Denom, diff)))
				if err != nil {
					return
				}
				borrowPos.BridgedAssetAmount.Amount = firstBridgedAssetQty
			} else {
				newDiff := firstBridgedAssetQty.Sub(borrowPos.BridgedAssetAmount.Amount)
				if newDiff.GT(sdk.ZeroInt()) {
					err := k.SendCoinFromModuleToModule(ctx, AssetInPool.ModuleName, AssetOutPool.ModuleName, sdk.NewCoins(sdk.NewCoin(borrowPos.BridgedAssetAmount.Denom, newDiff)))
					if err != nil {
						return
					}
					borrowPos.BridgedAssetAmount.Amount = firstBridgedAssetQty
				}
			}

		} else {
			secondBridgedAssetQty := amtIn.Quo(sdk.NewIntFromUint64(priceSecondBridgedAsset.Twa))
			diff := borrowPos.BridgedAssetAmount.Amount.Sub(secondBridgedAssetQty)
			if diff.GT(sdk.ZeroInt()) {
				err := k.SendCoinFromModuleToModule(ctx, AssetOutPool.ModuleName, AssetInPool.ModuleName, sdk.NewCoins(sdk.NewCoin(borrowPos.BridgedAssetAmount.Denom, diff)))
				if err != nil {
					return
				}
				borrowPos.BridgedAssetAmount.Amount = secondBridgedAssetQty
			} else {
				newDiff := secondBridgedAssetQty.Sub(borrowPos.BridgedAssetAmount.Amount)
				if newDiff.GT(sdk.ZeroInt()) {
					err := k.SendCoinFromModuleToModule(ctx, AssetInPool.ModuleName, AssetOutPool.ModuleName, sdk.NewCoins(sdk.NewCoin(borrowPos.BridgedAssetAmount.Denom, newDiff)))
					if err != nil {
						return
					}
					borrowPos.BridgedAssetAmount.Amount = secondBridgedAssetQty
				}
			}
		}

	}
	OriginalBorrowID := liqBorrow.OriginalVaultId
	err = k.UpdateLendIDToBorrowIDMapping(ctx, kind.LendingId, OriginalBorrowID, false)
	if err != nil {
		return
	}
	k.UpdateBorrowStats(ctx, pair, borrowPos, borrowPos.AmountOut.Amount, true)
	k.SetBorrow(ctx, borrowPos)
	k.SetUserBorrowIDHistory(ctx, borrowPos.ID)
	err = k.UpdateUserBorrowIDMapping(ctx, lendPos.Owner, borrowPos.ID, true)
	if err != nil {
		return
	}
	err = k.UpdateBorrowIDByOwnerAndPoolMapping(ctx, lendPos.Owner, borrowPos.ID, pair.AssetOutPoolID, true)
	if err != nil {
		return
	}
	err = k.UpdateBorrowIdsMapping(ctx, borrowPos.ID, true)
	if err != nil {
		return
	}
	err = k.UpdateUserBorrowIDMapping(ctx, lendPos.Owner, borrowPos.ID, true)
	if err != nil {
		return
	}
}

func (k Keeper) MsgCalculateBorrowInterest(ctx sdk.Context, borrowerAddr string, borrowID uint64) error {
	borrowPos, found := k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}

	lendPos, found := k.GetLend(ctx, borrowPos.LendingID)
	if !found {
		return types.ErrLendNotFound
	}
	killSwitchParams, _ := k.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}
	if lendPos.Owner != borrowerAddr {
		return types.ErrLendAccessUnauthorised
	}
	indexGlobalCurrent, reserveGlobalIndex, err := k.IterateBorrow(ctx, borrowID)
	if err != nil {
		return err
	}
	borrowPos, found = k.GetBorrow(ctx, borrowID)
	if !found {
		return types.ErrBorrowNotFound
	}
	borrowPos.GlobalIndex = indexGlobalCurrent
	borrowPos.ReserveGlobalIndex = reserveGlobalIndex
	borrowPos.LastInteractionTime = ctx.BlockTime()
	k.SetBorrow(ctx, borrowPos)
	return nil
}

func (k Keeper) MsgCalculateLendRewards(ctx sdk.Context, addr string, lendID uint64) error {
	lendPos, found := k.GetLend(ctx, lendID)
	if !found {
		return types.ErrLendNotFound
	}

	killSwitchParams, _ := k.GetKillSwitchData(ctx, lendPos.AppID)
	if killSwitchParams.BreakerEnable {
		return esmtypes.ErrCircuitBreakerEnabled
	}
	indexGlobalCurrent, err := k.IterateLends(ctx, lendID)
	if err != nil {
		return err
	}
	lendPos, _ = k.GetLend(ctx, lendID)
	lendPos.GlobalIndex = indexGlobalCurrent
	lendPos.LastInteractionTime = ctx.BlockTime()
	if lendPos.Owner != addr {
		return types.ErrLendAccessUnauthorised
	}
	k.SetLend(ctx, lendPos)
	return nil
}
