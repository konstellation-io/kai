FROM golang:1.20.3 as builder

# Build the binary statically.
#ENV CGO_ENABLED=0

#RUN apk add -U --no-cache btrfs-progs btrfs-progs-dev gpgme-dev device-mapper libassuan-dev
RUN apt update && apt install -y libbtrfs-dev libgpgme-dev libdevmapper-dev
#RUN apt-get install gcc pkgconf musl-dev glib-static libc-dev libseccomp-dev libselinux-dev ostree-dev

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o admin-api .
RUN #go build -o admin-api -tags=containers_image_openpgp .

FROM debian


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
#
RUN mkdir "/nonexistent"

#RUN apk add -U --no-cache ca-certificates btrfs-progs btrfs-progs-dev gpgme-dev device-mapper libassuan-dev
RUN apt update && apt install -y ca-certificates libbtrfs-dev libgpgme-dev libdevmapper-dev
RUN apt install -y containernetworking-plugins
#RUN apk add -U --no-cache  btrfs-progs btrfs-progs-dev gpgme-dev device-mapper libassuan-dev ostree-devgcc pkgconf musl-dev glib-static libc-dev libseccomp-dev libselinux-dev ostree-dev
RUN mkdir -p /var/log/app

WORKDIR /app
COPY --from=builder /app/admin-api .
COPY --from=builder /app/casbin_rbac_model.conf .
COPY --from=builder /app/casbin_rbac_policy.csv .
COPY config.yml .
#RUN apk add libcap && setcap 'cap_net_bind_service=+ep' /app/admin-api
RUN apt-get install -y libcap2-bin uidmap && setcap 'cap_net_bind_service=+ep' /app/admin-api
RUN setcap 'cap_setgid=ep' /usr/bin/newgidmap && setcap 'cap_setuid+ep' /usr/bin/newuidmap
RUN setcap 'cap_setgid=ep' /app/admin-api && setcap 'cap_setuid+ep' /app/admin-api

RUN chown -R kai:0 /app \
    && chmod -R g+w /app \
    && chown -R kai:0 /var/log/app \
    && chmod -R g+w /var/log/app \
    && chown -R kai:0 /nonexistent \
    && chmod -R g+w /nonexistent

USER kai

CMD ["sh","-c","/app/admin-api 2>&1 | tee -a /var/log/app/app.log"]
