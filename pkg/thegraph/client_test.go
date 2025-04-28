package thegraph

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
)

// MockRoundTripper to mock HTTP responses
type MockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

// Helper to create a mocked client
func newMockedClient(fn func(req *http.Request) (*http.Response, error)) *Client {
	return &Client{
		baseURL: "http://mocked",
		httpClient: &http.Client{
			Timeout:   defaulHttpTimeout,
			Transport: &MockRoundTripper{RoundTripFunc: fn},
		},
	}
}

func TestNewDefaultClient(t *testing.T) {
	client := NewDefaultClient()
	if client.baseURL != DEFAULT_URL {
		t.Errorf("expected baseURL %s, got %s", DEFAULT_URL, client.baseURL)
	}
	if client.httpClient.Timeout != defaulHttpTimeout {
		t.Errorf("expected timeout %s, got %s", defaulHttpTimeout, client.httpClient.Timeout)
	}
}

func TestNewClient(t *testing.T) {
	url := "http://example.com"
	client := NewClient(url)
	if client.baseURL != url {
		t.Errorf("expected baseURL %s, got %s", url, client.baseURL)
	}
}

func TestGetVouchersSuccess(t *testing.T) {
	mockResponse := `{
		"data": {
			"vouchers": [
				{
					"id": "1",
					"expiration": "123456",
					"value": "100",
					"balance": "50",
					"owner": {"id": "owner1"},
					"voucherType": {"id": "type1", "description": "desc1"}
				}
			]
		}
	}`

	client := newMockedClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockResponse)),
			Header:     make(http.Header),
		}, nil
	})

	vouchers, err := client.GetVouchers()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(vouchers.Data.Vouchers) != 1 {
		t.Fatalf("expected 1 voucher, got %d", len(vouchers.Data.Vouchers))
	}
	if vouchers.Data.Vouchers[0].ID != "1" {
		t.Errorf("expected voucher ID '1', got '%s'", vouchers.Data.Vouchers[0].ID)
	}
}

func TestFetchGraphQLDataMarshalError(t *testing.T) {
	client := NewClient("http://example.com")

	err := client.fetchGraphQLData("/whatever", "", nil)
	if !errors.Is(err, errOnTheGraph) {
		t.Errorf("expected errOnTheGraph, got %v", err)
	}
}

func TestFetchGraphQLDataRequestError(t *testing.T) {
	client := newMockedClient(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("network error")
	})

	err := client.fetchGraphQLData("/whatever", "query", nil)
	if !errors.Is(err, errOnTheGraph) {
		t.Errorf("expected errOnTheGraph, got %v", err)
	}
}

func TestFetchGraphQLDataReadError(t *testing.T) {
	client := newMockedClient(func(req *http.Request) (*http.Response, error) {
		r := io.NopCloser(&brokenReader{})
		return &http.Response{
			StatusCode: 200,
			Body:       r,
			Header:     make(http.Header),
		}, nil
	})

	err := client.fetchGraphQLData("/whatever", "query", nil)
	if !errors.Is(err, errOnTheGraph) {
		t.Errorf("expected errOnTheGraph, got %v", err)
	}
}

func TestFetchGraphQLDataUnmarshalError(t *testing.T) {
	client := newMockedClient(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString("invalid json")),
			Header:     make(http.Header),
		}, nil
	})

	err := client.fetchGraphQLData("/whatever", "query", nil)
	if !errors.Is(err, errOnTheGraph) {
		t.Errorf("expected errOnTheGraph, got %v", err)
	}
}

// brokenReader to simulate io.Reader errors
type brokenReader struct{}

func (b *brokenReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}
