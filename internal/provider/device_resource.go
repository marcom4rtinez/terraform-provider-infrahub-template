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
	client                              *graphql.Client
	Edges_node_id                       types.String `tfsdk:"id"`
	Edges_node_name_value               types.String `tfsdk:"name_value"`
	Edges_node_role_value               types.String `tfsdk:"role_value"`
	Edges_node_role_id                  types.String `tfsdk:"role_id"`
	Edges_node_asn_node_id              types.String `tfsdk:"asn_node_id"`
	Edges_node_description_id           types.String `tfsdk:"description_id"`
	Edges_node_description_value        types.String `tfsdk:"description_value"`
	Edges_node_device_type_node_id      types.String `tfsdk:"device_type_node_id"`
	Edges_node_location_node_id         types.String `tfsdk:"location_node_id"`
	Edges_node_platform_node_id         types.String `tfsdk:"platform_node_id"`
	Edges_node_primary_address_node_id  types.String `tfsdk:"primary_address_node_id"`
	Edges_node_status_id                types.String `tfsdk:"status_id"`
	Edges_node_status_value             types.String `tfsdk:"status_value"`
	Edges_node_topology_node_id         types.String `tfsdk:"topology_node_id"`
	Edges_node_topology_node_name_value types.String `tfsdk:"topology_node_name_value"`
}

