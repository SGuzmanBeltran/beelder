# Build stage
FROM oven/bun:1-alpine AS builder

WORKDIR /app

# Copy package files first
COPY package.json bun.lock* ./
RUN bun install --frozen-lockfile

# Copy source code
COPY . .

# Build args MUST be set right before the build command
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL

# Build the app - Vite will read VITE_API_URL from ENV
RUN bun run build

# Production stage with nginx
FROM nginx:alpine

# Copy built files
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx config (optional)
# COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
