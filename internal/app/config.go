package app

import (
	"encoding/json"
	"flag"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const (
	defRunAddress    string = "localhost:8081"
	defAccuralAdress string = "localhost:8080"
	defDatabase      string = "host=localhost user=postgres password=123 dbname=golang sslmode=disable"
)

type Config struct {
	RunAddress     string `json:"RUN_ADDRESS"`
	AccuralAddress string `json:"ACCRUAL_SYSTEM_ADDRESS"`
	DatabaseURI    string `json:"DATABASE_URI"`
}

func InitConfig(logger *zap.SugaredLogger) *Config {
	cfg := new(Config)

	err := godotenv.Load()
	if err != nil {
		logger.Info(".env file not found")
	}

	if envRunAddress, ok := os.LookupEnv("RUN_ADDRESS"); ok {
		cfg.RunAddress = envRunAddress
	} else {
		flag.StringVar(&cfg.RunAddress, "a", defRunAddress, "RUN_ADDRESS")
	}

	if envAccuralAdress, ok := os.LookupEnv("RUN_ADDRESS"); ok {
		cfg.AccuralAddress = envAccuralAdress
	} else {
		flag.StringVar(&cfg.AccuralAddress, "r", defAccuralAdress, "ACCRUAL_SYSTEM_ADDRESS")
	}

	if envDATABASE, ok := os.LookupEnv("DATABASE_URI"); ok {
		cfg.DatabaseURI = envDATABASE
	} else {
		flag.StringVar(&cfg.DatabaseURI, "d", defDatabase, "DATABASE_URI")
	}

	flag.Parse()

	jsonOptions, _ := json.Marshal(cfg)
	logger.Infof("Server run with this config: %s", jsonOptions)

	return cfg
}
