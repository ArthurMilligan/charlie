package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ArthurMilligan/charlie/api/internal/ai"
	"github.com/ArthurMilligan/charlie/api/internal/app/config"
	apphttp "github.com/ArthurMilligan/charlie/api/internal/app/http"
	"github.com/ArthurMilligan/charlie/api/internal/app/services"
	"github.com/ArthurMilligan/charlie/api/internal/app/storage"
	"github.com/ArthurMilligan/charlie/api/internal/telegram"
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
	notesStorage := storage.NewNotesStorage(db)
	tasksStorage := storage.NewTasksStorage(db)

	messagesService := services.NewMessagesService(messagesStorage)
	notesService := services.NewNotesService(notesStorage)
	tasksService := services.NewTasksService(tasksStorage)

	handler := apphttp.NewHttpServer(messagesService, notesService, tasksService)

	GPTClient := ai.NewGPTClient(conf.GPTToken, notesService, tasksService)

	go func() {
		telegram.ServeTelegramConnection(messagesService, notesService, conf.TelegramBotToken, GPTClient)
	}()

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
