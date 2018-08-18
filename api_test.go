package gmws

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestAPI(t *testing.T) {
	credientials := NewCredentials("fake-access-key", "fake-secret-key", "fake-auth-token", "fake-seller-id")

	t.Run("it initializes API correctly", func(t *testing.T) {
		api := NewAPI(credientials, "fake-marketplace-id")

		if api.Host != "mws.amazonservices.com" {
			t.Errorf("expected api host to be %s, but got %v", "mws.amazonservices.com", api.Host)
		}

		if api.Scheme != "https" {
			t.Errorf("expected api scheme to be %s, but got %v", "https", api.Host)
		}

		if api.Client == nil {
			t.Error("client should have been initialized")
		}
	})

	t.Run("making get request to fake endpoint", func(t *testing.T) {
		api := NewAPI(credientials, "fake-marketplace-id")

		client := NewTestClient(func(req *http.Request) *http.Response {
			// Test request parameters
			if !strings.HasPrefix(req.URL.String(), "https://mws.amazonservices.com/fake/get/request?") {
				t.Errorf("url is incorrect: %s", req.URL.String())
			}

			if !strings.Contains(req.URL.String(), "fake-key=fake-value") {
				t.Errorf("url %s does not contain param: %s", req.URL.String(), "fake-key=fake-value")
			}

			return &http.Response{
				StatusCode: 200,
				// Send response to be tested
				Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		})

		api.Client = client

		fr := &fakeGetRequest{
			Method:      http.MethodGet,
			Action:      "FakeGetRequest",
			ActionPath:  "/fake/get/request",
			QueryParams: map[string]string{"fake-key": "fake-value"},
		}

		response, err := api.Call(fr)

		if err != nil {
			t.Fatal(err)
		}

		if response != "OK" {
			t.Errorf("Expected response to contain 'OK' but got %s", response)
		}
	})
}

func assertEquals(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Failed asserting that two values are equal. Expected %v got %v", expected, actual)
	}
}
