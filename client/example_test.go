package hetrixtools_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	hetrixtools "github.com/xaf/terraform-provider-hetrixtools/client"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func exampleClient(responses map[string]string) (*hetrixtools.Client, func()) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Method + " " + r.URL.Path
		if r.URL.RawQuery != "" {
			key += "?" + r.URL.RawQuery
		}
		if response, ok := responses[key]; ok {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(response))
			return
		}
		http.Error(w, key, http.StatusNotFound)
	}))
	client := hetrixtools.NewClientWithBaseURL(server.URL, "token", hetrixtools.WithMinimumRequestInterval(0))
	return client, server.Close
}

func ExampleNewClient() {
	client := hetrixtools.NewClient("token", hetrixtools.WithMinimumRequestInterval(0))
	fmt.Println(client != nil)
	// Output: true
}

func ExampleNewClientWithBaseURL() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v3/account/limits": `{"api_calls_left":42}`,
	})
	defer closeServer()

	limits, _ := client.GetAccountLimits(context.Background())
	fmt.Println(limits.(map[string]any)["api_calls_left"])
	// Output: 42
}

func ExampleWithHTTPClient() {
	httpClient := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`{"api_calls_left":7}`)),
		}, nil
	})}
	client := hetrixtools.NewClient("token", hetrixtools.WithHTTPClient(httpClient), hetrixtools.WithMinimumRequestInterval(0))

	limits, _ := client.GetAccountLimits(context.Background())
	fmt.Println(limits.(map[string]any)["api_calls_left"])
	// Output: 7
}

func ExampleWithMinimumRequestInterval() {
	// Zero disables client-side pacing and is intended for tests, examples, or
	// callers that provide their own rate limiter.
	client := hetrixtools.NewClient("token", hetrixtools.WithMinimumRequestInterval(0))
	fmt.Println(client != nil)
	// Output: true
}

func ExampleWithV2RequestInterval() {
	client := hetrixtools.NewClient("token", hetrixtools.WithV2RequestInterval(0))
	fmt.Println(client != nil)
	// Output: true
}

func ExampleWithV3RequestInterval() {
	client := hetrixtools.NewClient("token", hetrixtools.WithV3RequestInterval(0))
	fmt.Println(client != nil)
	// Output: true
}

func ExampleWithV2BaseURL() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/custom-v2/token/blacklist-check/domain/example.com/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"SUCCESS","blacklisted_count":0}`))
	}))
	defer server.Close()

	client := hetrixtools.NewClientWithBaseURL("https://unused.invalid", "token",
		hetrixtools.WithV2BaseURL(server.URL+"/custom-v2"),
		hetrixtools.WithMinimumRequestInterval(0),
	)
	result, _ := client.CheckBlacklistDomain(context.Background(), "example.com")
	fmt.Println(result.BlacklistedCount)
	// Output: 0
}

func ExampleWithV3BaseURL() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/custom-v3/account/limits" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"api_calls_left":5}`))
	}))
	defer server.Close()

	client := hetrixtools.NewClientWithBaseURL("https://unused.invalid", "token",
		hetrixtools.WithV3BaseURL(server.URL+"/custom-v3"),
		hetrixtools.WithMinimumRequestInterval(0),
	)
	limits, _ := client.GetAccountLimits(context.Background())
	fmt.Println(limits.(map[string]any)["api_calls_left"])
	// Output: 5
}

func ExampleClient_GetAccountLimits() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v3/account/limits": `{"api_calls_left":42}`,
	})
	defer closeServer()

	limits, _ := client.GetAccountLimits(context.Background())
	fmt.Println(limits.(map[string]any)["api_calls_left"])
	// Output: 42
}

func ExampleClient_ListBlacklists() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v3/blacklists": `{"blacklists":["spamhaus"]}`,
	})
	defer closeServer()

	blacklists, _ := client.ListBlacklists(context.Background())
	fmt.Println(blacklists.(map[string]any)["blacklists"].([]any)[0])
	// Output: spamhaus
}

func ExampleClient_ListContactLists() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v3/contact-lists?page=1": `{"contact_lists":[{"id":"contacts-1"}]}`,
	})
	defer closeServer()

	contacts, _ := client.ListContactLists(context.Background(), map[string]string{"page": "1"})
	fmt.Println(contacts.(map[string]any)["contact_lists"].([]any)[0].(map[string]any)["id"])
	// Output: contacts-1
}

func ExampleClient_CreateBlacklistMonitor() {
	client, closeServer := exampleClient(map[string]string{
		"POST /v2/token/blacklist/add/": `{"status":"SUCCESS","monitor_id":"blacklist-1"}`,
	})
	defer closeServer()

	response, _ := client.CreateBlacklistMonitor(context.Background(), hetrixtools.BlacklistMonitorRequest{Target: "example.com"})
	fmt.Println(response.MonitorID)
	// Output: blacklist-1
}

