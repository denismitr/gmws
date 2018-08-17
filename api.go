package gmws

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// API - Amazon MWS API
type API struct {
	credentials   Credentials
	host          string
	marketplaceID string
	client        *http.Client
}

func (api *API) SetClient(client *http.Client) {
	api.client = client
}

// Call makes Amazon MWS API call
func (api *API) Call(mwsRequest MWSRequest) (string, error) {
	url, err := mwsRequest.GenerateURL(api.host)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate request query")
	}

	r, err := http.NewRequest(mwsRequest.Method(), url, mwsRequest.Body())
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to create request to %s", url))
	}

	if mwsRequest.Method() == http.MethodPost {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	res, err := api.client.Do(r)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("%s request to %s failed", mwsRequest.Method(), url))
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "Failed to read response body")
	}

	return string(body), nil
}

// NewAPI - creates new Amazon MWS API
func NewAPI(credentials Credentials, marketplaceID string) *API {
	return &API{
		credentials:   credentials,
		marketplaceID: marketplaceID,
		host:          "mws.amazonservices.com",
		client:        &http.Client{},
	}
}
