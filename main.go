package main

import (
	"context"
	"fmt"
	"github.com/core-go/config"
	"github.com/core-go/core"
	mid "github.com/core-go/log/middleware"
	"github.com/core-go/log/zap"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"go-service/internal/app"
)

func main() {
	var cfg app.Config
	err := config.Load(&cfg, "configs/config")
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()

	log.Initialize(cfg.Log)
	r.Use(mid.BuildContext)
	logger := mid.NewLogger()
	if log.IsInfoEnable() {
		r.Use(mid.Logger(cfg.MiddleWare, log.InfoFields, logger))
	}
	r.Use(mid.Recover(log.PanicMsg))

	err = app.Route(r, context.Background(), cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println(core.ServerInfo(cfg.Server))
	server := core.CreateServer(cfg.Server, r)
	if err = server.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
	}
}
