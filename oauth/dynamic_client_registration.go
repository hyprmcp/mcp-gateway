package oauth

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/dexidp/dex/api/v2"
	"github.com/hyprmcp/mcp-gateway/config"
	"github.com/hyprmcp/mcp-gateway/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const DynamicClientRegistrationPath = "/oauth/register"

type ClientInformation struct {
	ClientID              string   `json:"client_id"`
	ClientSecret          string   `json:"client_secret,omitempty"`
	ClientSecretExpiresAt int64    `json:"client_secret_expires_at,omitempty"`
	ClientName            string   `json:"client_name,omitempty"`
	RedirectURIs          []string `json:"redirect_uris"`
	LogoURI               string   `json:"logo_uri,omitempty"`
	Scope                 string   `json:"scope,omitempty"`
}

type DynamicClientRegistrationHandler interface {
	Handle(ctx context.Context, request ClientInformation) (*ClientInformation, error)
}

type dexDCRHandler struct {
	config    *config.Config
	dexClient api.DexClient
}

func (h *dexDCRHandler) Handle(ctx context.Context, request ClientInformation) (*ClientInformation, error) {
	client := api.Client{
		Id:           genRandom(),
		Name:         request.ClientName,
		LogoUrl:      request.LogoURI,
		RedirectUris: request.RedirectURIs,
		Public:       true,
	}

	if !h.config.Authorization.GetDynamicClientRegistration().PublicClient {
		client.Secret = genRandom()
	}

	clientResponse, err := h.dexClient.CreateClient(ctx, &api.CreateClientReq{Client: &client})
	if err != nil {
		return nil, err
	}

	return &ClientInformation{
		ClientID:     clientResponse.Client.Id,
		ClientSecret: clientResponse.Client.Secret,
		ClientName:   clientResponse.Client.Name,
		RedirectURIs: clientResponse.Client.RedirectUris,
		LogoURI:      clientResponse.Client.LogoUrl,
	}, nil
}

type fakeDCRHandler struct {
	clientId string
}

func (h *fakeDCRHandler) Handle(ctx context.Context, request ClientInformation) (*ClientInformation, error) {
	request.ClientID = h.clientId
	return &request, nil
}

func NewDynamicClientRegistrationHandler(config *config.Config, meta map[string]any) (http.Handler, error) {
	var dcrHandler DynamicClientRegistrationHandler
	if config.DexGRPCClient != nil {
		grpcClient, err := grpc.NewClient(
			config.DexGRPCClient.Addr,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, err
		}

		dcrHandler = &dexDCRHandler{
			config:    config,
			dexClient: api.NewDexClient(grpcClient),
		}
	} else if config.Authorization.ClientID != "" {
		dcrHandler = &fakeDCRHandler{
			clientId: config.Authorization.ClientID,
		}
	} else {
		return nil, errors.New("incomplete DCR config")
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body ClientInformation
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		log.Get(r.Context()).Info("Received dynamic client registration request", "body", body)

		resp, err := dcrHandler.Handle(r.Context(), body)
		if err != nil {
			log.Get(r.Context()).Error(err, "failed to create client")
			http.Error(w, "Failed to create client", http.StatusInternalServerError)
			return
		}

		if scopesSupported := getSupportedScopes(meta); len(scopesSupported) > 0 {
			resp.Scope = strings.Join(scopesSupported, " ")
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Get(r.Context()).Error(err, "Failed to encode response")
		}

		log.Get(r.Context()).Info("Client created successfully", "client_id", resp.ClientID)
	}), nil
}

func genRandom() string {
	return rand.Text()
}

func getSupportedScopes(meta map[string]any) []string {
	if scopesSupported, ok := meta["scopes_supported"].([]any); ok {
		scopesSupportedStr := make([]string, 0, len(scopesSupported))
		for _, v := range scopesSupported {
			if s, ok := v.(string); ok {
				scopesSupportedStr = append(scopesSupportedStr, s)
			}
		}

		return slices.Clip(scopesSupportedStr)
	}

	return nil
}
