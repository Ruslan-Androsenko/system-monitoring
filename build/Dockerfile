# Собираем в гошке
FROM golang:1.20 as build

RUN apt-get update \
    && apt-get -y --no-install-recommends install sysstat iproute2 net-tools \
    && apt-get clean; rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /usr/share/doc/*

ENV BIN_FILE /opt/system-monitoring/app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/server/*

# На выходе тонкий образ
#FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="system-monitoring"
LABEL MAINTAINERS="ruslan.developer.89@gmail.com"

#ENV BIN_FILE "/opt/system-monitoring/app"
#COPY --from=build ${BIN_FILE} ${BIN_FILE}

ARG SERVER_PORT
ENV SERVER_HOST 0.0.0.0
ENV CONFIG_FILE /etc/system-monitoring/config.toml
COPY ./configs/config.toml ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE} -host ${SERVER_HOST} -port ${SERVER_PORT}
