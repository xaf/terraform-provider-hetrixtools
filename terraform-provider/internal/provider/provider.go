package provider

import (
	"context"
	"os"

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
		MarkdownDescription: "Terraform provider for HetrixTools.",
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
				return c.ListContactLists(ctx, query)
			})
		},
		func() datasource.DataSource {
			return newGetDataSource("blacklist_monitors", "List HetrixTools blacklist monitors.", blacklistMonitorQueryAttributes(), func(ctx context.Context, c *hetrixtools.Client, query map[string]string) (any, error) {
				return c.ListBlacklistMonitors(ctx, query)
			})
		},
		func() datasource.DataSource {
			return newGetDataSource("blacklist_report", "Get a HetrixTools blacklist monitor report.", []queryAttribute{{Name: "identifier", Description: "Blacklist monitor ID, IP address, or hostname.", Required: true}, {Name: "date", Description: "Report date in YYYY-MM-DD format."}}, func(ctx context.Context, c *hetrixtools.Client, query map[string]string) (any, error) {
				identifier := takeQueryValue(query, "identifier")
				return c.GetBlacklistMonitorReport(ctx, identifier, query)
			})
		},
		func() datasource.DataSource {
			return newGetDataSource("uptime_monitors", "List HetrixTools uptime monitors.", uptimeMonitorQueryAttributes(), func(ctx context.Context, c *hetrixtools.Client, query map[string]string) (any, error) {
				return c.ListUptimeMonitors(ctx, query)
			})
		},
		func() datasource.DataSource {
			return newGetDataSource("uptime_report", "Get a HetrixTools uptime monitor report.", []queryAttribute{{Name: "monitor_id", Description: "Uptime monitor ID.", Required: true}, {Name: "days", Description: "Number of recent days to display."}, {Name: "month", Description: "Month in YYYY-MM format."}, {Name: "timezone", Description: "Report timezone."}, {Name: "hourly_stats", Description: "Whether to include hourly stats."}}, func(ctx context.Context, c *hetrixtools.Client, query map[string]string) (any, error) {
				monitorID := takeQueryValue(query, "monitor_id")
				return c.GetUptimeMonitorReport(ctx, monitorID, query)
			})
		},
		func() datasource.DataSource {
			return newGetDataSource("uptime_downtimes", "List HetrixTools uptime monitor downtimes.", []queryAttribute{{Name: "monitor_id", Description: "Uptime monitor ID.", Required: true}, {Name: "page", Description: "Page number."}, {Name: "per_page", Description: "Results per page."}, {Name: "start_before", Description: "Only downtimes started before this timestamp."}, {Name: "start_after", Description: "Only downtimes started after this timestamp."}}, func(ctx context.Context, c *hetrixtools.Client, query map[string]string) (any, error) {
				monitorID := takeQueryValue(query, "monitor_id")
				return c.ListUptimeMonitorDowntimes(ctx, monitorID, query)
			})
		},
		func() datasource.DataSource {
			return newGetDataSource("uptime_location_fail_log", "Get a HetrixTools uptime monitor location fail log.", []queryAttribute{{Name: "monitor_id", Description: "Uptime monitor ID.", Required: true}, {Name: "timestamp", Description: "Timestamp to start scanning from."}, {Name: "minutes", Description: "Number of minutes with log entries to return."}}, func(ctx context.Context, c *hetrixtools.Client, query map[string]string) (any, error) {
				monitorID := takeQueryValue(query, "monitor_id")
				return c.GetUptimeMonitorLocationFailLog(ctx, monitorID, query)
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
				return c.ListStatusPages(ctx, query)
			})
		},
		func() datasource.DataSource {
			return newGetDataSource("scheduled_maintenances", "List HetrixTools scheduled maintenances.", []queryAttribute{{Name: "page", Description: "Page number."}, {Name: "per_page", Description: "Results per page."}, {Name: "monitor_id", Description: "Optional monitor ID filter."}}, func(ctx context.Context, c *hetrixtools.Client, query map[string]string) (any, error) {
				return c.ListScheduledMaintenances(ctx, query)
			})
		},
	}
}

func takeQueryValue(query map[string]string, key string) string {
	value := query[key]
	delete(query, key)
	return value
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
