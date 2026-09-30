#!/bin/bash -e

DIR=$(cd "$(dirname "$0")" && pwd)
cd "$DIR"

ARTIFACT_SUBDIR=$1
SPEC=$2

if [[ -z "$SPEC" ]]; then
    echo "usage $0 artifact_subdir spec"
    exit 1
fi

export PLAYWRIGHT_FULL_DOMAIN=${PLAYWRIGHT_FULL_DOMAIN:-bookworm.com}
export PLAYWRIGHT_APP_DOMAIN=${PLAYWRIGHT_APP_DOMAIN:-jellyfin.${PLAYWRIGHT_FULL_DOMAIN}}
export PLAYWRIGHT_DEVICE_HOST=${PLAYWRIGHT_DEVICE_HOST:-${PLAYWRIGHT_APP_DOMAIN}}
export PLAYWRIGHT_DEVICE_USER=${PLAYWRIGHT_DEVICE_USER:-user}
export PLAYWRIGHT_DEVICE_PASSWORD=${PLAYWRIGHT_DEVICE_PASSWORD:-Password1}
export PLAYWRIGHT_SSH_USER=${PLAYWRIGHT_SSH_USER:-root}
export PLAYWRIGHT_SSH_PASSWORD=${PLAYWRIGHT_SSH_PASSWORD:-Password1}
export PLAYWRIGHT_ARTIFACT_DIR=${PLAYWRIGHT_ARTIFACT_DIR:-${DIR}/../../artifact/${ARTIFACT_SUBDIR}}

${DIR}/../../apt.sh sshpass openssh-client curl
${DIR}/wait-app.sh ${PLAYWRIGHT_APP_DOMAIN}
npm ci --no-audit --no-fund

for project in desktop mobile; do
  PLAYWRIGHT_PROJECT=${project} npx playwright test --project=${project} "$SPEC"
done
