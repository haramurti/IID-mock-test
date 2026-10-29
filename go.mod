// Root module exists only for the Vercel Go function in api/.
// All real code lives in ./backend (its own module).
module github.com/haramurti/IID-mock-test

go 1.25.5

require (
	github.com/gofiber/fiber/v2 v2.52.16
	github.com/haramurti/IID-mock-test/backend v0.0.0
)

require (
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.51.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
)

replace github.com/haramurti/IID-mock-test/backend => ./backend
