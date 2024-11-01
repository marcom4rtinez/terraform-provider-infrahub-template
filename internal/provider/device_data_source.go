// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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
	_ datasource.DataSource              = &deviceDataSource{}
	_ datasource.DataSourceWithConfigure = &deviceDataSource{}
)

// NewdeviceDataSource is a helper function to simplify the provider implementation.
func NewDeviceDataSource() datasource.DataSource {
	return &deviceDataSource{}
}

// deviceDataSource is the data source implementation.
type deviceDataSource struct {
	client     *graphql.Client
	Name       types.String `tfsdk:"name"`
	Id         types.String `tfsdk:"id"`
	Role       types.String `tfsdk:"role"`
	DeviceName types.String `tfsdk:"device_name"`
}

// Metadata returns the data source type name.
func (d *deviceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

// Schema defines the schema for the data source.
func (d *deviceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"device_name": schema.StringAttribute{
				Required: true, // This marks the attribute as required in the Terraform config
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"role": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *deviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading Device...\n")
	var config deviceDataSource

	// Read configuration into config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API with the specified device_name from the configuration
	device, err := infrahub_sdk.DeviceQuery(ctx, *d.client, config.DeviceName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read devices from Infrahub",
			err.Error(),
		)
		return
	}

	if len(device.InfraDevice.Edges) != 1 {
		resp.Diagnostics.AddError(
			"Didn't receive a single device, query didn't return exactly 1 device",
			"Expected exactly 1 device in response, got a different count.",
		)
		return
	}

	state := deviceDataSource{
		Name: types.StringValue(device.InfraDevice.Edges[0].Node.Name.Value),
		Id:   types.StringValue(device.InfraDevice.Edges[0].Node.Id),
		Role: types.StringValue(device.InfraDevice.Edges[0].Node.Role.Value),
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *deviceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
