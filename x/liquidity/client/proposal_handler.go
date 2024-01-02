package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/MonCatCat/comdex/v13/x/liquidity/client/cli"
)

var LiquidityProposalHandler = []govclient.ProposalHandler{
	govclient.NewProposalHandler(cli.NewCmdUpdateGenericParamsProposal),
	govclient.NewProposalHandler(cli.NewCmdCreateNewLiquidityPairProposal),
}
