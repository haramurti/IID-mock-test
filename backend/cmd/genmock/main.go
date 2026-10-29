// genmock writes data/leads.json with fake, reproducible leads.
//
//	go run ./cmd/genmock                     # anchor = today
//	go run ./cmd/genmock -anchor 2026-10-29  # fully reproducible output
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/haramurti/IID-mock-test/backend/internal/leads"
)

const day = 24 * time.Hour

var (
	firstNames = []string{"Andi", "Budi", "Citra", "Dewi", "Eka", "Fajar", "Gita", "Hana", "Indra", "Joko", "Kartika", "Lestari", "Made", "Nadia", "Oki", "Putri", "Rizky", "Sari", "Tono", "Wulan", "Yusuf", "Zahra"}
	lastNames  = []string{"Pratama", "Santoso", "Wijaya", "Saputra", "Lestari", "Nugroho", "Hidayat", "Kusuma", "Siregar", "Halim", "Rahman", "Setiawan"}
	programmes = []string{"GV", "GV", "GV", "GTa", "GTe"} // GV is most common
	homeLCs    = []string{"AIESEC in UB", "AIESEC in UI", "AIESEC in UGM", "AIESEC in ITB", "AIESEC in UNAIR", "AIESEC in UNDIP", "AIESEC in IPB", "AIESEC in UNPAD"}
)

// Funnel shape: weights per status (sum = 100).
var statusWeights = []struct {
	s leads.LeadStatus
	w int
}{
	{leads.StatusOpen, 20},
	{leads.StatusApplied, 25},
	{leads.StatusAccepted, 12},
	{leads.StatusApproved, 12},
	{leads.StatusRealized, 9},
	{leads.StatusFinished, 5},
	{leads.StatusCompleted, 4},
	{leads.StatusRejected, 8},
	{leads.StatusWithdrawn, 5},
}

func main() {
	n := flag.Int("n", 300, "number of leads")
	seed := flag.Int64("seed", 42, "random seed")
	anchorStr := flag.String("anchor", "", "newest possible date, YYYY-MM-DD (default: today)")
	out := flag.String("out", "data/leads.json", "output file")
	flag.Parse()

	anchor := time.Now().UTC()
	if *anchorStr != "" {
		t, err := time.Parse("2006-01-02", *anchorStr)
		if err != nil {
			log.Fatalf("bad -anchor: %v", err)
		}
		anchor = t.Add(23 * time.Hour) // late in the anchor day
	}

	r := rand.New(rand.NewSource(*seed))
	ls := make([]leads.Lead, 0, *n)
	for i := 1; i <= *n; i++ {
		ls = append(ls, makeLead(r, i, anchor))
	}

	b, err := json.MarshalIndent(ls, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(*out, append(b, '\n'), 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %d leads to %s\n", len(ls), *out)
}

func makeLead(r *rand.Rand, i int, anchor time.Time) leads.Lead {
	status := pickStatus(r)

	// Decide gaps between stages first, so the whole timeline fits before anchor.
	signupToApply := randDur(r, 0, 10)
	applyToApprove := randDur(r, 3, 20)
	approveToRealize := randDur(r, 10, 40)

	var span time.Duration
	switch status {
	case leads.StatusOpen:
	case leads.StatusApproved:
		span = signupToApply + applyToApprove
	case leads.StatusRealized, leads.StatusFinished, leads.StatusCompleted:
		span = signupToApply + applyToApprove + approveToRealize
	default: // applied, accepted, rejected, withdrawn
		span = signupToApply
	}

	// Sign-up somewhere in the last 6 months, leaving room for the span.
	window := 180*day - span
	created := anchor.Add(-span).Add(-time.Duration(r.Int63n(int64(window)))).Truncate(time.Minute)

	l := leads.Lead{
		ID:        fmt.Sprintf("L%05d", i),
		Name:      firstNames[r.Intn(len(firstNames))] + " " + lastNames[r.Intn(len(lastNames))],
		Email:     fmt.Sprintf("lead%d@example.com", i),
		Programme: programmes[r.Intn(len(programmes))],
		HomeLC:    homeLCs[r.Intn(len(homeLCs))],
		Status:    status,
		CreatedAt: created,
	}
	if status == leads.StatusOpen {
		return l // no application yet: appliedAt stays null
	}
	applied := created.Add(signupToApply)
	l.AppliedAt = &applied
	switch status {
	case leads.StatusApproved, leads.StatusRealized, leads.StatusFinished, leads.StatusCompleted:
		approved := applied.Add(applyToApprove)
		l.ApprovedAt = &approved
		if status != leads.StatusApproved {
			realized := approved.Add(approveToRealize)
			l.RealizedAt = &realized
		}
	}
	return l
}

func pickStatus(r *rand.Rand) leads.LeadStatus {
	x := r.Intn(100)
	for _, sw := range statusWeights {
		if x < sw.w {
			return sw.s
		}
		x -= sw.w
	}
	return leads.StatusOpen
}

// randDur returns a random duration between minDays and maxDays (plus random hours).
func randDur(r *rand.Rand, minDays, maxDays int) time.Duration {
	return time.Duration(minDays+r.Intn(maxDays-minDays+1))*day + time.Duration(r.Intn(24))*time.Hour
}
