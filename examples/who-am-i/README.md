# Demo: "MCP, Who Am I?"

This demo showcases a completele functional setup of the MCP Gateway including an instance of the Dex IdP used as
authorization server, as well as an upstream MCP server that returns information about the JWT that was used for
authentication.

## Setup

All components of this setup are configured through a compose file and can be started with Docker Compose.
Please ensure that you have a working Docker installation.

### GitHub As Federated OIDC Provider

First, create a new OAuth application on GitHub here: https://github.com/settings/applications/new
Use `http://localhost:5556/callback` as redirect URI.
Next, copy the file `.dex.secret.env.template` to `.dex.secret.env` and fill it with the client ID and secret of your new OAuth application.

### Starting the server

You can start all components with Docker Compose:

```shell
docker compose up -d
```

## Testing

Use your favourite MCP client to connect to the MCP server at `http://localhost:9000/who-am-i/mcp`.
You can also use the MCP inspector tool by runnign `npx @modelcontextprotocol/inspector`.
