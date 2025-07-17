package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func Send(ctx context.Context, method string, url string, payload WebhookPayload) error {
	if method == "" {
		method = http.MethodPost
	}
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return err
	}

	if req, err := http.NewRequestWithContext(ctx, method, url, &buf); err != nil {
		return err
	} else if resp, err := http.DefaultClient.Do(req); err != nil {
		return err
	} else if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unexpected http status: %v", resp.Status)
	} else {
		return nil
	}
}
