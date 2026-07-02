package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	hetrixtools "github.com/xaf/go-hetrixtools/client"
)

var _ resource.Resource = (*serverAgentResource)(nil)
var _ resource.ResourceWithConfigure = (*serverAgentResource)(nil)
var _ resource.ResourceWithImportState = (*serverAgentResource)(nil)

type serverAgentResource struct{ client *hetrixtools.Client }

type serverAgentModel struct {
	ID        types.String `tfsdk:"id"`
	MonitorID types.String `tfsdk:"monitor_id"`
	AgentID   types.String `tfsdk:"agent_id"`
}

func newServerAgentResource() resource.Resource { return &serverAgentResource{} }

func (r *serverAgentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_agent"
}

func (r *serverAgentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Attaches a HetrixTools server monitoring agent to an uptime monitor.",
		Attributes: map[string]schema.Attribute{
			"id":         schema.StringAttribute{Computed: true},
			"monitor_id": schema.StringAttribute{Required: true},
			"agent_id":   schema.StringAttribute{Computed: true},
		},
	}
}

func (r *serverAgentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*hetrixtools.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", "Expected *hetrixtools.Client.")
		return
	}
	r.client = c
}

func (r *serverAgentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serverAgentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	result, err := r.client.AttachServerAgent(ctx, plan.MonitorID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Attach server agent failed", err.Error())
		return
	}
	plan.ID = types.StringValue(plan.MonitorID.ValueString())
	if result.AgentID != nil {
		plan.AgentID = types.StringValue(*result.AgentID)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *serverAgentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serverAgentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	result, err := r.client.GetServerAgent(ctx, state.MonitorID.ValueString())
	if err != nil {
		if hetrixtools.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read server agent failed", err.Error())
		return
	}
	if result.AgentID == nil || *result.AgentID == "" {
		resp.State.RemoveResource(ctx)
		return
	}
	state.ID = types.StringValue(state.MonitorID.ValueString())
	state.AgentID = types.StringValue(*result.AgentID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *serverAgentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan serverAgentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ID = types.StringValue(plan.MonitorID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *serverAgentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state serverAgentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DetachServerAgent(ctx, state.MonitorID.ValueString()); err != nil && !hetrixtools.IsNotFound(err) {
		resp.Diagnostics.AddError("Detach server agent failed", err.Error())
	}
}

func (r *serverAgentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, pathRoot("monitor_id"), req, resp)
}
