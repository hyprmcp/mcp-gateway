package dcr

import "context"

type fakeRegistrar struct {
	clientId string
}

func NewFakeRegistrar(src ClientIDSource) ClientRegistrar {
	return &fakeRegistrar{clientId: src.GetClientID()}
}

// RegisterClient implements ClientRegistrar.
func (r *fakeRegistrar) RegisterClient(ctx context.Context, client Client) (*Client, error) {
	client.ClientID = r.clientId
	return &client, nil
}
