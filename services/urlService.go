package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/SergioVenicio/urlShortner/models"
	"github.com/SergioVenicio/urlShortner/repositories"
	"github.com/sirupsen/logrus"
)

type URLService struct {
	repository *repositories.URLRepository
	logger     *logrus.Logger
}

func (s *URLService) Add(u *models.URL) error {
	if u.Source == "" {
		return errors.New("invalid url")
	}
	s.logger.Debug("[URLService][Add] received url:", u)
	return s.repository.Add(u)
}

func (s *URLService) Get(id string, r *http.Request) (*models.URL, error) {
	s.logger.Debug(fmt.Sprintf("[URLService][Add] received id:%s", id))
	return s.repository.Get(id, r)
}

func NewURLService(r *repositories.URLRepository, logger *logrus.Logger) *URLService {
	return &URLService{repository: r, logger: logger}
}
