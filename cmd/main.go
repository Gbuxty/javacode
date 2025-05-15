package main

import (
	"context"
	"java-code/config"
	"java-code/internal/adapter/postgres"
	"java-code/internal/adapter/transport/http/handler"
	"java-code/internal/adapter/transport/http/server"
	"java-code/internal/adapter/txManager"
	"java-code/internal/service"
	"java-code/pkg/pgdb"
	"java-code/pkg/logger"
	"java-code/pkg/migrations"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	if err:=run();err!=nil{
		log.Fatalf("Failed app start:%v", err)
	}
}

func run() error {
	logger, err := logger.NewLogger()
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}

	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Fatal("failed to sync logger: %v", zap.Error(err))
		}
	}()

	cfg, err := config.New()
	if err != nil {
		logger.Fatal("create config ", zap.Error(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := connection.ConnectToDB(ctx, cfg.Postgres.ToDSN())
	if err != nil {
		logger.Fatal("Failed to connect to DB", zap.Error(err))
	}
	defer db.Close()

	if err := migrations.Run(db, logger); err != nil {
		logger.Fatal("Migrations failed", zap.Error(err))
	}
	executor := postgres.NewExecutor(db)
	walletRepo := postgres.NewWalletRepository(executor)
	txManager := txManager.NewTxManager(db)
	walletService := service.NewWallet(walletRepo, txManager)
	h := handler.NewWalletHandler(*logger, walletService)

	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	v1 := router.Group("/api/v1")
	{
		v1.POST("/wallet", h.OperationWithWallet)
		v1.GET("/wallets/:id", h.GetWallets)
	}

	srv := server.NewServer(cfg.HTTP, logger, router)

	srv.Run()

	return nil

}
	