package types_test

import (
	"github.com/comdex-official/comdex/x/vault/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesisState(),
			valid:    true,
		},
		{
			desc:     "valid genesis state",
			genState: &types.GenesisState{},
			valid:    true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
