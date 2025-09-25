package oauth

import (
	"encoding/json"
	"net/http"

	"github.com/hyprmcp/mcp-gateway/config"
	"github.com/hyprmcp/mcp-gateway/log"
	"github.com/hyprmcp/mcp-gateway/metadata"
)

func NewAuthorizationServerMetadataHandler(config *config.Config, metaSource metadata.MetadataSource, fn metadata.EditMetadataFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		meta, err := metaSource.GetMetadata(r.Context())
		if err != nil {
			log.Get(r.Context()).Error(err, "failed to get authorization server metadata from upstream")
			http.Error(w, "Failed to retrieve authorization server metadata", http.StatusInternalServerError)
		}

		if err := fn(meta); err != nil {
			log.Get(r.Context()).Error(err, "failed to edit authorization server metadata")
			http.Error(w, "Failed to edit authorization server metadata", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(meta); err != nil {
			log.Get(r.Context()).Error(err, "failed to encode authorization server metadata response")
		}
	})
}
