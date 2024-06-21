package svc

import (
	"bind9-manager-service/internal/config"
	"database/sql"
	"os"
)

type ServiceContext struct {
	Config config.Config
	DB     *sql.DB
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	db, err := InitDB(c.DB)
	if err != nil {
		return nil, err
	}

	return &ServiceContext{
		Config: c,
		DB:     db,
	}, nil
}

// 初始化DB，文件存在则跳过
func InitDB(filepath string) (*sql.DB, error) {
	createTables := `
	CREATE TABLE IF NOT EXISTS zones (
		domain TEXT PRIMARY KEY,
		soa_ttl INTEGER NOT NULL,
		soa_cache_ttl INTEGER NOT NULL,
		soa_expire INTEGER NOT NULL,
		soa_mail_address TEXT NOT NULL,
		soa_primary_name_server TEXT NOT NULL,
		soa_refresh INTEGER NOT NULL,
		soa_retry INTEGER NOT NULL,
		soa_serial INTEGER NOT NULL
	);
	CREATE TABLE IF NOT EXISTS records (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		domain TEXT NOT NULL,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		value TEXT NOT NULL,
		UNIQUE (domain, name, type, value)
	);
	`
	if !FileExists(filepath) {
		db, err := sql.Open("sqlite3", filepath)
		if err != nil {
			return nil, err
		}

		_, err = db.Exec(createTables)
		if err != nil {
			return nil, err
		}

		return db, nil
	}

	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func FileExists(filepath string) bool {
	info, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
