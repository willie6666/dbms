FROM docker.io/library/golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/vapor-auror-backend .

FROM docker.io/library/alpine:3.22

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=build /app/vapor-auror-backend /app/vapor-auror-backend
COPY --from=build /app/assets /app/assets

EXPOSE 8000
CMD ["/app/vapor-auror-backend"]
