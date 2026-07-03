package hetrixtools

import (
	"encoding/json"
	"net/url"
	"strconv"
)

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

// UnmarshalJSON accepts alternate field names returned by HetrixTools views.
func (m *BlacklistMonitor) UnmarshalJSON(body []byte) error {
	type blacklistMonitor BlacklistMonitor
	var aux struct {
		blacklistMonitor
		TargetIP     string   `json:"ip"`
		TargetDomain string   `json:"domain"`
		TargetHost   string   `json:"host"`
		Contacts     []string `json:"contact_lists"`
	}
	if err := json.Unmarshal(body, &aux); err != nil {
		return err
	}
	*m = BlacklistMonitor(aux.blacklistMonitor)
	if m.Target == "" {
		m.Target = firstNonEmpty(aux.TargetIP, aux.TargetDomain, aux.TargetHost)
	}
	if m.Contact == "" && len(aux.Contacts) > 0 {
		m.Contact = aux.Contacts[0]
	}
	return nil
}

// BlacklistMonitorsResponse is returned by ListBlacklistMonitors.
type BlacklistMonitorsResponse struct {
	BlacklistMonitors []BlacklistMonitor `json:"blacklist_monitors"`
	Meta              Meta               `json:"meta"`
}

// UnmarshalJSON accepts documented and legacy list envelope names.
func (r *BlacklistMonitorsResponse) UnmarshalJSON(body []byte) error {
	type blacklistMonitorsResponse BlacklistMonitorsResponse
	var aux struct {
		blacklistMonitorsResponse
		Monitors []BlacklistMonitor `json:"monitors"`
		Data     []BlacklistMonitor `json:"data"`
	}
	if err := json.Unmarshal(body, &aux); err != nil {
		return err
	}
	*r = BlacklistMonitorsResponse(aux.blacklistMonitorsResponse)
	if len(r.BlacklistMonitors) == 0 {
		r.BlacklistMonitors = firstNonEmptySlice(aux.Monitors, aux.Data)
	}
	return nil
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
	ID               string          `json:"id"`
	Type             int64           `json:"type"`
	Name             string          `json:"name"`
	Target           string          `json:"target"`
	Timeout          int64           `json:"timeout"`
	Frequency        int64           `json:"frequency"`
	FailsBeforeAlert int64           `json:"fails_before_alert"`
	FailedLocations  int64           `json:"failed_locations"`
	ContactListID    string          `json:"contact_list_id"`
	Category         string          `json:"category"`
	AlertAfter       string          `json:"alert_after"`
	RepeatTimes      int64           `json:"repeat_times"`
	RepeatEvery      string          `json:"repeat_every"`
	Public           *bool           `json:"public"`
	ShowTarget       *bool           `json:"show_target"`
	VerSSLCert       *bool           `json:"verify_ssl_certificate"`
	VerSSLHost       *bool           `json:"verify_ssl_host"`
	Locations        map[string]bool `json:"locations"`
	Grace            int64           `json:"grace"`
	InfoPublic       *bool           `json:"info_public"`
	CPUPublic        *bool           `json:"cpu_public"`
	RAMPublic        *bool           `json:"ram_public"`
	DiskPublic       *bool           `json:"disk_public"`
	NetPublic        *bool           `json:"net_public"`
	ServerID         string          `json:"server_id"`
	Extra            map[string]any  `json:"-"`
}

