package transactionprocessor

import (
	"github.com/golang/protobuf/proto"
	"gitlab.lrz.de/orderless/orderlesschain/internal/blockprocessor"
	"gitlab.lrz.de/orderless/orderlesschain/internal/config"
	"gitlab.lrz.de/orderless/orderlesschain/internal/connection/connpool"
	"gitlab.lrz.de/orderless/orderlesschain/internal/logger"
	"gitlab.lrz.de/orderless/orderlesschain/internal/profiling"
	"gitlab.lrz.de/orderless/orderlesschain/internal/transaction"
	"gitlab.lrz.de/orderless/orderlesschain/protos/goprotos"
	"reflect"
	"strconv"
	"sync"
)

type blockResponseSubscriber struct {
	stream   protos.OrdererService_SubscribeBlocksServer
	finished chan<- bool
}

type OrdererProcessor struct {
	blockProcessor                    *blockprocessor.Processor
	transactionChan                   chan *protos.Transaction
	subscribersLock                   *sync.RWMutex
	txJournal                         *transaction.NodeTransactionJournal
	blockSubscribers                  map[string]*blockResponseSubscriber
	inExperimentParticipatingOrderers []string
	inExperimentParticipatingNodes    []string
	inExperimentParticipatingClients  []string
	isThisOrdererParticipating        bool
}

func InitTransactionOrdererProcessor() *OrdererProcessor {
	tempProcessor := &OrdererProcessor{
		transactionChan:  make(chan *protos.Transaction),
		subscribersLock:  &sync.RWMutex{},
		blockSubscribers: map[string]*blockResponseSubscriber{},
	}
	tempProcessor.setInExperimentParticipatingComponents()
	tempProcessor.setIsThisOrdererParticipating()
	if !tempProcessor.isThisOrdererParticipating {
		logger.InfoLogger.Println("This orderer in NOT participating in the experiment")
		return tempProcessor
	}
	tempProcessor.txJournal = transaction.InitTransactionJournal()
	tempProcessor.blockProcessor = blockprocessor.InitBlockProcessor()
	go tempProcessor.txJournal.RunTransactionsQueueProcessorTicker()
	go tempProcessor.runTransactionProcessor()
	go tempProcessor.runTransactionQueueProcessing()
	return tempProcessor
}

func (p *OrdererProcessor) ProcessTransactionStream(tx *protos.Transaction) {
	p.txJournal.AddTransactionToQueue(tx)
}

func (p *OrdererProcessor) runTransactionQueueProcessing() {
	for {
		transactions := <-p.txJournal.DequeuedTransactionsChan
		go p.processDequeuedTransactions(transactions)
	}
}

func (p *OrdererProcessor) processDequeuedTransactions(transactions *transaction.DequeuedTransactions) {
	for _, tx := range transactions.DequeuedTransactions {
		p.blockProcessor.TransactionChan <- tx
	}
}

func (p *OrdererProcessor) runTransactionProcessor() {
	for {
		block := <-p.blockProcessor.MinedBlock
		go p.sendBlockToSubscriber(block.Block)
	}
}

func (p *OrdererProcessor) BlockSubscription(subscriberEvent *protos.BlockEventSubscription, stream protos.OrdererService_SubscribeBlocksServer) error {
	finished := make(chan bool)
	p.subscribersLock.Lock()
	p.blockSubscribers[subscriberEvent.NodeId] = &blockResponseSubscriber{
		stream:   stream,
		finished: finished,
	}
	p.subscribersLock.Unlock()
	cntx := stream.Context()
	for {
		select {
		case <-finished:
			return nil
		case <-cntx.Done():
			return nil
		}
	}
}

func (p *OrdererProcessor) sendBlockToSubscriber(block *protos.Block) {
	p.subscribersLock.RLock()
	nodes := reflect.ValueOf(p.blockSubscribers).MapKeys()
	p.subscribersLock.RUnlock()
	if profiling.IsBandwidthProfiling {
		_ = proto.Size(block)
	}
	for _, node := range nodes {
		nodeId := node.String()
		p.subscribersLock.RLock()
		streamer, ok := p.blockSubscribers[nodeId]
		p.subscribersLock.RUnlock()
		if !ok {
			logger.ErrorLogger.Println("Node was not found in the subscribers streams.", nodeId)
			return
		}
		if err := streamer.stream.Send(block); err != nil {
			streamer.finished <- true
			logger.ErrorLogger.Println("Could not send the response to the node " + nodeId)
			p.subscribersLock.Lock()
			delete(p.blockSubscribers, nodeId)
			p.subscribersLock.Unlock()
		}
	}
}

func (p *OrdererProcessor) setInExperimentParticipatingComponents() {
	for i := 0; i < config.Config.TotalOrdererCount; i++ {
		p.inExperimentParticipatingOrderers = append(p.inExperimentParticipatingOrderers, "orderer"+strconv.Itoa(i))
	}
	for i := 0; i < config.Config.TotalNodeCount; i++ {
		p.inExperimentParticipatingNodes = append(p.inExperimentParticipatingNodes, "node"+strconv.Itoa(i))
	}
	for i := 0; i < config.Config.TotalClientCount; i++ {
		p.inExperimentParticipatingClients = append(p.inExperimentParticipatingClients, "client"+strconv.Itoa(i))
	}
}

func (p *OrdererProcessor) setIsThisOrdererParticipating() {
	currentOrdererId := connpool.GetComponentPseudoName()
	for _, orderer := range p.inExperimentParticipatingOrderers {
		if currentOrdererId == orderer {
			p.isThisOrdererParticipating = true
		}
	}
}
