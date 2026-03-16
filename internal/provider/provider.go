package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure AwsIamPolicySplitterProvider satisfies various provider interfaces.
var _ provider.Provider = &AwsIamPolicySplitterProvider{}

// AwsIamPolicySplitterProvider defines the provider implementation.
type AwsIamPolicySplitterProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func (p *AwsIamPolicySplitterProvider) Metadata(
	_ context.Context,
	_ provider.MetadataRequest,
	resp *provider.MetadataResponse,
) {
	resp.TypeName = "aws-iam-policy-splitter"
	resp.Version = p.version
}

func (p *AwsIamPolicySplitterProvider) Schema(
	_ context.Context,
	_ provider.SchemaRequest,
	resp *provider.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{},
	}
}

func (p *AwsIamPolicySplitterProvider) Configure(
	_ context.Context,
	_ provider.ConfigureRequest,
	_ *provider.ConfigureResponse,
) {
}

func (p *AwsIamPolicySplitterProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *AwsIamPolicySplitterProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAwsIamPolicySplitterDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AwsIamPolicySplitterProvider{
			version: version,
		}
	}
}
