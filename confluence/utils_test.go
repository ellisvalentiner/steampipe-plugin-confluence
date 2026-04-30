package confluence

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
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

func TestCloneRequestForRetry(t *testing.T) {
	t.Run("nil body clones without error", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "https://example.com/", nil)
		clone, err := cloneRequestForRetry(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if clone.URL.String() != req.URL.String() {
			t.Errorf("clone URL = %q; want %q", clone.URL, req.URL)
		}
	})

	t.Run("replayable body is recreated", func(t *testing.T) {
		body := "hello"
		req, _ := http.NewRequest("POST", "https://example.com/", strings.NewReader(body))
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader(body)), nil
		}
		clone, err := cloneRequestForRetry(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		got, err := io.ReadAll(clone.Body)
		if err != nil {
			t.Fatalf("unexpected error reading clone body: %v", err)
		}
		if string(got) != body {
			t.Errorf("clone body = %q; want %q", got, body)
		}
	})

	t.Run("non-replayable body returns error", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "https://example.com/", nil)
		// Manually set Body without GetBody to simulate a non-replayable body.
		req.Body = io.NopCloser(strings.NewReader("data"))
		// GetBody is nil (not set), so it is not replayable
		_, err := cloneRequestForRetry(req)
		if err == nil {
			t.Error("expected error for non-replayable body, got nil")
		}
	})
}

func TestRetryTransportBodyClosed(t *testing.T) {
	var bodyClosed atomic.Bool
	attempts := 0

	stub := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		attempts++
		body := io.NopCloser(bytes.NewReader(nil))
		if attempts == 1 {
			body = &trackingCloser{ReadCloser: body, closed: &bodyClosed}
		}
		return &http.Response{
			StatusCode: 503,
			Body:       body,
			Header:     http.Header{},
		}, nil
	})

	transport := &retryTransport{wrapped: stub}
	req, _ := http.NewRequest("GET", "https://example.com/", nil)
	resp, _ := transport.RoundTrip(req)
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}

	if !bodyClosed.Load() {
		t.Error("expected previous response body to be closed before retry, but it was not")
	}
	if attempts < 2 {
		t.Errorf("expected at least 2 attempts, got %d", attempts)
	}
}

func TestRetryTransportNonReplayableBodyNotRetried(t *testing.T) {
	attempts := 0
	stub := roundTripFunc(func(r *http.Request) (*http.Response, error) {
		attempts++
		return &http.Response{
			StatusCode: 503,
			Body:       io.NopCloser(bytes.NewReader(nil)),
			Header:     http.Header{},
		}, nil
	})

	transport := &retryTransport{wrapped: stub}
	// Manually set Body without GetBody to simulate a non-replayable body.
	req, _ := http.NewRequest("POST", "https://example.com/", nil)
	req.Body = io.NopCloser(strings.NewReader("payload"))
	resp, _ := transport.RoundTrip(req)
	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}

	if attempts != 1 {
		t.Errorf("non-replayable body request should not be retried; got %d attempts", attempts)
	}
}

// trackingCloser records when Close is called.
type trackingCloser struct {
	io.ReadCloser
	closed *atomic.Bool
}

func (tc *trackingCloser) Close() error {
	tc.closed.Store(true)
	return tc.ReadCloser.Close()
}
