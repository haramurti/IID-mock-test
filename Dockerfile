# syntax=docker/dockerfile:1

# ---- 1. Build the React frontend ----
FROM node:20-alpine AS frontend
WORKDIR /src/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# ---- 2. Build the Go binary with the frontend embedded ----
FROM golang:1.25-alpine AS backend
WORKDIR /src/backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
COPY --from=frontend /src/frontend/dist/ ./web/dist/
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# ---- 3. Small runtime image: just the binary, runs as non-root ----
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=backend /out/server /server
ENV PORT=8080 DATA_SOURCE=mock
EXPOSE 8080
USER nonroot
ENTRYPOINT ["/server"]
