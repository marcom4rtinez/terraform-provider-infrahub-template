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
	_ datasource.DataSource              = &devicetypeDataSource{}
	_ datasource.DataSourceWithConfigure = &devicetypeDataSource{}
)

// NewDevicetypeDataSource is a helper function to simplify the provider implementation.
func NewDevicetypeDataSource() datasource.DataSource {
	return &devicetypeDataSource{}
}

type devicetypeDataSource struct {
	client                              *graphql.Client
	Device_type_name                    types.String `tfsdk:"device_type_name"`
	Edges_node_id                       types.String `tfsdk:"id"`
	Edges_node_platform_node_id         types.String `tfsdk:"platform_node_id"`
	Edges_node_platform_node_name_value types.String `tfsdk:"platform_node_name_value"`
	Edges_node_description_id           types.String `tfsdk:"description_id"`
	Edges_node_description_value        types.String `tfsdk:"description_value"`
	Edges_node_name_value               types.String `tfsdk:"name_value"`
	Edges_node_weight_value             types.String `tfsdk:"weight_value"`
}

func (d *devicetypeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devicetype"
}

func (d *devicetypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"device_type_name": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"platform_node_id": schema.StringAttribute{
				Computed: true,
			},
			"platform_node_name_value": schema.StringAttribute{
				Computed: true,
			},
			"description_id": schema.StringAttribute{
				Computed: true,
			},
			"description_value": schema.StringAttribute{
				Computed: true,
			},
			"name_value": schema.StringAttribute{
				Computed: true,
			},
			"weight_value": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *devicetypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading devicetype data...")
	var config devicetypeDataSource

	// Read configuration into config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	response, err := infrahub_sdk.Devicetype(ctx, *d.client, config.Device_type_name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read devicetype from Infrahub",
			err.Error(),
		)
		return
	}

	if len(response.InfraDeviceType.Edges) != 1 {
		resp.Diagnostics.AddError(
			"Didn't receive a single devicetype, query didn't return exactly 1 devicetype",
			"Expected exactly 1 devicetype in response, got a different count.",
		)
		return
	}

	state := devicetypeDataSource{
		Device_type_name:                    config.Device_type_name,
		Edges_node_id:                       types.StringValue(response.InfraDeviceType.Edges[0].Node.Id),
		Edges_node_platform_node_id:         types.StringValue(response.InfraDeviceType.Edges[0].Node.Platform.Node.Id),
		Edges_node_platform_node_name_value: types.StringValue(response.InfraDeviceType.Edges[0].Node.Platform.Node.Name.Value),
		Edges_node_description_id:           types.StringValue(response.InfraDeviceType.Edges[0].Node.Description.Id),
		Edges_node_description_value:        types.StringValue(response.InfraDeviceType.Edges[0].Node.Description.Value),
		Edges_node_name_value:               types.StringValue(response.InfraDeviceType.Edges[0].Node.Name.Value),
		Edges_node_weight_value:             types.StringValue(response.InfraDeviceType.Edges[0].Node.Weight.Value),
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *devicetypeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
