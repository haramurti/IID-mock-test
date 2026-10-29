// Package httpapi wires the HTTP routes under /api.
package httpapi

import (
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/haramurti/IID-mock-test/backend/internal/leads"
	"github.com/haramurti/IID-mock-test/backend/internal/metrics"
)

const (
	defaultPerPage = 20
	maxPerPage     = 100
	metricsTTL     = 5 * time.Minute
)

// API holds the dependencies the handlers need.
type API struct {
	source leads.LeadSource
	cache  *ttlCache[metricsResponse]
}

func New(source leads.LeadSource) *API {
	return &API{source: source, cache: newTTLCache[metricsResponse](metricsTTL)}
}

// Register mounts all API routes on the given router.
func (a *API) Register(r fiber.Router) {
	api := r.Group("/api")
	api.Get("/health", health)
	api.Get("/metrics", a.metrics)
	api.Get("/leads", a.leads)
	api.Use(func(c *fiber.Ctx) error { // unknown /api/* path
		return jsonError(c, fiber.StatusNotFound, "not found")
	})
}

func health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

type rangeInfo struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Timezone string `json:"timezone"`
}

type metricsResponse struct {
	Range rangeInfo `json:"range"`
	metrics.Summary
	LastUpdated time.Time `json:"lastUpdated"`
}

// GET /api/metrics?from=YYYY-MM-DD&to=YYYY-MM-DD
func (a *API) metrics(c *fiber.Ctx) error {
	r, err := metrics.ParseRange(c.Query("from"), c.Query("to"))
	if err != nil {
		return jsonError(c, fiber.StatusBadRequest, err.Error())
	}

	key := r.From + "|" + r.To
	if cached, ok := a.cache.Get(key); ok {
		c.Set("X-Cache", "HIT")
		return c.JSON(cached)
	}

	all, err := a.source.ListLeads(c.UserContext())
	if err != nil {
		log.Printf("metrics: list leads: %v", err)
		return jsonError(c, fiber.StatusBadGateway, "could not load leads")
	}
	resp := metricsResponse{
		Range:       rangeInfo{From: r.From, To: r.To, Timezone: metrics.WIB.String()},
		Summary:     metrics.Summarize(all, r),
		LastUpdated: time.Now().UTC().Truncate(time.Second),
	}
	a.cache.Set(key, resp)
	c.Set("X-Cache", "MISS")
	return c.JSON(resp)
}

type paging struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

type leadsResponse struct {
	Data   []leads.Lead `json:"data"`
	Paging paging       `json:"paging"`
}

// GET /api/leads?from&to&status&q&page&per_page
func (a *API) leads(c *fiber.Ctx) error {
	r, err := metrics.ParseRange(c.Query("from"), c.Query("to"))
	if err != nil {
		return jsonError(c, fiber.StatusBadRequest, err.Error())
	}
	page, err := intParam(c, "page", 1, 1, 1<<30)
	if err != nil {
		return jsonError(c, fiber.StatusBadRequest, err.Error())
	}
	perPage, err := intParam(c, "per_page", defaultPerPage, 1, maxPerPage)
	if err != nil {
		return jsonError(c, fiber.StatusBadRequest, err.Error())
	}
	status := leads.LeadStatus(strings.ToLower(c.Query("status")))
	if status != "" && !status.Valid() {
		return jsonError(c, fiber.StatusBadRequest, "invalid 'status' "+strconv.Quote(string(status)))
	}
	q := strings.ToLower(strings.TrimSpace(c.Query("q")))

	all, err := a.source.ListLeads(c.UserContext())
	if err != nil {
		log.Printf("leads: list leads: %v", err)
		return jsonError(c, fiber.StatusBadGateway, "could not load leads")
	}

	matched := make([]leads.Lead, 0, len(all))
	for _, l := range metrics.Filter(all, r) {
		if status != "" && l.Status != status {
			continue
		}
		if q != "" && !strings.Contains(strings.ToLower(l.Name), q) && !strings.Contains(strings.ToLower(l.ID), q) {
			continue
		}
		matched = append(matched, l)
	}
	sortNewestFirst(matched)

	total := len(matched)
	start := min((page-1)*perPage, total)
	end := min(start+perPage, total)
	return c.JSON(leadsResponse{
		Data: matched[start:end],
		Paging: paging{
			Page: page, PerPage: perPage, TotalItems: total,
			TotalPages: (total + perPage - 1) / perPage,
		},
	})
}

// sortNewestFirst orders by AppliedAt desc; leads without an application go last (by CreatedAt desc).
func sortNewestFirst(ls []leads.Lead) {
	sort.SliceStable(ls, func(i, j int) bool {
		a, b := ls[i].AppliedAt, ls[j].AppliedAt
		switch {
		case a != nil && b != nil:
			return a.After(*b)
		case a != nil:
			return true
		case b != nil:
			return false
		default:
			return ls[i].CreatedAt.After(ls[j].CreatedAt)
		}
	})
}

func intParam(c *fiber.Ctx, name string, def, lo, hi int) (int, error) {
	s := c.Query(name)
	if s == "" {
		return def, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < lo || n > hi {
		return 0, fiber.NewError(fiber.StatusBadRequest,
			"invalid '"+name+"': must be an integer between "+strconv.Itoa(lo)+" and "+strconv.Itoa(hi))
	}
	return n, nil
}

func jsonError(c *fiber.Ctx, code int, msg string) error {
	return c.Status(code).JSON(fiber.Map{"error": msg})
}
