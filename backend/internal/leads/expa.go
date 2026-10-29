package leads

import (
	"context"
	"errors"
)

// ExpaSource will load applications from the EXPA GraphQL API.
// STUB ONLY: not implemented and never called while DATA_SOURCE=mock.
//
// The token comes from the EXPA_TOKEN env var on the server (see config).
// It is sent only from this Go service to EXPA, never to the browser.
//
// Illustrative query (verify field names in Altair before implementing):
//
//	query Leads($page: Int, $from: DateTime, $to: DateTime) {
//	  allOpportunityApplication(
//	    page: $page
//	    per_page: 100
//	    filters: { created_at: { from: $from, to: $to } }
//	  ) {
//	    paging { total_items total_pages current_page }
//	    data { id status created_at date_approved date_realized }
//	  }
//	}
//
// Implementation plan: loop pages until current_page == total_pages,
// check len(results) == paging.total_items, map each item to Lead.
type ExpaSource struct {
	token string
}

func NewExpaSource(token string) *ExpaSource {
	return &ExpaSource{token: token}
}

func (e *ExpaSource) ListLeads(ctx context.Context) ([]Lead, error) {
	return nil, errors.New("ExpaSource not implemented yet; use DATA_SOURCE=mock")
}
