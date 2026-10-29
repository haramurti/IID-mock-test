package leads

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/haramurti/IID-mock-test/backend/data"
)

// MockSource serves leads from the embedded data/leads.json.
type MockSource struct {
	leads []Lead
}

func NewMockSource() (*MockSource, error) {
	var ls []Lead
	if err := json.Unmarshal(data.LeadsJSON, &ls); err != nil {
		return nil, fmt.Errorf("parse embedded leads.json: %w", err)
	}
	return &MockSource{leads: ls}, nil
}

func (m *MockSource) ListLeads(ctx context.Context) ([]Lead, error) {
	out := make([]Lead, len(m.leads))
	copy(out, m.leads) // callers can't mutate our copy
	return out, nil
}
