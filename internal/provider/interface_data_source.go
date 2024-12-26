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
	_ datasource.DataSource              = &interfaceDataSource{}
	_ datasource.DataSourceWithConfigure = &interfaceDataSource{}
)

// NewInterfaceDataSource is a helper function to simplify the provider implementation.
func NewInterfaceDataSource() datasource.DataSource {
	return &interfaceDataSource{}
}

type interfaceDataSource struct {
	client                       *graphql.Client
	Interface_name               types.String `tfsdk:"interface_name"`
	Edges_node_id                types.String `tfsdk:"id"`
	Edges_node_description_value types.String `tfsdk:"description_value"`
	Edges_node_address_ip        types.String `tfsdk:"address_ip"`
	Edges_node_address_value     types.String `tfsdk:"address_value"`
}

func (d *interfaceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_interface"
}

func (d *interfaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"interface_name": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"description_value": schema.StringAttribute{
				Computed: true,
			},
			"address_ip": schema.StringAttribute{
				Computed: true,
			},
			"address_value": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *interfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading interface data...")
	var config interfaceDataSource

	// Read configuration into config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	response, err := infrahub_sdk.Interface(ctx, *d.client, config.Interface_name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read interface from Infrahub",
			err.Error(),
		)
		return
	}

	if len(response.InfraIPAddress.Edges) != 1 {
		resp.Diagnostics.AddError(
			"Didn't receive a single interface, query didn't return exactly 1 interface",
			"Expected exactly 1 interface in response, got a different count.",
		)
		return
	}

	state := interfaceDataSource{
		Interface_name:               config.Interface_name,
		Edges_node_id:                types.StringValue(response.InfraIPAddress.Edges[0].Node.Id),
		Edges_node_description_value: types.StringValue(response.InfraIPAddress.Edges[0].Node.Description.Value),
		Edges_node_address_ip:        types.StringValue(response.InfraIPAddress.Edges[0].Node.Address.Ip),
		Edges_node_address_value:     types.StringValue(response.InfraIPAddress.Edges[0].Node.Address.Value),
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *interfaceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
