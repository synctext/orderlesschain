package contractinterface

import "gitlab.lrz.de/orderless/orderlesschain/protos/goprotos"

type ContractInterface interface {
	Invoke(ShimInterface, *protos.ProposalRequest) (*protos.ProposalResponse, error)
}

