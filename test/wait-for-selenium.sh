#!/bin/sh -e

URL=${1:-http://selenium:4444/wd/hub/status}
ATTEMPTS=${2:-120}

i=0
while [ $i -lt $ATTEMPTS ]; do
    if python3 - "$URL" <<'EOF' 2>/dev/null
import json, sys, urllib.request

with urllib.request.urlopen(sys.argv[1], timeout=5) as response:
    sys.exit(0 if json.load(response).get("value", {}).get("ready") else 1)
EOF
    then
        echo "selenium is ready"
        exit 0
    fi
    i=$((i + 1))
    sleep 5
done

echo "selenium did not become ready after $ATTEMPTS attempts" >&2
exit 1
