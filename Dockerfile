# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM node:20-alpine AS frontend-builder
WORKDIR /src/web/ui

COPY web/ui/package.json web/ui/yarn.lock ./
RUN corepack enable && yarn install --frozen-lockfile

COPY web/ui/ ./
RUN yarn build

FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS backend-builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-builder /src/web/ui/dist ./web/ui/dist

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags "-s -w" -o /out/vmcontroler .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=backend-builder /out/vmcontroler /app/vmcontroler

EXPOSE 8080
ENTRYPOINT ["/app/vmcontroler"]
