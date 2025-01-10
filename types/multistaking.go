package types

type MSValidatorInfo struct {
	ConsensusAddress string
	Denom            string
}

func NewMSValidatorInfo(consAddr, denom string) MSValidatorInfo {
	return MSValidatorInfo{
		ConsensusAddress: consAddr,
		Denom:            denom,
	}
}
