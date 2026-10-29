// Package handler is the Vercel serverless entry point. vercel.json rewrites
// every /api/* request here; we hand it to the same Fiber app the binary uses.
package handler

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2/middleware/adaptor"

	"github.com/haramurti/IID-mock-test/backend/app"
)

var (
	once    sync.Once
	fiberFn http.HandlerFunc
	initErr error
)

// Handler is called by Vercel for each request. The Fiber app is built once per
// warm instance and reused, so the in-memory metrics cache lives as long as the instance.
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		a, _, err := app.FromEnv()
		if err != nil {
			initErr = err
			return
		}
		fiberFn = adaptor.FiberApp(a)
	})
	if initErr != nil {
		log.Printf("init: %v", initErr)
		http.Error(w, `{"error":"server misconfigured"}`, http.StatusInternalServerError)
		return
	}

	// vercel.json rewrites /api/<path> to /api/index?__path=<path>; restore the original path.
	q := r.URL.Query()
	if p := q.Get("__path"); p != "" {
		q.Del("__path")
		r.URL.Path = "/api/" + strings.TrimPrefix(p, "/")
		r.URL.RawPath = ""
		r.URL.RawQuery = q.Encode()
		r.RequestURI = r.URL.RequestURI()
	}
	fiberFn(w, r)
}
