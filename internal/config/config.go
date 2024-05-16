package config

var dbCfg string

func LoadDB() string {
	dbCfg = "host=localhost user=postgres password=postgres dbname=server port=5433 sslmode=disable"
	return dbCfg
}
