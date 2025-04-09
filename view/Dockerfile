FROM node:18-alpine AS builder

ARG NEXT_PUBLIC_BASE_URL
ARG NEXT_PUBLIC_WEBSOCKET_URL
ARG PORT=7443

ENV NEXT_PUBLIC_BASE_URL=$NEXT_PUBLIC_BASE_URL
ENV NEXT_PUBLIC_WEBSOCKET_URL=$NEXT_PUBLIC_WEBSOCKET_URL
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=$PORT

WORKDIR /app

COPY package.json yarn.lock ./
RUN yarn install --frozen-lockfile --network-timeout 600000 --production=false

COPY . .

RUN yarn build && \
    yarn install --frozen-lockfile --network-timeout 600000 --production && \
    yarn cache clean

FROM node:18-alpine AS runner

ARG PORT=7443
ENV PORT=$PORT
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

RUN addgroup -g 1001 nixopus && \
    adduser -u 1001 -G nixopus -s /bin/sh -D nixopus && \
    apk add --no-cache docker-cli && \
    rm -rf /var/cache/apk/*

WORKDIR /app

COPY --from=builder --chown=nixopus:nixopus /app/next.config.ts ./
COPY --from=builder --chown=nixopus:nixopus /app/public ./public
COPY --from=builder --chown=nixopus:nixopus /app/.next ./.next
COPY --from=builder --chown=nixopus:nixopus /app/node_modules ./node_modules
COPY --from=builder --chown=nixopus:nixopus /app/package.json ./
COPY --from=builder --chown=nixopus:nixopus /app/yarn.lock ./

RUN mkdir -p .next/cache && \
    chown -R nixopus:nixopus . && \
    chmod -R 755 . && \
    rm -rf /app/node_modules/.cache

USER nixopus

EXPOSE ${PORT}

CMD ["yarn", "start"]