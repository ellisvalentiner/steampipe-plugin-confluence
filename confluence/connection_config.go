package confluence

import (
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/schema"
)

type confluenceConfig struct {
	BaseUrl  *string `cty:"base_url"`
	Username *string `cty:"username"`
	Token    *string `cty:"token"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"base_url": {
		Type: schema.TypeString,
	},
	"username": {
		Type: schema.TypeString,
	},
	"token": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &confluenceConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) confluenceConfig {
	if connection == nil || connection.Config == nil {
		return confluenceConfig{}
	}
	config, _ := connection.Config.(confluenceConfig)
	return config
}
