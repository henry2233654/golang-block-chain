package repositories

import (
	"golang-block-chain/entities"
	"gorm.io/gorm/clause"

	"gorm.io/gorm"
)

type TransactionFactory func(ctx *GormDBContext) ITransaction

type ITransaction interface {
	Get(txHash string) (transaction *entities.Transaction, err error)
	Save(transaction *entities.Transaction) (err error)
}

type Transaction struct {
	GormRepository
}

func NewTransaction(ctx *GormDBContext) ITransaction {
	repository := new(Transaction)
	repository.SetDBContext(ctx)
	return repository
}

func (repo *Transaction) Get(txHash string) (*entities.Transaction, error) {
	var transaction entities.Transaction
	err := repo.DB().Preload(clause.Associations).First(&transaction, "tx_hash = ?", txHash).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &transaction, err
}

func (repo *Transaction) Save(transaction *entities.Transaction) (err error) {
	err = repo.DB().Save(transaction).Error
	if err != nil && isCausedByUniqueConstraint(err) {
		err = NewUniqueConstrainError(err)
	}
	return
}
