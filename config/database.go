package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"strings"

	"tidb-gin-demo/models"

	mysqlDriver "github.com/go-sql-driver/mysql"
	gormMysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Read connection info from environment with sensible defaults
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "4000"
	}
	name := os.Getenv("DB_NAME")
	if name == "" {
		name = "test"
	}

	// Optional: allow user to provide full DSN in one env var (export once)
	// e.g. export DB_DSN='user:pass@tcp(host:port)/test?charset=utf8mb4&parseTime=True&loc=Local&tls=tidb'
	dsnEnv := os.Getenv("DB_DSN")
	var dsn string
	if dsnEnv != "" {
		dsn = dsnEnv
	} else {
		// TLS config for TiDB Cloud (optional)
		tlsEnabled := os.Getenv("TIDB_TLS") == "true"

		// Build DSN from components
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, name)
		if tlsEnabled {
			dsn += "&tls=tidb"
		}
	}

	// If DSN requests tls=tidb or TIDB_TLS env is set, register TLS config named 'tidb' BEFORE parsing DSN.
	// This avoids ParseDSN failing when a custom tls name (e.g. tidb) is present but not yet registered.
	needTLSEarly := os.Getenv("TIDB_TLS") == "true" || strings.Contains(dsn, "tls=tidb")
	if needTLSEarly {
		tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}
		// optional server name for verification (useful for TiDB Cloud)
		if serverName := os.Getenv("TIDB_TLS_SERVERNAME"); serverName != "" {
			tlsConfig.ServerName = serverName
		} else {
			// try to extract host from DSN using simple string parsing when server name not provided
			// expect pattern like '@tcp(host:port)'
			if idx := strings.Index(dsn, "tcp("); idx != -1 {
				start := idx + len("tcp(")
				if end := strings.Index(dsn[start:], ")"); end != -1 {
					addr := dsn[start : start+end]
					hostOnly := strings.Split(addr, ":")[0]
					tlsConfig.ServerName = hostOnly
				}
			}
		}
		// optional CA file path. If not provided, system root CAs will be used.
		if caPath := os.Getenv("TIDB_TLS_CA"); caPath != "" {
			caCert, err := os.ReadFile(caPath)
			if err != nil {
				log.Fatalf("Failed to read CA file: %v", err)
			}
			roots := x509.NewCertPool()
			if !roots.AppendCertsFromPEM(caCert) {
				log.Fatalf("Failed to append CA cert from %s", caPath)
			}
			tlsConfig.RootCAs = roots
		}

		if err := mysqlDriver.RegisterTLSConfig("tidb", tlsConfig); err != nil {
			log.Fatalf("Failed to register TLS config: %v", err)
		}
	}

	// Parse and normalize DSN to ensure parseTime=true and proper tls param
	cfg, err := mysqlDriver.ParseDSN(dsn)
	if err != nil {
		log.Fatalf("Failed to parse DSN: %v", err)
	}

	// Ensure parseTime is true so datetime columns scan into time.Time
	if cfg.Params == nil {
		cfg.Params = map[string]string{}
	}
	if _, ok := cfg.Params["parseTime"]; !ok {
		cfg.Params["parseTime"] = "true"
	}

	// If we registered TLS early, make sure cfg has tls=tidb so driver uses registered config
	if needTLSEarly {
		cfg.Params["tls"] = "tidb"
	}

	// Final normalized DSN
	dsn = cfg.FormatDSN()

	db, err := gorm.Open(gormMysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = db

	// Auto Migrate the User model
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database connected successfully!")
}
