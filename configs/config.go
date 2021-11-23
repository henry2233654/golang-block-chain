package configs

import (
	"os"
)

type Config struct {
	RelationalDB RelationalDB
}

func GetConfig() Config {
	var config Config
	if os.Getenv("IS_CONTAINER") == "TRUE" {
		config = getConfigFormEnv()
	} else {
		config = getConfigByDefault()
	}
	return config
}

func getConfigByDefault() Config {
	var config Config
	/// Postgres connection example///
	config.RelationalDB = RelationalDB{
		Type:     "postgres",
		Host:     "localhost",
		Port:     "5432",
		User:     "user",
		Password: "p@ssw0rd",
		Database: "golang_block_chain",
		SslMode:  "disable",
	}
	// config.RelationalDB = RelationalDB{
	// 	Type:     "sqlite3",
	// 	Database: "./temp.db",
	// }
	return config
}
