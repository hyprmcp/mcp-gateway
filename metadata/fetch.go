package metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"go.uber.org/multierr"
)

func FetchMetadataDiscovery(ctx context.Context, server string) (Metadata, error) {
	uris, err := getMetadataURIs(server)
	if err != nil {
		return nil, err
	}

	for _, u := range uris {
		if metadata, err1 := FetchMetadata(ctx, u); err1 != nil {
			multierr.AppendInto(&err, err1)
		} else {
			return metadata, nil
		}
	}

	return nil, err
}

func FetchMetadata(ctx context.Context, url string) (Metadata, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("failed to fetch %v: %v", url, resp.Status)
	}

	var metadata map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

func getMetadataURIs(server string) ([]string, error) {
	var uris []string

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, fmt.Errorf("failed to parse authorization server URL: %w", err)
	}

	serverURL.Path = OIDCMetadataPath
	uris = append(uris, serverURL.String())
	serverURL.Path = OAuth2MetadataPath
	uris = append(uris, serverURL.String())
	return uris, nil
}
