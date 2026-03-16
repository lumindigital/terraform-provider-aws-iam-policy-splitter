package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccAwsIamPolicySplitterDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "aws-iam-policy-splitter_data_source" "test" {
  json_policies = ["arn:aws:iam::123456789012:policy/A", "arn:aws:iam::123456789012:policy/B"]
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.aws-iam-policy-splitter_data_source.test",
						tfjsonpath.New("bin_size"),
						knownvalue.Int64Exact(6144),
					),
					statecheck.ExpectKnownValue(
						"data.aws-iam-policy-splitter_data_source.test",
						tfjsonpath.New("max_bins"),
						knownvalue.Int64Exact(10),
					),
					statecheck.ExpectKnownValue(
						"data.aws-iam-policy-splitter_data_source.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func TestAccAwsIamPolicySplitterDataSource_customBinSize(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "aws-iam-policy-splitter_data_source" "test" {
  json_policies = ["short", "medium-length", "a"]
  bin_size    = 20
  max_bins    = 3
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.aws-iam-policy-splitter_data_source.test",
						tfjsonpath.New("bin_size"),
						knownvalue.Int64Exact(20),
					),
					statecheck.ExpectKnownValue(
						"data.aws-iam-policy-splitter_data_source.test",
						tfjsonpath.New("max_bins"),
						knownvalue.Int64Exact(3),
					),
				},
			},
		},
	})
}

func TestAccAwsIamPolicySplitterDataSource_binSizeValidation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "aws-iam-policy-splitter_data_source" "test" {
  json_policies = ["policy-a"]
  bin_size    = 0
}
`,
				ExpectError: regexp.MustCompile(`must be at least 1`),
			},
		},
	})
}

func TestAccAwsIamPolicySplitterDataSource_maxBinsValidation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "aws-iam-policy-splitter_data_source" "test" {
  json_policies = ["policy-a"]
  max_bins    = 0
}
`,
				ExpectError: regexp.MustCompile(`must be at least 1`),
			},
		},
	})
}

func TestAccAwsIamPolicySplitterDataSource_negativeValues(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "aws-iam-policy-splitter_data_source" "test" {
  json_policies = ["policy-a"]
  bin_size    = -5
  max_bins    = -1
}
`,
				ExpectError: regexp.MustCompile(`must be at least 1`),
			},
		},
	})
}

func TestAccAwsIamPolicySplitterDataSource_emptyPolicies(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "aws-iam-policy-splitter_data_source" "test" {
  json_policies = []
}
`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.aws-iam-policy-splitter_data_source.test",
						tfjsonpath.New("policy_bins"),
						knownvalue.ListSizeExact(0),
					),
				},
			},
		},
	})
}
