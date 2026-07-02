package provider

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestProviderRegistersExpectedDataSources(t *testing.T) {
	t.Parallel()

	p := &hetrixToolsProvider{}
	dataSources := p.DataSources(context.Background())
	got := make(map[string]bool, len(dataSources))
	for _, factory := range dataSources {
		ds := factory()
		var resp datasource.MetadataResponse
		ds.Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "hetrixtools"}, &resp)
		got[resp.TypeName] = true
	}

	want := []string{
		"hetrixtools_account_limits",
		"hetrixtools_blacklists",
		"hetrixtools_contact_lists",
		"hetrixtools_blacklist_monitors",
		"hetrixtools_blacklist_report",
		"hetrixtools_uptime_monitors",
		"hetrixtools_uptime_report",
		"hetrixtools_uptime_downtimes",
		"hetrixtools_uptime_location_fail_log",
		"hetrixtools_uptime_server_agent",
		"hetrixtools_uptime_server_agent_warning_policies",
		"hetrixtools_status_pages",
		"hetrixtools_scheduled_maintenances",
	}
	for _, name := range want {
		if !got[name] {
			t.Fatalf("missing data source %s; got %#v", name, got)
		}
	}
	if got["hetrixtools_endpoint"] {
		t.Fatal("generic endpoint data source should not be registered")
	}
}

func TestProviderDoesNotCallRawClientHelpersOrHardCodeAPIVersionPaths(t *testing.T) {
	t.Parallel()

	root := "."
	files, err := filepath.Glob(filepath.Join(root, "*.go"))
	if err != nil {
		t.Fatalf("glob provider files: %s", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		body, err := os.ReadFile(file)
		if err != nil {
			t.Fatalf("read %s: %s", file, err)
		}
		content := string(body)
		for _, forbidden := range []string{"GetEndpoint", "getEndpoint", "DoV", "doV", "GetJSON", "PostJSON", "PutJSON", "DeleteJSON", "newEndpointDataSource"} {
			if strings.Contains(content, forbidden) {
				t.Fatalf("%s contains forbidden raw client usage %q", file, forbidden)
			}
		}
		for _, forbiddenPath := range []string{"/v1/", "/v2/", "/v3/", "/account/", "/blacklist", "/uptime", "/status", "/schedule-maintenance"} {
			if strings.Contains(content, forbiddenPath) {
				t.Fatalf("%s contains hard-coded API path %q", file, forbiddenPath)
			}
		}
	}
}
