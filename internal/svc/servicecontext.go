package svc

import (
	"bind9-manager-service/internal/config"
	"bind9-manager-service/internal/middleware"
	"database/sql"
	"os"
	"time"

	"github.com/zeromicro/go-zero/rest"
	"golang.org/x/crypto/bcrypt"
)

type ServiceContext struct {
	Config            config.Config
	DataSource        *sql.DB
	JwtAuthMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	jwtMiddleware := middleware.NewJwtAuthMiddleware(c.JwtAuth.AccessSecret)
	db, err := InitDB(c.DataSource)
	if err != nil {
		return nil, err
	}

	return &ServiceContext{
		Config:            c,
		DataSource:        db,
		JwtAuthMiddleware: jwtMiddleware.Handle,
	}, nil
}

// InitDB 初始化DB，文件存在则跳过
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
	CREATE TABLE IF NOT EXISTS config (
		key TEXT PRIMARY KEY,
		value TEXT NOT NULL
	);
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		role VARCHAR(50) NOT NULL,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL
	);
	CREATE TABLE operation_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(255) NOT NULL,
		operation VARCHAR(255) NOT NULL,
		context TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL
	);
	CREATE TABLE user_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(255) NOT NULL,
		operation TEXT NOT NULL,
		context TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL
	);
	CREATE TABLE login_logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username VARCHAR(255) NOT NULL,
		ip VARCHAR(255) NOT NULL,
		operation TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL
	);
	CREATE INDEX idx_records ON records (domain, name, type, value);
	`
	configContent := `options {
        directory "/var/cache/bind";

        // If there is a firewall between you and nameservers you want
        // to talk to, you may need to fix the firewall to allow multiple
        // ports to talk.  See http://www.kb.cert.org/vuls/id/800113

        // If your ISP provided one or more IP addresses for stable 
        // nameservers, you probably want to use them as forwarders.  
        // Uncomment the following block, and insert the addresses replacing 
        // the all-0's placeholder.

        // forwarders {
        //      0.0.0.0;
        // };

        //========================================================================
        // If BIND logs error messages about the root key being expired,
        // you will need to update your keys.  See https://www.isc.org/bind-keys
        //========================================================================
        // dnssec-validation auto;

        // listen-on-v6 { any; };
        forwarders {
            8.8.8.8;
            8.8.4.4;
            114.114.114.114;
        };
        recursion yes;
        allow-recursion { any; };
        allow-query { any; };
        listen-on port 53 { any; };
};`

	if !FileExists(filepath) {
		db, err := sql.Open("sqlite3", filepath)
		if err != nil {
			return nil, err
		}

		_, err = db.Exec(createTables)
		if err != nil {
			return nil, err
		}

		_, err = db.Exec("INSERT INTO config (key, value) VALUES (?, ?)", "named.conf.options", configContent)
		if err != nil {
			return nil, err
		}

		bytes, err := bcrypt.GenerateFromPassword([]byte("admin"), 14)
		if err != nil {
			return nil, err
		}
		password := string(bytes)

		// 初始化用户
		_, err = db.Exec("INSERT INTO users (username, password, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?)", "admin", password, "admin", time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
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
