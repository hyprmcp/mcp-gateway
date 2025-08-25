# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Development Commands

### Build and Run

```bash
# Run the server with the local config
mise run serve
```

### Development Setup
```bash
# Install development tools (requires mise)
mise install
```

## Architecture Overview

MCP-Gateway is an HTTP reverse proxy for MCP (Model Context Protocol) servers with OAuth authentication and webhook-based observability.

### Key Components

1. **HTTP Server** (`cmd/serve.go`): Cobra CLI with hot-reloading config via fsnotify
2. **OAuth Manager** (`oauth/`): JWT validation using JWK sets, integrates with Dex identity provider
3. **Proxy Handler** (`proxy/`):
   - Reverse proxy with MCP-aware transport
   - Intercepts JSON-RPC messages for session tracking
   - Supports SSE (Server-Sent Events) for streaming
4. **Webhook System** (`webhook/`): Async notifications with full request/response context


### Important Design Patterns
- **Config Hot-Reload**: Uses fsnotify to watch config file changes
- **Context Propagation**: Structured logging with request context
- **Transport Interception**: Custom RoundTripper for MCP protocol awareness
- **Per-Route Auth**: Authentication can be disabled for specific routes (e.g., public endpoints)
