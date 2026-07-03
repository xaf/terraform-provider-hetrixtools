package hetrixtools

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestClientBlacklistMonitorActionsUseTokenPathAndFormBody(t *testing.T) {
	t.Parallel()

	var calls []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		if got := r.Header.Get("Authorization"); got != "" {
			t.Fatalf("Authorization header = %q, want empty", got)
		}
		if got, want := r.Header.Get("Content-Type"), "application/x-www-form-urlencoded"; got != want {
			t.Fatalf("Content-Type = %q, want %q", got, want)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form: %s", err)
		}
		if got, want := r.Form.Get("target"), "example.com"; got != want {
			t.Fatalf("target = %q, want %q", got, want)
		}

		switch r.URL.Path {
		case "/v2/test-token/blacklist/add/", "/v2/test-token/blacklist/edit/":
			if got, want := r.Form.Get("label"), "Example"; got != want {
				t.Fatalf("label = %q, want %q", got, want)
			}
			if got, want := r.Form.Get("contact"), "contacts-1"; got != want {
				t.Fatalf("contact = %q, want %q", got, want)
			}
			_, _ = w.Write([]byte(`{"status":"SUCCESS","action":"ok"}`))
		case "/v2/test-token/blacklist/delete/":
			_, _ = w.Write([]byte(`{"status":"SUCCESS"}`))
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	c := NewClientWithBaseURL(server.URL+"/v3", "test-token")
	request := BlacklistMonitorRequest{Target: "example.com", Label: "Example", Contact: "contacts-1"}
	if _, err := c.CreateBlacklistMonitor(context.Background(), request); err != nil {
		t.Fatalf("CreateBlacklistMonitor returned error: %s", err)
	}
	if _, err := c.UpdateBlacklistMonitor(context.Background(), request); err != nil {
		t.Fatalf("UpdateBlacklistMonitor returned error: %s", err)
	}
	if err := c.DeleteBlacklistMonitor(context.Background(), "example.com"); err != nil {
		t.Fatalf("DeleteBlacklistMonitor returned error: %s", err)
	}

	want := []string{
		"POST /v2/test-token/blacklist/add/",
		"POST /v2/test-token/blacklist/edit/",
		"POST /v2/test-token/blacklist/delete/",
	}
	assertStringSlicesEqual(t, calls, want)
}

func TestClientUptimeMonitorActionsUseTokenPathAndJSONBody(t *testing.T) {
	t.Parallel()

	var calls []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		if got := r.Header.Get("Authorization"); got != "" {
			t.Fatalf("Authorization header = %q, want empty", got)
		}
		if got, want := r.Header.Get("Content-Type"), "application/json"; got != want {
			t.Fatalf("Content-Type = %q, want %q", got, want)
		}

		switch r.URL.Path {
		case "/v2/test-token/uptime/add/":
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode body: %s", err)
			}
			if got, want := body["Name"], "Homepage"; got != want {
				t.Fatalf("Name = %q, want %q", got, want)
			}
			if _, ok := body["extra_option"]; !ok {
				t.Fatalf("extra_option missing from body %#v", body)
			}
			_, _ = w.Write([]byte(`{"status":"SUCCESS","monitor_id":"mid-1","server_id":"srv-1"}`))
		case "/v2/test-token/uptime/delete/":
			var body map[string]string
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode delete body: %s", err)
			}
			if got, want := body["MID"], "mid-1"; got != want {
				t.Fatalf("MID = %q, want %q", got, want)
			}
			_, _ = w.Write([]byte(`{"status":"SUCCESS"}`))
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	c := NewClientWithBaseURL(server.URL+"/v3", "test-token")
	request := UptimeMonitorRequest{MID: "mid-1", Type: 1, Name: "Homepage", Extra: map[string]any{"extra_option": true}}
	created, err := c.CreateUptimeMonitor(context.Background(), request)
	if err != nil {
		t.Fatalf("CreateUptimeMonitor returned error: %s", err)
	}
	if got, want := created.MonitorID, "mid-1"; got != want {
		t.Fatalf("created monitor ID = %q, want %q", got, want)
	}
	updated, err := c.UpdateUptimeMonitor(context.Background(), request)
	if err != nil {
		t.Fatalf("UpdateUptimeMonitor returned error: %s", err)
	}
	if got, want := updated.ServerID, "srv-1"; got != want {
		t.Fatalf("updated server ID = %q, want %q", got, want)
	}
	if err := c.DeleteUptimeMonitor(context.Background(), "mid-1"); err != nil {
		t.Fatalf("DeleteUptimeMonitor returned error: %s", err)
	}

	want := []string{
		"POST /v2/test-token/uptime/add/",
		"POST /v2/test-token/uptime/add/",
		"POST /v2/test-token/uptime/delete/",
	}
	assertStringSlicesEqual(t, calls, want)
}

