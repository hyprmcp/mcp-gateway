package dcr

type Client struct {
	ClientID              string   `json:"client_id"`
	ClientSecret          string   `json:"client_secret,omitempty"`
	ClientSecretExpiresAt int64    `json:"client_secret_expires_at,omitempty"`
	ClientName            string   `json:"client_name,omitempty"`
	RedirectURIs          []string `json:"redirect_uris"`
	LogoURI               string   `json:"logo_uri,omitempty"`
	Scope                 string   `json:"scope,omitempty"`
}
