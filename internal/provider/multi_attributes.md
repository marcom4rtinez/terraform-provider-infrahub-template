```golang

type interfaceDataSource struct {
	client                       *graphql.Client
	Interface_name               types.String `tfsdk:"interface_name"`
	Edges_node_id                types.String `tfsdk:"edges_node_id"`
	Edges_node_hfid              types.List   `tfsdk:"edges_node_hfid"`
    .....
}

func (d *interfaceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"interface_name": schema.StringAttribute{
				Required: true,
			},
			"edges_node_id": schema.StringAttribute{
				Computed: true,
			},
			"edges_node_hfid": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
			},
            ....
		},
	}
}

func (d *interfaceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {


	x, _ := types.ListValueFrom(context.Background(), types.StringType, response.InfraIPAddress.Edges[0].Node.Hfid)

	state := interfaceDataSource{
		Interface_name:               config.Interface_name,
		Edges_node_id:                types.StringValue(response.InfraIPAddress.Edges[0].Node.Id),
		Edges_node_hfid:              x,
        ...
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
```