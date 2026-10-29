.PHONY: mock test dev build run docker clean

# Regenerate backend/data/leads.json (fake leads, fixed seed).
mock:
	cd backend && go run ./cmd/genmock

# Backend unit tests (metrics + HTTP handlers).
test:
	cd backend && go test ./...

# Backend on :8080 + Vite on :5173 (open http://localhost:5173). Ctrl+C stops both.
dev:
	@trap 'kill 0' INT TERM; \
	(cd backend && go run ./cmd/server) & \
	(cd frontend && npm run dev) & \
	wait

# Single binary at bin/server with the React app embedded.
build:
	cd frontend && npm ci && npm run build
	rm -rf backend/web/dist && mkdir -p backend/web/dist && touch backend/web/dist/.gitkeep
	cp -R frontend/dist/. backend/web/dist/
	cd backend && CGO_ENABLED=0 go build -o ../bin/server ./cmd/server

# Run the built binary (open http://localhost:8080).
run: build
	./bin/server

docker:
	docker build -t lead-dashboard .

clean:
	rm -rf bin frontend/dist
	find backend/web/dist -mindepth 1 ! -name .gitkeep -delete
