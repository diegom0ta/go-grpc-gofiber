package config

type DbConfig struct {
	dbCfg string
}

func (dbcfg DbConfig) LoadDB() string {
	dbcfg.dbCfg = "host=localhost user=postgres password=postgres dbname=myserver port=5433 sslmode=disable"
	return dbcfg.dbCfg
}
