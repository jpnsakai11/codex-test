package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	delivery "github.com/example/microservices-project/order-service/internal/delivery/http"
	"github.com/example/microservices-project/order-service/internal/infrastructure/client"
	"github.com/example/microservices-project/order-service/internal/infrastructure/config"
	"github.com/example/microservices-project/order-service/internal/infrastructure/database"
	"github.com/example/microservices-project/order-service/internal/infrastructure/logger"
	"github.com/example/microservices-project/order-service/internal/repository/postgres"
	"github.com/example/microservices-project/order-service/internal/usecase"
)

type Container struct {
	Server *http.Server
	Close  func(context.Context) error
}

func Build(ctx context.Context) (*Container, error) {
	cfg := config.Load()
	log, err := logger.New(cfg.LogLevel)
	if err != nil {
		return nil, err
	}
	db, err := database.NewPostgres(ctx, cfg.DBDSN)
	if err != nil {
		return nil, err
	}
	repo := postgres.NewOrderRepository(db)
	userClient := client.NewUserClient(cfg.UserServiceURL)
	uc := usecase.NewOrderUsecase(repo, userClient)
	h := delivery.NewHandler(uc, log, cfg.ServiceName)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.Port),
		Handler:           h.Router(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &Container{
		Server: srv,
		Close: func(_ context.Context) error {
			log.Sync()
			return db.Close()
		},
	}, nil
}