func TestClientUpsertMethodsChooseCreateOrUpdate(t *testing.T) {
	t.Parallel()

	var calls []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		switch r.URL.Path {
		case "/v3/blacklist-monitors":
			_, _ = w.Write([]byte(`{"blacklist_monitors":[{"id":"bm-1","target":"existing.example"}],"meta":{"pagination":{"current":1,"last":1}}}`))
		case "/v2/test-token/blacklist/add/", "/v2/test-token/blacklist/edit/":
			_, _ = w.Write([]byte(`{"status":"SUCCESS"}`))
		case "/v2/test-token/uptime/add/":
			_, _ = w.Write([]byte(`{"status":"SUCCESS"}`))
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	c := NewClientWithBaseURL(server.URL+"/v3", "test-token")
	if _, err := c.UpsertBlacklistMonitor(context.Background(), BlacklistMonitorRequest{Target: "new.example"}); err != nil {
		t.Fatalf("UpsertBlacklistMonitor create returned error: %s", err)
	}
	if _, err := c.UpsertBlacklistMonitor(context.Background(), BlacklistMonitorRequest{Target: "existing.example"}); err != nil {
		t.Fatalf("UpsertBlacklistMonitor update returned error: %s", err)
	}
	if _, err := c.UpsertUptimeMonitor(context.Background(), UptimeMonitorRequest{Name: "New"}); err != nil {
		t.Fatalf("UpsertUptimeMonitor create returned error: %s", err)
	}
	if _, err := c.UpsertUptimeMonitor(context.Background(), UptimeMonitorRequest{MID: "up-1", Name: "Existing"}); err != nil {
		t.Fatalf("UpsertUptimeMonitor update returned error: %s", err)
	}

	want := []string{
		"GET /v3/blacklist-monitors",
		"POST /v2/test-token/blacklist/add/",
		"GET /v3/blacklist-monitors",
		"POST /v2/test-token/blacklist/edit/",
		"POST /v2/test-token/uptime/add/",
		"POST /v2/test-token/uptime/add/",
	}
	assertStringSlicesEqual(t, calls, want)
}

func TestClientReadMethodsUseBearerAuthAndPaginate(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.Header.Get("Authorization"), "Bearer test-token"; got != want {
			t.Fatalf("Authorization header = %q, want %q", got, want)
		}

		switch r.URL.Path {
		case "/v3/blacklist-monitors":
			writePage(w, "blacklist_monitors", r.URL.Query().Get("page"), `{"id":"bm-1","target":"first.example"}`, `{"id":"bm-2","target":"target.example","label":"Target"}`)
		case "/v3/uptime-monitors":
			writePage(w, "uptime_monitors", r.URL.Query().Get("page"), `{"id":"up-1","name":"First"}`, `{"id":"up-2","name":"Second","target":"https://example.com","category":"prod"}`)
		case "/v3/schedule-maintenance":
			writePage(w, "scheduled_maintenances", r.URL.Query().Get("page"), `{"id":"sm-1","monitor_id":"up-1"}`, `{"id":"sm-2","monitor_id":"up-2","start":"2026-07-02 10:00"}`)
		case "/v3/status-pages":
			writePage(w, "status_pages", r.URL.Query().Get("page"), `{"id":"sp-1","monitors":["up-1"]}`, `{"id":"sp-2","monitors":["up-2"]}`)
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	c := NewClientWithBaseURL(server.URL+"/v3", "test-token")
	blacklistMonitor, err := c.GetBlacklistMonitor(context.Background(), "target.example")
	if err != nil {
		t.Fatalf("GetBlacklistMonitor returned error: %s", err)
	}
	if got, want := blacklistMonitor.ID, "bm-2"; got != want {
		t.Fatalf("blacklist monitor ID = %q, want %q", got, want)
	}
	uptimeMonitor, err := c.GetUptimeMonitor(context.Background(), "up-2")
	if err != nil {
		t.Fatalf("GetUptimeMonitor returned error: %s", err)
	}
	if got, want := uptimeMonitor.Category, "prod"; got != want {
		t.Fatalf("uptime monitor category = %q, want %q", got, want)
	}
	maintenance, err := c.GetScheduledMaintenance(context.Background(), "sm-2", "up-2")
	if err != nil {
		t.Fatalf("GetScheduledMaintenance returned error: %s", err)
	}
	if got, want := maintenance.Start, "2026-07-02 10:00"; got != want {
		t.Fatalf("maintenance start = %q, want %q", got, want)
	}
	statusPage, err := c.GetStatusPage(context.Background(), "sp-2")
	if err != nil {
		t.Fatalf("GetStatusPage returned error: %s", err)
	}
	if got, want := statusPage.Monitors[0], "up-2"; got != want {
		t.Fatalf("status page monitor = %q, want %q", got, want)
	}
}

func TestClientWriteMethodsUseBearerAuthAndJSONBody(t *testing.T) {
	t.Parallel()

	var calls []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.Method+" "+r.URL.Path)
		if got, want := r.Header.Get("Authorization"), "Bearer test-token"; got != want {
			t.Fatalf("Authorization header = %q, want %q", got, want)
		}

		switch r.URL.Path {
		case "/v3/schedule-maintenance":
			if r.Method != http.MethodPost {
				t.Fatalf("method = %q, want POST", r.Method)
			}
			var body ScheduledMaintenanceRequest
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode maintenance body: %s", err)
			}
			if got, want := body.MonitorID, "up-1"; got != want {
				t.Fatalf("monitor ID = %q, want %q", got, want)
			}
			_, _ = w.Write([]byte(`{"id":"sm-1","monitor_id":"up-1"}`))
		case "/v3/schedule-maintenance/sm-1":
			if r.Method != http.MethodDelete {
				t.Fatalf("method = %q, want DELETE", r.Method)
			}
			_, _ = w.Write([]byte(`{}`))
		case "/v3/status-pages/sp-1/monitors":
			if r.Method != http.MethodPost && r.Method != http.MethodDelete {
				t.Fatalf("method = %q, want POST or DELETE", r.Method)
			}
			var body []string
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode monitor IDs: %s", err)
			}
			if len(body) != 1 || body[0] != "up-1" {
				t.Fatalf("monitor IDs = %#v, want [up-1]", body)
			}
			_, _ = w.Write([]byte(`{}`))
		case "/v3/uptime-monitors/up-1/server-agent":
			switch r.Method {
			case http.MethodPost:
				_, _ = w.Write([]byte(`{"agent_id":"agent-1"}`))
			case http.MethodGet:
				_, _ = w.Write([]byte(`{"agent_id":"agent-1"}`))
			case http.MethodDelete:
				_, _ = w.Write([]byte(`{}`))
			default:
				t.Fatalf("method = %q, want POST, GET, or DELETE", r.Method)
			}
		case "/v3/uptime-monitors/up-1/server-agent/warning-policies":
			if r.Method == http.MethodGet {
				_, _ = w.Write([]byte(`{"cpu":{"warning":90}}`))
				return
			}
			if r.Method != http.MethodPut {
				t.Fatalf("method = %q, want PUT", r.Method)
			}
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode policies: %s", err)
			}
			_, _ = w.Write([]byte(`{}`))
		case "/v3/custom":
			if got, want := r.URL.Query().Get("x"), "1"; got != want {
				t.Fatalf("x query = %q, want %q", got, want)
			}
			_, _ = w.Write([]byte(`{"ok":true}`))
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
		}
	}))
	defer server.Close()

	c := NewClientWithBaseURL(server.URL+"/v3", "test-token")
	if _, err := c.CreateScheduledMaintenance(context.Background(), ScheduledMaintenanceRequest{MonitorID: "up-1"}); err != nil {
		t.Fatalf("CreateScheduledMaintenance returned error: %s", err)
	}
	if err := c.DeleteScheduledMaintenance(context.Background(), "sm-1"); err != nil {
		t.Fatalf("DeleteScheduledMaintenance returned error: %s", err)
	}
	if err := c.AddStatusPageMonitors(context.Background(), "sp-1", []string{"up-1"}); err != nil {
		t.Fatalf("AddStatusPageMonitors returned error: %s", err)
	}
	if err := c.RemoveStatusPageMonitors(context.Background(), "sp-1", []string{"up-1"}); err != nil {
		t.Fatalf("RemoveStatusPageMonitors returned error: %s", err)
	}
	serverAgent, err := c.AttachServerAgent(context.Background(), "up-1")
	if err != nil {
		t.Fatalf("AttachServerAgent returned error: %s", err)
	}
	if got, want := *serverAgent.AgentID, "agent-1"; got != want {
		t.Fatalf("agent ID = %q, want %q", got, want)
	}
	if _, err := c.GetServerAgent(context.Background(), "up-1"); err != nil {
		t.Fatalf("GetServerAgent returned error: %s", err)
	}
	if err := c.DetachServerAgent(context.Background(), "up-1"); err != nil {
		t.Fatalf("DetachServerAgent returned error: %s", err)
	}
	if _, err := c.GetServerAgentWarningPolicies(context.Background(), "up-1"); err != nil {
		t.Fatalf("GetServerAgentWarningPolicies returned error: %s", err)
	}
	if err := c.UpdateServerAgentWarningPolicies(context.Background(), "up-1", map[string]any{"cpu": map[string]any{"warning": 90}}); err != nil {
		t.Fatalf("UpdateServerAgentWarningPolicies returned error: %s", err)
	}
	if _, err := c.getEndpoint(context.Background(), "/custom", map[string]string{"x": "1"}); err != nil {
		t.Fatalf("getEndpoint returned error: %s", err)
	}

	want := []string{
		"POST /v3/schedule-maintenance",
		"DELETE /v3/schedule-maintenance/sm-1",
		"POST /v3/status-pages/sp-1/monitors",
		"DELETE /v3/status-pages/sp-1/monitors",
		"POST /v3/uptime-monitors/up-1/server-agent",
		"GET /v3/uptime-monitors/up-1/server-agent",
		"DELETE /v3/uptime-monitors/up-1/server-agent",
		"GET /v3/uptime-monitors/up-1/server-agent/warning-policies",
		"PUT /v3/uptime-monitors/up-1/server-agent/warning-policies",
		"GET /v3/custom",
	}
	assertStringSlicesEqual(t, calls, want)
}

