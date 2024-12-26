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
	_ datasource.DataSource              = &topologyDataSource{}
	_ datasource.DataSourceWithConfigure = &topologyDataSource{}
)

// NewTopologyDataSource is a helper function to simplify the provider implementation.
func NewTopologyDataSource() datasource.DataSource {
	return &topologyDataSource{}
}

type topologyDataSource struct {
	client                       *graphql.Client
	Topology_name                types.String `tfsdk:"topology_name"`
	Edges_node_id                types.String `tfsdk:"id"`
	Edges_node_display_label     types.String `tfsdk:"display_label"`
	Edges_node_description_value types.String `tfsdk:"description_value"`
	Edges_node_description_id    types.String `tfsdk:"description_id"`
	Edges_node_name_value        types.String `tfsdk:"name_value"`
	Edges_node_name_id           types.String `tfsdk:"name_id"`
}

func (d *topologyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_topology"
}

func (d *topologyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"topology_name": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"display_label": schema.StringAttribute{
				Computed: true,
			},
			"description_value": schema.StringAttribute{
				Computed: true,
			},
			"description_id": schema.StringAttribute{
				Computed: true,
			},
			"name_value": schema.StringAttribute{
				Computed: true,
			},
			"name_id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *topologyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading topology data...")
	var config topologyDataSource

	// Read configuration into config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	response, err := infrahub_sdk.Topology(ctx, *d.client, config.Topology_name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read topology from Infrahub",
			err.Error(),
		)
		return
	}

	if len(response.TopologyTopology.Edges) != 1 {
		resp.Diagnostics.AddError(
			"Didn't receive a single topology, query didn't return exactly 1 topology",
			"Expected exactly 1 topology in response, got a different count.",
		)
		return
	}

	state := topologyDataSource{
		Topology_name:                config.Topology_name,
		Edges_node_id:                types.StringValue(response.TopologyTopology.Edges[0].Node.Id),
		Edges_node_display_label:     types.StringValue(response.TopologyTopology.Edges[0].Node.Display_label),
		Edges_node_description_value: types.StringValue(response.TopologyTopology.Edges[0].Node.Description.Value),
		Edges_node_description_id:    types.StringValue(response.TopologyTopology.Edges[0].Node.Description.Id),
		Edges_node_name_value:        types.StringValue(response.TopologyTopology.Edges[0].Node.Name.Value),
		Edges_node_name_id:           types.StringValue(response.TopologyTopology.Edges[0].Node.Name.Id),
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *topologyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
