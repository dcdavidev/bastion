# Stage 1: Build the React frontend
FROM node:24-alpine AS frontend-builder

WORKDIR /app
COPY . .
# Use pnpm for faster builds in the monorepo
RUN corepack enable && pnpm install --frozen-lockfile --ignore-scripts
RUN pnpm build --filter @dcdavidev/bastion-web

# Stage 2: Build the Go backend
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app
COPY . .
# Download dependencies
RUN go mod download
# Build the server
RUN CGO_ENABLED=0 GOOS=linux go build -o bastion-server ./apps/server/main.go

# Stage 3: Final image
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary
COPY --from=backend-builder /app/bastion-server .
# Copy migrations
COPY --from=backend-builder /app/packages/core/db/migrations ./packages/core/db/migrations
# Copy frontend build
COPY --from=frontend-builder /app/apps/web/build/client ./ui

# Set environment variables
ENV BASTION_PORT=8287
ENV BASTION_UI_DIR=/app/ui

# Expose the unified port
EXPOSE 8287

# Command to run
CMD ["./bastion-server"]
