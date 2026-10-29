// Package metrics holds all dashboard business logic: the date filter,
// counts per status, and the conversion rate. The frontend only displays
// what this package computes.
//
// Rules (also documented in the README):
//   - The date filter uses AppliedAt. Leads without AppliedAt (status open)
//     are excluded when a range is active.
//   - from/to are dates in Asia/Jakarta (WIB); the whole `to` day is included.
//   - Funnel counting: a later stage also counts as the earlier ones.
//     Approved = approved + realized + finished + completed
//     Realized = realized + finished + completed
//     Applied  = every lead with AppliedAt in range
//   - ConversionRate = Approved / Applied * 100, rounded to 2 decimals;
//     nil when Applied = 0.
package metrics

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/haramurti/IID-mock-test/backend/internal/leads"
)

const DateLayout = "2006-01-02"

// WIB is the timezone all date ranges are interpreted in.
var WIB = mustLoadLocation("Asia/Jakarta")

func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		// Fallback for systems without tzdata; WIB has no DST.
		return time.FixedZone("WIB", 7*60*60)
	}
	return loc
}

// Range is a date range on AppliedAt. A nil bound means "unbounded".
// Start is inclusive, End is exclusive (start of the day after `to`).
type Range struct {
	From, To   string // original YYYY-MM-DD inputs ("" if not set)
	start, end *time.Time
}

// Active reports whether any bound is set.
func (r Range) Active() bool { return r.start != nil || r.end != nil }

// ParseRange validates from/to query values. Empty strings mean unbounded.
func ParseRange(from, to string) (Range, error) {
	r := Range{From: from, To: to}
	if from != "" {
		t, err := time.ParseInLocation(DateLayout, from, WIB)
		if err != nil {
			return Range{}, fmt.Errorf("invalid 'from' date %q: use YYYY-MM-DD", from)
		}
		r.start = &t
	}
	if to != "" {
		t, err := time.ParseInLocation(DateLayout, to, WIB)
		if err != nil {
			return Range{}, fmt.Errorf("invalid 'to' date %q: use YYYY-MM-DD", to)
		}
		end := t.AddDate(0, 0, 1) // include the whole `to` day
		r.end = &end
	}
	if r.start != nil && r.end != nil && !r.start.Before(*r.end) {
		return Range{}, errors.New("'from' must be on or before 'to'")
	}
	return r, nil
}

// Contains reports whether the lead's AppliedAt falls in the range.
// With no active range every lead matches (including open ones).
func (r Range) Contains(l leads.Lead) bool {
	if !r.Active() {
		return true
	}
	if l.AppliedAt == nil {
		return false
	}
	t := *l.AppliedAt
	if r.start != nil && t.Before(*r.start) {
		return false
	}
	if r.end != nil && !t.Before(*r.end) {
		return false
	}
	return true
}

// Filter returns the leads inside the range, keeping input order.
func Filter(ls []leads.Lead, r Range) []leads.Lead {
	out := make([]leads.Lead, 0, len(ls))
	for _, l := range ls {
		if r.Contains(l) {
			out = append(out, l)
		}
	}
	return out
}

type Summary struct {
	TotalApplications int                      `json:"totalApplications"`
	ByStatus          map[leads.LeadStatus]int `json:"byStatus"`
	Approved          int                      `json:"approved"`
	Realized          int                      `json:"realized"`
	ConversionRate    *float64                 `json:"conversionRate"` // null when no applications
}

// Summarize computes the KPI numbers for the leads inside r.
func Summarize(ls []leads.Lead, r Range) Summary {
	s := Summary{ByStatus: map[leads.LeadStatus]int{}}
	for _, st := range leads.AllStatuses {
		if st != leads.StatusOpen {
			s.ByStatus[st] = 0 // stable keys, even when zero
		}
	}
	for _, l := range ls {
		if l.AppliedAt == nil || !r.Contains(l) {
			continue // never applied: not an application
		}
		s.TotalApplications++
		s.ByStatus[l.Status]++
		switch l.Status {
		case leads.StatusApproved:
			s.Approved++
		case leads.StatusRealized, leads.StatusFinished, leads.StatusCompleted:
			s.Approved++
			s.Realized++
		}
	}
	s.ConversionRate = ConversionRate(s.Approved, s.TotalApplications)
	return s
}

// ConversionRate returns approved/applied*100 rounded to 2 decimals, or nil if applied is 0.
func ConversionRate(approved, applied int) *float64 {
	if applied == 0 {
		return nil
	}
	v := math.Round(float64(approved)/float64(applied)*100*100) / 100
	return &v
}
