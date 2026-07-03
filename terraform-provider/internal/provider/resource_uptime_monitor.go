package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	hetrixtools "github.com/xaf/terraform-provider-hetrixtools/client"
)

var _ resource.Resource = (*uptimeMonitorResource)(nil)
var _ resource.ResourceWithConfigure = (*uptimeMonitorResource)(nil)
var _ resource.ResourceWithValidateConfig = (*uptimeMonitorResource)(nil)
var _ resource.ResourceWithImportState = (*uptimeMonitorResource)(nil)

type uptimeMonitorResource struct{ client *hetrixtools.Client }

type uptimeMonitorModel struct {
	ID               types.String `tfsdk:"id"`
	Type             types.String `tfsdk:"type"`
	Name             types.String `tfsdk:"name"`
	Target           types.String `tfsdk:"target"`
	Port             types.Int64  `tfsdk:"port"`
	HTTPMethod       types.String `tfsdk:"http_method"`
	MaxRedirects     types.Int64  `tfsdk:"max_redirects"`
	SMTPUser         types.String `tfsdk:"smtp_user"`
	SMTPPass         types.String `tfsdk:"smtp_password"`
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
	Locations        types.Set    `tfsdk:"locations"`
	Keyword          types.String `tfsdk:"keyword"`
	HTTPCodes        types.List   `tfsdk:"accepted_http_codes"`
	Grace            types.Int64  `tfsdk:"grace"`
	InfoPublic       types.Bool   `tfsdk:"info_public"`
	CPUPublic        types.Bool   `tfsdk:"cpu_public"`
	RAMPublic        types.Bool   `tfsdk:"ram_public"`
	DiskPublic       types.Bool   `tfsdk:"disk_public"`
	NetPublic        types.Bool   `tfsdk:"net_public"`
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
			"type":                   schema.StringAttribute{Required: true, MarkdownDescription: "Monitor type: `http`, `ping`, `smtp`, or `heartbeat`. Changing this forces replacement.", PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()}},
			"name":                   schema.StringAttribute{Required: true},
			"target":                 optionalComputedString(),
			"port":                   schema.Int64Attribute{Optional: true, MarkdownDescription: "Port used for SMTP monitors. Required when `type = \"smtp\"`."},
			"http_method":            optionalComputedString(),
			"max_redirects":          optionalComputedInt64(),
			"smtp_user":              optionalComputedString(),
			"smtp_password":          optionalComputedSensitiveString(),
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
			"locations": schema.SetAttribute{
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "Set of canonical HetrixTools v3 location names enabled for this monitor, e.g. `[\"amsterdam\", \"new_york\"]`.",
				PlanModifiers:       []planmodifier.Set{setplanmodifier.UseStateForUnknown()},
			},
			"keyword":             optionalComputedString(),
			"accepted_http_codes": optionalComputedInt64List(),
			"grace":               optionalComputedInt64(),
			"info_public":         optionalComputedBool(),
			"cpu_public":          optionalComputedBool(),
			"ram_public":          optionalComputedBool(),
			"disk_public":         optionalComputedBool(),
			"net_public":          optionalComputedBool(),
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

func optionalComputedSensitiveString() schema.StringAttribute {
	return schema.StringAttribute{Optional: true, Computed: true, Sensitive: true, PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()}}
}

func optionalComputedInt64() schema.Int64Attribute {
	return schema.Int64Attribute{Optional: true, Computed: true, PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()}}
}

func optionalComputedBool() schema.BoolAttribute {
	return schema.BoolAttribute{Optional: true, Computed: true, PlanModifiers: []planmodifier.Bool{boolplanmodifier.UseStateForUnknown()}}
}

func optionalComputedInt64List() schema.ListAttribute {
	return schema.ListAttribute{Optional: true, Computed: true, ElementType: types.Int64Type, PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()}}
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

func (r *uptimeMonitorResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config uptimeMonitorModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() || config.Type.IsNull() || config.Type.IsUnknown() {
		return
	}
	validateUptimeMonitorModel(&resp.Diagnostics, config)
}

