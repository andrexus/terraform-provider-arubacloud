package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/andrexus/terraform-provider-arubacloud/arubacloud"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: arubacloud.Provider,
	})
}
