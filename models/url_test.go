package models_test

import (
	"testing"
	"time"

	"github.com/SergioVenicio/urlShortner/models"
)

func TestURLMetadata(t *testing.T) {
	t.Parallel()

	var hitDate time.Time

	metadata := models.URLMetadata{}
	newHit := models.Hit{
		Time:       hitDate.AddDate(2000, 0, 0),
		Host:       "localhost:8080",
		RequestURI: "/test",
		RemoteAddr: "127.0.0.1:35864",
	}
	metadata.AddHit(&newHit)

	if len(metadata.Hits) != 1 {
		t.Error("expected size 1 got size", len(metadata.Hits))
	}

	if metadata.Hits[0].Time != newHit.Time ||
		metadata.Hits[0].Host != newHit.Host ||
		metadata.Hits[0].RequestURI != newHit.RequestURI ||
		metadata.Hits[0].RemoteAddr != newHit.RemoteAddr {
		t.Errorf("expected %v got %v", metadata.Hits[0], newHit)
	}
}
