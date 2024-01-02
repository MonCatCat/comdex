package keeper_test

import (
	"testing"

	"github.com/MonCatCat/comdex/v13/app"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/MonCatCat/comdex/v13/x/rewards/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	comdexApp := app.Setup(t, false)
	ctx := comdexApp.BaseApp.NewContext(false, tmproto.Header{})

	params := types.DefaultParams()

	comdexApp.Rewardskeeper.SetParams(ctx, params)

	require.EqualValues(t, params, comdexApp.Rewardskeeper.GetParams(ctx))
}
