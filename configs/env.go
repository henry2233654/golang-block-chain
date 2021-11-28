package configs

import (
	"fmt"
	"os"
)

func getConfigFormEnv() Config {
	var config Config
	DbType := getEnvAndPanicIfNotExist("DB_TYPE")
	DbHost := getEnvAndPanicIfNotExist("DB_HOST")
	DbPort := getEnvAndPanicIfNotExist("DB_PORT")
	BlockChainUrl := getEnvAndPanicIfNotExist("BLOCK_CHAIN_URL")
	BlockChainFrom := getEnvAndPanicIfNotExist("BLOCK_CHAIN_FROM")

	config.RelationalDB = RelationalDB{
		Type:     DbType,
		Host:     DbHost,
		Port:     DbPort,
		User:     "user",
		Password: "p@ssw0rd",
		Database: "golang_block_chain",
		SslMode:  "disable",
	}
	config.BlockChain = BlockChain{
		RpcUrl: BlockChainUrl,
		From:   BlockChainFrom,
	}
	return config
}

func getEnvAndPanicIfNotExist(env string) string {
	variable := os.Getenv(env)
	if variable == "" {
		panic(fmt.Sprintf("Didn't set the environment variable [%s].", env))
	}
	return variable
}
