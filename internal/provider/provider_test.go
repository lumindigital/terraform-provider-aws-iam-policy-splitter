package provider_test

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	internalProvider "github.com/Lumin-Digital/terraform-aws-iam-policy-splitter/internal/provider"
)

// testAccProtoV6ProviderFactories is used to instantiate a provider during acceptance testing.
// The factory function is called for each Terraform CLI command to create a provider
// server that the CLI can connect to and interact with.
//
//nolint:gochecknoglobals // required by Terraform acceptance testing framework
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"aws-iam-policy-splitter": providerserver.NewProtocol6WithError(internalProvider.New("test")()),
}
