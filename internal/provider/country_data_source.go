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
	_ datasource.DataSource              = &countryDataSource{}
	_ datasource.DataSourceWithConfigure = &countryDataSource{}
)

// NewCountryDataSource is a helper function to simplify the provider implementation.
func NewCountryDataSource() datasource.DataSource {
	return &countryDataSource{}
}

type countryDataSource struct {
	client                       *graphql.Client
	Country_name                 types.String `tfsdk:"country_name"`
	Edges_node_id                types.String `tfsdk:"id"`
	Edges_node_display_label     types.String `tfsdk:"display_label"`
	Edges_node_name_value        types.String `tfsdk:"name_value"`
	Edges_node_description_value types.String `tfsdk:"description_value"`
}

func (d *countryDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_country"
}

func (d *countryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"country_name": schema.StringAttribute{
				Required: true,
			},
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
	}
}

func (d *countryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading country data...")
	var config countryDataSource

	// Read configuration into config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	response, err := infrahub_sdk.Country(ctx, *d.client, config.Country_name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read country from Infrahub",
			err.Error(),
		)
		return
	}

	if len(response.LocationCountry.Edges) != 1 {
		resp.Diagnostics.AddError(
			"Didn't receive a single country, query didn't return exactly 1 country",
			"Expected exactly 1 country in response, got a different count.",
		)
		return
	}

	state := countryDataSource{
		Country_name:                 config.Country_name,
		Edges_node_id:                types.StringValue(response.LocationCountry.Edges[0].Node.Id),
		Edges_node_display_label:     types.StringValue(response.LocationCountry.Edges[0].Node.Display_label),
		Edges_node_name_value:        types.StringValue(response.LocationCountry.Edges[0].Node.Name.Value),
		Edges_node_description_value: types.StringValue(response.LocationCountry.Edges[0].Node.Description.Value),
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *countryDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