func TestClientBlacklistCheckParsesResultAndRawJSON(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.Path, "/v2/test-token/blacklist-check/domain/example.com/"; got != want {
			t.Fatalf("path = %q, want %q", got, want)
		}
		_, _ = w.Write([]byte(`{"status":"SUCCESS","blacklisted_count":1,"blacklisted_on":[{"rbl":"rbl.example","delist":"https://delist.example"}],"links":{"report_link":"https://report.example"}}`))
	}))
	defer server.Close()

	c := NewClientWithBaseURL(server.URL+"/v3", "test-token")
	result, err := c.CheckBlacklistDomain(context.Background(), "example.com")
	if err != nil {
		t.Fatalf("CheckBlacklistDomain returned error: %s", err)
	}
	if got, want := result.BlacklistedCount, int64(1); got != want {
		t.Fatalf("blacklisted count = %d, want %d", got, want)
	}
	if got, want := result.BlacklistedOn[0].RBL, "rbl.example"; got != want {
		t.Fatalf("rbl = %q, want %q", got, want)
	}
	if !strings.Contains(string(result.RawJSON), "blacklisted_count") {
		t.Fatalf("raw JSON missing original body: %s", string(result.RawJSON))
	}
}

func TestClientReturnsErrorOnActionErrorStatus(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"status":"ERROR","error_message":"monitor already exists"}`))
	}))
	defer server.Close()

	c := NewClientWithBaseURL(server.URL+"/v3", "test-token")
	_, err := c.CreateUptimeMonitor(context.Background(), UptimeMonitorRequest{Name: "Homepage"})
	if err == nil {
		t.Fatal("expected action error")
	}
	actionErr, ok := err.(Error)
	if !ok {
		t.Fatalf("error type = %T, want Error", err)
	}
	if actionErr.Response == nil {
		t.Fatal("expected action response on Error")
	}
	if got, want := actionErr.Response.ErrorMessage, "monitor already exists"; got != want {
		t.Fatalf("error message = %q, want %q", got, want)
	}
}

func TestClientReturnsErrorAndRejectsAbsoluteEndpointPaths(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"message":"nope"}`, http.StatusForbidden)
	}))
	defer server.Close()

	c := NewClientWithBaseURL(server.URL, "test-token")
	_, err := c.getEndpoint(context.Background(), "/anything", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(Error)
	if !ok {
		t.Fatalf("error type = %T, want Error", err)
	}
	if got, want := apiErr.StatusCode, http.StatusForbidden; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
	if _, err := c.getEndpoint(context.Background(), "https://api.example.test/anything", nil); err == nil {
		t.Fatal("expected absolute endpoint path error")
	}
}

