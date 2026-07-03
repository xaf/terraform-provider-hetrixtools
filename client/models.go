package hetrixtools

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
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
	MID              string   `json:"MID,omitempty"`
	Type             string   `json:"-"`
	Name             string   `json:"Name,omitempty"`
	Target           string   `json:"Target,omitempty"`
	Port             int64    `json:"Port,omitempty"`
	HTTPMethod       string   `json:"Method,omitempty"`
	MaxRedirects     int64    `json:"MaxRedirects,omitempty"`
	SMTPUser         string   `json:"SMTPUser,omitempty"`
	SMTPPass         string   `json:"SMTPPass,omitempty"`
	Timeout          int64    `json:"Timeout,omitempty"`
	Frequency        int64    `json:"Frequency,omitempty"`
	FailsBeforeAlert int64    `json:"FailsBeforeAlert,omitempty"`
	FailedLocations  int64    `json:"FailedLocations,omitempty"`
	ContactList      string   `json:"ContactList,omitempty"`
	Category         string   `json:"Category,omitempty"`
	AlertAfter       string   `json:"AlertAfter,omitempty"`
	RepeatTimes      int64    `json:"RepeatTimes,omitempty"`
	RepeatEvery      string   `json:"RepeatEvery,omitempty"`
	Public           *bool    `json:"Public,omitempty"`
	ShowTarget       *bool    `json:"ShowTarget,omitempty"`
	VerSSLCert       *bool    `json:"VerSSLCert,omitempty"`
	VerSSLHost       *bool    `json:"VerSSLHost,omitempty"`
	Locations        []string `json:"-"`
	Keyword          string   `json:"-"`
	HTTPCodes        []int64  `json:"-"`
	Grace            int64    `json:"Grace,omitempty"`
	INFOPub          *bool    `json:"INFOPub,omitempty"`
	CPUPub           *bool    `json:"CPUPub,omitempty"`
	RAMPub           *bool    `json:"RAMPub,omitempty"`
	DISKPub          *bool    `json:"DISKPub,omitempty"`
	NETPub           *bool    `json:"NETPub,omitempty"`
}

