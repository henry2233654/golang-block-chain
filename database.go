package main

import (
	"errors"
	"fmt"
	"golang-block-chain/configs"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitializeDb(config configs.RelationalDB) (dsn string, db *gorm.DB) {
	CreateDatabase(config)
	dsn, db = ConnectDb(config)
	db.Config.FullSaveAssociations = true
	db.Config.DisableForeignKeyConstraintWhenMigrating = true
	db.Config.NamingStrategy = schema.NamingStrategy{
		SingularTable: true,
	}
	return
}

func CreateDatabase(config configs.RelationalDB) {
	database := config.Database
	switch config.Type {
	case "mssql":
		config.Database = "master"
	case "mysql":
		config.Database = "mysql"
	case "postgres":
		config.Database = "postgres"
	case "sqlite3":
		return
	default:
		panic(errors.New("Invalid database type."))
	}
	_, db := ConnectDb(config)
	db.Exec("CREATE DATABASE " + database)
	sqlDb, _ := db.DB()
	sqlDb.Close()
}

func ConnectDb(config configs.RelationalDB) (string, *gorm.DB) {
	var dsn string
	var db *gorm.DB
	var err error
	switch config.Type {
	case "mssql":
		dsn, db, err = ConnectMssqlDb(config)
	case "mysql":
		dsn, db, err = ConnectMysqlDb(config)
	case "postgres":
		dsn, db, err = ConnectPostgresDb(config)
	case "sqlite3":
		dsn, db, err = ConnectSqliteDb(config)
	default:
		panic(errors.New("Invalid database type."))
	}
	if err != nil {
		panic(err)
	}
	return dsn, db
}

func ConnectPostgresDb(config configs.RelationalDB) (string, *gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Database,
		config.Password,
		config.SslMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return dsn, db, err
}

func ConnectMysqlDb(config configs.RelationalDB) (string, *gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return dsn, db, err
}

func ConnectMssqlDb(config configs.RelationalDB) (string, *gorm.DB, error) {
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	return dsn, db, err
}

func ConnectSqliteDb(config configs.RelationalDB) (string, *gorm.DB, error) {
	dsn := config.Database
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	return dsn, db, err
}

func CloseDb(db *gorm.DB) {
	sqlDb, _ := db.DB()
	sqlDb.Close()
}
