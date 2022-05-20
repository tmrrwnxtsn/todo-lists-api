package main

import (
	"flag"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/apiserver"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/config"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/handler"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/service"
	"github.com/tmrrwnxtsn/todo-lists-api/internal/store"
	"log"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to the config path")
}

func main() {
	flag.Parse()

	cfg := config.NewConfig()
	if err := cfg.Load(configPath); err != nil {
		log.Fatalf("error occured while reading config file: %s", err.Error())
	}

	st := store.NewStore()
	services := service.NewService(st)
	router := handler.NewHandler(services)

	srv := apiserver.NewServer(cfg, router.InitRoutes())
	if err := srv.Run(); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
