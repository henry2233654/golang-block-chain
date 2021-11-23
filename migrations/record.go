package migrations

import (
	"time"
)

type Record struct {
	No        string    `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"Column:created_at"`
}
