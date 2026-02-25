# Build stage
FROM node:22-alpine AS builder

RUN npm install -g pnpm

WORKDIR /app

COPY package.json pnpm-lock.yaml* ./
RUN pnpm install

COPY . .
RUN pnpm run build

# Run stage
FROM nginx:alpine

# Copy only the client-side build folder
COPY --from=builder /app/build/client /usr/share/nginx/html

# Add a basic nginx config to handle SPA routing
RUN echo 'server { \
    listen 3500; \
    location / { \
        root /usr/share/nginx/html; \
        index index.html; \
        try_files $uri $uri/ /index.html; \
    } \
}' > /etc/nginx/conf.d/default.conf

EXPOSE 3500

CMD ["nginx", "-g", "daemon off;"]
