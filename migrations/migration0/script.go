package migration0

import (
	"gorm.io/gorm"
)

func Up(db *gorm.DB) (err error) {
	tx := db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	err = initTables(tx)
	if err != nil {
		return err
	}

	err = initData(tx)
	if err != nil {
		return err
	}
	return nil
}

func initTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&Block{},
		&Transaction{},
		&TransactionLog{},
	)
}

func initData(db *gorm.DB) (err error) {
	return nil
}