func TestClientV3RateLimitHeadersThrottleOnlyMatchingScope(t *testing.T) {
	var calls []string
	reset := time.Now().Add(2 * time.Second).Unix()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls = append(calls, r.URL.Path)
		if r.URL.Path == "/v3/limited" && len(calls) == 1 {
			w.Header().Set("ratelimit-remaining-endpoint", "0")
			w.Header().Set("ratelimit-reset-endpoint", strconv.FormatInt(reset, 10))
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	c := NewClientWithBaseURL(server.URL, "test-token", WithV3RequestInterval(0))
	if _, err := c.getEndpoint(context.Background(), "/limited", nil); err != nil {
		t.Fatalf("first limited request returned error: %s", err)
	}

	otherStart := time.Now()
	if _, err := c.getEndpoint(context.Background(), "/other", nil); err != nil {
		t.Fatalf("other request returned error: %s", err)
	}
	if elapsed := time.Since(otherStart); elapsed > 500*time.Millisecond {
		t.Fatalf("unrelated endpoint waited %s, want no endpoint-limit delay", elapsed)
	}

	limitedStart := time.Now()
	if _, err := c.getEndpoint(context.Background(), "/limited", nil); err != nil {
		t.Fatalf("second limited request returned error: %s", err)
	}
	if elapsed := time.Since(limitedStart); elapsed < 500*time.Millisecond {
		t.Fatalf("same endpoint waited %s, want reset delay", elapsed)
	}

	want := []string{"/v3/limited", "/v3/other", "/v3/limited"}
	assertStringSlicesEqual(t, calls, want)
}

func TestClientRetriesTooManyRequestsAfterDelay(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		if calls == 1 {
			w.Header().Set("Retry-After", "1")
			http.Error(w, `{"status":"too_many_requests","message":"user api rate limit exceeded"}`, http.StatusTooManyRequests)
			return
		}
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	c := NewClientWithBaseURL(server.URL, "test-token", WithV3RequestInterval(0))
	start := time.Now()
	if _, err := c.getEndpoint(context.Background(), "/limited", nil); err != nil {
		t.Fatalf("request returned error: %s", err)
	}
	if elapsed := time.Since(start); elapsed < time.Second {
		t.Fatalf("retry waited %s, want at least Retry-After delay", elapsed)
	}
	if calls != 2 {
		t.Fatalf("calls = %d, want 2", calls)
	}
}

func TestClientDerivesVersionedBaseURLs(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct {
		name string
		base string
		v2   string
		v3   string
	}{
		{name: "root", base: "https://api.example.test", v2: "https://api.example.test/v2", v3: "https://api.example.test/v3"},
		{name: "root trailing slash", base: "https://api.example.test/", v2: "https://api.example.test/v2", v3: "https://api.example.test/v3"},
		{name: "v2 compatibility", base: "https://api.example.test/v2", v2: "https://api.example.test/v2", v3: "https://api.example.test/v3"},
		{name: "v3 compatibility", base: "https://api.example.test/v3", v2: "https://api.example.test/v2", v3: "https://api.example.test/v3"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := NewClientWithBaseURL(tt.base, "test-token")
			if c.v2BaseURL != tt.v2 {
				t.Fatalf("v2BaseURL = %q, want %q", c.v2BaseURL, tt.v2)
			}
			if c.v3BaseURL != tt.v3 {
				t.Fatalf("v3BaseURL = %q, want %q", c.v3BaseURL, tt.v3)
			}
		})
	}
}