func validateUptimeMonitorModel(diagnostics interface{ AddError(string, string) }, config uptimeMonitorModel) {
	monitorType := config.Type.ValueString()
	switch monitorType {
	case "http", "ping", "smtp", "heartbeat":
	default:
		diagnostics.AddError("Invalid uptime monitor type", "type must be one of: http, ping, smtp, heartbeat")
		return
	}

	if monitorType != "http" {
		addIfSetString(diagnostics, config.HTTPMethod, "http_method is only supported for http uptime monitors")
		addIfSetInt64(diagnostics, config.MaxRedirects, "max_redirects is only supported for http uptime monitors")
		addIfSetString(diagnostics, config.Keyword, "keyword is only supported for http uptime monitors")
		if !config.HTTPCodes.IsNull() && !config.HTTPCodes.IsUnknown() {
			diagnostics.AddError("Invalid uptime monitor configuration", "accepted_http_codes is only supported for http uptime monitors")
		}
	}
	if monitorType != "smtp" {
		addIfSetInt64(diagnostics, config.Port, "port is only supported for smtp uptime monitors")
		addIfSetString(diagnostics, config.SMTPUser, "smtp_user is only supported for smtp uptime monitors")
		addIfSetString(diagnostics, config.SMTPPass, "smtp_password is only supported for smtp uptime monitors")
	}
	if monitorType == "smtp" && (config.Port.IsNull() || config.Port.IsUnknown()) {
		diagnostics.AddError("Invalid uptime monitor configuration", "port is required for smtp uptime monitors")
	}
	if setString(config.SMTPUser) != setString(config.SMTPPass) {
		diagnostics.AddError("Invalid uptime monitor configuration", "smtp_user and smtp_password must be set together")
	}
	if monitorType == "http" || monitorType == "ping" || monitorType == "smtp" {
		if !setString(config.Target) {
			diagnostics.AddError("Invalid uptime monitor configuration", fmt.Sprintf("target is required for %s uptime monitors", monitorType))
		}
	}
	if monitorType != "heartbeat" {
		addIfSetInt64(diagnostics, config.Grace, "grace is only supported for heartbeat uptime monitors")
		addIfSetBool(diagnostics, config.InfoPublic, "info_public is only supported for heartbeat uptime monitors")
		addIfSetBool(diagnostics, config.CPUPublic, "cpu_public is only supported for heartbeat uptime monitors")
		addIfSetBool(diagnostics, config.RAMPublic, "ram_public is only supported for heartbeat uptime monitors")
		addIfSetBool(diagnostics, config.DiskPublic, "disk_public is only supported for heartbeat uptime monitors")
		addIfSetBool(diagnostics, config.NetPublic, "net_public is only supported for heartbeat uptime monitors")
	}
	if monitorType == "heartbeat" {
		addIfSetString(diagnostics, config.Target, "target is not supported for heartbeat uptime monitors")
		addIfSetInt64(diagnostics, config.FailedLocations, "failed_locations is not supported for heartbeat uptime monitors")
		if !config.Locations.IsNull() && !config.Locations.IsUnknown() {
			diagnostics.AddError("Invalid uptime monitor configuration", "locations is not supported for heartbeat uptime monitors")
		}
	}
	if monitorType != "http" && monitorType != "smtp" {
		addIfSetBool(diagnostics, config.VerSSLCert, "verify_ssl_certificate is only supported for http and smtp uptime monitors")
		addIfSetBool(diagnostics, config.VerSSLHost, "verify_ssl_host is only supported for http and smtp uptime monitors")
	}
}

func addIfSetString(diagnostics interface{ AddError(string, string) }, value types.String, message string) {
	if setString(value) {
		diagnostics.AddError("Invalid uptime monitor configuration", message)
	}
}

func addIfSetInt64(diagnostics interface{ AddError(string, string) }, value types.Int64, message string) {
	if !value.IsNull() && !value.IsUnknown() {
		diagnostics.AddError("Invalid uptime monitor configuration", message)
	}
}

func addIfSetBool(diagnostics interface{ AddError(string, string) }, value types.Bool, message string) {
	if !value.IsNull() && !value.IsUnknown() {
		diagnostics.AddError("Invalid uptime monitor configuration", message)
	}
}

