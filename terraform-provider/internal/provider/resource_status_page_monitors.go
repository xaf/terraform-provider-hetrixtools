package provider

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	hetrixtools "github.com/xaf/terraform-provider-hetrixtools/client"
)

var _ resource.Resource = (*statusPageMonitorsResource)(nil)
var _ resource.ResourceWithConfigure = (*statusPageMonitorsResource)(nil)
var _ resource.ResourceWithImportState = (*statusPageMonitorsResource)(nil)

type statusPageMonitorsResource struct{ client *hetrixtools.Client }

type statusPageMonitorsModel struct {
	ID           types.String `tfsdk:"id"`
	StatusPageID types.String `tfsdk:"status_page_id"`
	MonitorIDs   types.Set    `tfsdk:"monitor_ids"`
}

func newStatusPageMonitorsResource() resource.Resource { return &statusPageMonitorsResource{} }

func (r *statusPageMonitorsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_status_page_monitors"
}

func (r *statusPageMonitorsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages the exact monitor set on a HetrixTools status page.",
		Attributes: map[string]schema.Attribute{
			"id":             schema.StringAttribute{Computed: true},
			"status_page_id": schema.StringAttribute{Required: true},
			"monitor_ids":    schema.SetAttribute{Required: true, ElementType: types.StringType},
		},
	}
}

func (r *statusPageMonitorsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *statusPageMonitorsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan statusPageMonitorsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	desired := stringSetValues(ctx, plan.MonitorIDs, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if len(desired) > 0 {
		if err := r.client.AddStatusPageMonitors(ctx, plan.StatusPageID.ValueString(), desired); err != nil {
			resp.Diagnostics.AddError("Add status page monitors failed", err.Error())
			return
		}
	}
	plan.ID = types.StringValue(plan.StatusPageID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *statusPageMonitorsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state statusPageMonitorsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	page, err := r.findStatusPage(ctx, state.StatusPageID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read status page failed", err.Error())
		return
	}
	if page == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	state.ID = types.StringValue(state.StatusPageID.ValueString())
	state.MonitorIDs = setFromStrings(ctx, page.Monitors, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *statusPageMonitorsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan statusPageMonitorsModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	desired := stringSetValues(ctx, plan.MonitorIDs, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	page, err := r.findStatusPage(ctx, plan.StatusPageID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read status page failed", err.Error())
		return
	}
	if page == nil {
		resp.Diagnostics.AddError("Status page not found", "HetrixTools did not return the configured status page.")
		return
	}

	add, remove := diffStrings(page.Monitors, desired)
	if len(add) > 0 {
		if err := r.client.AddStatusPageMonitors(ctx, plan.StatusPageID.ValueString(), add); err != nil {
			resp.Diagnostics.AddError("Add status page monitors failed", err.Error())
			return
		}
	}
	if len(remove) > 0 {
		if err := r.client.RemoveStatusPageMonitors(ctx, plan.StatusPageID.ValueString(), remove); err != nil {
			resp.Diagnostics.AddError("Remove status page monitors failed", err.Error())
			return
		}
	}
	plan.ID = types.StringValue(plan.StatusPageID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *statusPageMonitorsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state statusPageMonitorsModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	current := stringSetValues(ctx, state.MonitorIDs, &resp.Diagnostics)
	if resp.Diagnostics.HasError() || len(current) == 0 {
		return
	}
	if err := r.client.RemoveStatusPageMonitors(ctx, state.StatusPageID.ValueString(), current); err != nil && !hetrixtools.IsNotFound(err) {
		resp.Diagnostics.AddError("Remove status page monitors failed", err.Error())
	}
}

func (r *statusPageMonitorsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, pathRoot("status_page_id"), req, resp)
}

func (r *statusPageMonitorsResource) findStatusPage(ctx context.Context, id string) (*hetrixtools.StatusPage, error) {
	return r.client.GetStatusPage(ctx, id)
}

func diffStrings(current []string, desired []string) ([]string, []string) {
	currentSet := make(map[string]struct{}, len(current))
	desiredSet := make(map[string]struct{}, len(desired))
	for _, value := range current {
		currentSet[value] = struct{}{}
	}
	for _, value := range desired {
		desiredSet[value] = struct{}{}
	}
	var add []string
	var remove []string
	for value := range desiredSet {
		if _, ok := currentSet[value]; !ok {
			add = append(add, value)
		}
	}
	for value := range currentSet {
		if _, ok := desiredSet[value]; !ok {
			remove = append(remove, value)
		}
	}
	sort.Strings(add)
	sort.Strings(remove)
	return add, remove
}
