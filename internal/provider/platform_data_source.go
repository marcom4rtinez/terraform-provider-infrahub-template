package provider

import (
	"context"
	"fmt"

	infrahub_sdk "github.com/opsmill/infrahub-sdk-go"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &platformDataSource{}
	_ datasource.DataSourceWithConfigure = &platformDataSource{}
)

// NewPlatformDataSource is a helper function to simplify the provider implementation.
func NewPlatformDataSource() datasource.DataSource {
	return &platformDataSource{}
}

type platformDataSource struct {
	client                               *graphql.Client
	Platform_name                        types.String `tfsdk:"platform_name"`
	Edges_node_id                        types.String `tfsdk:"id"`
	Edges_node_description_value         types.String `tfsdk:"description_value"`
	Edges_node_containerlab_os_value     types.String `tfsdk:"containerlab_os_value"`
	Edges_node_name_value                types.String `tfsdk:"name_value"`
	Edges_node_nornir_platform_value     types.String `tfsdk:"nornir_platform_value"`
	Edges_node_netmiko_device_type_value types.String `tfsdk:"netmiko_device_type_value"`
	Edges_node_napalm_driver_value       types.String `tfsdk:"napalm_driver_value"`
}

func (d *platformDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_platform"
}

func (d *platformDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"platform_name": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"description_value": schema.StringAttribute{
				Computed: true,
			},
			"containerlab_os_value": schema.StringAttribute{
				Computed: true,
			},
			"name_value": schema.StringAttribute{
				Computed: true,
			},
			"nornir_platform_value": schema.StringAttribute{
				Computed: true,
			},
			"netmiko_device_type_value": schema.StringAttribute{
				Computed: true,
			},
			"napalm_driver_value": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *platformDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading platform data...")
	var config platformDataSource

	// Read configuration into config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	response, err := infrahub_sdk.Platform(ctx, *d.client, config.Platform_name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read platform from Infrahub",
			err.Error(),
		)
		return
	}

	if len(response.InfraPlatform.Edges) != 1 {
		resp.Diagnostics.AddError(
			"Didn't receive a single platform, query didn't return exactly 1 platform",
			"Expected exactly 1 platform in response, got a different count.",
		)
		return
	}

	state := platformDataSource{
		Platform_name:                        config.Platform_name,
		Edges_node_id:                        types.StringValue(response.InfraPlatform.Edges[0].Node.Id),
		Edges_node_description_value:         types.StringValue(response.InfraPlatform.Edges[0].Node.Description.Value),
		Edges_node_containerlab_os_value:     types.StringValue(response.InfraPlatform.Edges[0].Node.Containerlab_os.Value),
		Edges_node_name_value:                types.StringValue(response.InfraPlatform.Edges[0].Node.Name.Value),
		Edges_node_nornir_platform_value:     types.StringValue(response.InfraPlatform.Edges[0].Node.Nornir_platform.Value),
		Edges_node_netmiko_device_type_value: types.StringValue(response.InfraPlatform.Edges[0].Node.Netmiko_device_type.Value),
		Edges_node_napalm_driver_value:       types.StringValue(response.InfraPlatform.Edges[0].Node.Napalm_driver.Value),
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *platformDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(graphql.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *graphql.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = &client
}
