package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/example/microservices-project/user-service/internal/app"
)

func main() {
	ctx := context.Background()
	container, err := app.Build(ctx)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := container.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := container.Server.Shutdown(shutdownCtx); err != nil {
		zap.L().Error("server shutdown failed", zap.Error(err))
	}
	_ = container.Close(shutdownCtx)
}
