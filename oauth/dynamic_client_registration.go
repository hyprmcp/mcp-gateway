package oauth

import (
	"crypto/rand"
	"encoding/json"
	"net/http"
	"slices"
	"strings"

	"github.com/dexidp/dex/api/v2"
	"github.com/hyprmcp/mcp-gateway/config"
	"github.com/hyprmcp/mcp-gateway/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

func NewDynamicClientRegistrationHandler(config *config.Config, meta map[string]any) (http.Handler, error) {
	clientTLSConfig, err := config.DexGRPCClient.ClientTLSConfig()
	if err != nil {
		return nil, err
	}

	var creds credentials.TransportCredentials

	if clientTLSConfig != nil {
		creds = credentials.NewTLS(clientTLSConfig)
	} else {
		creds = insecure.NewCredentials()
	}

	grpcClient, err := grpc.NewClient(
		config.DexGRPCClient.Addr,
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		return nil, err
	}

	dexClient := api.NewDexClient(grpcClient)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body ClientInformation
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		log.Get(r.Context()).Info("Received dynamic client registration request", "body", body)

		client := api.Client{
			Id:           genRandom(),
			Name:         body.ClientName,
			LogoUrl:      body.LogoURI,
			RedirectUris: body.RedirectURIs,
			Public:       true,
		}

		if !config.Authorization.GetDynamicClientRegistration().PublicClient {
			client.Secret = genRandom()
		}

		clientResponse, err := dexClient.CreateClient(r.Context(), &api.CreateClientReq{Client: &client})
		if err != nil {
			log.Get(r.Context()).Error(err, "failed to create client")
			http.Error(w, "Failed to create client", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")

		resp := ClientInformation{
			ClientID:     clientResponse.Client.Id,
			ClientSecret: clientResponse.Client.Secret,
			ClientName:   clientResponse.Client.Name,
			RedirectURIs: clientResponse.Client.RedirectUris,
			LogoURI:      clientResponse.Client.LogoUrl,
		}

		if scopesSupported := getSupportedScopes(meta); len(scopesSupported) > 0 {
			resp.Scope = strings.Join(scopesSupported, " ")
		}

		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Get(r.Context()).Error(err, "Failed to encode response")
		}

		log.Get(r.Context()).Info("Client created successfully", "client_id", clientResponse.Client.Id)
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
