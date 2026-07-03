package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	hetrixtools "github.com/xaf/terraform-provider-hetrixtools/client"
)

var _ resource.Resource = (*warningPoliciesResource)(nil)
var _ resource.ResourceWithConfigure = (*warningPoliciesResource)(nil)
var _ resource.ResourceWithImportState = (*warningPoliciesResource)(nil)

type warningPoliciesResource struct{ client *hetrixtools.Client }

type warningPoliciesModel struct {
	ID        types.String `tfsdk:"id"`
	MonitorID types.String `tfsdk:"monitor_id"`
	Policies  types.String `tfsdk:"policies_json"`
}

func newWarningPoliciesResource() resource.Resource { return &warningPoliciesResource{} }

func (r *warningPoliciesResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_agent_warning_policies"
}

func (r *warningPoliciesResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages server-agent warning policies for an uptime monitor. The policy payload mirrors the HetrixTools API JSON schema.",
		Attributes: map[string]schema.Attribute{
			"id":            schema.StringAttribute{Computed: true},
			"monitor_id":    schema.StringAttribute{Required: true},
			"policies_json": schema.StringAttribute{Required: true, MarkdownDescription: "JSON object sent to HetrixTools to replace the monitor's server-agent warning policies."},
		},
	}
}

func (r *warningPoliciesResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *warningPoliciesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	r.upsert(ctx, req.Plan, &resp.Diagnostics, &resp.State)
}

func (r *warningPoliciesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state warningPoliciesModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	policies, err := r.client.GetServerAgentWarningPolicies(ctx, state.MonitorID.ValueString())
	if err != nil {
		if hetrixtools.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Read warning policies failed", err.Error())
		return
	}
	body, err := json.Marshal(policies)
	if err != nil {
		resp.Diagnostics.AddError("Read warning policies failed", err.Error())
		return
	}
	state.ID = types.StringValue(state.MonitorID.ValueString())
	state.Policies = types.StringValue(normalizeJSON(body))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *warningPoliciesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.upsert(ctx, req.Plan, &resp.Diagnostics, &resp.State)
}

func (r *warningPoliciesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning("Warning policies left in place", "HetrixTools does not expose a delete/reset endpoint for server-agent warning policies. Removing the Terraform resource only forgets it from state.")
}

func (r *warningPoliciesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, pathRoot("monitor_id"), req, resp)
}

func (r *warningPoliciesResource) upsert(ctx context.Context, planGetter interface {
	Get(context.Context, any) diag.Diagnostics
}, diagnostics *diag.Diagnostics, stateSetter interface {
	Set(context.Context, any) diag.Diagnostics
}) {
	var plan warningPoliciesModel
	diagnostics.Append(planGetter.Get(ctx, &plan)...)
	if diagnostics.HasError() {
		return
	}
	var payload any
	if err := json.Unmarshal([]byte(plan.Policies.ValueString()), &payload); err != nil {
		diagnostics.AddError("Invalid policies_json", fmt.Sprintf("policies_json must be a JSON object: %s", err))
		return
	}
	if err := r.client.UpdateServerAgentWarningPolicies(ctx, plan.MonitorID.ValueString(), payload); err != nil {
		diagnostics.AddError("Update warning policies failed", err.Error())
		return
	}
	plan.ID = types.StringValue(plan.MonitorID.ValueString())
	plan.Policies = types.StringValue(normalizeJSON([]byte(plan.Policies.ValueString())))
	diagnostics.Append(stateSetter.Set(ctx, &plan)...)
}
