FROM golang:1.20.3-alpine3.17 as builder

# Build the binary statically.
ENV CGO_ENABLED=0

WORKDIR /app
COPY . .
RUN cd cmd && go build -o ../mongo-writer .


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

RUN apk add -U --no-cache ca-certificates
RUN mkdir -p /var/log/app

WORKDIR /app
COPY --from=builder /app/mongo-writer .
COPY config.yml .
RUN chown -R kai:0 /app \
    && chmod -R g+w /app \
    && chown -R kai:0 /var/log/app \
    && chmod -R g+w /var/log/app

USER kai

CMD ["sh","-c","/app/mongo-writer 2>&1 | tee -a /var/log/app/app.log"]
