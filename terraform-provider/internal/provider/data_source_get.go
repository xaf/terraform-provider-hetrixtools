package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	hetrixtools "github.com/xaf/terraform-provider-hetrixtools/client"
)

var _ datasource.DataSource = (*getDataSource)(nil)
var _ datasource.DataSourceWithConfigure = (*getDataSource)(nil)

type queryAttribute struct {
	Name        string
	Description string
	Required    bool
}

type getDataSource struct {
	name        string
	description string
	queryAttrs  []queryAttribute
	read        func(context.Context, *hetrixtools.Client, map[string]string) (any, error)
	client      *hetrixtools.Client
}

func newGetDataSource(name string, description string, attrs []queryAttribute, read func(context.Context, *hetrixtools.Client, map[string]string) (any, error)) datasource.DataSource {
	return &getDataSource{name: name, description: description, queryAttrs: attrs, read: read}
}

func (d *getDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + d.name
}

func (d *getDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	attrs := map[string]schema.Attribute{
		"json": schema.StringAttribute{Computed: true, MarkdownDescription: "JSON response from the named HetrixTools client method."},
	}
	for _, attr := range d.queryAttrs {
		attrs[attr.Name] = schema.StringAttribute{Required: attr.Required, Optional: !attr.Required, MarkdownDescription: attr.Description}
	}
	resp.Schema = schema.Schema{MarkdownDescription: d.description, Attributes: attrs}
}

func (d *getDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*hetrixtools.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", "Expected *hetrixtools.Client.")
		return
	}
	d.client = c
}

func (d *getDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	values := map[string]types.String{}
	query := map[string]string{}
	for _, attr := range d.queryAttrs {
		var value types.String
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root(attr.Name), &value)...)
		values[attr.Name] = value
		if !value.IsNull() && !value.IsUnknown() && value.ValueString() != "" {
			query[attr.Name] = value.ValueString()
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.read(ctx, d.client, query)
	if err != nil {
		resp.Diagnostics.AddError("HetrixTools API request failed", err.Error())
		return
	}
	body, err := json.Marshal(result)
	if err != nil {
		resp.Diagnostics.AddError("HetrixTools API response encode failed", err.Error())
		return
	}

	for name, value := range values {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root(name), value)...)
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("json"), types.StringValue(normalizeJSON(body)))...)
}

func mapFromTerraformStringMap(ctx context.Context, input types.Map, diagnostics *diag.Diagnostics) map[string]string {
	if input.IsNull() || input.IsUnknown() {
		return nil
	}

	values := map[string]types.String{}
	diagnostics.Append(input.ElementsAs(ctx, &values, false)...)
	if diagnostics.HasError() {
		return nil
	}

	result := make(map[string]string, len(values))
	for key, value := range values {
		if !value.IsNull() && !value.IsUnknown() {
			result[key] = value.ValueString()
		}
	}
	return result
}
