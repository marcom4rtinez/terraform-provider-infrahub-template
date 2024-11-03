// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/Khan/genqlient/graphql"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	infrahub_sdk "github.com/opsmill/infrahub-sdk-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &deviceResource{}
	_ resource.ResourceWithConfigure = &deviceResource{}
)

// NewDeviceResource is a helper function to simplify the provider implementation.
func NewDeviceResource() resource.Resource {
	return &deviceResource{}
}

// deviceResource is the resource implementation.
type deviceResource struct {
	client         *graphql.Client
	Name           types.String `tfsdk:"name"`
	Id             types.String `tfsdk:"id"`
	Role           types.String `tfsdk:"role"`
	DeviceName     types.String `tfsdk:"device_name"`
	Location       types.String `tfsdk:"location"`
	Status         types.String `tfsdk:"status"`
	Asn            types.String `tfsdk:"asn"`
	PrimaryAddress types.String `tfsdk:"primary_address"`
	DeviceType     types.String `tfsdk:"device_type"`
	Platform       types.String `tfsdk:"platform"`
	Topology       types.String `tfsdk:"topology"`
}

// Metadata returns the resource type name.
func (r *deviceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

// Schema defines the schema for the resource.
func (r *deviceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"device_name": schema.StringAttribute{
				Required: true, // This marks the attribute as required in the Terraform config
			},
			"name": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"role": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"location": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"asn": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"primary_address": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"device_type": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"platform": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"topology": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *deviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan deviceResource
	tflog.Info(ctx, req.Config.Raw.String())
	tflog.Info(ctx, req.Plan.Raw.String())
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var defaultDeviceCreate infrahub_sdk.InfraDeviceCreateInput

	// Assign each field, using the helper function to handle defaults
	defaultDeviceCreate.Name.Value = plan.DeviceName.ValueString()
	defaultDeviceCreate.Role.Value = setDefault(plan.Role.ValueString(), "leaf")
	defaultDeviceCreate.Location.Id = setDefault(plan.Location.ValueString(), "1802dff5-0440-6832-3953-c51c60b72dfb")
	defaultDeviceCreate.Status.Value = setDefault(plan.Status.ValueString(), "active")
	defaultDeviceCreate.Asn.Id = setDefault(plan.Asn.ValueString(), "1802e1f2-b767-64d5-3952-c5185b886d16")
	defaultDeviceCreate.Primary_address.Id = setDefault(plan.PrimaryAddress.ValueString(), "1802e1f3-01ef-e7d7-395d-c51af858c4ce")
	defaultDeviceCreate.Device_type.Id = setDefault(plan.DeviceType.ValueString(), "1802dff2-a0db-e70e-3958-c518addd4232")
	defaultDeviceCreate.Platform.Id = setDefault(plan.Platform.ValueString(), "1802dff2-6d6a-5bfd-3956-c517263e9d03")
	defaultDeviceCreate.Topology.Id = setDefault(plan.Topology.ValueString(), "1802dffa-7ad7-fb58-3956-c51d2b3b03fb")

	tflog.Info(ctx, fmt.Sprint("Creating Device ", plan.DeviceName))

	device, err := infrahub_sdk.InfraDeviceCreate(ctx, *r.client, defaultDeviceCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create device in Infrahub",
			err.Error(),
		)
		return
	}

	plan.Name = types.StringValue(device.InfraDeviceCreate.Object.Name.Value)
	plan.Id = types.StringValue(device.InfraDeviceCreate.Object.Id)
	plan.Role = types.StringValue(device.InfraDeviceCreate.Object.Role.Value)
	plan.Platform = types.StringValue(device.InfraDeviceCreate.Object.Platform.Node.Id)
	plan.PrimaryAddress = types.StringValue(device.InfraDeviceCreate.Object.Primary_address.Node.Id)
	plan.Asn = types.StringValue(device.InfraDeviceCreate.Object.Asn.Node.Id)
	plan.DeviceType = types.StringValue(device.InfraDeviceCreate.Object.Device_type.Node.Id)
	plan.Location = types.StringValue(device.InfraDeviceCreate.Object.Location.Node.GetId())
	plan.Status = types.StringValue(device.InfraDeviceCreate.Object.Status.Value)
	plan.Topology = types.StringValue(device.InfraDeviceCreate.Object.Topology.Node.Id)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *deviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading Device...")
	var state deviceResource

	// Read configuration into config
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprint("Reading Device ", state.DeviceName))

	// Call the API with the specified device_name from the configuration
	device, err := infrahub_sdk.DeviceQuery(ctx, *r.client, state.DeviceName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read device from Infrahub",
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

	state.Name = types.StringValue(device.InfraDevice.Edges[0].Node.Name.Value)
	state.Id = types.StringValue(device.InfraDevice.Edges[0].Node.Id)
	state.Role = types.StringValue(device.InfraDevice.Edges[0].Node.Role.Value)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *deviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// // Retrieve values from plan
	// var plan orderResourceModel
	// diags := req.Plan.Get(ctx, &plan)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
	// // Generate API request body from plan
	// var hashicupsItems []hashicups.OrderItem
	// for _, item := range plan.Items {
	// 	hashicupsItems = append(hashicupsItems, hashicups.OrderItem{
	// 		Coffee: hashicups.Coffee{
	// 			ID: int(item.Coffee.ID.ValueInt64()),
	// 		},
	// 		Quantity: int(item.Quantity.ValueInt64()),
	// 	})
	// }
	// // Update existing order
	// _, err := r.client.UpdateOrder(plan.ID.ValueString(), hashicupsItems)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error Updating HashiCups Order",
	// 		"Could not update order, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }
	// // Fetch updated items from GetOrder as UpdateOrder items are not
	// // populated.
	// order, err := r.client.GetOrder(plan.ID.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error Reading HashiCups Order",
	// 		"Could not read HashiCups order ID "+plan.ID.ValueString()+": "+err.Error(),
	// 	)
	// 	return
	// }
	// // Update resource state with updated items and timestamp
	// plan.Items = []orderItemModel{}
	// for _, item := range order.Items {
	// 	plan.Items = append(plan.Items, orderItemModel{
	// 		Coffee: orderItemCoffeeModel{
	// 			ID:          types.Int64Value(int64(item.Coffee.ID)),
	// 			Name:        types.StringValue(item.Coffee.Name),
	// 			Teaser:      types.StringValue(item.Coffee.Teaser),
	// 			Description: types.StringValue(item.Coffee.Description),
	// 			Price:       types.Float64Value(item.Coffee.Price),
	// 			Image:       types.StringValue(item.Coffee.Image),
	// 		},
	// 		Quantity: types.Int64Value(int64(item.Quantity)),
	// 	})
	// }
	// plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	// diags = resp.State.Set(ctx, plan)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *deviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state deviceResource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	_, err := infrahub_sdk.InfraDeviceDelete(ctx, *r.client, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Infrahub Device",
			"Could not delete device, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *deviceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(graphql.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *graphql.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = &client
}

// Helper function to set a string value with a default if empty.
func setDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
