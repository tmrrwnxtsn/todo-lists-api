package main

import (
	"context"
	"flag"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/config"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/handler"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/server"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/service"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/store/postgres"
	"os"
	"os/signal"
	"syscall"
)

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

func main() {
	flag.Parse()

	logrus.SetFormatter(&logrus.JSONFormatter{})

	// load application configurations
	cfg, err := config.Load(*flagConfig)
	if err != nil {
		logrus.Fatalf("error occured while loading application configuration: %s", err.Error())
	}

	db, err := postgres.NewDB(cfg.DSN)
	if err != nil {
		logrus.Fatalf("error occured while connecting to database: %s", err.Error())
	}

	st := postgres.NewStore(db)
	services := service.NewService(st)
	router := handler.NewHandler(services)

	srv := server.NewServer(cfg.BindAddr, router.InitRoutes())
	go func() {
		if err = srv.Run(); err != nil {
			logrus.Fatalf("error occured while running API server: %s", err.Error())
		}
	}()

	logrus.Printf("API server is running at %v", cfg.BindAddr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("API server shutting down")

	if err = srv.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("error occured while shutting down API server: %s", err.Error())
	}
	if err = db.Close(); err != nil {
		logrus.Fatalf("error occured while closing database connection: %s", err.Error())
	}
}
