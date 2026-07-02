package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	hetrixtools "github.com/xaf/go-hetrixtools/client"
)

var _ resource.Resource = (*scheduledMaintenanceResource)(nil)
var _ resource.ResourceWithConfigure = (*scheduledMaintenanceResource)(nil)
var _ resource.ResourceWithImportState = (*scheduledMaintenanceResource)(nil)

type scheduledMaintenanceResource struct{ client *hetrixtools.Client }

type scheduledMaintenanceModel struct {
	ID                types.String `tfsdk:"id"`
	MonitorID         types.String `tfsdk:"monitor_id"`
	Start             types.String `tfsdk:"start"`
	End               types.String `tfsdk:"end"`
	Timezone          types.String `tfsdk:"timezone"`
	WithNotifications types.Bool   `tfsdk:"with_notifications"`
	Recurring         types.Bool   `tfsdk:"recurring"`
	RecurringTime     types.Int64  `tfsdk:"recurring_time"`
	RecurringTimeType types.String `tfsdk:"recurring_time_type"`
}

func newScheduledMaintenanceResource() resource.Resource { return &scheduledMaintenanceResource{} }

func (r *scheduledMaintenanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scheduled_maintenance"
}

func (r *scheduledMaintenanceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a HetrixTools scheduled maintenance window.",
		Attributes: map[string]schema.Attribute{
			"id":                  schema.StringAttribute{Computed: true},
			"monitor_id":          schema.StringAttribute{Required: true},
			"start":               schema.StringAttribute{Required: true, MarkdownDescription: "Maintenance start time as `YYYY-MM-DD HH:MM`."},
			"end":                 schema.StringAttribute{Required: true, MarkdownDescription: "Maintenance end time as `YYYY-MM-DD HH:MM`."},
			"timezone":            schema.StringAttribute{Required: true},
			"with_notifications":  schema.BoolAttribute{Optional: true, Computed: true},
			"recurring":           schema.BoolAttribute{Computed: true},
			"recurring_time":      schema.Int64Attribute{Optional: true, Computed: true},
			"recurring_time_type": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "One of `hour`, `day`, `week`, `month`, or `year` when recurring_time is set."},
		},
	}
}

func (r *scheduledMaintenanceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *scheduledMaintenanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan scheduledMaintenanceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	created, err := r.client.CreateScheduledMaintenance(ctx, scheduledMaintenanceRequestFromModel(plan))
	if err != nil {
		resp.Diagnostics.AddError("Create scheduled maintenance failed", err.Error())
		return
	}

	setScheduledMaintenanceState(&plan, *created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *scheduledMaintenanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state scheduledMaintenanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	found, err := r.find(ctx, state.ID.ValueString(), state.MonitorID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read scheduled maintenance failed", err.Error())
		return
	}
	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	setScheduledMaintenanceState(&state, *found)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *scheduledMaintenanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state scheduledMaintenanceModel
	var plan scheduledMaintenanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteScheduledMaintenance(ctx, state.ID.ValueString()); err != nil && !hetrixtools.IsNotFound(err) {
		resp.Diagnostics.AddError("Delete old scheduled maintenance failed", err.Error())
		return
	}

	created, err := r.client.CreateScheduledMaintenance(ctx, scheduledMaintenanceRequestFromModel(plan))
	if err != nil {
		resp.Diagnostics.AddError("Recreate scheduled maintenance failed", err.Error())
		return
	}
	setScheduledMaintenanceState(&plan, *created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *scheduledMaintenanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state scheduledMaintenanceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteScheduledMaintenance(ctx, state.ID.ValueString()); err != nil && !hetrixtools.IsNotFound(err) {
		resp.Diagnostics.AddError("Delete scheduled maintenance failed", err.Error())
	}
}

func (r *scheduledMaintenanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, pathRoot("id"), req, resp)
}

func (r *scheduledMaintenanceResource) find(ctx context.Context, id string, monitorID string) (*hetrixtools.ScheduledMaintenance, error) {
	return r.client.GetScheduledMaintenance(ctx, id, monitorID)
}

func setScheduledMaintenanceState(model *scheduledMaintenanceModel, maintenance hetrixtools.ScheduledMaintenance) {
	model.ID = types.StringValue(maintenance.ID)
	model.MonitorID = types.StringValue(maintenance.MonitorID)
	model.Start = types.StringValue(maintenance.Start)
	model.End = types.StringValue(maintenance.End)
	model.Timezone = types.StringValue(maintenance.Timezone)
	model.WithNotifications = types.BoolValue(maintenance.WithNotifications)
	model.Recurring = types.BoolValue(maintenance.Recurring)
	model.RecurringTime = types.Int64Value(maintenance.RecurringTime)
	model.RecurringTimeType = types.StringValue(maintenance.RecurringTimeType)
}

func scheduledMaintenanceRequestFromModel(model scheduledMaintenanceModel) hetrixtools.ScheduledMaintenanceRequest {
	return hetrixtools.ScheduledMaintenanceRequest{
		MonitorID:         model.MonitorID.ValueString(),
		Start:             model.Start.ValueString(),
		End:               model.End.ValueString(),
		Timezone:          model.Timezone.ValueString(),
		WithNotifications: boolValue(model.WithNotifications, false),
		RecurringTime:     int64Value(model.RecurringTime, 0),
		RecurringTimeType: stringValue(model.RecurringTimeType, ""),
	}
}

func boolValue(value types.Bool, fallback bool) bool {
	if value.IsNull() || value.IsUnknown() {
		return fallback
	}
	return value.ValueBool()
}

func int64Value(value types.Int64, fallback int64) int64 {
	if value.IsNull() || value.IsUnknown() {
		return fallback
	}
	return value.ValueInt64()
}

func stringValue(value types.String, fallback string) string {
	if value.IsNull() || value.IsUnknown() {
		return fallback
	}
	return value.ValueString()
}

func pathRoot(name string) path.Path { return path.Root(name) }
