package metrics

import (
	"testing"
	"time"

	"github.com/haramurti/IID-mock-test/backend/internal/leads"
)

func ts(s string) *time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return &t
}

func lead(id string, st leads.LeadStatus, applied string) leads.Lead {
	l := leads.Lead{ID: id, Status: st}
	if applied != "" {
		l.AppliedAt = ts(applied)
	}
	return l
}

func mustRange(t *testing.T, from, to string) Range {
	t.Helper()
	r, err := ParseRange(from, to)
	if err != nil {
		t.Fatalf("ParseRange(%q, %q): %v", from, to, err)
	}
	return r
}

func TestParseRangeErrors(t *testing.T) {
	tests := []struct {
		name, from, to string
		wantErr        bool
	}{
		{"both empty is all time", "", "", false},
		{"only from", "2026-10-01", "", false},
		{"only to", "", "2026-10-01", false},
		{"same day", "2026-10-01", "2026-10-01", false},
		{"bad from format", "01-10-2026", "", true},
		{"bad to format", "", "2026/10/01", true},
		{"impossible date", "2026-02-30", "", true},
		{"from after to", "2026-10-02", "2026-10-01", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseRange(tt.from, tt.to)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRangeContains(t *testing.T) {
	// 2026-10-01 00:00 WIB == 2026-09-30T17:00:00Z
	tests := []struct {
		name     string
		from, to string
		lead     leads.Lead
		want     bool
	}{
		{"no range includes open lead", "", "", lead("1", leads.StatusOpen, ""), true},
		{"active range excludes open lead", "2026-10-01", "2026-10-31", lead("1", leads.StatusOpen, ""), false},
		{"same-day range, inside", "2026-10-01", "2026-10-01", lead("1", leads.StatusApplied, "2026-10-01T05:00:00Z"), true},
		{"same-day range, day before", "2026-10-01", "2026-10-01", lead("1", leads.StatusApplied, "2026-09-30T05:00:00Z"), false},
		{"WIB midnight start, 1 min before", "2026-10-01", "", lead("1", leads.StatusApplied, "2026-09-30T16:59:00Z"), false},
		{"WIB midnight start, exactly", "2026-10-01", "", lead("1", leads.StatusApplied, "2026-09-30T17:00:00Z"), true},
		{"WIB end of to day, last second", "", "2026-10-01", lead("1", leads.StatusApplied, "2026-10-01T16:59:59Z"), true},
		{"WIB end of to day, next midnight", "", "2026-10-01", lead("1", leads.StatusApplied, "2026-10-01T17:00:00Z"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mustRange(t, tt.from, tt.to).Contains(tt.lead); got != tt.want {
				t.Fatalf("Contains = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSummarize(t *testing.T) {
	in := []leads.Lead{
		lead("o1", leads.StatusOpen, ""),
		lead("a1", leads.StatusApplied, "2026-10-02T03:00:00Z"),
		lead("a2", leads.StatusApplied, "2026-10-03T03:00:00Z"),
		lead("c1", leads.StatusAccepted, "2026-10-04T03:00:00Z"),
		lead("p1", leads.StatusApproved, "2026-10-05T03:00:00Z"),
		lead("r1", leads.StatusRealized, "2026-10-06T03:00:00Z"),
		lead("f1", leads.StatusFinished, "2026-10-07T03:00:00Z"),
		lead("x1", leads.StatusRejected, "2026-10-08T03:00:00Z"),
		lead("old", leads.StatusCompleted, "2026-08-01T03:00:00Z"), // outside October
	}

	tests := []struct {
		name                    string
		from, to                string
		wantTotal, wantApproved int
		wantRealized            int
		wantRate                *float64
		wantByStatus            map[leads.LeadStatus]int
	}{
		{
			name: "October only", from: "2026-10-01", to: "2026-10-31",
			wantTotal: 7, wantApproved: 3, wantRealized: 2, wantRate: f(42.86),
			wantByStatus: map[leads.LeadStatus]int{
				leads.StatusApplied: 2, leads.StatusAccepted: 1, leads.StatusApproved: 1,
				leads.StatusRealized: 1, leads.StatusFinished: 1, leads.StatusRejected: 1,
			},
		},
		{
			name:      "all time counts every application, never open",
			wantTotal: 8, wantApproved: 4, wantRealized: 3, wantRate: f(50),
			wantByStatus: map[leads.LeadStatus]int{
				leads.StatusApplied: 2, leads.StatusAccepted: 1, leads.StatusApproved: 1,
				leads.StatusRealized: 1, leads.StatusFinished: 1, leads.StatusCompleted: 1, leads.StatusRejected: 1,
			},
		},
		{
			name: "zero applications gives nil rate", from: "2025-01-01", to: "2025-01-31",
			wantTotal: 0, wantApproved: 0, wantRealized: 0, wantRate: nil,
			wantByStatus: map[leads.LeadStatus]int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Summarize(in, mustRange(t, tt.from, tt.to))
			if s.TotalApplications != tt.wantTotal || s.Approved != tt.wantApproved || s.Realized != tt.wantRealized {
				t.Fatalf("total/approved/realized = %d/%d/%d, want %d/%d/%d",
					s.TotalApplications, s.Approved, s.Realized, tt.wantTotal, tt.wantApproved, tt.wantRealized)
			}
			if !eqRate(s.ConversionRate, tt.wantRate) {
				t.Fatalf("rate = %v, want %v", deref(s.ConversionRate), deref(tt.wantRate))
			}
			if _, ok := s.ByStatus[leads.StatusOpen]; ok {
				t.Fatal("byStatus must not contain open")
			}
			for _, st := range leads.AllStatuses {
				if st == leads.StatusOpen {
					continue
				}
				if s.ByStatus[st] != tt.wantByStatus[st] {
					t.Errorf("byStatus[%s] = %d, want %d", st, s.ByStatus[st], tt.wantByStatus[st])
				}
			}
		})
	}
}

func TestConversionRate(t *testing.T) {
	tests := []struct {
		approved, applied int
		want              *float64
	}{
		{0, 0, nil},
		{5, 0, nil},
		{0, 10, f(0)},
		{1, 3, f(33.33)},
		{2, 3, f(66.67)},
		{72, 182, f(39.56)},
		{10, 10, f(100)},
	}
	for _, tt := range tests {
		if got := ConversionRate(tt.approved, tt.applied); !eqRate(got, tt.want) {
			t.Errorf("ConversionRate(%d, %d) = %v, want %v", tt.approved, tt.applied, deref(got), deref(tt.want))
		}
	}
}

func f(v float64) *float64 { return &v }

func eqRate(a, b *float64) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

func deref(p *float64) any {
	if p == nil {
		return "nil"
	}
	return *p
}
