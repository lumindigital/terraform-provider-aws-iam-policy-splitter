package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AwsIamPolicySplitterDataSource{}

func NewAwsIamPolicySplitterDataSource() datasource.DataSource {
	return &AwsIamPolicySplitterDataSource{}
}

// AwsIamPolicySplitterDataSource defines the data source implementation.
type AwsIamPolicySplitterDataSource struct {
	client *http.Client
}

// AwsIamPolicySplitterDataSourceModel describes the data source data model.
type AwsIamPolicySplitterDataSourceModel struct {
	ConfigurableAttribute types.String `tfsdk:"configurable_attribute"`
	ID                    types.String `tfsdk:"id"`
}

func (d *AwsIamPolicySplitterDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_example"
}

func (d *AwsIamPolicySplitterDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Takes a list of iam policy documents and splits them into multiple policies if they exceed the AWS size limits. Uses a binpacking algorithm to split the policies as efficiently as possible. Attempts to split policies in a way that keeps like policies together",

		Attributes: map[string]schema.Attribute{
			"configurable_attribute": schema.StringAttribute{
				MarkdownDescription: "Example configurable attribute",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
			},
		},
	}
}

func (d *AwsIamPolicySplitterDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf(
				"Expected *http.Client, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}

	d.client = client
}

func (d *AwsIamPolicySplitterDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data AwsIamPolicySplitterDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := d.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.ID = types.StringValue("example-id")

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
