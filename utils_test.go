package gmws

import (
	"encoding/base64"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestAddTimestamp(t *testing.T) {
	originalUrl, err := url.Parse("https://mws.amazonservices.com/som/path")
	if err != nil {
		t.Fatal(err)
	}

	currentTimestamp := time.Now().UTC().Format(time.RFC3339)
	err = addTimestamp(originalUrl)
	if err != nil {
		t.Fatal(err)
	}

	currentTimestampParam := base64.StdEncoding.EncodeToString([]byte("Timestamp=" + currentTimestamp))

	if strings.Contains(originalUrl.String(), currentTimestampParam) {
		t.Fatalf("Expected original url to contain current timestamp but got %s", originalUrl.String())
	}
}
