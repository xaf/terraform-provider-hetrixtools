package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	hetrixtools "github.com/xaf/go-hetrixtools/client"
)

var _ resource.Resource = (*blacklistMonitorResource)(nil)
var _ resource.ResourceWithConfigure = (*blacklistMonitorResource)(nil)
var _ resource.ResourceWithImportState = (*blacklistMonitorResource)(nil)

type blacklistMonitorResource struct{ client *hetrixtools.Client }

type blacklistMonitorModel struct {
	ID      types.String `tfsdk:"id"`
	Target  types.String `tfsdk:"target"`
	Label   types.String `tfsdk:"label"`
	Contact types.String `tfsdk:"contact_list_id"`
}

func newBlacklistMonitorResource() resource.Resource { return &blacklistMonitorResource{} }

func (r *blacklistMonitorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blacklist_monitor"
}

func (r *blacklistMonitorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a HetrixTools blacklist monitor.",
		Attributes: map[string]schema.Attribute{
			"id":              schema.StringAttribute{Computed: true},
			"target":          schema.StringAttribute{Required: true, MarkdownDescription: "IP address, IP range, CIDR block, or domain name."},
			"label":           schema.StringAttribute{Optional: true, Computed: true},
			"contact_list_id": schema.StringAttribute{Optional: true, Computed: true},
		},
	}
}

func (r *blacklistMonitorResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *blacklistMonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan blacklistMonitorModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.CreateBlacklistMonitor(ctx, hetrixtools.BlacklistMonitorRequest{Target: plan.Target.ValueString(), Label: stringValue(plan.Label, ""), Contact: stringValue(plan.Contact, "")})
	if err != nil {
		resp.Diagnostics.AddError("Create blacklist monitor failed", err.Error())
		return
	}

	found, err := r.find(ctx, plan.Target.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read created blacklist monitor failed", err.Error())
		return
	}
	if found != nil {
		setBlacklistMonitorState(&plan, *found)
	} else {
		plan.ID = types.StringValue(plan.Target.ValueString())
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *blacklistMonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state blacklistMonitorModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	found, err := r.find(ctx, state.Target.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read blacklist monitor failed", err.Error())
		return
	}
	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	setBlacklistMonitorState(&state, *found)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *blacklistMonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan blacklistMonitorModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	_, err := r.client.UpdateBlacklistMonitor(ctx, hetrixtools.BlacklistMonitorRequest{Target: plan.Target.ValueString(), Label: stringValue(plan.Label, ""), Contact: stringValue(plan.Contact, "")})
	if err != nil {
		resp.Diagnostics.AddError("Update blacklist monitor failed", err.Error())
		return
	}
	found, err := r.find(ctx, plan.Target.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read updated blacklist monitor failed", err.Error())
		return
	}
	if found != nil {
		setBlacklistMonitorState(&plan, *found)
	} else if plan.ID.IsNull() || plan.ID.IsUnknown() {
		plan.ID = types.StringValue(plan.Target.ValueString())
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *blacklistMonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state blacklistMonitorModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.DeleteBlacklistMonitor(ctx, state.Target.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete blacklist monitor failed", err.Error())
	}
}

func (r *blacklistMonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, pathRoot("target"), req, resp)
}

func (r *blacklistMonitorResource) find(ctx context.Context, target string) (*hetrixtools.BlacklistMonitor, error) {
	return r.client.GetBlacklistMonitor(ctx, target)
}

func setBlacklistMonitorState(model *blacklistMonitorModel, monitor hetrixtools.BlacklistMonitor) {
	model.ID = types.StringValue(firstNonEmpty(monitor.ID, monitor.Target))
	model.Target = types.StringValue(monitor.Target)
	model.Label = types.StringValue(firstNonEmpty(monitor.Label, monitor.Name))
	model.Contact = types.StringValue(monitor.Contact)
}

func decodeActionResponse(body []byte) (hetrixtools.ActionResponse, error) {
	var result hetrixtools.ActionResponse
	err := json.Unmarshal(body, &result)
	return result, err
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
