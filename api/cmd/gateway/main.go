package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ArthurMilligan/charlie/api/internal/app/config"
	apphttp "github.com/ArthurMilligan/charlie/api/internal/app/http"
	"github.com/ArthurMilligan/charlie/api/internal/app/services"
	"github.com/ArthurMilligan/charlie/api/internal/app/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	conf := config.MustLoad()

	ctx := context.Background()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		conf.DB.User,
		conf.DB.Password,
		conf.DB.Host,
		conf.DB.Port,
		conf.DB.DBName,
	)

	db, err := pgxpool.New(ctx, dsn)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		panic(err)
	}

	messagesStorage := storage.NewMessagesStorage(db)

	messagesService := services.NewMessagesService(messagesStorage)

	handler := apphttp.NewHttpServer(messagesService)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", conf.APIConfig.Port),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       conf.APIConfig.Timeout,
		WriteTimeout:      conf.APIConfig.Timeout,
		IdleTimeout:       60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
