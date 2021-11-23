package repositories

import (
	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type GormRepository struct {
	DBContext *GormDBContext
}

func (g *GormRepository) SetDBContext(ctx *GormDBContext) {
	g.DBContext = ctx
}

func (g *GormRepository) DB() *gorm.DB {
	return g.DBContext.DB()
}

func isCausedByUniqueConstraint(err error) bool {
	const (
		ErrMySQLDupEntry            = 1062
		ErrMySQLDupEntryWithKeyName = 1586
		ErrPostgresUniqueViolation  = "23505"
	)
	switch sureTypeErr := err.(type) {
	/// Open it if you using sqlite as db source
	// case sqlite3.Error:
	// 	if sureTypeErr.ExtendedCode == sqlite3.ErrConstraintUnique ||
	// 		sureTypeErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey {
	// 		return true
	// 	}
	case *mysql.MySQLError:
		if sureTypeErr.Number == ErrMySQLDupEntry ||
			sureTypeErr.Number == ErrMySQLDupEntryWithKeyName {
			return true
		}
	case *pgconn.PgError:
		if sureTypeErr.Code == ErrPostgresUniqueViolation {
			return true
		}
	default:
		return false
	}
	return false
}

type UniqueConstrainError struct {
	OriginalErr error
}

func NewUniqueConstrainError(err error) *UniqueConstrainError {
	return &UniqueConstrainError{
		OriginalErr: err,
	}
}

func (e *UniqueConstrainError) Error() string {
	return "UNIQUE constraint failed"
}
