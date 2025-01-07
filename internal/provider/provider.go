package provider

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure InfrahubProvider satisfies various provider interfaces.
var _ provider.Provider = &InfrahubProvider{}
var _ provider.ProviderWithFunctions = &InfrahubProvider{}

// InfrahubProvider defines the provider implementation.
type InfrahubProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// InfrahubProviderModel describes the provider data model.
type InfrahubProviderModel struct {
	ApiKey         types.String `tfsdk:"api_key"`
	InfrahubServer types.String `tfsdk:"infrahub_server"`
	Branch         types.String `tfsdk:"branch"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &InfrahubProvider{
			version: version,
		}
	}
}

func (p *InfrahubProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "infrahub"
	resp.Version = p.version
}

func (p *InfrahubProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API Key to access Infrahub",
				Optional:            true,
				Sensitive:           true,
			},
			"infrahub_server": schema.StringAttribute{
				MarkdownDescription: "Infrahub Server running API",
				Optional:            true,
			},
			"branch": schema.StringAttribute{
				MarkdownDescription: "Infrahub Branch",
				Optional:            true,
			},
		},
	}
}

func (p *InfrahubProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data InfrahubProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	if data.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown API Key for Infrahub",
			"The provider cannot read the Infrahub API Key as there is an unknown configuration value. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the INFRAHUB_API environment variable.",
		)
	}

	if data.InfrahubServer.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("infrahub_server"),
			"Unknown Infrahub API Endpoint",
			"The provider cannot read the Infrahub API address as there is an unknown configuration value for the API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the INFRAHUB_SERVER environment variable.",
		)
	}

	if data.Branch.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("branch"),
			"Unknown Infrahub Branch",
			"The provider cannot read the Infrahub API as there is an unknown branch configuration value for the API. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the INFRAHUB_BRANCH environment variable.",
		)
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	infrahubApi := os.Getenv("INFRAHUB_API")
	infrahub_server := os.Getenv("INFRAHUB_SERVER")
	branch := os.Getenv("INFRAHUB_BRANCH")

	if !data.ApiKey.IsNull() {
		infrahubApi = data.ApiKey.ValueString()
	}

	if !data.InfrahubServer.IsNull() {
		infrahub_server = data.InfrahubServer.ValueString()
	}

	if !data.Branch.IsNull() {
		branch = data.Branch.ValueString()
	}

	if infrahubApi == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing API Key",
			"The provider cannot find the Infrahub API key as there is a missing or empty value for the API Key. "+
				"Set the API Key value in the configuration or use the INFRAHUB_API environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if infrahub_server == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("infrahub_server"),
			"Missing Infrahub Server address",
			"The provider cannot find the Infrahub API Server address as there is a missing or empty value for the Server address. "+
				"Set the Infrahub Server address value in the configuration or use the INFRAHUB_SERVER environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if branch == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("branch"),
			"Missing Infrahub Branch",
			"The provider cannot find the Infrahub Branch as there is a missing or empty value. "+
				"Set the Infrahub Server address value in the configuration or use the INFRAHUB_SERVER environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	httpClient := &http.Client{
		Transport: &AuthTransport{
			Token:     infrahubApi,
			Transport: http.DefaultTransport,
		},
	}

	client := graphql.NewClient(fmt.Sprintf("%s/graphql/%s", infrahub_server, branch), httpClient)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *InfrahubProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDeviceResource,
	}
}

func (p *InfrahubProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAccountsDataSource,
		NewAutonomoussystemDataSource,
		NewBgpsessionsDataSource,
		NewCountriesDataSource,
		NewCountryDataSource,
		NewDevicequeryDataSource,
		NewDevicesDataSource,
		NewDevicetypeDataSource,
		NewInterfaceDataSource,
		NewIpaddressDataSource,
		NewPlatformDataSource,
		NewTopologyDataSource,
		NewArtifactDataSource,
	}
}

func (p *InfrahubProvider) Functions(ctx context.Context) []func() function.Function {
	return nil
}

type AuthTransport struct {
	Token     string
	Transport http.RoundTripper
}

// RoundTrip adds the authorization header and delegates the request to the original transport.
func (a *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-INFRAHUB-KEY", a.Token)
	return a.Transport.RoundTrip(req)
}

// Helper function to set a string value with a default if empty.
func setDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
