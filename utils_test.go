package gmws

import (
	"encoding/base64"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestAddTimestamp(t *testing.T) {
	originalURL, err := url.Parse("https://mws.amazonservices.com/som/path")
	if err != nil {
		t.Fatal(err)
	}

	currentTimestamp := time.Now().UTC().Format(time.RFC3339)
	err = addTimestamp(originalURL)
	if err != nil {
		t.Fatal(err)
	}

	currentTimestampParam := base64.StdEncoding.EncodeToString([]byte("Timestamp=" + currentTimestamp))

	if strings.Contains(originalURL.String(), currentTimestampParam) {
		t.Fatalf("Expected original url to contain current timestamp but got %s", originalURL.String())
	}
}
