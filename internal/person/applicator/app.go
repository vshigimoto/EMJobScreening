package applicator

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-jun/internal/person/config"
	"go-jun/internal/person/database"
	"go-jun/internal/person/repository"
	"go-jun/internal/person/server/http"
	"go-jun/internal/person/usecase"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

type Applicator struct {
	logger *zap.SugaredLogger
	config config.Config
}

func New(cfg config.Config, logger *zap.SugaredLogger) *Applicator {
	return &Applicator{
		config: cfg,
		logger: logger,
	}
}

func (a *Applicator) Run() {
	r := gin.Default()
	cfg := a.config
	l := a.logger
	ctx, cancel := context.WithCancel(context.TODO())
	_ = ctx

	mainDB, err := database.New(cfg.Database.Main)
	if err != nil {
		l.Fatalf("Error to connect DB '%s':%v", cfg.Database.Main.Host, err)
	}
	defer func() {
		if err := mainDB.Close(); err != nil {
			l.Panicf("failed close mainDB err: %v", err)
		}
		l.Info("mainDB closed")
	}()

	replicaDB, err := database.New(cfg.Database.Replica)
	if err != nil {
		l.Fatalf("Error to connect DB '%s':%v", cfg.Database.Replica.Host, err)
	}
	defer func() {
		if err := replicaDB.Close(); err != nil {
			l.Panicf("failed close replicaDB err: %v", err)
		}
		l.Info("replicaDB closed")
	}()

	rep := repository.New(mainDB, replicaDB)

	userUC := usecase.New(l, rep)
	http.UserRouter(r, *userUC, l, cfg)
	port := fmt.Sprintf(":%d", cfg.HttpServer.Port)

	l.Infof("User server on port %s, running", port)
	if err := r.Run(port); err != nil {
		l.Panicf("Error to run server %v", err)
	}

	a.gracefulShutdown(cancel)
}

func (a *Applicator) gracefulShutdown(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch
	signal.Stop(ch)
	cancel()
}
