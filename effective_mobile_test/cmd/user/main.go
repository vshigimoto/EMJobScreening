package main

import (
	"effective_mobile_test/internal/user/config"
	"effective_mobile_test/internal/user/database"
	"effective_mobile_test/internal/user/repository"
	"effective_mobile_test/internal/user/server/http"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" //
	"github.com/spf13/viper"
)

func main() {
	r := gin.Default()
	cfg, err := loadConfig("config/user")
	mainDb, err := database.New(cfg.Database.Main)
	replicaDb, err := database.New(cfg.Database.Replica)
	rep := repository.NewRepository(mainDb, replicaDb)
	http.InitRouter(r, rep)
	port := fmt.Sprintf(":%d", cfg.HttpServer.Port)
	err = r.Run(port)
	if err != nil {
		return
	}
}

func loadConfig(path string) (config config.Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return config, fmt.Errorf("failed to Read config err: %w", err)
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("failed to Unmarshal config err: %w", err)
	}
	return config, nil
}
