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

func TestUptimeMonitorModelFromAPIHydratesImportState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	showTarget := false
	trueValue := true
	falseValue := false
	diagnostics := &testDiagnostics{}
	state := uptimeMonitorModelFromAPI(ctx, uptimeMonitorModel{}, hetrixtools.UptimeMonitor{
		ID:               "up-1",
		Type:             1,
		Name:             "Homepage",
		Target:           "https://example.com",
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
		Locations:        map[string]bool{"ams": true, "nyc": false},
		Grace:            120,
		InfoPublic:       &trueValue,
		CPUPublic:        &falseValue,
		RAMPublic:        &trueValue,
		DiskPublic:       &falseValue,
		NetPublic:        &trueValue,
		ServerID:         "srv-1",
	}, diagnostics)
	if len(diagnostics.errors) > 0 {
		t.Fatalf("unexpected diagnostics: %#v", diagnostics.errors)
	}

	if got, want := state.ID.ValueString(), "up-1"; got != want {
		t.Fatalf("id = %q, want %q", got, want)
	}
	if got, want := state.Type.ValueInt64(), int64(1); got != want {
		t.Fatalf("type = %d, want %d", got, want)
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
	if got, want := state.ServerID.ValueString(), "srv-1"; got != want {
		t.Fatalf("server ID = %q, want %q", got, want)
	}

	locations := map[string]types.Bool{}
	diags := state.Locations.ElementsAs(ctx, &locations, false)
	if diags.HasError() {
		t.Fatalf("decode locations: %v", diags)
	}
	if got, want := locations["ams"].ValueBool(), true; got != want {
		t.Fatalf("ams location = %v, want %v", got, want)
	}
	if got, want := locations["nyc"].ValueBool(), false; got != want {
		t.Fatalf("nyc location = %v, want %v", got, want)
	}
}
