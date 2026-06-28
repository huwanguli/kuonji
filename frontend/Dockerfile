# Stage 1: build
FROM node:22-alpine AS builder

WORKDIR /src

COPY package.json package-lock.json ./
RUN npm install

COPY . .
RUN npm run build

# Stage 2: serve
FROM caddy:alpine

COPY --from=builder /src/dist /usr/share/caddy
COPY Caddyfile /etc/caddy/Caddyfile

EXPOSE 80 443
