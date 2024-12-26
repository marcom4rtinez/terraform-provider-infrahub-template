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
	_ datasource.DataSource              = &bgpsessionsDataSource{}
	_ datasource.DataSourceWithConfigure = &bgpsessionsDataSource{}
)

// NewBgpsessionsDataSource is a helper function to simplify the provider implementation.
func NewBgpsessionsDataSource() datasource.DataSource {
	return &bgpsessionsDataSource{}
}

type bgpsessionsDataSource struct {
	client      *graphql.Client
	Bgpsessions []bgpsessionsModel `tfsdk:"bgpsessions"`
}
type bgpsessionsModel struct {
	Edges_node_id                           types.String `tfsdk:"id"`
	Edges_node_display_label                types.String `tfsdk:"display_label"`
	Edges_node_description_value            types.String `tfsdk:"description_value"`
	Edges_node_remote_ip_node_address_value types.String `tfsdk:"remote_ip_node_address_value"`
}

func (d *bgpsessionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bgpsessions"
}

func (d *bgpsessionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"bgpsessions": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"display_label": schema.StringAttribute{
							Computed: true,
						},
						"description_value": schema.StringAttribute{
							Computed: true,
						},
						"remote_ip_node_address_value": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *bgpsessionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading bgpsessions data...")
	var config bgpsessionsDataSource

	// Read configuration into config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	response, err := infrahub_sdk.Bgpsessions(ctx, *d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read bgpsessions from Infrahub",
			err.Error(),
		)
		return
	}
	var state bgpsessionsDataSource
	for i := range response.InfraBGPSession.Edges {
		current := bgpsessionsModel{
			Edges_node_id:                           types.StringValue(response.InfraBGPSession.Edges[i].Node.Id),
			Edges_node_display_label:                types.StringValue(response.InfraBGPSession.Edges[i].Node.Display_label),
			Edges_node_description_value:            types.StringValue(response.InfraBGPSession.Edges[i].Node.Description.Value),
			Edges_node_remote_ip_node_address_value: types.StringValue(response.InfraBGPSession.Edges[i].Node.Remote_ip.Node.Address.Value),
		}
		state.Bgpsessions = append(state.Bgpsessions, current)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *bgpsessionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