// UnmarshalJSON accepts both v3 snake_case fields and legacy v2 camel-case names.
func (m *UptimeMonitor) UnmarshalJSON(body []byte) error {
	type uptimeMonitor UptimeMonitor
	var aux struct {
		uptimeMonitor
		IDMonitorID           string          `json:"monitor_id"`
		IDMID                 string          `json:"MID"`
		IDCamel               string          `json:"MonitorID"`
		TypeRaw               json.RawMessage `json:"type"`
		FrequencyV3           int64           `json:"check_frequency"`
		ContactLists          []string        `json:"contact_lists"`
		FailsBeforeAlertV3    int64           `json:"number_of_tries"`
		FailedLocationsV3     int64           `json:"triggering_locations"`
		AlertAfterMinutes     int64           `json:"alert_after_minutes"`
		RepeatTimesV3         int64           `json:"repeat_alert_times"`
		RepeatEveryV3         int64           `json:"repeat_alert_frequency"`
		PublicV3              *bool           `json:"public_report"`
		ShowTargetV3          *bool           `json:"public_target"`
		VerSSLHostV3          *bool           `json:"verify_ssl_hostname"`
		LocationsV3           map[string]any  `json:"locations"`
		KeywordV3             string          `json:"keyword"`
		HTTPCodesV3           []int64         `json:"accepted_http_codes"`
		ContactListCamel      string          `json:"ContactList"`
		FailsBeforeAlertCamel int64           `json:"FailsBeforeAlert"`
		FailedLocationsCamel  int64           `json:"FailedLocations"`
		AlertAfterCamel       string          `json:"AlertAfter"`
		RepeatTimesCamel      int64           `json:"RepeatTimes"`
		RepeatEveryCamel      string          `json:"RepeatEvery"`
		ShowTargetCamel       *bool           `json:"ShowTarget"`
		VerSSLCertCamel       *bool           `json:"VerSSLCert"`
		VerSSLHostCamel       *bool           `json:"VerSSLHost"`
		LocationsCamel        map[string]bool `json:"Locations"`
		InfoPublicCamel       *bool           `json:"INFOPub"`
		CPUPublicCamel        *bool           `json:"CPUPub"`
		RAMPublicCamel        *bool           `json:"RAMPub"`
		DiskPublicCamel       *bool           `json:"DISKPub"`
		NetPublicCamel        *bool           `json:"NETPub"`
	}
	if err := json.Unmarshal(body, &aux); err != nil {
		return err
	}
	*m = UptimeMonitor(aux.uptimeMonitor)
	if m.ID == "" {
		m.ID = firstNonEmpty(aux.IDMonitorID, aux.IDMID, aux.IDCamel)
	}
	if m.Type == 0 {
		m.Type = monitorTypeID(aux.TypeRaw)
	}
	if m.Frequency == 0 {
		m.Frequency = aux.FrequencyV3
	}
	if m.ContactListID == "" {
		if len(aux.ContactLists) > 0 {
			m.ContactListID = aux.ContactLists[0]
		} else {
			m.ContactListID = aux.ContactListCamel
		}
	}
	if m.FailsBeforeAlert == 0 {
		m.FailsBeforeAlert = firstNonZeroInt64(aux.FailsBeforeAlertV3, aux.FailsBeforeAlertCamel)
	}
	if m.FailedLocations == 0 {
		m.FailedLocations = firstNonZeroInt64(aux.FailedLocationsV3, aux.FailedLocationsCamel)
	}
	if m.AlertAfter == "" {
		m.AlertAfter = firstNonEmpty(durationMinutes(aux.AlertAfterMinutes), aux.AlertAfterCamel)
	}
	if m.RepeatTimes == 0 {
		m.RepeatTimes = firstNonZeroInt64(aux.RepeatTimesV3, aux.RepeatTimesCamel)
	}
	if m.RepeatEvery == "" {
		m.RepeatEvery = firstNonEmpty(durationMinutes(aux.RepeatEveryV3), aux.RepeatEveryCamel)
	}
	if m.Public == nil {
		m.Public = aux.PublicV3
	}
	if m.ShowTarget == nil {
		m.ShowTarget = aux.ShowTargetV3
	}
	if m.ShowTarget == nil {
		m.ShowTarget = aux.ShowTargetCamel
	}
	if m.VerSSLCert == nil {
		m.VerSSLCert = aux.VerSSLCertCamel
	}
	if m.VerSSLHost == nil {
		m.VerSSLHost = firstNonNilBool(aux.VerSSLHostV3, aux.VerSSLHostCamel)
	}
	if m.Locations == nil {
		if len(aux.LocationsV3) > 0 {
			m.Locations = map[string]bool{}
			for key := range aux.LocationsV3 {
				m.Locations[key] = true
			}
		} else {
			m.Locations = aux.LocationsCamel
		}
	}
	if m.InfoPublic == nil {
		m.InfoPublic = aux.InfoPublicCamel
	}
	if m.CPUPublic == nil {
		m.CPUPublic = aux.CPUPublicCamel
	}
	if m.RAMPublic == nil {
		m.RAMPublic = aux.RAMPublicCamel
	}
	if m.DiskPublic == nil {
		m.DiskPublic = aux.DiskPublicCamel
	}
	if m.NetPublic == nil {
		m.NetPublic = aux.NetPublicCamel
	}
	if m.Extra == nil {
		m.Extra = map[string]any{}
	}
	if aux.KeywordV3 != "" {
		m.Extra["Keyword"] = aux.KeywordV3
	}
	if len(aux.HTTPCodesV3) > 0 {
		m.Extra["HTTPCodes"] = aux.HTTPCodesV3
	}
	if len(m.Extra) == 0 {
		m.Extra = nil
	}
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func firstNonEmptySlice[T any](values ...[]T) []T {
	for _, value := range values {
		if len(value) > 0 {
			return value
		}
	}
	return nil
}

func firstNonZeroInt64(values ...int64) int64 {
	for _, value := range values {
		if value != 0 {
			return value
		}
	}
	return 0
}

func firstNonNilBool(values ...*bool) *bool {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}

func durationMinutes(minutes int64) string {
	if minutes == 0 {
		return ""
	}
	return strconv.FormatInt(minutes, 10) + "m"
}

func monitorTypeID(raw json.RawMessage) int64 {
	var id int64
	if err := json.Unmarshal(raw, &id); err == nil {
		return id
	}
	var name string
	if err := json.Unmarshal(raw, &name); err != nil {
		return 0
	}
	switch name {
	case "website":
		return 1
	case "ping", "service":
		return 2
	case "smtp":
		return 3
	case "server", "server_agent":
		return 9
	default:
		return 0
	}
}

// UptimeMonitorsResponse is returned by ListUptimeMonitors.
type UptimeMonitorsResponse struct {
	UptimeMonitors []UptimeMonitor `json:"uptime_monitors"`
	Meta           Meta            `json:"meta"`
}

// UnmarshalJSON accepts documented and legacy list envelope names.
func (r *UptimeMonitorsResponse) UnmarshalJSON(body []byte) error {
	type uptimeMonitorsResponse UptimeMonitorsResponse
	var aux struct {
		uptimeMonitorsResponse
		Monitors []UptimeMonitor `json:"monitors"`
		Data     []UptimeMonitor `json:"data"`
	}
	if err := json.Unmarshal(body, &aux); err != nil {
		return err
	}
	*r = UptimeMonitorsResponse(aux.uptimeMonitorsResponse)
	if len(r.UptimeMonitors) == 0 {
		r.UptimeMonitors = firstNonEmptySlice(aux.Monitors, aux.Data)
	}
	return nil
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
