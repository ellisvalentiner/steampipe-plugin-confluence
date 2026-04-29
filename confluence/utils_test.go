package confluence

import (
	"net/http"
	"testing"
	"time"
)

func TestPtrString(t *testing.T) {
	t.Run("nil pointer returns empty string", func(t *testing.T) {
		if got := ptrString(nil); got != "" {
			t.Errorf("ptrString(nil) = %q; want %q", got, "")
		}
	})
	t.Run("non-nil pointer returns value", func(t *testing.T) {
		s := "hello"
		if got := ptrString(&s); got != s {
			t.Errorf("ptrString(&%q) = %q; want %q", s, got, s)
		}
	})
	t.Run("empty string pointer returns empty string", func(t *testing.T) {
		s := ""
		if got := ptrString(&s); got != "" {
			t.Errorf("ptrString(&\"\") = %q; want %q", got, "")
		}
	})
}

func TestShouldRetry(t *testing.T) {
	tests := []struct {
		name string
		resp *http.Response
		want bool
	}{
		{"nil response never retries", nil, false},
		{"200 OK never retries", &http.Response{StatusCode: 200}, false},
		{"404 not found never retries", &http.Response{StatusCode: 404}, false},
		{"429 too many requests retries", &http.Response{StatusCode: 429}, true},
		{"500 internal server error retries", &http.Response{StatusCode: 500}, true},
		{"503 service unavailable retries", &http.Response{StatusCode: 503}, true},
		{"504 gateway timeout retries", &http.Response{StatusCode: 504}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldRetry(tt.resp)
			if got != tt.want {
				t.Errorf("shouldRetry(%v) = %v; want %v", tt.resp, got, tt.want)
			}
		})
	}
}

func TestRetryDelay(t *testing.T) {
	t.Run("nil response uses exponential backoff", func(t *testing.T) {
		got := retryDelay(nil, 1)
		if got != 2*time.Second {
			t.Errorf("retryDelay(nil, 1) = %v; want %v", got, 2*time.Second)
		}
		got = retryDelay(nil, 2)
		if got != 4*time.Second {
			t.Errorf("retryDelay(nil, 2) = %v; want %v", got, 4*time.Second)
		}
	})

	t.Run("Retry-After header is respected", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: 429,
			Header:     http.Header{"Retry-After": []string{"30"}},
		}
		got := retryDelay(resp, 1)
		if got != 30*time.Second {
			t.Errorf("retryDelay with Retry-After:30 = %v; want %v", got, 30*time.Second)
		}
	})

	t.Run("invalid Retry-After falls back to backoff", func(t *testing.T) {
		resp := &http.Response{
			StatusCode: 429,
			Header:     http.Header{"Retry-After": []string{"not-a-number"}},
		}
		got := retryDelay(resp, 2)
		if got != 4*time.Second {
			t.Errorf("retryDelay with bad Retry-After at attempt 2 = %v; want %v", got, 4*time.Second)
		}
	})

	t.Run("response without Retry-After uses backoff", func(t *testing.T) {
		resp := &http.Response{StatusCode: 503, Header: http.Header{}}
		got := retryDelay(resp, 3)
		if got != 8*time.Second {
			t.Errorf("retryDelay(503, 3) = %v; want %v", got, 8*time.Second)
		}
	})
}

func TestDataCenterTransportRewritesPath(t *testing.T) {
	var capturedPath string
	stub := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		capturedPath = r.URL.Path
		return &http.Response{StatusCode: 200}, nil
	})

	transport := &dataCenterTransport{wrapped: stub}
	req, _ := http.NewRequest("GET", "https://confluence.example.com/wiki/rest/api/space", nil)
	if _, err := transport.RoundTrip(req); err != nil {
		t.Fatal(err)
	}

	want := "/rest/api/space"
	if capturedPath != want {
		t.Errorf("dataCenterTransport rewrote path to %q; want %q", capturedPath, want)
	}
}

// roundTripFunc lets us create an http.RoundTripper from a plain function.
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
