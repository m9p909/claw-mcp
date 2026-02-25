# Builder stage
FROM golang:1.25-bookworm AS builder

# Install build dependencies
RUN apt-get update && apt-get install -y \
    git \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build

# Copy source code
COPY . .

# Build the binary
RUN go build -o /tmp/mcpclaw .

# Runtime stage
FROM ubuntu:22.04

# Install Node.js, Python, and batteries-included tools
RUN apt-get update && apt-get install -y \
    nodejs \
    npm \
    python3 \
    python3-pip \
    git \
    curl \
    wget \
    gawk \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create directories for workspace and data
RUN mkdir -p /.mcpclaw/workspace /.mcpclaw/data

# Copy compiled binary from builder
COPY --from=builder /tmp/mcpclaw /usr/local/bin/mcpclaw

# Copy entrypoint script
COPY scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Set working directory
WORKDIR /.mcpclaw

EXPOSE 8080

# Environment variable documentation:
# CLAW_TOKEN - Required. Bearer token for MCP endpoint authentication.
#              Must be set before container startup.
#              Example: docker run -e CLAW_TOKEN="your-secret-token" ...

ENTRYPOINT ["/entrypoint.sh"]
