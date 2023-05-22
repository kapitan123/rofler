package main

import (
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func main() {
	router := mux.NewRouter()

	setupRouter(router)

	srv := &http.Server{
		Handler: router,
		Addr:    "8080",
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.WithError(err).Fatal("Failed to start server")
	}
}
