package confluence

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ctreminiom/go-atlassian/v2/confluence"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

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
		return cachedData.(*confluence.Client), nil
	}

	isDataCenter := deploymentType == "datacenter"

	httpClient := &http.Client{Transport: http.DefaultTransport}
	if isDataCenter {
		httpClient.Transport = &dataCenterTransport{wrapped: http.DefaultTransport}
	}

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

func ptrString(v *string) string {
	if v == nil {
		return ""
	}

	return *v
}