func ExampleClient_UpdateBlacklistMonitor() {
	client, closeServer := exampleClient(map[string]string{
		"POST /v2/token/blacklist/edit/": `{"status":"SUCCESS","action":"updated"}`,
	})
	defer closeServer()

	response, _ := client.UpdateBlacklistMonitor(context.Background(), hetrixtools.BlacklistMonitorRequest{Target: "example.com", Label: "Example"})
	fmt.Println(response.Action)
	// Output: updated
}

func ExampleClient_UpsertBlacklistMonitor() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v3/blacklist-monitors?page=1&per_page=100": `{"blacklist_monitors":[],"meta":{"pagination":{"current":1,"last":1}}}`,
		"POST /v2/token/blacklist/add/":                  `{"status":"SUCCESS","monitor_id":"blacklist-1"}`,
	})
	defer closeServer()

	response, _ := client.UpsertBlacklistMonitor(context.Background(), hetrixtools.BlacklistMonitorRequest{Target: "example.com"})
	fmt.Println(response.MonitorID)
	// Output: blacklist-1
}

func ExampleClient_DeleteBlacklistMonitor() {
	client, closeServer := exampleClient(map[string]string{
		"POST /v2/token/blacklist/delete/": `{"status":"SUCCESS"}`,
	})
	defer closeServer()

	err := client.DeleteBlacklistMonitor(context.Background(), "example.com")
	fmt.Println(err == nil)
	// Output: true
}

func ExampleClient_ListBlacklistMonitors() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v3/blacklist-monitors?page=1": `{"blacklist_monitors":[{"id":"blacklist-1","target":"example.com"}],"meta":{"pagination":{"current":1,"last":1}}}`,
	})
	defer closeServer()

	monitors, _ := client.ListBlacklistMonitors(context.Background(), map[string]string{"page": "1"})
	fmt.Println(monitors.BlacklistMonitors[0].Target)
	// Output: example.com
}

func ExampleClient_GetBlacklistMonitor() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v3/blacklist-monitors?page=1&per_page=100": `{"blacklist_monitors":[{"id":"blacklist-1","target":"example.com"}],"meta":{"pagination":{"current":1,"last":1}}}`,
	})
	defer closeServer()

	monitor, _ := client.GetBlacklistMonitor(context.Background(), "example.com")
	fmt.Println(monitor.ID)
	// Output: blacklist-1
}

func ExampleClient_GetBlacklistMonitorReport() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v3/blacklist-monitors/blacklist-1/report": `{"target":"example.com"}`,
	})
	defer closeServer()

	report, _ := client.GetBlacklistMonitorReport(context.Background(), "blacklist-1", nil)
	fmt.Println(report.(map[string]any)["target"])
	// Output: example.com
}

func ExampleClient_CheckBlacklistIPv4() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v2/token/blacklist-check/ipv4/192.0.2.10/": `{"status":"SUCCESS","blacklisted_count":0}`,
	})
	defer closeServer()

	result, _ := client.CheckBlacklistIPv4(context.Background(), "192.0.2.10")
	fmt.Println(result.BlacklistedCount)
	// Output: 0
}

func ExampleClient_CheckBlacklistDomain() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v2/token/blacklist-check/domain/example.com/": `{"status":"SUCCESS","blacklisted_count":0}`,
	})
	defer closeServer()

	result, _ := client.CheckBlacklistDomain(context.Background(), "example.com")
	fmt.Println(result.BlacklistedCount)
	// Output: 0
}

func ExampleClient_CreateUptimeMonitor() {
	client, closeServer := exampleClient(map[string]string{
		"POST /v2/token/uptime/add/": `{"status":"SUCCESS","monitor_id":"monitor-1"}`,
	})
	defer closeServer()

	response, _ := client.CreateUptimeMonitor(context.Background(), hetrixtools.UptimeMonitorRequest{Type: "http", Name: "Homepage", Target: "https://example.com"})
	fmt.Println(response.MonitorID)
	// Output: monitor-1
}

func ExampleClient_UpdateUptimeMonitor() {
	client, closeServer := exampleClient(map[string]string{
		"POST /v2/token/uptime/add/": `{"status":"SUCCESS","action":"updated"}`,
	})
	defer closeServer()

	response, _ := client.UpdateUptimeMonitor(context.Background(), hetrixtools.UptimeMonitorRequest{MID: "monitor-1", Type: "http", Name: "Homepage", Target: "https://example.com"})
	fmt.Println(response.Action)
	// Output: updated
}

