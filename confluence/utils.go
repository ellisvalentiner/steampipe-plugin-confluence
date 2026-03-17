package confluence

import (
	"context"
	"net/http"
	"strings"

	"github.com/ctreminiom/go-atlassian/confluence"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
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

	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "atlassian-confluence"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*confluence.Client), nil
	}

	// Prefer config options given in Steampipe
	confluenceConfig := GetConfig(d.Connection)

	isDataCenter := confluenceConfig.DeploymentType != nil && *confluenceConfig.DeploymentType == "datacenter"

	var httpClient *http.Client
	if isDataCenter {
		httpClient = &http.Client{
			Transport: &dataCenterTransport{wrapped: http.DefaultTransport},
		}
	}

	instance, err := confluence.New(httpClient, *confluenceConfig.BaseUrl)
	if err != nil {
		return nil, err
	}

	if isDataCenter {
		instance.Auth.SetBearerToken(*confluenceConfig.Token)
	} else {
		instance.Auth.SetBasicAuth(*confluenceConfig.Username, *confluenceConfig.Token)
	}
	instance.Auth.SetUserAgent("curl/7.54.0")

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, instance)

	return instance, nil
}