func TestClientAllowsVersionedBaseURLOverrides(t *testing.T) {
	t.Parallel()

	c := NewClientWithBaseURL("https://api.example.test", "test-token", WithV2BaseURL("https://v2.example.test/"), WithV3BaseURL("https://v3.example.test/"))
	if c.v2BaseURL != "https://v2.example.test" {
		t.Fatalf("v2BaseURL = %q", c.v2BaseURL)
	}
	if c.v3BaseURL != "https://v3.example.test" {
		t.Fatalf("v3BaseURL = %q", c.v3BaseURL)
	}
}

func TestUptimeMonitorUnmarshalAcceptsLegacyCamelCaseFields(t *testing.T) {
	t.Parallel()

	var monitor UptimeMonitor
	body := []byte(`{
		"id":"up-1",
		"Type":1,
		"Name":"Homepage",
		"Target":"https://example.com",
		"Timeout":10,
		"Frequency":60,
		"FailsBeforeAlert":3,
		"FailedLocations":2,
		"ContactList":"contacts-1",
		"Category":"prod",
		"AlertAfter":"1m",
		"RepeatTimes":5,
		"RepeatEvery":"1h",
		"Public":true,
		"ShowTarget":false,
		"VerSSLCert":true,
		"VerSSLHost":true,
		"Locations":{"ams":true,"nyc":false},
		"Grace":120,
		"INFOPub":true,
		"CPUPub":false,
		"RAMPub":true,
		"DISKPub":false,
		"NETPub":true,
		"server_id":"srv-1"
	}`)
	if err := json.Unmarshal(body, &monitor); err != nil {
		t.Fatalf("unmarshal uptime monitor: %s", err)
	}

	if got, want := monitor.ContactListID, "contacts-1"; got != want {
		t.Fatalf("contact list = %q, want %q", got, want)
	}
	if got, want := monitor.FailsBeforeAlert, int64(3); got != want {
		t.Fatalf("fails before alert = %d, want %d", got, want)
	}
	if monitor.ShowTarget == nil || *monitor.ShowTarget {
		t.Fatalf("show target = %#v, want false", monitor.ShowTarget)
	}
	if monitor.VerSSLCert == nil || !*monitor.VerSSLCert {
		t.Fatalf("verify SSL certificate = %#v, want true", monitor.VerSSLCert)
	}
	if got, want := monitor.Locations["ams"], true; got != want {
		t.Fatalf("ams location = %v, want %v", got, want)
	}
	if monitor.InfoPublic == nil || !*monitor.InfoPublic {
		t.Fatalf("info public = %#v, want true", monitor.InfoPublic)
	}
}

