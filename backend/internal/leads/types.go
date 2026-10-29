// Package leads defines the lead model and where leads come from.
package leads

import "time"

// LeadStatus mirrors common EXPA application stages.
// To be confirmed against the live EXPA schema (Altair).
type LeadStatus string

const (
	StatusOpen      LeadStatus = "open" // signed up, no application yet
	StatusApplied   LeadStatus = "applied"
	StatusAccepted  LeadStatus = "accepted"
	StatusApproved  LeadStatus = "approved"
	StatusRealized  LeadStatus = "realized"
	StatusFinished  LeadStatus = "finished"
	StatusCompleted LeadStatus = "completed"
	StatusRejected  LeadStatus = "rejected"
	StatusWithdrawn LeadStatus = "withdrawn"
)

// AllStatuses lists every status in funnel order.
var AllStatuses = []LeadStatus{
	StatusOpen, StatusApplied, StatusAccepted, StatusApproved, StatusRealized,
	StatusFinished, StatusCompleted, StatusRejected, StatusWithdrawn,
}

// Valid reports whether s is a known status.
func (s LeadStatus) Valid() bool {
	for _, v := range AllStatuses {
		if s == v {
			return true
		}
	}
	return false
}

type Lead struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`      // fake names only
	Email      string     `json:"email"`     // fake, e.g. lead123@example.com
	Programme  string     `json:"programme"` // "GV" | "GTa" | "GTe"
	HomeLC     string     `json:"homeLC"`    // e.g. "AIESEC in UB"
	Status     LeadStatus `json:"status"`
	AppliedAt  *time.Time `json:"appliedAt"` // the field the date filter uses
	ApprovedAt *time.Time `json:"approvedAt"`
	RealizedAt *time.Time `json:"realizedAt"`
	CreatedAt  time.Time  `json:"createdAt"` // sign-up date
}
