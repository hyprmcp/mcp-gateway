# Hypr MCP Gateway Demo: "MCP, Who Am I?"

This demo showcases a complete and functional setup of the MCP Gateway including an instance of the Dex IdP used as
authorization server, as well as an upstream MCP server that returns information about the JWT that was used for
authentication.
The upstream MCP server is called "MCP, Who am I?" and can be found on GitHub:
[`hyprmcp/mcp-who-am-i`](http://github.com/hyprmcp/mcp-who-am-i/)

## Setup

All components of this setup are configured through a compose file and can be started with Docker Compose.
Please ensure that you have a working Docker installation.

### GitHub As Identity Provider

First, create a new OAuth application on GitHub here: https://github.com/settings/applications/new


| Field                      | Value                            |
|----------------------------|----------------------------------|
| Application name           | Hypr MCP Auth Demo               |
| Homepage URL               | https://hyprmcp.com              |
| Application description    | -                                |
| Authorization callback URL | `http://localhost:5556/callback` |
| Enable Device Flow         | `false` (unchecked)              |


After Creating the application make sure to `Generate a new client secret`.

You'll need the client ID and client secret for starting the server.

### Starting the server

Make sure to clone the repository locally:

```shell
git clone https://github.com/hyprmcp/mcp-gateway.git
````

Make sure to change into the who-am-i directory:

```shell
cd mcp-gateway/examples/who-am-i
````

Next, copy the file `.dex.secret.env.template` to `.dex.secret.env` and fill it with the client ID and client
secret of your new OAuth application.

Now you can start all components with Docker Compose:

```shell
docker compose up -d
```

## Testing

Use your favourite MCP client to connect to the MCP server at `http://localhost:9000/who-am-i/mcp`.
You can also use the MCP inspector tool by running `npx @modelcontextprotocol/inspector`.

You can either log in with your GitHub account or username password authentication with
`admin@example.com` and `password`.


If you want to bypass the authentication proxy you can directly call the "Who am I?" MCP server
at `http://localhost:3000/mcp` and will see that the request is not authenticated.

## HyprMCP Gateway

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


## Demo

Watch a demonstration of the MCP Gateway in action:

[![MCP Gateway Demo](https://img.youtube.com/vi/-oEzwJe1wac/maxresdefault.jpg)](https://youtu.be/-oEzwJe1wac)

## Hypr MCP Cloud

We also provide fully-managed MCP server and gateway hosting at Hypr MCP cloud, featuring
1-click OAuth authorization and [MCP analytics](https://hyprmcp.com/mcp-analytics/).

**Make sure to join our waitlist for early access:**

# <kbd>[**Join our waitlist**](https://hyprmcp.com/waitlist/)</kbd>
