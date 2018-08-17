package gmws

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"reflect"
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

	t.Run("api creation", func(t *testing.T) {
		api := NewAPI(credientials, "fake-marketplace-id")

		if api.host != "mws.amazonservices.com" {
			t.Errorf("Expected api host to be %s, but got %v", "mws.amazonservices.com", api.host)
		}
	})

	t.Run("making get request to fake endpoint", func(t *testing.T) {
		api := NewAPI(credientials, "fake-marketplace-id")

		client := NewTestClient(func(req *http.Request) *http.Response {
			// Test request parameters
			assertEquals(t, req.URL.String(), "https://mws.amazonservices.com/some/path")
			return &http.Response{
				StatusCode: 200,
				// Send response to be tested
				Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		})

		api.SetClient(client)

		fr := &fakeGetRequest{url: "/some/path", httpMethod: http.MethodGet}

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
