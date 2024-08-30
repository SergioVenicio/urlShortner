package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/SergioVenicio/urlShortner/models"
	"github.com/SergioVenicio/urlShortner/services"
	"github.com/sirupsen/logrus"
)

type URLController struct {
	service *services.URLService
	logger  *logrus.Logger
}

func (c *URLController) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	url, err := c.service.Get(id, r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("url not found"))
		c.logger.Debug("[URLController][GetByID] error:", err)
		return
	}

	http.Redirect(w, r, url.Source, http.StatusPermanentRedirect)
}

func (c *URLController) Add(w http.ResponseWriter, r *http.Request) {
	var u models.URL
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		c.logger.Debug("[URLController][Add] error:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.service.Add(&u)
	if err != nil {
		c.logger.Debug("[URLController] error:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func NewURLController(service *services.URLService, logger *logrus.Logger) *URLController {
	return &URLController{
		service: service,
		logger:  logger,
	}
}
