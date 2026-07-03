package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	hetrixtools "github.com/xaf/terraform-provider-hetrixtools/client"
)

var _ resource.Resource = (*uptimeMonitorResource)(nil)
var _ resource.ResourceWithConfigure = (*uptimeMonitorResource)(nil)
var _ resource.ResourceWithImportState = (*uptimeMonitorResource)(nil)

type uptimeMonitorResource struct{ client *hetrixtools.Client }

type uptimeMonitorModel struct {
	ID               types.String `tfsdk:"id"`
	Type             types.Int64  `tfsdk:"type"`
	Name             types.String `tfsdk:"name"`
	Target           types.String `tfsdk:"target"`
	Timeout          types.Int64  `tfsdk:"timeout"`
	Frequency        types.Int64  `tfsdk:"frequency"`
	FailsBeforeAlert types.Int64  `tfsdk:"fails_before_alert"`
	FailedLocations  types.Int64  `tfsdk:"failed_locations"`
	ContactList      types.String `tfsdk:"contact_list_id"`
	Category         types.String `tfsdk:"category"`
	AlertAfter       types.String `tfsdk:"alert_after"`
	RepeatTimes      types.Int64  `tfsdk:"repeat_times"`
	RepeatEvery      types.String `tfsdk:"repeat_every"`
	Public           types.Bool   `tfsdk:"public"`
	ShowTarget       types.Bool   `tfsdk:"show_target"`
	VerSSLCert       types.Bool   `tfsdk:"verify_ssl_certificate"`
	VerSSLHost       types.Bool   `tfsdk:"verify_ssl_host"`
	Locations        types.Map    `tfsdk:"locations"`
	Grace            types.Int64  `tfsdk:"grace"`
	InfoPublic       types.Bool   `tfsdk:"info_public"`
	CPUPublic        types.Bool   `tfsdk:"cpu_public"`
	RAMPublic        types.Bool   `tfsdk:"ram_public"`
	DiskPublic       types.Bool   `tfsdk:"disk_public"`
	NetPublic        types.Bool   `tfsdk:"net_public"`
	ExtraJSON        types.String `tfsdk:"extra_json"`
	ServerID         types.String `tfsdk:"server_id"`
}

func newUptimeMonitorResource() resource.Resource { return &uptimeMonitorResource{} }

func (r *uptimeMonitorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_uptime_monitor"
}

func (r *uptimeMonitorResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a HetrixTools uptime monitor.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"type":                   schema.Int64Attribute{Required: true, MarkdownDescription: "Monitor type: 1 website, 2 ping/service, 3 SMTP, 9 server agent."},
			"name":                   schema.StringAttribute{Required: true},
			"target":                 optionalComputedString(),
			"timeout":                optionalComputedInt64(),
			"frequency":              optionalComputedInt64(),
			"fails_before_alert":     optionalComputedInt64(),
			"failed_locations":       optionalComputedInt64(),
			"contact_list_id":        optionalComputedString(),
			"category":               optionalComputedString(),
			"alert_after":            optionalComputedString(),
			"repeat_times":           optionalComputedInt64(),
			"repeat_every":           optionalComputedString(),
			"public":                 optionalComputedBool(),
			"show_target":            optionalComputedBool(),
			"verify_ssl_certificate": optionalComputedBool(),
			"verify_ssl_host":        optionalComputedBool(),
			"locations": schema.MapAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.BoolType,
				MarkdownDescription: "Map of HetrixTools location code to enabled flag, e.g. `{ ams = true, nyc = false }`.",
				PlanModifiers:       []planmodifier.Map{mapplanmodifier.UseStateForUnknown()},
			},
			"grace":       optionalComputedInt64(),
			"info_public": optionalComputedBool(),
			"cpu_public":  optionalComputedBool(),
			"ram_public":  optionalComputedBool(),
			"disk_public": optionalComputedBool(),
			"net_public":  optionalComputedBool(),
			"extra_json": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "Additional JSON fields merged into the uptime monitor payload for type-specific options like Method, Keyword, HTTPCodes, SMTPUser, or SMTPPass.",
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"server_id": schema.StringAttribute{
				Computed:      true,
				Sensitive:     true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
		},
	}
}

