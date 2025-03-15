# syntax=docker/dockerfile:1

FROM golang:latest AS builder

WORKDIR /build


COPY . . 



RUN go mod download
RUN GOOS=linux go build -o . .

FROM gcr.io/distroless/base-debian12


WORKDIR /app
COPY --from=builder /build/ssh-server .
COPY --from=builder /build/markdown ./markdown
# COPY --from=builder /build/.env ./.env  


# COPY --from=builder /build/web/markdown ./web/markdown
# COPY --from=builder /build/web/static ./web/static



# EXPOSE 8080
EXPOSE 2225





CMD ["/app/ssh-server"]


