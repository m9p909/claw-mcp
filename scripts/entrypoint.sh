#!/bin/bash
set -e

# Default port
PORT=${PORT:-8080}

# Ensure workspace and data directories exist
mkdir -p ~/.mcpclaw/workspace
mkdir -p ~/.mcpclaw/data

# Start the server
exec /usr/local/bin/mcpclaw -port "$PORT"
