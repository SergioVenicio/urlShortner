package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/SergioVenicio/urlShortner/models"
	"github.com/SergioVenicio/urlShortner/repositories"
	"github.com/SergioVenicio/urlShortner/services"
	log "github.com/sirupsen/logrus"
)

var logger *log.Logger

func init() {
	logger = log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.DebugLevel)
}

func main() {
	urlRepository := repositories.NewURLRepository(logger)
	urlService := services.NewURLService(urlRepository, logger)

	srv := http.NewServeMux()
	srv.HandleFunc("GET /{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		url, err := urlService.Get(id, r)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("url not found"))
			logger.Debug("[HttpServer][GET][/] error:", err)
			return
		}

		http.Redirect(w, r, url.Source, http.StatusPermanentRedirect)
	})
	srv.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		var u models.URL
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			logger.Debug("[HttpServer][POST][/] error:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = urlService.Add(&u)
		if err != nil {
			logger.Debug("[HttpServer][POST][/] error:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
	})

	logger.Debug("starting http server on :8080")
	logger.Fatal(http.ListenAndServe(":8080", srv))
}
