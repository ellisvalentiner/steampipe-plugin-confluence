package confluence

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ctreminiom/go-atlassian/v2/confluence"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

const maxRetries = 3

// retryTransport retries requests on 429 and 5xx responses with exponential backoff.
type retryTransport struct {
	wrapped http.RoundTripper
}

func (t *retryTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	original := r

	for attempt := 0; attempt <= maxRetries; attempt++ {
		req := original
		if attempt > 0 {
			delay := retryDelay(resp, attempt)
			select {
			case <-original.Context().Done():
				return nil, original.Context().Err()
			case <-time.After(delay):
			}

			req, err = cloneRequestForRetry(original)
			if err != nil {
				return resp, err
			}
		}

		resp, err = t.wrapped.RoundTrip(req)
		if err != nil || !shouldRetry(resp) {
			break
		}

		if original.Body != nil && original.GetBody == nil {
			break
		}

		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}
	return resp, err
}

func cloneRequestForRetry(r *http.Request) (*http.Request, error) {
	clone := r.Clone(r.Context())
	if r.Body == nil {
		return clone, nil
	}
	if r.GetBody == nil {
		return nil, fmt.Errorf("request body is not replayable")
	}
	body, err := r.GetBody()
	if err != nil {
		return nil, err
	}
	clone.Body = body
	return clone, nil
}

func shouldRetry(resp *http.Response) bool {
	if resp == nil {
		return false
	}
	return resp.StatusCode == http.StatusTooManyRequests ||
		resp.StatusCode == http.StatusServiceUnavailable ||
		resp.StatusCode >= 500
}

// retryDelay returns the delay before the next attempt.
// Respects Retry-After if the server sends it, otherwise uses exponential backoff.
func retryDelay(resp *http.Response, attempt int) time.Duration {
	if resp != nil {
		if after := resp.Header.Get("Retry-After"); after != "" {
			if secs, err := strconv.Atoi(after); err == nil {
				return time.Duration(secs) * time.Second
			}
		}
	}
	return time.Duration(math.Pow(2, float64(attempt))) * time.Second
}

// dataCenterTransport rewrites the wiki/rest/api/ path prefix used by go-atlassian
// (which targets Confluence Cloud) to /rest/api/ for Confluence Data Center.
type dataCenterTransport struct {
	wrapped http.RoundTripper
}

func (t *dataCenterTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r = r.Clone(r.Context())
	r.URL.Path = strings.Replace(r.URL.Path, "/wiki/rest/api/", "/rest/api/", 1)
	return t.wrapped.RoundTrip(r)
}

func connect(_ context.Context, d *plugin.QueryData) (*confluence.Client, error) {
	confluenceConfig := GetConfig(d.Connection)

	baseURL := strings.TrimSpace(ptrString(confluenceConfig.BaseUrl))
	if baseURL == "" {
		baseURL = strings.TrimSpace(os.Getenv("CONFLUENCE_BASE_URL"))
	}

	deploymentType := strings.TrimSpace(ptrString(confluenceConfig.DeploymentType))
	if deploymentType == "" {
		deploymentType = strings.TrimSpace(os.Getenv("CONFLUENCE_DEPLOYMENT_TYPE"))
	}
	if deploymentType == "" {
		deploymentType = "cloud"
	}

	if baseURL == "" {
		return nil, fmt.Errorf("missing required config: set base_url or CONFLUENCE_BASE_URL")
	}

	token := strings.TrimSpace(ptrString(confluenceConfig.Token))
	if token == "" {
		token = strings.TrimSpace(os.Getenv("CONFLUENCE_TOKEN"))
	}
	if token == "" {
		return nil, fmt.Errorf("missing required config: set token or CONFLUENCE_TOKEN")
	}

	username := strings.TrimSpace(ptrString(confluenceConfig.Username))
	if username == "" {
		username = strings.TrimSpace(os.Getenv("CONFLUENCE_USERNAME"))
	}
	if deploymentType != "datacenter" && username == "" {
		return nil, fmt.Errorf("missing required config for cloud: set username or CONFLUENCE_USERNAME")
	}

	// Load connection from cache, which preserves throttling protection etc
	cacheKey := fmt.Sprintf("atlassian-confluence:%s:%s:%s", deploymentType, strings.TrimSuffix(baseURL, "/"), username)
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		if client, ok := cachedData.(*confluence.Client); ok {
			return client, nil
		}
	}

	isDataCenter := deploymentType == "datacenter"

	var transport http.RoundTripper = &retryTransport{wrapped: http.DefaultTransport}
	if isDataCenter {
		transport = &dataCenterTransport{wrapped: transport}
	}
	httpClient := &http.Client{Transport: transport}

	instance, err := confluence.New(httpClient, baseURL)
	if err != nil {
		return nil, err
	}

	if isDataCenter {
		instance.Auth.SetBearerToken(token)
	} else {
		instance.Auth.SetBasicAuth(username, token)
	}
	instance.Auth.SetUserAgent("steampipe-plugin-confluence")

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, instance)

	return instance, nil
}

type contentBody struct {
	ID             string `json:"id,omitempty"`
	Representation string `json:"representation,omitempty"`
	Value          string `json:"value,omitempty"`
}

func ptrString(v *string) string {
	if v == nil {
		return ""
	}

	return *v
}
