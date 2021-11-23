package migrations

import (
	"golang-block-chain/migrations/migration0"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	_ = db.AutoMigrate(&Record{})
	noToExist := loadRecords(db)
	fn := func(no string, up func(*gorm.DB) error) {
		if noToExist[no] {
			return
		}
		err := up(db)
		if err != nil {
			panic(err)
		}
		writeRecord(db, no)
	}
	fn("0", migration0.Up)
}

func loadRecords(db *gorm.DB) map[string]bool {
	var nos []string
	err := db.Model(&Record{}).Pluck("no", &nos).Error
	if err != nil {
		panic(err)
	}

	noToExist := make(map[string]bool)
	for _, no := range nos {
		noToExist[no] = true
	}
	return noToExist
}

func writeRecord(db *gorm.DB, no string) {
	r := new(Record)
	r.No = no
	db.Save(r)
}
