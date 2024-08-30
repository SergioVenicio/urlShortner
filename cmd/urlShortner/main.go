package main

import (
	"net/http"
	"os"

	"github.com/SergioVenicio/urlShortner/controllers"
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
	urlController := controllers.NewURLController(urlService, logger)

	srv := http.NewServeMux()
	srv.HandleFunc("GET /{id}", urlController.GetByID)
	srv.HandleFunc("POST /", urlController.Add)

	logger.Debug("starting http server on :8080")
	logger.Fatal(http.ListenAndServe(":8080", srv))
}