// Metadata returns the resource type name.
func (r *deviceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

// Schema defines the schema for the resource.
func (r *deviceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"role_id": schema.StringAttribute{
				Computed: true,
			},
			"description_id": schema.StringAttribute{
				Computed: true,
			},
			"status_id": schema.StringAttribute{
				Computed: true,
			},
			"topology_node_name_value": schema.StringAttribute{
				Computed: true,
			},
			"name_value": schema.StringAttribute{
				Required: true,
			},
			"role_value": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"asn_node_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"description_value": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"device_type_node_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"location_node_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"platform_node_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"primary_address_node_id": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"status_value": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"topology_node_id": schema.StringAttribute{
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

	var defaultDevice infrahub_sdk.InfraDeviceCreateInput

	// Assign each field, using the helper function to handle defaults
	defaultDevice.Name.Value = plan.Edges_node_name_value.ValueString()
	defaultDevice.Role.Value = plan.Edges_node_role_value.ValueString()
	defaultDevice.Asn.Id = plan.Edges_node_asn_node_id.ValueString()
	defaultDevice.Description.Value = plan.Edges_node_description_value.ValueString()
	defaultDevice.Device_type.Id = plan.Edges_node_device_type_node_id.ValueString()
	defaultDevice.Location.Id = plan.Edges_node_location_node_id.ValueString()
	defaultDevice.Platform.Id = plan.Edges_node_platform_node_id.ValueString()
	defaultDevice.Primary_address.Id = plan.Edges_node_primary_address_node_id.ValueString()
	defaultDevice.Status.Value = plan.Edges_node_status_value.ValueString()
	defaultDevice.Topology.Id = plan.Edges_node_topology_node_id.ValueString()

	tflog.Info(ctx, fmt.Sprint("Creating Device ", plan.Edges_node_name_value))

	response, err := infrahub_sdk.DeviceCreate(ctx, *r.client, defaultDevice)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create device in Infrahub",
			err.Error(),
		)
		return
	}
	plan.Edges_node_id = types.StringValue(response.InfraDeviceCreate.Object.GetId())
	plan.Edges_node_name_value = types.StringValue(response.InfraDeviceCreate.Object.Name.Value)
	plan.Edges_node_role_value = types.StringValue(response.InfraDeviceCreate.Object.Role.Value)
	plan.Edges_node_role_id = types.StringValue(response.InfraDeviceCreate.Object.Role.GetId())
	plan.Edges_node_asn_node_id = types.StringValue(response.InfraDeviceCreate.Object.Asn.Node.GetId())
	plan.Edges_node_description_id = types.StringValue(response.InfraDeviceCreate.Object.Description.GetId())
	plan.Edges_node_description_value = types.StringValue(response.InfraDeviceCreate.Object.Description.Value)
	plan.Edges_node_device_type_node_id = types.StringValue(response.InfraDeviceCreate.Object.Device_type.Node.GetId())
	plan.Edges_node_location_node_id = types.StringValue(response.InfraDeviceCreate.Object.Location.Node.GetId())
	plan.Edges_node_platform_node_id = types.StringValue(response.InfraDeviceCreate.Object.Platform.Node.GetId())
	plan.Edges_node_primary_address_node_id = types.StringValue(response.InfraDeviceCreate.Object.Primary_address.Node.GetId())
	plan.Edges_node_status_id = types.StringValue(response.InfraDeviceCreate.Object.Status.GetId())
	plan.Edges_node_status_value = types.StringValue(response.InfraDeviceCreate.Object.Status.Value)
	plan.Edges_node_topology_node_id = types.StringValue(response.InfraDeviceCreate.Object.Topology.Node.GetId())
	plan.Edges_node_topology_node_name_value = types.StringValue(response.InfraDeviceCreate.Object.Topology.Node.Name.Value)

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

	tflog.Info(ctx, fmt.Sprint("Reading Device ", state.Edges_node_name_value))

	// Call the API with the specified device_name from the configuration
	response, err := infrahub_sdk.Device(ctx, *r.client, state.Edges_node_name_value.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read device from Infrahub",
			err.Error(),
		)
		return
	}

	if len(response.InfraDevice.Edges) != 1 {
		resp.Diagnostics.AddError(
			"Didn't receive a single device, query didn't return exactly 1 device",
			"Expected exactly 1 device in response, got a different count.",
		)
		return
	}
	state.Edges_node_id = types.StringValue(response.InfraDevice.Edges[0].Node.GetId())
	state.Edges_node_name_value = types.StringValue(response.InfraDevice.Edges[0].Node.Name.Value)
	state.Edges_node_role_value = types.StringValue(response.InfraDevice.Edges[0].Node.Role.Value)
	state.Edges_node_role_id = types.StringValue(response.InfraDevice.Edges[0].Node.Role.GetId())
	state.Edges_node_asn_node_id = types.StringValue(response.InfraDevice.Edges[0].Node.Asn.Node.GetId())
	state.Edges_node_description_id = types.StringValue(response.InfraDevice.Edges[0].Node.Description.GetId())
	state.Edges_node_description_value = types.StringValue(response.InfraDevice.Edges[0].Node.Description.Value)
	state.Edges_node_device_type_node_id = types.StringValue(response.InfraDevice.Edges[0].Node.Device_type.Node.GetId())
	state.Edges_node_location_node_id = types.StringValue(response.InfraDevice.Edges[0].Node.Location.Node.GetId())
	state.Edges_node_platform_node_id = types.StringValue(response.InfraDevice.Edges[0].Node.Platform.Node.GetId())
	state.Edges_node_primary_address_node_id = types.StringValue(response.InfraDevice.Edges[0].Node.Primary_address.Node.GetId())
	state.Edges_node_status_id = types.StringValue(response.InfraDevice.Edges[0].Node.Status.GetId())
	state.Edges_node_status_value = types.StringValue(response.InfraDevice.Edges[0].Node.Status.Value)
	state.Edges_node_topology_node_id = types.StringValue(response.InfraDevice.Edges[0].Node.Topology.Node.GetId())
	state.Edges_node_topology_node_name_value = types.StringValue(response.InfraDevice.Edges[0].Node.Topology.Node.Name.Value)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *deviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve the planned configuration values from Terraform
	var plan deviceResource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the current state
	var state deviceResource
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var updateInput infrahub_sdk.InfraDeviceUpsertInput

	// Prepare the update input using values from the plan and applying defaults
	updateInput.Name.Value = setDefault(plan.Edges_node_name_value.ValueString(), state.Edges_node_name_value.ValueString())
	updateInput.Role.Value = setDefault(plan.Edges_node_role_value.ValueString(), state.Edges_node_role_value.ValueString())
	updateInput.Asn.Id = setDefault(plan.Edges_node_asn_node_id.ValueString(), state.Edges_node_asn_node_id.ValueString())
	updateInput.Description.Value = setDefault(plan.Edges_node_description_value.ValueString(), state.Edges_node_description_value.ValueString())
	updateInput.Device_type.Id = setDefault(plan.Edges_node_device_type_node_id.ValueString(), state.Edges_node_device_type_node_id.ValueString())
	updateInput.Location.Id = setDefault(plan.Edges_node_location_node_id.ValueString(), state.Edges_node_location_node_id.ValueString())
	updateInput.Platform.Id = setDefault(plan.Edges_node_platform_node_id.ValueString(), state.Edges_node_platform_node_id.ValueString())
	updateInput.Primary_address.Id = setDefault(plan.Edges_node_primary_address_node_id.ValueString(), state.Edges_node_primary_address_node_id.ValueString())
	updateInput.Status.Value = setDefault(plan.Edges_node_status_value.ValueString(), state.Edges_node_status_value.ValueString())
	updateInput.Topology.Id = setDefault(plan.Edges_node_topology_node_id.ValueString(), state.Edges_node_topology_node_id.ValueString())
	updateInput.Id = state.Edges_node_id.ValueString()

	// Log the update operation
	tflog.Info(ctx, fmt.Sprintf("Updating Device %s", state.Edges_node_name_value.ValueString()))

	// Send the update request to the API
	response, err := infrahub_sdk.DeviceUpsert(ctx, *r.client, updateInput)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update device in Infrahub",
			err.Error(),
		)
		return
	}
	plan.Edges_node_id = types.StringValue(response.InfraDeviceUpsert.Object.GetId())
	plan.Edges_node_name_value = types.StringValue(response.InfraDeviceUpsert.Object.Name.Value)
	plan.Edges_node_role_value = types.StringValue(response.InfraDeviceUpsert.Object.Role.Value)
	plan.Edges_node_role_id = types.StringValue(response.InfraDeviceUpsert.Object.Role.GetId())
	plan.Edges_node_asn_node_id = types.StringValue(response.InfraDeviceUpsert.Object.Asn.Node.GetId())
	plan.Edges_node_description_id = types.StringValue(response.InfraDeviceUpsert.Object.Description.GetId())
	plan.Edges_node_description_value = types.StringValue(response.InfraDeviceUpsert.Object.Description.Value)
	plan.Edges_node_device_type_node_id = types.StringValue(response.InfraDeviceUpsert.Object.Device_type.Node.GetId())
	plan.Edges_node_location_node_id = types.StringValue(response.InfraDeviceUpsert.Object.Location.Node.GetId())
	plan.Edges_node_platform_node_id = types.StringValue(response.InfraDeviceUpsert.Object.Platform.Node.GetId())
	plan.Edges_node_primary_address_node_id = types.StringValue(response.InfraDeviceUpsert.Object.Primary_address.Node.GetId())
	plan.Edges_node_status_id = types.StringValue(response.InfraDeviceUpsert.Object.Status.GetId())
	plan.Edges_node_status_value = types.StringValue(response.InfraDeviceUpsert.Object.Status.Value)
	plan.Edges_node_topology_node_id = types.StringValue(response.InfraDeviceUpsert.Object.Topology.Node.GetId())
	plan.Edges_node_topology_node_name_value = types.StringValue(response.InfraDeviceUpsert.Object.Topology.Node.Name.Value)

	// Set the updated state with the latest data
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
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
	_, err := infrahub_sdk.DeviceDelete(ctx, *r.client, state.Edges_node_id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Device",
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
