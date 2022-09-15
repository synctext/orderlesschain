package contracts

import (
	"errors"
	"gitlab.lrz.de/orderless/orderlesschain/contractsbenchmarks/contracts/englishauctioncontractfabric"
	"gitlab.lrz.de/orderless/orderlesschain/contractsbenchmarks/contracts/englishauctioncontractfabriccrdt"
	"gitlab.lrz.de/orderless/orderlesschain/contractsbenchmarks/contracts/englishauctioncontractorderlesschain"
	"gitlab.lrz.de/orderless/orderlesschain/contractsbenchmarks/contracts/evotingcontractfabric"
	"gitlab.lrz.de/orderless/orderlesschain/contractsbenchmarks/contracts/evotingcontractfabriccrdt"
	"gitlab.lrz.de/orderless/orderlesschain/contractsbenchmarks/contracts/evotingcontractorderlesschain"
	"gitlab.lrz.de/orderless/orderlesschain/contractsbenchmarks/contracts/syntheticcontractorderlesschain"
	"gitlab.lrz.de/orderless/orderlesschain/internal/contract/contractinterface"
	"strings"
)

func GetContract(contractName string) (contractinterface.ContractInterface, error) {
	contractName = strings.ToLower(contractName)
	switch contractName {
	case "evotingcontractfabric":
		return evotingcontractfabric.NewContract(), nil
	case "evotingcontractorderlesschain":
		return evotingcontractorderlesschain.NewContract(), nil
	case "evotingcontractfabriccrdt":
		return evotingcontractfabriccrdt.NewContract(), nil
	case "englishauctioncontractfabric":
		return englishauctioncontractfabric.NewContract(), nil
	case "englishauctioncontractorderlesschain":
		return englishauctioncontractorderlesschain.NewContract(), nil
	case "englishauctioncontractfabriccrdt":
		return englishauctioncontractfabriccrdt.NewContract(), nil
	case "syntheticcontractorderlesschain":
		return syntheticcontractorderlesschain.NewContract(), nil
	default:
		return nil, errors.New("contract not found")
	}
}

func GetContractNames() []string {
	return []string{
		"englishauctioncontractorderlesschain",
		"englishauctioncontractfabric",
		"englishauctioncontractfabriccrdt",
		"evotingcontractfabric",
		"evotingcontractfabriccrdt",
		"evotingcontractorderlesschain",
		"syntheticcontractorderlesschain",
	}
}
