package provider

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/Lumin-Digital/terraform-aws-iam-policy-splitter/internal/binpack"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AwsIamPolicySplitterDataSource{}

func NewAwsIamPolicySplitterDataSource() datasource.DataSource {
	return &AwsIamPolicySplitterDataSource{}
}

// AwsIamPolicySplitterDataSource defines the data source implementation.
type AwsIamPolicySplitterDataSource struct{}

func (d *AwsIamPolicySplitterDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_data_source"
}

func (d *AwsIamPolicySplitterDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Takes a list of IAM policy JSON strings and splits them into bins using a best-fit decreasing bin packing algorithm. Useful for distributing policies across IAM roles while respecting AWS size limits.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for this data source instance",
				Computed:            true,
			},
			"json_policies": schema.ListAttribute{
				MarkdownDescription: "List of IAM policy JSON strings to split across bins",
				Required:            true,
				ElementType:         types.StringType,
			},
			"bin_size": schema.Int64Attribute{
				MarkdownDescription: "Maximum total character length per bin (default: 6144)",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"max_bins": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of bins to distribute policies into (default: 10)",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"policy_bins": schema.ListNestedAttribute{
				MarkdownDescription: "Result of bin packing: a list of bins, each containing a list of policy ARN strings",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"arns": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

// PolicyBinModel represents a single bin in the output.
type PolicyBinModel struct {
	Arns []types.String `tfsdk:"arns"`
}

// AwsIamPolicySplitterModel describes the data source data model.
type AwsIamPolicySplitterModel struct {
	ID           types.String     `tfsdk:"id"`
	JsonPolicies []types.String   `tfsdk:"json_policies"`
	BinSize      types.Int64      `tfsdk:"bin_size"`
	MaxBins      types.Int64      `tfsdk:"max_bins"`
	PolicyBins   []PolicyBinModel `tfsdk:"policy_bins"`
}

func (d *AwsIamPolicySplitterDataSource) Configure(
	_ context.Context,
	_ datasource.ConfigureRequest,
	_ *datasource.ConfigureResponse,
) {
}

func (d *AwsIamPolicySplitterDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data AwsIamPolicySplitterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Apply defaults
	binSize := int64(6144)
	if !data.BinSize.IsNull() && !data.BinSize.IsUnknown() {
		binSize = data.BinSize.ValueInt64()
	}
	data.BinSize = types.Int64Value(binSize)

	maxBins := int64(10)
	if !data.MaxBins.IsNull() && !data.MaxBins.IsUnknown() {
		maxBins = data.MaxBins.ValueInt64()
	}
	data.MaxBins = types.Int64Value(maxBins)

	// Convert policy ARNs to plain strings
	items := make([]string, len(data.JsonPolicies))
	for i, arn := range data.JsonPolicies {
		items[i] = arn.ValueString()
	}

	tflog.Debug(ctx, "binpacking policy ARNs", map[string]interface{}{
		"count":    len(items),
		"bin_size": binSize,
		"max_bins": maxBins,
	})

	// Run bin packing
	bins, err := binpack.BestFit(items, int(maxBins), int(binSize))
	if err != nil {
		resp.Diagnostics.AddError(
			"Bin Packing Error",
			fmt.Sprintf("Failed to distribute policy ARNs into bins: %s", err),
		)
		return
	}

	// Convert bins to output model (skip empty bins)
	data.PolicyBins = []PolicyBinModel{}
	for _, bin := range bins {
		if len(bin.Items) == 0 {
			continue
		}
		arnValues := make([]types.String, len(bin.Items))
		for i, item := range bin.Items {
			arnValues[i] = types.StringValue(item)
		}
		data.PolicyBins = append(data.PolicyBins, PolicyBinModel{Arns: arnValues})
	}

	// Generate deterministic ID from inputs
	h := sha256.New()
	for _, item := range items {
		h.Write([]byte(item))
	}
	data.ID = types.StringValue(fmt.Sprintf("%x", h.Sum(nil))[:16])

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
