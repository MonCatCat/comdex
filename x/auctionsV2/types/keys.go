package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the module name
	ModuleName = "auctionsV2"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_newauc"
)

var (
	TypePlaceMarketBidRequest    = ModuleName + ":market-bid-request"
	TypePlaceLimitBidRequest     = ModuleName + ":deposit-limit-bid-request"
	TypeCancelLimitBidRequest    = ModuleName + ":cancel-limit-bid-request"
	TypeWithdrawLimitBidRequest  = ModuleName + ":withdraw-limit-bid-request"
	AuctionIDKey                 = []byte{0x01}
	AuctionKeyPrefix             = []byte{0x02}
	LimitAuctionBidIDKey         = []byte{0x03}
	AuctionParamsKey             = []byte{0x04}
	UserLimitBidMappingKeyPrefix = []byte{0x05}
)

func AuctionKey(auctionID uint64) []byte {
	return append(append(AuctionKeyPrefix, sdk.Uint64ToBigEndian(auctionID)...))
}

func UserLimitBidKey(debtTokenID, collateralTokenID uint64, premium, address string) []byte {
	return append(append(append(append(UserLimitBidMappingKeyPrefix, sdk.Uint64ToBigEndian(debtTokenID)...), sdk.Uint64ToBigEndian(collateralTokenID)...), premium...), address...)
}

func UserLimitBidKeyForPremium(debtTokenID, collateralTokenID uint64, premium string) []byte {
	return append(append(append(UserLimitBidMappingKeyPrefix, sdk.Uint64ToBigEndian(debtTokenID)...), sdk.Uint64ToBigEndian(collateralTokenID)...), premium...)
}
