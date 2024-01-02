package tokenmint_test

import (
	"testing"

	"github.com/MonCatCat/comdex/v13/app"
	"github.com/MonCatCat/comdex/v13/x/tokenmint"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/MonCatCat/comdex/v13/x/tokenmint/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	comdexApp := app.Setup(t, false)
	ctx := comdexApp.BaseApp.NewContext(false, tmproto.Header{})

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	tokenmint.InitGenesis(ctx, comdexApp.TokenmintKeeper, &genesisState)
	got := tokenmint.ExportGenesis(ctx, comdexApp.TokenmintKeeper)
	require.NotNil(t, got)
}
