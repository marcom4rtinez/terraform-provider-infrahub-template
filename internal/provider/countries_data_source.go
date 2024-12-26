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
	_ datasource.DataSource              = &countriesDataSource{}
	_ datasource.DataSourceWithConfigure = &countriesDataSource{}
)

// NewCountriesDataSource is a helper function to simplify the provider implementation.
func NewCountriesDataSource() datasource.DataSource {
	return &countriesDataSource{}
}

type countriesDataSource struct {
	client    *graphql.Client
	Countries []countriesModel `tfsdk:"countries"`
}
type countriesModel struct {
	Edges_node_id                types.String `tfsdk:"id"`
	Edges_node_display_label     types.String `tfsdk:"display_label"`
	Edges_node_name_value        types.String `tfsdk:"name_value"`
	Edges_node_description_value types.String `tfsdk:"description_value"`
}

func (d *countriesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_countries"
}

func (d *countriesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"countries": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"display_label": schema.StringAttribute{
							Computed: true,
						},
						"name_value": schema.StringAttribute{
							Computed: true,
						},
						"description_value": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *countriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading countries data...")
	var config countriesDataSource

	// Read configuration into config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	response, err := infrahub_sdk.Countries(ctx, *d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read countries from Infrahub",
			err.Error(),
		)
		return
	}
	var state countriesDataSource
	for i := range response.LocationCountry.Edges {
		current := countriesModel{
			Edges_node_id:                types.StringValue(response.LocationCountry.Edges[i].Node.Id),
			Edges_node_display_label:     types.StringValue(response.LocationCountry.Edges[i].Node.Display_label),
			Edges_node_name_value:        types.StringValue(response.LocationCountry.Edges[i].Node.Name.Value),
			Edges_node_description_value: types.StringValue(response.LocationCountry.Edges[i].Node.Description.Value),
		}
		state.Countries = append(state.Countries, current)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *countriesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
