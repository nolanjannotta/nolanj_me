# syntax=docker/dockerfile:1

FROM golang:latest AS builder

WORKDIR /build


COPY . . 



RUN go mod download
RUN GOOS=linux go build -o . .

FROM gcr.io/distroless/base-debian12


WORKDIR /app
COPY --from=builder /build/nolanj_me .
COPY --from=builder /build/static ./static
# COPY --from=builder /build/.env ./.env  


# COPY --from=builder /build/web/markdown ./web/markdown
# COPY --from=builder /build/web/static ./web/static



EXPOSE 8080
# EXPOSE 23234





CMD ["/app/nolanj_me"]


