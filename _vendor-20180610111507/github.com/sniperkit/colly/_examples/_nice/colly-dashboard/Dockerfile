FROM golang:1.10-alpine AS builder

# container - meta
MAINTAINER <rosco_pecoltran@msn.com> Rosco Pecoltran

# container - deps
ARG APK_BUILD=${APK_BUILD:="make"}
ARG APK_DEBUG=${APK_DEBUG:=""}

# golang - args
ARG GO_SRC=${GO_PATH:-"/go/src"}

# repo - args
ARG REPO_VCS=${REPO_VCS:-"github.com"}
ARG REPO_OWNER=${REPO_OWNER:-"sniperkit"}
ARG REPO_NAME=${REPO_NAME:-"colly"}
ARG REPO_URI="${REPO_VCS}/${REPO_OWNER}/${REPO_NAME}"
ARG REPO_DIR_ABS="${GO_SRC}/${REPO_URI}"

# build - args
ARG BUILD_BIN_NAME=${BUILD_BIN_NAME:-"colly-dashboard"}
ARG BUILD_BIN_PREFIX_PATH=${BUILD_BIN_PREFIX_PATH:-"/usr/local/bin"}
ARG BUILD_BIN_FULL_PATH="${BUILD_BIN_PREFIX_PATH}/${BUILD_BIN_NAME}"
ARG BUILD_BIN_DIR=${BUILD_BIN_DIR:-"/usr/local/bin"}
ARG BUILD_PREFIX_PATH=${BUILD_PREFIX_PATH:-""}
ARG BUILD_SUFFIX_PATH=${BUILD_SUFFIX_PATH:-"./_example/_nice/colly-dashboard"}
ARG BUILD_DIR_WORK="${REPO_DIR_ABS}/${BUILD_PREFIX_PATH}${BUILD_SUFFIX_PATH}"

RUN apk add --no-cache make

COPY .  ${REPO_DIR_ABS}
WORKDIR ${REPO_DIR_ABS}

RUN cd ${REPO_DIR_ABS} && go build -v -o ${BUILD_BIN_FULL_PATH}

FROM alpine:3.7 AS runtime

# container - meta
MAINTAINER <rosco_pecoltran@msn.com> Rosco Pecoltran

# container - deps
ARG APK_RUNTIME=${APK_RUNTIME:="ca-certificates libssh2 openssl"}

# app - env vars
ENV APP_OWNER=${APP_OWNER:-"sniperkit"} \
	APP_PROJECT=${APP_PROJECT:-"colly"} \
	APP_NAME=${APP_NAME:-"colly-dashboard"} \
	APP_WORKDIR=${APP_WORKDIR:-"/app"} \
	APP_SUFFIX_CMD=${APP_SUFFIX_CMD:-"tini -g --"}

# app - build args
ARG APP_EXECUTABLE_FILEPATH="${APP_EXECUTABLE_DIR}/${APP_NAME}"
ARG APP_EXECUTABLE_ARGS=${APP_EXECUTABLE_ARGS:-"-cfg_dir=${APP_WORKDIR}/config"}
ARG APP_DATADIR_DEFAULT="${APP_WORKDIR}/shared"
ARG APP_DATADIR=${APP_DATADIR:-"${APP_DATADIR_DEFAULT}"}
ARG APP_SCRIPT_ENTRYPOINT=${APP_SCRIPT_ENTRYPOINT:-"${APP_WORKDIR}/entrypoint.sh"}

# app - deps
ARG TINI_VERSION=${TINI_VERSION:-"v0.17.0"}
ARG TINI_ARCH=${TINI_ARCH:-"amd64"}
ARG TINI_VCS_URL=${TINI_VCS_URL:-"https://github.com/krallin/tini"}
ARG TINI_VCS_RELEASE_DOWNLOADL_URL="${TINI_REPO_URL}/releases/download/${TINI_VERSION}/tini-muslc-${TINI_ARCH}"
ARG TINI_PREFIX_PATH=${TINI_PREFIX_PATH:-"/usr/local/sbin"}
ARG TINI_BINDIR="${TINI_PREFIX_PATH}/tini"

# Install tini to /usr/local/sbin
ADD ${TINI_VCS_RELEASE_DOWNLOADL_URL} ${TINI_BINDIR}

# Install runtime dependencies & create runtime user
RUN apk --no-cache --no-progress add ${APK_RUNTIME} \
 	&& chmod +x ${TINI_BINDIR} && mkdir -p ${APP_WORKDIR} \
 	&& adduser -D ${APP_OWNER} -h ${APP_WORKDIR} -s /bin/sh \
 	&& su ${APP_OWNER} -c 'cd ${APP_WORKDIR}; mkdir -p bin config shared'

# Switch to user context
USER ${APP_OWNER}
WORKDIR ${APP_WORKDIR}

# app - env vars
ENV	APP_EXECUTABLE_DIR="${APP_WORKDIR}/bin" \
	APP_ENTRYPOINT_SCRIPT="${APP_WORKDIR}/entrypoint.sh" \
	APP_ENTRYPOINT_TARGET_DEFAULT=${APP_ENTRYPOINT_TARGET_DEFAULT:-"master"} \
	PATH="$PATH:${APP_EXECUTABLE_DIR}" 

# Copy app binary to /opt/colly/bin
COPY --from=builder ${REPO_DIR_ABS}/git2etcd ${APP_WORKDIR}/bin/${APP_NAME}

# Container configuration (create a gotemplate with vars)
# EXPOSE {{APP_PORTS}}
VOLUME ["${APP_DIR_SHARED}"]

# Start container command
ENTRYPOINT ["${APP_ENTRYPOINT_SCRIPT}"]
CMD ["${APP_ENTRYPOINT_TARGET_DEFAULT}"]

# Run direct entry point command
# ENTRYPOINT ["tini", "-g", "--"]
# CMD ["${APP_EXECUTABLE_FILEPATH}", "${APP_EXECUTABLE_ARGS}"]

# Snippets
# docker build -t sniperkit/colly-dashboard --no-cache .
# docker run -rm sniperkit/colly-dashboard -v `pwd`/shared:/app/shared -v `pwd`/conf:/app/conf
