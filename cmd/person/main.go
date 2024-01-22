package main

import (
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/lib/pq"
	"go-jun/internal/person/applicator"
	"go-jun/internal/person/config"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	l := logger.Sugar()
	l = l.With(zap.String("app", "go-jun"))
	var cfg config.Config
	err := cleanenv.ReadConfig("config/person/config.yaml", &cfg)
	if err != nil {
		l.Fatalf("Failed to load config: %v", err)
	}
	app := applicator.New(cfg, l)
	app.Run()
}
