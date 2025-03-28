FROM ubuntu:20.04 AS builder

# Install dependencies
RUN echo "Installing dependencies"; \
  set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends \
    ca-certificates \
    openssl \
    curl \
		git \
    wget; \
  rm -rf /var/lib/apt/lists/*; \
  apt-get clean;

# Get github.com and proxy.golang.org certificates
ARG cert_location=/usr/local/share/ca-certificates
RUN echo "Getting required CA certifates"; \
  openssl s_client -showcerts -connect github.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/github.crt; \
  openssl s_client -showcerts -connect proxy.golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM >  ${cert_location}/proxy.golang.crt; \
  update-ca-certificates;

# Install Go
ENV GO_VERSION=1.13.4
ENV PATH=$PATH:/usr/local/go/bin
RUN echo "Installing Go"; \
  wget https://golang.org/dl/go$GO_VERSION.linux-amd64.tar.gz; \
  tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz; \
  go version; \
  rm go$GO_VERSION.linux-amd64.tar.gz

# Build program commands
# CMD_NAME will be passed through the docker compose file.
# Possible values are 'telebot', 'resource' and 'observer'.
# Default value is set to 'observer'.
ENV GOPATH=/go
ARG CMD_NAME=observer
WORKDIR $GOPATH/src/crypto-balance-bot
COPY . .
RUN echo "Building ${CMD_NAME}"; \
  set -ux; \
  # mkdir -p $GOPATH/src && cd $GOPATH/src; \
  # git clone https://github.com/psychoplasma/crypto-balance-bot.git; \
  # cd /crypto-balance-bot; \
  go mod download; \
  BUILDTIME=$(date --iso-8601=seconds); \
  GITCOMMIT=$(git rev-parse --short HEAD); \
  LDFLAGS="-X github.com/psychoplasma/crypto-balance-bot/common.version=master -X github.com/psychoplasma/crypto-balance-bot/common.gitcommit=${GITCOMMIT} -X github.com/psychoplasma/crypto-balance-bot/common.buildtime=${BUILDTIME}"; \
  cd cmd/$CMD_NAME; \
  go build -ldflags="-s -w ${LDFLAGS}"


FROM ubuntu:20.04
# CMD_NAME will be passed through the docker compose file.
# Possible values are 'telebot', 'resource' and 'observer'.
# Default value is set to 'observer'.
ARG CMD_NAME=observer

# ENV APP_WORKING_DIR=/cryptobalancebot
ENV APP_EXEC=$CMD_NAME

# Install ca-certificates
RUN echo "Installing ca-certificates"; \
  set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends ca-certificates; \
  rm -rf /var/lib/apt/lists/*; \
  apt-get clean;

# Create missing directories
# RUN mkdir -p $APP_WORKING_DIR

# Copy executables from the build stage
COPY --from=builder /go/src/crypto-balance-bot/cmd/$CMD_NAME/$CMD_NAME /usr/local/bin/
COPY docker/entrypoint.sh /usr/local/bin/

# Make entrypoint script and app executables executable
RUN chmod -R 755 /usr/local/bin/entrypoint.sh; \
  chmod -R 755 /usr/local/bin/$CMD_NAME

# VOLUME $APP_WORKING_DIR
# WORKDIR $APP_WORKING_DIR
VOLUME /cryptobalancebot
WORKDIR /cryptobalancebot
ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
