package provider

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure AwsIamPolicySplitterProvider satisfies various provider interfaces.
var _ provider.Provider = &AwsIamPolicySplitterProvider{}
var _ provider.ProviderWithFunctions = &AwsIamPolicySplitterProvider{}
var _ provider.ProviderWithEphemeralResources = &AwsIamPolicySplitterProvider{}
var _ provider.ProviderWithActions = &AwsIamPolicySplitterProvider{}

// AwsIamPolicySplitterProvider defines the provider implementation.
type AwsIamPolicySplitterProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// AwsIamPolicySplitterProviderModel describes the provider data model.
type AwsIamPolicySplitterProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
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
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	var data AwsIamPolicySplitterProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *AwsIamPolicySplitterProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *AwsIamPolicySplitterProvider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *AwsIamPolicySplitterProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAwsIamPolicySplitterDataSource,
	}
}

func (p *AwsIamPolicySplitterProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{}
}

func (p *AwsIamPolicySplitterProvider) Actions(_ context.Context) []func() action.Action {
	return []func() action.Action{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AwsIamPolicySplitterProvider{
			version: version,
		}
	}
}
