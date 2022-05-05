package main

import (
	"github.com/ellisvalentiner/steampipe-plugin-confluence/confluence"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: confluence.Plugin})
}
