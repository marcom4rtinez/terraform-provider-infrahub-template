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
	_ datasource.DataSource              = &accountsDataSource{}
	_ datasource.DataSourceWithConfigure = &accountsDataSource{}
)

// NewAccountsDataSource is a helper function to simplify the provider implementation.
func NewAccountsDataSource() datasource.DataSource {
	return &accountsDataSource{}
}

type accountsDataSource struct {
	client   *graphql.Client
	Accounts []accountsModel `tfsdk:"accounts"`
}
type accountsModel struct {
	Edges_node_id                 types.String `tfsdk:"id"`
	Edges_node_status_id          types.String `tfsdk:"status_id"`
	Edges_node_status_description types.String `tfsdk:"status_description"`
	Edges_node_status_color       types.String `tfsdk:"status_color"`
	Edges_node_status_value       types.String `tfsdk:"status_value"`
}

func (d *accountsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_accounts"
}

func (d *accountsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"accounts": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"status_id": schema.StringAttribute{
							Computed: true,
						},
						"status_description": schema.StringAttribute{
							Computed: true,
						},
						"status_color": schema.StringAttribute{
							Computed: true,
						},
						"status_value": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *accountsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading accounts data...")
	var config accountsDataSource

	// Read configuration into config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	response, err := infrahub_sdk.Accounts(ctx, *d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read accounts from Infrahub",
			err.Error(),
		)
		return
	}
	var state accountsDataSource
	for i := range response.CoreAccount.Edges {
		current := accountsModel{
			Edges_node_id:                 types.StringValue(response.CoreAccount.Edges[i].Node.Id),
			Edges_node_status_id:          types.StringValue(response.CoreAccount.Edges[i].Node.Status.Id),
			Edges_node_status_description: types.StringValue(response.CoreAccount.Edges[i].Node.Status.Description),
			Edges_node_status_color:       types.StringValue(response.CoreAccount.Edges[i].Node.Status.Color),
			Edges_node_status_value:       types.StringValue(response.CoreAccount.Edges[i].Node.Status.Value),
		}
		state.Accounts = append(state.Accounts, current)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *accountsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
