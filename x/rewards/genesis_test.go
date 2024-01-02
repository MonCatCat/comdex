package rewards_test

import (
	"testing"

	"github.com/MonCatCat/comdex/v13/app"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/MonCatCat/comdex/v13/x/rewards"
	"github.com/MonCatCat/comdex/v13/x/rewards/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	comdexApp := app.Setup(t, false)
	ctx := comdexApp.BaseApp.NewContext(false, tmproto.Header{})

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	rewards.InitGenesis(ctx, comdexApp.Rewardskeeper, &genesisState)
	got := rewards.ExportGenesis(ctx, comdexApp.Rewardskeeper)
	require.NotNil(t, got)
}
