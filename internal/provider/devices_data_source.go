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
	_ datasource.DataSource              = &devicesDataSource{}
	_ datasource.DataSourceWithConfigure = &devicesDataSource{}
)

// NewDevicesDataSource is a helper function to simplify the provider implementation.
func NewDevicesDataSource() datasource.DataSource {
	return &devicesDataSource{}
}

type devicesDataSource struct {
	client  *graphql.Client
	Devices []devicesModel `tfsdk:"devices"`
}
type devicesModel struct {
	Edges_node_id         types.String `tfsdk:"id"`
	Edges_node_name_value types.String `tfsdk:"name_value"`
	Edges_node_role_value types.String `tfsdk:"role_value"`
	Edges_node_role_color types.String `tfsdk:"role_color"`
}

func (d *devicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices"
}

func (d *devicesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"devices": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
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
					},
				},
			},
		},
	}
}

func (d *devicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading devices data...")
	var config devicesDataSource

	// Read configuration into config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	response, err := infrahub_sdk.Devices(ctx, *d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read devices from Infrahub",
			err.Error(),
		)
		return
	}
	var state devicesDataSource
	for i := range response.InfraDevice.Edges {
		current := devicesModel{
			Edges_node_id:         types.StringValue(response.InfraDevice.Edges[i].Node.Id),
			Edges_node_name_value: types.StringValue(response.InfraDevice.Edges[i].Node.Name.Value),
			Edges_node_role_value: types.StringValue(response.InfraDevice.Edges[i].Node.Role.Value),
			Edges_node_role_color: types.StringValue(response.InfraDevice.Edges[i].Node.Role.Color),
		}
		state.Devices = append(state.Devices, current)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *devicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
