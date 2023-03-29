package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/MonCatCat/comdex/v9/x/bandoracle/client/cli"
	"github.com/MonCatCat/comdex/v9/x/bandoracle/client/rest"
)

var AddFetchPriceHandler = govclient.NewProposalHandler(cli.NewCmdSubmitFetchPriceProposal, rest.SubmitFetchPriceProposalRESTHandler)
