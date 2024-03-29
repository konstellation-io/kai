# Generate proto
FROM namely/protoc-all as protobuf

WORKDIR /app

COPY scripts scripts

COPY ./internal/infrastructure/grpc/proto ./internal/infrastructure/grpc/proto

RUN ./scripts/generate_proto.sh


# Build k8s-manager
FROM golang:1.22.0-alpine3.19 as builder

ENV CGO_ENABLED=0

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN rm -rf /app/internal/infrastructure/grpc/proto
COPY --from=protobuf /app/internal/infrastructure/grpc/proto/ /app/internal/infrastructure/grpc/proto/

RUN go build -o k8s-manager cmd/main.go

# Final image
FROM alpine:3.10.2

# Create kai user.
ENV USER=kai
ENV UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

RUN mkdir -p /var/log/app

WORKDIR /app
COPY --from=builder /app/k8s-manager .
COPY config.yml .
RUN chown -R kai:0 /app \
    && chmod -R g+w /app \
    && chown -R kai:0 /var/log/app \
    && chmod -R g+w /var/log/app

USER kai

CMD ["sh","-c","/app/k8s-manager 2>&1 | tee -a /var/log/app/app.log"]
