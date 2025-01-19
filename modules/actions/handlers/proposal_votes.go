package handlers

import (
	"fmt"

	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/forbole/callisto/v4/modules/actions/types"
	"github.com/rs/zerolog/log"
)

func ProposalVotesHandler(ctx *types.Context, payload *types.Payload) (interface{}, error) {
	log.Debug().Uint64("proposal id", payload.Input.PropsalID).
		Msg("executing account proposal votes action")

	height, err := ctx.GetHeight(payload)
	if err != nil {
		return nil, err
	}

	res, err := ctx.Sources.GovSource.Votes(height, payload.Input.PropsalID, payload.GetPagination())
	if err != nil {
		return nil, fmt.Errorf("error while getting votes: %s", err)
	}

	return types.ProposalVotesResponse{
		Votes:      convertVotes(res.Votes),
		Pagination: res.Pagination,
	}, nil
}

func convertVotes(votes govtypesv1.Votes) []types.Vote {
	converted := make([]types.Vote, len(votes))
	for index, vote := range votes {
		converted[index] = types.Vote{
			ProposalId: vote.ProposalId,
			Voter:      vote.Voter,
			Options:    vote.Options,
			Metadata:   vote.Metadata,
		}
	}
	return converted
}
