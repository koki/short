#!/usr/bin/env bash
docker tag short-server ${KOKI_SHORT_SERVER_IMAGE}
docker push ${KOKI_SHORT_SERVER_IMAGE}