func TestUptimeMonitorsResponseUnmarshalAcceptsV3MonitorsEnvelope(t *testing.T) {
	t.Parallel()

	var response UptimeMonitorsResponse
	body := []byte(`{
		"monitors":[{
			"id":"up-1",
			"name":"Homepage",
			"type":"website",
			"target":"https://example.com",
			"timeout":10,
			"check_frequency":1,
			"contact_lists":["contacts-1"],
			"locations":{"new_york":{"uptime_status":"up"}},
			"public_report":false,
			"public_target":false,
			"verify_ssl_certificate":true,
			"verify_ssl_hostname":true,
			"number_of_tries":3,
			"triggering_locations":2,
			"alert_after_minutes":5,
			"repeat_alert_times":1,
			"repeat_alert_frequency":60
		}]
	}`)
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("unmarshal uptime monitors response: %s", err)
	}

	if got, want := len(response.UptimeMonitors), 1; got != want {
		t.Fatalf("uptime monitors length = %d, want %d", got, want)
	}
	monitor := response.UptimeMonitors[0]
	if got, want := monitor.Type, int64(1); got != want {
		t.Fatalf("type = %d, want %d", got, want)
	}
	if got, want := monitor.Frequency, int64(1); got != want {
		t.Fatalf("frequency = %d, want %d", got, want)
	}
	if got, want := monitor.ContactListID, "contacts-1"; got != want {
		t.Fatalf("contact list = %q, want %q", got, want)
	}
	if got, want := monitor.AlertAfter, "5m"; got != want {
		t.Fatalf("alert after = %q, want %q", got, want)
	}
	if got, want := monitor.RepeatEvery, "60m"; got != want {
		t.Fatalf("repeat every = %q, want %q", got, want)
	}
	if !monitor.Locations["new_york"] {
		t.Fatalf("new_york location missing or false: %#v", monitor.Locations)
	}
}

