FROM golang:1.22.0-alpine3.19 as builder

# Build the binary statically.
ENV CGO_ENABLED=0

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o admin-api .


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
COPY --from=builder /app/admin-api .
COPY --from=builder /app/casbin_rbac_model.conf .
COPY --from=builder /app/casbin_rbac_policy.csv .
COPY config.yml .
RUN apk add libcap && setcap 'cap_net_bind_service=+ep' /app/admin-api
RUN chown -R kai:0 /app \
    && chmod -R g+w /app \
    && chown -R kai:0 /var/log/app \
    && chmod -R g+w /var/log/app

USER kai

CMD ["sh","-c","/app/admin-api 2>&1 | tee -a /var/log/app/app.log"]
