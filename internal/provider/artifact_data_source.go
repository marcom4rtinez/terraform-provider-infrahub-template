package provider

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &artifactDataSource{}
	_ datasource.DataSourceWithConfigure = &artifactDataSource{}
)

// NewArtifactDataSource is a helper function to simplify the provider implementation.
func NewArtifactDataSource() datasource.DataSource {
	return &artifactDataSource{}
}

// artifactDataSource is the data source implementation.
type artifactDataSource struct {
	client         *graphql.Client
	ArtifactId     types.String `tfsdk:"artifact_id"`
	InfrahubServer types.String `tfsdk:"infrahub_server"`
	Content        types.String `tfsdk:"content"`
}

// Metadata returns the data source type name.
func (d *artifactDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_artifact"
}

// Schema defines the schema for the data source.
func (d *artifactDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"artifact_id": schema.StringAttribute{
				Required: true,
			},
			"infrahub_server": schema.StringAttribute{
				Required: true,
			},
			"content": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *artifactDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Artifact...\n")
	var config artifactDataSource

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	url := fmt.Sprintf("%s/api/storage/object/%s", config.InfrahubServer, config.ArtifactId.ValueString())
	artifactRequest, err := http.NewRequest("GET", url, nil)
	httpClient := &http.Client{}
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Request",
			err.Error(),
		)
		return
	}

	artifactResponse, err := httpClient.Do(artifactRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read artifact from Infrahub",
			err.Error(),
		)
		return
	}
	defer artifactResponse.Body.Close()
	if artifactResponse.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"Got non 200 HTTP Status Code",
			fmt.Sprintf("Non-200 response: %v\n", artifactResponse.Status))
		return
	}

	body, err := io.ReadAll(artifactResponse.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading response body: %v\n", err),
			fmt.Sprintf("Error reading response body: %v\n", err),
		)
		return
	}

	state := artifactDataSource{
		ArtifactId: types.StringValue(config.ArtifactId.String()),
		Content:    types.StringValue(string(body)),
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *artifactDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