func TestBlacklistMonitorsResponseUnmarshalAcceptsV3MonitorsEnvelope(t *testing.T) {
	t.Parallel()

	var response BlacklistMonitorsResponse
	body := []byte(`{
		"monitors":[{
			"id":"blacklist-1",
			"target":"example.com",
			"contact_lists":["contacts-1"]
		}]
	}`)
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("unmarshal blacklist monitors response: %s", err)
	}

	if got, want := len(response.BlacklistMonitors), 1; got != want {
		t.Fatalf("blacklist monitors length = %d, want %d", got, want)
	}
	monitor := response.BlacklistMonitors[0]
	if got, want := monitor.Target, "example.com"; got != want {
		t.Fatalf("target = %q, want %q", got, want)
	}
	if got, want := monitor.Contact, "contacts-1"; got != want {
		t.Fatalf("contact = %q, want %q", got, want)
	}
}

func writePage(w http.ResponseWriter, name string, page string, first string, second string) {
	w.Header().Set("Content-Type", "application/json")
	if page == "" || page == "1" {
		_, _ = w.Write([]byte(`{"` + name + `":[` + first + `],"meta":{"pagination":{"current":1,"last":2,"next":2}}}`))
		return
	}
	_, _ = w.Write([]byte(`{"` + name + `":[` + second + `],"meta":{"pagination":{"current":2,"last":2}}}`))
}

func assertStringSlicesEqual(t *testing.T, got []string, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("len(got) = %d, want %d; got %#v", len(got), len(want), got)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %q, want %q; got all %#v", i, got[i], want[i], got)
		}
	}
}
