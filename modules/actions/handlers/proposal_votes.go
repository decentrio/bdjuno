package handlers

import (
	"fmt"

	"github.com/forbole/callisto/v4/modules/actions/types"
	"github.com/rs/zerolog/log"
)

func ProposalVotes(ctx *types.Context, payload *types.Payload) (interface{}, error) {
	log.Debug().Uint64("proposal id", payload.Input.PropsalID).
		Msg("executing account balance action")

	height, err := ctx.GetHeight(payload)
	if err != nil {
		return nil, err
	}

	res, err := ctx.Sources.GovSource.Votes(height, payload.Input.PropsalID, payload.GetPagination())
	if err != nil {
		return nil, fmt.Errorf("error while getting account balance: %s", err)
	}

	return types.ProposalVotesResponse{
		Votes:      res.Votes,
		Pagination: res.Pagination,
	}, nil
}
