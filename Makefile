.PHONY: mock test dev build

mock:
	cd backend && go run ./cmd/genmock

test:
	cd backend && go test ./...

dev:
	(cd backend && go run ./cmd/server) & (cd frontend && npm run dev); wait

build:
	cd frontend && npm ci && npm run build
	cd backend && go build -o ../bin/server ./cmd/server
