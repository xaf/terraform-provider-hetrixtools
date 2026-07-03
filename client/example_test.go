package hetrixtools_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"

	hetrixtools "github.com/xaf/terraform-provider-hetrixtools/client"
)

func ExampleClient_GetUptimeMonitor() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"uptime_monitors":[{"id":"monitor-1","name":"Homepage"}],"meta":{"pagination":{"current":1,"last":1}}}`))
	}))
	defer server.Close()

	client := hetrixtools.NewClientWithBaseURL(server.URL, "token")
	monitor, _ := client.GetUptimeMonitor(context.Background(), "monitor-1")

	fmt.Println(monitor.Name)
	// Output: Homepage
}

func ExampleError() {
	err := hetrixtools.Error{StatusCode: http.StatusForbidden, Body: `{"message":"forbidden"}`}

	var apiErr hetrixtools.Error
	if errors.As(err, &apiErr) {
		fmt.Println(apiErr.StatusCode)
	}
	// Output: 403
}

func ExampleClient_UpsertUptimeMonitor() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"status":"SUCCESS","monitor_id":"monitor-1"}`))
	}))
	defer server.Close()

	client := hetrixtools.NewClientWithBaseURL(server.URL+"/v3", "token")
	response, _ := client.UpsertUptimeMonitor(context.Background(), hetrixtools.UptimeMonitorRequest{
		MID:  "monitor-1",
		Name: "Homepage",
	})

	fmt.Println(response.Status)
	// Output: SUCCESS
}
