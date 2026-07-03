package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
			"id":                     schema.StringAttribute{Computed: true},
			"type":                   schema.Int64Attribute{Required: true, MarkdownDescription: "Monitor type: 1 website, 2 ping/service, 3 SMTP, 9 server agent."},
			"name":                   schema.StringAttribute{Required: true},
			"target":                 schema.StringAttribute{Optional: true, Computed: true},
			"timeout":                schema.Int64Attribute{Optional: true, Computed: true},
			"frequency":              schema.Int64Attribute{Optional: true, Computed: true},
			"fails_before_alert":     schema.Int64Attribute{Optional: true, Computed: true},
			"failed_locations":       schema.Int64Attribute{Optional: true, Computed: true},
			"contact_list_id":        schema.StringAttribute{Optional: true, Computed: true},
			"category":               schema.StringAttribute{Optional: true, Computed: true},
			"alert_after":            schema.StringAttribute{Optional: true, Computed: true},
			"repeat_times":           schema.Int64Attribute{Optional: true, Computed: true},
			"repeat_every":           schema.StringAttribute{Optional: true, Computed: true},
			"public":                 schema.BoolAttribute{Optional: true, Computed: true},
			"show_target":            schema.BoolAttribute{Optional: true, Computed: true},
			"verify_ssl_certificate": schema.BoolAttribute{Optional: true, Computed: true},
			"verify_ssl_host":        schema.BoolAttribute{Optional: true, Computed: true},
			"locations":              schema.MapAttribute{Optional: true, ElementType: types.BoolType, MarkdownDescription: "Map of HetrixTools location code to enabled flag, e.g. `{ ams = true, nyc = false }`."},
			"grace":                  schema.Int64Attribute{Optional: true, Computed: true},
			"info_public":            schema.BoolAttribute{Optional: true, Computed: true},
			"cpu_public":             schema.BoolAttribute{Optional: true, Computed: true},
			"ram_public":             schema.BoolAttribute{Optional: true, Computed: true},
			"disk_public":            schema.BoolAttribute{Optional: true, Computed: true},
			"net_public":             schema.BoolAttribute{Optional: true, Computed: true},
			"extra_json":             schema.StringAttribute{Optional: true, Sensitive: true, MarkdownDescription: "Additional JSON fields merged into the uptime monitor payload for type-specific options like Method, Keyword, HTTPCodes, SMTPUser, or SMTPPass."},
			"server_id":              schema.StringAttribute{Computed: true, Sensitive: true},
		},
	}
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
	state.Name = types.StringValue(found.Name)
	state.Target = types.StringValue(found.Target)
	state.Category = types.StringValue(found.Category)
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
