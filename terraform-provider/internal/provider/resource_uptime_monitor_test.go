package provider

import (
	"context"
	"strings"
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

func TestUptimeMonitorModelFromAPIHydratesImportState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	showTarget := false
	trueValue := true
	diagnostics := &testDiagnostics{}
	state := uptimeMonitorModelFromAPI(ctx, uptimeMonitorModel{}, hetrixtools.UptimeMonitor{
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
	if got, want := state.Type.ValueString(), "http"; got != want {
		t.Fatalf("type = %q, want %q", got, want)
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
	if !state.Port.IsNull() {
		t.Fatalf("port = %#v, want null for http monitor", state.Port)
	}
	if !state.ServerID.IsNull() {
		t.Fatalf("server ID = %#v, want null for http monitor", state.ServerID)
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

func TestUptimeMonitorModelFromAPIPreservesSMTPPort(t *testing.T) {
	t.Parallel()

	state := uptimeMonitorModelFromAPI(context.Background(), uptimeMonitorModel{}, hetrixtools.UptimeMonitor{
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

func TestUptimeMonitorModelFromAPIPreservesHeartbeatServerID(t *testing.T) {
	t.Parallel()

	state := uptimeMonitorModelFromAPI(context.Background(), uptimeMonitorModel{}, hetrixtools.UptimeMonitor{
		ID:       "heartbeat-1",
		Type:     "heartbeat",
		Name:     "Heartbeat",
		ServerID: pointerToString("srv-1"),
	}, &testDiagnostics{})
	if got, want := state.ServerID.ValueString(), "srv-1"; got != want {
		t.Fatalf("server ID = %q, want %q", got, want)
	}
}

func TestUptimeMonitorRequestFromModelUsesCanonicalLocationNames(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	diagnostics := &testDiagnostics{}
	locations, diags := types.SetValueFrom(ctx, types.StringType, []string{"new_york", "san_francisco", "dallas", "amsterdam"})
	if diags.HasError() {
		t.Fatalf("build locations: %v", diags)
	}

	request := uptimeMonitorRequestFromModel(ctx, uptimeMonitorModel{
		ID:              types.StringValue("up-1"),
		Type:            types.StringValue("http"),
		Name:            types.StringValue("Homepage"),
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

func TestValidateUptimeMonitorModelRequiresTypeSpecificFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		model      uptimeMonitorModel
		wantErrors []string
	}{
		{
			name:       "http requires target",
			model:      uptimeMonitorModel{Type: types.StringValue("http")},
			wantErrors: []string{"target is required for http uptime monitors"},
		},
		{
			name:       "ping requires target",
			model:      uptimeMonitorModel{Type: types.StringValue("ping")},
			wantErrors: []string{"target is required for ping uptime monitors"},
		},
		{
			name:       "smtp requires target and port",
			model:      uptimeMonitorModel{Type: types.StringValue("smtp")},
			wantErrors: []string{"port is required for smtp uptime monitors", "target is required for smtp uptime monitors"},
		},
		{
			name:       "smtp auth must be paired",
			model:      uptimeMonitorModel{Type: types.StringValue("smtp"), Target: types.StringValue("smtp.example.com"), Port: types.Int64Value(587), SMTPUser: types.StringValue("user")},
			wantErrors: []string{"smtp_user and smtp_password must be set together"},
		},
		{
			name:       "heartbeat rejects target",
			model:      uptimeMonitorModel{Type: types.StringValue("heartbeat"), Target: types.StringValue("https://example.com")},
			wantErrors: []string{"target is not supported for heartbeat uptime monitors"},
		},
		{
			name:       "ping rejects http fields",
			model:      uptimeMonitorModel{Type: types.StringValue("ping"), Target: types.StringValue("example.com"), Keyword: types.StringValue("healthy")},
			wantErrors: []string{"keyword is only supported for http uptime monitors"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			diagnostics := &testDiagnostics{}
			validateUptimeMonitorModel(diagnostics, test.model)
			for _, want := range test.wantErrors {
				if !diagnosticContains(diagnostics.errors, want) {
					t.Fatalf("diagnostics = %#v, want containing %q", diagnostics.errors, want)
				}
			}
		})
	}
}

func diagnosticContains(errors []string, target string) bool {
	for _, err := range errors {
		if strings.Contains(err, target) {
			return true
		}
	}
	return false
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
