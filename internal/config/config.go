package config

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/sonikq/gophermart/pkg/generator"
	"github.com/spf13/cast"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	RunAddress string

	AccrualSystemAddress string

	DatabaseURI   string
	DBPoolWorkers int

	CtxTimeOut  time.Duration
	ServiceName string

	TokenSecretKey  string
	TokenExpiration time.Duration
}

func Load(envFiles ...string) (Config, error) {
	const (
		tokenSecretKeyLength = 10
	)
	if len(envFiles) != 0 {
		if err := godotenv.Load(envFiles...); err != nil {
			return Config{}, err
		}
	}

	runAddress := flag.String("a", "", "run address defines on what port and host the server will be started")
	databaseURI := flag.String("d", "", "defines the database connection address")
	accrualSystemAddress := flag.String("r", "", "address of the accrual calculation system")

	var cfg = Config{}

	cfg.RunAddress = getEnvString(*runAddress, "RUN_ADDRESS")

	cfg.AccrualSystemAddress = getEnvString(*accrualSystemAddress, "ACCRUAL_SYSTEM_ADDRESS")

	cfg.DatabaseURI = getEnvString(*databaseURI, "DATABASE_URI")
	cfg.DBPoolWorkers = cast.ToInt(os.Getenv("DB_POOL_WORKERS"))

	cfg.CtxTimeOut = time.Millisecond * time.Duration(cast.ToInt(os.Getenv("CTX_TIMEOUT")))
	cfg.ServiceName = cast.ToString(os.Getenv("SERVICE_NAME"))

	secretKey, err := generator.GeneratePassword(tokenSecretKeyLength)
	if err != nil {
		return Config{}, fmt.Errorf("config.Load: failed generate secret key %w", err)
	}
	cfg.TokenSecretKey = getEnvString(secretKey, "TOKEN_SECRET_KEY")
	cfg.TokenExpiration = time.Hour * time.Duration(cast.ToInt(os.Getenv("TOKEN_EXPIRATION")))

	return cfg, nil
}

func getEnvString(flagValue string, envKey string) string {
	envValue, exists := os.LookupEnv(envKey)
	if exists {
		return envValue
	}
	return flagValue
}

func getEnvInt(flagValue int, envKey string) int {
	envValue, exists := os.LookupEnv(envKey)
	if exists {
		intVal, err := strconv.Atoi(envValue)
		if err != nil {
			log.Printf("cant convert env-key: %s to int", envValue)
			return 1
		}

		return intVal
	}

	return flagValue
}
