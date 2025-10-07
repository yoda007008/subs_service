package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"subs/subservice/cmd/migrator"
	"subs/subservice/internal/config"
	connect "subs/subservice/internal/db"
	"subs/subservice/internal/router"
	"syscall"
	"time"
)

func main() {
	// todo init logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// todo init config
	configPath := flag.String("config", "/home/kirill/GolandProjects/subs/subservice/internal/config/config.yaml", "path to config file")

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		logger.Error("error load config path", "err", "path not found")
		os.Exit(1)
	}

	// todo init base urls and paths
	DATABASE_URL := cfg.DatabaseConfig.Url
	MIGRATIONS_PATH := cfg.MigrationsPath.Path
	//fmt.Println(DATABASE_URL)
	//fmt.Println(MIGRATIONS_PATH)

	// todo run migrations
	if err := migrator.RunMigrations(DATABASE_URL, MIGRATIONS_PATH); err != nil { // run migrations
		logger.Error("Migraions is not access", "error", err)
		os.Exit(1)
	}

	// todo connect database
	logger.Info("connect to database")
	db, err := connect.NewPostgres(DATABASE_URL)
	if err != nil {
		logger.Info("connect db")
		os.Exit(1)
	}
	defer db.Close()

	// todo init handlers
	router := router.Router(&cfg, db)

	// todo run service with graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := router.ListenAndServe(); err != nil {
			if err.Error() != "http: Server closed" {
				logger.Error("don't run this http server")
			}
		}
	}()

	logger.Info("server started on %s", cfg.SubServiceConfig.Port)

	<-ctx.Done()

	logger.Info("shutting down")
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := router.Shutdown(ctxShutdown); err != nil {
		logger.Info("shutdown")
	}

	logger.Info("graceful shutdown complete")
}
