package provider

type pagination struct {
	Current  int  `json:"current"`
	Last     int  `json:"last"`
	Previous *int `json:"previous"`
	Next     *int `json:"next"`
}

type meta struct {
	Total         int        `json:"total"`
	TotalFiltered int        `json:"total_filtered"`
	Returned      int        `json:"returned"`
	Pagination    pagination `json:"pagination"`
}

type scheduledMaintenancesResponse struct {
	ScheduledMaintenances []scheduledMaintenance `json:"scheduled_maintenances"`
	Meta                  meta                   `json:"meta"`
}

type scheduledMaintenance struct {
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

type statusPagesResponse struct {
	StatusPages []statusPage `json:"status_pages"`
	Meta        meta         `json:"meta"`
}

type statusPage struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Type     string   `json:"type"`
	Monitors []string `json:"monitors"`
}

type serverAgentResponse struct {
	AgentID *string `json:"agent_id"`
}
