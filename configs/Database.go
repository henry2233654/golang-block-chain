package configs

type RelationalDB struct {
	Type     string
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SslMode  string
}
