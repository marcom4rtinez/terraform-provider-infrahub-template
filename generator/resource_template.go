// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

const resourceTemplateContent = `package provider

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
	_ resource.Resource              = &{{.QueryName}}Resource{}
	_ resource.ResourceWithConfigure = &{{.QueryName}}Resource{}
)

// New{{.QueryName | title }}Resource is a helper function to simplify the provider implementation.
func New{{.QueryName | title }}Resource() resource.Resource {
	return &{{.QueryName}}Resource{}
}

// {{.QueryName }}Resource is the resource implementation.
type {{.QueryName }}Resource struct {
	client         *graphql.Client
	{{.Required | title }} types.String ` + "`tfsdk:\"{{.Required}}\"`" + `
	{{- range .GenqlientFields }}
	{{ .Name | title }} types.String ` + "`tfsdk:\"{{ .HumanReadableName }}\"`" + `
	{{- end }}
}

// Metadata returns the resource type name.
func (r *{{.QueryName}}Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_{{.QueryName}}"
}

// Schema defines the schema for the resource.
func (r *{{.QueryName}}Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"{{.Required}}": schema.StringAttribute{
				Required: true,
			},
			{{- range .GenqlientFieldsReadOnly }}
			"{{ .HumanReadableName }}": schema.StringAttribute{
				Computed: true,
			},
			{{- end }}
			{{- range .GenqlientFieldsModify }}
			"{{ .HumanReadableName }}": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			{{- end }}
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *{{.QueryName}}Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan {{.QueryName}}Resource
	tflog.Info(ctx, req.Config.Raw.String())
	tflog.Info(ctx, req.Plan.Raw.String())
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var default{{ .QueryName | title }} infrahub_sdk.{{ .ObjectName }}CreateInput

	// Assign each field, using the helper function to handle defaults
	{{- $defaultCreate :=  .QueryName | title  }}
	{{- range .GenqlientFieldsModify }}
	default{{$defaultCreate}}.{{ .InputObjectNames }} = plan.{{ .Name | title }}.ValueString()
	{{- end }}

	tflog.Info(ctx, fmt.Sprint("Creating {{ .QueryName | title }} ", plan.{{.Required | title }}))

	response, err := infrahub_sdk.{{ .QueryName | title }}Create(ctx, *r.client, default{{ .QueryName | title }})
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create {{ .QueryName }} in Infrahub",
			err.Error(),
		)
		return
	}

	{{- $defaultCreateObject :=  .ObjectName }}
	{{- range .GenqlientFields }}
	plan.{{ .Name | title }} = types.StringValue(response.{{ $defaultCreateObject }}Create.Object.{{ .PlainObject }})
	{{- end }}


	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *{{.QueryName}}Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "Reading {{.QueryName | title }}...")
	var state {{ .QueryName }}Resource

	// Read configuration into config
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprint("Reading {{ .QueryName | title }} ", state.{{ .Required | title }}))

	// Call the API with the specified device_name from the configuration
	response, err := infrahub_sdk.{{ .QueryName | title }}(ctx, *r.client, state.{{ .Required | title }}.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read {{ .QueryName }} from Infrahub",
			err.Error(),
		)
		return
	}

	if len(response.{{ .ObjectName }}.Edges) != 1 {
		resp.Diagnostics.AddError(
			"Didn't receive a single {{ .QueryName }}, query didn't return exactly 1 {{ .QueryName }}",
			"Expected exactly 1 {{ .QueryName }} in response, got a different count.",
		)
		return
	}


	{{- $defaultObject :=  .ObjectName }}
	{{- range .GenqlientFields }}
	state.{{ .Name | title }} = types.StringValue(response.{{ .Query }})
	{{- end }}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *{{.QueryName}}Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve the planned configuration values from Terraform
	var plan {{ .QueryName }}Resource
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve the current state
	var state {{ .QueryName }}Resource
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var updateInput infrahub_sdk.{{ .ObjectName }}UpsertInput

	// Prepare the update input using values from the plan and applying defaults
	{{- range .GenqlientFieldsModify }}
	updateInput.{{ .InputObjectNames }} = setDefault(plan.{{ .Name | title }}.ValueString(), state.{{ .Name | title }}.ValueString())
	{{- end }}

	// Log the update operation
	tflog.Info(ctx, fmt.Sprintf("Updating {{ .QueryName | title }} %s", state.{{ .Required | title }}.ValueString()))

	// Send the update request to the API
	response, err := infrahub_sdk.{{ .QueryName | title }}Upsert(ctx, *r.client, updateInput)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to update device in Infrahub",
			err.Error(),
		)
		return
	}

	{{- $defaultUpsertObject :=  .ObjectName }}
	{{- range .GenqlientFields }}
	plan.{{ .Name | title }} = types.StringValue(response.{{ $defaultUpsertObject }}Upsert.Object.{{ .PlainObject }})
	{{- end }}

	// Set the updated state with the latest data
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *{{.QueryName}}Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state {{ .QueryName }}Resource
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	{{- $firstId :=  (index .GenqlientFieldsReadOnly 0).Name | title  }}
	_, err := infrahub_sdk.{{ .QueryName | title }}Delete(ctx, *r.client, state.{{$firstId}}.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting {{ .QueryName | title }}",
			"Could not delete {{ .QueryName }}, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *{{.QueryName}}Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
`
