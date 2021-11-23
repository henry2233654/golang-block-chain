package controllers

import (
	"github.com/gin-gonic/gin"
	"golang-block-chain/services"
	"golang-block-chain/services/contexts"
	"golang-block-chain/web/cores"
	"golang-block-chain/web/cores/di"
	"golang-block-chain/web/inputmodels"
	"golang-block-chain/web/viewmodels"
	"net/http"
)

type Block struct {
	Service    services.IBlock
	CtxFactory *contexts.BlockFactory
}

func (ctr *Block) Route(e *gin.Engine) {
	block := e.Group("/blocks")
	{
		block.GET("", di.AutoDi(ctr.ListBlocks)...)
		block.GET("/:id", di.AutoDi(ctr.GetBlock)...)
	}
}

// ListBlocks
// @Summary ListBlocks Blocks
// @Description ListBlocks blocks by list parameters
// @ID list-blocks
// @Tags Block
// @Param ListBlocksParam query inputmodels.ListBlocksParam true "ListBlocksParam"
// @Success 200 {object} viewmodels.ListBlocks "Success."
// @Failure 500 {object} viewmodels.Error "Internal error."
// @Router /blocks [get]
func (ctr *Block) ListBlocks(listBlocksParam inputmodels.ListBlocksParam) (*cores.Response, error) {
	ctx := ctr.CtxFactory.NewContext()
	blocks, err := ctr.Service.ListBlocks(ctx, listBlocksParam.Limit)
	if err != nil {
		return useDefaultErrorHandler(err)
	}

	respBody := viewmodels.ListBlocks{
		Blocks: make([]viewmodels.ListBlockItem, len(blocks)),
	}
	for i, block := range blocks {
		respBody.Blocks[i] = viewmodels.ListBlockItem{
			BlockNum:   block.BlockNum,
			BlockHash:  block.BlockHash,
			BlockTime:  block.BlockTime,
			ParentHash: block.ParentHash,
		}
	}
	resp := cores.NewResponse(http.StatusOK, respBody)
	return resp, nil
}

// GetBlock
// @Summary Get Block
// @Description Get block by block number
// @ID get-block
// @Tags Block
// @Param block_num path string true "Block Number"
// @Success 200 {object} viewmodels.GetSingleBlock "Success."
// @Success 404 {object} viewmodels.Error "Not found."
// @Failure 500 {object} viewmodels.Error "Internal error."
// @Router /blocks/{block_num} [get]
func (ctr *Block) GetBlock(blockNum uint64) (*cores.Response, error) {
	ctx := ctr.CtxFactory.NewContext()
	block, err := ctr.Service.RetrieveBlock(ctx, blockNum)
	if err != nil {
		return useDefaultErrorHandler(err)
	}
	transactions := make([]string, len(block.Transactions))
	for i, transaction := range block.Transactions {
		transactions[i] = transaction.TxHash
	}
	respBody := viewmodels.GetSingleBlock{
		BlockNum:     block.BlockNum,
		BlockHash:    block.BlockHash,
		BlockTime:    block.BlockTime,
		ParentHash:   block.ParentHash,
		Transactions: transactions,
	}
	return cores.NewResponse(http.StatusOK, respBody), nil
}
