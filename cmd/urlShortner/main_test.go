package main_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/SergioVenicio/urlShortner/controllers"
	"github.com/SergioVenicio/urlShortner/models"
	"github.com/SergioVenicio/urlShortner/repositories"
	"github.com/SergioVenicio/urlShortner/services"
	log "github.com/sirupsen/logrus"
)

func TestURL(t *testing.T) {
	t.Parallel()

	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.DebugLevel)

	tt := []struct {
		body string
		name string
	}{
		{
			name: "with google.com",
			body: `{"source": "https://google.com"}`,
		},
		{
			name: "with youtube.com",
			body: `{"source": "https://youtube.com"}`,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/", strings.NewReader(tc.body))
			responseRecorder := httptest.NewRecorder()

			urlRepository := repositories.NewURLRepository(logger)
			urlService := services.NewURLService(urlRepository, logger)
			urlController := controllers.NewURLController(urlService, logger)

			urlController.Add(responseRecorder, request)

			if responseRecorder.Code != http.StatusCreated {
				t.Fatalf("expected %v got %v", http.StatusCreated, responseRecorder.Code)
			}

			responseValue := strings.TrimSpace(responseRecorder.Body.String())
			fmt.Println(responseRecorder.Body.String(), responseValue)

			var received models.URL
			err := json.Unmarshal(responseRecorder.Body.Bytes(), &received)
			if err != nil {
				t.Fatal("invalid response payload, err:", err)
			}

			if received.ID == "" {
				t.Fatal("empty URL id")
			}

			type expectedValue struct {
				Source string
			}

			var expected expectedValue
			json.Unmarshal([]byte(tc.body), &expected)
			if received.Source != expected.Source {
				t.Fatalf("invalid url source expected: %v received:%v", received.Source, expected.Source)
			}
		})
	}
}

func TestURLWithInvalidValues(t *testing.T) {
	t.Parallel()

	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.DebugLevel)

	tt := []struct {
		body string
		name string
	}{
		{
			name: "without value on source",
			body: `{"source": ""}`,
		},
		{
			name: "without payload",
			body: `{}`,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest("POST", "/", strings.NewReader(tc.body))
			responseRecorder := httptest.NewRecorder()

			urlRepository := repositories.NewURLRepository(logger)
			urlService := services.NewURLService(urlRepository, logger)
			urlController := controllers.NewURLController(urlService, logger)

			urlController.Add(responseRecorder, request)
			if responseRecorder.Code != http.StatusBadRequest {
				t.Fatalf("expected %v got %v", http.StatusBadRequest, responseRecorder.Code)
			}
		})
	}
}