func ExampleClient_UpsertUptimeMonitor() {
	client, closeServer := exampleClient(map[string]string{
		"POST /v2/token/uptime/add/": `{"status":"SUCCESS","monitor_id":"monitor-1"}`,
	})
	defer closeServer()

	response, _ := client.UpsertUptimeMonitor(context.Background(), hetrixtools.UptimeMonitorRequest{Type: "http", Name: "Homepage", Target: "https://example.com"})
	fmt.Println(response.Status)
	// Output: SUCCESS
}

func ExampleClient_DeleteUptimeMonitor() {
	client, closeServer := exampleClient(map[string]string{
		"POST /v2/token/uptime/delete/": `{"status":"SUCCESS"}`,
	})
	defer closeServer()

	err := client.DeleteUptimeMonitor(context.Background(), "monitor-1")
	fmt.Println(err == nil)
	// Output: true
}

func ExampleClient_ListUptimeMonitors() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v3/uptime-monitors?page=1": `{"uptime_monitors":[{"id":"monitor-1","name":"Homepage","type":"website"}],"meta":{"pagination":{"current":1,"last":1}}}`,
	})
	defer closeServer()

	monitors, _ := client.ListUptimeMonitors(context.Background(), map[string]string{"page": "1"})
	fmt.Println(monitors.UptimeMonitors[0].Type)
	// Output: http
}

func ExampleClient_GetUptimeMonitor() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v3/uptime-monitors?page=1&per_page=100": `{"uptime_monitors":[{"id":"monitor-1","name":"Homepage","type":"website"}],"meta":{"pagination":{"current":1,"last":1}}}`,
	})
	defer closeServer()

	monitor, _ := client.GetUptimeMonitor(context.Background(), "monitor-1")
	fmt.Println(monitor.Name)
	// Output: Homepage
}

func ExampleClient_GetUptimeMonitorReport() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v3/uptime-monitors/monitor-1/report": `{"uptime":"99.99"}`,
	})
	defer closeServer()

	report, _ := client.GetUptimeMonitorReport(context.Background(), "monitor-1", nil)
	fmt.Println(report.(map[string]any)["uptime"])
	// Output: 99.99
}

func ExampleClient_ListUptimeMonitorDowntimes() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v3/uptime-monitors/monitor-1/downtimes": `{"downtimes":[{"id":"down-1"}]}`,
	})
	defer closeServer()

	downtimes, _ := client.ListUptimeMonitorDowntimes(context.Background(), "monitor-1", nil)
	fmt.Println(downtimes.(map[string]any)["downtimes"].([]any)[0].(map[string]any)["id"])
	// Output: down-1
}

func ExampleClient_GetUptimeMonitorLocationFailLog() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v3/uptime-monitors/monitor-1/location-fail-log": `{"locations":["new_york"]}`,
	})
	defer closeServer()

	log, _ := client.GetUptimeMonitorLocationFailLog(context.Background(), "monitor-1", nil)
	fmt.Println(log.(map[string]any)["locations"].([]any)[0])
	// Output: new_york
}

func ExampleClient_CreateScheduledMaintenance() {
	client, closeServer := exampleClient(map[string]string{
		"POST /v3/schedule-maintenance": `{"id":"maintenance-1","monitor_id":"monitor-1"}`,
	})
	defer closeServer()

	maintenance, _ := client.CreateScheduledMaintenance(context.Background(), hetrixtools.ScheduledMaintenanceRequest{MonitorID: "monitor-1"})
	fmt.Println(maintenance.ID)
	// Output: maintenance-1
}

func ExampleClient_DeleteScheduledMaintenance() {
	client, closeServer := exampleClient(map[string]string{
		"DELETE /v3/schedule-maintenance/maintenance-1": `{}`,
	})
	defer closeServer()

	err := client.DeleteScheduledMaintenance(context.Background(), "maintenance-1")
	fmt.Println(err == nil)
	// Output: true
}

func ExampleClient_ListScheduledMaintenances() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v3/schedule-maintenance?page=1": `{"scheduled_maintenances":[{"id":"maintenance-1"}],"meta":{"pagination":{"current":1,"last":1}}}`,
	})
	defer closeServer()

	maintenances, _ := client.ListScheduledMaintenances(context.Background(), map[string]string{"page": "1"})
	fmt.Println(maintenances.ScheduledMaintenances[0].ID)
	// Output: maintenance-1
}

func ExampleClient_GetScheduledMaintenance() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v3/schedule-maintenance?page=1&per_page=200": `{"scheduled_maintenances":[{"id":"maintenance-1"}],"meta":{"pagination":{"current":1,"last":1}}}`,
	})
	defer closeServer()

	maintenance, _ := client.GetScheduledMaintenance(context.Background(), "maintenance-1", "")
	fmt.Println(maintenance.ID)
	// Output: maintenance-1
}

