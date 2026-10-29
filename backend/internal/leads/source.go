package leads

import (
	"context"
	"fmt"
)

// LeadSource is the only thing the rest of the app knows about where leads
// come from. Swap MockSource for ExpaSource without touching metrics or HTTP code.
type LeadSource interface {
	ListLeads(ctx context.Context) ([]Lead, error)
}

// NewSource picks a source by name ("mock" or "expa").
func NewSource(name, expaToken string) (LeadSource, error) {
	switch name {
	case "mock":
		return NewMockSource()
	case "expa":
		return NewExpaSource(expaToken), nil
	default:
		return nil, fmt.Errorf("unknown data source %q", name)
	}
}
