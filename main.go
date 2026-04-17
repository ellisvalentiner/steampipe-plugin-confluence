package main

import (
	"github.com/ellisvalentiner/steampipe-plugin-confluence/confluence"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: confluence.Plugin})
}
