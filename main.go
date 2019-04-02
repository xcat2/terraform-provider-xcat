package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.ibm.com/yangsbj/terraform-provider-xcat/xcat"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: xcat.Provider})
}