// MarshalJSON translates the canonical client model into the v2 add/update API shape.
func (r UptimeMonitorRequest) MarshalJSON() ([]byte, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	type uptimeMonitorRequest UptimeMonitorRequest
	body, err := json.Marshal(uptimeMonitorRequest(r))
	if err != nil {
		return nil, err
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	if len(r.Locations) > 0 {
		locations, _ := uptimeLocationCodes(r.Locations)
		payload["Locations"] = locations
	}
	if r.Type != "" {
		typeID, _ := uptimeMonitorTypeID(r.Type)
		payload["Type"] = typeID
	}
	if r.Keyword != "" {
		payload["Keyword"] = r.Keyword
	}
	if len(r.HTTPCodes) > 0 {
		payload["HTTPCodes"] = r.HTTPCodes
	}
	return json.Marshal(payload)
}

// Validate rejects combinations the HetrixTools uptime API cannot represent.
func (r UptimeMonitorRequest) Validate() error {
	monitorType := r.Type
	if monitorType == "" {
		return nil
	}
	if _, err := uptimeMonitorTypeID(monitorType); err != nil {
		return err
	}
	if _, err := uptimeLocationCodes(r.Locations); err != nil {
		return err
	}

	if monitorType != "http" {
		if r.HTTPMethod != "" {
			return fmt.Errorf("http_method is only supported for http uptime monitors")
		}
		if r.MaxRedirects != 0 {
			return fmt.Errorf("max_redirects is only supported for http uptime monitors")
		}
		if r.Keyword != "" {
			return fmt.Errorf("keyword is only supported for http uptime monitors")
		}
		if len(r.HTTPCodes) > 0 {
			return fmt.Errorf("accepted_http_codes is only supported for http uptime monitors")
		}
	}

	if monitorType != "smtp" {
		if r.Port != 0 {
			return fmt.Errorf("port is only supported for smtp uptime monitors")
		}
		if r.SMTPUser != "" {
			return fmt.Errorf("smtp_user is only supported for smtp uptime monitors")
		}
		if r.SMTPPass != "" {
			return fmt.Errorf("smtp_password is only supported for smtp uptime monitors")
		}
	}
	if monitorType == "smtp" && r.Port == 0 {
		return fmt.Errorf("port is required for smtp uptime monitors")
	}
	if (r.SMTPUser == "") != (r.SMTPPass == "") {
		return fmt.Errorf("smtp_user and smtp_password must be set together")
	}
	if monitorType == "http" || monitorType == "ping" || monitorType == "smtp" {
		if r.Target == "" {
			return fmt.Errorf("target is required for %s uptime monitors", monitorType)
		}
	}

	if monitorType != "heartbeat" {
		if r.Grace != 0 || r.INFOPub != nil || r.CPUPub != nil || r.RAMPub != nil || r.DISKPub != nil || r.NETPub != nil {
			return fmt.Errorf("heartbeat visibility and grace settings are only supported for heartbeat uptime monitors")
		}
	}
	if monitorType == "heartbeat" {
		if r.Target != "" {
			return fmt.Errorf("target is not supported for heartbeat uptime monitors")
		}
		if len(r.Locations) > 0 || r.FailedLocations != 0 {
			return fmt.Errorf("locations and failed_locations are not supported for heartbeat uptime monitors")
		}
	}

	if monitorType != "http" && monitorType != "smtp" {
		if r.VerSSLCert != nil || r.VerSSLHost != nil {
			return fmt.Errorf("SSL verification settings are only supported for http and smtp uptime monitors")
		}
	}
	return nil
}

// UptimeMonitor describes a HetrixTools uptime monitor.
type UptimeMonitor struct {
	ID               string   `json:"id"`
	Type             string   `json:"-"`
	Name             string   `json:"name"`
	Target           string   `json:"target"`
	Port             *int64   `json:"port"`
	HTTPMethod       string   `json:"http_method"`
	MaxRedirects     int64    `json:"max_redirects"`
	SMTPUser         string   `json:"smtp_user"`
	Timeout          int64    `json:"timeout"`
	Frequency        int64    `json:"frequency"`
	FailsBeforeAlert int64    `json:"fails_before_alert"`
	FailedLocations  int64    `json:"failed_locations"`
	ContactListID    string   `json:"contact_list_id"`
	Category         string   `json:"category"`
	AlertAfter       string   `json:"alert_after"`
	RepeatTimes      int64    `json:"repeat_times"`
	RepeatEvery      string   `json:"repeat_every"`
	Public           *bool    `json:"public"`
	ShowTarget       *bool    `json:"show_target"`
	VerSSLCert       *bool    `json:"verify_ssl_certificate"`
	VerSSLHost       *bool    `json:"verify_ssl_host"`
	Locations        []string `json:"-"`
	Keyword          string   `json:"keyword"`
	HTTPCodes        []int64  `json:"accepted_http_codes"`
	Grace            int64    `json:"grace"`
	InfoPublic       *bool    `json:"info_public"`
	CPUPublic        *bool    `json:"cpu_public"`
	RAMPublic        *bool    `json:"ram_public"`
	DiskPublic       *bool    `json:"disk_public"`
	NetPublic        *bool    `json:"net_public"`
	ServerID         *string  `json:"server_id"`
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
		TypeCamel             json.RawMessage `json:"Type"`
		AgentID               *string         `json:"agent_id"`
		PortCamel             *int64          `json:"Port"`
		HTTPMethodCamel       string          `json:"Method"`
		MaxRedirectsCamel     int64           `json:"MaxRedirects"`
		SMTPUserCamel         string          `json:"SMTPUser"`
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
	if m.Type == "" {
		m.Type = uptimeMonitorTypeName(firstNonEmptyRawMessage(aux.TypeRaw, aux.TypeCamel))
	}
	if m.Frequency == 0 {
		m.Frequency = aux.FrequencyV3
	}
	if m.ServerID == nil {
		m.ServerID = aux.AgentID
	}
	if m.Type != "heartbeat" {
		m.ServerID = nil
	}
	if m.Port == nil {
		m.Port = aux.PortCamel
	}
	if m.Type != "smtp" {
		m.Port = nil
	}
	if m.HTTPMethod == "" {
		m.HTTPMethod = aux.HTTPMethodCamel
	}
	if m.MaxRedirects == 0 {
		m.MaxRedirects = aux.MaxRedirectsCamel
	}
	if m.SMTPUser == "" {
		m.SMTPUser = aux.SMTPUserCamel
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
			m.Locations = make([]string, 0, len(aux.LocationsV3))
			for key := range aux.LocationsV3 {
				m.Locations = append(m.Locations, uptimeLocationName(key))
			}
			sort.Strings(m.Locations)
		} else {
			m.Locations = uptimeLocationNames(aux.LocationsCamel)
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
	return nil
}

func uptimeLocationNames(locations map[string]bool) []string {
	if locations == nil {
		return nil
	}
	normalized := make([]string, 0, len(locations))
	for key, value := range locations {
		if value {
			normalized = append(normalized, uptimeLocationName(key))
		}
	}
	sort.Strings(normalized)
	return normalized
}

func uptimeLocationCodes(locations []string) (map[string]bool, error) {
	if locations == nil {
		return nil, nil
	}
	normalized := map[string]bool{}
	for _, location := range locations {
		code, ok := uptimeLocationCode(location)
		if !ok {
			return nil, fmt.Errorf("unknown uptime monitor location %q", location)
		}
		normalized[code] = true
	}
	return normalized, nil
}

func uptimeLocationName(location string) string {
	switch location {
	case "nyc":
		return "new_york"
	case "sfo":
		return "san_francisco"
	case "dal":
		return "dallas"
	case "ams":
		return "amsterdam"
	case "lon":
		return "london"
	case "fra":
		return "frankfurt"
	case "sgp":
		return "singapore"
	case "syd":
		return "sydney"
	case "sao":
		return "sao_paulo"
	case "tok":
		return "tokyo"
	case "mba":
		return "mumbai"
	case "waw":
		return "warsaw"
	default:
		return location
	}
}

func uptimeLocationCode(location string) (string, bool) {
	switch location {
	case "new_york":
		return "nyc", true
	case "san_francisco":
		return "sfo", true
	case "dallas":
		return "dal", true
	case "amsterdam":
		return "ams", true
	case "london":
		return "lon", true
	case "frankfurt":
		return "fra", true
	case "singapore":
		return "sgp", true
	case "sydney":
		return "syd", true
	case "sao_paulo":
		return "sao", true
	case "tokyo":
		return "tok", true
	case "mumbai":
		return "mba", true
	case "warsaw":
		return "waw", true
	default:
		return "", false
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func firstNonEmptyRawMessage(values ...json.RawMessage) json.RawMessage {
	for _, value := range values {
		if len(value) > 0 {
			return value
		}
	}
	return nil
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

func uptimeMonitorTypeName(raw json.RawMessage) string {
	var id int64
	if err := json.Unmarshal(raw, &id); err == nil {
		switch id {
		case 1:
			return "http"
		case 2:
			return "ping"
		case 3:
			return "smtp"
		case 9:
			return "heartbeat"
		default:
			return ""
		}
	}
	var name string
	if err := json.Unmarshal(raw, &name); err != nil {
		return ""
	}
	switch name {
	case "website", "http":
		return "http"
	case "ping", "service":
		return "ping"
	case "smtp":
		return "smtp"
	case "server", "server_agent", "heartbeat":
		return "heartbeat"
	default:
		return name
	}
}

func uptimeMonitorTypeID(name string) (int64, error) {
	switch name {
	case "http":
		return 1, nil
	case "ping":
		return 2, nil
	case "smtp":
		return 3, nil
	case "heartbeat":
		return 9, nil
	default:
		return 0, fmt.Errorf("unknown uptime monitor type %q", name)
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
