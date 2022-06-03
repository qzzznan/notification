package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"notification/config"
	v1 "notification/internal/controller/http/v1"
	"notification/internal/usecase"
	"notification/internal/usecase/repo/bark"
	"notification/internal/usecase/repo/pushdeer"
	"notification/internal/usecase/webapi"
	"notification/pkg/httpserver"
	"notification/pkg/logger"
	"notification/pkg/postgres"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	// **** log
	l, err := logger.New(cfg.Log.Level)
	if err != nil {
		log.Fatalln(err)
	}
	// **** log
	// *** cache
	/*
		r := rdb.New(&redis.Options{
			Addr:     cfg.RDB.Addr,
			Password: cfg.RDB.Password,
		}, rdb.WithKeyPrefix("app.notify"))
		err = r.Cli.Ping(context.Background()).Err()
		if err != nil {
			l.Fatalln(err)
		}
	*/
	// *** cache

	// **** db
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable",
		cfg.PG.User, cfg.PG.Password, cfg.PG.URL, cfg.PG.Port)
	pg, err := postgres.New(dsn) //postgres.WithCache(r))
	if err != nil {
		l.Fatalln(err)
	}
	defer pg.Close()
	// **** db

	barkApns, err := webapi.NewBarkAPNs(l)
	if err != nil {
		l.Fatalln(err)
	}
	pdApns, err := webapi.NewPushDeerAPNs(l)
	if err != nil {
		l.Fatalln(err)
	}

	// **** service
	barkUseCase := usecase.NewBark(bark.New(pg, l), barkApns)
	pushDeerUseCase := usecase.NewPushDeer(pushdeer.New(pg, l), pdApns)
	// **** service

	e := gin.New()
	v1.NewRouter(e, l, &v1.UseCase{
		B: barkUseCase,
		P: pushDeerUseCase,
	})
	ctx := context.WithValue(context.Background(), "logger", l)
	httpServer := httpserver.New(e,
		httpserver.WithContext(ctx),
		httpserver.WithAddr(cfg.Http.Addr),
	)

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-sig:
		l.Warningln("signal received:", s)

	case s := <-httpServer.Notify():
		l.Warningln("server notify:", s)
	}

	err = httpServer.Shutdown()
	if err != nil {
		l.Fatalln(err)
	}
}
