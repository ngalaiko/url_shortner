#!/usr/bin/env sh

STATIC_DIR=/app/static
SHARED_DIR=/data/shortner

BIN_PATH=/app/url_shortner
CONFIG_PATH=/app/config.yaml

mv ${STATIC_DIR} ${SHARED_DIR}

${BIN_PATH} --config=${CONFIG_PATH}
