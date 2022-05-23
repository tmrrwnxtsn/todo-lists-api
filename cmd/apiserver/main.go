package main

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/handler"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/server"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/service"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/store/postgres"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	db, err := postgres.NewDB(postgres.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5433"),
		Username: getEnv("DB_USERNAME", "postgres"),
		Password: getEnv("DB_PASSWORD", "qwerty"),
		DBName:   getEnv("DB_NAME", "postgres"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	})
	if err != nil {
		logrus.Fatalf("error occured while connecting to database: %s", err.Error())
	}

	st := postgres.NewStore(db)
	services := service.NewService(st)
	router := handler.NewHandler(services)

	srv := server.NewServer(server.Config{
		BindAddr:       getEnv("BIND_ADDR", ":8080"),
		MaxHeaderBytes: 1,
		ReadTimeout:    10,
		WriteTimeout:   10,
	}, router.InitRoutes())
	go func() {
		if err = srv.Run(); err != nil {
			logrus.Fatalf("error occured while running API server: %s", err.Error())
		}
	}()

	logrus.Println("API server started!")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Println("API shutting down")

	if err = srv.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("error occured while shutting down API server: %s", err.Error())
	}
	if err = db.Close(); err != nil {
		logrus.Fatalf("error occured while closing database connection: %s", err.Error())
	}
}

// getEnv is a simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
