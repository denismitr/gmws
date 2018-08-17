package gmws

import (
	"net/url"
	"time"
)

func addTimestamp(originalURL *url.URL) (err error) {
	values, err := url.ParseQuery(originalURL.RawQuery)
	if err != nil {
		return err
	}

	values.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))
	originalURL.RawQuery = values.Encode()

	return nil
}
