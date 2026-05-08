package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"megadoge1337/maxon/internal/config"
	"megadoge1337/maxon/internal/handler"
	"megadoge1337/maxon/internal/infrastructure"
	maxonmw "megadoge1337/maxon/internal/middleware"
	"megadoge1337/maxon/internal/repository"
	"megadoge1337/maxon/internal/usecase/auth"
	"megadoge1337/maxon/internal/usecase/user"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
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
	var cfg config.Config

	err := config.LoadConfig(&cfg)
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	setupLogger(cfg.Mode)
	slog.Info("setup logger", slog.String("mode", cfg.Mode))

	driver := cfg.Database.Driver
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SslMode)

	db, err := infrastructure.NewPostgresConnection(driver, dsn)
	if err != nil {
		slog.Error("failed to create new postgres connection", slog.Any("error", err))
		os.Exit(1)
	}
	defer db.Close()

	// migrate tables
	goose.SetBaseFS(nil)
	goose.SetTableName(cfg.Database.MigrationsTable)

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
	var roleRepo repository.RoleRepository = infrastructure.NewSqlBoilerRoleRepository(db)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

	authMiddleware := maxonmw.AuthMiddleware(cfg.Auth.JwtSecret)
	adminMiddleware := maxonmw.AdminMiddleware()

	userHandler := handler.NewUserHandler(handler.UserHandlerDeps{
		GetAllUC:        user.NewGetAllUsersUseCase(userRepo),
		GetByIdUC:       user.NewGetUserByIdUseCase(userRepo),
		GetByUsernameUC: user.NewGetUserByUsernameUseCase(userRepo),
		CreateUC:        user.NewCreateUserUseCase(userRepo),
		UpdateUC:        user.NewUpdateUserUseCase(userRepo),
		DeleteUC:        user.NewDeleteUserUseCase(userRepo),
		AuthMiddleware:  authMiddleware,
		AdminMiddleware: adminMiddleware,
	})
	authHandler := handler.NewAuthHandler(handler.AuthHandlerDeps{
		LoginUC: auth.NewLoginUseCase(
			auth.LoginUseCaseDeps{
				AuthRepo: authRepo,
				UserRepo: userRepo,
				RoleRepo: roleRepo,
				Config: auth.LoginUseCaseConfig{
					JwtSecret: cfg.Auth.JwtSecret,
				},
			},
		),
		RegisterUC: auth.NewRegisterUseCase(
			auth.RegisterUseCaseDeps{
				AuthRepo: authRepo,
				UserRepo: userRepo,
				RoleRepo: roleRepo,
			},
		),
		RefreshUC: auth.NewRefreshUseCase(
			auth.RefreshUseCaseDeps{
				Repo: authRepo,
				Config: auth.RefreshUseCaseConfig{
					JwtSecret: cfg.Auth.JwtSecret,
				},
			},
		),
		AuthMiddleware: authMiddleware,
	})

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
