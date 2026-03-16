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
		Address: "registry.terraform.io/Lumin-Digital/aws-iam-policy-splitter",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
