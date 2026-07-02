package hetrixtools

import "net/url"

// ActionResponse is the common response returned by HetrixTools mutation endpoints.
type ActionResponse struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
	MonitorID    string `json:"monitor_id"`
	ServerID     string `json:"server_id"`
	Action       string `json:"action"`
}

// BlacklistCheckResult is the result of a one-off blacklist check.
type BlacklistCheckResult struct {
	Status                    string              `json:"status"`
	ErrorMessage              string              `json:"error_message"`
	APICallsLeft              int64               `json:"api_calls_left"`
	BlacklistCheckCreditsLeft int64               `json:"blacklist_check_credits_left"`
	BlacklistedCount          int64               `json:"blacklisted_count"`
	BlacklistedOn             []BlacklistListing  `json:"blacklisted_on"`
	Links                     BlacklistCheckLinks `json:"links"`
	RawJSON                   []byte              `json:"-"`
}

// BlacklistListing identifies one blacklist where a target is listed.
type BlacklistListing struct {
	RBL    string `json:"rbl"`
	Delist string `json:"delist"`
}

// BlacklistCheckLinks contains report URLs returned by a blacklist check.
type BlacklistCheckLinks struct {
	ReportLink            string `json:"report_link"`
	WhitelabelReportLink  string `json:"whitelabel_report_link"`
	APIReportLink         string `json:"api_report_link"`
	APIBlacklistCheckLink string `json:"api_blacklist_check_link"`
}

// BlacklistMonitorRequest describes a blacklist monitor create or update request.
type BlacklistMonitorRequest struct {
	Target  string
	Label   string
	Contact string
}

// Pagination describes HetrixTools paginated list metadata.
type Pagination struct {
	Current  int  `json:"current"`
	Last     int  `json:"last"`
	Previous *int `json:"previous"`
	Next     *int `json:"next"`
}

// Meta contains HetrixTools list response metadata.
type Meta struct {
	Total         int        `json:"total"`
	TotalFiltered int        `json:"total_filtered"`
	Returned      int        `json:"returned"`
	Pagination    Pagination `json:"pagination"`
}

// BlacklistMonitor describes a HetrixTools blacklist monitor.
type BlacklistMonitor struct {
	ID      string `json:"id"`
	Target  string `json:"target"`
	Label   string `json:"label"`
	Name    string `json:"name"`
	Contact string `json:"contact"`
}

// BlacklistMonitorsResponse is returned by ListBlacklistMonitors.
type BlacklistMonitorsResponse struct {
	BlacklistMonitors []BlacklistMonitor `json:"blacklist_monitors"`
	Meta              Meta               `json:"meta"`
}

func (r BlacklistMonitorRequest) form() url.Values {
	values := url.Values{"target": {r.Target}}
	if r.Label != "" {
		values.Set("label", r.Label)
	}
	if r.Contact != "" {
		values.Set("contact", r.Contact)
	}
	return values
}

// UptimeMonitorRequest describes an uptime monitor create or update request.
type UptimeMonitorRequest struct {
	MID              string          `json:"MID,omitempty"`
	Type             int64           `json:"Type,omitempty"`
	Name             string          `json:"Name,omitempty"`
	Target           string          `json:"Target,omitempty"`
	Port             int64           `json:"Port,omitempty"`
	Timeout          int64           `json:"Timeout,omitempty"`
	Frequency        int64           `json:"Frequency,omitempty"`
	FailsBeforeAlert int64           `json:"FailsBeforeAlert,omitempty"`
	FailedLocations  int64           `json:"FailedLocations,omitempty"`
	ContactList      string          `json:"ContactList,omitempty"`
	Category         string          `json:"Category,omitempty"`
	AlertAfter       string          `json:"AlertAfter,omitempty"`
	RepeatTimes      int64           `json:"RepeatTimes,omitempty"`
	RepeatEvery      string          `json:"RepeatEvery,omitempty"`
	Public           *bool           `json:"Public,omitempty"`
	ShowTarget       *bool           `json:"ShowTarget,omitempty"`
	VerSSLCert       *bool           `json:"VerSSLCert,omitempty"`
	VerSSLHost       *bool           `json:"VerSSLHost,omitempty"`
	Locations        map[string]bool `json:"Locations,omitempty"`
	Grace            int64           `json:"Grace,omitempty"`
	INFOPub          *bool           `json:"INFOPub,omitempty"`
	CPUPub           *bool           `json:"CPUPub,omitempty"`
	RAMPub           *bool           `json:"RAMPub,omitempty"`
	DISKPub          *bool           `json:"DISKPub,omitempty"`
	NETPub           *bool           `json:"NETPub,omitempty"`
	Extra            map[string]any  `json:"-"`
}

// UptimeMonitor describes a HetrixTools uptime monitor.
type UptimeMonitor struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Target   string `json:"target"`
	Category string `json:"category"`
}

// UptimeMonitorsResponse is returned by ListUptimeMonitors.
type UptimeMonitorsResponse struct {
	UptimeMonitors []UptimeMonitor `json:"uptime_monitors"`
	Meta           Meta            `json:"meta"`
}

// ScheduledMaintenanceRequest describes a scheduled maintenance create request.
type ScheduledMaintenanceRequest struct {
	MonitorID         string `json:"monitor_id"`
	Start             string `json:"start"`
	End               string `json:"end"`
	Timezone          string `json:"timezone"`
	WithNotifications bool   `json:"with_notifications"`
	RecurringTime     int64  `json:"recurring_time"`
	RecurringTimeType string `json:"recurring_time_type"`
}

// ScheduledMaintenance describes a HetrixTools scheduled maintenance window.
type ScheduledMaintenance struct {
	ID                string `json:"id"`
	MonitorID         string `json:"monitor_id"`
	Start             string `json:"start"`
	End               string `json:"end"`
	Timezone          string `json:"timezone"`
	WithNotifications bool   `json:"with_notifications"`
	Recurring         bool   `json:"recurring"`
	RecurringTime     int64  `json:"recurring_time"`
	RecurringTimeType string `json:"recurring_time_type"`
}

// ScheduledMaintenancesResponse is returned by ListScheduledMaintenances.
type ScheduledMaintenancesResponse struct {
	ScheduledMaintenances []ScheduledMaintenance `json:"scheduled_maintenances"`
	Meta                  Meta                   `json:"meta"`
}

// StatusPage describes a HetrixTools status page.
type StatusPage struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Monitors []string `json:"monitors"`
}

// StatusPagesResponse is returned by ListStatusPages.
type StatusPagesResponse struct {
	StatusPages []StatusPage `json:"status_pages"`
	Meta        Meta         `json:"meta"`
}

// ServerAgentResponse describes the server agent attached to an uptime monitor.
type ServerAgentResponse struct {
	AgentID *string `json:"agent_id"`
}

// MarshalJSON merges Extra fields into the uptime monitor request payload.
func (r UptimeMonitorRequest) MarshalJSON() ([]byte, error) {
	type alias UptimeMonitorRequest
	base, err := marshalWithoutExtra(alias(r))
	if err != nil {
		return nil, err
	}
	for key, value := range r.Extra {
		base[key] = value
	}
	return marshalMap(base)
}
