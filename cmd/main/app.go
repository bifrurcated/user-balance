package main

import (
	"context"
	"fmt"
	"github.com/bifrurcated/user-balance/internal/balance"
	"github.com/bifrurcated/user-balance/internal/balance/db"
	"github.com/bifrurcated/user-balance/internal/config"
	"github.com/bifrurcated/user-balance/internal/history"
	historydb "github.com/bifrurcated/user-balance/internal/history/db"
	"github.com/bifrurcated/user-balance/internal/reserve"
	"github.com/bifrurcated/user-balance/internal/reserve/db"
	"github.com/bifrurcated/user-balance/pkg/client/postgresql"
	"github.com/bifrurcated/user-balance/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	client, err := postgresql.NewClient(context.TODO(), cfg.Storage)
	if err != nil {
		logger.Fatal(err)
	}

	err = executeSQLScript(context.TODO(), client)
	if err != nil {
		logger.Fatal(err)
	}

	historyRepository := historydb.NewRepository(client, logger)
	historyService := history.NewService(historyRepository, logger)
	historyHandler := history.NewHandler(historyService, logger)
	historyHandler.Register(router)

	balanceRepository := balancedb.NewRepository(client, logger)
	balanceService := balance.NewService(balanceRepository, historyRepository, logger)
	balanceHandler := balance.NewHandler(balanceService, logger)
	balanceHandler.Register(router)

	reserveRepository := reservedb.NewRepository(client, logger)
	reserveService := reserve.NewService(reserveRepository, balanceRepository, historyRepository, logger)
	reserveHandler := reserve.NewHandler(reserveService, logger)
	reserveHandler.Register(router)

	start(router, cfg)
}

func executeSQLScript(ctx context.Context, client postgresql.Client) error {
	logger := logging.GetLogger()
	path := filepath.Join("data.sql")

	logger.Info("read sql script file")
	c, ioErr := os.ReadFile(path)
	if ioErr != nil {
		return ioErr
	}
	sql := string(c)
	logger.Info("execute sql script")
	res, err := client.Exec(ctx, sql)
	if err != nil {
		return err
	}
	logger.Debug(res)
	return nil
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	logger.Info("listen tcp")
	listener, listenErr := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Infof("server is listening address to %s and network to %s", listener.Addr().String(), listener.Addr().Network())

	logger.Fatal(server.Serve(listener))
}
