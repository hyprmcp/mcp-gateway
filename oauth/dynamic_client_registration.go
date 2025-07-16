package oauth

import (
	"crypto/rand"
	"encoding/json"
	"net/http"

	"github.com/dexidp/dex/api/v2"
	"github.com/jetski-sh/mcp-proxy/config"
	"github.com/jetski-sh/mcp-proxy/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const DynamicClientRegistrationPath = "/oauth/register"

type ClientInformation struct {
	ClientID              string   `json:"client_id"`
	ClientSecret          string   `json:"client_secret"`
	ClientSecretExpiresAt int64    `json:"client_secret_expires_at"`
	ClientName            string   `json:"client_name,omitempty"`
	RedirectURIs          []string `json:"redirect_uris"`
	LogoURI               string   `json:"logo_uri,omitempty"`
}

func NewDynamicClientRegistrationHandler(config *config.Config) (http.Handler, error) {
	grpcClient, err := grpc.NewClient(
		config.DexGRPCClient.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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
			Secret:       genRandom(),
			Name:         body.ClientName,
			LogoUrl:      body.LogoURI,
			RedirectUris: body.RedirectURIs,
		}

		clientResponse, err := dexClient.CreateClient(r.Context(), &api.CreateClientReq{Client: &client})
		if err != nil {
			http.Error(w, "Failed to create client", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(ClientInformation{
			ClientID:     clientResponse.Client.Id,
			ClientSecret: clientResponse.Client.Secret,
			ClientName:   clientResponse.Client.Name,
			RedirectURIs: clientResponse.Client.RedirectUris,
			LogoURI:      clientResponse.Client.LogoUrl,
		})
		if err != nil {
			log.Get(r.Context()).Error(err, "Failed to encode response")
		}

		log.Get(r.Context()).Info("Client created successfully", "client_id", clientResponse.Client.Id)
	}), nil
}

func genRandom() string {
	return rand.Text()
}
