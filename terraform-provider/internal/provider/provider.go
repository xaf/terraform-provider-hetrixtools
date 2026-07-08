package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	hetrixtools "github.com/xaf/terraform-provider-hetrixtools/client"
)

var _ provider.Provider = (*hetrixToolsProvider)(nil)

type hetrixToolsProvider struct {
	version string
}

type providerModel struct {
	APIToken  types.String `tfsdk:"api_token"`
	BaseURL   types.String `tfsdk:"base_url"`
	BaseURLV2 types.String `tfsdk:"base_url_v2"`
	BaseURLV3 types.String `tfsdk:"base_url_v3"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &hetrixToolsProvider{version: version}
	}
}

func (p *hetrixToolsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hetrixtools"
	resp.Version = p.version
}

func (p *hetrixToolsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Terraform provider for HetrixTools. This project is not affiliated with, endorsed by, or supported by HetrixTools. It is provided without any guarantee that the HetrixTools API or this provider will behave as expected for your account.",
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				MarkdownDescription: "HetrixTools API bearer token. Can also be set with `HETRIXTOOLS_API_TOKEN`.",
				Optional:            true,
				Sensitive:           true,
			},
			"base_url": schema.StringAttribute{
				MarkdownDescription: "HetrixTools API root URL. The provider appends `/v2` and `/v3`. Can also be set with `HETRIXTOOLS_BASE_URL`.",
				Optional:            true,
			},
			"base_url_v2": schema.StringAttribute{
				MarkdownDescription: "HetrixTools v2 API base URL. Overrides `base_url` for token-path endpoints. Can also be set with `HETRIXTOOLS_BASE_URL_V2`.",
				Optional:            true,
			},
			"base_url_v3": schema.StringAttribute{
				MarkdownDescription: "HetrixTools v3 API base URL. Overrides `base_url` for REST endpoints. Can also be set with `HETRIXTOOLS_BASE_URL_V3`.",
				Optional:            true,
			},
		},
	}
}

func (p *hetrixToolsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	token := os.Getenv("HETRIXTOOLS_API_TOKEN")
	if !config.APIToken.IsNull() && !config.APIToken.IsUnknown() {
		token = config.APIToken.ValueString()
	}
	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing HetrixTools API token",
			"Set api_token in the provider configuration or HETRIXTOOLS_API_TOKEN in the environment.",
		)
		return
	}

	baseURL := os.Getenv("HETRIXTOOLS_BASE_URL")
	if !config.BaseURL.IsNull() && !config.BaseURL.IsUnknown() {
		baseURL = config.BaseURL.ValueString()
	}
	if baseURL == "" {
		baseURL = hetrixtools.DefaultBaseURL
	}
	v2BaseURL := os.Getenv("HETRIXTOOLS_BASE_URL_V2")
	if !config.BaseURLV2.IsNull() && !config.BaseURLV2.IsUnknown() {
		v2BaseURL = config.BaseURLV2.ValueString()
	}
	v3BaseURL := os.Getenv("HETRIXTOOLS_BASE_URL_V3")
	if !config.BaseURLV3.IsNull() && !config.BaseURLV3.IsUnknown() {
		v3BaseURL = config.BaseURLV3.ValueString()
	}

	c := hetrixtools.NewClientWithBaseURL(baseURL, token, hetrixtools.WithV2BaseURL(v2BaseURL), hetrixtools.WithV3BaseURL(v3BaseURL))
	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *hetrixToolsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource {
			return newGetDataSource("account_limits", "Read HetrixTools account limits.", nil, func(ctx context.Context, c *hetrixtools.Client, _ map[string]string) (any, error) {
				return c.GetAccountLimits(ctx)
			})
		},
		func() datasource.DataSource {
			return newGetDataSource("blacklists", "List HetrixTools blacklist providers.", nil, func(ctx context.Context, c *hetrixtools.Client, _ map[string]string) (any, error) {
				return c.ListBlacklists(ctx)
			})
		},
		func() datasource.DataSource {
			return newGetDataSource("contact_lists", "List HetrixTools contact lists.", []queryAttribute{{Name: "page", Description: "Page number to return."}, {Name: "per_page", Description: "Number of contact lists per page."}}, func(ctx context.Context, c *hetrixtools.Client, query map[string]string) (any, error) {
				pagination, err := paginationRequest(query)
				if err != nil {
					return nil, err
				}
				return c.ListContactLists(ctx, hetrixtools.ListContactListsRequest{PaginationRequest: pagination})
			})
		},
		func() datasource.DataSource {
			return newGetDataSource("blacklist_monitors", "List HetrixTools blacklist monitors.", blacklistMonitorQueryAttributes(), func(ctx context.Context, c *hetrixtools.Client, query map[string]string) (any, error) {
				typed, err := blacklistMonitorRequest(query)
				if err != nil {
					return nil, err
				}
				return c.ListBlacklistMonitors(ctx, typed)
			})
		},
		func() datasource.DataSource {
			return newGetDataSource("blacklist_report", "Get a HetrixTools blacklist monitor report.", []queryAttribute{{Name: "identifier", Description: "Blacklist monitor ID, IP address, or hostname.", Required: true}, {Name: "date", Description: "Report date in YYYY-MM-DD format."}}, func(ctx context.Context, c *hetrixtools.Client, query map[string]string) (any, error) {
				identifier := takeQueryValue(query, "identifier")
				return c.GetBlacklistMonitorReport(ctx, identifier, hetrixtools.GetBlacklistMonitorReportRequest{Date: query["date"]})
			})
		},
		func() datasource.DataSource {
			return newGetDataSource("uptime_monitors", "List HetrixTools uptime monitors.", uptimeMonitorQueryAttributes(), func(ctx context.Context, c *hetrixtools.Client, query map[string]string) (any, error) {
				typed, err := uptimeMonitorRequest(query)
				if err != nil {
					return nil, err
				}
				return c.ListUptimeMonitors(ctx, typed)
			})
		},
		func() datasource.DataSource {
			return newGetDataSource("uptime_report", "Get a HetrixTools uptime monitor report.", []queryAttribute{{Name: "monitor_id", Description: "Uptime monitor ID.", Required: true}, {Name: "days", Description: "Number of recent days to display."}, {Name: "month", Description: "Month in YYYY-MM format."}, {Name: "timezone", Description: "Report timezone."}, {Name: "hourly_stats", Description: "Whether to include hourly stats."}}, func(ctx context.Context, c *hetrixtools.Client, query map[string]string) (any, error) {
				monitorID := takeQueryValue(query, "monitor_id")
				typed, err := uptimeReportRequest(query)
				if err != nil {
					return nil, err
				}
				return c.GetUptimeMonitorReport(ctx, monitorID, typed)
			})
		},
		func() datasource.DataSource {
			return newGetDataSource("uptime_downtimes", "List HetrixTools uptime monitor downtimes.", []queryAttribute{{Name: "monitor_id", Description: "Uptime monitor ID.", Required: true}, {Name: "page", Description: "Page number."}, {Name: "per_page", Description: "Results per page."}, {Name: "start_before", Description: "Only downtimes started before this timestamp."}, {Name: "start_after", Description: "Only downtimes started after this timestamp."}}, func(ctx context.Context, c *hetrixtools.Client, query map[string]string) (any, error) {
				monitorID := takeQueryValue(query, "monitor_id")
				typed, err := uptimeDowntimesRequest(query)
				if err != nil {
					return nil, err
				}
				return c.ListUptimeMonitorDowntimes(ctx, monitorID, typed)
			})
		},
		func() datasource.DataSource {
			return newGetDataSource("uptime_location_fail_log", "Get a HetrixTools uptime monitor location fail log.", []queryAttribute{{Name: "monitor_id", Description: "Uptime monitor ID.", Required: true}, {Name: "timestamp", Description: "Timestamp to start scanning from."}, {Name: "minutes", Description: "Number of minutes with log entries to return."}}, func(ctx context.Context, c *hetrixtools.Client, query map[string]string) (any, error) {
				monitorID := takeQueryValue(query, "monitor_id")
				typed, err := uptimeLocationFailLogRequest(query)
				if err != nil {
					return nil, err
				}
				return c.GetUptimeMonitorLocationFailLog(ctx, monitorID, typed)
			})
		},
		func() datasource.DataSource {
			return newGetDataSource("uptime_server_agent", "Get the HetrixTools server agent attached to an uptime monitor.", []queryAttribute{{Name: "monitor_id", Description: "Uptime monitor ID.", Required: true}}, func(ctx context.Context, c *hetrixtools.Client, query map[string]string) (any, error) {
				return c.GetServerAgent(ctx, query["monitor_id"])
			})
		},
		func() datasource.DataSource {
			return newGetDataSource("uptime_server_agent_warning_policies", "Get HetrixTools server agent warning policies for an uptime monitor.", []queryAttribute{{Name: "monitor_id", Description: "Uptime monitor ID.", Required: true}}, func(ctx context.Context, c *hetrixtools.Client, query map[string]string) (any, error) {
				return c.GetServerAgentWarningPolicies(ctx, query["monitor_id"])
			})
		},
		func() datasource.DataSource {
			return newGetDataSource("status_pages", "List HetrixTools status pages.", []queryAttribute{{Name: "page", Description: "Page number."}, {Name: "per_page", Description: "Results per page."}, {Name: "name", Description: "Status page name filter."}}, func(ctx context.Context, c *hetrixtools.Client, query map[string]string) (any, error) {
				typed, err := statusPagesRequest(query)
				if err != nil {
					return nil, err
				}
				return c.ListStatusPages(ctx, typed)
			})
		},
		func() datasource.DataSource {
			return newGetDataSource("scheduled_maintenances", "List HetrixTools scheduled maintenances.", []queryAttribute{{Name: "page", Description: "Page number."}, {Name: "per_page", Description: "Results per page."}, {Name: "monitor_id", Description: "Optional monitor ID filter."}}, func(ctx context.Context, c *hetrixtools.Client, query map[string]string) (any, error) {
				typed, err := scheduledMaintenancesRequest(query)
				if err != nil {
					return nil, err
				}
				return c.ListScheduledMaintenances(ctx, typed)
			})
		},
	}
}

func takeQueryValue(query map[string]string, key string) string {
	value := query[key]
	delete(query, key)
	return value
}

func paginationRequest(query map[string]string) (hetrixtools.PaginationRequest, error) {
	page, err := optionalIntQuery(query, "page")
	if err != nil {
		return hetrixtools.PaginationRequest{}, err
	}
	perPage, err := optionalIntQuery(query, "per_page")
	if err != nil {
		return hetrixtools.PaginationRequest{}, err
	}
	return hetrixtools.PaginationRequest{Page: page, PerPage: perPage}, nil
}

func blacklistMonitorRequest(query map[string]string) (hetrixtools.ListBlacklistMonitorsRequest, error) {
	pagination, err := paginationRequest(query)
	if err != nil {
		return hetrixtools.ListBlacklistMonitorsRequest{}, err
	}
	exactName, err := optionalBoolQuery(query, "exact_name")
	if err != nil {
		return hetrixtools.ListBlacklistMonitorsRequest{}, err
	}
	exactTarget, err := optionalBoolQuery(query, "exact_target")
	if err != nil {
		return hetrixtools.ListBlacklistMonitorsRequest{}, err
	}
	cidr, err := optionalIntQuery(query, "cidr")
	if err != nil {
		return hetrixtools.ListBlacklistMonitorsRequest{}, err
	}
	listed, err := optionalBoolQuery(query, "listed")
	if err != nil {
		return hetrixtools.ListBlacklistMonitorsRequest{}, err
	}
	return hetrixtools.ListBlacklistMonitorsRequest{
		PaginationRequest: pagination,
		Name:              query["name"],
		ExactName:         exactName,
		Target:            query["target"],
		ExactTarget:       exactTarget,
		CIDR:              cidr,
		Type:              query["type"],
		Listed:            listed,
		Order:             query["order"],
		OrderBy:           query["order_by"],
	}, nil
}

func uptimeMonitorRequest(query map[string]string) (hetrixtools.ListUptimeMonitorsRequest, error) {
	pagination, err := paginationRequest(query)
	if err != nil {
		return hetrixtools.ListUptimeMonitorsRequest{}, err
	}
	return hetrixtools.ListUptimeMonitorsRequest{
		PaginationRequest: pagination,
		ID:                query["id"],
		Name:              query["name"],
		Target:            query["target"],
		Category:          query["category"],
		Type:              query["type"],
		UptimeStatus:      query["uptime_status"],
		MonitorStatus:     query["monitor_status"],
		Order:             query["order"],
		OrderBy:           query["order_by"],
	}, nil
}

func uptimeReportRequest(query map[string]string) (hetrixtools.GetUptimeMonitorReportRequest, error) {
	days, err := optionalIntQuery(query, "days")
	if err != nil {
		return hetrixtools.GetUptimeMonitorReportRequest{}, err
	}
	hourlyStats, err := optionalBoolQuery(query, "hourly_stats")
	if err != nil {
		return hetrixtools.GetUptimeMonitorReportRequest{}, err
	}
	return hetrixtools.GetUptimeMonitorReportRequest{Days: days, Month: query["month"], Timezone: query["timezone"], HourlyStats: hourlyStats}, nil
}

func uptimeDowntimesRequest(query map[string]string) (hetrixtools.ListUptimeMonitorDowntimesRequest, error) {
	pagination, err := paginationRequest(query)
	if err != nil {
		return hetrixtools.ListUptimeMonitorDowntimesRequest{}, err
	}
	startBefore, err := optionalInt64Query(query, "start_before")
	if err != nil {
		return hetrixtools.ListUptimeMonitorDowntimesRequest{}, err
	}
	startAfter, err := optionalInt64Query(query, "start_after")
	if err != nil {
		return hetrixtools.ListUptimeMonitorDowntimesRequest{}, err
	}
	return hetrixtools.ListUptimeMonitorDowntimesRequest{PaginationRequest: pagination, StartBefore: startBefore, StartAfter: startAfter}, nil
}

func uptimeLocationFailLogRequest(query map[string]string) (hetrixtools.GetUptimeMonitorLocationFailLogRequest, error) {
	timestamp, err := optionalInt64Query(query, "timestamp")
	if err != nil {
		return hetrixtools.GetUptimeMonitorLocationFailLogRequest{}, err
	}
	minutes, err := optionalIntQuery(query, "minutes")
	if err != nil {
		return hetrixtools.GetUptimeMonitorLocationFailLogRequest{}, err
	}
	return hetrixtools.GetUptimeMonitorLocationFailLogRequest{Timestamp: timestamp, Minutes: minutes}, nil
}

func statusPagesRequest(query map[string]string) (hetrixtools.ListStatusPagesRequest, error) {
	pagination, err := paginationRequest(query)
	if err != nil {
		return hetrixtools.ListStatusPagesRequest{}, err
	}
	return hetrixtools.ListStatusPagesRequest{PaginationRequest: pagination, Name: query["name"]}, nil
}

func scheduledMaintenancesRequest(query map[string]string) (hetrixtools.ListScheduledMaintenancesRequest, error) {
	pagination, err := paginationRequest(query)
	if err != nil {
		return hetrixtools.ListScheduledMaintenancesRequest{}, err
	}
	return hetrixtools.ListScheduledMaintenancesRequest{PaginationRequest: pagination, MonitorID: query["monitor_id"]}, nil
}

func optionalIntQuery(query map[string]string, key string) (int, error) {
	value := query[key]
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", key, err)
	}
	return parsed, nil
}

func optionalInt64Query(query map[string]string, key string) (int64, error) {
	value := query[key]
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", key, err)
	}
	return parsed, nil
}

func optionalBoolQuery(query map[string]string, key string) (*bool, error) {
	value := query[key]
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf("%s must be a boolean: %w", key, err)
	}
	return &parsed, nil
}

func (p *hetrixToolsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newBlacklistMonitorResource,
		newUptimeHTTPMonitorResource,
		newUptimePingMonitorResource,
		newUptimeSMTPMonitorResource,
		newUptimeHeartbeatMonitorResource,
		newScheduledMaintenanceResource,
		newStatusPageMonitorsResource,
		newServerAgentResource,
		newWarningPoliciesResource,
	}
}

func uptimeMonitorQueryAttributes() []queryAttribute {
	return []queryAttribute{{Name: "page", Description: "Page number."}, {Name: "per_page", Description: "Results per page."}, {Name: "id", Description: "Monitor ID."}, {Name: "name", Description: "Name filter."}, {Name: "target", Description: "Target filter."}, {Name: "category", Description: "Category filter."}, {Name: "type", Description: "Monitor type."}, {Name: "uptime_status", Description: "Uptime status."}, {Name: "monitor_status", Description: "Monitor status."}, {Name: "order", Description: "Sort order."}, {Name: "order_by", Description: "Sort field."}}
}

func blacklistMonitorQueryAttributes() []queryAttribute {
	return []queryAttribute{{Name: "page", Description: "Page number."}, {Name: "per_page", Description: "Results per page."}, {Name: "name", Description: "Name filter."}, {Name: "exact_name", Description: "Exact name filter."}, {Name: "target", Description: "Target filter."}, {Name: "exact_target", Description: "Exact target filter."}, {Name: "cidr", Description: "CIDR prefix for target filter."}, {Name: "type", Description: "Monitor type."}, {Name: "listed", Description: "Listed status."}, {Name: "order", Description: "Sort order."}, {Name: "order_by", Description: "Sort field."}}
}
