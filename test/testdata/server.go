package testdata

import (
	"context"
	"errors"
	"fmt"
	"github.com/bifrurcated/user-balance/internal/balance"
	balancedb "github.com/bifrurcated/user-balance/internal/balance/db"
	"github.com/bifrurcated/user-balance/internal/config"
	historydb "github.com/bifrurcated/user-balance/internal/history/db"
	"github.com/bifrurcated/user-balance/pkg/client/postgresql"
	"github.com/bifrurcated/user-balance/pkg/logging"
	"github.com/jackc/pgx/v5"
	"github.com/julienschmidt/httprouter"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync"
)

type Server struct {
	Test  *httptest.Server
	Store postgresql.Client
}

var one sync.Once
var server *Server

func GetTestServer() *Server {
	one.Do(func() {
		logger := logging.GetLogger()
		router := httprouter.New()

		storageConfig := config.StorageConfig{
			Host:     "localhost",
			Port:     "5432",
			Database: "test_db",
			Username: "postgres",
			Password: "123",
		}
		createTestDB(context.TODO(), storageConfig)

		client, err := postgresql.NewClient(context.TODO(), storageConfig)
		if err != nil {
			logger.Fatal(err)
		}

		err = ExecuteSQLScript(context.TODO(), client, "create.sql")
		if err != nil {
			logger.Fatal(err)
		}

		historyRepository := historydb.NewRepository(client, logger)
		balanceRepository := balancedb.NewRepository(client, logger)
		service := balance.NewService(balanceRepository, historyRepository, logger)
		handler := balance.NewHandler(service, logger)
		handler.Register(router)

		server = &Server{
			Test:  httptest.NewServer(router),
			Store: client,
		}
	})

	return server
}

func ExecuteSQLScript(ctx context.Context, client postgresql.Client, fileName string) error {
	logger := logging.GetLogger()
	path := filepath.Join("../../test/testdata/" + fileName)

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

func createTestDB(ctx context.Context, sc config.StorageConfig) {
	logger := logging.GetLogger()
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/",
		sc.Username,
		sc.Password,
		sc.Host,
		sc.Port)
	db, err := pgx.Connect(ctx, dsn)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close(ctx)
	var dbName string
	err = db.QueryRow(ctx, `SELECT datname FROM pg_database WHERE datname = $1`, sc.Database).Scan(&dbName)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			logger.Fatal(err)
		}
		logger.Info("create test_db because not exist")
		q := fmt.Sprintf("CREATE DATABASE %s WITH OWNER = %s ENCODING = 'UTF8' CONNECTION LIMIT = -1", sc.Database, sc.Username)
		_, err = db.Exec(ctx, q)
		if err != nil {
			logger.Fatal(err)
		}
	} else {
		logger.Infof("database (%s) exist", dbName)
	}
}
