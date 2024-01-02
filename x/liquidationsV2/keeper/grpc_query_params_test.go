package keeper_test

// import (
// 	"testing"

// 	testkeeper "github.com/MonCatCat/comdex/v13/testutil/keeper"
// 	"github.com/MonCatCat/comdex/v13/x/liquidationsV2/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/stretchr/testify/require"
// )

// func TestParamsQuery(t *testing.T) {
// 	keeper, ctx := testkeeper.NewliqKeeper(t)
// 	wctx := sdk.WrapSDKContext(ctx)
// 	params := types.DefaultParams()
// 	keeper.SetParams(ctx, params)

// 	response, err := keeper.Params(wctx, &types.QueryParamsRequest{})
// 	require.NoError(t, err)
// 	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
// }
