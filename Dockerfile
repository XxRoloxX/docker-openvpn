# Original credit: https://github.com/jpetazzo/dockvpn

# Smallest base image
FROM alpine:latest as base

LABEL maintainer="Kyle Manna <kyle@kylemanna.com>"

# Testing: pamtester
RUN echo "http://dl-cdn.alpinelinux.org/alpine/edge/testing/" >> /etc/apk/repositories && \
    apk add --update openvpn iptables bash easy-rsa openvpn-auth-pam google-authenticator pamtester libqrencode && \
    ln -s /usr/share/easy-rsa/easyrsa /usr/local/bin && \
    rm -rf /tmp/* /var/tmp/* /var/cache/apk/* /var/cache/distfiles/*

# Needed by scripts
ENV OPENVPN=/etc/openvpn
ENV EASYRSA=/usr/share/easy-rsa \
    EASYRSA_CRL_DAYS=3650 \
    EASYRSA_PKI=$OPENVPN/pki

VOLUME ["/etc/openvpn"]

# Internally uses port 1194/udp, remap using `docker run -p 443:1194/tcp`
EXPOSE 1194/udp

CMD ["ovpn_run"]

ADD ./bin /usr/local/bin
RUN chmod a+x /usr/local/bin/*

# Add support for OTP authentication using a PAM module
ADD ./otp/openvpn /etc/pam.d/

# Add Go dependencies in a multi-stage build for minimal size
FROM golang:1.23.4-alpine as builder

# Set up Go environment and working directory
WORKDIR /go/src/app

COPY ./go .

# Build the Go microservice
RUN cd ./managed-openvpn && go build -o /go/bin/microservice cmd/managed-openvpn/main.go

# Combine OpenVPN base image and Go service
FROM base

# Install Go runtime dependencies
RUN apk add --no-cache libc6-compat

# Copy Go binary from builder stage
COPY --from=builder /go/bin/microservice /usr/local/bin/microservice

# Entry point to run both OpenVPN and the Go microservice
CMD ["sh", "-c", "/usr/local/bin/ovpn_run & /usr/local/bin/microservice"]
