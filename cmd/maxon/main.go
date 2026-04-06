package main

import (
	"database/sql"
	"log/slog"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
)

const (
	modeProd  = "prod"
	modeDebug = "debug"
)

type Config struct {
	Mode       string     `mapstructure:"mode"`
	Env        string     `mapstructure:"env"`
	HTTPServer HTTPServer `mapstructure:"http_server"`
}

type HTTPServer struct {
	Address     string        `mapstructure:"address"`
	Timeout     time.Duration `mapstructure:"timeout"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
}

func main() {
	// read configs
	viperConfig := viper.New()
	viperConfig.SetConfigType("yaml")
	viperConfig.SetConfigFile("./config/config.yml")
	if err := viperConfig.ReadInConfig(); err != nil {
		slog.Error("failed to read config file", slog.Any("error", err))
		os.Exit(1)
	}

	// unmarshal config into struct
	var cfg Config
	if err := viperConfig.Unmarshal(&cfg); err != nil {
		slog.Error("failed to unmarshal config", slog.Any("error", err))
		os.Exit(1)
	}

	// setup logger based on mode from config
	mode := viperConfig.GetString("mode")
	setupLogger(mode)
	slog.Info("setup logger", slog.String("mode", mode))
	slog.Info("loaded config",
		slog.String("http_address", cfg.HTTPServer.Address),
		slog.Duration("timeout", cfg.HTTPServer.Timeout),
		slog.Duration("idle_timeout", cfg.HTTPServer.IdleTimeout))

	// read environment variables
	environment := viper.New()
	environment.SetConfigType("env")
	environment.SetConfigFile(viperConfig.GetString("env"))
	if err := environment.ReadInConfig(); err != nil {
		slog.Error("failed to read env file", slog.Any("error", err))
		os.Exit(1)
	}

	// open db connection
	db, err := sql.Open("pgx", environment.GetString("DBINFO"))
	if err != nil {
		slog.Error("failed to open database connection", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	// test database connection
	if err := db.Ping(); err != nil {
		slog.Error("failed to ping database", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("database connection established")
}

func setupLogger(mode string) {
	switch mode {
	case modeProd:
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
		slog.SetDefault(logger)
	case modeDebug:
		logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
		slog.SetDefault(logger)
	}
}
