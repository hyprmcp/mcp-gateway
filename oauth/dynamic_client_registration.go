package oauth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hyprmcp/mcp-gateway/log"
	"github.com/hyprmcp/mcp-gateway/metadata"
	"github.com/hyprmcp/mcp-gateway/oauth/dcr"
)

const DynamicClientRegistrationPath = "/oauth/register"

func NewDynamicClientRegistrationHandler(registrar dcr.ClientRegistrar, meta metadata.Metadata) (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var clientInformation dcr.Client
		if err := json.NewDecoder(r.Body).Decode(&clientInformation); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		log.Get(r.Context()).Info("Received dynamic client registration request", "body", clientInformation)

		resp, err := registrar.RegisterClient(r.Context(), clientInformation)
		if err != nil {
			log.Get(r.Context()).Error(err, "failed to create client")
			http.Error(w, "Failed to create client", http.StatusInternalServerError)
			return
		}

		if scopesSupported := meta.GetSupportedScopes(); len(scopesSupported) > 0 {
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
