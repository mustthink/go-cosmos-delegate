package app

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/mustthink/go-cosmos-delegate/internal/config"
	"github.com/mustthink/go-cosmos-delegate/internal/db"
	"github.com/mustthink/go-cosmos-delegate/internal/parser"
)

type App struct {
	conn   *grpc.ClientConn
	config *config.Config
	parser parser.Parser
	db     db.Storage
	logger *logrus.Logger
}

func New() *App {
	cfg := config.MustLoad()

	log := logrus.New()
	log.SetLevel(cfg.LogLevel())
	log.Debugf("logger successfully initialized")

	creds, err := cfg.Grpc.Credentials()
	if err != nil {
		log.Fatalf("couldn't get grpc credentials: %v", err)
	}

	conn, err := grpc.Dial(cfg.Grpc.Address(), grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("couldn't connect: %v", err)
	}

	cosmosParser := parser.New(conn)
	log.Debugf("cosmos parser successfully initialized")

	storage := db.MustNew(cfg.Database)

	return &App{
		conn:   conn,
		config: cfg,
		parser: cosmosParser,
		db:     storage,
		logger: log,
	}
}

func (a *App) Run(ctx context.Context) {
	ticker := time.NewTicker(a.config.BlockTime)

	for range ticker.C {
		transactions, err := a.parser.Parse(ctx)
		if err != nil {
			a.logger.Errorf("failed to parse transactions: %v", err)
			continue
		}

		for _, t := range transactions {
			if len(t.Messages) == 0 {
				continue
			}

			transactionID, err := a.db.CreateTransaction(ctx, t)
			if err != nil {
				a.logger.Errorf("failed to create transaction: %v", err)
			}
			if transactionID == 0 {
				continue
			}

			err = a.db.CreateDelegateMessages(ctx, transactionID, t.Messages)
			if err != nil {
				a.logger.Errorf("failed to create delegate messages: %v", err)
			}
		}
	}
}
