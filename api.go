package gmws

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

// API - Amazon MWS API
type API struct {
	Credentials   Credentials
	Host          string
	Scheme        string
	MarketplaceID string
	Client        *http.Client
}

// Call makes Amazon MWS API call
func (api *API) Call(mwsRequest MWSRequest) (string, error) {
	requestParams := mwsRequest.GetParams()

	urlBuilder, err := api.GenerateGetURL(requestParams)
	if err != nil {
		return "", err
	}

	urlString := urlBuilder.String()

	r, err := http.NewRequest(requestParams.Method, urlString, nil)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to create request to %s", urlString))
	}

	if requestParams.Method == http.MethodPost {
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	res, err := api.Client.Do(r)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("%s request to %s failed", requestParams.Method, urlString))
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "Failed to read response body")
	}

	return string(body), nil
}

func (api *API) GenerateGetURL(requestParams *RequestParams) (*url.URL, error) {
	urlBuilder, err := url.Parse(requestParams.ActionPath)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create valid params for given request action")
	}

	urlBuilder.Scheme = api.Scheme
	urlBuilder.Host = api.Host

	query := url.Values{}

	query.Add("Action", requestParams.Action)

	if api.Credentials.AuthToken != "" {
		query.Add("MWSAuthToken", api.Credentials.AuthToken)
	}

	query.Add("AWSAccessKeyId", api.Credentials.AccessKey)
	query.Add("SellerId", api.Credentials.SellerID)
	query.Add("SignatureVersion", "2")
	query.Add("SignatureMethod", "HmacSHA256")
	query.Add("Version", "2011-10-01")

	for k, v := range requestParams.QueryParams {
		query.Set(k, v)
	}

	urlBuilder.RawQuery = query.Encode()

	return urlBuilder, nil
}

func (api *API) sign(method string, url *url.URL) (string, error) {
	escapedURL := strings.Replace(url.RawQuery, ",", "%2C", -1)
	escapedURL = strings.Replace(escapedURL, ":", "%3A", -1)

	params := strings.Split(escapedURL, "&")
	paramsMap := make(map[string]string)
	keys := make([]string, len(params))

	for k, v := range params {
		pair := strings.Split(v, "=")
		paramsMap[pair[0]] = pair[1]
		keys[k] = pair[0]
	}

	sort.Strings(keys)

	sortedParams := make([]string, len(params))

	for k, _ := range params {
		var buffer bytes.Buffer
		buffer.WriteString(keys[k])
		buffer.WriteString("=")
		buffer.WriteString(paramsMap[keys[k]])
		sortedParams[k] = buffer.String()
	}

	paramString := strings.Join(sortedParams, "&")

	urlToSign := fmt.Sprintf("%s\n%s\n%s\n%s", method, url.Host, url.Path, paramString)

	hasher := hmac.New(sha256.New, []byte(api.Credentials.SecretKey))
	_, err := hasher.Write([]byte(urlToSign))
	if err != nil {
		return "", errors.Wrap(err, "could not write hash to sign url")
	}

	hash := base64.StdEncoding.EncodeToString(hasher.Sum(nil))

	return hash, nil
}

// NewAPI - creates new Amazon MWS API
func NewAPI(credentials Credentials, marketplaceID string) *API {
	return &API{
		Credentials:   credentials,
		MarketplaceID: marketplaceID,
		Host:          "mws.amazonservices.com",
		Scheme:        "https",
		Client:        &http.Client{},
	}
}
