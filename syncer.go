package main

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang-block-chain/entities"
	"golang-block-chain/services"
	"golang-block-chain/services/contexts"
	"math/big"
	"time"
)

type Syncer struct {
	client              *ethclient.Client
	ctx                 context.Context
	service             services.IBlock
	blockContextFactory contexts.BlockFactory
	startFrom           int64
}

func NewSyncer(client *ethclient.Client, service services.IBlock, blockContextFactory contexts.BlockFactory, startFrom int64) *Syncer {
	return &Syncer{
		client:              client,
		ctx:                 context.TODO(),
		service:             service,
		blockContextFactory: blockContextFactory,
		startFrom:           startFrom,
	}
}

func (syncer *Syncer) Start() error {
	index := syncer.startFrom
	blocks, err := syncer.service.ListBlocks(syncer.blockContextFactory.NewContext(), 1)
	if err != nil {
		return err
	}
	if len(blocks) != 0 {
		lastBlockInStore := blocks[0]
		if lastBlockInStore.BlockNum >= index {
			index = lastBlockInStore.BlockNum + 1
		}
	}
	for {
		block, err := syncer.getBlockByNumber(index)
		if err != nil {
			return err
		}
		if block == nil {
			time.Sleep(5 * time.Second)
			continue
		}
		// save block
		err = syncer.service.SaveBlock(syncer.blockContextFactory.NewContext(), block)
		if err != nil {
			return err
		}
		index++
	}
}

func (syncer *Syncer) getBlockByNumber(number int64) (*entities.Block, error) {
	n := big.NewInt(number)
	//n = nil
	b, err := syncer.client.BlockByNumber(syncer.ctx, n)
	if err != nil {
		if err == ethereum.NotFound {
			return nil, nil
		}
		return nil, err
	}
	block, err := syncer.parseBlockHeader(b)
	if err != nil {
		return nil, err
	}
	return block, nil
}

func (syncer *Syncer) parseBlockHeader(b *types.Block) (*entities.Block, error) {
	header := b.Header()
	block := &entities.Block{
		BlockNum:     header.Number.Int64(),
		BlockHash:    header.Hash().String(),
		BlockTime:    header.Time,
		ParentHash:   header.ParentHash.String(),
		Transactions: make([]*entities.Transaction, len(b.Transactions())),
	}

	for i, tx := range b.Transactions() {
		receipt, err := syncer.client.TransactionReceipt(syncer.ctx, tx.Hash())
		if err != nil {
			return nil, err
		}
		// the transactions in block must include a receipt
		if receipt != nil {
			msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), b.BaseFee())
			if err != nil {
				return nil, err
			}

			to := ""
			if msg.To() != nil {
				to = msg.To().String()
			}

			transaction := &entities.Transaction{
				TxHash:   tx.Hash().String(),
				BlockNum: block.BlockNum,
				From:     msg.From().String(),
				To:       to,
				Nonce:    msg.Nonce(),
				Data:     msg.Data(),
				Value:    tx.Value().String(),
				Logs:     make([]*entities.TransactionLog, len(receipt.Logs)),
			}
			for j, log := range receipt.Logs {
				transaction.Logs[j] = &entities.TransactionLog{
					TxHash: log.TxHash.String(),
					Index:  log.Index,
					Data:   log.Data,
				}
			}
			block.Transactions[i] = transaction
		}
	}
	return block, nil
}
