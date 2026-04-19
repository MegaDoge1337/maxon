package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"time"

	"megadoge1337/maxon/internal/handler"
	"megadoge1337/maxon/internal/infrastructure"
	maxonmw "megadoge1337/maxon/internal/middleware"
	"megadoge1337/maxon/internal/repository"
	"megadoge1337/maxon/internal/service"
	"megadoge1337/maxon/internal/usecase/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
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

	// migrate tables
	goose.SetBaseFS(nil)
	goose.SetTableName(environment.GetString("GOOSE_TABLE"))

	if err := goose.SetDialect("postgres"); err != nil {
		slog.Error("failed to set goose dialect", slog.Any("error", err))
		os.Exit(1)
	}

	if err := goose.Up(db, "./migrations"); err != nil {
		slog.Error("failed to run migrations", slog.Any("error", err))
		os.Exit(1)
	}
	slog.Info("migrations completed successfully")

	var userRepo repository.UserRepository = infrastructure.NewSqlBoilerUserRepository(db)
	var authRepo repository.AuthRepository = infrastructure.NewSqlBoilerAuthRepository(db)

	authService := service.NewAuthService(service.AuthSerivceDeps{
		AuthRepo: authRepo,
		UserRepo: userRepo,
		Config: service.AuthServiceConfig{
			JwtSecret: environment.GetString("JWT_SECRET"),
		},
	})

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

	authMiddleware := maxonmw.AuthMiddleware(environment.GetString("JWT_SECRET"))

	userHandler := handler.NewUserHandler(handler.UserHandlerDeps{
		GetAllUC:        user.NewGetAllUsersUseCase(userRepo),
		GetByIdUC:       user.NewGetUserByIdUseCase(userRepo),
		GetByUsernameUC: user.NewGetUserByUsernameUseCase(userRepo),
		CreateUC:        user.NewCreateUserUseCase(userRepo),
		UpdateUC:        user.NewUpdateUserUseCase(userRepo),
		DeleteUC:        user.NewDeleteUserUseCase(userRepo),
		AuthMiddleware:  authMiddleware,
	})
	authHandler := handler.NewAuthHandler(handler.AuthHandlerDeps{AuthService: *authService, AuthMiddleware: authMiddleware})

	// api routes
	router.Route("/api", func(r chi.Router) {
		// v1 routes
		r.Route("/v1", func(r chi.Router) {
			r.Mount("/auth", authHandler.RoutesV1())
			r.Mount("/users", userHandler.RoutesV1())
		})
	})

	// create HTTP server with config settings
	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	slog.Info("starting HTTP server", slog.String("address", cfg.HTTPServer.Address))

	if err := server.ListenAndServe(); err != nil {
		slog.Error("server failed to start", slog.Any("error", err))
		os.Exit(1)
	}
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
