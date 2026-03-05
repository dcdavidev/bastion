# Stage 1: Build the React frontend
FROM node:24-alpine AS frontend-builder

WORKDIR /app

# Enable corepack and copy dependency files first for caching
RUN corepack enable
COPY pnpm-lock.yaml pnpm-workspace.yaml package.json ./
# Copy workspace package.json files
COPY apps/web/package.json ./apps/web/
COPY apps/cli/package.json ./apps/cli/
COPY apps/server/package.json ./apps/server/
COPY packages/api/package.json ./packages/api/
COPY packages/auth/package.json ./packages/auth/
COPY packages/config/package.json ./packages/config/
COPY packages/crypto/package.json ./packages/crypto/
COPY packages/db/package.json ./packages/db/
COPY packages/models/package.json ./packages/models/
COPY packages/version/package.json ./packages/version/

RUN CI=true pnpm install --frozen-lockfile --ignore-scripts

# Copy source and build
COPY . .
RUN pnpm build --filter @dcdavidev/bastion-web

# Stage 2: Build the Go backend
FROM golang:1.25-alpine AS backend-builder

WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bastion-server ./apps/server/main.go

# Stage 3: Final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary
COPY --from=backend-builder /app/bastion-server .
# Copy migrations
COPY --from=backend-builder /app/packages/db/migrations ./packages/db/migrations
# Copy frontend build
COPY --from=frontend-builder /app/apps/web/build/client ./ui

# Set environment variables
ENV BASTION_PORT=8287
ENV BASTION_UI_DIR=/app/ui

# Expose the unified port
EXPOSE 8287

# Command to run
CMD ["./bastion-server"]