func setString(value types.String) bool {
	return !value.IsNull() && !value.IsUnknown() && value.ValueString() != ""
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
	state.Type = types.StringValue(monitor.Type)
	state.Name = types.StringValue(monitor.Name)
	state.Target = types.StringValue(monitor.Target)
	state.Port = types.Int64Value(monitor.Port)
	state.HTTPMethod = stringNullIfEmpty(monitor.HTTPMethod)
	state.MaxRedirects = types.Int64Value(monitor.MaxRedirects)
	state.SMTPUser = stringNullIfEmpty(monitor.SMTPUser)
	if state.SMTPPass.IsNull() || state.SMTPPass.IsUnknown() {
		state.SMTPPass = types.StringNull()
	}
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
	state.Keyword = stringNullIfEmpty(monitor.Keyword)
	state.HTTPCodes = int64ListFromAPI(ctx, state.HTTPCodes, monitor.HTTPCodes, diagnostics)
	state.Grace = types.Int64Value(monitor.Grace)
	state.InfoPublic = boolFromPointer(monitor.InfoPublic)
	state.CPUPublic = boolFromPointer(monitor.CPUPublic)
	state.RAMPublic = boolFromPointer(monitor.RAMPublic)
	state.DiskPublic = boolFromPointer(monitor.DiskPublic)
	state.NetPublic = boolFromPointer(monitor.NetPublic)
	state.ServerID = stringNullIfEmpty(monitor.ServerID)

	if monitor.Locations == nil {
		state.Locations = types.SetNull(types.StringType)
	} else {
		locations, diags := types.SetValueFrom(ctx, types.StringType, monitor.Locations)
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

func int64ListFromAPI(ctx context.Context, current types.List, values []int64, diagnostics interface{ AddError(string, string) }) types.List {
	if len(values) == 0 {
		if current.IsNull() || current.IsUnknown() {
			return types.ListNull(types.Int64Type)
		}
		return current
	}
	list, diags := types.ListValueFrom(ctx, types.Int64Type, values)
	if diags.HasError() {
		diagnostics.AddError("Invalid uptime monitor integer list", fmt.Sprintf("Could not encode integer list: %v", diags))
		return current
	}
	return list
}

func uptimeMonitorRequestFromModel(ctx context.Context, model uptimeMonitorModel, diagnostics interface{ AddError(string, string) }) hetrixtools.UptimeMonitorRequest {
	var locationValues []string
	if !model.Locations.IsNull() && !model.Locations.IsUnknown() {
		var locations []types.String
		if diags := model.Locations.ElementsAs(ctx, &locations, false); diags.HasError() {
			diagnostics.AddError("Invalid locations", fmt.Sprintf("Could not decode locations: %v", diags))
			return hetrixtools.UptimeMonitorRequest{}
		}
		for _, location := range locations {
			if !location.IsNull() && !location.IsUnknown() {
				locationValues = append(locationValues, location.ValueString())
			}
		}
	}

	var httpCodes []int64
	if !model.HTTPCodes.IsNull() && !model.HTTPCodes.IsUnknown() {
		if diags := model.HTTPCodes.ElementsAs(ctx, &httpCodes, false); diags.HasError() {
			diagnostics.AddError("Invalid accepted_http_codes", fmt.Sprintf("Could not decode accepted_http_codes: %v", diags))
			return hetrixtools.UptimeMonitorRequest{}
		}
	}

	return hetrixtools.UptimeMonitorRequest{
		MID:              stringValue(model.ID, ""),
		Type:             model.Type.ValueString(),
		Name:             model.Name.ValueString(),
		Target:           stringValue(model.Target, ""),
		Port:             int64Value(model.Port, 0),
		HTTPMethod:       stringValue(model.HTTPMethod, ""),
		MaxRedirects:     int64Value(model.MaxRedirects, 0),
		SMTPUser:         stringValue(model.SMTPUser, ""),
		SMTPPass:         stringValue(model.SMTPPass, ""),
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
		Keyword:          stringValue(model.Keyword, ""),
		HTTPCodes:        httpCodes,
		Grace:            int64Value(model.Grace, 0),
		INFOPub:          boolPointer(model.InfoPublic),
		CPUPub:           boolPointer(model.CPUPublic),
		RAMPub:           boolPointer(model.RAMPublic),
		DISKPub:          boolPointer(model.DiskPublic),
		NETPub:           boolPointer(model.NetPublic),
	}
}

func boolPointer(value types.Bool) *bool {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	result := value.ValueBool()
	return &result
}
