package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kapitan123/telegrofler/config"
	"github.com/kapitan123/telegrofler/internal/bot"
	"github.com/kapitan123/telegrofler/internal/data/post"
	"github.com/kapitan123/telegrofler/internal/routes"

	log "github.com/sirupsen/logrus"
)

// Main entry point. Starts HTTP service
func main() {
	log.Info("Telegrofler: starting...")

	router := mux.NewRouter()

	api := routes.API{
		Bot:        bot.New(),
		PostsStore: post.NewStore(),
	}

	// AK TODO should close the whole api
	defer api.Close()

	api.AddRoutes(router)
	api.AddCallback(router)

	log.Info("Telegrofler: listening on: ", config.ServerPort)

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", config.ServerPort),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}
