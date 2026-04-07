#!/usr/bin/env bash
set -e

NODE_IP="${1:?Usage: $0 <NODE_IP> [QUANTITY]}"
QUANTITY="${2:-100}"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
GENERATOR="$SCRIPT_DIR/../go/generator/generator"

if [ ! -f "$GENERATOR" ]; then
    echo "Building generator..."
    cd "$SCRIPT_DIR/../go"
    go build -o generator/generator ./generator/
    cd "$SCRIPT_DIR"
fi

echo "Sending $QUANTITY UDP packets to $NODE_IP..."
"$GENERATOR" Do Udp_port=30000 Destination="$NODE_IP" Port=30000 Quantity="$QUANTITY"