func optionalComputedString() schema.StringAttribute {
	return schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}}
}

func optionalComputedInt64() schema.Int64Attribute {
	return schema.Int64Attribute{Optional: true, Computed: true, PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()}}
}

func optionalComputedBool() schema.BoolAttribute {
	return schema.BoolAttribute{Optional: true, Computed: true, PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()}}
}

func (r *uptimeMonitorResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *uptimeMonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan uptimeMonitorModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	action, err := r.client.CreateUptimeMonitor(ctx, uptimeMonitorRequestFromModel(ctx, plan, &resp.Diagnostics))
	if resp.Diagnostics.HasError() {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Create uptime monitor failed", err.Error())
		return
	}
	if action != nil {
		if action.MonitorID != "" {
			plan.ID = types.StringValue(action.MonitorID)
		}
		if action.ServerID != "" {
			plan.ServerID = types.StringValue(action.ServerID)
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *uptimeMonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state uptimeMonitorModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	found, err := r.find(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Read uptime monitor failed", err.Error())
		return
	}
	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}
	state = uptimeMonitorModelFromAPI(ctx, state, *found, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *uptimeMonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan uptimeMonitorModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	action, err := r.client.UpdateUptimeMonitor(ctx, uptimeMonitorRequestFromModel(ctx, plan, &resp.Diagnostics))
	if resp.Diagnostics.HasError() {
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Update uptime monitor failed", err.Error())
		return
	}
	if action != nil && action.ServerID != "" {
		plan.ServerID = types.StringValue(action.ServerID)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *uptimeMonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state uptimeMonitorModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.DeleteUptimeMonitor(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Delete uptime monitor failed", err.Error())
	}
}

func (r *uptimeMonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, pathRoot("id"), req, resp)
}

func (r *uptimeMonitorResource) find(ctx context.Context, id string) (*hetrixtools.UptimeMonitor, error) {
	return r.client.GetUptimeMonitor(ctx, id)
}

func uptimeMonitorModelFromAPI(ctx context.Context, state uptimeMonitorModel, monitor hetrixtools.UptimeMonitor, diagnostics interface{ AddError(string, string) }) uptimeMonitorModel {
	state.ID = types.StringValue(monitor.ID)
	state.Type = types.Int64Value(monitor.Type)
	state.Name = types.StringValue(monitor.Name)
	state.Target = types.StringValue(monitor.Target)
	state.Timeout = types.Int64Value(monitor.Timeout)
	state.Frequency = types.Int64Value(monitor.Frequency)
	state.FailsBeforeAlert = types.Int64Value(monitor.FailsBeforeAlert)
	state.FailedLocations = types.Int64Value(monitor.FailedLocations)
	state.ContactList = stringNullIfEmpty(monitor.ContactListID)
	state.Category = stringNullIfEmpty(monitor.Category)
	state.AlertAfter = stringNullIfEmpty(monitor.AlertAfter)
	state.RepeatTimes = types.Int64Value(monitor.RepeatTimes)
	state.RepeatEvery = stringNullIfEmpty(monitor.RepeatEvery)
	state.Public = boolFromPointer(monitor.Public)
	state.ShowTarget = boolFromPointer(monitor.ShowTarget)
	state.VerSSLCert = boolFromPointer(monitor.VerSSLCert)
	state.VerSSLHost = boolFromPointer(monitor.VerSSLHost)
	state.Grace = types.Int64Value(monitor.Grace)
	state.InfoPublic = boolFromPointer(monitor.InfoPublic)
	state.CPUPublic = boolFromPointer(monitor.CPUPublic)
	state.RAMPublic = boolFromPointer(monitor.RAMPublic)
	state.DiskPublic = boolFromPointer(monitor.DiskPublic)
	state.NetPublic = boolFromPointer(monitor.NetPublic)
	state.ServerID = stringNullIfEmpty(monitor.ServerID)
	state.ExtraJSON = extraJSONFromAPI(state.ExtraJSON, monitor.Extra, diagnostics)

	if monitor.Locations == nil {
		state.Locations = types.MapNull(types.BoolType)
	} else {
		locations, diags := types.MapValueFrom(ctx, types.BoolType, monitor.Locations)
		if diags.HasError() {
			diagnostics.AddError("Invalid uptime monitor locations", fmt.Sprintf("Could not encode locations: %v", diags))
			return state
		}
		state.Locations = locations
	}
	return state
}

func boolFromPointer(value *bool) types.Bool {
	if value == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*value)
}

func stringNullIfEmpty(value string) types.String {
	if value == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}

func extraJSONFromAPI(current types.String, extra map[string]any, diagnostics interface{ AddError(string, string) }) types.String {
	if len(extra) == 0 {
		if current.IsNull() || current.IsUnknown() {
			return types.StringNull()
		}
		return current
	}

	merged := map[string]any{}
	if !current.IsNull() && !current.IsUnknown() {
		if err := json.Unmarshal([]byte(current.ValueString()), &merged); err != nil {
			diagnostics.AddError("Invalid uptime monitor extra_json", fmt.Sprintf("Could not decode extra_json: %v", err))
			return current
		}
	}
	for key, value := range extra {
		merged[key] = value
	}

	body, err := json.Marshal(merged)
	if err != nil {
		diagnostics.AddError("Invalid uptime monitor extra fields", fmt.Sprintf("Could not encode extra fields: %v", err))
		return current
	}
	return types.StringValue(string(body))
}

func uptimeMonitorRequestFromModel(ctx context.Context, model uptimeMonitorModel, diagnostics interface{ AddError(string, string) }) hetrixtools.UptimeMonitorRequest {
	locations := map[string]types.Bool{}
	if !model.Locations.IsNull() && !model.Locations.IsUnknown() {
		if diags := model.Locations.ElementsAs(ctx, &locations, false); diags.HasError() {
			diagnostics.AddError("Invalid locations", fmt.Sprintf("Could not decode locations: %v", diags))
			return hetrixtools.UptimeMonitorRequest{}
		}
	}
	locationValues := map[string]bool{}
	for key, value := range locations {
		if !value.IsNull() && !value.IsUnknown() {
			locationValues[key] = value.ValueBool()
		}
	}

	extra := map[string]any{}
	if !model.ExtraJSON.IsNull() && !model.ExtraJSON.IsUnknown() && model.ExtraJSON.ValueString() != "" {
		if err := json.Unmarshal([]byte(model.ExtraJSON.ValueString()), &extra); err != nil {
			diagnostics.AddError("Invalid extra_json", "extra_json must be a JSON object: "+err.Error())
			return hetrixtools.UptimeMonitorRequest{}
		}
	}

	return hetrixtools.UptimeMonitorRequest{
		MID:              stringValue(model.ID, ""),
		Type:             int64Value(model.Type, 0),
		Name:             model.Name.ValueString(),
		Target:           stringValue(model.Target, ""),
		Timeout:          int64Value(model.Timeout, 0),
		Frequency:        int64Value(model.Frequency, 0),
		FailsBeforeAlert: int64Value(model.FailsBeforeAlert, 0),
		FailedLocations:  int64Value(model.FailedLocations, 0),
		ContactList:      stringValue(model.ContactList, ""),
		Category:         stringValue(model.Category, ""),
		AlertAfter:       stringValue(model.AlertAfter, ""),
		RepeatTimes:      int64Value(model.RepeatTimes, 0),
		RepeatEvery:      stringValue(model.RepeatEvery, ""),
		Public:           boolPointer(model.Public),
		ShowTarget:       boolPointer(model.ShowTarget),
		VerSSLCert:       boolPointer(model.VerSSLCert),
		VerSSLHost:       boolPointer(model.VerSSLHost),
		Locations:        locationValues,
		Grace:            int64Value(model.Grace, 0),
		INFOPub:          boolPointer(model.InfoPublic),
		CPUPub:           boolPointer(model.CPUPublic),
		RAMPub:           boolPointer(model.RAMPublic),
		DISKPub:          boolPointer(model.DiskPublic),
		NETPub:           boolPointer(model.NetPublic),
		Extra:            extra,
	}
}

func boolPointer(value types.Bool) *bool {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	result := value.ValueBool()
	return &result
}
