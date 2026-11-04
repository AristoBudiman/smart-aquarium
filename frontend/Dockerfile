FROM node:22-alpine AS build
WORKDIR /src
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build-only

FROM caddy:2-alpine
COPY --from=build /src/dist /srv

