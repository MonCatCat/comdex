package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/MonCatCat/comdex/v13/x/bandoracle/client/cli"
)

var AddFetchPriceHandler = govclient.NewProposalHandler(cli.NewCmdSubmitFetchPriceProposal)
