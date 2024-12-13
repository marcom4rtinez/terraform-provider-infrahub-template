// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

const datasourceTemplateContent = `package provider

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
	_ datasource.DataSource              = &{{.QueryName}}DataSource{}
	_ datasource.DataSourceWithConfigure = &{{.QueryName}}DataSource{}
)

// New{{.QueryName | title }}DataSource is a helper function to simplify the provider implementation.
func New{{.QueryName | title }}DataSource() datasource.DataSource {
	return &{{.QueryName}}DataSource{}
}


type {{.StructName}} struct {
	client     *graphql.Client
	{{- if .Required }}
	{{.Required | title }} types.String ` + "`tfsdk:\"{{.Required}}\"`" + `
	{{- range .GenqlientFields }}
	{{ .Name | title }} types.String ` + "`tfsdk:\"{{ .HumanReadableName }}\"`" + `
	{{- end }}
	{{- else }}
	{{ .QueryName | title }} []{{ .QueryName }}Model ` + "`tfsdk:\"{{ .QueryName }}\"`" + `
	{{- end }}
}

{{- if not .Required }}
type {{ .QueryName}}Model struct {
	{{- range .GenqlientFields }}
	{{ .Name | title }} types.String ` + "`tfsdk:\"{{ .HumanReadableName }}\"`" + `
	{{- end }}
}
{{- end }}

func (d *{{.QueryName}}DataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_{{.QueryName}}"
}

func (d *{{.QueryName}}DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			{{- if .Required }}
			"{{.Required}}": schema.StringAttribute{
				Required: true,
			},
			{{- range .GenqlientFields }}
			"{{ .HumanReadableName }}": schema.StringAttribute{
				Computed: true,
			},
			{{- end }}
			{{- else}}
			"{{ .QueryName }}": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						{{- range .GenqlientFields }}
						"{{ .HumanReadableName }}": schema.StringAttribute{
							Computed: true,
						},
						{{- end }}
					},
				},
			},
			{{- end }}
		},
	}
}

func (d *{{.QueryName }}DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading {{.QueryName}} data...")
	var config {{.StructName}}

	// Read configuration into config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	{{- if .Required }}
	response, err := infrahub_sdk.{{.QueryName | title}}(ctx, *d.client, config.{{.Required | title }}.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read {{.QueryName}} from Infrahub",
			err.Error(),
		)
		return
	}

	if len(response.{{.ObjectName}}.Edges) != 1 {
		resp.Diagnostics.AddError(
			"Didn't receive a single {{.QueryName}}, query didn't return exactly 1 {{.QueryName}}",
			"Expected exactly 1 {{.QueryName}} in response, got a different count.",
		)
		return
	}

	state := {{.StructName}}{
		{{.Required | title}}: config.{{.Required | title }},
		{{- range .GenqlientFields }}
		{{ .Name | title }}: types.StringValue(response.{{ .Query }}),
		{{- end }}
	}
	{{- else }}
	response, err := infrahub_sdk.{{.QueryName | title}}(ctx, *d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read {{.QueryName}} from Infrahub",
			err.Error(),
		)
		return
	}
	var state {{.StructName}}
	for i, _ := range response.{{.ObjectName}}.Edges {
		current := {{.QueryName}}Model{
			{{- range .GenqlientFields }}
			{{ .Name | title }}: types.StringValue(response.{{ .Query }}),
			{{- end }}
		}
		state.{{.QueryName | title }} = append(state.{{.QueryName| title }}, current)
	}
	{{- end}}


	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *{{.QueryName}}DataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
`