func ExampleClient_ListStatusPages() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v3/status-pages?page=1": `{"status_pages":[{"id":"status-1","name":"Status"}],"meta":{"pagination":{"current":1,"last":1}}}`,
	})
	defer closeServer()

	pages, _ := client.ListStatusPages(context.Background(), map[string]string{"page": "1"})
	fmt.Println(pages.StatusPages[0].Name)
	// Output: Status
}

func ExampleClient_GetStatusPage() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v3/status-pages?page=1&per_page=100": `{"status_pages":[{"id":"status-1","name":"Status"}],"meta":{"pagination":{"current":1,"last":1}}}`,
	})
	defer closeServer()

	page, _ := client.GetStatusPage(context.Background(), "status-1")
	fmt.Println(page.Name)
	// Output: Status
}

func ExampleClient_AddStatusPageMonitors() {
	client, closeServer := exampleClient(map[string]string{
		"POST /v3/status-pages/status-1/monitors": `{}`,
	})
	defer closeServer()

	err := client.AddStatusPageMonitors(context.Background(), "status-1", []string{"monitor-1"})
	fmt.Println(err == nil)
	// Output: true
}

func ExampleClient_RemoveStatusPageMonitors() {
	client, closeServer := exampleClient(map[string]string{
		"DELETE /v3/status-pages/status-1/monitors": `{}`,
	})
	defer closeServer()

	err := client.RemoveStatusPageMonitors(context.Background(), "status-1", []string{"monitor-1"})
	fmt.Println(err == nil)
	// Output: true
}

func ExampleClient_AttachServerAgent() {
	client, closeServer := exampleClient(map[string]string{
		"POST /v3/uptime-monitors/monitor-1/server-agent": `{"agent_id":"agent-1"}`,
	})
	defer closeServer()

	agent, _ := client.AttachServerAgent(context.Background(), "monitor-1")
	fmt.Println(*agent.AgentID)
	// Output: agent-1
}

func ExampleClient_GetServerAgent() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v3/uptime-monitors/monitor-1/server-agent": `{"agent_id":"agent-1"}`,
	})
	defer closeServer()

	agent, _ := client.GetServerAgent(context.Background(), "monitor-1")
	fmt.Println(*agent.AgentID)
	// Output: agent-1
}

func ExampleClient_DetachServerAgent() {
	client, closeServer := exampleClient(map[string]string{
		"DELETE /v3/uptime-monitors/monitor-1/server-agent": `{}`,
	})
	defer closeServer()

	err := client.DetachServerAgent(context.Background(), "monitor-1")
	fmt.Println(err == nil)
	// Output: true
}

func ExampleClient_GetServerAgentWarningPolicies() {
	client, closeServer := exampleClient(map[string]string{
		"GET /v3/uptime-monitors/monitor-1/server-agent/warning-policies": `{"cpu":{"warning":90}}`,
	})
	defer closeServer()

	policies, _ := client.GetServerAgentWarningPolicies(context.Background(), "monitor-1")
	fmt.Println(policies.(map[string]any)["cpu"].(map[string]any)["warning"])
	// Output: 90
}

func ExampleClient_UpdateServerAgentWarningPolicies() {
	client, closeServer := exampleClient(map[string]string{
		"PUT /v3/uptime-monitors/monitor-1/server-agent/warning-policies": `{}`,
	})
	defer closeServer()

	err := client.UpdateServerAgentWarningPolicies(context.Background(), "monitor-1", map[string]any{"cpu": map[string]int{"warning": 90}})
	fmt.Println(err == nil)
	// Output: true
}

func ExampleUptimeMonitorRequest_MarshalJSON() {
	body, _ := json.Marshal(hetrixtools.UptimeMonitorRequest{
		Type:      "http",
		Name:      "Homepage",
		Target:    "https://example.com",
		Locations: []string{"new_york"},
	})
	fmt.Println(strings.Contains(string(body), `"Type":1`))
	// Output: true
}

func ExampleUptimeMonitorRequest_Validate() {
	err := hetrixtools.UptimeMonitorRequest{Type: "smtp", Target: "smtp.example.com"}.Validate()
	fmt.Println(err != nil)
	// Output: true
}

func ExampleError() {
	err := hetrixtools.Error{StatusCode: http.StatusForbidden, Body: `{"message":"forbidden"}`}

	var apiErr hetrixtools.Error
	if errors.As(err, &apiErr) {
		fmt.Println(apiErr.StatusCode)
	}
	// Output: 403
}

func ExampleIsNotFound() {
	fmt.Println(hetrixtools.IsNotFound(hetrixtools.Error{StatusCode: http.StatusNotFound}))
	// Output: true
}
