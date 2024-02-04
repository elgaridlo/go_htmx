package configs

import (
	"database/sql"
	"os"
	"project/models"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	PSQL models.PostgresConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func LoadEnvConfig(path string) (Config, error) {
	var cfg Config
	err := godotenv.Load(path)
	if err != nil {
		return cfg, err
	}

	cfg.PSQL = models.DefaultPostgresConfig()

	cfg.CSRF.Key = os.Getenv("CSRF_KEY")

	cfg.CSRF.Secure, err = strconv.ParseBool(os.Getenv("CSRF_SECURE"))
	if err != nil {
		panic(err) // atau penanganan error sesuai kebutuhan Anda
	}

	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")

	return cfg, nil
}

func SetupDatabase(source models.PostgresConfig) *sql.DB {
	db, err := models.Open(source)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxIdleConns(10)

	return db
}
