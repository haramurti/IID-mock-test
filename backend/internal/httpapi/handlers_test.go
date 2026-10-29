package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/haramurti/IID-mock-test/backend/internal/leads"
)

type fakeSource struct {
	leads []leads.Lead
	calls int
}

func (f *fakeSource) ListLeads(ctx context.Context) ([]leads.Lead, error) {
	f.calls++
	return f.leads, nil
}

// 25 applied leads on consecutive days from 2026-10-01, plus 1 open lead.
func newTestApp() (*fiber.App, *fakeSource) {
	src := &fakeSource{}
	base := time.Date(2026, 10, 1, 3, 0, 0, 0, time.UTC)
	for i := 0; i < 25; i++ {
		t := base.AddDate(0, 0, i)
		st := leads.StatusApplied
		if i%5 == 0 {
			st = leads.StatusApproved
		}
		src.leads = append(src.leads, leads.Lead{
			ID: fmt.Sprintf("L%03d", i), Name: fmt.Sprintf("Lead %d", i), Status: st, AppliedAt: &t, CreatedAt: t,
		})
	}
	src.leads = append(src.leads, leads.Lead{ID: "OPEN", Name: "Open Lead", Status: leads.StatusOpen, CreatedAt: base})

	app := fiber.New()
	New(src).Register(app)
	return app, src
}

func get(t *testing.T, app *fiber.App, url string, out any) int {
	t.Helper()
	resp, err := app.Test(httptest.NewRequest("GET", url, nil))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			t.Fatalf("decode %s: %v", url, err)
		}
	}
	return resp.StatusCode
}

func TestInvalidParamsReturn400(t *testing.T) {
	app, _ := newTestApp()
	urls := []string{
		"/api/metrics?from=2026-13-01",
		"/api/metrics?from=yesterday",
		"/api/metrics?from=2026-10-10&to=2026-10-01",
		"/api/leads?to=10/01/2026",
		"/api/leads?page=0",
		"/api/leads?page=abc",
		"/api/leads?per_page=101",
		"/api/leads?per_page=0",
		"/api/leads?status=hired",
	}
	for _, u := range urls {
		t.Run(u, func(t *testing.T) {
			var body map[string]string
			if code := get(t, app, u, &body); code != fiber.StatusBadRequest {
				t.Fatalf("status = %d, want 400", code)
			}
			if body["error"] == "" {
				t.Fatalf("missing JSON error message, got %v", body)
			}
		})
	}
}

func TestMetricsEndpointAndCache(t *testing.T) {
	app, src := newTestApp()
	var m metricsResponse
	if code := get(t, app, "/api/metrics?from=2026-10-01&to=2026-10-10", &m); code != 200 {
		t.Fatalf("status = %d", code)
	}
	if m.TotalApplications != 10 || m.Approved != 2 || m.ConversionRate == nil || *m.ConversionRate != 20 {
		t.Fatalf("got total=%d approved=%d rate=%v", m.TotalApplications, m.Approved, m.ConversionRate)
	}
	if m.Range.Timezone != "Asia/Jakarta" {
		t.Fatalf("timezone = %q", m.Range.Timezone)
	}
	get(t, app, "/api/metrics?from=2026-10-01&to=2026-10-10", nil)
	if src.calls != 1 {
		t.Fatalf("source called %d times, want 1 (second request should hit cache)", src.calls)
	}
}

func TestLeadsPagingAndFilters(t *testing.T) {
	tests := []struct {
		url                string
		wantLen, wantTotal int
		wantPages          int
		wantFirstID        string
	}{
		{"/api/leads", 20, 26, 2, "L024"}, // newest first, open lead included
		{"/api/leads?page=2", 6, 26, 2, "L004"},
		{"/api/leads?page=9", 0, 26, 2, ""},                             // past the end: empty, not an error
		{"/api/leads?from=2026-10-01&to=2026-10-31", 20, 25, 2, "L024"}, // range excludes open
		{"/api/leads?per_page=100", 26, 26, 1, "L024"},
		{"/api/leads?status=approved", 5, 5, 1, "L020"},
		{"/api/leads?q=lead%201", 11, 11, 1, "L019"}, // "Lead 1", "Lead 10".."Lead 19"
	}
	app, _ := newTestApp()
	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			var r leadsResponse
			if code := get(t, app, tt.url, &r); code != 200 {
				t.Fatalf("status = %d", code)
			}
			if len(r.Data) != tt.wantLen || r.Paging.TotalItems != tt.wantTotal || r.Paging.TotalPages != tt.wantPages {
				t.Fatalf("len=%d total=%d pages=%d, want %d/%d/%d",
					len(r.Data), r.Paging.TotalItems, r.Paging.TotalPages, tt.wantLen, tt.wantTotal, tt.wantPages)
			}
			if tt.wantFirstID != "" && r.Data[0].ID != tt.wantFirstID {
				t.Fatalf("first id = %s, want %s", r.Data[0].ID, tt.wantFirstID)
			}
		})
	}
}
