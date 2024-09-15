package app

import (
	"log"

	"github.com/bbquite/go-loyalty/internal/handlers"
	"github.com/bbquite/go-loyalty/internal/services"
	"github.com/bbquite/go-loyalty/internal/storage"
)

func Run() {
	// ctx := context.Background()

	appLogger, err := InitLogger()
	if err != nil {
		log.Fatalf("app logger init error: %v", err)
	}

	appCfg := InitConfig(appLogger)

	db, err := storage.NewDBStorage(appCfg.DatabaseURI)
	if err != nil {
		log.Fatalf("database connection error: %v", err)
	}

	appService := services.NewAppService(db)

	handler, err := handlers.NewHandler(appService, appLogger)
	if err != nil {
		log.Fatalf("handler construction error: %v", err)
	}

	srv := new(Server)
	err = srv.Run(appCfg, handler.InitRoutes(), appLogger)
	if err != nil {
		log.Fatalf("server run error: %v", err)
	}
}
