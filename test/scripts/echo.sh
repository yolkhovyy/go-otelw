#!/bin/bash

TRACE_ID=$(openssl rand -hex 16)
SPAN_ID=$(openssl rand -hex 8)
TRACE_FLAGS="01"
TRACEPARENT="00-${TRACE_ID}-${SPAN_ID}-${TRACE_FLAGS}"

BODY="${1:-echo}"
COUNT="${2:-5}"
URL="http://localhost:8080/echo?count=${COUNT}"

echo "${BODY} with traceparent: ${TRACEPARENT}"

curl "${URL}" \
  -H "traceparent: ${TRACEPARENT}" \
  -d "${BODY}"
