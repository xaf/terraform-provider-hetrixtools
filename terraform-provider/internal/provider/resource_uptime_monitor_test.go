package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	hetrixtools "github.com/xaf/terraform-provider-hetrixtools/client"
)

type testDiagnostics struct{ errors []string }

func (d *testDiagnostics) AddError(summary string, detail string) {
	d.errors = append(d.errors, summary+": "+detail)
}

func pointerToInt64(value int64) *int64 { return &value }

func pointerToString(value string) *string { return &value }

func TestUptimeHTTPMonitorModelFromAPIHydratesImportState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	showTarget := false
	trueValue := true
	diagnostics := &testDiagnostics{}
	state := uptimeHTTPMonitorModelFromAPI(ctx, uptimeHTTPMonitorModel{}, hetrixtools.UptimeMonitor{
		ID:               "up-1",
		Type:             "http",
		Name:             "Homepage",
		Target:           "https://example.com",
		HTTPMethod:       "GET",
		MaxRedirects:     5,
		Timeout:          10,
		Frequency:        60,
		FailsBeforeAlert: 3,
		FailedLocations:  2,
		ContactListID:    "contacts-1",
		Category:         "prod",
		AlertAfter:       "1m",
		RepeatTimes:      5,
		RepeatEvery:      "1h",
		Public:           &trueValue,
		ShowTarget:       &showTarget,
		VerSSLCert:       &trueValue,
		VerSSLHost:       &trueValue,
		Locations:        []string{"amsterdam", "new_york"},
		Keyword:          "healthy",
		HTTPCodes:        []int64{200, 204},
	}, diagnostics)
	if len(diagnostics.errors) > 0 {
		t.Fatalf("unexpected diagnostics: %#v", diagnostics.errors)
	}

	if got, want := state.ID.ValueString(), "up-1"; got != want {
		t.Fatalf("id = %q, want %q", got, want)
	}
	if got, want := state.ContactList.ValueString(), "contacts-1"; got != want {
		t.Fatalf("contact list = %q, want %q", got, want)
	}
	if got, want := state.ShowTarget.ValueBool(), false; got != want {
		t.Fatalf("show target = %v, want %v", got, want)
	}
	if got, want := state.VerSSLCert.ValueBool(), true; got != want {
		t.Fatalf("verify ssl certificate = %v, want %v", got, want)
	}
	if got, want := state.Keyword.ValueString(), "healthy"; got != want {
		t.Fatalf("keyword = %q, want %q", got, want)
	}
	if got, want := state.HTTPMethod.ValueString(), "GET"; got != want {
		t.Fatalf("http_method = %q, want %q", got, want)
	}
	if got, want := state.MaxRedirects.ValueInt64(), int64(5); got != want {
		t.Fatalf("max_redirects = %d, want %d", got, want)
	}
	var httpCodes []int64
	if diags := state.HTTPCodes.ElementsAs(ctx, &httpCodes, false); diags.HasError() {
		t.Fatalf("decode accepted_http_codes: %v", diags)
	}
	if len(httpCodes) != 2 || httpCodes[0] != 200 || httpCodes[1] != 204 {
		t.Fatalf("accepted_http_codes = %#v, want [200 204]", httpCodes)
	}

	var locations []string
	diags := state.Locations.ElementsAs(ctx, &locations, false)
	if diags.HasError() {
		t.Fatalf("decode locations: %v", diags)
	}
	if got, want := locations, []string{"amsterdam", "new_york"}; len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("locations = %#v, want %#v", got, want)
	}
}

func TestUptimeSMTPMonitorModelFromAPIPreservesPort(t *testing.T) {
	t.Parallel()

	state := uptimeSMTPMonitorModelFromAPI(context.Background(), uptimeSMTPMonitorModel{}, hetrixtools.UptimeMonitor{
		ID:     "smtp-1",
		Type:   "smtp",
		Name:   "SMTP",
		Target: "smtp.example.com",
		Port:   pointerToInt64(587),
	}, &testDiagnostics{})
	if got, want := state.Port.ValueInt64(), int64(587); got != want {
		t.Fatalf("port = %d, want %d", got, want)
	}
}

func TestUptimeHeartbeatMonitorModelFromAPIPreservesServerID(t *testing.T) {
	t.Parallel()

	state := uptimeHeartbeatMonitorModelFromAPI(uptimeHeartbeatMonitorModel{}, hetrixtools.UptimeMonitor{
		ID:       "heartbeat-1",
		Type:     "heartbeat",
		Name:     "Heartbeat",
		ServerID: pointerToString("srv-1"),
	})
	if got, want := state.ServerID.ValueString(), "srv-1"; got != want {
		t.Fatalf("server ID = %q, want %q", got, want)
	}
}

func TestUptimeHTTPMonitorRequestFromModelUsesCanonicalLocationNames(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	diagnostics := &testDiagnostics{}
	locations, diags := types.SetValueFrom(ctx, types.StringType, []string{"new_york", "san_francisco", "dallas", "amsterdam"})
	if diags.HasError() {
		t.Fatalf("build locations: %v", diags)
	}

	request := uptimeHTTPMonitorRequestFromModel(ctx, uptimeHTTPMonitorModel{
		uptimeCommonModel: uptimeCommonModel{
			ID:   types.StringValue("up-1"),
			Name: types.StringValue("Homepage"),
		},
		Target:          types.StringValue("https://example.com"),
		HTTPMethod:      types.StringValue("GET"),
		MaxRedirects:    types.Int64Value(5),
		FailedLocations: types.Int64Value(3),
		Locations:       locations,
		Keyword:         types.StringValue("healthy"),
	}, diagnostics)
	if len(diagnostics.errors) > 0 {
		t.Fatalf("unexpected diagnostics: %#v", diagnostics.errors)
	}

	for _, location := range []string{"new_york", "san_francisco", "dallas", "amsterdam"} {
		if !containsString(request.Locations, location) {
			t.Fatalf("location %q missing from request: %#v", location, request.Locations)
		}
	}
	for _, location := range []string{"nyc", "sfo", "dal", "ams"} {
		if containsString(request.Locations, location) {
			t.Fatalf("v2 location code %q leaked into provider request model: %#v", location, request.Locations)
		}
	}
	if got, want := request.Type, "http"; got != want {
		t.Fatalf("type = %q, want %q", got, want)
	}
	if got, want := request.Keyword, "healthy"; got != want {
		t.Fatalf("keyword = %q, want %q", got, want)
	}
	if got, want := request.HTTPMethod, "GET"; got != want {
		t.Fatalf("http_method = %q, want %q", got, want)
	}
	if got, want := request.MaxRedirects, int64(5); got != want {
		t.Fatalf("max_redirects = %d, want %d", got, want)
	}
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
