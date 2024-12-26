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
	_ datasource.DataSource              = &devicequeryDataSource{}
	_ datasource.DataSourceWithConfigure = &devicequeryDataSource{}
)

// NewDevicequeryDataSource is a helper function to simplify the provider implementation.
func NewDevicequeryDataSource() datasource.DataSource {
	return &devicequeryDataSource{}
}

type devicequeryDataSource struct {
	client                             *graphql.Client
	Device_name                        types.String `tfsdk:"device_name"`
	Edges_node_id                      types.String `tfsdk:"id"`
	Edges_node_name_value              types.String `tfsdk:"name_value"`
	Edges_node_role_value              types.String `tfsdk:"role_value"`
	Edges_node_role_color              types.String `tfsdk:"role_color"`
	Edges_node_role_description        types.String `tfsdk:"role_description"`
	Edges_node_role_id                 types.String `tfsdk:"role_id"`
	Edges_node_platform_node_id        types.String `tfsdk:"platform_node_id"`
	Edges_node_primary_address_node_id types.String `tfsdk:"primary_address_node_id"`
	Edges_node_status_id               types.String `tfsdk:"status_id"`
	Edges_node_topology_node_id        types.String `tfsdk:"topology_node_id"`
	Edges_node_device_type_node_id     types.String `tfsdk:"device_type_node_id"`
	Edges_node_asn_node_asn_id         types.String `tfsdk:"asn_node_asn_id"`
	Edges_node_description_value       types.String `tfsdk:"description_value"`
}

func (d *devicequeryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devicequery"
}

func (d *devicequeryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"device_name": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name_value": schema.StringAttribute{
				Computed: true,
			},
			"role_value": schema.StringAttribute{
				Computed: true,
			},
			"role_color": schema.StringAttribute{
				Computed: true,
			},
			"role_description": schema.StringAttribute{
				Computed: true,
			},
			"role_id": schema.StringAttribute{
				Computed: true,
			},
			"platform_node_id": schema.StringAttribute{
				Computed: true,
			},
			"primary_address_node_id": schema.StringAttribute{
				Computed: true,
			},
			"status_id": schema.StringAttribute{
				Computed: true,
			},
			"topology_node_id": schema.StringAttribute{
				Computed: true,
			},
			"device_type_node_id": schema.StringAttribute{
				Computed: true,
			},
			"asn_node_asn_id": schema.StringAttribute{
				Computed: true,
			},
			"description_value": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *devicequeryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading devicequery data...")
	var config devicequeryDataSource

	// Read configuration into config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	response, err := infrahub_sdk.Devicequery(ctx, *d.client, config.Device_name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read devicequery from Infrahub",
			err.Error(),
		)
		return
	}

	if len(response.InfraDevice.Edges) != 1 {
		resp.Diagnostics.AddError(
			"Didn't receive a single devicequery, query didn't return exactly 1 devicequery",
			"Expected exactly 1 devicequery in response, got a different count.",
		)
		return
	}

	state := devicequeryDataSource{
		Device_name:                        config.Device_name,
		Edges_node_id:                      types.StringValue(response.InfraDevice.Edges[0].Node.Id),
		Edges_node_name_value:              types.StringValue(response.InfraDevice.Edges[0].Node.Name.Value),
		Edges_node_role_value:              types.StringValue(response.InfraDevice.Edges[0].Node.Role.Value),
		Edges_node_role_color:              types.StringValue(response.InfraDevice.Edges[0].Node.Role.Color),
		Edges_node_role_description:        types.StringValue(response.InfraDevice.Edges[0].Node.Role.Description),
		Edges_node_role_id:                 types.StringValue(response.InfraDevice.Edges[0].Node.Role.Id),
		Edges_node_platform_node_id:        types.StringValue(response.InfraDevice.Edges[0].Node.Platform.Node.Id),
		Edges_node_primary_address_node_id: types.StringValue(response.InfraDevice.Edges[0].Node.Primary_address.Node.Id),
		Edges_node_status_id:               types.StringValue(response.InfraDevice.Edges[0].Node.Status.Id),
		Edges_node_topology_node_id:        types.StringValue(response.InfraDevice.Edges[0].Node.Topology.Node.Id),
		Edges_node_device_type_node_id:     types.StringValue(response.InfraDevice.Edges[0].Node.Device_type.Node.Id),
		Edges_node_asn_node_asn_id:         types.StringValue(response.InfraDevice.Edges[0].Node.Asn.Node.Asn.Id),
		Edges_node_description_value:       types.StringValue(response.InfraDevice.Edges[0].Node.Description.Value),
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *devicequeryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
