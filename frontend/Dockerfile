# Build stage
FROM node:22-alpine AS builder

# Install pnpm
RUN npm install -g pnpm

WORKDIR /app

# Copy package files
COPY package.json pnpm-lock.yaml* ./

# Install dependencies
RUN pnpm install

# Copy source code
COPY . .

# Build the application
RUN pnpm run build

# Run stage
FROM node:22-alpine

RUN npm install -g pnpm

WORKDIR /app

COPY --from=builder /app/package.json /app/pnpm-lock.yaml* ./
RUN pnpm install --prod

COPY --from=builder /app/build ./build
COPY --from=builder /app/react-router.config.ts ./

# Expose port
EXPOSE 3000

# Start the application
CMD ["pnpm", "start"]
