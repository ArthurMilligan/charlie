package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/effector-fullstack/api-gateway/internal/app/config"
	apphttp "github.com/effector-fullstack/api-gateway/internal/app/http"
)

func main() {
	conf := config.MustLoad()

	handler := apphttp.NewHttpServer()

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
