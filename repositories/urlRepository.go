package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SergioVenicio/urlShortner/models"
	"github.com/aidarkhanov/nanoid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type URLRepository struct {
	rdb    *redis.Client
	ctx    context.Context
	logger *logrus.Logger
}

func (r *URLRepository) buildKey(pattern string, id string) string {
	return fmt.Sprintf("%s:%s", pattern, id)
}

func (r *URLRepository) Add(u *models.URL) error {
	if u.ID == "" {
		u.ID = nanoid.New()
	}

	jsonData, err := json.Marshal(u)
	if err != nil {
		r.logger.Warn("[URLRepository][Add] error on unmarshal redis data err:", err)
		return err
	}
	err = r.rdb.Set(r.ctx, r.buildKey("URLS", u.ID), jsonData, 0).Err()
	if err != nil {
		r.logger.Warn("[URLRepository][Add] error on set redis data err:", err)
		return err
	}
	return nil
}

func (r *URLRepository) GetMedatada(url string) (models.URLMetadata, error) {
	m := models.URLMetadata{URL: url}
	jsonValue, err := r.rdb.Get(r.ctx, r.buildKey("URL_MEDATADA", url)).Result()
	if err != nil {
		r.logger.Warn("[URLRepository][GetMedatada] error on get redis data err:", err)
		return m, err
	}

	err = json.Unmarshal([]byte(jsonValue), &m)
	if err != nil {
		r.logger.Warn("[URLRepository][GetMedatada] error on unmarshal redis data err:", err)
		return m, err
	}

	return m, nil
}

func (r *URLRepository) AddMetadata(m *models.URLMetadata) error {
	jsonData, err := json.Marshal(m)
	if err != nil {
		r.logger.Warn("[URLRepository][AddMetadata] error on unmarshal redis data err:", err)
		return err
	}
	err = r.rdb.Set(r.ctx, r.buildKey("URL_MEDATADA", m.URL), jsonData, 0).Err()
	if err != nil {
		r.logger.Warn("[URLRepository][AddMetadata] error on set redis data err:", err)
		return err
	}
	return nil
}

func (r *URLRepository) Get(id string, request *http.Request) (*models.URL, error) {
	jsonValue, err := r.rdb.Get(r.ctx, r.buildKey("URLS", id)).Result()
	if err != nil {
		r.logger.Warn("[URLRepository][Get] error on get redis data err:", err)
		return nil, err
	}

	var u models.URL
	err = json.Unmarshal([]byte(jsonValue), &u)
	if err != nil {
		r.logger.Warn("[URLRepository][Get] error on unmarshal redis data err:", err)
		return nil, err
	}

	metadata, _ := r.GetMedatada(u.ID)
	metadata.AddHit(&models.Hit{
		Time:       time.Now(),
		RequestURI: request.RequestURI,
		RemoteAddr: request.RemoteAddr,
		Host:       request.Host,
		Header:     request.Header,
		Method:     request.Method,
	})
	r.AddMetadata(&metadata)
	return &u, nil
}

func NewURLRepository(logger *logrus.Logger) *URLRepository {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	return &URLRepository{rdb: rdb, ctx: ctx, logger: logger}
}
