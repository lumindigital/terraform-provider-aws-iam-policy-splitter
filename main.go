package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/Lumin-Digital/terraform-aws-iam-policy-splitter/internal/provider"
)

var (
	version = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		// TODO: Update this string with the published name of your provider.
		// Also update the tfplugindocs generate command to either remove the
		// -provider-name flag or set its value to the updated provider name.
		Address: "registry.terraform.io/hashicorp/scaffolding",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
