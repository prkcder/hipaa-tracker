# ======================
# Crawler - Production
# ======================
FROM node:22-alpine

WORKDIR /app


ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD=true \
    PUPPETEER_EXECUTABLE_PATH=/usr/bin/chromium-browser


RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    freetype-dev \
    harfbuzz \
    ca-certificates \
    ttf-freefont \
    curl


COPY package*.json ./

RUN npm install --omit=dev

COPY . .

RUN npm run build

EXPOSE 4000

CMD ["node", "dist/server.js"]