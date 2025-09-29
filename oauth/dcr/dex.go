package dcr

import (
	"context"
	"crypto/rand"

	"github.com/dexidp/dex/api/v2"
	"github.com/hyprmcp/mcp-gateway/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type dexRegistrar struct {
	config    *config.AuthorizationDex
	dexClient api.DexClient
}

func NewDexRegistrar(config *config.AuthorizationDex) (ClientRegistrar, error) {
	grpcClient, err := grpc.NewClient(config.GRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	} else {
		return &dexRegistrar{config: config, dexClient: api.NewDexClient(grpcClient)}, nil
	}
}

// RegisterClient implements ClientRegistrar.
func (r *dexRegistrar) RegisterClient(ctx context.Context, request Client) (*Client, error) {
	client := api.Client{
		Id:           rand.Text(),
		Name:         request.ClientName,
		LogoUrl:      request.LogoURI,
		RedirectUris: request.RedirectURIs,
		Public:       true,
	}

	if r.config.DynamicClientRegistration.ClientType.IsConfidential() {
		client.Secret = rand.Text()
	}

	clientResponse, err := r.dexClient.CreateClient(ctx, &api.CreateClientReq{Client: &client})
	if err != nil {
		return nil, err
	}

	return &Client{
		ClientID:     clientResponse.Client.Id,
		ClientSecret: clientResponse.Client.Secret,
		ClientName:   clientResponse.Client.Name,
		RedirectURIs: clientResponse.Client.RedirectUris,
		LogoURI:      clientResponse.Client.LogoUrl,
	}, nil
}
