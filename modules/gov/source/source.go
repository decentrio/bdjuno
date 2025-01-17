package source

import (
	"github.com/cosmos/cosmos-sdk/types/query"
	govtypesv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

type Source interface {
	Proposal(height int64, id uint64) (*govtypesv1.Proposal, error)
	ProposalDeposit(height int64, id uint64, depositor string) (*govtypesv1.Deposit, error)
	TallyResult(height int64, proposalID uint64) (*govtypesv1.TallyResult, error)
	Params(height int64) (*govtypesv1.Params, error)
	Votes(height int64, proposalId uint64, pagination *query.PageRequest) (govtypesv1.Votes, error)
}
