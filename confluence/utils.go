package confluence

import (
	"context"

	"github.com/ctreminiom/go-atlassian/confluence"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
)

func connect(_ context.Context, d *plugin.QueryData) (*confluence.Client, error) {

	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "atlassian-confluence"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*confluence.Client), nil
	}

	// Prefer config options given in Steampipe
	confluenceConfig := GetConfig(d.Connection)

	instance, err := confluence.New(nil, *confluenceConfig.BaseUrl)
	if err != nil {
		return nil, err
	}
	instance.Auth.SetBasicAuth(*confluenceConfig.Username, *confluenceConfig.Token)
	instance.Auth.SetUserAgent("curl/7.54.0")

	// Save to cache
	d.ConnectionManager.Cache.Set(cacheKey, instance)

	return instance, nil
}
