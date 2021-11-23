package controllers

import (
	"github.com/gin-gonic/gin"
	"golang-block-chain/services"
	"golang-block-chain/services/contexts"
	"golang-block-chain/web/cores"
	"golang-block-chain/web/cores/di"
	"golang-block-chain/web/viewmodels"
	"net/http"
)

type Transaction struct {
	Service    services.IBlock
	CtxFactory *contexts.BlockFactory
}

func (ctr *Transaction) Route(e *gin.Engine) {
	transaction := e.Group("/transaction")
	{
		transaction.GET("/:id", di.AutoDi(ctr.GetTransaction)...)
	}
}

// GetTransaction
// @Summary Get Transaction
// @Description GetBlock block by block number
// @ID get-transaction
// @Tags Transaction
// @Param tx_hash path string true "transaction hash"
// @Success 200 {object} viewmodels.GetTransaction "Success."
// @Success 404 {object} viewmodels.Error "Not found."
// @Failure 500 {object} viewmodels.Error "Internal error."
// @Router /transaction/{tx_hash} [get]
func (ctr *Transaction) GetTransaction(txHash string) (*cores.Response, error) {
	ctx := ctr.CtxFactory.NewContext()
	transaction, err := ctr.Service.RetrieveTransaction(ctx, txHash)
	if err != nil {
		return useDefaultErrorHandler(err)
	}

	logs := make([]viewmodels.TransactionLog, len(transaction.Logs))
	for i, log := range transaction.Logs {
		logs[i] = viewmodels.TransactionLog{
			TxHash: log.TxHash,
			Index:  log.Index,
			Data:   log.Data,
		}
	}
	respBody := viewmodels.GetTransaction{
		TxHash:    transaction.TxHash,
		BlockHash: transaction.BlockHash,
		From:      transaction.From,
		To:        transaction.To,
		Nonce:     transaction.Nonce,
		Data:      transaction.Data,
		Value:     transaction.Value,
		Logs:      logs,
	}
	return cores.NewResponse(http.StatusOK, respBody), nil
}
