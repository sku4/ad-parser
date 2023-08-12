package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/sku4/ad-parser/configs"
	"github.com/sku4/ad-parser/internal/repository"
	"github.com/sku4/ad-parser/internal/service"
	"github.com/sku4/ad-parser/pkg/logger"
	"github.com/tarantool/go-tarantool/v2"
	"github.com/tarantool/go-tarantool/v2/pool"
)

func main() {
	// init config
	log := logger.Get()
	cfg, err := configs.Init()
	if err != nil {
		log.Fatalf("error init config: %s", err)
	}

	// init tarantool
	conn, err := pool.Connect(cfg.Tarantool.Servers, tarantool.Opts{
		Timeout:   cfg.Tarantool.Timeout,
		Reconnect: cfg.Tarantool.ReconnectInterval,
	})
	if err != nil {
		log.Fatalf("error tarantool connection refused: %s", err)
	}
	defer func() {
		errs := conn.Close()
		for _, e := range errs {
			log.Errorf("error close connection pool: %s", e)
		}
	}()

	// init context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	ctx = configs.Set(ctx, cfg)

	repos := repository.NewRepository(conn)
	services := service.NewService(repos)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	log.Infof("App Started")

	go func() {
		if err = services.Parser.Run(ctx); err != nil {
			log.Errorf("error parser run: %s", err)
		}
	}()

	// graceful shutdown
	log.Infof("Got signal %v, attempting graceful shutdown", <-quit)
	cancel()
	log.Info("Context is stopped")

	err = services.Parser.Shutdown()
	if err != nil {
		log.Errorf("error parser shutdown: %s", err)
	} else {
		log.Info("Parser stopped")
	}

	errs := conn.CloseGraceful()
	for _, e := range errs {
		log.Errorf("error close graceful connection pool: %s", e)
	}

	log.Info("App Shutting Down")
}
