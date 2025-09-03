# Hypr MCP Gateway

[![Go Report Card](https://goreportcard.com/badge/github.com/hyprmcp/mcp-gateway)](https://goreportcard.com/report/github.com/hyprmcp/mcp-gateway)
[![GoDoc](https://godoc.org/github.com/hyprmcp/mcp-gateway?status.svg)](https://godoc.org/github.com/hyprmcp/mcp-gateway)

Hypr MCP Gateway featuring 1-click plug-in OAuth authorization including dynamic client registration
and MCP prompt analytics for streamable HTTP MCP servers.

## Main Features

- OAuth Proxy (incl. dynamic client registration)
- Prompt Telemetry
- MCP request logging and payload inspection

```
┌──────────────┐     OAuth2       ┌──────────────┐
│   End User   │ ───────────────▶ │   Hypr MCP   │
└──────────────┘  Single Login    │   Gateway    │
                                  └──────┬───────┘
                                         │
                               ┌─────────┼───────────┐
                               │         │           │
                          ┌────▼───┐ ┌───▼────┐ ┌────▼───┐
                          │ Tool A │ │ Tool B │ │ Tool C │
                          └────────┘ └────────┘ └────────┘
```

## Examples

Do you want to see the Hypr MCP gateway in action?

Check out our [`examples/who-am-i`](examples/who-am-i/README.md) featuring an upstream MCP server that
is able to return the authorization state.

## Why did we build Hypr MCP Gateway?

Adding OAuth2 support to an MCP server has lots of footguns and poses a significant challenge with to many developers.
By shifting the task to a separate infrastructure component, your MCP server can remain lean and be fully spec-compliant without you having to write a single line of code.
We have created a write-up on [_Building Supabase-like OAuth Authentication For MCP Servers_](https://hyprmcp.com/blog/mcp-server-authentication/)
on our [blog](https://hyprmcp.com/blog/) that goes into the details on MCP Server authentication. 


## Contributing & Local development

Check out our [CONTRIBUTING.md](./CONTRIBUTING.md) for detailed instructions.

## Hypr MCP Cloud

We also provide fully-managed MCP server and gateway hosting at Hypr MCP cloud, featuring
1-click OAuth authorization and [MCP analytics](https://hyprmcp.com/mcp-analytics/).

**Make sure to join our waitlist for early access:**

<kbd>[**Join our waitlist**](https://hyprmcp.com/waitlist/)</kbd>

