package main

import (
	"flag"

	c "github.com/Mongey/terraform-provider-kafka-connect/connect"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "Set to true to run the provider with support for debuggers.")
	flag.Parse()

	opts := &plugin.ServeOpts{
		Debug:        debug,
		ProviderAddr: "Mongey/kafka-connect",
		ProviderFunc: c.Provider,
	}

	plugin.Serve(opts)
}
